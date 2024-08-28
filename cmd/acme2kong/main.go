package main

import (
	"github.com/mynameismaxz/acme2kong/app"
	"github.com/mynameismaxz/acme2kong/config"
	"github.com/mynameismaxz/acme2kong/pkg/logger"
)

func main() {
	// Initialize logger
	lg := logger.New()

	// Initialize config
	cfg, err := config.Initialize()
	if err != nil {
		lg.Error(err.Error())
	}

	// Initialize app
	app := app.New(lg, cfg)

	if err := app.Run(); err != nil {
		lg.Error(err.Error())
	}
}
