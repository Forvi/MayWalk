package oauth

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type OAuthConfig struct {
	YaClientID     string `env:"YA_CLIENT_ID,required"`
	YaClientSecret string `env:"YA_CLIENT_SECRET,required"`
	YaRedirectURL  string `env:"YA_REDIRECT_URL,required"`
}

func LoadConfig() (*OAuthConfig, error) {
	cfg := &OAuthConfig{}
	opts := env.Options{
		Prefix: "OAUTH_",
	}
	if err := env.ParseWithOptions(cfg, opts); err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}
	return cfg, nil
}
