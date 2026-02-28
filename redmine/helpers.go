package redmine

import "time"

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
