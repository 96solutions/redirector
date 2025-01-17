package config

import (
	"fmt"
	"time"
)

type DBConf struct {
	Host string `mapstructure:"db_host"`
	Port string `mapstructure:"db_port"`

	User     string `mapstructure:"db_username"`
	Password string `mapstructure:"db_password"`

	Database string `mapstructure:"db_database"`

	ConnectionMaxLife  int `mapstructure:"db_conn_max_life"`
	MaxIdleConnections int `mapstructure:"db_max_idle_conn"`
	MaxOpenConnections int `mapstructure:"db_max_open_conn"`
}

func (m *DBConf) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", m.User, m.Password, m.Host, m.Port, m.Database)
}

func (m *DBConf) ConnectionMaxLifeDuration() time.Duration {
	return time.Duration(m.ConnectionMaxLife) * time.Second
}
