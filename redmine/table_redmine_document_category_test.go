package redmine

import (
	"testing"

	rm "github.com/nixys/nxs-go-redmine/v5"
)

func TestDocumentCategoryRowFromObject(t *testing.T) {
	tests := []struct {
		name   string
		input  rm.EnumerationDocumentCategoryObject
		checks func(t *testing.T, row documentCategoryRow)
	}{
		{
			name: "default category",
			input: rm.EnumerationDocumentCategoryObject{
				ID:        6,
				Name:      "User Documentation",
				IsDefault: true,
				Active:    true,
			},
			checks: func(t *testing.T, row documentCategoryRow) {
				if row.ID != 6 {
					t.Errorf("expected ID 6, got %d", row.ID)
				}
				if row.Name != "User Documentation" {
					t.Errorf("expected Name User Documentation, got %s", row.Name)
				}
				if !row.IsDefault {
					t.Error("expected IsDefault true")
				}
				if !row.Active {
					t.Error("expected Active true")
				}
				if len(row.Akas) != 1 || row.Akas[0] != "/enumerations/document_categories/6" {
					t.Errorf("unexpected Akas: %v", row.Akas)
				}
			},
		},
		{
			name: "inactive category",
			input: rm.EnumerationDocumentCategoryObject{
				ID:        7,
				Name:      "Technical Documentation",
				IsDefault: false,
				Active:    false,
			},
			checks: func(t *testing.T, row documentCategoryRow) {
				if row.ID != 7 {
					t.Errorf("expected ID 7, got %d", row.ID)
				}
				if row.IsDefault {
					t.Error("expected IsDefault false")
				}
				if row.Active {
					t.Error("expected Active false")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			row := documentCategoryRowFromObject(tt.input)
			tt.checks(t, row)
		})
	}
}
