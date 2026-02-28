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
