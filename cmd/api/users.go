package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/alvarolucio2007/GoSocial/internal/store"
	"github.com/go-chi/chi/v5"
)

type userKey string

const userCtx userKey = "user"

type CreateUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateUserPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
	}
	user := &store.User{
		Username: payload.Username,
		Email:    payload.Email,
		Password: payload.Password,
	}
	ctx := r.Context()
	if err := app.storage.Users.Create(ctx, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := app.jsonResponse(w, http.StatusCreated, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) readUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)
	if err := app.jsonResponse(w, http.StatusOK, &user); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

type UpdateUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (app *application) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "userID")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}
	var payload UpdateUserPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
	}
	user := &store.User{
		ID:       id,
		Username: payload.Username,
		Password: payload.Password,
		Email:    payload.Email,
	}
	ctx := r.Context()
	if err := app.storage.Users.Update(ctx, user); err != nil {
		switch {
		case errors.Is(err, store.ErrUserNotFound):
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}
	if err := app.jsonResponse(w, http.StatusOK, &user); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "userID")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}
	ctx := r.Context()
	if err := app.storage.Users.Delete(ctx, int(id)); err != nil {
		switch {
		case errors.Is(err, store.ErrPostNotFound):
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type FollowUser struct {
	UserID int64 `json:"user_id"`
}

func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	followerUser := getUserFromContext(r)
	var payload FollowUser
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	ctx := r.Context()
	if err := app.storage.Followers.Follow(ctx, followerUser.ID, payload.UserID); err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)
	if err := app.jsonResponse(w, http.StatusNoContent, &user); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "userID")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			app.badRequestError(w, r, err)
			return
		}
		ctx := r.Context()
		user, err := app.storage.Users.Read(ctx, int(id))
		if err != nil {
			switch {
			case errors.Is(err, store.ErrUserNotFound):
				app.notFoundError(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}
		ctx = context.WithValue(ctx, userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserFromContext(r *http.Request) *store.User {
	user, _ := r.Context().Value(userCtx).(*store.User)
	return user
}
