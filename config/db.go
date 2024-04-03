package config

import "fmt"

type dbConf struct {
	Host string `mapstructure:"db_host"`
	Port string `mapstructure:"db_port"`

	User     string `mapstructure:"db_username"`
	Password string `mapstructure:"db_password"`

	Database string `mapstructure:"db_database"`
}

func (m *dbConf) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", m.User, m.Password, m.Host, m.Port, m.Database)
}
