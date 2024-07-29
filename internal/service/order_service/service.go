package order_service

import (
	"context"
	"encoding/json"

	"github.com/ShevelevEvgeniy/kafkaManager/internal/clients/kafka"
	"github.com/ShevelevEvgeniy/kafkaManager/internal/clients/kafka/statuses"
	"github.com/ShevelevEvgeniy/kafkaManager/internal/clients/kafka/topics"
	"github.com/ShevelevEvgeniy/kafkaManager/internal/convertor"
	"github.com/ShevelevEvgeniy/kafkaManager/internal/dto"
	repoInterfaces "github.com/ShevelevEvgeniy/kafkaManager/internal/postgres/repository/repository_interfaces"
	"github.com/pkg/errors"
)

type OrderService struct {
	repository repoInterfaces.MessageTrackerRepository
	kafka      *kafka.Kafka
}

func NewOrderService(repository repoInterfaces.MessageTrackerRepository, kafka *kafka.Kafka) *OrderService {
	return &OrderService{
		repository: repository,
		kafka:      kafka,
	}
}

func (s *OrderService) SaveOrderMessage(ctx context.Context, order dto.Order, requestId string) error {
	message, err := json.Marshal(order)
	if err != nil {
		return errors.Wrap(err, "failed to marshal order")
	}
	key := []byte(requestId)

	status := statuses.Sent

	err = s.kafka.SendMessage(ctx, key, message, topics.OrderTopic)
	if err != nil {
		status = statuses.Failed
	}

	model := convertor.OrderDtoToTrackingModel(message, requestId, status)

	err = s.repository.SaveMessage(ctx, model)
	if err != nil {
		return errors.Wrap(err, "failed to save message tracker")
	}

	return nil
}
