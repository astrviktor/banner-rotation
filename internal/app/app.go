package app

import (
	"github.com/astrviktor/banner-rotation/internal/config"
	internalhttp "github.com/astrviktor/banner-rotation/internal/server/http"
	sqlstorage "github.com/astrviktor/banner-rotation/internal/storage/sql"
)

type App struct {
	config config.Config
	server *internalhttp.Server
}

func New(config config.Config) *App {
	storage := sqlstorage.New("postgres://user:password@postgres:5432/banner_rotation")
	server := internalhttp.NewServer(config.HTTPServer.Host, config.HTTPServer.Port, storage)
	return &App{config, server}
}

func (a *App) Start() {
	a.server.Start()
}

func (a *App) Stop() {
	a.server.Stop()
}
