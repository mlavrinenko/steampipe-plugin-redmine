package redmine

import (
	"errors"
	"testing"
)

func TestIsNotFoundError(t *testing.T) {
	predicate := isNotFoundError([]string{"404", "NotFound"})

	tests := map[string]struct {
		err      error
		expected bool
	}{
		"404 error": {
			err:      errors.New("unexpected status code: 404"),
			expected: true,
		},
		"NotFound error": {
			err:      errors.New("resource NotFound"),
			expected: true,
		},
		"500 error": {
			err:      errors.New("unexpected status code: 500"),
			expected: false,
		},
		"nil error": {
			err:      nil,
			expected: false,
		},
		"empty error": {
			err:      errors.New(""),
			expected: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := predicate(tc.err)
			if result != tc.expected {
				t.Errorf("isNotFoundError(%v) = %v, want %v", tc.err, result, tc.expected)
			}
		})
	}
}
