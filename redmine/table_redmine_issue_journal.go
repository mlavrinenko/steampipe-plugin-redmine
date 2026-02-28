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

// issueJournalRow is a denormalized row combining issue metadata with a single journal entry.
type issueJournalRow struct {
	IssueID      int64
	IssueSubject string
	ProjectID    int64
	ProjectName  string
	JournalID    int64
	Notes        string
	CreatedOn    *time.Time
	UserID       int64
	UserName     string
	PrivateNotes bool
	Details      []rm.IssueJournalDetailObject
}

func tableRedmineIssueJournal() *plugin.Table {
	return &plugin.Table{
		Name:        "redmine_issue_journal",
		Description: "Journal entries (comments and field changes) on Redmine issues where the current user has participated.",
		List: &plugin.ListConfig{
			Hydrate: listIssueJournals,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "created_on", Require: plugin.Required, Operators: []string{">=", ">", "<", "<="}},
				{Name: "issue_id", Require: plugin.Optional},
				{Name: "project_id", Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			// Key columns first
			{Name: "issue_id", Type: proto.ColumnType_INT, Description: "The issue ID."},
			{Name: "project_id", Type: proto.ColumnType_INT, Description: "The project ID."},
			{Name: "created_on", Type: proto.ColumnType_TIMESTAMP, Description: "When the journal entry was created."},
			// Remaining columns alphabetically
			{Name: "details", Type: proto.ColumnType_JSON, Description: "Field change details.", Transform: transform.FromField("Details")},
			{Name: "issue_subject", Type: proto.ColumnType_STRING, Description: "The issue subject."},
			{Name: "journal_id", Type: proto.ColumnType_INT, Description: "The journal entry ID."},
			{Name: "notes", Type: proto.ColumnType_STRING, Description: "The journal notes/comment text."},
			{Name: "private_notes", Type: proto.ColumnType_BOOL, Description: "Whether the note is private."},
			{Name: "project_name", Type: proto.ColumnType_STRING, Description: "The project name."},
			{Name: "user_id", Type: proto.ColumnType_INT, Description: "ID of the user who created the journal entry."},
			{Name: "user_name", Type: proto.ColumnType_STRING, Description: "Name of the user who created the journal entry."},
		},
	}
}

//// HELPER FUNCTIONS

// dateRange represents a half-open time interval [from, to).
type dateRange struct {
	from *time.Time
	to   *time.Time
}

// extractDateRange parses created_on qualifiers into a dateRange.
func extractDateRange(quals plugin.KeyColumnQualMap) dateRange {
	var dr dateRange

	if quals["created_on"] == nil {
		return dr
	}

	for _, q := range quals["created_on"].Quals {
		ts := q.Value.GetTimestampValue().AsTime()
		switch q.Operator {
		case ">=":
			t := ts
			dr.from = &t
		case ">":
			t := ts.Add(time.Second)
			dr.from = &t
		case "<=":
			t := ts.Add(time.Second)
			dr.to = &t
		case "<":
			t := ts
			dr.to = &t
		}
	}

	return dr
}

// journalInRange checks if a journal's created_on falls within the date range.
func journalInRange(journalCreatedOn string, dr dateRange) bool {
	t, err := time.Parse(time.RFC3339, journalCreatedOn)
	if err != nil {
		// Try alternative format without timezone
		t, err = time.Parse("2006-01-02T15:04:05Z", journalCreatedOn)
		if err != nil {
			return false
		}
	}

	if dr.from != nil && t.Before(*dr.from) {
		return false
	}
	if dr.to != nil && !t.Before(*dr.to) {
		return false
	}

	return true
}

// userParticipated checks if the given user ID authored any journal entry in the list.
func userParticipated(journals []rm.IssueJournalObject, userID int64) bool {
	for _, j := range journals {
		if j.User.ID == userID {
			return true
		}
	}
	return false
}

// buildUpdatedOnFilter converts a dateRange into a Redmine API updated_on filter string.
// Uses the >< (between) operator when both bounds exist, or >= / <= for single bounds.
func buildUpdatedOnFilter(dr dateRange) string {
	layout := "2006-01-02T15:04:05Z"

	if dr.from != nil && dr.to != nil {
		return "><" + dr.from.Format(layout) + "|" + dr.to.Format(layout)
	}
	if dr.from != nil {
		return ">=" + dr.from.Format(layout)
	}
	if dr.to != nil {
		return "<=" + dr.to.Format(layout)
	}

	return ""
}

func parseJournalTime(s string) *time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		t, err = time.Parse("2006-01-02T15:04:05Z", s)
		if err != nil {
			return nil
		}
	}
	return &t
}

//// HYDRATE FUNCTIONS

func listIssueJournals(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	currentUserID, err := getCurrentUserID(ctx, d)
	if err != nil {
		return nil, err
	}

	dr := extractDateRange(d.Quals)

	// If a specific issue_id is provided, fetch just that issue
	if d.EqualsQuals["issue_id"] != nil {
		issueID := d.EqualsQuals["issue_id"].GetInt64Value()
		return fetchAndStreamIssueJournals(ctx, d, client, issueID, currentUserID, dr)
	}

	// Build filters for issue listing
	filters := rm.IssueGetRequestFiltersInit()

	// Use updated_on to narrow down candidate issues
	updatedOnFilter := buildUpdatedOnFilter(dr)
	if updatedOnFilter != "" {
		filters.FieldAdd("updated_on", updatedOnFilter)
	}

	if d.EqualsQuals["project_id"] != nil {
		filters.FieldAdd("project_id", fmt.Sprintf("%d", d.EqualsQuals["project_id"].GetInt64Value()))
	}

	// Include all statuses (open and closed)
	filters.FieldAdd("status_id", "*")

	// Sort by updated_on descending for most relevant results first
	sort := rm.IssueGetRequestSortInit().Set("updated_on", true)

	// Paginate through issues; reduce page size if query has a small limit
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
			_, err := fetchAndStreamIssueJournals(ctx, d, client, issue.ID, currentUserID, dr)
			if err != nil {
				plugin.Logger(ctx).Error("listIssueJournals", "issue_id", issue.ID, "error", err)
				continue
			}

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

func fetchAndStreamIssueJournals(
	ctx context.Context,
	d *plugin.QueryData,
	client *rm.Context,
	issueID int64,
	currentUserID int64,
	dr dateRange,
) (interface{}, error) {
	issue, _, err := client.IssueSingleGet(issueID, rm.IssueSingleGetRequest{
		Includes: []rm.IssueInclude{rm.IssueIncludeJournals},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get issue %d: %w", issueID, err)
	}

	if issue.Journals == nil {
		return nil, nil
	}

	journals := *issue.Journals

	// Only include this issue if the current user participated
	if !userParticipated(journals, currentUserID) {
		return nil, nil
	}

	// Stream all journals within the date range
	for _, journal := range journals {
		if !journalInRange(journal.CreatedOn, dr) {
			continue
		}

		row := issueJournalRow{
			IssueID:      issue.ID,
			IssueSubject: issue.Subject,
			ProjectID:    issue.Project.ID,
			ProjectName:  issue.Project.Name,
			JournalID:    journal.ID,
			Notes:        journal.Notes,
			CreatedOn:    parseJournalTime(journal.CreatedOn),
			UserID:       journal.User.ID,
			UserName:     journal.User.Name,
			PrivateNotes: journal.PrivateNotes,
			Details:      journal.Details,
		}

		d.StreamListItem(ctx, row)

		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}
