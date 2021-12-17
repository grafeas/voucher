package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	goauth2 "google.golang.org/api/oauth2/v2"
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

	idTokenFromToken, ok := token.Extra(idTokenKey).(string)
	if !ok {
		return nil, fmt.Errorf("token did not contain an id_token")
	}

	return &oauth2.Token{
		AccessToken: idTokenFromToken,
		TokenType:   "Bearer",
		Expiry:      token.Expiry,
	}, nil
}

func (s *idTokenSource) HasIdToken() bool {
	// Check for idToken
	if token, err := s.Token(); err == nil {
		_, ok := token.Extra(idTokenKey).(string)
		return ok
	}
	return false
}

func getDefaultTokenSourceClient(ctx context.Context) (*http.Client, error) {
	src, err := google.DefaultTokenSource(ctx)
	if err != nil {
		return nil, fmt.Errorf("error creating token source: %w", err)
	}

	tokenSource := &idTokenSource{TokenSource: src}
	if !tokenSource.HasIdToken() {
		// Make Client to generate token
		// generate idtoken
		generateIdToken(ctx)
	}

	ts := oauth2.ReuseTokenSource(nil, tokenSource)
	transport, err := htransport.NewTransport(ctx, http.DefaultTransport, option.WithTokenSource(ts))
	if err != nil {
		return nil, fmt.Errorf("error creating client: %w", err)
	}
	return &http.Client{Transport: transport}, nil
}

func generateIdToken(ctx context.Context) {
	googleOauth2Service, err := goauth2.NewService(ctx)
	if err != nil {
		fmt.Errorf("connecting to google oauth2 service to generate idToken: %w", err)
	}
	tokenInfo, err := googleOauth2Service.Tokeninfo().Do()
	if err != nil {
		fmt.Errorf("get token info: %w", err)
	}

	requestBody := map[string]string{
		"audience": defaultConfig.Server,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Errorf("get id token into: %w", err)
	}

	client, err := google.DefaultClient(ctx, "https://www.googleapis.com/auth/iam")
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("https://iamcredentials.googleapis.com/v1/projects/-/serviceAccounts/%s:generateIdToken", tokenInfo.Email),
		bytes.NewBuffer(body),
	)
	if err != nil {
		fmt.Errorf("get id token info: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Errorf("get id token info: %w", err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
}
