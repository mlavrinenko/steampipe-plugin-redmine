package redmine

import (
	"testing"

	rm "github.com/nixys/nxs-go-redmine/v5"
)

func TestMyAccountRowFromObject(t *testing.T) {
	status := rm.UserStatus(1)
	twofa := "totp"

	tests := []struct {
		name   string
		input  rm.UserObject
		checks func(t *testing.T, row myAccountRow)
	}{
		{
			name: "full object",
			input: rm.UserObject{
				ID:              42,
				Login:           "admin",
				Admin:           true,
				FirstName:       "Admin",
				LastName:        "User",
				Mail:            "admin@example.com",
				CreatedOn:       "2026-01-01T00:00:00Z",
				LastLoginOn:     "2026-02-28T12:00:00Z",
				PasswdChangedOn: "2026-01-15T00:00:00Z",
				TwofaScheme:     &twofa,
				Status:          &status,
			},
			checks: func(t *testing.T, row myAccountRow) {
				if row.ID != 42 {
					t.Errorf("expected ID 42, got %d", row.ID)
				}
				if row.Login != "admin" {
					t.Errorf("expected Login admin, got %s", row.Login)
				}
				if !row.Admin {
					t.Error("expected Admin true")
				}
				if row.Status != 1 {
					t.Errorf("expected Status 1, got %d", row.Status)
				}
				if row.TwofaScheme == nil || *row.TwofaScheme != "totp" {
					t.Error("expected TwofaScheme totp")
				}
			},
		},
		{
			name: "nil optional fields",
			input: rm.UserObject{
				ID:    1,
				Login: "user",
			},
			checks: func(t *testing.T, row myAccountRow) {
				if row.Status != 0 {
					t.Errorf("expected Status 0 for nil, got %d", row.Status)
				}
				if row.TwofaScheme != nil {
					t.Error("expected nil TwofaScheme")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			row := myAccountRowFromObject(tt.input)
			tt.checks(t, row)
		})
	}
}
