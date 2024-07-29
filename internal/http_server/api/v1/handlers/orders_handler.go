package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	DTOs "github.com/ShevelevEvgeniy/kafkaManager/internal/dto"
	"github.com/ShevelevEvgeniy/kafkaManager/internal/service/service_interfaces"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type OrdersHandler struct {
	log       *zap.Logger
	service   service_interfaces.OrderService
	validator *validator.Validate
}

func NewOrdersHandler(log *zap.Logger, service service_interfaces.OrderService, validator *validator.Validate) *OrdersHandler {
	return &OrdersHandler{
		log:       log,
		service:   service,
		validator: validator,
	}
}

func (h *OrdersHandler) CreateOrder(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.log.Info("Received HTTP POST request: " + r.RequestURI)

		var dto DTOs.Order

		err := json.NewDecoder(r.Body).Decode(&dto)
		if err != nil {
			h.log.Error("failed to decode request body", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": err.Error()})
			return
		}

		err = h.validator.Struct(dto)
		if err != nil {
			h.log.Error("failed to validate request body", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": err.Error()})
			return
		}

		queryParams := r.URL.Query()
		requestId := queryParams.Get("request_id")
		if requestId == "" {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "request_id is required"})
			return
		}

		err = h.service.SaveOrderMessage(ctx, dto, requestId)
		if err != nil {
			h.log.Error("failed to save orders", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "failed to save orders"})
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, map[string]string{"message": "success"})
	}
}
