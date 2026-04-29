package redmine

import (
	"testing"

	rm "github.com/nixys/nxs-go-redmine/v5"
)

func TestGroupMemberRowFromObject(t *testing.T) {
	tests := []struct {
		name   string
		group  rm.GroupObject
		user   rm.IDName
		checks func(t *testing.T, row groupMemberRow)
	}{
		{
			name:  "standard membership",
			group: rm.GroupObject{ID: 3, Name: "Developers"},
			user:  rm.IDName{ID: 42, Name: "Alice"},
			checks: func(t *testing.T, row groupMemberRow) {
				if row.GroupID != 3 {
					t.Errorf("GroupID = %d, want 3", row.GroupID)
				}
				if row.GroupName != "Developers" {
					t.Errorf("GroupName = %q, want %q", row.GroupName, "Developers")
				}
				if row.UserID != 42 {
					t.Errorf("UserID = %d, want 42", row.UserID)
				}
				if row.UserName != "Alice" {
					t.Errorf("UserName = %q, want %q", row.UserName, "Alice")
				}
				if len(row.Akas) != 1 || row.Akas[0] != "/groups/3/users/42" {
					t.Errorf("unexpected Akas: %v", row.Akas)
				}
			},
		},
		{
			name:  "zero ids",
			group: rm.GroupObject{Name: "Anonymous"},
			user:  rm.IDName{Name: "Nobody"},
			checks: func(t *testing.T, row groupMemberRow) {
				if row.GroupID != 0 || row.UserID != 0 {
					t.Errorf("expected zero IDs, got group=%d user=%d", row.GroupID, row.UserID)
				}
				if row.Akas[0] != "/groups/0/users/0" {
					t.Errorf("unexpected Akas: %v", row.Akas)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			row := groupMemberRowFromObject(tt.group, tt.user)
			tt.checks(t, row)
		})
	}
}

func TestGroupMemberRowFromObjectIterableShape(t *testing.T) {
	t.Run("group with users iterates each user", func(t *testing.T) {
		users := []rm.IDName{
			{ID: 1, Name: "Alice"},
			{ID: 2, Name: "Bob"},
		}
		group := rm.GroupObject{ID: 7, Name: "QA", Users: &users}

		var rows []groupMemberRow
		if group.Users != nil {
			for _, u := range *group.Users {
				rows = append(rows, groupMemberRowFromObject(group, u))
			}
		}

		if len(rows) != 2 {
			t.Fatalf("expected 2 rows, got %d", len(rows))
		}
		if rows[0].UserID != 1 || rows[1].UserID != 2 {
			t.Errorf("rows out of order: %+v", rows)
		}
	})

	t.Run("group with nil Users yields zero rows without panic", func(t *testing.T) {
		group := rm.GroupObject{ID: 7, Name: "Empty", Users: nil}

		var rows []groupMemberRow
		if group.Users != nil {
			for _, u := range *group.Users {
				rows = append(rows, groupMemberRowFromObject(group, u))
			}
		}

		if len(rows) != 0 {
			t.Errorf("expected 0 rows for nil Users, got %d", len(rows))
		}
	})

	t.Run("group with empty Users slice yields zero rows", func(t *testing.T) {
		empty := []rm.IDName{}
		group := rm.GroupObject{ID: 7, Name: "Empty", Users: &empty}

		var rows []groupMemberRow
		if group.Users != nil {
			for _, u := range *group.Users {
				rows = append(rows, groupMemberRowFromObject(group, u))
			}
		}

		if len(rows) != 0 {
			t.Errorf("expected 0 rows for empty Users, got %d", len(rows))
		}
	})
}
