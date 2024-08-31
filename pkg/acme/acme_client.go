package acme

import (
	"fmt"
	"os"

	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge/dns01"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/providers/dns"
	"github.com/go-acme/lego/v4/registration"
	"github.com/mynameismaxz/acme2kong/pkg/logger"
)

var (
	DNS_RESOLVER = []string{"1.1.1.1", "1.0.0.1"}
)

type ACME struct {
	User        *User
	DNSProvider string
	DomainName  []string
	CertPath    string

	client *lego.Client
	logger *logger.Logger
}

func NewClient(user *User, provider string, domainName []string, certPath string, logger *logger.Logger) *ACME {
	config := lego.NewConfig(user)
	// for development, use the staging environment.
	config.CADirURL = lego.LEDirectoryStaging
	legoClient, err := lego.NewClient(config)
	if err != nil {
		return nil
	}

	return &ACME{
		User:        user,
		DNSProvider: provider,
		DomainName:  domainName,
		CertPath:    certPath,
		client:      legoClient,
		logger:      logger,
	}
}

// GenerateNewCertificate generates a new certificate using the ACME client.
// It returns an error if the certificate generation fails.
func (ac *ACME) GenerateNewCertificate() error {
	// create provider
	provider, err := dns.NewDNSChallengeProviderByName("cloudflare")
	if err != nil {
		return err
	}

	if err = ac.client.Challenge.SetDNS01Provider(
		provider,
		dns01.CondOption((len(ac.DomainName) > 0),
			dns01.AddRecursiveNameservers(dns01.ParseNameservers(DNS_RESOLVER)))); err != nil {
		return err
	}

	// check registration
	if ac.User.Registration != nil {
		ac.logger.Info("User already registered, skip registration.")
	} else {
		ac.logger.Info("Registering...")
		reg, err := ac.client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
		if err != nil {
			return err
		}
		ac.User.Registration = reg
	}

	// obtain a new certificate
	request := certificate.ObtainRequest{
		Domains: ac.DomainName,
		Bundle:  true,
	}

	cert, err := ac.client.Certificate.Obtain(request)
	if err != nil {
		return err
	}

	// Save the certificate to the path
	if err := ac.saveCertificate(cert, ac.CertPath); err != nil {
		return err
	}
	ac.logger.Info(fmt.Sprintf("Certificate generated and saved to %s", ac.CertPath))

	return nil
}

// TODO: implement RenewCertificate
func (ac *ACME) RenewCertificate() error {
	return nil
}

// Save certificate to disk
func (ac *ACME) saveCertificate(cert *certificate.Resource, path string) error {
	// save the certificate
	if err := os.WriteFile(fmt.Sprintf("%s/certificate.crt", path), cert.Certificate, 0644); err != nil {
		return err
	}

	// save the private key
	if err := os.WriteFile(fmt.Sprintf("%s/private.key", path), cert.PrivateKey, 0644); err != nil {
		return err
	}

	// save the issuer certificate
	if err := os.WriteFile(fmt.Sprintf("%s/issuer.crt", path), cert.IssuerCertificate, 0644); err != nil {
		return err
	}

	return nil
}
