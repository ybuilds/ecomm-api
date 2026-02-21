package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ybuilds/ecomm-api/internal/products"
)

// mount
func (api *api) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID) //for rate limiting
	r.Use(middleware.RealIP)    //for rate limiting, analytics and tracing
	r.Use(middleware.Logger)    //better terminal logging
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("system ok"))
	})

	productsService := products.NewService()
	productsHandler := products.NewHandler(productsService)
	r.Get("/products", productsHandler.ListProducts)

	return r
}

// run
func (api *api) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         api.config.addr,
		Handler:      h,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
	}

	log.Printf("server listening at addr %s", api.config.addr)

	return srv.ListenAndServe()
}

type api struct {
	config config
}

type config struct {
	addr string
	db   dbConfig
}

type dbConfig struct {
	dsn string
}
