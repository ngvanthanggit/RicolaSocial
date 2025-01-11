package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ngvanthanggit/RicolaSocial/internal/store"
)

type application struct {
	config config
	store  store.Storage
}

type config struct {
	addr string
	db   dbConfig
	env  string
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

// returning the mux for the server
func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		// authenticate user
		r.Use(app.userAuthentication)

		r.Get("/health", app.healthCheckHandler)

		/// POST/GET... v1/posts/
		r.Route("/posts", func(r chi.Router) {
			r.Post("/", app.createPostHandler)

			r.Route("/{postID}", func(r chi.Router) {
				// using the post middleware to get the post, pass it to the request
				r.Use(app.postsContextMiddleware)

				r.Get("/", app.getPostHandler)
				r.Delete("/", app.deletePostHandler)
				r.Patch("/", app.updatePostHandler)

				r.Route("/comment", func(r chi.Router) {
					r.Post("/", app.commentOnPostHandler)
				})
			})
		})
	})

	return r
}

// run the server and return the error status
func (app *application) run(mux http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
	}

	log.Printf("Server is running at %s", app.config.addr)

	return srv.ListenAndServe()
}
