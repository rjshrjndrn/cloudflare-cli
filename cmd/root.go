package cmd

import (
	"fmt"
	"os"

	"github.com/skyline/cfcli/internal/config"
	"github.com/spf13/cobra"
)

var (
	cfgFile    string
	email      string
	token      string
	account    string
	domain     string
	recordType string
	newType    string
	priority   int64
	ttl        int64
	activate   bool
	format     string
	query      string

	cfg *config.Config
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "cfcli",
	Short: "Cloudflare CLI - Manage your Cloudflare DNS records",
	Long: `cfcli is a command-line interface for managing Cloudflare DNS records.
It supports CRUD operations on DNS records with a simple, intuitive syntax.`,
	Version: version,
}

func SetVersion(v, c, d string) {
	version = v
	commit = c
	date = d
	rootCmd.Version = fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date)
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.config/cfcli/config.yaml)")
	rootCmd.PersistentFlags().StringVarP(&email, "email", "e", "", "Email of your cloudflare account")
	rootCmd.PersistentFlags().StringVarP(&token, "token", "k", "", "API token for your cloudflare account")
	rootCmd.PersistentFlags().StringVarP(&account, "account", "u", "", "Named account from config file")
	rootCmd.PersistentFlags().StringVarP(&domain, "domain", "d", "", "Domain to operate on")
	rootCmd.PersistentFlags().StringVarP(&recordType, "type", "t", "", "Type of DNS record (A, AAAA, CNAME, MX, TXT, NS, SRV)")
	rootCmd.PersistentFlags().StringVarP(&newType, "newtype", "n", "", "New type when editing a record")
	rootCmd.PersistentFlags().Int64VarP(&priority, "priority", "p", 0, "Priority for MX or SRV records")
	rootCmd.PersistentFlags().Int64VarP(&ttl, "ttl", "l", 1, "TTL in seconds (1 for auto, 120-86400)")
	rootCmd.PersistentFlags().BoolVarP(&activate, "activate", "a", false, "Activate cloudflare (enable proxy) after creating record")
	rootCmd.PersistentFlags().StringVarP(&format, "format", "f", "table", "Output format: table, json, csv")
	rootCmd.PersistentFlags().StringVarP(&query, "query", "q", "", "Comma-separated filters (e.g., content:1.1.1.1,type:A)")
}

func initConfig() {
	var err error
	cfg, err = config.LoadConfig(cfgFile, account)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Could not load config file: %v\n", err)
		cfg = &config.Config{}
	}

	// Override config with command-line flags
	if token != "" {
		cfg.Token = token
	}
	if email != "" {
		cfg.Email = email
	}
	if domain != "" {
		cfg.Domain = domain
	}
}
