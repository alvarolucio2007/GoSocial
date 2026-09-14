package auth_test

import (
	"encoding/base64"
	"testing"
	"time"

	"github.com/alvarolucio2007/GoSocial/internal/auth"
	"github.com/go-openapi/testify/require"
)

func TestPaseto(t *testing.T) {
	keyB64 := "0wKIJsaVqOKCE9rWYnTZD8ASPn8Ks9YCXpQfwDjQ/b4="
	key, err := base64.StdEncoding.DecodeString(keyB64)
	require.NoError(t, err)
	pasAuth, err := auth.NewPasetoAuthenticator(key)
	require.NoError(t, err)
	t.Run("happy path", func(t *testing.T) {
		claims, err := auth.NewClaims("test@gmail.com", 5*time.Minute, 1)
		require.NoError(t, err)
		token, err := pasAuth.CreateToken(*claims)
		require.NoError(t, err)
		verToken, err := pasAuth.VerifyToken(token)
		require.NoError(t, err)
		require.EqualValues(t, claims, verToken)
	})
	t.Run("expired token", func(t *testing.T) {
		claims, err := auth.NewClaims("test@gmail.com", 0, 1)
		require.NoError(t, err)
		token, err := pasAuth.CreateToken(*claims)
		require.NoError(t, err)
		_, err = pasAuth.VerifyToken(token)
		require.ErrorIs(t, err, auth.ErrExpiredToken)
	})
	t.Run("invalid token", func(t *testing.T) {
		claims, err := auth.NewClaims("test@gmail.com", 1*time.Minute, 1)
		require.NoError(t, err)
		token, err := pasAuth.CreateToken(*claims)
		require.NoError(t, err)
		infectedToken := token + "1"
		_, err = pasAuth.VerifyToken(infectedToken)
		require.ErrorIs(t, err, auth.ErrInvalidToken)
	})
}
