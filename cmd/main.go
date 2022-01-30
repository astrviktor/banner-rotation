package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/astrviktor/banner-rotation/internal/app"
	"github.com/astrviktor/banner-rotation/internal/config"
)

func main() {
	config := config.DefaultConfig()

	app := app.New(config)

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	app.Start()

	<-exit

	app.Stop()
}
