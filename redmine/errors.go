package redmine

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// isRetryableError returns true if the error message indicates a retryable condition
// (rate limiting or service unavailable).
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "429") || strings.Contains(msg, "503")
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
		if err != nil {
			for _, item := range notFoundErrors {
				if strings.Contains(err.Error(), item) {
					return true
				}
			}
		}
		return false
	}
}
