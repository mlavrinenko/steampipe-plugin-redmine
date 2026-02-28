package redmine

import "testing"

func TestParseRedmineTime(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected bool // whether parsing should succeed (non-nil result)
	}{
		"RFC3339":        {input: "2026-02-15T10:00:00Z", expected: true},
		"RFC3339 offset": {input: "2026-02-15T10:00:00+02:00", expected: true},
		"invalid":        {input: "not-a-date", expected: false},
		"empty":          {input: "", expected: false},
		"date only":      {input: "2026-02-15", expected: false},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := parseRedmineTime(tc.input)
			if (result != nil) != tc.expected {
				t.Errorf("parseRedmineTime(%q) = %v, want non-nil=%v", tc.input, result, tc.expected)
			}
		})
	}
}
