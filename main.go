package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/lroman242/redirector/config"
	"github.com/spf13/viper"
)

func main() {
	conf := config.InitConfig()

	fmt.Printf("%+v\n", conf)
	fmt.Printf("%+v\n", viper.AllSettings())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c
	os.Exit(1)
}
