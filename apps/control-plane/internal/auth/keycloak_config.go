package auth

import (
	"context"
	"log"

	"github.com/Provision-Labs/ProvisionHub/apps/control-plane/internal/config"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

func setupOAuth(cfg *config.Config) (*oidc.Provider, *oauth2.Config) {
	ctx := context.Background()

	// Keycloak auto-exposes /.well-known/openid-configuration
	issuer := cfg.Issuer
	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		log.Fatal("Failed to get provider:", err)
	}

	config := &oauth2.Config{
		ClientID:     cfg.ClientId,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURL,
		Scopes:       cfg.Scopes,
		Endpoint:     provider.Endpoint(), // Endpoint auto-configurable
	}

	return provider, config
}
