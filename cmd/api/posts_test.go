package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alvarolucio2007/GoSocial/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestPOSTPostHandler(t *testing.T) {
	mockStorage := store.NewMockStorage(map[int]*store.Post{}, nil)
	app := &application{storage: store.Storage(mockStorage)}
	body := `{"title":"TEST","content":"Content","tags":["test"]}`
	req := httptest.NewRequest("POST", "/v1/posts", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	app.createPostHandler(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
	var post store.Post
	err := json.NewDecoder(w.Body).Decode(&post)
	require.NoError(t, err)
	require.Equal(t, "TEST", post.Title)
	require.Equal(t, "Content", post.Content)
	require.Equal(t, []string{"test"}, post.Tags)
}

func TestGetPostHandler(t *testing.T) {
	mockStorage := store.NewMockStorage(
		map[int]*store.Post{
			1: {ID: 1, Title: "Test", Content: "Content"},
		},
		nil,
	)
	app := &application{storage: store.Storage(mockStorage)}
	req := httptest.NewRequest("GET", "/v1/posts/1", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("postID", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()
	app.readPostHandler(w, req)
	require.Equal(t, w.Code, http.StatusCreated)
	var post store.Post
	err := json.NewDecoder(w.Body).Decode(&post)
	require.NoError(t, err)
	require.Equal(t, "Test", post.Title)
	require.Equal(t, "Content", post.Content)
}
