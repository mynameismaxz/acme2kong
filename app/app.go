package app

import (
	"crypto/rand"
	"crypto/rsa"
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

	// TODO: implement the logic to check that have certificate in the path before generate a new certificate
	privKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}

	userProfile := &acme.User{
		Email:        a.cfg.RegistrationEmail,
		Registration: nil,
		Key:          privKey,
	}

	acmeClient := acme.NewClient(userProfile, a.cfg.ChallengeProvider, a.log)
	if acmeClient == nil {
		return fmt.Errorf("failed to create acme client")
	}

	if err := acmeClient.GenerateNewCertificate(); err != nil {
		return err
	}

	return nil
}
