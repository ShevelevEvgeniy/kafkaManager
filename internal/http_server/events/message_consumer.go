package events

import (
	"context"
	"encoding/json"

	"github.com/ShevelevEvgeniy/kafkaManager/config"
	kafka2 "github.com/ShevelevEvgeniy/kafkaManager/internal/clients/kafka"
	"github.com/ShevelevEvgeniy/kafkaManager/internal/clients/kafka/topics"
	DTOs "github.com/ShevelevEvgeniy/kafkaManager/internal/dto"
	servInterfaces "github.com/ShevelevEvgeniy/kafkaManager/internal/service/service_interfaces"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type MessageConsumerEvent struct {
	log         *zap.Logger
	kafkaClient kafka2.ClientKafka
	service     servInterfaces.OrderService
	validator   *validator.Validate
}

func NewMessageConsumerEvent(log *zap.Logger, kafka kafka2.ClientKafka, service servInterfaces.OrderService, validator *validator.Validate) *MessageConsumerEvent {
	return &MessageConsumerEvent{
		log:         log,
		kafkaClient: kafka,
		service:     service,
		validator:   validator,
	}
}

func (e *MessageConsumerEvent) Start(ctx context.Context, cfg config.Kafka) {
	go e.listenAndHandleMessages(ctx, cfg)
}

func (e *MessageConsumerEvent) listenAndHandleMessages(ctx context.Context, cfg config.Kafka) {
	err := e.kafkaClient.SubscribeToTopics(ctx, cfg.Topic)
	if err != nil {
		e.log.Error("error while subscribing to topics", zap.Error(err))
		return
	}

	messageChan, errChan := e.kafkaClient.ListenToTopic(ctx, topics.OrderStatusTopic)

	for {
		select {
		case <-ctx.Done():
			e.kafkaClient.Close()
			e.log.Info("shutting down message consumer")
			return
		case err = <-errChan:
			if err != nil {
				e.log.Error("error while listening to topic", zap.Error(err))
				return
			}
		case msg, ok := <-messageChan:
			if !ok {
				e.log.Error("message channel closed")
				return
			}

			e.log.Info("received message", zap.String("message", string(msg.Value)))

			err = e.handleMessages(msg, ctx)
			if err != nil {
				e.log.Error("error while handling message", zap.Error(err))
				return
			}
		}
	}
}

func (e *MessageConsumerEvent) handleMessages(msg *kafka.Message, ctx context.Context) error {
	var dto DTOs.OrderMessageResponse

	err := json.Unmarshal(msg.Value, &dto)
	if err != nil {
		e.log.Error("failed to unmarshal order message", zap.Error(err))
		return errors.Wrap(err, "failed to unmarshal order message")
	}

	if err := e.validator.Struct(dto); err != nil {
		e.log.Error("failed to validate order message", zap.Error(err))
		return errors.Wrap(err, "failed to validate order message")
	}

	err = e.service.UpdateStatusOrderMessage(ctx, dto)
	if err != nil {
		e.log.Error("failed to update order message status", zap.Error(err))
		return errors.Wrap(err, "failed to update order message status")
	}

	e.log.Info("updated order message status", zap.String("message", string(msg.Value)))

	return nil
}
