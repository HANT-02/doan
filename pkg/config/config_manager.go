package config

import (
	"context"
	"time"
)

type Manager interface {
	// Start begins watching Consul for config changes
	Start(ctx context.Context)

	// WatchKey watches a specific configuration key for changes
	WatchKey(key string) <-chan struct{}

	// UnwatchKey stops watching a specific configuration key
	UnwatchKey(key string)

	// Get retrieves a configuration value
	Get(key string) interface{}

	// GetString retrieves a string configuration value
	GetString(key string) string

	// GetInt retrieves an integer configuration value
	GetInt(key string) int

	// GetBool retrieves a boolean configuration value
	GetBool(key string) bool

	// GetDuration retrieves a duration configuration value
	GetDuration(key string) time.Duration

	// GetStringSlice retrieves a string slice configuration value
	GetStringSlice(key string) []string

	// GetStringMap retrieves a string map configuration value
	GetStringMap(key string) map[string]interface{}

	// UnmarshalKey unmarshal a configuration key into a struct
	UnmarshalKey(key string, rawVal interface{}) error

	// Unmarshal unmarshall the entire configuration into a struct
	Unmarshal(rawVal interface{}) error

	// Set sets a configuration value
	Set(key string, value interface{})

	// SetDefault sets a default configuration value
	SetDefault(key string, value interface{})

	// IsSet checks if a configuration key is set
	IsSet(key string) bool

	// AllKeys returns all configuration keys
	AllKeys() []string

	// AllSettings returns all configuration settings
	AllSettings() map[string]interface{}
}
