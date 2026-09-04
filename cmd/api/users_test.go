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

func TestReadUserHandler(t *testing.T) {
	mockUser := &store.User{ID: 1, Username: "Test", Email: "Test@gmail.com"}
	mockUser.Password.Set("password")
	mockStorage := store.NewMockStorage(nil, map[int]*store.User{1: mockUser}, nil)
	app := &application{storage: store.Storage(mockStorage)}
	req := httptest.NewRequest("GET", "/v1/users/1", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("userID", "1")
	userContext := &store.User{ID: 1, Username: "Test", Email: "Test@gmail.com"}
	userContext.Password.Set("password")
	req = req.WithContext(context.WithValue(req.Context(), userCtx, userContext))
	w := httptest.NewRecorder()
	app.readUserHandler(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	var user struct {
		Data store.User `json:"data"`
	}
	err := json.NewDecoder(w.Body).Decode(&user)
	require.NoError(t, err)
	require.Equal(t, "Test", user.Data.Username)
	require.Equal(t, "Test@gmail.com", user.Data.Email)
}

func TestEditUserHandler(t *testing.T) {
	mockUser := &store.User{ID: 1, Username: "Test", Email: "test@gmail.com"}
	mockUser.Password.Set("...")
	mockStorage := store.NewMockStorage(nil, map[int]*store.User{
		1: mockUser,
	}, nil)
	app := &application{storage: store.Storage(mockStorage)}
	body := `{"username":"TEST2","email":"test2@gmail.com","password":"12345"}`
	req := httptest.NewRequest("PUT", "/v1/users/1", strings.NewReader(body))
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("userID", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	app.updateUserHandler(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	user, err := mockStorage.Users.Read(context.TODO(), 1)
	require.NoError(t, err)
	require.Equal(t, "TEST2", user.Username)
	require.Equal(t, "test2@gmail.com", user.Email)
	samePassword, err := user.Password.Compare("12345")
	require.NoError(t, err)
	require.True(t, samePassword)
}

func TestDeleteUserHandler(t *testing.T) {
	mockStorage := store.NewMockStorage(nil,
		map[int]*store.User{
			1: {ID: 1, Username: "TEST", Email: "TEST@GMAIL.COM"},
		}, nil)
	app := &application{storage: store.Storage(mockStorage)}
	req := httptest.NewRequest("DELETE", "/v1/users/1", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("userID", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	app.deleteUserHandler(w, req)
	require.Equal(t, http.StatusNoContent, w.Code)
	user, err := mockStorage.Users.Read(context.TODO(), 1)
	require.ErrorIs(t, err, store.ErrUserNotFound)
	require.Nil(t, user)
}
