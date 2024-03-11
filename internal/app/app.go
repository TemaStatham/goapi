package app

import (
	"goapi/internal/app/logger"
	"log/slog"
)

type App struct {
	log *slog.Logger
}

func New(env string) *App {
	return &App{
		log: logger.SetupLogger(env),
	}
}

func (a *App) Run() {

}
