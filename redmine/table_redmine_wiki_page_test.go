package redmine

import (
	"testing"

	rm "github.com/nixys/nxs-go-redmine/v5"
)

func TestWikiPageRowFromMultiObject(t *testing.T) {
	tests := []struct {
		name      string
		projectID string
		input     rm.WikiMultiObject
		checks    func(t *testing.T, row wikiPageRow)
	}{
		{
			name:      "full object with parent",
			projectID: "articles",
			input: rm.WikiMultiObject{
				Title:     "Getting_Started",
				Parent:    &rm.WikiParentObject{Title: "Wiki"},
				Version:   3,
				CreatedOn: "2026-01-15T10:00:00Z",
				UpdatedOn: "2026-02-20T14:30:00Z",
			},
			checks: func(t *testing.T, row wikiPageRow) {
				if row.ProjectID != "articles" {
					t.Errorf("expected ProjectID articles, got %s", row.ProjectID)
				}
				if row.Title != "Getting_Started" {
					t.Errorf("expected Title Getting_Started, got %s", row.Title)
				}
				if row.ParentTitle == nil || *row.ParentTitle != "Wiki" {
					t.Error("expected ParentTitle Wiki")
				}
				if row.Version != 3 {
					t.Errorf("expected Version 3, got %d", row.Version)
				}
				if row.CreatedOn == nil {
					t.Error("expected CreatedOn to be set")
				}
				if row.UpdatedOn == nil {
					t.Error("expected UpdatedOn to be set")
				}
				if row.Text != nil {
					t.Error("expected Text to be nil for multi object")
				}
				if row.AuthorID != nil {
					t.Error("expected AuthorID to be nil for multi object")
				}
				if len(row.Akas) != 1 || row.Akas[0] != "/projects/articles/wiki/Getting_Started" {
					t.Errorf("unexpected Akas: %v", row.Akas)
				}
			},
		},
		{
			name:      "no parent",
			projectID: "myproject",
			input: rm.WikiMultiObject{
				Title:     "Wiki",
				Version:   1,
				CreatedOn: "2026-01-01T00:00:00Z",
				UpdatedOn: "2026-01-01T00:00:00Z",
			},
			checks: func(t *testing.T, row wikiPageRow) {
				if row.ParentTitle != nil {
					t.Errorf("expected nil ParentTitle, got %v", row.ParentTitle)
				}
				if row.ProjectID != "myproject" {
					t.Errorf("expected ProjectID myproject, got %s", row.ProjectID)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			row := wikiPageRowFromMultiObject(tt.projectID, tt.input)
			tt.checks(t, row)
		})
	}
}

func TestWikiPageRowFromObject(t *testing.T) {
	tests := []struct {
		name      string
		projectID string
		input     rm.WikiObject
		checks    func(t *testing.T, row wikiPageRow)
	}{
		{
			name:      "full object",
			projectID: "articles",
			input: rm.WikiObject{
				Title:     "Getting_Started",
				Parent:    &rm.WikiParentObject{Title: "Wiki"},
				Text:      "h1. Getting Started\n\nWelcome!",
				Version:   5,
				Author:    rm.IDName{ID: 42, Name: "Alice"},
				Comments:  "Updated intro section",
				CreatedOn: "2026-01-15T10:00:00Z",
				UpdatedOn: "2026-03-01T09:00:00Z",
			},
			checks: func(t *testing.T, row wikiPageRow) {
				if row.ProjectID != "articles" {
					t.Errorf("expected ProjectID articles, got %s", row.ProjectID)
				}
				if row.Title != "Getting_Started" {
					t.Errorf("expected Title Getting_Started, got %s", row.Title)
				}
				if row.Text == nil || *row.Text != "h1. Getting Started\n\nWelcome!" {
					t.Error("expected Text to match")
				}
				if row.AuthorID == nil || *row.AuthorID != 42 {
					t.Error("expected AuthorID 42")
				}
				if row.AuthorName == nil || *row.AuthorName != "Alice" {
					t.Error("expected AuthorName Alice")
				}
				if row.Comments == nil || *row.Comments != "Updated intro section" {
					t.Error("expected Comments to match")
				}
				if row.ParentTitle == nil || *row.ParentTitle != "Wiki" {
					t.Error("expected ParentTitle Wiki")
				}
				if row.Version != 5 {
					t.Errorf("expected Version 5, got %d", row.Version)
				}
			},
		},
		{
			name:      "minimal object",
			projectID: "test",
			input: rm.WikiObject{
				Title:     "Wiki",
				Version:   1,
				CreatedOn: "2026-01-01T00:00:00Z",
				UpdatedOn: "2026-01-01T00:00:00Z",
			},
			checks: func(t *testing.T, row wikiPageRow) {
				if row.ParentTitle != nil {
					t.Error("expected nil ParentTitle")
				}
				if row.Text != nil {
					t.Error("expected nil Text")
				}
				if row.AuthorID != nil {
					t.Error("expected nil AuthorID")
				}
				if row.AuthorName != nil {
					t.Error("expected nil AuthorName")
				}
				if row.Comments != nil {
					t.Error("expected nil Comments")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			row := wikiPageRowFromObject(tt.projectID, tt.input)
			tt.checks(t, row)
		})
	}
}
