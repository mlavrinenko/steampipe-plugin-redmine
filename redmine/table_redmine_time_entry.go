package redmine

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	rm "github.com/nixys/nxs-go-redmine/v5"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

type timeEntryRow struct {
	ID           int64
	ProjectID    int64
	ProjectName  string
	IssueID      int64
	UserID       int64
	UserName     string
	ActivityID   int64
	ActivityName string
	Hours        float64
	Comments     string
	SpentOn      *time.Time
	CreatedOn    *time.Time
	UpdatedOn    *time.Time
	Title        string
}

func tableRedmineTimeEntry() *plugin.Table {
	return &plugin.Table{
		Name:        "redmine_time_entry",
		Description: "Time entries in the Redmine instance.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getTimeEntry,
		},
		List: &plugin.ListConfig{
			Hydrate: listTimeEntries,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "project_id", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "issue_id", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "user_id", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "activity_id", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "spent_on", Require: plugin.Optional, Operators: []string{">=", ">", "<", "<="}},
			},
		},
		Columns: []*plugin.Column{
			// Key columns first
			{Name: "id", Type: proto.ColumnType_INT, Description: "The time entry ID."},
			{Name: "project_id", Type: proto.ColumnType_INT, Description: "The project ID."},
			{Name: "user_id", Type: proto.ColumnType_INT, Description: "The user ID."},
			{Name: "activity_id", Type: proto.ColumnType_INT, Description: "The activity ID."},
			{Name: "spent_on", Type: proto.ColumnType_TIMESTAMP, Description: "The date the time was spent."},
			// Remaining columns alphabetically
			{Name: "activity_name", Type: proto.ColumnType_STRING, Description: "The activity name."},
			{Name: "comments", Type: proto.ColumnType_STRING, Description: "Comments for the time entry."},
			{Name: "created_on", Type: proto.ColumnType_TIMESTAMP, Description: "When the time entry was created."},
			{Name: "hours", Type: proto.ColumnType_DOUBLE, Description: "Hours logged."},
			{Name: "issue_id", Type: proto.ColumnType_INT, Description: "The issue ID."},
			{Name: "project_name", Type: proto.ColumnType_STRING, Description: "The project name."},
			{Name: "updated_on", Type: proto.ColumnType_TIMESTAMP, Description: "When the time entry was last updated."},
			{Name: "user_name", Type: proto.ColumnType_STRING, Description: "The user name."},
			// Standard columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: "The display name for this resource."},
		},
	}
}

//// HELPER FUNCTIONS

func timeEntryRowFromObject(te rm.TimeEntryObject) timeEntryRow {
	title := fmt.Sprintf("%.2fh on %s", te.Hours, te.Project.Name)
	if te.Comments != "" {
		title = fmt.Sprintf("%.2fh on %s: %s", te.Hours, te.Project.Name, te.Comments)
	}

	return timeEntryRow{
		ID:           te.ID,
		ProjectID:    te.Project.ID,
		ProjectName:  te.Project.Name,
		IssueID:      te.Issue.ID,
		UserID:       te.User.ID,
		UserName:     te.User.Name,
		ActivityID:   te.Activity.ID,
		ActivityName: te.Activity.Name,
		Hours:        te.Hours,
		Comments:     te.Comments,
		SpentOn:      parseRedmineDate(&te.SpentOn),
		CreatedOn:    parseRedmineTime(te.CreatedOn),
		UpdatedOn:    parseRedmineTime(te.UpdatedOn),
		Title:        title,
	}
}

// extractSpentOnRange parses spent_on qualifiers into from/to date strings (YYYY-MM-DD).
func extractSpentOnRange(quals plugin.KeyColumnQualMap) (from, to string) {
	if quals["spent_on"] == nil {
		return "", ""
	}

	for _, q := range quals["spent_on"].Quals {
		ts := q.Value.GetTimestampValue().AsTime()
		date, isFrom := adjustDateBound(q.Operator, ts)
		if isFrom {
			from = date
		} else {
			to = date
		}
	}

	return from, to
}

//// HYDRATE FUNCTIONS

func getTimeEntry(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	id := d.EqualsQuals["id"].GetInt64Value()

	te, _, err := client.TimeEntrySingleGet(id, rm.TimeEntrySingleGetRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to get time entry %d: %w", id, err)
	}

	return timeEntryRowFromObject(te), nil
}

func listTimeEntries(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	// We use raw client.Get() with manual pagination instead of the library's
	// TimeEntryAllGet() for two reasons:
	// 1. TimeEntryAllGet fetches all pages into memory (no early exit via RowsRemaining)
	// 2. The library's TimeEntryGetRequestFilters doesn't support issue_id filtering
	params := url.Values{}

	if d.EqualsQuals["project_id"] != nil {
		params.Set("project_id", strconv.FormatInt(d.EqualsQuals["project_id"].GetInt64Value(), 10))
	}
	if d.EqualsQuals["issue_id"] != nil {
		params.Set("issue_id", strconv.FormatInt(d.EqualsQuals["issue_id"].GetInt64Value(), 10))
	}
	if d.EqualsQuals["user_id"] != nil {
		params.Set("user_id", strconv.FormatInt(d.EqualsQuals["user_id"].GetInt64Value(), 10))
	}
	if d.EqualsQuals["activity_id"] != nil {
		params.Set("activity_id", strconv.FormatInt(d.EqualsQuals["activity_id"].GetInt64Value(), 10))
	}

	from, to := extractSpentOnRange(d.Quals)
	if from != "" {
		params.Set("from", from)
	}
	if to != "" {
		params.Set("to", to)
	}

	var offset int64
	var pageSize int64 = 100
	if d.QueryContext.Limit != nil && *d.QueryContext.Limit < pageSize {
		pageSize = *d.QueryContext.Limit
	}

	params.Set("limit", strconv.FormatInt(pageSize, 10))

	for {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		params.Set("offset", strconv.FormatInt(offset, 10))

		var result rm.TimeEntryResult
		_, err := client.Get(
			&result,
			url.URL{
				Path:     "/time_entries.json",
				RawQuery: params.Encode(),
			},
			http.StatusOK,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to list time entries: %w", err)
		}

		for _, te := range result.TimeEntries {
			d.StreamListItem(ctx, timeEntryRowFromObject(te))

			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if int64(len(result.TimeEntries)) < pageSize {
			break
		}
		offset += pageSize
	}

	return nil, nil
}
