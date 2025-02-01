package config

type LoggerConf struct {
	Level          string `mapstructure:"log_level"`
	IsJSON         bool   `mapstructure:"log_is_json"`
	AddSource      bool   `mapstructure:"log_add_source"`
	ReplaceDefault bool   `mapstructure:"log_replace_default"`

	OpenSearchHost  string `mapstructure:"log_open_search_host"`
	OpenSearchPort  string `mapstructure:"log_open_search_port"`
	OpenSearchIndex string `mapstructure:"log_open_search_index"`
	OpenSearchUser  string `mapstructure:"log_open_search_user"`
	OpenSearchPass  string `mapstructure:"log_open_search_pass"`
}
