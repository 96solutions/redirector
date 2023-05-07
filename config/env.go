package config

import (
	"log"

	"github.com/spf13/viper"
)

// Call to load the variables from env
func loadEnvVariables() {
	// Tell viper the path/location of your env file
	//viper.AddConfigPath(".")

	// Tell viper the name of your file
	viper.SetConfigFile("config.env")

	// Tell viper the type of your file
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	// Viper reads all the variables from env file and log error if any found
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading env file ", err)
	}

	// Viper unmarshals the loaded env varialbes into the struct
	if err := viper.Unmarshal(&applicationConfig); err != nil {
		log.Fatal(err)
	}

	return
}
