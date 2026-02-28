package redmine

import (
	"testing"

	rm "github.com/nixys/nxs-go-redmine/v5"
)

func TestTrackerRowFromObject(t *testing.T) {
	t.Run("full object", func(t *testing.T) {
		desc := "Bug tracker"
		obj := rm.TrackerObject{
			ID:                    1,
			Name:                  "Bug",
			DefaultStatus:         rm.IDName{ID: 1, Name: "New"},
			Description:           &desc,
			EnabledStandardFields: []string{"assigned_to_id", "category_id"},
		}

		row := trackerRowFromObject(obj)

		if row.ID != 1 {
			t.Errorf("ID = %d, want 1", row.ID)
		}
		if row.Name != "Bug" {
			t.Errorf("Name = %q, want %q", row.Name, "Bug")
		}
		if row.DefaultStatusID != 1 {
			t.Errorf("DefaultStatusID = %d, want 1", row.DefaultStatusID)
		}
		if row.Description == nil || *row.Description != "Bug tracker" {
			t.Errorf("Description = %v, want %q", row.Description, "Bug tracker")
		}
		if len(row.EnabledStandardFields) != 2 {
			t.Errorf("EnabledStandardFields length = %d, want 2", len(row.EnabledStandardFields))
		}
	})

	t.Run("nil optional fields", func(t *testing.T) {
		obj := rm.TrackerObject{
			ID:   2,
			Name: "Feature",
		}

		row := trackerRowFromObject(obj)

		if row.Description != nil {
			t.Errorf("Description = %v, want nil", row.Description)
		}
		if row.EnabledStandardFields != nil {
			t.Errorf("EnabledStandardFields = %v, want nil", row.EnabledStandardFields)
		}
	})
}
