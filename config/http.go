// Package config contains structures that represent configs for different application modules.
package config

import (
	"fmt"
	"time"
)

// HTTPServerConf describes HTTP Server configuration settings.
type HTTPServerConf struct {
	// Host is the server hostname to listen on
	Host string `mapstructure:"http_server_host"`
	// Port is the server port to listen on
	Port int `mapstructure:"http_server_port"`
	// ReadTimeout is the maximum duration for reading the entire request (in seconds)
	ReadTimeout int `mapstructure:"http_server_read_timeout"`
	// WriteTimeout is the maximum duration for writing the response (in seconds)
	WriteTimeout int `mapstructure:"http_server_write_timeout"`
	// ConnectionTimeout is the maximum duration for waiting for new connections (in seconds)
	ConnectionTimeout int `mapstructure:"http_server_connection_timeout"`
	// ShutdownTimeout is the maximum duration to wait for server shutdown (in seconds)
	ShutdownTimeout int `mapstructure:"http_server_shutdown_timeout"`
	// SSLCertPath is the path to SSL certificate file
	SSLCertPath string `mapstructure:"http_server_cert"`
}

// GetHTTPReadTimeout return server read timeout configuration value.
func (c *HTTPServerConf) GetHTTPReadTimeout() time.Duration {
	return time.Duration(c.ReadTimeout) * time.Second
}

// GetHTTPWriteTimeout return server write timeout configuration value.
func (c *HTTPServerConf) GetHTTPWriteTimeout() time.Duration {
	return time.Duration(c.WriteTimeout) * time.Second
}

// GetHTTPAddress return server host string.
func (c *HTTPServerConf) GetHTTPAddress() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// GetShutdownTimeout return server shutdown timeout configuration value.
func (c *HTTPServerConf) GetShutdownTimeout() time.Duration {
	return time.Duration(c.ShutdownTimeout) * time.Second
}
