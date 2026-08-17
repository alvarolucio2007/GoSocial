package main

import (
	"net/http"

	"github.com/alvarolucio2007/GoSocial/internal/store"
)

type CreateUserPayload struct {
	Username string
	Email    string
	Password string
}

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateUserPayload
	if err := readJSON(w, r, &payload); err != nil {
		_ = writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	user := &store.User{
		Username: payload.Username,
		Email:    payload.Email,
		Password: payload.Password,
	}
	ctx := r.Context()
	if err := app.storage.Users.Create(ctx, user); err != nil {
		_ = writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := writeJSON(w, http.StatusCreated, user); err != nil {
		_ = writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
