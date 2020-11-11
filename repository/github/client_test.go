package github

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grafeas/voucher/repository"
)

func TestNewClient(t *testing.T) {
	var auth *repository.Auth

	t.Run("Test with valid auth method", func(t *testing.T) {
		auth = &repository.Auth{
			Token: "asdf1234",
		}

		client, err := NewClient(context.Background(), auth)

		require.NoError(t, err)
		assert.NotNil(t, client)
		assert.Implements(t, (*repository.Client)(nil), client)
	})

	t.Run("Test with invalid auth method", func(t *testing.T) {
		auth = &repository.Auth{
			Username: "user",
			Password: "pass",
		}

		client, err := NewClient(context.Background(), auth)

		require.Error(t, err)
		assert.Nil(t, client)
		assert.Equal(t, "unsupported auth type: userpassword", err.Error())
	})

	t.Run("Test with nil auth method", func(t *testing.T) {
		auth = nil

		client, err := NewClient(context.Background(), auth)

		require.Error(t, err)
		assert.Nil(t, client)
		assert.Equal(t, "must provide authentication", err.Error())
	})
}
