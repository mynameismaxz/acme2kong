package app

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"time"

	"github.com/mynameismaxz/acme2kong/config"
	"github.com/mynameismaxz/acme2kong/pkg/acme"
	"github.com/mynameismaxz/acme2kong/pkg/kong"
	"github.com/mynameismaxz/acme2kong/pkg/logger"
)

const (
	DEFAULT_RSA_BIT_SIZE    = 4096
	PEM_BLOCK_TYPE          = "RSA PRIVATE KEY"
	CERTIFICATE_FILE_NAME   = "certificate.crt"
	PRIVATE_KEY_FILE_NAME   = "private.key"
	ISSUER_FILE_NAME        = "issuer.crt"
	CERT_RESOURCE_FILE_NAME = "cert_resource.json"
	INTERVAL_CHECK_CERT     = 24 // in hours
	CERTIFICATE_RENEW_AFTER = 30 // in days
)

type App struct {
	log  *logger.Logger
	cfg  *config.Config
	kong *kong.Kong
}

func New(l *logger.Logger, conf *config.Config) *App {
	// cast domain name
	domain := []string{conf.DomainName}

	kongClient, err := kong.New(conf.KongEndpoint, domain, l)
	if err != nil {
		l.Error("failed to create kong client")
		return nil
	}

	return &App{
		log:  l,
		cfg:  conf,
		kong: kongClient,
	}
}

func (a *App) Run() error {
	var newCertRegister bool
	a.log.Info("acme2kong is running")

	// create user profile
	userProfile := &acme.User{
		Email:        a.cfg.RegistrationEmail,
		Registration: nil,
		Key:          nil,
	}

	certificateListedFile := []string{
		CERTIFICATE_FILE_NAME,
		PRIVATE_KEY_FILE_NAME,
		ISSUER_FILE_NAME,
		CERT_RESOURCE_FILE_NAME,
	}

	// check from certificate directory have file from certificatedListedFile or not
	for _, file := range certificateListedFile {
		if _, err := os.Stat(fmt.Sprintf("%s/%s", a.cfg.CertPath, file)); os.IsNotExist(err) {
			newCertRegister = true
			break
		}
	}

	if newCertRegister {
		// New register certificate logic.
		a.log.Info("New certificate registration")
		if err := a.generateNewCertificate(userProfile); err != nil {
			return err
		}

		// Update Kong API Gateway
		if err := a.updateKongAPIGateway(); err != nil {
			return err
		}
	} else {
		// Check the cert is expired or not.
		// If expired, renew the certificate immediately.
		cert, err := a.readCertificate()
		if err != nil {
			return err
		}

		if a.checkCertificateExpired(cert) {
			if err := a.renewCertificate(userProfile); err != nil {
				return err
			}

			// Update Kong API Gateway
			if err := a.updateKongAPIGateway(); err != nil {
				return err
			}
		} else {
			a.log.Info("Certificate is not expired")
		}
	}

	interval := time.NewTicker(time.Hour * INTERVAL_CHECK_CERT)
	a.log.Info(fmt.Sprintf("Start checking certificate interval in every %d hours", INTERVAL_CHECK_CERT))

	for {
		select {
		case <-interval.C:
			a.log.Info("Checking certificate interval...")
			cert, err := a.readCertificate()
			if err != nil {
				return err
			}

			if a.checkCertificateExpired(cert) {
				// if true, expired
				userProfile.Key = nil
				userProfile.Registration = nil
				if err := a.renewCertificate(userProfile); err != nil {
					return err
				}

				// Update Kong API Gateway
				if err := a.updateKongAPIGateway(); err != nil {
					return err
				}
			} else {
				a.log.Info("Certificate is not expired")
			}
		default:
		}
	}
}

// GenerateNewCertificate will generate a new certificate.
func (a *App) generateNewCertificate(user *acme.User) error {
	// check that have certificate in the path before generate a new certificate
	privKey, err := a.getPrivateKey()
	if err != nil {
		return err
	}
	user.Key = privKey
	domains := []string{a.cfg.DomainName}

	acmeClient := acme.NewClient(user, a.cfg.ChallengeProvider, domains, a.cfg.CertPath, a.log)
	if acmeClient == nil {
		return fmt.Errorf("failed to create acme client")
	}

	if err := acmeClient.GenerateNewCertificate(); err != nil {
		return err
	}

	return nil
}

// TODO: Renew the certificate.
func (a *App) renewCertificate(user *acme.User) error {
	if err := a.generateNewCertificate(user); err != nil {
		return err
	}

	return nil
}

// GetPrivateKey will return the private key (if not have) will generate new private key.
func (a *App) getPrivateKey() (*rsa.PrivateKey, error) {
	// check the cert path is have the private key or not
	if _, err := os.Stat(fmt.Sprintf("%s/%s", a.cfg.CertPath, PRIVATE_KEY_FILE_NAME)); os.IsNotExist(err) {
		// if not exist, create the new one
		a.log.Info("Generating new private key")
		privKey, err := rsa.GenerateKey(rand.Reader, DEFAULT_RSA_BIT_SIZE)
		if err != nil {
			return nil, err
		}
		return privKey, nil
	}

	// if exist, read the private key
	a.log.Info("Reading existing private key")
	privKeyBytes, err := os.ReadFile(fmt.Sprintf("%s/%s", a.cfg.CertPath, PRIVATE_KEY_FILE_NAME))
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

// Read the certificate
func (a *App) readCertificate() (*x509.Certificate, error) {
	pemData, err := os.ReadFile(fmt.Sprintf("%s/%s", a.cfg.CertPath, CERTIFICATE_FILE_NAME))
	if err != nil {
		return nil, err
	}

	// Decode the PEM block
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block containing the certificate")
	}

	// Parse the certificate to x509.Certificate
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}

	return cert, nil
}

// Update to Kong API Gateway
func (a *App) updateKongAPIGateway() error {
	// read the certificate
	cert, err := os.ReadFile(fmt.Sprintf("%s/%s", a.cfg.CertPath, CERTIFICATE_FILE_NAME))
	if err != nil {
		return err
	}
	privCertKey, err := os.ReadFile(fmt.Sprintf("%s/%s", a.cfg.CertPath, PRIVATE_KEY_FILE_NAME))
	if err != nil {
		return err
	}

	// update the certificate to Kong API Gateway
	if err := a.kong.UpdateCertificate(cert, privCertKey); err != nil {
		return err
	}

	return nil
}

func (a *App) checkCertificateExpired(cert *x509.Certificate) bool {
	// check the certificate is expired or not
	timeLeft := cert.NotAfter.Sub(time.Now().UTC())
	// convert timeleft to days
	daysRemaining := timeLeft.Hours() / INTERVAL_CHECK_CERT
	a.log.Info(fmt.Sprintf("Certificate will be expired in %d days", int(daysRemaining)))

	if daysRemaining <= CERTIFICATE_RENEW_AFTER {
		return true
	}

	return false
}
