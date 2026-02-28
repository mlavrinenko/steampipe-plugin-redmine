package redmine

import (
	"testing"

	rm "github.com/nixys/nxs-go-redmine/v5"
)

func TestIssuePriorityRowFromObject(t *testing.T) {
	tests := []struct {
		name   string
		input  rm.EnumerationPriorityObject
		checks func(t *testing.T, row issuePriorityRow)
	}{
		{
			name: "default priority",
			input: rm.EnumerationPriorityObject{
				ID:        2,
				Name:      "Normal",
				IsDefault: true,
				Active:    true,
			},
			checks: func(t *testing.T, row issuePriorityRow) {
				if row.ID != 2 {
					t.Errorf("expected ID 2, got %d", row.ID)
				}
				if row.Name != "Normal" {
					t.Errorf("expected Name Normal, got %s", row.Name)
				}
				if !row.IsDefault {
					t.Error("expected IsDefault true")
				}
				if !row.Active {
					t.Error("expected Active true")
				}
			},
		},
		{
			name: "inactive priority",
			input: rm.EnumerationPriorityObject{
				ID:        5,
				Name:      "Immediate",
				IsDefault: false,
				Active:    false,
			},
			checks: func(t *testing.T, row issuePriorityRow) {
				if row.ID != 5 {
					t.Errorf("expected ID 5, got %d", row.ID)
				}
				if row.IsDefault {
					t.Error("expected IsDefault false")
				}
				if row.Active {
					t.Error("expected Active false")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			row := issuePriorityRowFromObject(tt.input)
			tt.checks(t, row)
		})
	}
}
