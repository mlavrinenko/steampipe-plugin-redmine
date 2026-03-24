package redmine

import (
	"testing"

	rm "github.com/nixys/nxs-go-redmine/v5"
)

func TestAttachmentRowFromObject(t *testing.T) {
	tests := []struct {
		name   string
		input  rm.AttachmentObject
		checks func(t *testing.T, row attachmentRow)
	}{
		{
			name: "full object",
			input: rm.AttachmentObject{
				ID:           42,
				FileName:     "screenshot.png",
				FileSize:     "128456",
				ContentType:  "image/png",
				Description:  "Bug reproduction screenshot",
				ContentURL:   "https://redmine.example.com/attachments/download/42/screenshot.png",
				ThumbnailURL: "https://redmine.example.com/attachments/thumbnail/42/200",
				Author:       rm.IDName{ID: 5, Name: "Alice"},
				CreatedOn:    "2026-02-15T10:30:00Z",
			},
			checks: func(t *testing.T, row attachmentRow) {
				if row.ID != 42 {
					t.Errorf("expected ID 42, got %d", row.ID)
				}
				if row.FileName != "screenshot.png" {
					t.Errorf("expected FileName screenshot.png, got %s", row.FileName)
				}
				if row.FileSize != "128456" {
					t.Errorf("expected FileSize 128456, got %s", row.FileSize)
				}
				if row.ContentType != "image/png" {
					t.Errorf("expected ContentType image/png, got %s", row.ContentType)
				}
				if row.Description != "Bug reproduction screenshot" {
					t.Errorf("expected Description to match, got %s", row.Description)
				}
				if row.AuthorID != 5 {
					t.Errorf("expected AuthorID 5, got %d", row.AuthorID)
				}
				if row.AuthorName != "Alice" {
					t.Errorf("expected AuthorName Alice, got %s", row.AuthorName)
				}
				if row.CreatedOn == nil {
					t.Error("expected CreatedOn to be set")
				}
				if row.ThumbnailURL != "https://redmine.example.com/attachments/thumbnail/42/200" {
					t.Errorf("unexpected ThumbnailURL: %s", row.ThumbnailURL)
				}
				if len(row.Akas) != 1 || row.Akas[0] != "/attachments/42" {
					t.Errorf("unexpected Akas: %v", row.Akas)
				}
			},
		},
		{
			name: "minimal object",
			input: rm.AttachmentObject{
				ID:       99,
				FileName: "data.csv",
				FileSize: "0",
				Author:   rm.IDName{ID: 1, Name: "Admin"},
			},
			checks: func(t *testing.T, row attachmentRow) {
				if row.ID != 99 {
					t.Errorf("expected ID 99, got %d", row.ID)
				}
				if row.FileName != "data.csv" {
					t.Errorf("expected FileName data.csv, got %s", row.FileName)
				}
				if row.ContentType != "" {
					t.Errorf("expected empty ContentType, got %s", row.ContentType)
				}
				if row.Description != "" {
					t.Errorf("expected empty Description, got %s", row.Description)
				}
				if row.ThumbnailURL != "" {
					t.Errorf("expected empty ThumbnailURL, got %s", row.ThumbnailURL)
				}
				if row.CreatedOn != nil {
					t.Error("expected nil CreatedOn for empty timestamp")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			row := attachmentRowFromObject(tt.input)
			tt.checks(t, row)
		})
	}
}
