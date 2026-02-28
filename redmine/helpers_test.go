package redmine

import (
	"testing"
	"time"
)

func TestParseRedmineDate(t *testing.T) {
	date := "2026-02-15"
	empty := ""

	tests := map[string]struct {
		input    *string
		expected bool // whether parsing should succeed (non-nil result)
	}{
		"valid date": {input: &date, expected: true},
		"nil":        {input: nil, expected: false},
		"empty":      {input: &empty, expected: false},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := parseRedmineDate(tc.input)
			if (result != nil) != tc.expected {
				t.Errorf("parseRedmineDate(%v) = %v, want non-nil=%v", tc.input, result, tc.expected)
			}
			if result != nil && result.Format("2006-01-02") != "2026-02-15" {
				t.Errorf("parseRedmineDate(%v) = %v, want 2026-02-15", tc.input, result)
			}
		})
	}
}

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

func TestAdjustTimestampBound(t *testing.T) {
	ts := time.Date(2026, 2, 15, 10, 0, 0, 0, time.UTC)

	tests := map[string]struct {
		operator   string
		wantBound  time.Time
		wantIsFrom bool
	}{
		">= is lower bound, unchanged": {
			operator:   ">=",
			wantBound:  ts,
			wantIsFrom: true,
		},
		"> is lower bound, +1s": {
			operator:   ">",
			wantBound:  ts.Add(time.Second),
			wantIsFrom: true,
		},
		"<= is upper bound, +1s": {
			operator:   "<=",
			wantBound:  ts.Add(time.Second),
			wantIsFrom: false,
		},
		"< is upper bound, unchanged": {
			operator:   "<",
			wantBound:  ts,
			wantIsFrom: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			bound, isFrom := adjustTimestampBound(tc.operator, ts)
			if !bound.Equal(tc.wantBound) {
				t.Errorf("adjustTimestampBound(%q, %v) bound = %v, want %v", tc.operator, ts, bound, tc.wantBound)
			}
			if isFrom != tc.wantIsFrom {
				t.Errorf("adjustTimestampBound(%q, %v) isFrom = %v, want %v", tc.operator, ts, isFrom, tc.wantIsFrom)
			}
		})
	}
}

func TestAdjustDateBound(t *testing.T) {
	ts := time.Date(2026, 2, 15, 0, 0, 0, 0, time.UTC)

	tests := map[string]struct {
		operator   string
		wantDate   string
		wantIsFrom bool
	}{
		">= is lower bound, same date": {
			operator:   ">=",
			wantDate:   "2026-02-15",
			wantIsFrom: true,
		},
		"> is lower bound, next day": {
			operator:   ">",
			wantDate:   "2026-02-16",
			wantIsFrom: true,
		},
		"<= is upper bound, same date": {
			operator:   "<=",
			wantDate:   "2026-02-15",
			wantIsFrom: false,
		},
		"< is upper bound, previous day": {
			operator:   "<",
			wantDate:   "2026-02-14",
			wantIsFrom: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			date, isFrom := adjustDateBound(tc.operator, ts)
			if date != tc.wantDate {
				t.Errorf("adjustDateBound(%q, %v) date = %q, want %q", tc.operator, ts, date, tc.wantDate)
			}
			if isFrom != tc.wantIsFrom {
				t.Errorf("adjustDateBound(%q, %v) isFrom = %v, want %v", tc.operator, ts, isFrom, tc.wantIsFrom)
			}
		})
	}
}
