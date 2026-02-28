package redmine

import (
	"context"
	"fmt"

	rm "github.com/nixys/nxs-go-redmine/v5"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

type issueStatusRow struct {
	ID       int64
	Name     string
	IsClosed bool
}

func tableRedmineIssueStatus() *plugin.Table {
	return &plugin.Table{
		Name:        "redmine_issue_status",
		Description: "Issue statuses in the Redmine instance.",
		List: &plugin.ListConfig{
			Hydrate: listIssueStatuses,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_INT, Description: "The status ID."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The status name."},
			{Name: "is_closed", Type: proto.ColumnType_BOOL, Description: "Whether this status represents a closed state."},
		},
	}
}

//// HELPER FUNCTIONS

func issueStatusRowFromObject(s rm.IssueStatusObject) issueStatusRow {
	return issueStatusRow{
		ID:       s.ID,
		Name:     s.Name,
		IsClosed: s.IsClosed,
	}
}

//// HYDRATE FUNCTIONS

func listIssueStatuses(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	statuses, _, err := client.IssueStatusAllGet()
	if err != nil {
		return nil, fmt.Errorf("failed to list issue statuses: %w", err)
	}

	for _, status := range statuses {
		d.StreamListItem(ctx, issueStatusRowFromObject(status))

		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}
