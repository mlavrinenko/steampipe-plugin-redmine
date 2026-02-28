package redmine

import (
	"testing"

	rm "github.com/nixys/nxs-go-redmine/v5"
)

func TestUserRowFromObject(t *testing.T) {
	t.Run("full object", func(t *testing.T) {
		active := rm.UserStatusActive
		scheme := "totp"
		obj := rm.UserObject{
			ID:              1,
			Login:           "alice",
			Admin:           true,
			FirstName:       "Alice",
			LastName:        "Smith",
			Mail:            "alice@example.com",
			CreatedOn:       "2026-01-01T00:00:00Z",
			LastLoginOn:     "2026-02-15T10:00:00Z",
			PasswdChangedOn: "2026-01-15T00:00:00Z",
			TwofaScheme:     &scheme,
			Status:          &active,
			Groups:          &[]rm.IDName{{ID: 1, Name: "Developers"}},
		}

		row := userRowFromObject(obj)

		if row.ID != 1 {
			t.Errorf("ID = %d, want 1", row.ID)
		}
		if row.Mail != "alice@example.com" {
			t.Errorf("Mail = %q, want %q", row.Mail, "alice@example.com")
		}
		if row.Status != 1 {
			t.Errorf("Status = %d, want 1", row.Status)
		}
		if row.TwofaScheme == nil || *row.TwofaScheme != "totp" {
			t.Errorf("TwofaScheme = %v, want 'totp'", row.TwofaScheme)
		}
		if row.Groups == nil || len(*row.Groups) != 1 {
			t.Error("Groups should have 1 entry")
		}
		if row.CreatedOn == nil {
			t.Error("CreatedOn should not be nil")
		}
	})

	t.Run("nil status and optional fields", func(t *testing.T) {
		obj := rm.UserObject{
			ID:    2,
			Login: "bob",
			Mail:  "bob@example.com",
		}

		row := userRowFromObject(obj)

		if row.Status != 0 {
			t.Errorf("Status = %d, want 0", row.Status)
		}
		if row.TwofaScheme != nil {
			t.Errorf("TwofaScheme = %v, want nil", row.TwofaScheme)
		}
		if row.Groups != nil {
			t.Errorf("Groups = %v, want nil", row.Groups)
		}
		if row.Memberships != nil {
			t.Errorf("Memberships = %v, want nil", row.Memberships)
		}
	})
}
