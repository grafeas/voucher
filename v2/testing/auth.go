package vtesting

import (
	"context"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/docker/distribution/reference"
	"golang.org/x/oauth2"

	voucher "github.com/grafeas/voucher/v2"
	"github.com/grafeas/voucher/v2/auth"
)

type testTokenSource struct {
}

func (tkSrc *testTokenSource) Token() (*oauth2.Token, error) {
	return &oauth2.Token{
		AccessToken:  "abcd",
		TokenType:    "Bearer",
		RefreshToken: "",
		Expiry:       time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC),
	}, nil
}

type testAuth struct {
	server *httptest.Server
}

// GetTokenSource gets the default oauth2.TokenSource for connecting to OAuth2
// protected systems, based on the runtime environment, or returns error if there's
// an issue getting the token source.
func (a *testAuth) GetTokenSource(ctx context.Context, ref reference.Named) (oauth2.TokenSource, error) {
	return new(testTokenSource), nil
}

// ToClient returns a new http.Client with the authentication details setup by
// Auth.GetTokenSource.
func (a *testAuth) ToClient(ctx context.Context, image reference.Named) (*http.Client, error) {
	if !a.IsForDomain(image) {
		return nil, auth.NewAuthError("does not match domain", image)
	}
	tokenSource, err := a.GetTokenSource(ctx, image)
	if nil != err {
		return nil, err
	}

	client := oauth2.NewClient(ctx, tokenSource)
	err = UpdateClient(client, a.server)

	return client, err
}

func (a *testAuth) IsForDomain(image reference.Named) bool {
	return reference.Domain(image) == "localhost"
}

// NewAuth creates a new Auth suitable for testing with.
func NewAuth(server *httptest.Server) voucher.Auth {
	auth := new(testAuth)
	auth.server = server
	return auth
}
