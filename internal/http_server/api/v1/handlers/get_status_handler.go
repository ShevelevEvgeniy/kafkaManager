package handlers

import (
	"context"
	"net/http"

	"github.com/ShevelevEvgeniy/kafkaManager/internal/http_server/api/v1/response"
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

// GetStatus @Summary Get Status
// @Description Retrieve the status for a given request_id
// @Accept json
// @Produce json
// @Param request_id query string true "Request ID"
// @Success 200 {object} response.SuccessResponse "Status response"
// @Failure 400 {object} response.Response "Invalid input"
// @Failure 500 {object} response.Response "Server error"
// @Router /api/v1/get_status [get]
func (h *GetStatusHandler) GetStatus(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.log.Info("Received HTTP GET request: " + r.RequestURI)

		queryParams := r.URL.Query()
		requestId := queryParams.Get("request_id")
		if requestId == "" {
			h.log.Error("request_id is empty")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.BadRequest("request_id is required"))
			return
		}

		h.log.Info("request id", zap.Any("request_id", requestId))

		status, err := h.service.GetStatus(ctx, requestId)
		if err != nil {
			h.log.Error("failed to get status", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.InternalServerError())
			return
		}

		h.log.Info("got status", zap.Any("status", status))

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, response.OK())
	}
}
