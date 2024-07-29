package app

import (
	"context"

	_ "github.com/ShevelevEvgeniy/kafkaManager/docs"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	ApiV1Group = "/api/v1"
	Order      = "/orders"
	GetStatus  = "/get_status"
	Swagger    = "/swagger/*"
)

func initRouter(ctx context.Context, di DiContainer) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)

	router.Route(ApiV1Group, func(router chi.Router) {
		router.Post(Order, di.OrdersHandler(ctx).CreateOrder(ctx))
		router.Get(GetStatus, di.GetStatusHandler(ctx).GetStatus(ctx))
	})

	router.Get(Swagger, httpSwagger.WrapHandler)

	return router
}
