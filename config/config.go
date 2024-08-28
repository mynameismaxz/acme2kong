package config

import "github.com/spf13/viper"

type Config struct {
	KongEndpoint      string
	DomainName        string
	RegistrationEmail string
	ChallengeType     string // by default is dns
	ChallengeProvider string // by default is cloudflare (only support cloudflare for now)
	CloudflareEmail   string
	CloudflareAPIKey  string
}

func Initialize() (*Config, error) {
	viper.AutomaticEnv()

	// Set default values
	viper.SetDefault("KONG_ENDPOINT", "http://localhost:8001")
	viper.SetDefault("DOMAIN_NAME", "*.mymacz.com")
	viper.SetDefault("REGISTRATION_EMAIL", "")
	viper.SetDefault("CHALLENGE_TYPE", "dns")
	viper.SetDefault("CHALLENGE_PROVIDER", "cloudflare")
	viper.SetDefault("CLOUDFLARE_EMAIL", "")
	viper.SetDefault("CLOUDFLARE_API_KEY", "")

	config := &Config{
		KongEndpoint:      viper.GetString("KONG_ENDPOINT"),
		DomainName:        viper.GetString("DOMAIN_NAME"),
		RegistrationEmail: viper.GetString("REGISTRATION_EMAIL"),
		ChallengeType:     viper.GetString("CHALLENGE_TYPE"),
		ChallengeProvider: viper.GetString("CHALLENGE_PROVIDER"),
		CloudflareEmail:   viper.GetString("CLOUDFLARE_EMAIL"),
		CloudflareAPIKey:  viper.GetString("CLOUDFLARE_API_KEY"),
	}

	return config, nil
}
