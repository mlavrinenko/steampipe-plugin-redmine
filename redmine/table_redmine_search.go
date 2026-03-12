package redmine

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

// searchResponse represents the JSON response from /search.json.
type searchResponse struct {
	Results    []searchResultObject `json:"results"`
	TotalCount int64                `json:"total_count"`
	Offset     int64                `json:"offset"`
	Limit      int64                `json:"limit"`
}

type searchResultObject struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Type        string `json:"type"`
	URL         string `json:"url"`
	Description string `json:"description"`
	Datetime    string `json:"datetime"`
}

type searchRow struct {
	ID          int64
	Type        string
	URL         string
	Description string
	Datetime    *time.Time
	Akas        []string
	Title       string
}

func tableRedmineSearch() *plugin.Table {
	return &plugin.Table{
		Name:        "redmine_search",
		Description: "Search results from the Redmine instance.",
		List: &plugin.ListConfig{
			Hydrate: listSearchResults,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "q", Require: plugin.Required, Operators: []string{"="}},
				{Name: "scope", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "all_words", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "titles_only", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "open_issues", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "attachments", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "issues", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "news", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "documents", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "changesets", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "wiki_pages", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "messages", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "projects", Require: plugin.Optional, Operators: []string{"="}},
			},
		},
		Columns: []*plugin.Column{
			// Key columns first
			{Name: "q", Type: proto.ColumnType_STRING, Description: "The search query string (required).", Transform: transform.FromQual("q")},
			{Name: "id", Type: proto.ColumnType_INT, Description: "The resource ID of the search result."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "The resource type (e.g. issue, wiki-page, news, changeset, message, project, document)."},
			// Filter-only columns
			{Name: "scope", Type: proto.ColumnType_STRING, Description: "Filter-only: search scope (all, my_project, subprojects).", Transform: transform.FromQual("scope")},
			{Name: "all_words", Type: proto.ColumnType_BOOL, Description: "Filter-only: true to match all words, false for any word.", Transform: transform.FromQual("all_words")},
			{Name: "titles_only", Type: proto.ColumnType_BOOL, Description: "Filter-only: true to search titles only.", Transform: transform.FromQual("titles_only")},
			{Name: "open_issues", Type: proto.ColumnType_BOOL, Description: "Filter-only: true to return only open issues.", Transform: transform.FromQual("open_issues")},
			{Name: "attachments", Type: proto.ColumnType_STRING, Description: "Filter-only: attachment search mode (0=description only, 1=description+attachments, only=attachments only).", Transform: transform.FromQual("attachments")},
			{Name: "issues", Type: proto.ColumnType_BOOL, Description: "Filter-only: include issues in results.", Transform: transform.FromQual("issues")},
			{Name: "news", Type: proto.ColumnType_BOOL, Description: "Filter-only: include news in results.", Transform: transform.FromQual("news")},
			{Name: "documents", Type: proto.ColumnType_BOOL, Description: "Filter-only: include documents in results.", Transform: transform.FromQual("documents")},
			{Name: "changesets", Type: proto.ColumnType_BOOL, Description: "Filter-only: include changesets in results.", Transform: transform.FromQual("changesets")},
			{Name: "wiki_pages", Type: proto.ColumnType_BOOL, Description: "Filter-only: include wiki pages in results.", Transform: transform.FromQual("wiki_pages")},
			{Name: "messages", Type: proto.ColumnType_BOOL, Description: "Filter-only: include forum messages in results.", Transform: transform.FromQual("messages")},
			{Name: "projects", Type: proto.ColumnType_BOOL, Description: "Filter-only: include projects in results.", Transform: transform.FromQual("projects")},
			// Data columns
			{Name: "datetime", Type: proto.ColumnType_TIMESTAMP, Description: "The timestamp of the search result."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "A text snippet or description of the search result."},
			{Name: "url", Type: proto.ColumnType_STRING, Description: "The URL path to the resource."},
			// Standard columns
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings for the resource."},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "The display name for this resource."},
		},
	}
}

//// HELPER FUNCTIONS

func searchRowFromObject(obj searchResultObject) searchRow {
	return searchRow{
		ID:          obj.ID,
		Title:       obj.Title,
		Type:        obj.Type,
		URL:         obj.URL,
		Description: obj.Description,
		Datetime:    parseRedmineTime(obj.Datetime),
		Akas:        []string{obj.URL},
	}
}

//// HYDRATE FUNCTIONS

func listSearchResults(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	params := url.Values{}

	params.Set("q", d.EqualsQuals["q"].GetStringValue())

	if d.EqualsQuals["scope"] != nil {
		params.Set("scope", d.EqualsQuals["scope"].GetStringValue())
	}
	if d.EqualsQuals["attachments"] != nil {
		params.Set("attachments", d.EqualsQuals["attachments"].GetStringValue())
	}

	// Boolean filter params
	boolParams := map[string]string{
		"all_words":   "all_words",
		"titles_only": "titles_only",
		"open_issues": "open_issues",
		"issues":      "issues",
		"news":        "news",
		"documents":   "documents",
		"changesets":  "changesets",
		"wiki_pages":  "wiki_pages",
		"messages":    "messages",
		"projects":    "projects",
	}
	for col, param := range boolParams {
		if d.EqualsQuals[col] != nil {
			if d.EqualsQuals[col].GetBoolValue() {
				params.Set(param, "1")
			} else {
				params.Set(param, "0")
			}
		}
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

		var result searchResponse
		_, err := client.Get(
			&result,
			url.URL{
				Path:     "/search.json",
				RawQuery: params.Encode(),
			},
			http.StatusOK,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to search: %w", err)
		}

		for _, item := range result.Results {
			d.StreamListItem(ctx, searchRowFromObject(item))

			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if int64(len(result.Results)) < pageSize {
			break
		}
		offset += pageSize
	}

	return nil, nil
}
