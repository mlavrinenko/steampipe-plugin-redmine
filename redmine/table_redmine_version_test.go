package redmine

import (
	"testing"

	rm "github.com/nixys/nxs-go-redmine/v5"
)

func TestVersionRowFromObject(t *testing.T) {
	dueDate := "2026-03-15"

	tests := []struct {
		name   string
		input  versionObject
		checks func(t *testing.T, row versionRow)
	}{
		{
			name: "full object",
			input: versionObject{
				ID:            10,
				Project:       rm.IDName{ID: 1, Name: "Test Project"},
				Name:          "v1.0",
				Description:   "First release",
				Status:        "open",
				DueDate:       &dueDate,
				Sharing:       "none",
				WikiPageTitle: "Release_v1",
				CreatedOn:     "2026-01-01T00:00:00Z",
				UpdatedOn:     "2026-02-15T12:00:00Z",
			},
			checks: func(t *testing.T, row versionRow) {
				if row.ID != 10 {
					t.Errorf("expected ID 10, got %d", row.ID)
				}
				if row.ProjectID != 1 {
					t.Errorf("expected ProjectID 1, got %d", row.ProjectID)
				}
				if row.ProjectName != "Test Project" {
					t.Errorf("expected ProjectName 'Test Project', got %s", row.ProjectName)
				}
				if row.Name != "v1.0" {
					t.Errorf("expected Name v1.0, got %s", row.Name)
				}
				if row.Status != "open" {
					t.Errorf("expected Status open, got %s", row.Status)
				}
				if row.Sharing != "none" {
					t.Errorf("expected Sharing none, got %s", row.Sharing)
				}
				if row.WikiPageTitle != "Release_v1" {
					t.Errorf("expected WikiPageTitle Release_v1, got %s", row.WikiPageTitle)
				}
				if row.DueDate == nil {
					t.Error("expected DueDate to be set")
				}
				if row.CreatedOn == nil {
					t.Error("expected CreatedOn to be set")
				}
			},
		},
		{
			name: "nil optional fields",
			input: versionObject{
				ID:      20,
				Project: rm.IDName{ID: 2, Name: "Other"},
				Name:    "v2.0",
				Status:  "closed",
			},
			checks: func(t *testing.T, row versionRow) {
				if row.ID != 20 {
					t.Errorf("expected ID 20, got %d", row.ID)
				}
				if row.DueDate != nil {
					t.Error("expected nil DueDate")
				}
				if row.Status != "closed" {
					t.Errorf("expected Status closed, got %s", row.Status)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			row := versionRowFromObject(tt.input)
			tt.checks(t, row)
		})
	}
}
