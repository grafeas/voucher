package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/grafeas/voucher/v2/client"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iamcredentials/v1"
	googleoauth2 "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
	htransport "google.golang.org/api/transport/http"
)

// idTokenSource is an oauth2.TokenSource that wraps another
// It takes the id_token from TokenSource and passes that on as a bearer token
// Implementation from: https://github.com/googleapis/google-api-go-client/issues/873
type idTokenSource struct {
	TokenSource oauth2.TokenSource
}

func (s *idTokenSource) Token() (*oauth2.Token, error) {
	tok, err := s.TokenSource.Token()
	if err != nil {
		return nil, err
	}

	// If the token includes an identity token, use that:
	if idTok, ok := extractID(tok); ok {
		return &oauth2.Token{
			AccessToken: idTok,
			TokenType:   "Bearer",
			Expiry:      tok.Expiry,
		}, nil
	}

	// If the token does not include an identity token, attempt to generate one:
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	idTok, err := generateIDToken(ctx, tok)
	if err != nil {
		return nil, fmt.Errorf("error generating id token: %w", err)
	}
	return &oauth2.Token{
		AccessToken: idTok,
		TokenType:   "Bearer",
		Expiry:      tok.Expiry,
	}, nil
}

const idTokenKey = "id_token"

func extractID(token *oauth2.Token) (ret string, ok bool) {
	ret, ok = token.Extra(idTokenKey).(string)
	return
}

// NewAuthClientWithToken creates an auth client using the token created from ADC
func NewAuthClientWithToken(ctx context.Context, voucherURL string) (*client.Client, error) {
	ts, err := getDefaultTokenSource(ctx)
	if err != nil {
		return nil, err
	}

	c, err := newHTTPClient(ctx, ts)
	if err != nil {
		return nil, err
	}
	return client.NewCustomClient(voucherURL, c)
}

func getDefaultTokenSource(ctx context.Context) (oauth2.TokenSource, error) {
	src, err := google.DefaultTokenSource(ctx)
	if err != nil {
		return nil, fmt.Errorf("error creating token source: %w", err)
	}

	ts := oauth2.ReuseTokenSource(nil, &idTokenSource{TokenSource: src})
	return ts, nil
}

func generateIDToken(ctx context.Context, tok *oauth2.Token) (string, error) {
	clientOpts := []option.ClientOption{option.WithTokenSource(oauth2.StaticTokenSource(tok))}

	// Get the serviceAccount from the current token
	googOauth, err := googleoauth2.NewService(ctx, clientOpts...)
	if err != nil {
		return "", fmt.Errorf("error creating google oauth client: %w", err)
	}
	tokenInfo, err := googOauth.Tokeninfo().Do()
	if err != nil {
		return "", fmt.Errorf("error fetching token info: %w", err)
	}

	// Generate an id token for that serviceAccount
	// This is expensive, but we are wrapped in a ReuseTokenSource() for caching.
	svc, err := iamcredentials.NewService(ctx, clientOpts...)
	if err != nil {
		return "", fmt.Errorf("error creating iamcredentials service: %w", err)
	}
	tokenName := fmt.Sprintf("projects/-/serviceAccounts/%s", tokenInfo.Email)
	generatedTok, err := svc.Projects.ServiceAccounts.GenerateIdToken(tokenName, &iamcredentials.GenerateIdTokenRequest{
		Audience: defaultConfig.Server,
	}).Do()
	if err != nil {
		return "", fmt.Errorf("error calling GenerateIdToken: %w", err)
	}
	return generatedTok.Token, nil
}

func newHTTPClient(ctx context.Context, ts oauth2.TokenSource) (*http.Client, error) {
	t, err := htransport.NewTransport(ctx, http.DefaultTransport, option.WithTokenSource(ts))
	if err != nil {
		return nil, fmt.Errorf("error creating client: %w", err)
	}
	return &http.Client{Transport: t}, nil
}
