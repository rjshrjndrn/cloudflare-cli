package cmd

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/skyline/cfcli/internal/cloudflare"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "listrecords"},
	Short:   "List DNS records for the domain",
	RunE: func(cmd *cobra.Command, args []string) error {
		if cfg.Token == "" {
			return fmt.Errorf("API token is required (use -k or set CF_API_KEY)")
		}
		if cfg.Domain == "" {
			return fmt.Errorf("domain is required (use -d or set CF_API_DOMAIN)")
		}

		client, err := cloudflare.NewClient(cfg.Token, cfg.Email)
		if err != nil {
			return err
		}
		ctx := context.Background()

		if err := client.SetZone(ctx, cfg.Domain); err != nil {
			return err
		}

		records, err := client.ListDNSRecords(ctx)
		if err != nil {
			return err
		}

		// Filter records if query is provided
		if query != "" {
			records = filterRecords(records)
		}

		// Output in requested format
		switch strings.ToLower(format) {
		case "json":
			return outputJSON(records)
		case "csv":
			return outputCSV(records)
		default:
			return outputTable(records)
		}
	},
}

func filterRecords(records []cloudflare.DNSRecord) []cloudflare.DNSRecord {
	filters := parseQuery(query)
	var filtered []cloudflare.DNSRecord

	for _, record := range records {
		match := true
		for key, value := range filters {
			switch strings.ToLower(key) {
			case "type":
				if !strings.EqualFold(record.Type, value) {
					match = false
				}
			case "name":
				if !strings.Contains(strings.ToLower(record.Name), strings.ToLower(value)) {
					match = false
				}
			case "content":
				if !strings.Contains(strings.ToLower(record.Content), strings.ToLower(value)) {
					match = false
				}
			}
		}
		if match {
			filtered = append(filtered, record)
		}
	}
	return filtered
}

func parseQuery(q string) map[string]string {
	filters := make(map[string]string)
	if q == "" {
		return filters
	}

	parts := strings.Split(q, ",")
	for _, part := range parts {
		kv := strings.SplitN(part, ":", 2)
		if len(kv) == 2 {
			filters[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
		}
	}
	return filters
}

func outputJSON(records []cloudflare.DNSRecord) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(records)
}

func outputCSV(records []cloudflare.DNSRecord) error {
	writer := csv.NewWriter(os.Stdout)
	defer writer.Flush()

	// Write header
	if err := writer.Write([]string{"ID", "Type", "Name", "Content", "TTL", "Priority", "Proxied"}); err != nil {
		return err
	}

	// Write records
	for _, record := range records {
		priority := ""
		if record.Priority != nil {
			priority = strconv.Itoa(int(*record.Priority))
		}
		proxiedStr := "false"
		if record.Proxied != nil && *record.Proxied {
			proxiedStr = "true"
		}
		row := []string{
			record.ID,
			record.Type,
			record.Name,
			record.Content,
			strconv.Itoa(record.TTL),
			priority,
			proxiedStr,
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}
	return nil
}

func outputTable(records []cloudflare.DNSRecord) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"Type", "Name", "Content", "TTL", "Priority", "Proxied"})

	for _, record := range records {
		priority := "-"
		if record.Priority != nil {
			priority = strconv.Itoa(int(*record.Priority))
		}
		ttlStr := strconv.Itoa(record.TTL)
		if record.TTL == 1 {
			ttlStr = "auto"
		}
		proxied := " "
		if record.Proxied != nil && *record.Proxied {
			proxied = "✓"
		}
		if err := table.Append(record.Type, record.Name, record.Content, ttlStr, priority, proxied); err != nil {
			return err
		}
	}
	return table.Render()
}

func init() {
	rootCmd.AddCommand(listCmd)
}
