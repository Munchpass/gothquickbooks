// Package quickbooks implements the OAuth protocol for authenticating users through quickbooks.
// This package can be used as a reference implementation of an OAuth provider for Goth.
package gothquickbooks

import (
	"context"
	"fmt"
	"net/http"

	"github.com/markbates/goth"
	"golang.org/x/oauth2"
)

type quickbooksGlobalConfig struct {
	IsProd bool
}

var globalConfig = &quickbooksGlobalConfig{IsProd: true}

func SetIsProd(isProd bool) {
	globalConfig.IsProd = isProd
}

// For the time being, the urls are the same are the same.
const (
	sandboxAuthURL     string = "https://appcenter.intuit.com/connect/oauth2"
	sandboxTokenURL    string = "https://oauth.platform.intuit.com/oauth2/v1/tokens/bearer"
	productionAuthURL  string = "https://appcenter.intuit.com/connect/oauth2"
	productionTokenURL string = "https://oauth.platform.intuit.com/oauth2/v1/tokens/bearer"
	revocationURL      string = "https://developer.api.intuit.com/v2/oauth2/tokens/revoke"
)

// New creates a new quickbooks provider, and sets up important connection details.
// You should always call `quickbooks.New` to get a new Provider. Never try to create
// one manually.
func New(clientKey, secret, callbackURL string, scopes ...string) *Provider {
	p := &Provider{
		ClientKey:    clientKey,
		Secret:       secret,
		CallbackURL:  callbackURL,
		providerName: "quickbooks",
	}
	p.config = newConfig(p, scopes)
	fmt.Printf("config: %+v\n", p.config)
	return p
}

// Provider is the implementation of `goth.Provider` for accessing quickbooks.
type Provider struct {
	ClientKey    string
	Secret       string
	CallbackURL  string
	HTTPClient   *http.Client
	config       *oauth2.Config
	providerName string
}

// Name is the name used to retrieve this provider later.
func (p *Provider) Name() string {
	return p.providerName
}

// SetName is to update the name of the provider (needed in case of multiple providers of 1 type)
func (p *Provider) SetName(name string) {
	p.providerName = name
}

func (p *Provider) Client() *http.Client {
	return goth.HTTPClientWithFallBack(p.HTTPClient)
}

// Debug is a no-op for the quickbooks package.
func (p *Provider) Debug(debug bool) {}

// BeginAuth asks quickbooks for an authentication end-point.
func (p *Provider) BeginAuth(state string) (goth.Session, error) {
	url := p.config.AuthCodeURL(state)
	session := &Session{
		AuthURL: url,
	}
	return session, nil
}

// FetchUser will go to quickbooks and access basic information about the user.
func (p *Provider) FetchUser(session goth.Session) (goth.User, error) {
	s := session.(*Session)
	user := goth.User{
		AccessToken:  s.AccessToken,
		Provider:     p.Name(),
		RefreshToken: s.RefreshToken,
		ExpiresAt:    s.ExpiresAt,
		UserID:       s.UserID,
	}

	if user.AccessToken == "" {
		// data is not yet retrieved since accessToken is still empty
		return user, fmt.Errorf("%s cannot get user information without accessToken", p.providerName)
	}

	return user, nil
}

func newConfig(provider *Provider, scopes []string) *oauth2.Config {
	authURL := productionAuthURL
	tokenURL := productionTokenURL
	if !globalConfig.IsProd {
		authURL = sandboxAuthURL
		tokenURL = sandboxTokenURL
	}

	c := &oauth2.Config{
		ClientID:     provider.ClientKey,
		ClientSecret: provider.Secret,
		RedirectURL:  provider.CallbackURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
		Scopes: scopes,
	}
	if len(c.Scopes) == 0 {
		c.Scopes = []string{ScopeOpenID, ScopeProfile, ScopeEmail, ScopePhone} // default scopes are used if none are provided in the constructor
	}

	return c
}

// RefreshToken get new access token based on the refresh token
func (p *Provider) RefreshToken(refreshToken string) (*oauth2.Token, error) {
	token := &oauth2.Token{RefreshToken: refreshToken}
	ts := p.config.TokenSource(context.Background(), token)
	newToken, err := ts.Token()
	if err != nil {
		return nil, err
	}
	return newToken, err
}

// RefreshTokenAvailable refresh token is not provided by quickbooks
func (p *Provider) RefreshTokenAvailable() bool {
	return true
}
