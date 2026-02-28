package redmine

import (
	"os"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/schema"
)

type redmineConfig struct {
	Endpoint *string `hcl:"endpoint"`
	APIKey   *string `hcl:"api_key"`
}

var ConfigSchema = map[string]*schema.Attribute{
	"endpoint": {Type: schema.TypeString},
	"api_key":  {Type: schema.TypeString},
}

func ConfigInstance() interface{} {
	return &redmineConfig{}
}

func GetConfig(connection *plugin.Connection) *redmineConfig {
	if connection == nil || connection.Config == nil {
		return &redmineConfig{}
	}
	config, ok := connection.Config.(*redmineConfig)
	if !ok {
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
