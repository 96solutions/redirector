// Package config contains structures that represent configs for different application modules.
package config

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var lock = &sync.Mutex{}
var applicationConfig *AppConfig

type AppConfig struct {
	DBConf         *DBConf
	HttpServerConf *HttpServerConf
	LogConf        *LoggerConf

	GeoIP2DBPath string `mapstructure:"geoip2_db_path"`
}

func (cfg *AppConfig) String() string {
	cfgJSON, _ := json.MarshalIndent(cfg, "", "    ")

	return string(cfgJSON)
}

func GetConfig() *AppConfig {
	if applicationConfig == nil {
		lock.Lock()
		defer lock.Unlock()
		if applicationConfig == nil {
			applicationConfig = initConfig()
		}
	}

	return applicationConfig
}

func initConfig() *AppConfig {
	// Tell viper the path/location of your env file
	// viper.AddConfigPath(".")

	// Tell viper the name of your file
	viper.SetConfigFile(viper.Get("config").(string))

	// Tell viper the type of your file
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	// Viper reads all the variables from env file and log error if any found
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("error reading env file. error: %w", err))
	}

	cfg := new(AppConfig)
	cfg.HttpServerConf = new(HttpServerConf)
	cfg.DBConf = new(DBConf)
	cfg.LogConf = new(LoggerConf)

	// Viper unmarshal the loaded env variables into the config structs
	if err := viper.Unmarshal(&cfg.HttpServerConf); err != nil {
		panic(fmt.Errorf("cannot unmarshal HttpServerConf. error: %w", err))
	}
	if err := viper.Unmarshal(&cfg.DBConf); err != nil {
		panic(fmt.Errorf("cannot unmarshal DBConf. error: %s", err))
	}
	if err := viper.Unmarshal(&cfg.LogConf); err != nil {
		panic(fmt.Errorf("cannot unmarshal LogConf. error: %s", err))
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		panic(fmt.Errorf("cannot unmarshal GeoIP2DBPath. error: %s", err))
	}

	// add watcher on init
	if applicationConfig == nil {
		viper.OnConfigChange(func(e fsnotify.Event) {
			applicationConfig = initConfig()
		})
		viper.WatchConfig()
	}

	return cfg
}
