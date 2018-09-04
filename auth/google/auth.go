package google

import (
	"context"
	"fmt"

	"github.com/Shopify/voucher"
	"github.com/docker/distribution/reference"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// GoogleAuth wraps the Google OAuth2 code.
type auth struct {
}

// GetTokenSource gets the default oauth2.TokenSource for connecting to Google's,
// OAuth2 protected systems, based on the runtime environment, or returns error
// if there's an issue getting the token source.
func (a *auth) GetTokenSource(ctx context.Context, reference reference.Named) (oauth2.TokenSource, error) {
	repository := "repository:" + reference.String() + ":*"
	source, err := google.DefaultTokenSource(ctx, repository)
	if nil != err {
		err = fmt.Errorf("failed to get Google Auth token source: %s", err)
	}

	return source, err

}

// NewAuth returns a new voucher.Auth to access Google specific resources.
func NewAuth() voucher.Auth {
	return new(auth)
}
