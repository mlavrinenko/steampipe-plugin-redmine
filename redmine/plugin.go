// Package redmine implements a Steampipe plugin for querying Redmine instances.
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
		DefaultIgnoreConfig: &plugin.IgnoreConfig{
			ShouldIgnoreError: isNotFoundError([]string{"404", "not found"}),
		},
		TableMap: map[string]*plugin.Table{
			"redmine_issue":         tableRedmineIssue(),
			"redmine_issue_journal": tableRedmineIssueJournal(),
			"redmine_project":       tableRedmineProject(),
			"redmine_user":          tableRedmineUser(),
		},
	}
	return p
}
