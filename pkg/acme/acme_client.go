package acme

import (
	"fmt"

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

	client *lego.Client
	logger *logger.Logger
}

func NewClient(user *User, provider string, logger *logger.Logger) *ACME {
	config := lego.NewConfig(user)
	legoClient, err := lego.NewClient(config)
	if err != nil {
		return nil
	}

	return &ACME{
		User:        user,
		DNSProvider: provider,
		client:      legoClient,
		logger:      logger,
	}
}

// TODO: implement GenerateNewCertificate
func (ac *ACME) GenerateNewCertificate() error {
	tmp := []string{"*.tha.mymacz.com"}

	// create provider
	provider, err := dns.NewDNSChallengeProviderByName("cloudflare")
	if err != nil {
		return err
	}

	if err = ac.client.Challenge.SetDNS01Provider(
		provider,
		dns01.CondOption((len(tmp) > 0),
			dns01.AddRecursiveNameservers(dns01.ParseNameservers(DNS_RESOLVER)))); err != nil {
		return err
	}

	// check registration
	if ac.User.Registration == nil {
		ac.logger.Info("Registering...")
		reg, err := ac.client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
		if err != nil {
			return err
		}
		ac.User.Registration = reg
	} else {
		ac.logger.Info("User already registered, skip registration.")
	}

	// obtain a new certificate
	request := certificate.ObtainRequest{
		Domains: tmp,
		Bundle:  true,
	}

	cert, err := ac.client.Certificate.Obtain(request)
	if err != nil {
		return err
	}

	ac.logger.Info(fmt.Sprintf("Certificate obtained: %s", cert.Domain))

	return nil
}

// TODO: implement RenewCertificate
func (ac *ACME) RenewCertificate() error {
	return nil
}

// TODO: Implement ObtainNewCertificate when the registration of user is none.
func (ac *ACME) ObtainNewCertificate() error {
	return nil
}

// TODO: Implement ObtainRenewCertificate when the registration of user is not none.
func (ac *ACME) ObtainRenewCertificate() error {
	return nil
}
