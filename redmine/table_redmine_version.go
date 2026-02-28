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
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

// versionObject mirrors the Redmine Version JSON structure.
// nxs-go-redmine does not provide a versions API, so we define types here.
type versionObject struct {
	ID            int64                     `json:"id"`
	Project       rm.IDName                 `json:"project"`
	Name          string                    `json:"name"`
	Description   string                    `json:"description"`
	Status        string                    `json:"status"`
	DueDate       *string                   `json:"due_date"`
	Sharing       string                    `json:"sharing"`
	WikiPageTitle string                    `json:"wiki_page_title"`
	CustomFields  []rm.CustomFieldGetObject `json:"custom_fields"`
	CreatedOn     string                    `json:"created_on"`
	UpdatedOn     string                    `json:"updated_on"`
}

type versionSingleResult struct {
	Version versionObject `json:"version"`
}

type versionsResult struct {
	Versions   []versionObject `json:"versions"`
	TotalCount int64           `json:"total_count"`
}

type versionRow struct {
	ID            int64
	ProjectID     int64
	ProjectName   string
	Name          string
	Description   string
	Status        string
	DueDate       *time.Time
	Sharing       string
	WikiPageTitle string
	CustomFields  []rm.CustomFieldGetObject
	CreatedOn     *time.Time
	UpdatedOn     *time.Time
}

func tableRedmineVersion() *plugin.Table {
	return &plugin.Table{
		Name:        "redmine_version",
		Description: "Versions (milestones) in the Redmine instance. Listing requires a project_id.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getVersion,
		},
		List: &plugin.ListConfig{
			Hydrate:    listVersions,
			KeyColumns: plugin.SingleColumn("project_id"),
		},
		Columns: []*plugin.Column{
			// Key columns first
			{Name: "id", Type: proto.ColumnType_INT, Description: "The version ID."},
			{Name: "project_id", Type: proto.ColumnType_INT, Description: "The project ID."},
			// Remaining columns alphabetically
			{Name: "created_on", Type: proto.ColumnType_TIMESTAMP, Description: "When the version was created."},
			{Name: "custom_fields", Type: proto.ColumnType_JSON, Description: "Custom field values."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "The version description."},
			{Name: "due_date", Type: proto.ColumnType_TIMESTAMP, Description: "The version due date."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The version name."},
			{Name: "project_name", Type: proto.ColumnType_STRING, Description: "The project name."},
			{Name: "sharing", Type: proto.ColumnType_STRING, Description: "Version sharing scope (none, descendants, hierarchy, tree, system)."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "The version status (open, locked, closed)."},
			{Name: "updated_on", Type: proto.ColumnType_TIMESTAMP, Description: "When the version was last updated."},
			{Name: "wiki_page_title", Type: proto.ColumnType_STRING, Description: "Associated wiki page title."},
			// Standard columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: "The display name for this resource.", Transform: transform.FromField("Name")},
		},
	}
}

//// HELPER FUNCTIONS

func versionRowFromObject(v versionObject) versionRow {
	return versionRow{
		ID:            v.ID,
		ProjectID:     v.Project.ID,
		ProjectName:   v.Project.Name,
		Name:          v.Name,
		Description:   v.Description,
		Status:        v.Status,
		DueDate:       parseRedmineDate(v.DueDate),
		Sharing:       v.Sharing,
		WikiPageTitle: v.WikiPageTitle,
		CustomFields:  v.CustomFields,
		CreatedOn:     parseRedmineTime(v.CreatedOn),
		UpdatedOn:     parseRedmineTime(v.UpdatedOn),
	}
}

//// HYDRATE FUNCTIONS

func getVersion(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	versionID := d.EqualsQuals["id"].GetInt64Value()

	var result versionSingleResult
	_, err = client.Get(
		&result,
		url.URL{
			Path: "/versions/" + strconv.FormatInt(versionID, 10) + ".json",
		},
		http.StatusOK,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get version %d: %w", versionID, err)
	}

	return versionRowFromObject(result.Version), nil
}

func listVersions(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	projectID := d.EqualsQuals["project_id"].GetInt64Value()

	// The Redmine versions API does not support pagination; it returns all versions for a project.
	var result versionsResult
	_, err = client.Get(
		&result,
		url.URL{
			Path: "/projects/" + strconv.FormatInt(projectID, 10) + "/versions.json",
		},
		http.StatusOK,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list versions for project %d: %w", projectID, err)
	}

	for _, version := range result.Versions {
		d.StreamListItem(ctx, versionRowFromObject(version))

		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}
