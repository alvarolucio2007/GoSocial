package main

import (
	"net/http"
	"testing"

	"github.com/go-openapi/testify/require"
)

func TestGetUser(t *testing.T) {
	app := newTestApplication(t)
	mux := app.mount()
	t.Run("should not allow unauthenticated requests", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		require.NoError(t, err)
		rr := executeRequest(req, mux)
		require.Equal(t, http.StatusUnauthorized, rr.Code)
	})
	t.Run("should allow authenticated requests", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		require.NoError(t, err)
		rr := executeRequest(req, mux)
		require.Equal(t, http.StatusOK, rr.Code)
	})
	t.Run("should not allow negative inputs", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/v1/users/-1", nil)
		require.NoError(t, err)
		rr := executeRequest(req, mux)
		require.Equal(t, http.StatusBadRequest, rr.Code)
	})
}
