package app

import (
	"fmt"

	"github.com/mynameismaxz/acme2kong/config"
	"github.com/mynameismaxz/acme2kong/pkg/acme"
	"github.com/mynameismaxz/acme2kong/pkg/logger"
)

type App struct {
	log *logger.Logger
	cfg *config.Config
}

func New(l *logger.Logger, conf *config.Config) *App {
	return &App{
		log: l,
		cfg: conf,
	}
}

func (a *App) Run() error {
	a.log.Info("acme2kong is running")

	// test section
	acmeClient := acme.NewClient("https://acme-staging-v02.api.letsencrypt.org/directory", a.cfg.DomainName, a.cfg.RegistrationEmail, a.log)
	// check nil
	if acmeClient == nil {
		return fmt.Errorf("failed to create acme client")
	}

	if err := acmeClient.Register(); err != nil {
		return err
	}

	return nil
}
