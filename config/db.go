// Package config contains structures that represent configs for different application modules.
package config

import (
	"fmt"
	"time"
)

// DBConf contains database connection configuration settings.
type DBConf struct {
	// Host is the database server hostname
	Host string `mapstructure:"db_host"`
	// Port is the database server port
	Port string `mapstructure:"db_port"`

	// User is the database username
	User string `mapstructure:"db_username"`
	// Password is the database password
	Password string `mapstructure:"db_password"`

	// Database is the name of the database to connect to
	Database string `mapstructure:"db_database"`

	// ConnectionMaxLife is the maximum amount of time a connection may be reused (in seconds)
	ConnectionMaxLife int `mapstructure:"db_conn_max_life"`
	// MaxIdleConnections is the maximum number of idle connections in the pool
	MaxIdleConnections int `mapstructure:"db_max_idle_conn"`
	// MaxOpenConnections is the maximum number of open connections to the database
	MaxOpenConnections int `mapstructure:"db_max_open_conn"`
}

// DSN returns the database connection string.
func (m *DBConf) DSN() string {
	return fmt.Sprintf("%s:%s@%s:%s/%s", m.User, m.Password, m.Host, m.Port, m.Database)
}

// ConnectionMaxLifeDuration returns the maximum amount of time a connection may be reused.
func (m *DBConf) ConnectionMaxLifeDuration() time.Duration {
	return time.Duration(m.ConnectionMaxLife) * time.Second
}
