package oauth

import (
	"context"
	"encoding/json"
	"fmt"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/yandex"
)

type AuthorizedUser struct {
	ID           string   `json:"id"`
	Login        string   `json:"login"`
	ClientID     string   `json:"client_id"`
	DisplayName  string   `json:"display_name"`
	RealName     string   `json:"real_name"`
	FirstName    string   `json:"first_name"`
	LastName     string   `json:"last_name"`
	Sex          string   `json:"sex"`
	DefaultEmail string   `json:"default_email"`
	Emails       []string `json:"emails"`
	PSUID        string   `json:"psuid"`
}

type YaOAuth struct {
	config *oauth2.Config
}

func NewYaOAuth(cfg *OAuthConfig) *YaOAuth {
	return &YaOAuth{
		config: &oauth2.Config{
			ClientID:     cfg.YaClientID,
			ClientSecret: cfg.YaClientSecret,
			RedirectURL:  cfg.YaRedirectURL,
			Scopes:       []string{"login:info", "login:email"},
			Endpoint:     yandex.Endpoint,
		},
	}
}

func (auth *YaOAuth) GeneratePKCE() string {
	return oauth2.GenerateVerifier()
}

func (auth *YaOAuth) Authorize(ctx context.Context, code, verifier string) (*AuthorizedUser, error) {
	tok, err := auth.config.Exchange(ctx, code, oauth2.VerifierOption(verifier))
	if err != nil {
		return nil, fmt.Errorf("excange token: %w", err)
	}

	client := auth.config.Client(ctx, tok)
	resp, err := client.Get("https://login.yandex.ru/info?format=json")
	if err != nil {
		return nil, fmt.Errorf("api request: %w", err)
	}
	defer resp.Body.Close()

	var authorizedUser AuthorizedUser
	if err := json.NewDecoder(resp.Body).Decode(&authorizedUser); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &authorizedUser, nil
}

func (auth *YaOAuth) SignIn(ctx context.Context, verifier string) string {
	return auth.config.AuthCodeURL("state", oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(verifier))
}
