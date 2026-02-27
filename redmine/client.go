package redmine

import (
	"context"
	"fmt"

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
		return nil, fmt.Errorf("'endpoint' must be set in the connection configuration")
	}
	if config.APIKey == nil || *config.APIKey == "" {
		return nil, fmt.Errorf("'api_key' must be set in the connection configuration")
	}

	client := rm.Init(rm.Settings{
		Endpoint: *config.Endpoint,
		APIKey:   *config.APIKey,
	})

	d.ConnectionManager.Cache.Set(cacheKey, client)

	return client, nil
}

func getCurrentUserID(ctx context.Context, d *plugin.QueryData) (int64, error) {
	cacheKey := "current_user_id"
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(int64), nil
	}

	client, err := connect(ctx, d)
	if err != nil {
		return 0, err
	}

	user, _, err := client.UserCurrentGet(rm.UserCurrentGetRequest{})
	if err != nil {
		return 0, fmt.Errorf("failed to get current user: %w", err)
	}

	d.ConnectionManager.Cache.Set(cacheKey, user.ID)

	return user.ID, nil
}
