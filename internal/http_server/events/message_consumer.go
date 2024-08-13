package events

import (
	"context"

	kafka2 "github.com/ShevelevEvgeniy/kafkaManager/internal/clients/kafka"
	"github.com/ShevelevEvgeniy/kafkaManager/internal/convertor"
	"github.com/ShevelevEvgeniy/kafkaManager/pkg/event_dispatcher"
	"go.uber.org/zap"
)

type MessageConsumerEvent struct {
	log             *zap.Logger
	kafkaClient     kafka2.ClientKafka
	eventDispatcher event_dispatcher.EventDispatcherInterface
}

func NewMessageConsumerEvent(log *zap.Logger, kafka kafka2.ClientKafka, eventDispatcher event_dispatcher.EventDispatcherInterface) *MessageConsumerEvent {
	return &MessageConsumerEvent{
		log:             log,
		kafkaClient:     kafka,
		eventDispatcher: eventDispatcher,
	}
}

func (e *MessageConsumerEvent) Start(ctx context.Context) {
	go e.listenAndHandleMessages(ctx)
}

func (e *MessageConsumerEvent) listenAndHandleMessages(ctx context.Context) {
	messageChan, errChan := e.kafkaClient.ListenToTopic(ctx)

	for {
		select {
		case <-ctx.Done():
			e.kafkaClient.Close()
			e.log.Info("shutting down message consumer")
			return
		case err := <-errChan:
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

			e.log.Info("topic for event", zap.String("topic", *msg.TopicPartition.Topic))

			e.eventDispatcher.Trigger(*msg.TopicPartition.Topic, convertor.KafkaMessageToMessage(msg), ctx)
		}
	}
}
