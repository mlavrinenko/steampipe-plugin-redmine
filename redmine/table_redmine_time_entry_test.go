package redmine

import (
	"testing"

	rm "github.com/nixys/nxs-go-redmine/v5"
)

func TestTimeEntryRowFromObject(t *testing.T) {
	t.Run("full object", func(t *testing.T) {
		obj := rm.TimeEntryObject{
			ID:        1,
			Project:   rm.IDName{ID: 10, Name: "Project"},
			Issue:     rm.TimeEntryIssueObject{ID: 100},
			User:      rm.IDName{ID: 5, Name: "Alice"},
			Activity:  rm.IDName{ID: 3, Name: "Development"},
			Hours:     2.5,
			Comments:  "Worked on feature",
			SpentOn:   "2026-02-15",
			CreatedOn: "2026-02-15T10:00:00Z",
			UpdatedOn: "2026-02-15T10:00:00Z",
		}

		row := timeEntryRowFromObject(obj)

		if row.ID != 1 {
			t.Errorf("ID = %d, want 1", row.ID)
		}
		if row.ProjectID != 10 {
			t.Errorf("ProjectID = %d, want 10", row.ProjectID)
		}
		if row.IssueID != 100 {
			t.Errorf("IssueID = %d, want 100", row.IssueID)
		}
		if row.Hours != 2.5 {
			t.Errorf("Hours = %f, want 2.5", row.Hours)
		}
		if row.Comments != "Worked on feature" {
			t.Errorf("Comments = %q, want %q", row.Comments, "Worked on feature")
		}
		if row.SpentOn == nil || row.SpentOn.Format("2006-01-02") != "2026-02-15" {
			t.Errorf("SpentOn = %v, want 2026-02-15", row.SpentOn)
		}
		if row.CreatedOn == nil {
			t.Error("CreatedOn should not be nil")
		}
		if row.Title != "2.50h on Project: Worked on feature" {
			t.Errorf("Title = %q, want %q", row.Title, "2.50h on Project: Worked on feature")
		}
	})

	t.Run("no issue", func(t *testing.T) {
		obj := rm.TimeEntryObject{
			ID:       2,
			Project:  rm.IDName{ID: 10, Name: "Project"},
			User:     rm.IDName{ID: 5, Name: "Alice"},
			Activity: rm.IDName{ID: 3, Name: "Development"},
			Hours:    1.0,
			SpentOn:  "2026-02-15",
		}

		row := timeEntryRowFromObject(obj)

		if row.IssueID != 0 {
			t.Errorf("IssueID = %d, want 0", row.IssueID)
		}
		if row.Title != "1.00h on Project" {
			t.Errorf("Title = %q, want %q", row.Title, "1.00h on Project")
		}
	})
}
