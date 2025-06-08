package main

import (
	"net/http"

	"github.com/ngvanthanggit/RicolaSocial/internal/store"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

// registerUserHandler godoc
//	@Summary		Register a user
//	@Description	Register a new user with username, email and password
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		RegisterUserPayload	true	"User credentials"
//	@Success		201		{object}	store.User			"User registered"
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

	// store user

	if err := app.jsonResponse(w, http.StatusCreated, nil); err != nil {
		app.internalServerResponse(w, r, err)
		return
	}
}
