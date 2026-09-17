package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alvarolucio2007/GoSocial/internal/auth"
	"github.com/alvarolucio2007/GoSocial/internal/store"
	"go.uber.org/zap"
)

func newTestApplication(t *testing.T) *application {
	t.Helper()
	logger := zap.Must(zap.NewProduction()).Sugar()

	mapPost := make(map[int64]*store.Post)
	mapUser := make(map[int64]*store.User)
	mapComments := make(map[int64]*store.Comment)
	mapFollowers := make(map[store.FollowKey]struct{})
	mapRoles := make(map[int]*store.Role)
	mockStorage := store.NewMockStorage(mapPost, mapUser, mapComments, mapFollowers, mapRoles)

	mockAuth := auth.NewMockAuthenticator()

	return &application{
		logger:        logger,
		storage:       mockStorage,
		authenticator: mockAuth,
	}
}

func executeRequest(req *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	return rr
}
