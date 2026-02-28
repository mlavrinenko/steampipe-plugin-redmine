package redmine

import (
	"context"
	"fmt"
	"strconv"
	"time"

	rm "github.com/nixys/nxs-go-redmine/v5"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

type projectRow struct {
	ID                  int64
	Name                string
	Identifier          string
	Description         string
	Homepage            *string
	ParentID            int64
	ParentName          string
	Status              int64
	IsPublic            bool
	InheritMembers      bool
	DefaultVersionID    int64
	DefaultVersionName  string
	DefaultAssigneeID   int64
	DefaultAssigneeName string
	CustomFields        []rm.CustomFieldGetObject
	Trackers            *[]rm.IDName
	IssueCategories     *[]rm.IDName
	TimeEntryActivities *[]rm.IDName
	EnabledModules      *[]rm.IDName
	CreatedOn           *time.Time
	UpdatedOn           *time.Time
}

func tableRedmineProject() *plugin.Table {
	return &plugin.Table{
		Name:        "redmine_project",
		Description: "Projects in the Redmine instance.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AnyColumn([]string{"id", "identifier"}),
			Hydrate:    getProject,
		},
		List: &plugin.ListConfig{
			Hydrate: listProjects,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "status", Require: plugin.Optional, Operators: []string{"="}},
			},
		},
		Columns: []*plugin.Column{
			// Key columns first
			{Name: "id", Type: proto.ColumnType_INT, Description: "The project ID."},
			{Name: "identifier", Type: proto.ColumnType_STRING, Description: "The project identifier slug."},
			{Name: "status", Type: proto.ColumnType_INT, Description: "The project status (1=active, 5=closed, 9=archived)."},
			// Remaining columns alphabetically
			{Name: "created_on", Type: proto.ColumnType_TIMESTAMP, Description: "When the project was created."},
			{Name: "custom_fields", Type: proto.ColumnType_JSON, Description: "Custom field values."},
			{Name: "default_assignee_id", Type: proto.ColumnType_INT, Description: "Default assignee user ID."},
			{Name: "default_assignee_name", Type: proto.ColumnType_STRING, Description: "Default assignee user name."},
			{Name: "default_version_id", Type: proto.ColumnType_INT, Description: "Default version ID."},
			{Name: "default_version_name", Type: proto.ColumnType_STRING, Description: "Default version name."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "The project description."},
			{Name: "enabled_modules", Type: proto.ColumnType_JSON, Description: "Enabled modules."},
			{Name: "homepage", Type: proto.ColumnType_STRING, Description: "The project homepage URL."},
			{Name: "inherit_members", Type: proto.ColumnType_BOOL, Description: "Whether the project inherits members from parent."},
			{Name: "is_public", Type: proto.ColumnType_BOOL, Description: "Whether the project is public."},
			{Name: "issue_categories", Type: proto.ColumnType_JSON, Description: "Issue categories."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The project name."},
			{Name: "parent_id", Type: proto.ColumnType_INT, Description: "Parent project ID."},
			{Name: "parent_name", Type: proto.ColumnType_STRING, Description: "Parent project name."},
			{Name: "time_entry_activities", Type: proto.ColumnType_JSON, Description: "Time entry activities."},
			{Name: "trackers", Type: proto.ColumnType_JSON, Description: "Trackers."},
			{Name: "updated_on", Type: proto.ColumnType_TIMESTAMP, Description: "When the project was last updated."},
		},
	}
}

//// HELPER FUNCTIONS

func projectRowFromObject(p rm.ProjectObject) projectRow {
	row := projectRow{
		ID:                  p.ID,
		Name:                p.Name,
		Identifier:          p.Identifier,
		Description:         p.Description,
		Homepage:            p.Homepage,
		ParentID:            p.Parent.ID,
		ParentName:          p.Parent.Name,
		Status:              int64(p.Status),
		IsPublic:            p.IsPublic,
		InheritMembers:      p.InheritMembers,
		CustomFields:        p.CustomFields,
		Trackers:            p.Trackers,
		IssueCategories:     p.IssueCategories,
		TimeEntryActivities: p.TimeEntryActivities,
		EnabledModules:      p.EnabledModules,
		CreatedOn:           parseRedmineTime(p.CreatedOn),
		UpdatedOn:           parseRedmineTime(p.UpdatedOn),
	}

	if p.DefaultVersion != nil {
		row.DefaultVersionID = p.DefaultVersion.ID
		row.DefaultVersionName = p.DefaultVersion.Name
	}
	if p.DefaultAssignee != nil {
		row.DefaultAssigneeID = p.DefaultAssignee.ID
		row.DefaultAssigneeName = p.DefaultAssignee.Name
	}

	return row
}

//// HYDRATE FUNCTIONS

func getProject(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	var idStr string
	if d.EqualsQuals["id"] != nil {
		idStr = strconv.FormatInt(d.EqualsQuals["id"].GetInt64Value(), 10)
	} else if d.EqualsQuals["identifier"] != nil {
		idStr = d.EqualsQuals["identifier"].GetStringValue()
	} else {
		return nil, nil
	}

	project, _, err := client.ProjectSingleGet(idStr, rm.ProjectSingleGetRequest{
		Includes: []rm.ProjectInclude{
			rm.ProjectIncludeTrackers,
			rm.ProjectIncludeIssueCategories,
			rm.ProjectIncludeEnabledModules,
			rm.ProjectIncludeTimeEntryActivities,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get project %s: %w", idStr, err)
	}

	return projectRowFromObject(project), nil
}

func listProjects(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	includes := []rm.ProjectInclude{
		rm.ProjectIncludeTrackers,
		rm.ProjectIncludeIssueCategories,
		rm.ProjectIncludeEnabledModules,
		rm.ProjectIncludeTimeEntryActivities,
	}

	var filters *rm.ProjectGetRequestFilters
	if d.EqualsQuals["status"] != nil {
		filters = rm.ProjectGetRequestFiltersInit().
			StatusSet(rm.ProjectStatus(d.EqualsQuals["status"].GetInt64Value()))
	}

	var offset int64
	var pageSize int64 = 100
	if d.QueryContext.Limit != nil && *d.QueryContext.Limit < pageSize {
		pageSize = *d.QueryContext.Limit
	}

	for {
		result, _, err := client.ProjectMultiGet(rm.ProjectMultiGetRequest{
			Includes: includes,
			Filters:  filters,
			Offset:   offset,
			Limit:    pageSize,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list projects: %w", err)
		}

		for _, project := range result.Projects {
			d.StreamListItem(ctx, projectRowFromObject(project))

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
