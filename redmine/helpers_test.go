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

func TestTimestampInRange(t *testing.T) {
	beforeRef := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)
	afterRef := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)

	tests := map[string]struct {
		timestamp string
		dr        dateRange
		expected  bool
	}{
		"within range": {
			timestamp: "2026-02-15T10:00:00Z",
			dr:        dateRange{from: &beforeRef, to: &afterRef},
			expected:  true,
		},
		"at range start (inclusive)": {
			timestamp: "2026-02-01T00:00:00Z",
			dr:        dateRange{from: &beforeRef, to: &afterRef},
			expected:  true,
		},
		"at range end (exclusive)": {
			timestamp: "2026-03-01T00:00:00Z",
			dr:        dateRange{from: &beforeRef, to: &afterRef},
			expected:  false,
		},
		"before range": {
			timestamp: "2026-01-31T23:59:59Z",
			dr:        dateRange{from: &beforeRef, to: &afterRef},
			expected:  false,
		},
		"after range": {
			timestamp: "2026-03-01T00:00:01Z",
			dr:        dateRange{from: &beforeRef, to: &afterRef},
			expected:  false,
		},
		"no lower bound": {
			timestamp: "2020-01-01T00:00:00Z",
			dr:        dateRange{from: nil, to: &afterRef},
			expected:  true,
		},
		"no upper bound": {
			timestamp: "2030-12-31T23:59:59Z",
			dr:        dateRange{from: &beforeRef, to: nil},
			expected:  true,
		},
		"no bounds": {
			timestamp: "2026-02-15T10:00:00Z",
			dr:        dateRange{from: nil, to: nil},
			expected:  true,
		},
		"invalid timestamp": {
			timestamp: "not-a-date",
			dr:        dateRange{from: nil, to: nil},
			expected:  false,
		},
		"RFC3339 with offset": {
			timestamp: "2026-02-15T10:00:00+00:00",
			dr:        dateRange{from: &beforeRef, to: &afterRef},
			expected:  true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := timestampInRange(tc.timestamp, tc.dr)
			if result != tc.expected {
				t.Errorf("timestampInRange(%q, %+v) = %v, want %v", tc.timestamp, tc.dr, result, tc.expected)
			}
		})
	}
}

func TestBuildTimestampFilter(t *testing.T) {
	from := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)

	tests := map[string]struct {
		dr       dateRange
		expected string
	}{
		"both bounds": {
			dr:       dateRange{from: &from, to: &to},
			expected: "><2026-02-01T00:00:00Z|2026-03-01T00:00:00Z",
		},
		"only from": {
			dr:       dateRange{from: &from, to: nil},
			expected: ">=2026-02-01T00:00:00Z",
		},
		"only to": {
			dr:       dateRange{from: nil, to: &to},
			expected: "<=2026-03-01T00:00:00Z",
		},
		"no bounds": {
			dr:       dateRange{from: nil, to: nil},
			expected: "",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := buildTimestampFilter(tc.dr)
			if result != tc.expected {
				t.Errorf("buildTimestampFilter(%+v) = %q, want %q", tc.dr, result, tc.expected)
			}
		})
	}
}
