package app

import (
	"github.com/astrviktor/banner-rotation/internal/config"
	internalhttp "github.com/astrviktor/banner-rotation/internal/server/http"
	"github.com/astrviktor/banner-rotation/internal/storage"
	memorystorage "github.com/astrviktor/banner-rotation/internal/storage/memory"
	sqlstorage "github.com/astrviktor/banner-rotation/internal/storage/sql"
)

type App struct {
	config config.Config
	server *internalhttp.Server
}

func New(conf config.Config) *App {
	var stor storage.Storage
	if conf.DB.Mode == config.DBMemoryMode {
		stor = memorystorage.New()
	} else {
		stor = sqlstorage.New(conf)
	}

	server := internalhttp.NewServer(conf.HTTPServer.Host, conf.HTTPServer.Port, stor)
	return &App{conf, server}
}

func (a *App) Start() {
	a.server.Start()
}

func (a *App) Stop() {
	a.server.Stop()
}
