package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/astrviktor/banner-rotation/internal/config"
	internalhttp "github.com/astrviktor/banner-rotation/internal/server/http"
)

func main() {
	config.GlobalConfig = config.DefaultConfig()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	httpServer := internalhttp.NewServer(config.GlobalConfig.HTTPServer.Host, config.GlobalConfig.HTTPServer.Port)

	httpServer.Start()

	<-exit

	httpServer.Stop()
}
