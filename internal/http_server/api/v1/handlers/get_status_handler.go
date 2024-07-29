package handlers

import (
	"context"
	"net/http"

	"github.com/go-chi/render"
	"go.uber.org/zap"
)

type StatusService interface {
	GetStatus(ctx context.Context, requestId string) (string, error)
}

type GetStatusHandler struct {
	service StatusService
	log     *zap.Logger
}

func NewGetStatusHandler(log *zap.Logger, service StatusService) *GetStatusHandler {
	return &GetStatusHandler{
		service: service,
		log:     log,
	}
}

func (h *GetStatusHandler) GetStatus(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.log.Info("Received HTTP GET request: " + r.RequestURI)

		queryParams := r.URL.Query()
		requestId := queryParams.Get("request_id")

		h.log.Info("request id", zap.Any("request_id", requestId))

		status, err := h.service.GetStatus(ctx, requestId)
		if err != nil {
			h.log.Error("failed to get status", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
		}

		h.log.Info("got status", zap.Any("status", status))

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, status)
	}
}
