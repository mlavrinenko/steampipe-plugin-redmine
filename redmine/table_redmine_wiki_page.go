package redmine

import (
	"context"
	"fmt"
	"time"

	rm "github.com/nixys/nxs-go-redmine/v5"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

type wikiPageRow struct {
	ProjectID   string
	Title       string
	ParentTitle *string
	Text        *string
	Version     int64
	AuthorID    *int64
	AuthorName  *string
	Comments    *string
	CreatedOn   *time.Time
	UpdatedOn   *time.Time
	Akas        []string
}

func tableRedmineWikiPage() *plugin.Table {
	return &plugin.Table{
		Name:        "redmine_wiki_page",
		Description: "Wiki pages in the Redmine instance. Listing requires a project_id (identifier string).",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"project_id", "title"}),
			Hydrate:    getWikiPage,
		},
		List: &plugin.ListConfig{
			Hydrate:    listWikiPages,
			KeyColumns: plugin.SingleColumn("project_id"),
		},
		Columns: []*plugin.Column{
			// Key columns first
			{Name: "project_id", Type: proto.ColumnType_STRING, Description: "The project identifier (string slug, not numeric ID)."},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "The wiki page title."},
			// Remaining columns alphabetically
			{Name: "author_id", Type: proto.ColumnType_INT, Description: "The author user ID (available on single-get only)."},
			{Name: "author_name", Type: proto.ColumnType_STRING, Description: "The author name (available on single-get only)."},
			{Name: "comments", Type: proto.ColumnType_STRING, Description: "Version comments (available on single-get only)."},
			{Name: "created_on", Type: proto.ColumnType_TIMESTAMP, Description: "When the wiki page was created."},
			{Name: "parent_title", Type: proto.ColumnType_STRING, Description: "The parent wiki page title."},
			{Name: "text", Type: proto.ColumnType_STRING, Description: "The wiki page content in textile/markdown (available on single-get only)."},
			{Name: "updated_on", Type: proto.ColumnType_TIMESTAMP, Description: "When the wiki page was last updated."},
			{Name: "version", Type: proto.ColumnType_INT, Description: "The wiki page version number."},
			// Standard columns
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings for the resource."},
		},
	}
}

//// HELPER FUNCTIONS

func wikiPageRowFromMultiObject(projectID string, w rm.WikiMultiObject) wikiPageRow {
	var parentTitle *string
	if w.Parent != nil {
		parentTitle = &w.Parent.Title
	}

	return wikiPageRow{
		ProjectID:   projectID,
		Title:       w.Title,
		ParentTitle: parentTitle,
		Version:     w.Version,
		CreatedOn:   parseRedmineTime(w.CreatedOn),
		UpdatedOn:   parseRedmineTime(w.UpdatedOn),
		Akas:        []string{fmt.Sprintf("/projects/%s/wiki/%s", projectID, w.Title)},
	}
}

func wikiPageRowFromObject(projectID string, w rm.WikiObject) wikiPageRow {
	var parentTitle *string
	if w.Parent != nil {
		parentTitle = &w.Parent.Title
	}

	var authorID *int64
	var authorName *string
	if w.Author.ID != 0 {
		authorID = &w.Author.ID
		authorName = &w.Author.Name
	}

	var comments *string
	if w.Comments != "" {
		comments = &w.Comments
	}

	var text *string
	if w.Text != "" {
		text = &w.Text
	}

	return wikiPageRow{
		ProjectID:   projectID,
		Title:       w.Title,
		ParentTitle: parentTitle,
		Text:        text,
		Version:     w.Version,
		AuthorID:    authorID,
		AuthorName:  authorName,
		Comments:    comments,
		CreatedOn:   parseRedmineTime(w.CreatedOn),
		UpdatedOn:   parseRedmineTime(w.UpdatedOn),
		Akas:        []string{fmt.Sprintf("/projects/%s/wiki/%s", projectID, w.Title)},
	}
}

//// HYDRATE FUNCTIONS

func getWikiPage(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	projectID := d.EqualsQuals["project_id"].GetStringValue()
	title := d.EqualsQuals["title"].GetStringValue()

	wiki, _, err := client.WikiSingleGet(projectID, title, rm.WikiSingleGetRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to get wiki page %s/%s: %w", projectID, title, err)
	}

	return wikiPageRowFromObject(projectID, wiki), nil
}

func listWikiPages(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	projectID := d.EqualsQuals["project_id"].GetStringValue()

	pages, _, err := client.WikiAllGet(projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to list wiki pages for project %s: %w", projectID, err)
	}

	for _, page := range pages {
		d.StreamListItem(ctx, wikiPageRowFromMultiObject(projectID, page))

		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}
