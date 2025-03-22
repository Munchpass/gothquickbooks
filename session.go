package gothquickbooks

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/markbates/goth"
	"golang.org/x/oauth2"
)

// Session stores data during the auth process with quickbooks.
type Session struct {
	AuthURL string
	// Each token is associated with a Quickbooks realm
	AccessToken string
	// Each token is associated with a Quickbooks realm
	RefreshToken          string
	AccessTokenExpiresAt  time.Time
	RefreshTokenExpiresAt time.Time

	// Each token is associated with a Quickbooks realm
	IDToken string

	// Only specified if IDToken is specified
	ParsedIDToken *IDToken
}

// GetAuthURL will return the URL set by calling the `BeginAuth` function on the
// quickbooks provider.
func (s Session) GetAuthURL() (string, error) {
	if s.AuthURL == "" {
		return "", errors.New(goth.NoAuthUrlErrorMessage)
	}
	return s.AuthURL, nil
}

// decodeJWTPart decodes a base64url encoded part of the JWT token
func decodeJWTPart(part string) ([]byte, error) {
	// Add padding if necessary
	if len(part)%4 != 0 {
		part += strings.Repeat("=", 4-len(part)%4)
	}

	// Replace URL encoding specific characters
	part = strings.ReplaceAll(part, "-", "+")
	part = strings.ReplaceAll(part, "_", "/")

	// Decode base64
	decodedBytes, err := base64.StdEncoding.DecodeString(part)
	if err != nil {
		return nil, err
	}

	return decodedBytes, nil
}

// Authorize completes the authorization with quickbooks and returns the access
// token to be stored for future use.
func (s *Session) Authorize(provider goth.Provider, params goth.Params) (string, error) {
	p := provider.(*Provider)
	token, err := p.config.Exchange(context.Background(), params.Get("code"), oauth2.SetAuthURLParam("code_verifier", params.Get("code_verifier")))
	if err != nil {
		return "", err
	}
	s.AccessToken = token.AccessToken
	s.RefreshToken = token.RefreshToken
	s.AccessTokenExpiresAt = token.Expiry
	if idToken := token.Extra("id_token"); idToken != nil {
		s.IDToken = idToken.(string)

		var idToken IDToken
		if s.IDToken != "" {
			// Split the token into its three parts: header, payload, and signature
			parts := strings.Split(s.IDToken, ".")
			if len(parts) != 3 {
				return "", errors.New("invalid JWT ID token format. Expected format: header.payload.signature")
			}
			// Decode only the payload
			rawPayload, err := decodeJWTPart(parts[1])
			if err != nil {
				return "", fmt.Errorf("failed to decode JWT ID token: %w", err)
			}

			err = json.Unmarshal(rawPayload, &idToken)
			if err == nil {
				s.ParsedIDToken = &idToken
			}
		}
	}

	if refreshTokenExpiresIn := token.Extra("x_refresh_token_expires_in"); refreshTokenExpiresIn != nil {
		s.RefreshTokenExpiresAt = time.Now().Add(time.Duration(refreshTokenExpiresIn.(float64)) * time.Second)
	}

	return token.AccessToken, err
}

// Marshal marshals a session into a JSON string.
func (s Session) Marshal() string {
	j, _ := json.Marshal(s)
	return string(j)
}

// String is equivalent to Marshal.  It returns a JSON representation of the session.
func (s Session) String() string {
	return s.Marshal()
}

// UnmarshalSession will unmarshal a JSON string into a session.
func (p *Provider) UnmarshalSession(data string) (goth.Session, error) {
	s := Session{}
	err := json.Unmarshal([]byte(data), &s)
	return &s, err
}
