package acme

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"

	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge/dns01"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/providers/dns"
	"github.com/go-acme/lego/v4/registration"
	"github.com/mynameismaxz/acme2kong/pkg/logger"
)

type ACME struct {
	CADirectoryUrl string
	Domain         string
	Email          string

	userAccount *User
	client      *lego.Client
	logger      *logger.Logger
}

func NewClient(caDirectoryUrl, domain, email string, l *logger.Logger) *ACME {
	// generate private key
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil
	}

	user := User{
		Email: email,
		key:   privateKey,
	}

	config := lego.NewConfig(&user)
	legoClient, err := lego.NewClient(config)
	if err != nil {
		return nil
	}

	return &ACME{
		CADirectoryUrl: caDirectoryUrl,
		Domain:         domain,
		Email:          email,
		userAccount:    &user,
		client:         legoClient,
		logger:         l,
	}
}

func (ac *ACME) Register() error {
	tmp := []string{"*.tha.mymacz.com"}
	nameservers := []string{"1.1.1.1", "8.8.8.8"}

	// register
	provider, err := dns.NewDNSChallengeProviderByName("cloudflare")
	if err != nil {
		return err
	}

	if err = ac.client.Challenge.SetDNS01Provider(
		provider,
		dns01.CondOption((len(tmp) > 0),
			dns01.AddRecursiveNameservers(dns01.ParseNameservers(nameservers)))); err != nil {
		return err
	}

	// check registration
	// TODO: Rewrite this part to check if the registration is already exist
	ac.logger.Info("Registering...")
	reg, err := ac.client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		return err
	}

	ac.userAccount.Registration = reg

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
