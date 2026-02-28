package redmine

import (
	"errors"
	"testing"
)

func TestIsRetryableError(t *testing.T) {
	tests := map[string]struct {
		err      error
		expected bool
	}{
		"429 rate limit": {
			err:      errors.New("unexpected status code has been returned (expected: 200, returned: 429, url: http://example.com, method: GET)"),
			expected: true,
		},
		"503 unavailable": {
			err:      errors.New("unexpected status code has been returned (expected: 200, returned: 503, url: http://example.com, method: GET)"),
			expected: true,
		},
		"500 server error": {
			err:      errors.New("unexpected status code has been returned (expected: 200, returned: 500, url: http://example.com, method: GET)"),
			expected: false,
		},
		"nil error": {
			err:      nil,
			expected: false,
		},
		"generic error": {
			err:      errors.New("connection refused"),
			expected: false,
		},
		"incidental 429 in text": {
			err:      errors.New("processed 429 records"),
			expected: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := isRetryableError(tc.err)
			if result != tc.expected {
				t.Errorf("isRetryableError(%v) = %v, want %v", tc.err, result, tc.expected)
			}
		})
	}
}

func TestIsNotFoundError(t *testing.T) {
	predicate := isNotFoundError([]string{"returned: 404", "returned: 422"})

	tests := map[string]struct {
		err      error
		expected bool
	}{
		"404 error": {
			err:      errors.New("unexpected status code has been returned (expected: 200, returned: 404, url: http://example.com, method: GET)"),
			expected: true,
		},
		"422 error": {
			err:      errors.New("unexpected status code has been returned (expected: 200, returned: 422, url: http://example.com, method: GET)"),
			expected: true,
		},
		"500 error": {
			err:      errors.New("unexpected status code has been returned (expected: 200, returned: 500, url: http://example.com, method: GET)"),
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
		"incidental 404 in text": {
			err:      errors.New("found 404 items in database"),
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
