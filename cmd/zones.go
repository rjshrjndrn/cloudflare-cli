package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/skyline/cfcli/internal/cloudflare"
	"github.com/spf13/cobra"
)

var zonesCmd = &cobra.Command{
	Use:   "zones",
	Short: "List all zones in your Cloudflare account",
	RunE: func(cmd *cobra.Command, args []string) error {
		if cfg.Token == "" {
			return fmt.Errorf("API token is required (use -k or set CF_API_KEY)")
		}

		client, err := cloudflare.NewClient(cfg.Token, cfg.Email)
		if err != nil {
			return err
		}
		ctx := context.Background()

		zones, err := client.ListZones(ctx)
		if err != nil {
			return err
		}

		if len(zones) == 0 {
			fmt.Println("No zones found")
			return nil
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.Header("Name", "Status", "ID")

		for _, zone := range zones {
			table.Append(zone.Name, string(zone.Status), zone.ID)
		}

		table.Render()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(zonesCmd)
}
