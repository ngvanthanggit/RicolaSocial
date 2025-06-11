package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/ngvanthanggit/RicolaSocial/internal/store"
)

type userKey string

var userCtx userKey = "user"

func (app *application) usersContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// get the user id from the URL
		userIDParam := chi.URLParam(r, "userID")
		userID, err := strconv.ParseInt(userIDParam, 10, 64)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		// retrieve the user from database
		user, err := app.store.Users.GetUserById(ctx, userID)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerResponse(w, r, err)
			}
			return
		}

		// pass the user to the request context
		ctx = context.WithValue(ctx, userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserFromCtx(r *http.Request) *store.User {
	return r.Context().Value(userCtx).(*store.User)
}

// GetUser godoc
//
//	@Summary		Fetches a user profile
//	@Description	Fetches a user profile by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	store.User
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/{id} [get]
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerResponse(w, r, err)
		return
	}
}

type FollowUser struct {
	UserID int64 ` json:"user_id" validate:"required"`
}

// FollowUser godoc
//
//	@Summary		Follow a user
//	@Description	Follow a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int		true	"User ID"
//	@Success		204		{string}	string	"User followed"
//	@Failure		400		{object}	error	"User payload missing"
//	@Failure		404		{object}	error	"User not found"
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/follow [put]
func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)

	// TODO: need to change to get the current user id instead of reading from request like this
	var payload FollowUser
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if payload.UserID == user.ID {
		app.badRequestResponse(w, r, fmt.Errorf("not allowed to follow your own profile"))
		return
	}

	// the user sending request follows the user from context
	if err := app.store.Followers.Follow(r.Context(), payload.UserID, user.ID); err != nil {
		switch err {
		case store.ErrConflict:
			app.conflictResponse(w, r, err)
		default:
			app.internalServerResponse(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerResponse(w, r, err)
		return
	}
}

// UnfollowUser godoc
//
//	@Summary		Unfollow a user
//	@Description	Unfollow a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int		true	"User ID"
//	@Success		204		{string}	string	"User unfollowed"
//	@Failure		400		{object}	error	"User payload missing"
//	@Failure		404		{object}	error	"User not found"
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/unfollow [put]
func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)

	// TODO: need to change to get the current user id instead of reading from request like this
	var payload FollowUser
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// the user sending request follows the user from context
	if err := app.store.Followers.Unfollow(r.Context(), payload.UserID, user.ID); err != nil {
		switch err {
		case store.ErrConflict:
			app.conflictResponse(w, r, err)
		default:
			app.internalServerResponse(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerResponse(w, r, err)
		return
	}
}

// ActivateUserHandler godoc
//
//	@Summary		Activate a User
//	@Description	Activate a User by token
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			token	path	string	true	"Token"
//	@Security		ApiKeyAuth
//	@Router			/users/activate/{token} [put]
func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	// get the token from the URL
	token := chi.URLParam(r, "token")

	// find the token inside user_invitations db
	if err := app.store.Users.Activate(r.Context(), token); err != nil {
		switch err {
		case store.ErrNotFound:
			app.badRequestResponse(w, r, err)
		default:
			app.internalServerResponse(w, r, err)
		}
		return
	}
	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerResponse(w, r, err)
	}
}
