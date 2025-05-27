package main

import (
	"net/http"

	"github.com/ngvanthanggit/RicolaSocial/internal/store"
)

// getUserFeedHandler godoc
//
//	@Summary		Fetches the user feed
//	@Description	Fetches the user feed
//	@Tags			feed
//	@Accept			json
//	@Produce		json
//	@Param			since	query		string	false	"Since"
//	@Param			until	query		string	false	"Until"
//	@Param			limit	query		int		false	"Limit"
//	@Param			offset	query		int		false	"Offset"
//	@Param			sort	query		string	false	"Sort"
//	@Param			tags	query		string	false	"Tags"
//	@Param			search	query		string	false	"Search"
//	@Success		200		{object}	[]store.PostWithMetaData
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/feed [get]
func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	feedQuery := store.PaginatedFeedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
		Since:  store.DefaultStartTime,
		Until:  store.TimeNowString(),
	}

	feedQuery, err := feedQuery.Parse(r)
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
