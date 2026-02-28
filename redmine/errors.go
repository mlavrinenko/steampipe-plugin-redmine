package redmine

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// isRetryableError returns true if the error message indicates a retryable condition
// (rate limiting or service unavailable).
// NOTE: relies on nxs-go-redmine formatting errors as "returned: <status>".
// If the library changes its error format, these predicates must be updated.
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "returned: 429") || strings.Contains(msg, "returned: 503")
}

func shouldRetryError(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, err error) bool {
	retry := isRetryableError(err)
	if retry {
		plugin.Logger(ctx).Debug("shouldRetryError", "retrying", err)
	}
	return retry
}

func retryConfig() *plugin.RetryConfig {
	return &plugin.RetryConfig{
		ShouldRetryErrorFunc: shouldRetryError,
		MaxAttempts:          5,
		BackoffAlgorithm:     "Exponential",
		RetryInterval:        1000,
		CappedDuration:       30000,
	}
}

func isNotFoundError(notFoundErrors []string) plugin.ErrorPredicate {
	return func(err error) bool {
		if err == nil {
			return false
		}
		msg := err.Error()
		for _, item := range notFoundErrors {
			if strings.Contains(msg, item) {
				return true
			}
		}
		return false
	}
}
