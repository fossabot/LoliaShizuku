package services

import (
	"context"
	"fmt"
	"strings"

	"golang.org/x/oauth2"
)

func loadOrRefreshOAuthToken(ctx context.Context) (*oauth2.Token, error) {
	token, err := LoadOAuthToken()
	if err != nil {
		return nil, fmt.Errorf("load oauth token: %w", err)
	}

	if token == nil || strings.TrimSpace(token.AccessToken) == "" {
		return nil, fmt.Errorf("oauth token is empty")
	}
	if token.Valid() {
		return token, nil
	}
	if strings.TrimSpace(token.RefreshToken) == "" {
		return nil, fmt.Errorf("oauth access token expired and no refresh token available")
	}

	oauthCfg, cfgErr := resolveOAuthConfig()
	if cfgErr != nil {
		return nil, fmt.Errorf("load oauth config for refresh: %w", cfgErr)
	}

	refreshedToken, refreshErr := oauthCfg.TokenSource(ctx, token).Token()
	if refreshErr != nil {
		return nil, fmt.Errorf("refresh oauth token: %w", refreshErr)
	}

	if saveErr := SaveOAuthToken(refreshedToken); saveErr != nil {
		return nil, fmt.Errorf("save refreshed oauth token: %w", saveErr)
	}

	return refreshedToken, nil
}
