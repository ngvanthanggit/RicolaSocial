package main

import (
	"net/http"
	"time"

	"github.com/ngvanthanggit/RicolaSocial/internal/store"
)

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	dbTimeZone, err := app.store.DBTime.GetDBTimeZone()
	if err != nil {
		app.internalServerResponse(w, r, err)
		return
	}
	location, err := time.LoadLocation(dbTimeZone)
	if err != nil {
		app.internalServerResponse(w, r, err)
		return
	}
	time.Local = location

	feedQuery := store.PaginatedFeedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
		Since:  time.Date(2003, time.January, 1, 0, 0, 0, 0, time.UTC).Format(time.RFC3339),
		Until:  time.Now().Format(time.RFC3339),
	}

	feedQuery, err = feedQuery.Parse(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if err := Validate.Struct(feedQuery); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	// TODO: get the current logged user id after doing authorization
	curUserID := int64(200)

	feed, err := app.store.Posts.GetUserFeed(ctx, curUserID, feedQuery)
	if err != nil {
		app.internalServerResponse(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, feed); err != nil {
		app.internalServerResponse(w, r, err)
		return
	}
}
