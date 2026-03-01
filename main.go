package main

import (
	"github.com/mlavrinenko/steampipe-plugin-redmine/redmine"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{PluginFunc: redmine.Plugin})
}
