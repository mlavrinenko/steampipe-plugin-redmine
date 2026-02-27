package redmine

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

type redmineConfig struct {
	Endpoint *string `hcl:"endpoint"`
	APIKey   *string `hcl:"api_key"`
}

func ConfigInstance() interface{} {
	return &redmineConfig{}
}

func GetConfig(connection *plugin.Connection) redmineConfig {
	if connection == nil || connection.Config == nil {
		return redmineConfig{}
	}
	config, _ := connection.Config.(redmineConfig)
	return config
}
