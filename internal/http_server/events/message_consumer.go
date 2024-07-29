package events

import (
	"context"

	kafka2 "github.com/ShevelevEvgeniy/kafkaManager/internal/clients/kafka"
	"github.com/ShevelevEvgeniy/kafkaManager/internal/clients/kafka/topics"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type MessageConsumerEvent struct {
	log         *zap.Logger
	kafkaClient *kafka2.Kafka
}

func NewMessageConsumerEvent(log *zap.Logger, kafka *kafka2.Kafka) *MessageConsumerEvent {
	return &MessageConsumerEvent{
		log:         log,
		kafkaClient: kafka,
	}
}

func (e *MessageConsumerEvent) Start(ctx context.Context) error {
	err := e.kafkaClient.SubscribeToTopics(ctx, topics.OrderStatusTopic)
	if err != nil {
		e.log.Error("error while subscribing to topics", zap.Error(err))
		return errors.Wrap(err, "error while subscribing to topics")
	}

	messageChan, errChan := e.kafkaClient.ListenToTopic(ctx, topics.OrderStatusTopic)

	for {
		select {
		case <-ctx.Done():
			e.kafkaClient.CloseConsumer()
			e.log.Info("shutting down message consumer")
			return ctx.Err()
		case err = <-errChan:
			if err != nil {
				e.log.Error("error while listening to topic", zap.Error(err))
				return errors.Wrap(err, "error while listening to topic")
			}
		case msg, ok := <-messageChan:
			if !ok {
				e.log.Error("message channel closed")
				return errors.New("message channel closed")
			}

			e.log.Info("received message", zap.String("message", string(msg.Value)))

			err = e.handleMessages(msg, ctx)
			if err != nil {
				e.log.Error("error while handling message", zap.Error(err))
				return errors.Wrap(err, "error while handling message")
			}
		}
	}
}

func (e *MessageConsumerEvent) handleMessages(msg *kafka.Message, ctx context.Context) error {
	//TODO

	return nil
}
