package service_interfaces

import (
	"context"

	"github.com/ShevelevEvgeniy/kafkaManager/internal/dto"
)

type OrderService interface {
	SaveOrderMessage(ctx context.Context, order dto.Order, requestId string) error
	UpdateStatusOrderMessage(ctx context.Context, dto dto.OrderMessageResponse) error
}
