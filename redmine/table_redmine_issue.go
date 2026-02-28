package redmine

import (
	"context"
	"fmt"
	"time"

	rm "github.com/nixys/nxs-go-redmine/v5"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

type issueRow struct {
	ID                  int64
	ProjectID           int64
	ProjectName         string
	TrackerID           int64
	TrackerName         string
	StatusID            int64
	StatusName          string
	StatusIsClosed      bool
	PriorityID          int64
	PriorityName        string
	AuthorID            int64
	AuthorName          string
	AssignedToID        int64
	AssignedToName      string
	CategoryID          int64
	CategoryName        string
	FixedVersionID      int64
	FixedVersionName    string
	ParentID            int64
	Subject             string
	Description         string
	StartDate           *time.Time
	DueDate             *time.Time
	DoneRatio           int64
	IsPrivate           bool
	EstimatedHours      *float64
	TotalEstimatedHours *float64
	SpentHours          float64
	TotalSpentHours     float64
	CustomFields        []rm.CustomFieldGetObject
	CreatedOn           *time.Time
	UpdatedOn           *time.Time
	ClosedOn            *time.Time
}

func tableRedmineIssue() *plugin.Table {
	return &plugin.Table{
		Name:        "redmine_issue",
		Description: "Issues in the Redmine instance.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getIssue,
		},
		List: &plugin.ListConfig{
			Hydrate: listIssues,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "project_id", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "tracker_id", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "status_id", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "assigned_to_id", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "assigned_to_me", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "created_on", Require: plugin.Optional, Operators: []string{">=", ">", "<", "<="}},
				{Name: "updated_on", Require: plugin.Optional, Operators: []string{">=", ">", "<", "<="}},
			},
		},
		Columns: []*plugin.Column{
			// Key columns first
			{Name: "id", Type: proto.ColumnType_INT, Description: "The issue ID."},
			{Name: "project_id", Type: proto.ColumnType_INT, Description: "The project ID."},
			{Name: "tracker_id", Type: proto.ColumnType_INT, Description: "The tracker ID."},
			{Name: "status_id", Type: proto.ColumnType_INT, Description: "The issue status ID."},
			{Name: "assigned_to_id", Type: proto.ColumnType_INT, Description: "The assigned user ID."},
			{Name: "assigned_to_me", Type: proto.ColumnType_BOOL, Description: "If true, filter to issues assigned to the API key owner. Only useful as a filter qualifier.", Transform: transform.FromConstant(false)},
			// Remaining columns alphabetically
			{Name: "assigned_to_name", Type: proto.ColumnType_STRING, Description: "The assigned user name."},
			{Name: "author_id", Type: proto.ColumnType_INT, Description: "The author user ID."},
			{Name: "author_name", Type: proto.ColumnType_STRING, Description: "The author user name."},
			{Name: "category_id", Type: proto.ColumnType_INT, Description: "The category ID."},
			{Name: "category_name", Type: proto.ColumnType_STRING, Description: "The category name."},
			{Name: "closed_on", Type: proto.ColumnType_TIMESTAMP, Description: "When the issue was closed."},
			{Name: "created_on", Type: proto.ColumnType_TIMESTAMP, Description: "When the issue was created."},
			{Name: "custom_fields", Type: proto.ColumnType_JSON, Description: "Custom field values."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "The issue description."},
			{Name: "done_ratio", Type: proto.ColumnType_INT, Description: "The percentage done (0-100)."},
			{Name: "due_date", Type: proto.ColumnType_TIMESTAMP, Description: "The due date."},
			{Name: "estimated_hours", Type: proto.ColumnType_DOUBLE, Description: "Estimated hours for the issue."},
			{Name: "fixed_version_id", Type: proto.ColumnType_INT, Description: "The target version ID."},
			{Name: "fixed_version_name", Type: proto.ColumnType_STRING, Description: "The target version name."},
			{Name: "is_private", Type: proto.ColumnType_BOOL, Description: "Whether the issue is private."},
			{Name: "parent_id", Type: proto.ColumnType_INT, Description: "The parent issue ID."},
			{Name: "priority_id", Type: proto.ColumnType_INT, Description: "The priority ID."},
			{Name: "priority_name", Type: proto.ColumnType_STRING, Description: "The priority name."},
			{Name: "project_name", Type: proto.ColumnType_STRING, Description: "The project name."},
			{Name: "spent_hours", Type: proto.ColumnType_DOUBLE, Description: "Hours spent on the issue."},
			{Name: "start_date", Type: proto.ColumnType_TIMESTAMP, Description: "The start date."},
			{Name: "status_is_closed", Type: proto.ColumnType_BOOL, Description: "Whether the status represents a closed state."},
			{Name: "status_name", Type: proto.ColumnType_STRING, Description: "The issue status name."},
			{Name: "subject", Type: proto.ColumnType_STRING, Description: "The issue subject/title."},
			{Name: "total_estimated_hours", Type: proto.ColumnType_DOUBLE, Description: "Total estimated hours including subtasks."},
			{Name: "total_spent_hours", Type: proto.ColumnType_DOUBLE, Description: "Total spent hours including subtasks."},
			{Name: "tracker_name", Type: proto.ColumnType_STRING, Description: "The tracker name."},
			{Name: "updated_on", Type: proto.ColumnType_TIMESTAMP, Description: "When the issue was last updated."},
		},
	}
}

//// HELPER FUNCTIONS

func issueRowFromObject(issue rm.IssueObject) issueRow {
	row := issueRow{
		ID:                  issue.ID,
		ProjectID:           issue.Project.ID,
		ProjectName:         issue.Project.Name,
		TrackerID:           issue.Tracker.ID,
		TrackerName:         issue.Tracker.Name,
		StatusID:            issue.Status.ID,
		StatusName:          issue.Status.Name,
		StatusIsClosed:      issue.Status.IsClosed,
		PriorityID:          issue.Priority.ID,
		PriorityName:        issue.Priority.Name,
		AuthorID:            issue.Author.ID,
		AuthorName:          issue.Author.Name,
		Subject:             issue.Subject,
		Description:         issue.Description,
		StartDate:           parseRedmineDate(issue.StartDate),
		DueDate:             parseRedmineDate(issue.DueDate),
		DoneRatio:           issue.DoneRatio,
		IsPrivate:           issue.IsPrivate != 0,
		EstimatedHours:      issue.EstimatedHours,
		TotalEstimatedHours: issue.TotalEstimatedHours,
		SpentHours:          issue.SpentHours,
		TotalSpentHours:     issue.TotalSpentHours,
		CustomFields:        issue.CustomFields,
		CreatedOn:           parseRedmineTime(issue.CreatedOn),
		UpdatedOn:           parseRedmineTime(issue.UpdatedOn),
		ClosedOn:            parseRedmineTime(issue.ClosedOn),
	}

	if issue.AssignedTo != nil {
		row.AssignedToID = issue.AssignedTo.ID
		row.AssignedToName = issue.AssignedTo.Name
	}
	if issue.Category != nil {
		row.CategoryID = issue.Category.ID
		row.CategoryName = issue.Category.Name
	}
	if issue.FixedVersion != nil {
		row.FixedVersionID = issue.FixedVersion.ID
		row.FixedVersionName = issue.FixedVersion.Name
	}
	if issue.Parent != nil {
		row.ParentID = issue.Parent.ID
	}

	return row
}

//// HYDRATE FUNCTIONS

func getIssue(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	issueID := d.EqualsQuals["id"].GetInt64Value()

	issue, _, err := client.IssueSingleGet(issueID, rm.IssueSingleGetRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to get issue %d: %w", issueID, err)
	}

	return issueRowFromObject(issue), nil
}

func listIssues(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	filters := rm.IssueGetRequestFiltersInit()

	if d.EqualsQuals["project_id"] != nil {
		filters.FieldAdd("project_id", fmt.Sprintf("%d", d.EqualsQuals["project_id"].GetInt64Value()))
	}
	if d.EqualsQuals["tracker_id"] != nil {
		filters.FieldAdd("tracker_id", fmt.Sprintf("%d", d.EqualsQuals["tracker_id"].GetInt64Value()))
	}
	if d.EqualsQuals["status_id"] != nil {
		filters.FieldAdd("status_id", fmt.Sprintf("%d", d.EqualsQuals["status_id"].GetInt64Value()))
	} else {
		// Include all statuses (Redmine defaults to open-only)
		filters.FieldAdd("status_id", "*")
	}
	if d.EqualsQuals["assigned_to_id"] != nil {
		filters.FieldAdd("assigned_to_id", fmt.Sprintf("%d", d.EqualsQuals["assigned_to_id"].GetInt64Value()))
	}
	if d.EqualsQuals["assigned_to_me"] != nil && d.EqualsQuals["assigned_to_me"].GetBoolValue() {
		filters.FieldAdd("assigned_to_id", "me")
	}

	// Date range filters
	if d.Quals["created_on"] != nil {
		dr := extractDateRange(d.Quals)
		if f := buildDateFilter(dr); f != "" {
			filters.FieldAdd("created_on", f)
		}
	}
	if d.Quals["updated_on"] != nil {
		dr := extractDateRange(d.Quals, "updated_on")
		if f := buildDateFilter(dr); f != "" {
			filters.FieldAdd("updated_on", f)
		}
	}

	sort := rm.IssueGetRequestSortInit().Set("updated_on", true)

	var offset int64
	var pageSize int64 = 100
	if d.QueryContext.Limit != nil && *d.QueryContext.Limit < pageSize {
		pageSize = *d.QueryContext.Limit
	}

	for {
		result, _, err := client.IssuesMultiGet(rm.IssueMultiGetRequest{
			Filters: filters,
			Sort:    sort,
			Offset:  offset,
			Limit:   pageSize,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list issues: %w", err)
		}

		for _, issue := range result.Issues {
			d.StreamListItem(ctx, issueRowFromObject(issue))

			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if offset+result.Limit >= result.TotalCount {
			break
		}
		offset += pageSize
	}

	return nil, nil
}
