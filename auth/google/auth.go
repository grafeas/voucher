package google

import (
	"context"
	"fmt"
	"net/http"

	"github.com/docker/distribution/reference"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/grafeas/voucher"
	"github.com/grafeas/voucher/auth"
)

const gcrScope = "https://www.googleapis.com/auth/cloud-platform"

// GoogleAuth wraps the Google OAuth2 code.
type gAuth struct {
}

// GetTokenSource gets the default oauth2.TokenSource for connecting to Google's,
// OAuth2 protected systems, based on the runtime environment, or returns error
// if there's an issue getting the token source.
func (a *gAuth) GetTokenSource(ctx context.Context, ref reference.Named) (oauth2.TokenSource, error) {
	source, err := google.DefaultTokenSource(ctx, gcrScope)
	if nil != err {
		err = fmt.Errorf("failed to get Google Auth token source: %s", err)
	}

	return source, err
}

// ToClient returns a new http.Client with the authentication details setup by
// Auth.GetTokenSource.
func (a *gAuth) ToClient(ctx context.Context, image reference.Named) (*http.Client, error) {
	if !a.IsForDomain(image) {
		return nil, auth.NewAuthError("does not match domain", image)
	}

	tokenSource, err := a.GetTokenSource(ctx, image)
	if nil != err {
		return nil, err
	}

	client := oauth2.NewClient(ctx, tokenSource)
	err = auth.UpdateIdleConnectionsTimeout(client)

	return client, err
}

// IsForDomain validates the domain part of the Named image reference
func (a *gAuth) IsForDomain(image reference.Named) bool {
	return "gcr.io" == reference.Domain(image)
}

// NewAuth returns a new voucher.Auth to access Google specific resources.
func NewAuth() voucher.Auth {
	return new(gAuth)
}
