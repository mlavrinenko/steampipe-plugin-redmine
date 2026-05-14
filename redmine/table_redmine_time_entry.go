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
	Akas         []string
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
				{Name: "spent_on", Require: plugin.Optional, Operators: []string{"=", ">=", ">", "<", "<="}},
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
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings for the resource."},
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
		Akas:         []string{fmt.Sprintf("/time_entries/%d", te.ID)},
	}
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
	baseParams := url.Values{}

	if d.EqualsQuals["issue_id"] != nil {
		baseParams.Set("issue_id", strconv.FormatInt(d.EqualsQuals["issue_id"].GetInt64Value(), 10))
	}
	if d.EqualsQuals["user_id"] != nil {
		baseParams.Set("user_id", strconv.FormatInt(d.EqualsQuals["user_id"].GetInt64Value(), 10))
	}
	if d.EqualsQuals["activity_id"] != nil {
		baseParams.Set("activity_id", strconv.FormatInt(d.EqualsQuals["activity_id"].GetInt64Value(), 10))
	}

	from, to := extractSpentOnRange(d.Quals)
	if from != "" {
		baseParams.Set("from", from)
	}
	if to != "" {
		baseParams.Set("to", to)
	}

	// Redmine's GET /time_entries.json?project_id=P returns entries from P AND
	// all its descendant projects. When the user writes `project_id IN (P1..Pn)`,
	// the SDK fan-out runs this list hydrate once per IN value in its own
	// goroutine, and the same entry can come back from the call for its leaf
	// project, its parent, its grandparent, etc. — so the naive merged result
	// double-counts.
	//
	// Fix: read the original IN list from QueryContext.UnsafeQuals (shared
	// across fan-out goroutines), elect a single "leader" (the goroutine whose
	// scalar EqualsQuals["project_id"] matches the smallest IN value), and
	// have it issue one REST sequence per distinct project_id while deduping
	// the merged stream by entry ID. Non-leaders no-op so we still fire one
	// call per project_id but stream each entry at most once.
	//
	// When project_id IN (...) co-occurs with another IN qual (e.g. user_id),
	// the SDK skips fan-out entirely (no required IN qual to pick) and calls
	// this function once with EqualsQuals["project_id"] == nil; the same
	// expand-and-dedup path still applies.
	if projectIDs := extractInt64InList(d.QueryContext.UnsafeQuals, "project_id"); projectIDs != nil {
		if currentQ := d.EqualsQuals["project_id"]; currentQ != nil && currentQ.GetInt64Value() != projectIDs[0] {
			return nil, nil
		}
		seen := make(map[int64]struct{})
		for _, pid := range projectIDs {
			params := cloneURLValues(baseParams)
			params.Set("project_id", strconv.FormatInt(pid, 10))
			stopped, err := streamTimeEntries(ctx, d, client, params, seen)
			if err != nil {
				return nil, err
			}
			if stopped {
				return nil, nil
			}
		}
		return nil, nil
	}

	if d.EqualsQuals["project_id"] != nil {
		baseParams.Set("project_id", strconv.FormatInt(d.EqualsQuals["project_id"].GetInt64Value(), 10))
	}
	if _, err := streamTimeEntries(ctx, d, client, baseParams, nil); err != nil {
		return nil, err
	}
	return nil, nil
}

// streamTimeEntries paginates GET /time_entries.json with the given params and
// streams each entry to the SDK. It returns (stopped, err): stopped is true
// when the row-limit has been reached and no further pages should be fetched.
// When `seen` is non-nil, entries whose IDs are already in the set are skipped
// (and the rest are added to it) so successive calls with overlapping result
// sets stream each entry at most once.
func streamTimeEntries(ctx context.Context, d *plugin.QueryData, client *rm.Context, params url.Values, seen map[int64]struct{}) (bool, error) {
	var offset int64
	var pageSize int64 = 100
	if d.QueryContext.Limit != nil && *d.QueryContext.Limit < pageSize {
		pageSize = *d.QueryContext.Limit
	}
	params.Set("limit", strconv.FormatInt(pageSize, 10))

	for {
		if ctx.Err() != nil {
			return false, ctx.Err()
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
			return false, fmt.Errorf("failed to list time entries: %w", err)
		}

		for _, te := range result.TimeEntries {
			if seen != nil {
				if _, dup := seen[te.ID]; dup {
					continue
				}
				seen[te.ID] = struct{}{}
			}
			d.StreamListItem(ctx, timeEntryRowFromObject(te))
			if d.RowsRemaining(ctx) == 0 {
				return true, nil
			}
		}

		if int64(len(result.TimeEntries)) < pageSize {
			return false, nil
		}
		offset += pageSize
	}
}

func cloneURLValues(v url.Values) url.Values {
	out := make(url.Values, len(v))
	for k, vs := range v {
		out[k] = append([]string(nil), vs...)
	}
	return out
}
