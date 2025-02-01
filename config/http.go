package config

import (
	"fmt"
	"time"
)

// HTTPServerConf describe HTTP Server configuration.
type HTTPServerConf struct {
	Host              string `mapstructure:"http_server_host"`
	Port              int    `mapstructure:"http_server_port"`
	ReadTimeout       int    `mapstructure:"http_server_read_timeout"`
	WriteTimeout      int    `mapstructure:"http_server_write_timeout"`
	ConnectionTimeout int    `mapstructure:"http_server_connection_timeout"`
	ShutdownTimeout   int    `mapstructure:"http_server_shutdown_timeout"`
	SSLCertPath       string `mapstructure:"http_server_cert"`
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
