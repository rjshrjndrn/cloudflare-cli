package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Token   string
	Email   string
	Domain  string
	Account string
}

type AccountConfig struct {
	Token  string `mapstructure:"token"`
	Email  string `mapstructure:"email"`
	Domain string `mapstructure:"domain"`
}

type ConfigFile struct {
	Defaults struct {
		Token   string `mapstructure:"token"`
		Email   string `mapstructure:"email"`
		Domain  string `mapstructure:"domain"`
		Account string `mapstructure:"account"`
	} `mapstructure:"defaults"`
	Accounts map[string]AccountConfig `mapstructure:"accounts"`
}

func LoadConfig(configPath, accountName string) (*Config, error) {
	v := viper.New()

	// Set config file path
	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		v.AddConfigPath(home)
		v.SetConfigName(".cfcli")
		v.SetConfigType("yml")
	}

	// Read config file (optional)
	var cfg ConfigFile
	if err := v.ReadInConfig(); err == nil {
		if err := v.Unmarshal(&cfg); err != nil {
			return nil, err
		}
	}

	// Build final config with precedence: CLI flags > Env vars > Config file
	config := &Config{}

	// Check environment variables first
	if token := os.Getenv("CF_API_KEY"); token != "" {
		config.Token = token
	}
	if email := os.Getenv("CF_API_EMAIL"); email != "" {
		config.Email = email
	}
	if domain := os.Getenv("CF_API_DOMAIN"); domain != "" {
		config.Domain = domain
	}

	// If no env vars, use config file
	if config.Token == "" {
		if accountName != "" {
			if acc, ok := cfg.Accounts[accountName]; ok {
				config.Token = acc.Token
				config.Email = acc.Email
				config.Domain = acc.Domain
			}
		} else if cfg.Defaults.Account != "" {
			// Use default account
			if acc, ok := cfg.Accounts[cfg.Defaults.Account]; ok {
				config.Token = acc.Token
				config.Email = acc.Email
				config.Domain = acc.Domain
			}
		} else {
			// Use defaults directly
			config.Token = cfg.Defaults.Token
			config.Email = cfg.Defaults.Email
			config.Domain = cfg.Defaults.Domain
		}
	}

	return config, nil
}

func GetDefaultConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".cfcli.yml")
}
