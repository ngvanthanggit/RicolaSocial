package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/ngvanthanggit/RicolaSocial/internal/store"
)

type userKey string

var userCtx userKey = "user"

var (
	ErrInvalidToken = errors.New("invalid token")
)

func validateToken(ctx context.Context, store *store.Storage, token string) (*store.User, error) {
	// implement token validation logic and retrieve user information
	if token == "" {
		//return nil, ErrInvalidToken
	}

	user, err := store.Users.GetUserById(ctx, 1)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (app *application) userAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		token := r.Header.Get("Authorization")

		user, err := validateToken(ctx, &app.store, token)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		ctx = context.WithValue(ctx, userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
