package config_test

import (
	"os"
	"testing"

	cfg "longtrail-api/config"
)

func TestEnvConfig(t *testing.T) {
	os.Setenv("COGNITO_USER_POOL_ID", "COGNITO_USER_POOL_ID")
	os.Setenv("EVENTS_TABLE_NAME", "EVENTS_TABLE_NAME")

	config := cfg.NewFromEnv()

	if config == nil {
		t.Error("Config object nil")
	}
	if config.UserPoolID != "COGNITO_USER_POOL_ID" {
		t.Log(config.UserPoolID)
		t.Error("UserPoolID incorrect")
	}
	if config.EventsTableName != "EVENTS_TABLE_NAME" {
		t.Log(config.EventsTableName)
		t.Error("EventsTableName incorrect")
	}
}
