package main

import (
	"context"
	"errors"
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

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := ctx.Value(userCtx)

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerResponse(w, r, err)
		return
	}
}
