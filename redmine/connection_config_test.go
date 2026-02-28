package redmine

import (
	"testing"
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
