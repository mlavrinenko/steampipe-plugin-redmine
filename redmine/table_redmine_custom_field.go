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

type customFieldRow struct {
	ID             int64
	Name           string
	CustomizedType string
	FieldFormat    string
	Regexp         string
	MinLength      int64
	MaxLength      int64
	IsRequired     bool
	IsFilter       bool
	Searchable     bool
	Multiple       bool
	DefaultValue   *string
	Visible        bool
	Trackers       []rm.IDName
	PossibleValues *[]rm.CustomFieldPossibleValueObject
	Roles          []rm.IDName
	Akas           []string
}

func tableRedmineCustomField() *plugin.Table {
	return &plugin.Table{
		Name:        "redmine_custom_field",
		Description: "Custom field definitions in the Redmine instance.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getCustomField,
		},
		List: &plugin.ListConfig{
			Hydrate: listCustomFields,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_INT, Description: "The custom field ID."},
			{Name: "customized_type", Type: proto.ColumnType_STRING, Description: "The type of object this field applies to (issue, project, user, etc.)."},
			{Name: "default_value", Type: proto.ColumnType_STRING, Description: "The default value."},
			{Name: "field_format", Type: proto.ColumnType_STRING, Description: "The field format (string, int, list, date, etc.)."},
			{Name: "is_filter", Type: proto.ColumnType_BOOL, Description: "Whether the field can be used as a filter."},
			{Name: "is_required", Type: proto.ColumnType_BOOL, Description: "Whether the field is required."},
			{Name: "max_length", Type: proto.ColumnType_INT, Description: "Maximum length for the field value."},
			{Name: "min_length", Type: proto.ColumnType_INT, Description: "Minimum length for the field value."},
			{Name: "multiple", Type: proto.ColumnType_BOOL, Description: "Whether the field allows multiple values."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The custom field name."},
			{Name: "possible_values", Type: proto.ColumnType_JSON, Description: "Possible values for list-type fields."},
			{Name: "regexp", Type: proto.ColumnType_STRING, Description: "Validation regexp pattern."},
			{Name: "roles", Type: proto.ColumnType_JSON, Description: "Roles that can see the field."},
			{Name: "searchable", Type: proto.ColumnType_BOOL, Description: "Whether the field is searchable."},
			{Name: "trackers", Type: proto.ColumnType_JSON, Description: "Trackers the field applies to."},
			{Name: "visible", Type: proto.ColumnType_BOOL, Description: "Whether the field is visible to all users."},
			// Standard columns
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings for the resource."},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "The display name for this resource.", Transform: transform.FromField("Name")},
		},
	}
}

//// HELPER FUNCTIONS

func customFieldRowFromObject(cf rm.CustomFieldObject) customFieldRow {
	return customFieldRow{
		ID:             cf.ID,
		Name:           cf.Name,
		CustomizedType: cf.CustomizedType,
		FieldFormat:    cf.FieldFormat,
		Regexp:         cf.Regexp,
		MinLength:      cf.MinLength,
		MaxLength:      cf.MaxLength,
		IsRequired:     cf.IsRequired,
		IsFilter:       cf.IsFilter,
		Searchable:     cf.Searchable,
		Multiple:       cf.Multiple,
		DefaultValue:   cf.DefaultValue,
		Visible:        cf.Visible,
		Trackers:       cf.Trackers,
		PossibleValues: cf.PossibleValues,
		Roles:          cf.Roles,
		Akas:           []string{fmt.Sprintf("/custom_fields/%d", cf.ID)},
	}
}

//// HYDRATE FUNCTIONS

func getAllCustomFields(ctx context.Context, d *plugin.QueryData) ([]rm.CustomFieldObject, error) {
	cacheKey := "redmine_custom_fields"
	if cached, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cached.([]rm.CustomFieldObject), nil
	}

	client, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	fields, _, err := client.CustomFieldAllGet()
	if err != nil {
		return nil, fmt.Errorf("failed to list custom fields: %w", err)
	}

	d.ConnectionManager.Cache.Set(cacheKey, fields)
	return fields, nil
}

func getCustomField(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	fieldID := d.EqualsQuals["id"].GetInt64Value()

	fields, err := getAllCustomFields(ctx, d)
	if err != nil {
		return nil, err
	}

	for _, field := range fields {
		if field.ID == fieldID {
			return customFieldRowFromObject(field), nil
		}
	}

	return nil, nil
}

func listCustomFields(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	fields, err := getAllCustomFields(ctx, d)
	if err != nil {
		return nil, err
	}

	for _, field := range fields {
		d.StreamListItem(ctx, customFieldRowFromObject(field))

		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}
