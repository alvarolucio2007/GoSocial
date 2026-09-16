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

func TestCreatePostHandler(t *testing.T) {
	app := newTestApplication(t)
	mux := app.mount()
	user := store.User{ID: 1, Username: "test", Email: "test@gmail.com"}
	err := app.storage.Users.Create(context.Background(), nil, &user)
	require.NoError(t, err)
	t.Run("should not allow unauthenticated access", func(t *testing.T) {
		body := CreatePostPayload{Title: "title", Content: "content", Tags: []string{"tag 1", "tag 2"}}
		jsonBody, err := json.Marshal(body)
		require.NoError(t, err)
		req, err := http.NewRequest(http.MethodPost, "/v1/posts", bytes.NewBuffer(jsonBody))
		require.NoError(t, err)

		rr := executeRequest(req, mux)
		require.Equal(t, http.StatusUnauthorized, rr.Code)
	})
	t.Run("should allow authenticated requests and creating posts", func(t *testing.T) {
		body := CreatePostPayload{Title: "title", Content: "content", Tags: []string{"tag 1", "tag 2"}}
		jsonBody, err := json.Marshal(body)
		require.NoError(t, err)
		req, err := http.NewRequest(http.MethodPost, "/v1/posts", bytes.NewBuffer(jsonBody))
		require.NoError(t, err)

		claimsToken, err := auth.NewClaims("test@gmail.com", 10*time.Second, 1)
		require.NoError(t, err)
		testToken, err := app.authenticator.CreateToken(*claimsToken)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+testToken)

		rr := executeRequest(req, mux)
		require.Equal(t, http.StatusCreated, rr.Code)

		var postResponse struct {
			Data store.Post `json:"data"`
		}
		err = json.NewDecoder(rr.Body).Decode(&postResponse)
		require.NoError(t, err)

		fetchedPost, err := app.storage.Posts.Read(context.Background(), int(postResponse.Data.ID))
		require.NoError(t, err)

		require.EqualValues(t, *fetchedPost, postResponse.Data)
	})
	t.Run("should not allow invalid payload", func(t *testing.T) {
		body := CreatePostPayload{Title: "title", Content: "", Tags: []string{}}
		jsonBody, err := json.Marshal(body)
		require.NoError(t, err)
		req, err := http.NewRequest(http.MethodPost, "/v1/posts", bytes.NewBuffer(jsonBody))
		require.NoError(t, err)

		claimsToken, err := auth.NewClaims("test@gmail.com", 10*time.Second, 1)
		require.NoError(t, err)
		testToken, err := app.authenticator.CreateToken(*claimsToken)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+testToken)

		rr := executeRequest(req, mux)
		require.Equal(t, http.StatusBadRequest, rr.Code)
	})
}

func TestReadPostHandler(t *testing.T) {
	app := newTestApplication(t)
	mux := app.mount()
	user := store.User{ID: 1, Username: "test", Email: "test@gmail.com"}
	err := app.storage.Users.Create(context.Background(), nil, &user)
	require.NoError(t, err)
	post := store.Post{ID: 1, Content: "test", Title: "Test", UserID: 1, Tags: []string{}, User: user}
	err = app.storage.Posts.Create(context.Background(), &post)
	require.NoError(t, err)
	t.Run("should allow authenticated requests and read posts correctly", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/v1/posts/1", nil)
		require.NoError(t, err)

		claimsToken, err := auth.NewClaims("test@gmail.com", 10*time.Second, 1)
		require.NoError(t, err)
		testToken, err := app.authenticator.CreateToken(*claimsToken)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+testToken)

		rr := executeRequest(req, mux)
		require.Equal(t, http.StatusOK, rr.Code)

		var postResponse struct {
			Data store.Post `json:"data"`
		}
		err = json.NewDecoder(rr.Body).Decode(&postResponse)
		require.NoError(t, err)

		require.EqualValues(t, post, postResponse.Data)
	})
}
