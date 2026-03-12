package redmine

import (
	"testing"
)

func TestSearchRowFromObject(t *testing.T) {
	t.Run("full object", func(t *testing.T) {
		obj := searchResultObject{
			ID:          42,
			Title:       "Fix login redirect",
			Type:        "issue",
			URL:         "/issues/42",
			Description: "Users are redirected to the wrong page after login...",
			Datetime:    "2026-03-01T14:30:00Z",
		}

		row := searchRowFromObject(obj)

		if row.ID != 42 {
			t.Errorf("ID = %d, want 42", row.ID)
		}
		if row.Title != "Fix login redirect" {
			t.Errorf("Title = %q, want %q", row.Title, "Fix login redirect")
		}
		if row.Type != "issue" {
			t.Errorf("Type = %q, want %q", row.Type, "issue")
		}
		if row.URL != "/issues/42" {
			t.Errorf("URL = %q, want %q", row.URL, "/issues/42")
		}
		if row.Description != "Users are redirected to the wrong page after login..." {
			t.Errorf("Description = %q, want %q", row.Description, "Users are redirected to the wrong page after login...")
		}
		if row.Datetime == nil {
			t.Fatal("Datetime should not be nil")
		}
		if row.Datetime.Format("2006-01-02") != "2026-03-01" {
			t.Errorf("Datetime = %v, want 2026-03-01", row.Datetime)
		}
		if len(row.Akas) != 1 || row.Akas[0] != "/issues/42" {
			t.Errorf("Akas = %v, want [/issues/42]", row.Akas)
		}
	})

	t.Run("minimal object", func(t *testing.T) {
		obj := searchResultObject{
			ID:    7,
			Title: "Architecture",
			Type:  "wiki-page",
			URL:   "/projects/myproject/wiki/Architecture",
		}

		row := searchRowFromObject(obj)

		if row.ID != 7 {
			t.Errorf("ID = %d, want 7", row.ID)
		}
		if row.Type != "wiki-page" {
			t.Errorf("Type = %q, want %q", row.Type, "wiki-page")
		}
		if row.Description != "" {
			t.Errorf("Description = %q, want empty", row.Description)
		}
		if row.Datetime != nil {
			t.Errorf("Datetime = %v, want nil", row.Datetime)
		}
		if len(row.Akas) != 1 || row.Akas[0] != "/projects/myproject/wiki/Architecture" {
			t.Errorf("Akas = %v, want [/projects/myproject/wiki/Architecture]", row.Akas)
		}
	})
}
