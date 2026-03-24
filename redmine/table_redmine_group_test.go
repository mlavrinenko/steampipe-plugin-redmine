package redmine

import (
	"testing"

	rm "github.com/nixys/nxs-go-redmine/v5"
)

func TestGroupRowFromObject(t *testing.T) {
	tests := []struct {
		name   string
		input  rm.GroupObject
		checks func(t *testing.T, row groupRow)
	}{
		{
			name: "standard group",
			input: rm.GroupObject{
				ID:   3,
				Name: "Developers",
			},
			checks: func(t *testing.T, row groupRow) {
				if row.ID != 3 {
					t.Errorf("expected ID 3, got %d", row.ID)
				}
				if row.Name != "Developers" {
					t.Errorf("expected Name Developers, got %s", row.Name)
				}
				if len(row.Akas) != 1 || row.Akas[0] != "/groups/3" {
					t.Errorf("unexpected Akas: %v", row.Akas)
				}
			},
		},
		{
			name: "group with zero ID",
			input: rm.GroupObject{
				Name: "Anonymous",
			},
			checks: func(t *testing.T, row groupRow) {
				if row.ID != 0 {
					t.Errorf("expected ID 0, got %d", row.ID)
				}
				if row.Akas[0] != "/groups/0" {
					t.Errorf("unexpected Akas: %v", row.Akas)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			row := groupRowFromObject(tt.input)
			tt.checks(t, row)
		})
	}
}
