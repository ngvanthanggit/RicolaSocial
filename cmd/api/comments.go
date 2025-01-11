package main

import (
	"net/http"

	"github.com/ngvanthanggit/RicolaSocial/internal/store"
)

type CommentPayload struct {
	Content string `json:"content" validate:"required"`
}

func (app *application) commentOnPostHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// retrieve user
	user := ctx.Value(userCtx).(*store.User)

	// retrieve post
	post := ctx.Value(postCtx).(*store.Post)

	var payload CommentPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	comment := store.Comment{
		PostID:  post.ID,
		UserID:  user.ID,
		Content: payload.Content,
		User:    *user,
	}
	if err := app.store.Comments.Create(ctx, &comment); err != nil {
		app.internalServerResponse(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, comment); err != nil {
		app.internalServerResponse(w, r, err)
		return
	}
}
