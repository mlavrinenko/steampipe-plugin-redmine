package redmine

import (
	"testing"

	rm "github.com/nixys/nxs-go-redmine/v5"
)

func TestProjectMembershipRowFromObject(t *testing.T) {
	t.Run("user membership", func(t *testing.T) {
		user := rm.IDName{ID: 10, Name: "Alice"}
		obj := rm.MembershipObject{
			ID:      1,
			Project: rm.IDName{ID: 5, Name: "My Project"},
			User:    &user,
			Roles:   []rm.MembershipRoleObject{{ID: 3, Name: "Developer", Inherited: false}},
		}

		row := projectMembershipRowFromObject(obj, "my-project")

		if row.ID != 1 {
			t.Errorf("ID = %d, want 1", row.ID)
		}
		if row.ProjectID != 5 {
			t.Errorf("ProjectID = %d, want 5", row.ProjectID)
		}
		if row.ProjectName != "My Project" {
			t.Errorf("ProjectName = %q, want %q", row.ProjectName, "My Project")
		}
		if row.ProjectIdentifier != "my-project" {
			t.Errorf("ProjectIdentifier = %q, want %q", row.ProjectIdentifier, "my-project")
		}
		if row.UserID != 10 {
			t.Errorf("UserID = %d, want 10", row.UserID)
		}
		if row.UserName != "Alice" {
			t.Errorf("UserName = %q, want %q", row.UserName, "Alice")
		}
		if row.GroupID != 0 {
			t.Errorf("GroupID = %d, want 0", row.GroupID)
		}
		if len(row.Roles) != 1 || row.Roles[0].ID != 3 {
			t.Errorf("Roles = %v, want [{ID:3 Name:Developer}]", row.Roles)
		}
	})

	t.Run("group membership", func(t *testing.T) {
		group := rm.IDName{ID: 20, Name: "Managers"}
		obj := rm.MembershipObject{
			ID:      2,
			Project: rm.IDName{ID: 5, Name: "My Project"},
			Group:   &group,
			Roles:   []rm.MembershipRoleObject{{ID: 4, Name: "Manager", Inherited: false}},
		}

		row := projectMembershipRowFromObject(obj, "")

		if row.UserID != 0 {
			t.Errorf("UserID = %d, want 0", row.UserID)
		}
		if row.GroupID != 20 {
			t.Errorf("GroupID = %d, want 20", row.GroupID)
		}
		if row.GroupName != "Managers" {
			t.Errorf("GroupName = %q, want %q", row.GroupName, "Managers")
		}
		if row.ProjectIdentifier != "" {
			t.Errorf("ProjectIdentifier = %q, want empty", row.ProjectIdentifier)
		}
	})

	t.Run("nil user and group", func(t *testing.T) {
		obj := rm.MembershipObject{
			ID:      3,
			Project: rm.IDName{ID: 5, Name: "My Project"},
			Roles:   []rm.MembershipRoleObject{},
		}

		row := projectMembershipRowFromObject(obj, "")

		if row.UserID != 0 {
			t.Errorf("UserID = %d, want 0", row.UserID)
		}
		if row.GroupID != 0 {
			t.Errorf("GroupID = %d, want 0", row.GroupID)
		}
	})
}
