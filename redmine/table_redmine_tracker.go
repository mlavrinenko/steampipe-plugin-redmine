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

type trackerRow struct {
	ID                    int64
	Name                  string
	DefaultStatusID       int64
	DefaultStatusName     string
	Description           *string
	EnabledStandardFields []string
}

func tableRedmineTracker() *plugin.Table {
	return &plugin.Table{
		Name:        "redmine_tracker",
		Description: "Trackers in the Redmine instance.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getTracker,
		},
		List: &plugin.ListConfig{
			Hydrate: listTrackers,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_INT, Description: "The tracker ID."},
			{Name: "default_status_id", Type: proto.ColumnType_INT, Description: "Default status ID for new issues."},
			{Name: "default_status_name", Type: proto.ColumnType_STRING, Description: "Default status name for new issues."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "The tracker description."},
			{Name: "enabled_standard_fields", Type: proto.ColumnType_JSON, Description: "Standard fields enabled for this tracker."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The tracker name."},
			// Standard columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: "The display name for this resource.", Transform: transform.FromField("Name")},
		},
	}
}

//// HELPER FUNCTIONS

func trackerRowFromObject(t rm.TrackerObject) trackerRow {
	return trackerRow{
		ID:                    t.ID,
		Name:                  t.Name,
		DefaultStatusID:       t.DefaultStatus.ID,
		DefaultStatusName:     t.DefaultStatus.Name,
		Description:           t.Description,
		EnabledStandardFields: t.EnabledStandardFields,
	}
}

//// HYDRATE FUNCTIONS

func getAllTrackers(ctx context.Context, d *plugin.QueryData) ([]rm.TrackerObject, error) {
	cacheKey := "redmine_trackers"
	if cached, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cached.([]rm.TrackerObject), nil
	}

	client, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	trackers, _, err := client.TrackerAllGet()
	if err != nil {
		return nil, fmt.Errorf("failed to list trackers: %w", err)
	}

	d.ConnectionManager.Cache.Set(cacheKey, trackers)
	return trackers, nil
}

func getTracker(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	trackerID := d.EqualsQuals["id"].GetInt64Value()

	trackers, err := getAllTrackers(ctx, d)
	if err != nil {
		return nil, err
	}

	for _, tracker := range trackers {
		if tracker.ID == trackerID {
			return trackerRowFromObject(tracker), nil
		}
	}

	return nil, nil
}

func listTrackers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	trackers, err := getAllTrackers(ctx, d)
	if err != nil {
		return nil, err
	}

	for _, tracker := range trackers {
		d.StreamListItem(ctx, trackerRowFromObject(tracker))

		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}
