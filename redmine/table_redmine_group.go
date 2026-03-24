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

type groupRow struct {
	ID   int64
	Name string
	Akas []string
}

func tableRedmineGroup() *plugin.Table {
	return &plugin.Table{
		Name:        "redmine_group",
		Description: "Groups in the Redmine instance.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getGroup,
		},
		List: &plugin.ListConfig{
			Hydrate: listGroups,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_INT, Description: "The group ID."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The group name."},
			// Standard columns
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings for the resource."},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "The display name for this resource.", Transform: transform.FromField("Name")},
		},
	}
}

//// HELPER FUNCTIONS

func groupRowFromObject(g rm.GroupObject) groupRow {
	return groupRow{
		ID:   g.ID,
		Name: g.Name,
		Akas: []string{fmt.Sprintf("/groups/%d", g.ID)},
	}
}

//// HYDRATE FUNCTIONS

func getAllGroups(ctx context.Context, d *plugin.QueryData) ([]rm.GroupObject, error) {
	cacheKey := "redmine_groups"
	if cached, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cached.([]rm.GroupObject), nil
	}

	client, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	result, _, err := client.GroupAllGet()
	if err != nil {
		return nil, fmt.Errorf("failed to list groups: %w", err)
	}

	d.ConnectionManager.Cache.Set(cacheKey, result.Groups)
	return result.Groups, nil
}

func getGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	groupID := d.EqualsQuals["id"].GetInt64Value()

	groups, err := getAllGroups(ctx, d)
	if err != nil {
		return nil, err
	}

	for _, group := range groups {
		if group.ID == groupID {
			return groupRowFromObject(group), nil
		}
	}

	return nil, nil
}

func listGroups(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	groups, err := getAllGroups(ctx, d)
	if err != nil {
		return nil, err
	}

	for _, group := range groups {
		d.StreamListItem(ctx, groupRowFromObject(group))

		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}
