package cmd

import (
	"context"
	"fmt"

	"github.com/skyline/cfcli/internal/cloudflare"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit <name> <content>",
	Short: "Edit a DNS record",
	Long: `Edit an existing DNS record.
	
Examples:
  cfcli -d example.com -t A edit mail 5.6.7.8
  cfcli -d example.com -t A -n CNAME edit test example.com  # Change type
  cfcli -d example.com -t A --ttl 300 edit mail 1.2.3.4     # Set TTL`,
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

		// Determine the actual type to update to
		updateType := recordType
		if newType != "" {
			updateType = newType
		}

		client, err := cloudflare.NewClient(cfg.Token, cfg.Email)
		if err != nil {
			return err
		}
		ctx := context.Background()

		if err := client.SetZone(ctx, cfg.Domain); err != nil {
			return err
		}

		// Find the record to update
		records, err := client.FindDNSRecord(ctx, name, "", recordType)
		if err != nil {
			return err
		}

		if len(records) == 0 {
			return fmt.Errorf("no record found matching name=%s type=%s", name, recordType)
		}

		if len(records) > 1 {
			return fmt.Errorf("multiple records found (%d), please be more specific with -q", len(records))
		}

		record := records[0]

		var priorityPtr *uint16
		if priority > 0 {
			p := uint16(priority)
			priorityPtr = &p
		} else if record.Priority != nil {
			priorityPtr = record.Priority
		}

		// Use existing proxy setting if not specified
		proxied := false
		if record.Proxied != nil {
			proxied = *record.Proxied
		}
		if activate {
			proxied = true
		}

		updated, err := client.UpdateDNSRecord(ctx, record.ID, updateType, name, content, int(ttl), priorityPtr, proxied)
		if err != nil {
			return err
		}

		fmt.Printf("✓ Updated %s record: %s -> %s", updated.Type, updated.Name, updated.Content)
		if updated.Proxied != nil && *updated.Proxied {
			fmt.Print(" (proxied)")
		}
		fmt.Println()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
