package voucher

import (
	"context"
	"errors"
	"net/http"

	"github.com/docker/distribution/reference"
	"golang.org/x/oauth2"
)

// ErrNoAuth should be returned when something that depends on an Auth does not
// have one.
var ErrNoAuth = errors.New("no configured Auth")

// Auth is an interface that wraps an to an OAuth2 system, to simplify the path
// from having an image reference to getting access to the data that makes up
// that image from the registry it lives in.
type Auth interface {
	GetTokenSource(context.Context, reference.Named) (oauth2.TokenSource, error)
	ToClient(ctx context.Context, image reference.Named) (*http.Client, error)
}

// AuthToClient takes a struct implementing Auth and returns a new http.Client
// with the authentication details setup by Auth.GetTokenSource.
//
// DEPRECATED: This function has been superceded by Auth.ToClient. This function
// now calls that method directly.
func AuthToClient(ctx context.Context, auth Auth, image reference.Named) (*http.Client, error) {
	return auth.ToClient(ctx, image)
}
