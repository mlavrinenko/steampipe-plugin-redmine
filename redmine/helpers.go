package redmine

import (
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// parseRedmineTime parses a Redmine API timestamp string into *time.Time.
// Returns nil if parsing fails.
func parseRedmineTime(s string) *time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		t, err = time.Parse("2006-01-02T15:04:05Z", s)
		if err != nil {
			return nil
		}
	}
	return &t
}

// parseRedmineDate parses a Redmine API date string (YYYY-MM-DD) into *time.Time.
// Returns nil if the input is nil, empty, or unparseable.
func parseRedmineDate(s *string) *time.Time {
	if s == nil || *s == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", *s)
	if err != nil {
		return nil
	}
	return &t
}

// adjustTimestampBound converts a SQL comparison operator and timestamp into the
// normalized bound for a half-open [from, to) interval.
// Returns the adjusted time and whether it is a lower bound (true) or upper bound (false).
//
// The "=" operator is not handled here because it sets both bounds at once;
// callers must handle "=" by setting from=ts and to=ts+1s directly.
func adjustTimestampBound(operator string, ts time.Time) (bound time.Time, isFrom bool) {
	switch operator {
	case ">=":
		return ts, true
	case ">":
		return ts.Add(time.Second), true
	case "<=":
		return ts.Add(time.Second), false
	case "<":
		return ts, false
	default:
		return ts, true
	}
}

// adjustDateBound converts a SQL comparison operator and date into normalized
// inclusive from/to date strings (YYYY-MM-DD) suitable for the Redmine time entry API.
//
// The "=" operator is not handled here because it sets both bounds at once;
// callers must handle "=" by setting from=date and to=date directly.
func adjustDateBound(operator string, ts time.Time) (date string, isFrom bool) {
	switch operator {
	case ">=":
		return ts.Format("2006-01-02"), true
	case ">":
		return ts.AddDate(0, 0, 1).Format("2006-01-02"), true
	case "<=":
		return ts.Format("2006-01-02"), false
	case "<":
		return ts.AddDate(0, 0, -1).Format("2006-01-02"), false
	default:
		return ts.Format("2006-01-02"), true
	}
}

// dateRange represents a half-open time interval [from, to).
type dateRange struct {
	from *time.Time
	to   *time.Time
}

// extractDateRange parses timestamp qualifiers for the given column into a dateRange.
// Defaults to "created_on" if no column name is provided.
func extractDateRange(quals plugin.KeyColumnQualMap, column ...string) dateRange {
	col := "created_on"
	if len(column) > 0 {
		col = column[0]
	}

	var dr dateRange

	if quals[col] == nil {
		return dr
	}

	for _, q := range quals[col].Quals {
		ts := q.Value.GetTimestampValue().AsTime()

		if q.Operator == "=" {
			from := ts
			to := ts.Add(time.Second)
			dr.from = &from
			dr.to = &to
			continue
		}

		bound, isFrom := adjustTimestampBound(q.Operator, ts)
		if isFrom {
			t := bound
			dr.from = &t
		} else {
			t := bound
			dr.to = &t
		}
	}

	return dr
}

// timestampInRange checks if a parsed timestamp falls within the half-open [from, to) date range.
func timestampInRange(s string, dr dateRange) bool {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		t, err = time.Parse("2006-01-02T15:04:05Z", s)
		if err != nil {
			return false
		}
	}

	if dr.from != nil && t.Before(*dr.from) {
		return false
	}
	if dr.to != nil && !t.Before(*dr.to) {
		return false
	}

	return true
}

// buildTimestampFilter converts a dateRange into a Redmine API date filter string.
// Uses the >< (between) operator when both bounds exist, or >= / <= for single bounds.
func buildTimestampFilter(dr dateRange) string {
	layout := "2006-01-02T15:04:05Z"

	if dr.from != nil && dr.to != nil {
		return "><" + dr.from.Format(layout) + "|" + dr.to.Format(layout)
	}
	if dr.from != nil {
		return ">=" + dr.from.Format(layout)
	}
	if dr.to != nil {
		return "<=" + dr.to.Format(layout)
	}

	return ""
}

// extractSpentOnRange parses spent_on qualifiers into from/to date strings (YYYY-MM-DD)
// suitable for the Redmine time entry API's `from` and `to` parameters.
func extractSpentOnRange(quals plugin.KeyColumnQualMap) (from, to string) {
	if quals["spent_on"] == nil {
		return "", ""
	}

	for _, q := range quals["spent_on"].Quals {
		ts := q.Value.GetTimestampValue().AsTime()

		if q.Operator == "=" {
			d := ts.Format("2006-01-02")
			return d, d
		}

		date, isFrom := adjustDateBound(q.Operator, ts)
		if isFrom {
			from = date
		} else {
			to = date
		}
	}

	return from, to
}
