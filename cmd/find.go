package cmd

import (
	"context"
	"fmt"

	"github.com/skyline/cfcli/internal/cloudflare"
	"github.com/spf13/cobra"
)

var findCmd = &cobra.Command{
	Use:   "find <name> [content]",
	Short: "Find DNS record(s)",
	Long: `Find DNS records matching the specified criteria.
	
Examples:
  cfcli -d example.com find test
  cfcli -d example.com -t A find test
  cfcli -d example.com find -q content:1.1.1.1`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if cfg.Token == "" {
			return fmt.Errorf("API token is required (use -k or set CF_API_KEY)")
		}
		if cfg.Domain == "" {
			return fmt.Errorf("domain is required (use -d or set CF_API_DOMAIN)")
		}

		name := args[0]
		content := ""
		if len(args) > 1 {
			content = args[1]
		}

		client, err := cloudflare.NewClient(cfg.Token, cfg.Email)
		if err != nil {
			return err
		}
		ctx := context.Background()

		if err := client.SetZone(ctx, cfg.Domain); err != nil {
			return err
		}

		records, err := client.FindDNSRecord(ctx, name, content, recordType)
		if err != nil {
			return err
		}

		// Apply query filters if provided
		if query != "" {
			records = filterRecords(records)
		}

		if len(records) == 0 {
			fmt.Println("No records found")
			return nil
		}

		fmt.Printf("Found %d record(s):\n\n", len(records))
		return outputTable(records)
	},
}

func init() {
	rootCmd.AddCommand(findCmd)
}
