package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

type userKeyType string

const userExample userKeyType = "username"

func (app *application) BasicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			app.unauthorizedErrorResponse(w, r, fmt.Errorf("no authorization header"))
			return
		}

		// we got the authorization header -> convert to the authorizing values
		// token
		splitedStrings := strings.Split(authHeader, " ")
		if len(splitedStrings) < 1 {
			app.badRequestResponse(w, r, fmt.Errorf("invalid authorization token"))
			return
		}

		token := splitedStrings[len(splitedStrings)-1]
		// decode
		decoded, err := base64.StdEncoding.DecodeString(token)
		if err != nil {
			app.unauthorizedErrorResponse(w, r, err)
			return
		}

		// check credentials
		username := app.config.auth.basic.user
		password := app.config.auth.basic.pass

		creds := strings.SplitN(string(decoded), ":", 2)

		if len(creds) != 2 || creds[0] != username || creds[1] != password {
			app.unauthorizedErrorResponse(w, r, fmt.Errorf("invalid credentials"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
