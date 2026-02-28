package redmine

import (
	"testing"

	rm "github.com/nixys/nxs-go-redmine/v5"
)

func TestTimeEntryActivityRowFromObject(t *testing.T) {
	tests := []struct {
		name   string
		input  rm.EnumerationTimeEntryActivityObject
		checks func(t *testing.T, row timeEntryActivityRow)
	}{
		{
			name: "default activity",
			input: rm.EnumerationTimeEntryActivityObject{
				ID:        9,
				Name:      "Development",
				IsDefault: true,
				Active:    true,
			},
			checks: func(t *testing.T, row timeEntryActivityRow) {
				if row.ID != 9 {
					t.Errorf("expected ID 9, got %d", row.ID)
				}
				if row.Name != "Development" {
					t.Errorf("expected Name Development, got %s", row.Name)
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
			name: "inactive activity",
			input: rm.EnumerationTimeEntryActivityObject{
				ID:     10,
				Name:   "Deprecated",
				Active: false,
			},
			checks: func(t *testing.T, row timeEntryActivityRow) {
				if row.ID != 10 {
					t.Errorf("expected ID 10, got %d", row.ID)
				}
				if row.Active {
					t.Error("expected Active false")
				}
				if row.IsDefault {
					t.Error("expected IsDefault false")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			row := timeEntryActivityRowFromObject(tt.input)
			tt.checks(t, row)
		})
	}
}
