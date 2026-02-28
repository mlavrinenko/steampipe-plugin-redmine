package redmine

import (
	"testing"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func TestGetConfig_NilConnection(t *testing.T) {
	config := GetConfig(nil)
	if config == nil {
		t.Fatal("expected non-nil config, got nil")
	}
	if config.Endpoint != nil {
		t.Errorf("expected nil Endpoint, got %v", config.Endpoint)
	}
	if config.APIKey != nil {
		t.Errorf("expected nil APIKey, got %v", config.APIKey)
	}
}

func TestGetConfig_NilConfig(t *testing.T) {
	conn := &plugin.Connection{}
	config := GetConfig(conn)
	if config == nil {
		t.Fatal("expected non-nil config, got nil")
	}
	if config.Endpoint != nil {
		t.Errorf("expected nil Endpoint, got %v", config.Endpoint)
	}
}

func TestGetConfig_PointerType(t *testing.T) {
	endpoint := "https://example.com"
	apiKey := "abc123"
	conn := &plugin.Connection{
		Config: &redmineConfig{
			Endpoint: &endpoint,
			APIKey:   &apiKey,
		},
	}

	config := GetConfig(conn)

	if config.Endpoint == nil || *config.Endpoint != "https://example.com" {
		t.Errorf("expected Endpoint 'https://example.com', got %v", config.Endpoint)
	}
	if config.APIKey == nil || *config.APIKey != "abc123" {
		t.Errorf("expected APIKey 'abc123', got %v", config.APIKey)
	}
}

func TestGetConfig_ValueType(t *testing.T) {
	// The SDK may store config as a value (not pointer) after HCL parsing.
	endpoint := "https://example.com"
	apiKey := "abc123"
	conn := &plugin.Connection{
		Config: redmineConfig{
			Endpoint: &endpoint,
			APIKey:   &apiKey,
		},
	}

	config := GetConfig(conn)

	if config.Endpoint == nil || *config.Endpoint != "https://example.com" {
		t.Errorf("expected Endpoint 'https://example.com', got %v", config.Endpoint)
	}
	if config.APIKey == nil || *config.APIKey != "abc123" {
		t.Errorf("expected APIKey 'abc123', got %v", config.APIKey)
	}
}

func TestGetConfig_UnknownType(t *testing.T) {
	conn := &plugin.Connection{
		Config: "not a redmineConfig",
	}

	config := GetConfig(conn)

	if config.Endpoint != nil {
		t.Errorf("expected nil Endpoint for unknown type, got %v", config.Endpoint)
	}
}

func TestConfigInstance(t *testing.T) {
	instance := ConfigInstance()
	if instance == nil {
		t.Fatal("ConfigInstance() returned nil")
	}
	config, ok := instance.(*redmineConfig)
	if !ok {
		t.Fatalf("ConfigInstance() returned %T, want *redmineConfig", instance)
	}
	if config.Endpoint != nil {
		t.Errorf("expected nil Endpoint, got %v", config.Endpoint)
	}
	if config.APIKey != nil {
		t.Errorf("expected nil APIKey, got %v", config.APIKey)
	}
}
