package config

import "fmt"

type mysqlConf struct {
	Host string `mapstructure:"mysql_host"`
	Port string `mapstructure:"mysql_port"`

	User     string `mapstructure:"mysql_username"`
	Password string `mapstructure:"mysql_password"`

	Database string `mapstructure:"mysql_database"`
}

func (m *mysqlConf) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", m.User, m.Password, m.Host, m.Port, m.Database)
}
