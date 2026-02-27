package redmine

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func shouldRetryError(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()

	// Retry on rate limiting
	if strings.Contains(msg, "429") {
		plugin.Logger(ctx).Debug("shouldRetryError", "rate_limit", err)
		return true
	}

	// Retry on service unavailable
	if strings.Contains(msg, "503") {
		plugin.Logger(ctx).Debug("shouldRetryError", "service_unavailable", err)
		return true
	}

	return false
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
