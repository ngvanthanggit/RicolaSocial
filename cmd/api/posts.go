package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/ngvanthanggit/RicolaSocial/internal/store"
)

type postKey string

const postCtx postKey = "post"

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// checking for the validation of the JSON data read
	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		UserID:  1,
	}

	ctx := r.Context()

	if err := app.store.Posts.Create(ctx, post); err != nil {
		app.internalServerResponse(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusCreated, post); err != nil {
		app.internalServerResponse(w, r, err)
		return
	}
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	// retrieve the post from the request context
	post := getPostFromCtx(r)

	comments, err := app.store.Comments.GetByPostID(r.Context(), post.ID)
	if err != nil {
		app.internalServerResponse(w, r, err)
		return
	}

	post.Comments = comments

	if err := writeJSON(w, http.StatusOK, post); err != nil {
		app.internalServerResponse(w, r, err)
		return
	}
}

func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	// retrieve post from request context
	post := getPostFromCtx(r)

	// delete the post with id = postID
	err := app.store.Posts.Delete(r.Context(), post.ID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
			return
		default:
			app.internalServerResponse(w, r, err)
		}
		return
	}
	if err := writeJSONResponseMessage(w, http.StatusOK, "successful deleted the post"); err != nil {
		app.internalServerResponse(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// "omitempty" to avoid auto-defining the value if the variable is not declared (unknown)
// applying the pointer to the string type is for the case where we don't want to modify
// one of the fields (title or content), then they will be nil pointer, instead of "" string
type UpdatePostPayload struct {
	Title   *string `json:"title" validate:"omitempty,max=100"`
	Content *string `json:"content" validate:"omitempty,max=1000"`
}

func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	var payload UpdatePostPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// retrieve post from request context
	post := getPostFromCtx(r)

	// modify the post instance
	if payload.Title != nil {
		post.Title = *(payload.Title)
	}
	if payload.Content != nil {
		post.Content = *(payload.Content)
	}

	// udpate the post to the database
	if err := app.store.Posts.Update(r.Context(), post); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.badRequestResponse(w, r, err)
		default:
			app.internalServerResponse(w, r, err)
		}
		return
	}
}

func (app *application) postsContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get the post id from the url
		postIDParam := chi.URLParam(r, "postID")
		// convert the post id to integer type
		postID, err := strconv.ParseInt(postIDParam, 10, 64)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		// get the post by the postID
		ctx := r.Context()
		post, err := app.store.Posts.GetPostById(ctx, postID)

		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerResponse(w, r, err)
			}
			return
		}

		// update the post to the context of the request
		// postCtx = "post"
		ctx = context.WithValue(ctx, postCtx, post)
		// passing the new context to the middleware chain
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getPostFromCtx(r *http.Request) *store.Post {
	return r.Context().Value(postCtx).(*store.Post)
}
