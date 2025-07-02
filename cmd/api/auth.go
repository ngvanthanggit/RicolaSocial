package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/ngvanthanggit/RicolaSocial/internal/mailer"
	"github.com/ngvanthanggit/RicolaSocial/internal/store"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

type UserWithToken struct {
	*store.User
	Token string
}

// registerUserHandler godoc
//
//	@Summary		Register a user
//	@Description	Register a new user with username, email and password
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
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// create a temporary user with active status is false
	new_user := &store.User{
		Username: payload.Username,
		Email:    payload.Email,
	}

	// hash the password
	if err := new_user.Password.Set(payload.Password); err != nil {
		app.internalServerResponse(w, r, err)
		return
	}

	ctx := r.Context()

	// create new plain token
	plainToken := uuid.New().String()
	// store token into database
	hash := sha256.Sum256([]byte(plainToken))
	hashedToken := hex.EncodeToString(hash[:])

	// store user
	err := app.store.Users.CreateAndInvite(ctx, new_user, hashedToken, app.config.mail.exp)
	if err != nil {
		switch err {
		case store.ErrDuplicateEmail:
			app.badRequestResponse(w, r, err)
		case store.ErrDuplicateUsername:
			app.badRequestResponse(w, r, err)
		default:
			app.internalServerResponse(w, r, err)
		}
		return
	}

	// send email
	isProdEnv := app.config.env == "production"

	activationURL := fmt.Sprintf("%s/confirm/%s", app.config.frontendURl, plainToken)
	vars := struct {
		Username      string
		ActivationURL string
	}{
		Username:      new_user.Username,
		ActivationURL: activationURL,
	}

	err = app.mailer.Send(mailer.UserInvitationTemplate, new_user.Username, new_user.Email, vars, !isProdEnv)
	if err != nil {
		app.logger.Errorw("error sending invitation email", "error", err)

		// roll back user creation if failed to send email (SAGA pattern)
		if err := app.store.Users.Delete(ctx, new_user.ID); err != nil {
			app.logger.Errorw("error deletign user", "error", err)
		}

		app.internalServerResponse(w, r, err)
		return
	}

	userWithToken := &UserWithToken{
		User:  new_user,
		Token: plainToken,
	}

	if err := app.jsonResponse(w, http.StatusCreated, userWithToken); err != nil {
		app.internalServerResponse(w, r, err)
		return
	}
}
