package main

import (
	"context"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	htransport "google.golang.org/api/transport/http"
)

const idTokenKey = "id_token"

// idTokenSource is an oauth2.TokenSource that wraps another
// It takes the id_token from TokenSource and passes that on as a bearer token
// Implementation from: https://github.com/googleapis/google-api-go-client/issues/873
type idTokenSource struct {
	TokenSource oauth2.TokenSource
}

func (s *idTokenSource) Token() (*oauth2.Token, error) {
	token, err := s.TokenSource.Token()
	if err != nil {
		return nil, err
	}

	idToken, ok := token.Extra(idTokenKey).(string)
	if !ok {
		return nil, fmt.Errorf("token did not contain an id_token")
	}

	return &oauth2.Token{
		AccessToken: idToken,
		TokenType:   "Bearer",
		Expiry:      token.Expiry,
	}, nil
}

func getDefaultTokenSourceClient(ctx context.Context) (*http.Client, error) {
	src, err := google.DefaultTokenSource(ctx)
	if err != nil {
		return nil, fmt.Errorf("error creating token source: %w", err)
	}
	ts := oauth2.ReuseTokenSource(nil, &idTokenSource{TokenSource: src})

	transport, err := htransport.NewTransport(ctx, http.DefaultTransport, option.WithTokenSource(ts))
	if err != nil {
		return nil, fmt.Errorf("error creating client: %w", err)
	}
	return &http.Client{Transport: transport}, nil
}
