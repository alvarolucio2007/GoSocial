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

func TestCreateUserHandler(t *testing.T) {
	mockStorage := store.NewMockStorage(nil, map[int]*store.User{}, nil)
	app := &application{storage: store.Storage(mockStorage)}
	body := `{"username":"Test","email":"test@email.com","password":"password"}`
	req := httptest.NewRequest("POST", "/v1/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	app.createUserHandler(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
	var user store.User
	err := json.NewDecoder(w.Body).Decode(&user)
	require.NoError(t, err)
	require.Equal(t, "Test", user.Username)
	require.Equal(t, "test@email.com", user.Email)
}

func TestReadUserHandler(t *testing.T) {
	mockStorage := store.NewMockStorage(nil, map[int]*store.User{1: {ID: 1, Username: "Test", Email: "Test@gmail.com", Password: "password"}}, nil)
	app := &application{storage: store.Storage(mockStorage)}
	req := httptest.NewRequest("GET", "/v1/users/1", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("userID", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()
	app.readUserHandler(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	var user store.User
	err := json.NewDecoder(w.Body).Decode(&user)
	require.NoError(t, err)
	require.Equal(t, "Test", user.Username)
	require.Equal(t, "Test@gmail.com", user.Email)
}

func TestEditUserHandler(t *testing.T) {
	mockStorage := store.NewMockStorage(nil, map[int]*store.User{
		1: {ID: 1, Username: "Test", Email: "test@gmail.com", Password: "..."},
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
	require.Equal(t, "12345", user.Password)
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
