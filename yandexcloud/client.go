package yandexcloud

import (
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

// Strict error type for client.go
type ClientError string

func (e ClientError) Error() string { return string(e) }

// GetConfig is syntactic sugar to cast connection config.
func GetConfig(conn *plugin.Connection) *Config {
	if conn == nil || conn.Config == nil {
		return nil
	}
	// If already a pointer
	if cfg, ok := conn.Config.(*Config); ok {
		return cfg
	}
	// If value (struct), convert to pointer
	if cfg, ok := conn.Config.(Config); ok {
		return &cfg
	}
	return nil
}

// --- KMS Key types and REST client ---
// (All code related to KMS has been removed. It is now in kms_client.go)
