package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/astrviktor/banner-rotation/internal/app"
	"github.com/astrviktor/banner-rotation/internal/config"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "config_local.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	config := config.NewConfig(configFile)
	//config := config.NewConfig("/home/astrviktor/golang/src/banner-rotation/configs/config_local.yaml")

	app := app.New(config)

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	app.Start()

	<-exit

	app.Stop()
}
