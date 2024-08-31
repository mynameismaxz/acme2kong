package app

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/mynameismaxz/acme2kong/config"
	"github.com/mynameismaxz/acme2kong/pkg/acme"
	"github.com/mynameismaxz/acme2kong/pkg/kong"
	"github.com/mynameismaxz/acme2kong/pkg/logger"
)

const (
	DEFAULT_RSA_BIT_SIZE = 4096
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
	kongClient, err := kong.New(a.cfg.KongEndpoint, a.log)
	if err != nil {
		return err
	}

	// TODO: implement the logic to check that have certificate in the path before generate a new certificate
	// privKey, err := rsa.GenerateKey(rand.Reader, 4096)
	// if err != nil {
	// 	return err
	// }
	privKey, err := a.getPrivateKey(a.cfg.CertPath)
	if err != nil {
		return err
	}

	domains := []string{a.cfg.DomainName}
	userProfile := &acme.User{
		Email:        a.cfg.RegistrationEmail,
		Registration: nil,
		Key:          privKey,
	}

	acmeClient := acme.NewClient(userProfile, a.cfg.ChallengeProvider, domains, a.cfg.CertPath, a.log)
	if acmeClient == nil {
		return fmt.Errorf("failed to create acme client")
	}

	if err := acmeClient.GenerateNewCertificate(); err != nil {
		return err
	}

	// read the certificate
	cert, err := os.ReadFile(fmt.Sprintf("%s/certificate.crt", a.cfg.CertPath))
	if err != nil {
		return err
	}
	privCertKey, err := os.ReadFile(fmt.Sprintf("%s/private.key", a.cfg.CertPath))
	if err != nil {
		return err
	}

	// TODO: implement the logic to update the Kong API Gateway with the new certificate
	if err := kongClient.UpdateCertificate(cert, privCertKey); err != nil {
		return err
	}

	return nil
}

func (a *App) getPrivateKey(certPath string) (*rsa.PrivateKey, error) {
	// check the cert path is have the private key or not
	if _, err := os.Stat(fmt.Sprintf("%s/private.key", certPath)); os.IsNotExist(err) {
		// if not exist, create the new one
		a.log.Info("Generating new private key")
		privKey, err := rsa.GenerateKey(rand.Reader, DEFAULT_RSA_BIT_SIZE)
		if err != nil {
			return nil, err
		}

		// save the private key to the cert path
		privKeyBytes := x509.MarshalPKCS1PrivateKey(privKey)
		if err := os.WriteFile(fmt.Sprintf("%s/private.key", certPath), privKeyBytes, 0644); err != nil {
			return nil, err
		}

		return privKey, nil
	}

	// if exist, read the private key
	privKeyBytes, err := os.ReadFile(fmt.Sprintf("%s/private.key", certPath))
	if err != nil {
		return nil, err
	}

	// decode the PEM block
	block, _ := pem.Decode(privKeyBytes)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing private key")
	}

	// convert the private key bytes to rsa.PrivateKey
	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privKey, nil
}
