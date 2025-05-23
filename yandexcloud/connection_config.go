package yandexcloud

import (
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/schema"
)

// Strict error type for configs
type ConfigError string

func (e ConfigError) Error() string { return string(e) }

// Strict types for identifiers and parameters
type CloudID string

// FolderID is declared only here, do not duplicate
// type FolderID string
type Token string
type UserAgent string
type EndpointOverride string
type LogLevel string

const (
	LogLevelError LogLevel = "error"
	LogLevelInfo  LogLevel = "info"
	LogLevelDebug LogLevel = "debug"
)

func connectionConfig() *plugin.ConnectionConfigSchema {
	return &plugin.ConnectionConfigSchema{
		NewInstance: func() interface{} { return &Config{} },
		Schema: map[string]*schema.Attribute{
			"token":                    {Type: schema.TypeString},
			"service_account_key_file": {Type: schema.TypeString},
			"cloud_id":                 {Type: schema.TypeString},
			"folder_id":                {Type: schema.TypeString},
			"timeout":                  {Type: schema.TypeInt},
			"retry":                    {Type: schema.TypeInt},
			"user_agent":               {Type: schema.TypeString},
			"endpoint_override":        {Type: schema.TypeString},
			"log_level":                {Type: schema.TypeString},
		},
	}
}

type Config struct {
	Token                 *Token            `cty:"token"`
	ServiceAccountKeyFile *string           `cty:"service_account_key_file"`
	CloudID               *CloudID          `cty:"cloud_id"`
	FolderID              *FolderID         `cty:"folder_id"`
	Timeout               *int              `cty:"timeout"`
	Retry                 *int              `cty:"retry"`
	UserAgent             *UserAgent        `cty:"user_agent"`
	EndpointOverride      *EndpointOverride `cty:"endpoint_override"`
	LogLevel              *LogLevel         `cty:"log_level"`
}

// ValidateConfig checks required and conflicting config parameters.
func ValidateConfig(cfg *Config) error {
	if cfg == nil {
		return ConfigError("config is nil")
	}
	if (cfg.Token == nil || *cfg.Token == "") && (cfg.ServiceAccountKeyFile == nil || *cfg.ServiceAccountKeyFile == "") {
		return ConfigError("either token or service_account_key_file must be set in connection config")
	}
	if cfg.Token != nil && *cfg.Token != "" && cfg.ServiceAccountKeyFile != nil && *cfg.ServiceAccountKeyFile != "" {
		return ConfigError("only one of token or service_account_key_file should be set, not both")
	}
	if cfg.Timeout != nil && *cfg.Timeout < 0 {
		return ConfigError("timeout must be >= 0")
	}
	if cfg.Retry != nil && *cfg.Retry < 1 {
		return ConfigError("retry must be >= 1")
	}
	return nil
}
