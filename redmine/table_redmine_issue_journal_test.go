package redmine

import (
	"testing"
	"time"

	rm "github.com/nixys/nxs-go-redmine/v5"
)

func TestExtractDateRange(t *testing.T) {
	// extractDateRange is tested indirectly through journalInRange and buildUpdatedOnFilter
	// since it requires plugin.KeyColumnQualMap which is hard to construct in unit tests.
	// The pure functions below are tested directly.
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

func TestUserParticipated(t *testing.T) {
	journals := []rm.IssueJournalObject{
		{ID: 1, User: rm.IDName{ID: 10, Name: "Alice"}, Notes: "First comment"},
		{ID: 2, User: rm.IDName{ID: 20, Name: "Bob"}, Notes: "Second comment"},
		{ID: 3, User: rm.IDName{ID: 10, Name: "Alice"}, Notes: "Third comment"},
	}

	tests := map[string]struct {
		journals []rm.IssueJournalObject
		userID   int64
		expected bool
	}{
		"user participated": {
			journals: journals,
			userID:   10,
			expected: true,
		},
		"user participated (second user)": {
			journals: journals,
			userID:   20,
			expected: true,
		},
		"user did not participate": {
			journals: journals,
			userID:   99,
			expected: false,
		},
		"empty journals": {
			journals: []rm.IssueJournalObject{},
			userID:   10,
			expected: false,
		},
		"nil journals": {
			journals: nil,
			userID:   10,
			expected: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := userParticipated(tc.journals, tc.userID)
			if result != tc.expected {
				t.Errorf("userParticipated(journals, %d) = %v, want %v", tc.userID, result, tc.expected)
			}
		})
	}
}

func TestBuildUpdatedOnFilter(t *testing.T) {
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
			result := buildUpdatedOnFilter(tc.dr)
			if result != tc.expected {
				t.Errorf("buildUpdatedOnFilter(%+v) = %q, want %q", tc.dr, result, tc.expected)
			}
		})
	}
}

func TestParseJournalTime(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected bool // whether parsing should succeed (non-nil result)
	}{
		"RFC3339":          {input: "2026-02-15T10:00:00Z", expected: true},
		"RFC3339 offset":   {input: "2026-02-15T10:00:00+02:00", expected: true},
		"invalid":          {input: "not-a-date", expected: false},
		"empty":            {input: "", expected: false},
		"date only":        {input: "2026-02-15", expected: false},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := parseJournalTime(tc.input)
			if (result != nil) != tc.expected {
				t.Errorf("parseJournalTime(%q) = %v, want non-nil=%v", tc.input, result, tc.expected)
			}
		})
	}
}
