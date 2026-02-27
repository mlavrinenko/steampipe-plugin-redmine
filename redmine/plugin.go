package redmine

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name: "steampipe-plugin-redmine",
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
		},
		DefaultTransform:   transform.FromGo().NullIfZero(),
		DefaultRetryConfig: retryConfig(),
		TableMap: map[string]*plugin.Table{
			"redmine_issue_journal": tableRedmineIssueJournal(),
		},
	}
	return p
}
