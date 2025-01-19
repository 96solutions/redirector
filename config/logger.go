package config

type LoggerConf struct {
	Level          string `mapstructure:"log_level"`
	IsJSON         bool   `mapstructure:"log_is_json"`
	AddSource      bool   `mapstructure:"log_add_source"`
	ReplaceDefault bool   `mapstructure:"log_replace_default"`
}
