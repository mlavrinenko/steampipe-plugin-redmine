package redmine

import (
	"context"
	"fmt"
	"strings"

	rm "github.com/nixys/nxs-go-redmine/v5"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func connect(ctx context.Context, d *plugin.QueryData) (*rm.Context, error) {
	cacheKey := "redmine_client"
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(*rm.Context), nil
	}

	config := GetConfig(d.Connection)

	if config.Endpoint == nil || *config.Endpoint == "" {
		return nil, fmt.Errorf("'endpoint' must be set in the connection configuration or REDMINE_ENDPOINT environment variable")
	}
	if config.APIKey == nil || *config.APIKey == "" {
		return nil, fmt.Errorf("'api_key' must be set in the connection configuration or REDMINE_API_KEY environment variable")
	}

	endpoint := strings.TrimRight(*config.Endpoint, "/")

	client := rm.Init(rm.Settings{
		Endpoint: endpoint,
		APIKey:   *config.APIKey,
	})

	d.ConnectionManager.Cache.Set(cacheKey, client)

	return client, nil
}
