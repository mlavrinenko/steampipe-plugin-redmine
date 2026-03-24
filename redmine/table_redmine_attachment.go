package redmine

import (
	"context"
	"fmt"
	"time"

	rm "github.com/nixys/nxs-go-redmine/v5"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

type attachmentRow struct {
	ID           int64
	FileName     string
	FileSize     string
	ContentType  string
	Description  string
	ContentURL   string
	ThumbnailURL string
	AuthorID     int64
	AuthorName   string
	CreatedOn    *time.Time
	Akas         []string
}

func tableRedmineAttachment() *plugin.Table {
	return &plugin.Table{
		Name:        "redmine_attachment",
		Description: "Attachments in the Redmine instance. Get-only: Redmine has no list-all-attachments endpoint.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getAttachment,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_INT, Description: "The attachment ID."},
			{Name: "author_id", Type: proto.ColumnType_INT, Description: "The author user ID."},
			{Name: "author_name", Type: proto.ColumnType_STRING, Description: "The author name."},
			{Name: "content_type", Type: proto.ColumnType_STRING, Description: "The MIME content type."},
			{Name: "content_url", Type: proto.ColumnType_STRING, Description: "The URL to download the attachment."},
			{Name: "created_on", Type: proto.ColumnType_TIMESTAMP, Description: "When the attachment was uploaded."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "The attachment description."},
			{Name: "file_name", Type: proto.ColumnType_STRING, Description: "The attachment file name."},
			{Name: "file_size", Type: proto.ColumnType_STRING, Description: "The file size in bytes."},
			{Name: "thumbnail_url", Type: proto.ColumnType_STRING, Description: "The thumbnail URL (for images)."},
			// Standard columns
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings for the resource."},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "The display name for this resource.", Transform: transform.FromField("FileName")},
		},
	}
}

//// HELPER FUNCTIONS

func attachmentRowFromObject(a rm.AttachmentObject) attachmentRow {
	return attachmentRow{
		ID:           a.ID,
		FileName:     a.FileName,
		FileSize:     a.FileSize,
		ContentType:  a.ContentType,
		Description:  a.Description,
		ContentURL:   a.ContentURL,
		ThumbnailURL: a.ThumbnailURL,
		AuthorID:     a.Author.ID,
		AuthorName:   a.Author.Name,
		CreatedOn:    parseRedmineTime(a.CreatedOn),
		Akas:         []string{fmt.Sprintf("/attachments/%d", a.ID)},
	}
}

//// HYDRATE FUNCTIONS

func getAttachment(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	attachmentID := d.EqualsQuals["id"].GetInt64Value()

	attachment, _, err := client.AttachmentSingleGet(attachmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get attachment %d: %w", attachmentID, err)
	}

	return attachmentRowFromObject(attachment), nil
}
