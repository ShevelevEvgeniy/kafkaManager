package events

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/ShevelevEvgeniy/kafkaManager/config"
	"github.com/ShevelevEvgeniy/kafkaManager/internal/clients/kafka/topics"
	"github.com/ShevelevEvgeniy/kafkaManager/internal/http_server/events/mocks"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// mockery --name ClientKafka --dir internal/clients/kafka --output internal/http_server/events/mocks --outpkg mocks
// mockery --name OrderService --dir internal/service/service_interfaces --output internal/http_server/events/mocks --outpkg mocks

type testCase struct {
	name           string
	subscribeError error
	listenError    error
	message        *kafka.Message
	handleError    error
	expectedError  string
}

func TestListenAndHandleMessages(t *testing.T) {
	tests := []testCase{
		{
			name:           "Success",
			subscribeError: nil,
			listenError:    nil,
			message:        &kafka.Message{Value: json.RawMessage(`{"request_id":"123", "status":"success"}`)},
			handleError:    nil,
			expectedError:  "",
		},
		{
			name:           "Subscription Error",
			subscribeError: errors.New("subscription error"),
			listenError:    nil,
			message:        nil,
			handleError:    nil,
			expectedError:  "error while subscribing to topics: subscription error",
		},
		{
			name:           "Listen Error",
			subscribeError: nil,
			listenError:    errors.New("listen error"),
			message:        nil,
			handleError:    nil,
			expectedError:  "error while listening to topic: listen error",
		},
		{
			name:           "Message Handling Error",
			subscribeError: nil,
			listenError:    nil,
			message:        &kafka.Message{Value: json.RawMessage(`{"request_id":"123", "status":"success"}`)},
			handleError:    errors.New("handle error"),
			expectedError:  "error while handling message: handle error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockKafkaClient := new(mocks.ClientKafka)
			mockOrderService := new(mocks.OrderService)
			mockValidator := validator.New()

			mockKafkaClient.On("SubscribeToTopics", mock.Anything, mock.Anything).Return(tt.subscribeError)

			if tt.subscribeError == nil {

				mockKafkaClient.On("ListenToTopic", mock.Anything, topics.OrderStatusTopic).Return(
					listenToTopic(tt.message, tt.listenError),
				)
			}

			if tt.handleError != nil {
				mockOrderService.On("UpdateStatusOrderMessage", mock.Anything, mock.Anything).Return(tt.handleError)
			} else {
				mockOrderService.On("UpdateStatusOrderMessage", mock.Anything, mock.Anything).Return(nil)
			}

			logger := zap.NewNop()

			consumerEvent := NewMessageConsumerEvent(logger, mockKafkaClient, mockOrderService, mockValidator)

			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			consumerEvent.listenAndHandleMessages(ctx, config.Kafka{Topic: topics.OrderStatusTopic})
		})
	}
}

func listenToTopic(msg *kafka.Message, listenError error) (<-chan *kafka.Message, <-chan error) {
	messageChan := make(chan *kafka.Message)
	errChan := make(chan error, 1)

	go func() {
		defer close(messageChan)
		defer close(errChan)

		if listenError != nil {
			errChan <- listenError
		} else {
			if msg != nil {
				messageChan <- msg
			}
		}
	}()

	return messageChan, errChan
}
