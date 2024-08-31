package config

import (
	"errors"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	CertPath              string
	KongEndpoint          string
	DomainName            string
	RegistrationEmail     string
	ChallengeType         string // by default is dns
	ChallengeProvider     string // by default is cloudflare (only support cloudflare for now)
	CloudflareZoneReadKey string
	CloudflareDNSWriteKey string
}

var (
	ErrCloudflareZoneReadKeyNotSet = errors.New("CF_ZONE_API_TOKEN is not set")
	ErrCloudflareDNSWriteKeyNotSet = errors.New("CF_DNS_API_TOKEN is not set")
)

func Initialize() (*Config, error) {
	viper.AutomaticEnv()

	// Set default values
	viper.SetDefault("CERT_PATH", "./certs")
	viper.SetDefault("KONG_ENDPOINT", "http://localhost:8001")
	viper.SetDefault("DOMAIN_NAME", "")
	viper.SetDefault("REGISTRATION_EMAIL", "")
	viper.SetDefault("CHALLENGE_TYPE", "dns")
	// With the cloudflare-dns challenge provider, need to use the Cloudflare API to create the DNS record

	viper.SetDefault("CHALLENGE_PROVIDER", "cloudflare")
	viper.SetDefault("CF_ZONE_API_TOKEN", "")
	viper.SetDefault("CF_DNS_API_TOKEN", "")

	requiredEnv := []string{
		"DOMAIN_NAME",
		"REGISTRATION_EMAIL",
		"CF_ZONE_API_TOKEN",
		"CF_DNS_API_TOKEN",
	}

	for _, env := range requiredEnv {
		if viper.GetString(env) == "" {
			return nil, errors.New(env + " is not set")
		}
	}

	// check directory cert path is exist or not
	if _, err := os.Stat(viper.GetString("CERT_PATH")); os.IsNotExist(err) {
		// if not exist, create the directory
		if err := os.Mkdir(viper.GetString("CERT_PATH"), 0755); err != nil {
			return nil, err
		}
	}

	config := &Config{
		CertPath:              viper.GetString("CERT_PATH"),
		KongEndpoint:          viper.GetString("KONG_ENDPOINT"),
		DomainName:            viper.GetString("DOMAIN_NAME"),
		RegistrationEmail:     viper.GetString("REGISTRATION_EMAIL"),
		ChallengeType:         viper.GetString("CHALLENGE_TYPE"),
		ChallengeProvider:     viper.GetString("CHALLENGE_PROVIDER"),
		CloudflareZoneReadKey: viper.GetString("CF_ZONE_API_TOKEN"),
		CloudflareDNSWriteKey: viper.GetString("CF_DNS_API_TOKEN"),
	}

	return config, nil
}
