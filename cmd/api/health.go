package main

import "net/http"

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Health Check: OK!"))

	app.store.Posts.Create(r.Context())
	app.store.Users.Create(r.Context())
}
