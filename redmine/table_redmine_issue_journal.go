package redmine

import (
	"context"
	"fmt"
	"time"

	rm "github.com/nixys/nxs-go-redmine/v5"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"golang.org/x/sync/errgroup"
)

// maxConcurrentIssueFetches limits parallel IssueSingleGet calls for journal retrieval.
const maxConcurrentIssueFetches = 5

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
		Description: "Journal entries (comments and field changes) on Redmine issues.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"issue_id", "journal_id"}),
			Hydrate:    getIssueJournal,
		},
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
			// Standard columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: "The display name for this resource.", Transform: transform.FromField("IssueSubject")},
		},
	}
}

//// HELPER FUNCTIONS

// dateRange represents a half-open time interval [from, to).
type dateRange struct {
	from *time.Time
	to   *time.Time
}

// extractDateRange parses timestamp qualifiers for the given column into a dateRange.
func extractDateRange(quals plugin.KeyColumnQualMap, column ...string) dateRange {
	col := "created_on"
	if len(column) > 0 {
		col = column[0]
	}

	var dr dateRange

	if quals[col] == nil {
		return dr
	}

	for _, q := range quals[col].Quals {
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

// buildDateFilter converts a dateRange into a Redmine API date filter string.
// Uses the >< (between) operator when both bounds exist, or >= / <= for single bounds.
func buildDateFilter(dr dateRange) string {
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

//// HYDRATE FUNCTIONS

func getIssueJournal(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	issueID := d.EqualsQuals["issue_id"].GetInt64Value()
	journalID := d.EqualsQuals["journal_id"].GetInt64Value()

	issue, _, err := client.IssueSingleGet(issueID, rm.IssueSingleGetRequest{
		Includes: []rm.IssueInclude{rm.IssueIncludeJournals},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get issue %d: %w", issueID, err)
	}

	if issue.Journals == nil {
		return nil, nil
	}

	for _, journal := range *issue.Journals {
		if journal.ID == journalID {
			return issueJournalRow{
				IssueID:      issue.ID,
				IssueSubject: issue.Subject,
				ProjectID:    issue.Project.ID,
				ProjectName:  issue.Project.Name,
				JournalID:    journal.ID,
				Notes:        journal.Notes,
				CreatedOn:    parseRedmineTime(journal.CreatedOn),
				UserID:       journal.User.ID,
				UserName:     journal.User.Name,
				PrivateNotes: journal.PrivateNotes,
				Details:      journal.Details,
			}, nil
		}
	}

	return nil, nil
}

func listIssueJournals(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	dr := extractDateRange(d.Quals)

	// If a specific issue_id is provided, fetch just that issue
	if d.EqualsQuals["issue_id"] != nil {
		issueID := d.EqualsQuals["issue_id"].GetInt64Value()
		return fetchAndStreamIssueJournals(ctx, d, client, issueID, dr)
	}

	// Build filters for issue listing
	filters := rm.IssueGetRequestFiltersInit()

	// Use updated_on to narrow down candidate issues
	updatedOnFilter := buildDateFilter(dr)
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

		// Fetch journals concurrently with bounded parallelism and context cancellation
		g, gctx := errgroup.WithContext(ctx)
		g.SetLimit(maxConcurrentIssueFetches)

		for _, issue := range result.Issues {
			if d.RowsRemaining(ctx) == 0 {
				break
			}

			issueID := issue.ID
			g.Go(func() error {
				if gctx.Err() != nil {
					return gctx.Err()
				}
				_, err := fetchAndStreamIssueJournals(gctx, d, client, issueID, dr)
				if err != nil {
					// Swallow 404s: issue may have been deleted between list and get
					if isNotFoundError([]string{"404", "not found"})(err) {
						plugin.Logger(ctx).Warn("listIssueJournals", "issue_id", issueID, "msg", "issue not found, skipping", "error", err)
						return nil
					}
					plugin.Logger(ctx).Error("listIssueJournals", "issue_id", issueID, "error", err)
				}
				return err
			})
		}

		if err := g.Wait(); err != nil {
			return nil, err
		}

		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}

		if int64(len(result.Issues)) < pageSize {
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
			CreatedOn:    parseRedmineTime(journal.CreatedOn),
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
