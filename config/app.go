package config

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var lock = &sync.Mutex{}
var applicationConfig *appConfig

type appConfig struct {
	MySQLConf      *mysqlConf
	LoggerConf     *loggerConf
	HttpServerConf *httpServerConf
}

func (cfg *appConfig) String() string {
	cfgJSON, _ := json.MarshalIndent(cfg, "", "    ")

	return string(cfgJSON)
}

func GetConfig() *appConfig {
	if applicationConfig == nil {
		lock.Lock()
		defer lock.Unlock()
		if applicationConfig == nil {
			applicationConfig = initConfig()
		}
	}

	return applicationConfig
}

func initConfig() *appConfig {
	// Tell viper the path/location of your env file
	//viper.AddConfigPath(".")

	// Tell viper the name of your file
	viper.SetConfigFile(viper.Get("config").(string))

	// Tell viper the type of your file
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	// Viper reads all the variables from env file and log error if any found
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading env file ", err)
	}

	cfg := new(appConfig)
	cfg.HttpServerConf = new(httpServerConf)
	cfg.MySQLConf = new(mysqlConf)
	cfg.LoggerConf = new(loggerConf)

	// Viper unmarshals the loaded env varialbes into the config structs
	if err := viper.Unmarshal(&cfg.HttpServerConf); err != nil {
		log.Fatalf("cannot unmarshal HttpServerConf. error: %s", err)
	}
	if err := viper.Unmarshal(&cfg.MySQLConf); err != nil {
		log.Fatalf("cannot unmarshal MySQLConf. error: %s", err)
	}
	if err := viper.Unmarshal(&cfg.LoggerConf); err != nil {
		log.Fatalf("cannot unmarshal LoggerConf. error: %s", err)
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
