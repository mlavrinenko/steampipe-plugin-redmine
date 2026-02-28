package redmine

import (
	"testing"

	rm "github.com/nixys/nxs-go-redmine/v5"
)

func TestIssueRowFromObject(t *testing.T) {
	t.Run("full object", func(t *testing.T) {
		hours := 8.5
		totalHours := 10.0
		startDate := "2026-02-01"
		dueDate := "2026-02-28"
		obj := rm.IssueObject{
			ID:                  100,
			Project:             rm.IDName{ID: 1, Name: "Project"},
			Tracker:             rm.IDName{ID: 2, Name: "Bug"},
			Status:              rm.IssueStatusObject{ID: 3, Name: "In Progress", IsClosed: false},
			Priority:            rm.IDName{ID: 4, Name: "High"},
			Author:              rm.IDName{ID: 5, Name: "Alice"},
			AssignedTo:          &rm.IDName{ID: 6, Name: "Bob"},
			Category:            &rm.IDName{ID: 7, Name: "Backend"},
			FixedVersion:        &rm.IDName{ID: 8, Name: "v2.0"},
			Parent:              &rm.IssueParentObject{ID: 50},
			Subject:             "Fix the bug",
			Description:         "Detailed description",
			StartDate:           &startDate,
			DueDate:             &dueDate,
			DoneRatio:           50,
			IsPrivate:           0,
			EstimatedHours:      &hours,
			TotalEstimatedHours: &totalHours,
			SpentHours:          3.5,
			TotalSpentHours:     5.0,
			CreatedOn:           "2026-02-01T00:00:00Z",
			UpdatedOn:           "2026-02-15T10:00:00Z",
			ClosedOn:            "",
		}

		row := issueRowFromObject(obj)

		if row.ID != 100 {
			t.Errorf("ID = %d, want 100", row.ID)
		}
		if row.ProjectID != 1 {
			t.Errorf("ProjectID = %d, want 1", row.ProjectID)
		}
		if row.StatusIsClosed {
			t.Error("StatusIsClosed should be false")
		}
		if row.AssignedToID != 6 {
			t.Errorf("AssignedToID = %d, want 6", row.AssignedToID)
		}
		if row.CategoryID != 7 {
			t.Errorf("CategoryID = %d, want 7", row.CategoryID)
		}
		if row.FixedVersionID != 8 {
			t.Errorf("FixedVersionID = %d, want 8", row.FixedVersionID)
		}
		if row.ParentID != 50 {
			t.Errorf("ParentID = %d, want 50", row.ParentID)
		}
		if row.StartDate == nil || *row.StartDate != "2026-02-01" {
			t.Errorf("StartDate = %v, want '2026-02-01'", row.StartDate)
		}
		if row.EstimatedHours == nil || *row.EstimatedHours != 8.5 {
			t.Errorf("EstimatedHours = %v, want 8.5", row.EstimatedHours)
		}
		if row.ClosedOn != nil {
			t.Errorf("ClosedOn = %v, want nil (empty string)", row.ClosedOn)
		}
		if row.CreatedOn == nil {
			t.Error("CreatedOn should not be nil")
		}
	})

	t.Run("nil pointer fields", func(t *testing.T) {
		obj := rm.IssueObject{
			ID:       200,
			Project:  rm.IDName{ID: 1, Name: "Project"},
			Tracker:  rm.IDName{ID: 2, Name: "Feature"},
			Status:   rm.IssueStatusObject{ID: 5, Name: "Closed", IsClosed: true},
			Priority: rm.IDName{ID: 3, Name: "Normal"},
			Author:   rm.IDName{ID: 10, Name: "Charlie"},
			Subject:  "Simple issue",
			ClosedOn: "2026-02-20T12:00:00Z",
		}

		row := issueRowFromObject(obj)

		if row.AssignedToID != 0 {
			t.Errorf("AssignedToID = %d, want 0", row.AssignedToID)
		}
		if row.CategoryID != 0 {
			t.Errorf("CategoryID = %d, want 0", row.CategoryID)
		}
		if row.ParentID != 0 {
			t.Errorf("ParentID = %d, want 0", row.ParentID)
		}
		if !row.StatusIsClosed {
			t.Error("StatusIsClosed should be true")
		}
		if row.ClosedOn == nil {
			t.Error("ClosedOn should not be nil")
		}
	})
}
