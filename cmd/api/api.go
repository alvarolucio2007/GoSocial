package main

import (
	"log"
	"net/http"
	"time"

	"github.com/alvarolucio2007/GoSocial/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	config  config
	storage store.Storage
}
type config struct {
	addr string
	db   dbConfig
	env  string
}
type dbConfig struct {
	addr        string
	maxOpenConn int
	maxIdleConn int
	maxIdleTime time.Duration
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.ClientIPFromRemoteAddr)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)

		r.Route("/posts", func(r chi.Router) {
			r.Post("/", app.createPostHandler)
			r.Route("/{postID}", func(r chi.Router) {
				r.Use(app.postContextMiddleware)

				r.Get("/", app.readPostHandler)
				r.Put("/", app.updatePostHandler)
				r.Delete("/", app.deletePostHandler)
			})
		})
		r.Route("/users", func(r chi.Router) {
			r.Post("/", app.createUserHandler)
			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.userContextMiddleware)

				r.Get("/", app.readUserHandler)
				r.Delete("/", app.deleteUserHandler)
				r.Put("/", app.updateUserHandler)

				r.Put("/follow", app.followUserHandler)
				r.Put("/unfollow", app.unfollowUserHandler)
			})
		})
	})
	return r
}

func (app *application) run(mux http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
	}
	log.Printf("server has started at %s", app.config.addr)
	return srv.ListenAndServe()
}
