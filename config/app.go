package config

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var applicationConfig *appConfig

type appConfig struct {
	Name string `mapstructure:"APP_NAME"`
	Mode string `mapstructure:"APP_MODE"`

	LogLevel string `mapstructure:"LOG_LEVEL"`
	LogDir   string `mapstructure:"LOG_DIR"`
	LogFile  string `mapstructure:"LOG_FILE"`

	StorageHost     string `mapstructure:"STORAGE_HOST"`
	StoragePort     string `mapstructure:"STORAGE_PORT"`
	StorageUsername string `mapstructure:"STORAGE_USERNAME"`
	StoragePassword string `mapstructure:"STORAGE_PASSWORD"`
	StorageDatabase string `mapstructure:"STORAGE_DATABASE"`

	ServerHost string `mapstructure:"SERVER_HOST"`
	ServerPort string `mapstructure:"SERVER_PORT"`
}

func InitConfig() *appConfig {
	loadEnvVariables()

	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		loadEnvVariables()
	})
	viper.WatchConfig()

	return applicationConfig
}
