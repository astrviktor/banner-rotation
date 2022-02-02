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
	flag.StringVar(&configFile, "config", "config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	config := config.NewConfig(configFile)

	app := app.New(config)

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	app.Start()

	<-exit

	app.Stop()
}
