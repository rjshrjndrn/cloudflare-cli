package cmd

import (
	"context"
	"fmt"

	"github.com/skyline/cfcli/internal/cloudflare"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:     "rm <name> [content]",
	Aliases: []string{"remove", "delete"},
	Short:   "Remove DNS record(s)",
	Long: `Remove one or more DNS records matching the criteria.
	
Examples:
  cfcli -d example.com rm test                           # Remove all records named 'test'
  cfcli -d example.com -t A rm test                      # Remove A records named 'test'
  cfcli -d example.com -t A rm test -q content:1.1.1.1   # Remove specific record`,
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

		// Parse query filters
		filters := parseQuery(query)
		if content != "" {
			filters["content"] = content
		}

		client, err := cloudflare.NewClient(cfg.Token, cfg.Email)
		if err != nil {
			return err
		}
		ctx := context.Background()

		if err := client.SetZone(ctx, cfg.Domain); err != nil {
			return err
		}

		// Find matching records
		queryContent := filters["content"]
		queryType := recordType
		if filters["type"] != "" {
			queryType = filters["type"]
		}

		records, err := client.FindDNSRecord(ctx, name, queryContent, queryType)
		if err != nil {
			return err
		}

		if len(records) == 0 {
			return fmt.Errorf("no records found matching the criteria")
		}

		// Apply additional filters
		if len(filters) > 0 {
			records = filterRecords(records)
		}

		if len(records) == 0 {
			return fmt.Errorf("no records found matching all filters")
		}

		// Delete all matching records
		for _, record := range records {
			if err := client.DeleteDNSRecord(ctx, record.ID); err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "Failed to delete %s record %s: %v\n", record.Type, record.Name, err)
				continue
			}
			fmt.Printf("✓ Deleted %s record: %s -> %s\n", record.Type, record.Name, record.Content)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
