package redmine

import (
	"testing"
	"time"
)

func TestExtractDateRange(t *testing.T) {
	// extractDateRange requires plugin.KeyColumnQualMap which is hard to construct in unit tests.
	// The operator logic is extracted into adjustTimestampBound and tested in helpers_test.go.
}

func TestJournalInRange(t *testing.T) {
	refTime := time.Date(2026, 2, 15, 10, 0, 0, 0, time.UTC)
	beforeRef := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)
	afterRef := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)

	tests := map[string]struct {
		journalTime string
		dr          dateRange
		expected    bool
	}{
		"within range": {
			journalTime: "2026-02-15T10:00:00Z",
			dr:          dateRange{from: &beforeRef, to: &afterRef},
			expected:    true,
		},
		"at range start": {
			journalTime: "2026-02-01T00:00:00Z",
			dr:          dateRange{from: &beforeRef, to: &afterRef},
			expected:    true,
		},
		"at range end (exclusive)": {
			journalTime: "2026-03-01T00:00:00Z",
			dr:          dateRange{from: &beforeRef, to: &afterRef},
			expected:    false,
		},
		"before range": {
			journalTime: "2026-01-31T23:59:59Z",
			dr:          dateRange{from: &beforeRef, to: &afterRef},
			expected:    false,
		},
		"after range": {
			journalTime: "2026-03-01T00:00:01Z",
			dr:          dateRange{from: &beforeRef, to: &afterRef},
			expected:    false,
		},
		"no lower bound": {
			journalTime: "2020-01-01T00:00:00Z",
			dr:          dateRange{from: nil, to: &afterRef},
			expected:    true,
		},
		"no upper bound": {
			journalTime: "2030-12-31T23:59:59Z",
			dr:          dateRange{from: &beforeRef, to: nil},
			expected:    true,
		},
		"no bounds": {
			journalTime: "2026-02-15T10:00:00Z",
			dr:          dateRange{from: nil, to: nil},
			expected:    true,
		},
		"invalid timestamp": {
			journalTime: "not-a-date",
			dr:          dateRange{from: nil, to: nil},
			expected:    false,
		},
		"RFC3339 with offset": {
			journalTime: refTime.Format(time.RFC3339),
			dr:          dateRange{from: &beforeRef, to: &afterRef},
			expected:    true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := journalInRange(tc.journalTime, tc.dr)
			if result != tc.expected {
				t.Errorf("journalInRange(%q, %+v) = %v, want %v", tc.journalTime, tc.dr, result, tc.expected)
			}
		})
	}
}

func TestBuildDateFilter(t *testing.T) {
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
			result := buildDateFilter(tc.dr)
			if result != tc.expected {
				t.Errorf("buildDateFilter(%+v) = %q, want %q", tc.dr, result, tc.expected)
			}
		})
	}
}
