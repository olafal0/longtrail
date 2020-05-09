package config

import "os"

const userPoolEnvKey = "COGNITO_USER_POOL_ID"
const eventsTableNameEnvKey = "EVENTS_TABLE_NAME"

// LongtrailConfig is an object containing global configuration for the app.
type LongtrailConfig struct {
	UserPoolID      string
	EventsTableName string
}

// NewFromEnv returns a LongtrailConfig object initialized from environment variables.
func NewFromEnv() *LongtrailConfig {
	config := &LongtrailConfig{}
	config.UserPoolID = os.Getenv(userPoolEnvKey)
	config.EventsTableName = os.Getenv(eventsTableNameEnvKey)

	return config
}
