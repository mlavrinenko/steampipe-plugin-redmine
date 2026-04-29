package redmine

import (
	"context"
	"fmt"

	rm "github.com/nixys/nxs-go-redmine/v5"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

type groupMemberRow struct {
	GroupID   int64
	GroupName string
	UserID    int64
	UserName  string
	Akas      []string
}

func tableRedmineGroupMember() *plugin.Table {
	return &plugin.Table{
		Name:        "redmine_group_member",
		Description: "Group memberships in Redmine. Each row is one user in a group. Listing requires group_id in the WHERE clause.",
		List: &plugin.ListConfig{
			Hydrate: listGroupMembers,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "group_id", Require: plugin.Required, Operators: []string{"="}},
			},
		},
		Columns: []*plugin.Column{
			// Key columns first
			{Name: "group_id", Type: proto.ColumnType_INT, Description: "The group ID."},
			// Remaining columns alphabetically
			{Name: "group_name", Type: proto.ColumnType_STRING, Description: "The group name."},
			{Name: "user_id", Type: proto.ColumnType_INT, Description: "The user ID."},
			{Name: "user_name", Type: proto.ColumnType_STRING, Description: "The user name."},
			// Standard columns
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings for the resource."},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "The display name for this resource.", Transform: transform.FromField("UserName")},
		},
	}
}

//// HELPER FUNCTIONS

func groupMemberRowFromObject(g rm.GroupObject, u rm.IDName) groupMemberRow {
	return groupMemberRow{
		GroupID:   g.ID,
		GroupName: g.Name,
		UserID:    u.ID,
		UserName:  u.Name,
		Akas:      []string{fmt.Sprintf("/groups/%d/users/%d", g.ID, u.ID)},
	}
}

//// HYDRATE FUNCTIONS

func getGroupWithUsers(ctx context.Context, d *plugin.QueryData, groupID int64) (*rm.GroupObject, error) {
	cacheKey := fmt.Sprintf("redmine_group_member:%d", groupID)
	if cached, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cached.(*rm.GroupObject), nil
	}

	client, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	group, _, err := client.GroupSingleGet(groupID, rm.GroupSingleGetRequest{
		Includes: []rm.GroupInclude{rm.GroupIncludeUsers},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get group %d: %w", groupID, err)
	}

	d.ConnectionManager.Cache.Set(cacheKey, &group)
	return &group, nil
}

func listGroupMembers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	groupID := d.EqualsQuals["group_id"].GetInt64Value()

	group, err := getGroupWithUsers(ctx, d, groupID)
	if err != nil {
		return nil, err
	}
	if group == nil || group.Users == nil {
		return nil, nil
	}

	for _, user := range *group.Users {
		d.StreamListItem(ctx, groupMemberRowFromObject(*group, user))

		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}
