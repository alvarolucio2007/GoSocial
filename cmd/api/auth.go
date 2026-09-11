package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"uuid"

	"github.com/alvarolucio2007/GoSocial/internal/auth"
	"github.com/alvarolucio2007/GoSocial/internal/mailer"
	"github.com/alvarolucio2007/GoSocial/internal/store"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"email,required,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

type UserWithToken struct {
	*store.User
	Token string `json:"token"`
}

// registerUserHandler godoc
//
//	@Summary		Registers an user
//	@Description	Registers an user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		RegisterUserPayload	true	"User credentials"
//	@Success		201		{object}	UserWithToken		"User registered"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Router			/authentication/user [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	user := &store.User{
		Username: payload.Username,
		Email:    payload.Email,
	}
	if err := user.Password.Set(payload.Password); err != nil {
		app.internalServerError(w, r, err)
		return
	}
	ctx := r.Context()
	plainToken := uuid.New().String()
	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])
	err := app.storage.Users.CreateAndInvite(ctx, user, hashToken, app.config.mail.exp)
	if err != nil {
		switch err {
		case store.ErrDuplicateEmail:
			app.badRequestError(w, r, err)
		case store.ErrDuplicateUsername:
			app.badRequestError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}
	userWithToken := UserWithToken{
		User:  user,
		Token: plainToken,
	}
	isProdEnv := app.config.env == "production"
	activationURL := fmt.Sprintf("%s/confirm/%s", app.config.frontendURL, plainToken)
	vars := struct {
		Username      string
		ActivationURL string
	}{
		Username:      user.Username,
		ActivationURL: activationURL,
	}

	status, err := app.mailer.Send(mailer.UserWelcomeTemplate, user.Username, user.Email, vars, !isProdEnv)
	if err != nil {
		app.logger.Errorw("error sending welcome email", "error", err)
		if err := app.storage.Users.Delete(ctx, user.ID); err != nil {
			app.logger.Errorw("error deleting user", "error", err)
		}
		app.internalServerError(w, r, err)
		return
	}
	app.logger.Infow("Email sent", "status code", status)
	if err := app.jsonResponse(w, http.StatusCreated, userWithToken); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

type CreateUserTokenPayload struct {
	Email    string `json:"email" validate:"required,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

// createTokenHandler godoc
//
//	@Summary		Creates a token
//	@Description	Creates a token for the user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreateUserTokenPayload	true	"User credentials"
//	@Success		200		{object}	string					"User registered"
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Router			/authentication/token [post]
func (app *application) createTokenHandler(w http.ResponseWriter, r *http.Request) {
	// parse payload cred
	var payload CreateUserTokenPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
	}
	// fetch user
	user, err := app.storage.Users.GetByEmail(r.Context(), payload.Email)
	if err != nil {
		switch err {
		case store.ErrUserNotFound:
			app.unauthorizedError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}
	isSame, err := user.Password.Compare(payload.Password)
	if err != nil {
		app.unauthorizedError(w, r, err)
		return
	}
	if !isSame {
		app.unauthorizedError(w, r, fmt.Errorf("unauthorized"))
	}
	// generate token
	claims, err := auth.NewClaims(user.Email, app.config.auth.token.exp) // TODO: Implement custom durations later
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	token, err := app.authenticator.CreateToken(*claims)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := app.jsonResponse(w, http.StatusCreated, token); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
