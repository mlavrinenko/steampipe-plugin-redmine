// Package redmine implements a Steampipe plugin for querying Redmine instances.
package redmine

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-sdk/v5/rate_limiter"
)

func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name: "steampipe-plugin-redmine",
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
			Schema:      ConfigSchema,
		},
		ConnectionConfigChangedFunc: configChanged,
		DefaultTransform:            transform.FromGo().NullIfZero(),
		DefaultRetryConfig:          retryConfig(),
		RateLimiters: []*rate_limiter.Definition{
			{
				Name:       "redmine_global",
				FillRate:   10,
				BucketSize: 10,
				Scope:      []string{"connection"},
			},
		},
		DefaultIgnoreConfig: &plugin.IgnoreConfig{
			ShouldIgnoreError: isNotFoundError([]string{"404", "not found"}),
		},
		TableMap: map[string]*plugin.Table{
			"redmine_issue":          tableRedmineIssue(),
			"redmine_issue_journal":  tableRedmineIssueJournal(),
			"redmine_issue_priority": tableRedmineIssuePriority(),
			"redmine_issue_status":   tableRedmineIssueStatus(),
			"redmine_my_account":     tableRedmineMyAccount(),
			"redmine_project":        tableRedmineProject(),
			"redmine_time_entry":     tableRedmineTimeEntry(),
			"redmine_tracker":        tableRedmineTracker(),
			"redmine_user":           tableRedmineUser(),
			"redmine_version":        tableRedmineVersion(),
		},
	}
	return p
}

func configChanged(ctx context.Context, p *plugin.Plugin, old, new *plugin.Connection) error {
	config := GetConfig(new)

	if config.Endpoint == nil || *config.Endpoint == "" {
		return fmt.Errorf("'endpoint' must be set in the connection configuration or REDMINE_ENDPOINT environment variable")
	}
	if config.APIKey == nil || *config.APIKey == "" {
		return fmt.Errorf("'api_key' must be set in the connection configuration or REDMINE_API_KEY environment variable")
	}

	return nil
}
