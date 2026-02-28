package redmine

import (
	"testing"

	rm "github.com/nixys/nxs-go-redmine/v5"
)

func TestIssueStatusRowFromObject(t *testing.T) {
	t.Run("open status", func(t *testing.T) {
		obj := rm.IssueStatusObject{
			ID:       1,
			Name:     "New",
			IsClosed: false,
		}

		row := issueStatusRowFromObject(obj)

		if row.ID != 1 {
			t.Errorf("ID = %d, want 1", row.ID)
		}
		if row.Name != "New" {
			t.Errorf("Name = %q, want %q", row.Name, "New")
		}
		if row.IsClosed {
			t.Error("IsClosed should be false")
		}
	})

	t.Run("closed status", func(t *testing.T) {
		obj := rm.IssueStatusObject{
			ID:       5,
			Name:     "Closed",
			IsClosed: true,
		}

		row := issueStatusRowFromObject(obj)

		if row.ID != 5 {
			t.Errorf("ID = %d, want 5", row.ID)
		}
		if !row.IsClosed {
			t.Error("IsClosed should be true")
		}
	})
}
