package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/alvarolucio2007/GoSocial/internal/store"
)

func (app *application) BasicAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				app.unauthorizedBasicError(w, r, fmt.Errorf("authorization header is missing"))
				return
			}
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Basic" {
				app.unauthorizedBasicError(w, r, fmt.Errorf("authorization header is malformed"))
				return
			}
			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				app.unauthorizedBasicError(w, r, err)
				return
			}
			username := app.config.auth.basic.user
			pass := app.config.auth.basic.pass
			creds := strings.SplitN(string(decoded), ":", 2)
			if len(creds) != 2 || creds[0] != username || creds[1] != pass {
				app.unauthorizedBasicError(w, r, fmt.Errorf("invalid credentials"))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (app *application) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			app.unauthorizedError(w, r, fmt.Errorf("authorization header is missing"))
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			app.unauthorizedError(w, r, fmt.Errorf("authorization header is malformed"))
			return
		}
		token := parts[1]
		pasetoToken, err := app.authenticator.VerifyToken(token)
		if err != nil {
			app.unauthorizedError(w, r, err)
			return
		}
		ctx := r.Context()
		user, err := app.storage.Users.GetByEmail(ctx, pasetoToken.Email)
		if err != nil {
			app.unauthorizedError(w, r, err)
			return
		}
		ctx = context.WithValue(ctx, userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) checkPostOwnership(role string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := getUserFromContext(r)
		post := getPostFromCtx(r)
		if post.UserID == user.ID {
			next.ServeHTTP(w, r)
			return
		}
		allowed, err := app.checkRolePrecedecence(r.Context(), user, role)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}
		if !allowed {
			app.forbiddenError(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (app *application) checkRolePrecedecence(ctx context.Context, user *store.User, roleName string) (bool, error) {
	role, err := app.storage.Roles.GetByName(ctx, roleName)
	if err != nil {
		return false, err
	}
	return user.Role.Level >= role.Level, nil
}
