package config

type httpServerConf struct {
	Host        string `mapstructure:"http_server_host"`
	Port        string `mapstructure:"http_server_port"`
	SSL         bool   `mapstructure:"http_server_ssl"`
	SSLCertPath string `mapstructure:"http_server_cert"`
}
