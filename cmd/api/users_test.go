package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/alvarolucio2007/GoSocial/internal/auth"
	"github.com/alvarolucio2007/GoSocial/internal/store"
	"github.com/go-openapi/testify/require"
)

func TestGetUser(t *testing.T) {
	app := newTestApplication(t)
	mux := app.mount()
	user := store.User{ID: 1, Username: "Test", Email: "test@gmail.com"}
	err := app.storage.Users.Create(context.Background(), nil, &user)
	require.NoError(t, err)
	t.Run("should not allow unauthenticated requests", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		require.NoError(t, err)
		rr := executeRequest(req, mux)
		require.Equal(t, http.StatusUnauthorized, rr.Code)
	})
	t.Run("should allow authenticated requests and return same value", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		require.NoError(t, err)
		claimsToken, err := auth.NewClaims("test@gmail.com", 10*time.Second, 1)
		require.NoError(t, err)

		testToken, err := app.authenticator.CreateToken(*claimsToken)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+testToken)
		rr := executeRequest(req, mux)
		require.Equal(t, http.StatusOK, rr.Code)
		var userResponse struct {
			Data store.User `json:"data"`
		}
		err = json.NewDecoder(rr.Body).Decode(&userResponse)
		require.NoError(t, err)
		require.EqualValues(t, user, userResponse.Data)
	})
}

func TestUpdateUser(t *testing.T) {
	app := newTestApplication(t)
	mux := app.mount()
	user := store.User{ID: 1, Username: "Test", Email: "test@gmail.com"}
	err := app.storage.Users.Create(context.Background(), nil, &user)
	require.NoError(t, err)
	t.Run("should not allow unauthenticated requests", func(t *testing.T) {
		body := UpdateUserPayload{Username: "Updated", Email: "updated@gmail.com", Password: "Updated"}
		jsonBody, err := json.Marshal(body)
		require.NoError(t, err)
		req, err := http.NewRequest(http.MethodPut, "/v1/users/1", bytes.NewBuffer(jsonBody))
		require.NoError(t, err)
		rr := executeRequest(req, mux)
		require.Equal(t, http.StatusUnauthorized, rr.Code)
	})
	t.Run("should allow authenticated requests and edit the user", func(t *testing.T) {
		body := UpdateUserPayload{Username: "Updated", Email: "updated@gmail.com", Password: "Updated"}
		jsonBody, err := json.Marshal(body)
		require.NoError(t, err)
		req, err := http.NewRequest(http.MethodPut, "/v1/users/1", bytes.NewBuffer(jsonBody))
		require.NoError(t, err)

		claimsToken, err := auth.NewClaims("test@gmail.com", 10*time.Second, 1)
		require.NoError(t, err)
		testToken, err := app.authenticator.CreateToken(*claimsToken)
		require.NoError(t, err)

		req.Header.Set("Authorization", "Bearer "+testToken)
		rr := executeRequest(req, mux)
		var userResponse struct {
			Data store.User `json:"data"`
		}
		err = json.NewDecoder(rr.Body).Decode(&userResponse)
		require.NoError(t, err)
		editedUser, err := app.storage.Users.Read(context.Background(), 1)
		require.NoError(t, err)
		require.Equal(t, editedUser.Username, userResponse.Data.Username)
		require.Equal(t, editedUser.Email, userResponse.Data.Email)
		isSame, err := editedUser.Password.Compare("Updated")
		require.NoError(t, err)
		require.True(t, isSame)
	})
}
