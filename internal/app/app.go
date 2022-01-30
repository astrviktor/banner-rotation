package app

import (
	"github.com/astrviktor/banner-rotation/internal/config"
	internalhttp "github.com/astrviktor/banner-rotation/internal/server/http"
	memorystorage "github.com/astrviktor/banner-rotation/internal/storage/memory"
)

type App struct {
	config config.Config
	server *internalhttp.Server
}

func New(config config.Config) *App {
	storage := memorystorage.New()
	server := internalhttp.NewServer(config.HTTPServer.Host, config.HTTPServer.Port, storage)
	return &App{config, server}
}

func (a *App) Start() {
	a.server.Start()
}

func (a *App) Stop() {
	a.server.Stop()
}
