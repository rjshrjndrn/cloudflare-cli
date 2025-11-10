package cmd

import (
	"context"
	"fmt"

	"github.com/skyline/cfcli/internal/cloudflare"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <name> <content>",
	Short: "Add a DNS record",
	Long: `Add a new DNS record to your domain.
	
Examples:
  cfcli -d example.com -t A add mail 1.2.3.4
  cfcli -d example.com -t CNAME add www example.com
  cfcli -d example.com -t MX -p 10 add @ mail.example.com
  cfcli -d example.com -t A -a add test 1.1.1.1  # -a enables proxy`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if cfg.Token == "" {
			return fmt.Errorf("API token is required (use -k or set CF_API_KEY)")
		}
		if cfg.Domain == "" {
			return fmt.Errorf("domain is required (use -d or set CF_API_DOMAIN)")
		}
		if recordType == "" {
			return fmt.Errorf("record type is required (use -t)")
		}

		name := args[0]
		content := args[1]

		client, err := cloudflare.NewClient(cfg.Token, cfg.Email)
		if err != nil {
			return err
		}
		ctx := context.Background()

		if err := client.SetZone(ctx, cfg.Domain); err != nil {
			return err
		}

		var priorityPtr *uint16
		if priority > 0 {
			p := uint16(priority)
			priorityPtr = &p
		}

		record, err := client.AddDNSRecord(ctx, recordType, name, content, int(ttl), priorityPtr, activate)
		if err != nil {
			return err
		}

		fmt.Printf("✓ Created %s record: %s -> %s", record.Type, record.Name, record.Content)
		if record.Proxied != nil && *record.Proxied {
			fmt.Print(" (proxied)")
		}
		fmt.Println()
		fmt.Printf("  Record ID: %s\n", record.ID)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
