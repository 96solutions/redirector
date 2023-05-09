package config

type loggerConf struct {
	LogLevel string `mapstructure:"log_level"`
	LogFile  string `mapstructure:"log_file"`
	LogDir   string `mapstructure:"log_dir"`
}
