package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grafeas/voucher/v2/client"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	htransport "google.golang.org/api/transport/http"
)

// Implementation from: https://github.com/googleapis/google-api-go-client/issues/873
// idTokenSource is an oauth2.TokenSource that wraps another
// It takes the id_token from TokenSource and passes that on as a bearer token
type idTokenSource struct {
	TokenSource oauth2.TokenSource
}

func (s *idTokenSource) Token() (*oauth2.Token, error) {
	token, err := s.TokenSource.Token()
	if err != nil {
		return nil, err
	}

	idToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, fmt.Errorf("token did not contain an id_token")
	}

	return &oauth2.Token{
		AccessToken: idToken,
		TokenType:   "Bearer",
		Expiry:      token.Expiry,
	}, nil
}

func NewAuthClientWithToken(voucherURL string) (*client.Client, error) {
	c, err := newHttpClient()
	if err != nil {
		return nil, err
	}
	return client.NewCustomClient(voucherURL, c)
}

func getDefaultTokenSource() (oauth2.TokenSource, error) {
	src, err := google.DefaultTokenSource(context.Background())
	if err != nil {
		return nil, err
	}
	ts := oauth2.ReuseTokenSource(nil, &idTokenSource{TokenSource: src})
	return ts, nil
}

func newHttpClient() (*http.Client, error) {
	ts, err := getDefaultTokenSource()
	if err != nil {
		return nil, fmt.Errorf("error creating client: %w", err)
	}

	t, err := htransport.NewTransport(context.Background(), http.DefaultTransport, option.WithTokenSource(ts))
	if err != nil {
		return nil, fmt.Errorf("error creating client: %w", err)
	}
	return &http.Client{Transport: t}, nil
}
