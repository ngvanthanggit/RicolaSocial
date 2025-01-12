package main

import (
	"net/http"
)

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// TODO: get the current logged user id after doing authorization
	curUserID := int64(200)

	feed, err := app.store.Posts.GetUserFeed(ctx, curUserID)
	if err != nil {
		app.internalServerResponse(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, feed); err != nil {
		app.internalServerResponse(w, r, err)
		return
	}
}
