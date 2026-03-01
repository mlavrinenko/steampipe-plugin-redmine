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

type timeEntryActivityRow struct {
	ID        int64
	Name      string
	IsDefault bool
	Active    bool
	Akas      []string
}

func tableRedmineTimeEntryActivity() *plugin.Table {
	return &plugin.Table{
		Name:        "redmine_time_entry_activity",
		Description: "Time entry activity definitions in the Redmine instance.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getTimeEntryActivity,
		},
		List: &plugin.ListConfig{
			Hydrate: listTimeEntryActivities,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_INT, Description: "The activity ID."},
			{Name: "active", Type: proto.ColumnType_BOOL, Description: "Whether the activity is active."},
			{Name: "is_default", Type: proto.ColumnType_BOOL, Description: "Whether this is the default activity."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The activity name."},
			// Standard columns
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings for the resource."},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "The display name for this resource.", Transform: transform.FromField("Name")},
		},
	}
}

//// HELPER FUNCTIONS

func timeEntryActivityRowFromObject(a rm.EnumerationTimeEntryActivityObject) timeEntryActivityRow {
	return timeEntryActivityRow{
		ID:        a.ID,
		Name:      a.Name,
		IsDefault: a.IsDefault,
		Active:    a.Active,
		Akas:      []string{fmt.Sprintf("/enumerations/time_entry_activities/%d", a.ID)},
	}
}

//// HYDRATE FUNCTIONS

func getAllTimeEntryActivities(ctx context.Context, d *plugin.QueryData) ([]rm.EnumerationTimeEntryActivityObject, error) {
	cacheKey := "redmine_time_entry_activities"
	if cached, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cached.([]rm.EnumerationTimeEntryActivityObject), nil
	}

	client, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	activities, _, err := client.EnumerationTimeEntryActivitiesAllGet()
	if err != nil {
		return nil, fmt.Errorf("failed to list time entry activities: %w", err)
	}

	d.ConnectionManager.Cache.Set(cacheKey, activities)
	return activities, nil
}

func getTimeEntryActivity(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	activityID := d.EqualsQuals["id"].GetInt64Value()

	activities, err := getAllTimeEntryActivities(ctx, d)
	if err != nil {
		return nil, err
	}

	for _, activity := range activities {
		if activity.ID == activityID {
			return timeEntryActivityRowFromObject(activity), nil
		}
	}

	return nil, nil
}

func listTimeEntryActivities(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	activities, err := getAllTimeEntryActivities(ctx, d)
	if err != nil {
		return nil, err
	}

	for _, activity := range activities {
		d.StreamListItem(ctx, timeEntryActivityRowFromObject(activity))

		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}
