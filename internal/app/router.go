package app

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	ApiV1Group = "/api/v1"
)

func initRouter(ctx context.Context, di DiContainer) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)

	router.Route(ApiV1Group, func(router chi.Router) {
	})

	return router
}
