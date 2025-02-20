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

// AppConfig type represents application config.
type AppConfig struct {
	DBConf         *DBConf
	HTTPServerConf *HTTPServerConf
	LogConf        *LoggerConf

	GeoIP2DBPath string `mapstructure:"geoip2_db_path"`
}

// String function implements Stringer interfaces and used to represent
// application configuration as string, mostly used for logging.
func (cfg *AppConfig) String() string {
	cfgJSON, _ := json.MarshalIndent(cfg, "", "    ")

	return string(cfgJSON)
}

// GetConfig function provides access to the application config.
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

	viper.AutomaticEnv()

	if viper.Get("config").(string) != "" {
		// Tell viper the name of your file
		viper.SetConfigFile(viper.Get("config").(string))
		// Tell viper the type of your file
		viper.SetConfigType("env")
		// Viper reads all the variables from env file and log error if any found
		if err := viper.ReadInConfig(); err != nil {
			panic(fmt.Errorf("error reading env file. error: %w", err))
		}
	}

	cfg := new(AppConfig)
	cfg.HTTPServerConf = new(HTTPServerConf)
	cfg.DBConf = new(DBConf)
	cfg.LogConf = new(LoggerConf)

	// Viper unmarshal the loaded env variables into the config structs
	if err := viper.Unmarshal(&cfg.HTTPServerConf); err != nil {
		panic(fmt.Errorf("cannot unmarshal HttpServerConf. error: %w", err))
	}
	if err := viper.Unmarshal(&cfg.DBConf); err != nil {
		panic(fmt.Errorf("cannot unmarshal DBConf. error: %w", err))
	}
	if err := viper.Unmarshal(&cfg.LogConf); err != nil {
		panic(fmt.Errorf("cannot unmarshal LogConf. error: %w", err))
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		panic(fmt.Errorf("cannot unmarshal GeoIP2DBPath. error: %w", err))
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
