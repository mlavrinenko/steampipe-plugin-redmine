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

type documentCategoryRow struct {
	ID        int64
	Name      string
	IsDefault bool
	Active    bool
	Akas      []string
}

func tableRedmineDocumentCategory() *plugin.Table {
	return &plugin.Table{
		Name:        "redmine_document_category",
		Description: "Document category definitions in the Redmine instance.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getDocumentCategory,
		},
		List: &plugin.ListConfig{
			Hydrate: listDocumentCategories,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_INT, Description: "The document category ID."},
			{Name: "active", Type: proto.ColumnType_BOOL, Description: "Whether the category is active."},
			{Name: "is_default", Type: proto.ColumnType_BOOL, Description: "Whether this is the default category."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The category name."},
			// Standard columns
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings for the resource."},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "The display name for this resource.", Transform: transform.FromField("Name")},
		},
	}
}

//// HELPER FUNCTIONS

func documentCategoryRowFromObject(dc rm.EnumerationDocumentCategoryObject) documentCategoryRow {
	return documentCategoryRow{
		ID:        dc.ID,
		Name:      dc.Name,
		IsDefault: dc.IsDefault,
		Active:    dc.Active,
		Akas:      []string{fmt.Sprintf("/enumerations/document_categories/%d", dc.ID)},
	}
}

//// HYDRATE FUNCTIONS

func getAllDocumentCategories(ctx context.Context, d *plugin.QueryData) ([]rm.EnumerationDocumentCategoryObject, error) {
	cacheKey := "redmine_document_categories"
	if cached, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cached.([]rm.EnumerationDocumentCategoryObject), nil
	}

	client, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	categories, _, err := client.EnumerationDocumentCategoriesAllGet()
	if err != nil {
		return nil, fmt.Errorf("failed to list document categories: %w", err)
	}

	d.ConnectionManager.Cache.Set(cacheKey, categories)
	return categories, nil
}

func getDocumentCategory(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	categoryID := d.EqualsQuals["id"].GetInt64Value()

	categories, err := getAllDocumentCategories(ctx, d)
	if err != nil {
		return nil, err
	}

	for _, category := range categories {
		if category.ID == categoryID {
			return documentCategoryRowFromObject(category), nil
		}
	}

	return nil, nil
}

func listDocumentCategories(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	categories, err := getAllDocumentCategories(ctx, d)
	if err != nil {
		return nil, err
	}

	for _, category := range categories {
		d.StreamListItem(ctx, documentCategoryRowFromObject(category))

		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}
