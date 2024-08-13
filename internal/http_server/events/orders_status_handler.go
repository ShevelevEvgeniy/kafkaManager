package events

import (
	"context"
	"encoding/json"

	DTOs "github.com/ShevelevEvgeniy/kafkaManager/internal/dto"
	servInterfaces "github.com/ShevelevEvgeniy/kafkaManager/internal/service/service_interfaces"
	"github.com/ShevelevEvgeniy/kafkaManager/pkg/event_dispatcher"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type MessageStatusHandler struct {
	log       *zap.Logger
	service   servInterfaces.OrderService
	validator *validator.Validate
}

func NewMessageStatusHandler(log *zap.Logger, service servInterfaces.OrderService, validator *validator.Validate) *MessageStatusHandler {
	return &MessageStatusHandler{
		log:       log,
		service:   service,
		validator: validator,
	}
}

func (ms *MessageStatusHandler) Handle(ctx context.Context, msg event_dispatcher.Message) {
	var dto DTOs.OrderMessageResponse

	err := json.Unmarshal(msg.Value, &dto)
	if err != nil {
		ms.log.Error("failed to unmarshal order message", zap.Error(err))
		return
	}

	if err := ms.validator.Struct(dto); err != nil {
		ms.log.Error("failed to validate order message", zap.Error(err))
		return
	}

	err = ms.service.UpdateStatusOrderMessage(ctx, dto)
	if err != nil {
		ms.log.Error("failed to update order message status", zap.Error(err))
		return
	}

	ms.log.Info("updated order message status", zap.String("message", string(msg.Value)))
}
