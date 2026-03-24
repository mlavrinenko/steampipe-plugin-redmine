package redmine

import (
	"testing"

	rm "github.com/nixys/nxs-go-redmine/v5"
)

func TestCustomFieldRowFromObject(t *testing.T) {
	defaultVal := "open"

	tests := []struct {
		name   string
		input  rm.CustomFieldObject
		checks func(t *testing.T, row customFieldRow)
	}{
		{
			name: "full list field",
			input: rm.CustomFieldObject{
				ID:             15,
				Name:           "Difficulty",
				CustomizedType: "issue",
				FieldFormat:    "list",
				IsRequired:     false,
				IsFilter:       true,
				Searchable:     false,
				Multiple:       false,
				DefaultValue:   &defaultVal,
				Visible:        false,
				Trackers:       []rm.IDName{{ID: 1, Name: "Bug"}, {ID: 2, Name: "Feature"}},
				PossibleValues: &[]rm.CustomFieldPossibleValueObject{
					{Value: "Easy", Label: ""},
					{Value: "Medium", Label: ""},
					{Value: "Hard", Label: ""},
				},
				Roles: []rm.IDName{{ID: 3, Name: "Manager"}},
			},
			checks: func(t *testing.T, row customFieldRow) {
				if row.ID != 15 {
					t.Errorf("expected ID 15, got %d", row.ID)
				}
				if row.Name != "Difficulty" {
					t.Errorf("expected Name Difficulty, got %s", row.Name)
				}
				if row.CustomizedType != "issue" {
					t.Errorf("expected CustomizedType issue, got %s", row.CustomizedType)
				}
				if row.FieldFormat != "list" {
					t.Errorf("expected FieldFormat list, got %s", row.FieldFormat)
				}
				if !row.IsFilter {
					t.Error("expected IsFilter true")
				}
				if row.IsRequired {
					t.Error("expected IsRequired false")
				}
				if row.DefaultValue == nil || *row.DefaultValue != "open" {
					t.Error("expected DefaultValue open")
				}
				if len(row.Trackers) != 2 {
					t.Errorf("expected 2 trackers, got %d", len(row.Trackers))
				}
				if row.PossibleValues == nil || len(*row.PossibleValues) != 3 {
					t.Error("expected 3 possible values")
				}
				if len(row.Roles) != 1 {
					t.Errorf("expected 1 role, got %d", len(row.Roles))
				}
				if len(row.Akas) != 1 || row.Akas[0] != "/custom_fields/15" {
					t.Errorf("unexpected Akas: %v", row.Akas)
				}
			},
		},
		{
			name: "minimal bool field",
			input: rm.CustomFieldObject{
				ID:             7,
				Name:           "Weekend Work",
				CustomizedType: "time_entry",
				FieldFormat:    "bool",
				IsRequired:     true,
				Visible:        true,
			},
			checks: func(t *testing.T, row customFieldRow) {
				if row.ID != 7 {
					t.Errorf("expected ID 7, got %d", row.ID)
				}
				if !row.IsRequired {
					t.Error("expected IsRequired true")
				}
				if !row.Visible {
					t.Error("expected Visible true")
				}
				if row.DefaultValue != nil {
					t.Error("expected nil DefaultValue")
				}
				if row.PossibleValues != nil {
					t.Error("expected nil PossibleValues")
				}
				if row.Trackers != nil {
					t.Error("expected nil Trackers")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			row := customFieldRowFromObject(tt.input)
			tt.checks(t, row)
		})
	}
}
