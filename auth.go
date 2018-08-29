package voucher

import (
	"context"
	"net/http"

	"github.com/docker/distribution/reference"
	"golang.org/x/oauth2"
)

// Auth is an interface that wraps an to an OAuth2 system, to simplify the path
// from having an image reference to getting access to the data that makes up
// that image from the registry it lives in.
type Auth interface {
	GetTokenSource(context.Context, reference.Named) (oauth2.TokenSource, error)
}

// AuthToClient takes a struct implementing Auth and returns a new http.Client
// with the authentication details setup by Auth.GetTokenSource.
func AuthToClient(ctx context.Context, auth Auth, image reference.Named) (*http.Client, error) {
	tokenSource, err := auth.GetTokenSource(ctx, image)
	if nil != err {
		return nil, err
	}

	client := oauth2.NewClient(ctx, tokenSource)
	return client, nil
}
