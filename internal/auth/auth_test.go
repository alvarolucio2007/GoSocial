package auth_test

import (
	"encoding/base64"
	"testing"

	"github.com/alvarolucio2007/GoSocial/internal/auth"
	"github.com/go-openapi/testify/require"
)

func TestGenerateSymmetricKey(t *testing.T) {
	key, err := auth.GenerateSymmetricKey()
	require.NoError(t, err)
	require.Len(t, key, 32)
}

func TestNewPasetoAuthenticator(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		symKey := "0wKIJsaVqOKCE9rWYnTZD8ASPn8Ks9YCXpQfwDjQ/b4="
		key, err := base64.StdEncoding.DecodeString(symKey)
		require.NoError(t, err)
		_, err = auth.NewPasetoAuthenticator(key)
		require.NoError(t, err)
	})
}
