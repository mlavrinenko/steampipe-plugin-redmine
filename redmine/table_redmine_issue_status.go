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

type issueStatusRow struct {
	ID       int64
	Name     string
	IsClosed bool
	Akas     []string
}

func tableRedmineIssueStatus() *plugin.Table {
	return &plugin.Table{
		Name:        "redmine_issue_status",
		Description: "Issue statuses in the Redmine instance.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getIssueStatus,
		},
		List: &plugin.ListConfig{
			Hydrate: listIssueStatuses,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_INT, Description: "The status ID."},
			{Name: "is_closed", Type: proto.ColumnType_BOOL, Description: "Whether this status represents a closed state."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The status name."},
			// Standard columns
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings for the resource."},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "The display name for this resource.", Transform: transform.FromField("Name")},
		},
	}
}

//// HELPER FUNCTIONS

func issueStatusRowFromObject(s rm.IssueStatusObject) issueStatusRow {
	return issueStatusRow{
		ID:       s.ID,
		Name:     s.Name,
		IsClosed: s.IsClosed,
		Akas:     []string{fmt.Sprintf("/issue_statuses/%d", s.ID)},
	}
}

//// HYDRATE FUNCTIONS

func getAllIssueStatuses(ctx context.Context, d *plugin.QueryData) ([]rm.IssueStatusObject, error) {
	cacheKey := "redmine_issue_statuses"
	if cached, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cached.([]rm.IssueStatusObject), nil
	}

	client, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	statuses, _, err := client.IssueStatusAllGet()
	if err != nil {
		return nil, fmt.Errorf("failed to list issue statuses: %w", err)
	}

	d.ConnectionManager.Cache.Set(cacheKey, statuses)
	return statuses, nil
}

func getIssueStatus(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	statusID := d.EqualsQuals["id"].GetInt64Value()

	statuses, err := getAllIssueStatuses(ctx, d)
	if err != nil {
		return nil, err
	}

	for _, status := range statuses {
		if status.ID == statusID {
			return issueStatusRowFromObject(status), nil
		}
	}

	return nil, nil
}

func listIssueStatuses(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	statuses, err := getAllIssueStatuses(ctx, d)
	if err != nil {
		return nil, err
	}

	for _, status := range statuses {
		d.StreamListItem(ctx, issueStatusRowFromObject(status))

		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}
