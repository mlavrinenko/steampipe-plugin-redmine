package redmine

import (
	"context"
	"fmt"
	"strconv"

	rm "github.com/nixys/nxs-go-redmine/v5"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

type projectMembershipRow struct {
	ID                int64
	ProjectID         int64
	ProjectName       string
	ProjectIdentifier string
	UserID            int64
	UserName          string
	GroupID           int64
	GroupName         string
	Roles             []rm.MembershipRoleObject
	Akas              []string
}

func tableRedmineProjectMembership() *plugin.Table {
	return &plugin.Table{
		Name:        "redmine_project_membership",
		Description: "Project memberships in Redmine. Each row is one user or group assigned to a project with specific roles.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getProjectMembership,
		},
		List: &plugin.ListConfig{
			Hydrate: listProjectMemberships,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "project_id", Require: plugin.AnyOf, Operators: []string{"="}},
				{Name: "project_identifier", Require: plugin.AnyOf, Operators: []string{"="}},
			},
		},
		Columns: []*plugin.Column{
			// Key columns first
			{Name: "id", Type: proto.ColumnType_INT, Description: "The membership ID."},
			{Name: "project_id", Type: proto.ColumnType_INT, Description: "The project ID."},
			{Name: "project_identifier", Type: proto.ColumnType_STRING, Description: "The project identifier slug (only populated when listing by project_identifier)."},
			// Remaining columns alphabetically
			{Name: "group_id", Type: proto.ColumnType_INT, Description: "The group ID (if membership belongs to a group)."},
			{Name: "group_name", Type: proto.ColumnType_STRING, Description: "The group name (if membership belongs to a group)."},
			{Name: "project_name", Type: proto.ColumnType_STRING, Description: "The project name."},
			{Name: "roles", Type: proto.ColumnType_JSON, Description: "Roles assigned to this member."},
			{Name: "user_id", Type: proto.ColumnType_INT, Description: "The user ID (if membership belongs to a user)."},
			{Name: "user_name", Type: proto.ColumnType_STRING, Description: "The user name (if membership belongs to a user)."},
			// Standard columns
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings for the resource."},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "The display name for this resource.", Transform: transform.FromField("UserName")},
		},
	}
}

//// HELPER FUNCTIONS

func projectMembershipRowFromObject(m rm.MembershipObject, projectIdentifier string) projectMembershipRow {
	row := projectMembershipRow{
		ID:                m.ID,
		ProjectID:         m.Project.ID,
		ProjectName:       m.Project.Name,
		ProjectIdentifier: projectIdentifier,
		Roles:             m.Roles,
		Akas:              []string{fmt.Sprintf("/memberships/%d", m.ID)},
	}

	if m.User != nil {
		row.UserID = m.User.ID
		row.UserName = m.User.Name
	}
	if m.Group != nil {
		row.GroupID = m.Group.ID
		row.GroupName = m.Group.Name
	}

	return row
}

//// HYDRATE FUNCTIONS

func getProjectMembership(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	membershipID := d.EqualsQuals["id"].GetInt64Value()

	membership, _, err := client.MembershipSingleGet(membershipID)
	if err != nil {
		return nil, fmt.Errorf("failed to get membership %d: %w", membershipID, err)
	}

	return projectMembershipRowFromObject(membership, ""), nil
}

func listProjectMemberships(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	var projectKey string
	var projectIdentifier string

	if d.EqualsQuals["project_identifier"] != nil {
		projectIdentifier = d.EqualsQuals["project_identifier"].GetStringValue()
		projectKey = projectIdentifier
	} else if d.EqualsQuals["project_id"] != nil {
		projectKey = strconv.FormatInt(d.EqualsQuals["project_id"].GetInt64Value(), 10)
	} else {
		return nil, nil
	}

	var offset int64
	var pageSize int64 = 100
	if d.QueryContext.Limit != nil && *d.QueryContext.Limit < pageSize {
		pageSize = *d.QueryContext.Limit
	}

	for {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		result, _, err := client.MembershipMultiGet(projectKey, rm.MembershipMultiGetRequest{
			Offset: offset,
			Limit:  pageSize,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list memberships for project %s: %w", projectKey, err)
		}

		for _, membership := range result.Memberships {
			d.StreamListItem(ctx, projectMembershipRowFromObject(membership, projectIdentifier))

			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if int64(len(result.Memberships)) < pageSize {
			break
		}
		offset += pageSize
	}

	return nil, nil
}
