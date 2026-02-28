package redmine

import (
	"testing"

	rm "github.com/nixys/nxs-go-redmine/v5"
)

func TestProjectRowFromObject(t *testing.T) {
	t.Run("full object", func(t *testing.T) {
		obj := rm.ProjectObject{
			ID:              1,
			Name:            "Test Project",
			Identifier:      "test-project",
			Description:     "A test project",
			Status:          rm.ProjectStatusActive,
			IsPublic:        true,
			Parent:          rm.IDName{ID: 10, Name: "Parent"},
			DefaultVersion:  &rm.IDName{ID: 5, Name: "v1.0"},
			DefaultAssignee: &rm.IDName{ID: 42, Name: "Alice"},
			CreatedOn:       "2026-01-01T00:00:00Z",
			UpdatedOn:       "2026-02-01T00:00:00Z",
		}

		row := projectRowFromObject(obj)

		if row.ID != 1 {
			t.Errorf("ID = %d, want 1", row.ID)
		}
		if row.Identifier != "test-project" {
			t.Errorf("Identifier = %q, want %q", row.Identifier, "test-project")
		}
		if row.Status != 1 {
			t.Errorf("Status = %d, want 1", row.Status)
		}
		if row.ParentID != 10 {
			t.Errorf("ParentID = %d, want 10", row.ParentID)
		}
		if row.DefaultVersionID != 5 {
			t.Errorf("DefaultVersionID = %d, want 5", row.DefaultVersionID)
		}
		if row.DefaultAssigneeID != 42 {
			t.Errorf("DefaultAssigneeID = %d, want 42", row.DefaultAssigneeID)
		}
		if row.CreatedOn == nil {
			t.Error("CreatedOn should not be nil")
		}
	})

	t.Run("nil pointer fields", func(t *testing.T) {
		obj := rm.ProjectObject{
			ID:         2,
			Name:       "Minimal",
			Identifier: "minimal",
			Status:     rm.ProjectStatusClosed,
		}

		row := projectRowFromObject(obj)

		if row.DefaultVersionID != 0 {
			t.Errorf("DefaultVersionID = %d, want 0", row.DefaultVersionID)
		}
		if row.DefaultAssigneeID != 0 {
			t.Errorf("DefaultAssigneeID = %d, want 0", row.DefaultAssigneeID)
		}
		if row.Status != 5 {
			t.Errorf("Status = %d, want 5", row.Status)
		}
	})
}
