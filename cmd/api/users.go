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
