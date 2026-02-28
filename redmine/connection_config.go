package redmine

import (
	"os"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

type redmineConfig struct {
	Endpoint *string `hcl:"endpoint"`
	APIKey   *string `hcl:"api_key"`
}

func ConfigInstance() interface{} {
	return &redmineConfig{}
}

func GetConfig(connection *plugin.Connection) *redmineConfig {
	if connection == nil || connection.Config == nil {
		return &redmineConfig{}
	}

	// The SDK may store config as either *redmineConfig or redmineConfig (by value)
	// depending on the parsing path used.
	var config *redmineConfig
	switch c := connection.Config.(type) {
	case *redmineConfig:
		config = c
	case redmineConfig:
		config = &c
	default:
		return &redmineConfig{}
	}

	// Fall back to environment variables
	if config.Endpoint == nil || *config.Endpoint == "" {
		if v := os.Getenv("REDMINE_ENDPOINT"); v != "" {
			config.Endpoint = &v
		}
	}
	if config.APIKey == nil || *config.APIKey == "" {
		if v := os.Getenv("REDMINE_API_KEY"); v != "" {
			config.APIKey = &v
		}
	}

	return config
}
