package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alvarolucio2007/GoSocial/internal/auth"
	"github.com/alvarolucio2007/GoSocial/internal/store"
	"github.com/alvarolucio2007/GoSocial/internal/store/cache"
	"go.uber.org/zap"
)

func newTestApplication(t *testing.T) *application {
	t.Helper()
	logger := zap.NewNop().Sugar()

	mapPost := make(map[int]*store.Post)
	mapUser := make(map[int]*store.User)
	mapComments := make(map[int]*store.Comment)
	mapFollowers := make(map[int]*store.Follower)
	mapRoles := make(map[int]*store.Role)
	mockStorage := store.NewMockStorage(mapPost, mapUser, mapComments, mapFollowers, mapRoles)

	mockCacheStorage := cache.NewMockCache(cache.MockCacheStore{})

	mockAuth := auth.NewMockAuthenticator()
	return &application{
		logger:        logger,
		storage:       mockStorage,
		cacheStorage:  mockCacheStorage,
		authenticator: mockAuth,
	}
}

func executeRequest(req *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	return rr
}
