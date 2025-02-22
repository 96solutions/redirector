// Package config contains structures that represent configs for different application modules.
package config

// LoggerConf contains logging configuration settings.
type LoggerConf struct {
	// Level sets the minimum logging level (debug, info, warn, error)
	Level string `mapstructure:"log_level"`
	// IsJSON determines if logs should be formatted as JSON
	IsJSON bool `mapstructure:"log_is_json"`
	// AddSource determines if source code location should be added to log entries
	AddSource bool `mapstructure:"log_add_source"`
	// ReplaceDefault determines if this logger should replace the default logger
	ReplaceDefault bool `mapstructure:"log_replace_default"`

	// OpenSearchHost is the OpenSearch server hostname
	OpenSearchHost string `mapstructure:"log_open_search_host"`
	// OpenSearchPort is the OpenSearch server port
	OpenSearchPort string `mapstructure:"log_open_search_port"`
	// OpenSearchIndex is the name of the OpenSearch index to use
	OpenSearchIndex string `mapstructure:"log_open_search_index"`
	// OpenSearchUser is the OpenSearch username
	OpenSearchUser string `mapstructure:"log_open_search_user"`
	// OpenSearchPass is the OpenSearch password
	OpenSearchPass string `mapstructure:"log_open_search_pass"`
}
