package redmine

import (
	"context"
	"fmt"

	rm "github.com/nixys/nxs-go-redmine/v5"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

type issuePriorityRow struct {
	ID        int64
	Name      string
	IsDefault bool
	Active    bool
}

func tableRedmineIssuePriority() *plugin.Table {
	return &plugin.Table{
		Name:        "redmine_issue_priority",
		Description: "Issue priority definitions in the Redmine instance.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getIssuePriority,
		},
		List: &plugin.ListConfig{
			Hydrate: listIssuePriorities,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_INT, Description: "The priority ID."},
			{Name: "active", Type: proto.ColumnType_BOOL, Description: "Whether the priority is active."},
			{Name: "is_default", Type: proto.ColumnType_BOOL, Description: "Whether this is the default priority."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The priority name."},
			// Standard columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: "The display name for this resource.", Transform: transform.FromField("Name")},
		},
	}
}

//// HELPER FUNCTIONS

func issuePriorityRowFromObject(p rm.EnumerationPriorityObject) issuePriorityRow {
	return issuePriorityRow{
		ID:        p.ID,
		Name:      p.Name,
		IsDefault: p.IsDefault,
		Active:    p.Active,
	}
}

//// HYDRATE FUNCTIONS

func getIssuePriority(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	priorityID := d.EqualsQuals["id"].GetInt64Value()

	priorities, _, err := client.EnumerationPrioritiesAllGet()
	if err != nil {
		return nil, fmt.Errorf("failed to list issue priorities: %w", err)
	}

	for _, priority := range priorities {
		if priority.ID == priorityID {
			return issuePriorityRowFromObject(priority), nil
		}
	}

	return nil, nil
}

func listIssuePriorities(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	priorities, _, err := client.EnumerationPrioritiesAllGet()
	if err != nil {
		return nil, fmt.Errorf("failed to list issue priorities: %w", err)
	}

	for _, priority := range priorities {
		d.StreamListItem(ctx, issuePriorityRowFromObject(priority))

		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}
