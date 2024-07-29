package order_service

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/ShevelevEvgeniy/kafkaManager/internal/clients/kafka/statuses"
	"github.com/ShevelevEvgeniy/kafkaManager/internal/clients/kafka/topics"
	"github.com/ShevelevEvgeniy/kafkaManager/internal/convertor"
	"github.com/ShevelevEvgeniy/kafkaManager/internal/dto"
	"github.com/ShevelevEvgeniy/kafkaManager/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockery --name MessageTrackerRepository --dir internal/postgres/repository/repository_interfaces --output internal/service/mocks --outpkg mocks
// mockery --name ClientKafka --dir internal/clients/kafka --output internal/service/mocks --outpkg mocks

func TestSaveOrderMessage(t *testing.T) {
	tests := []struct {
		name           string
		order          dto.Order
		requestId      string
		kafkaError     error
		repoError      error
		expectedStatus string
		expectedError  error
	}{
		{
			name:           "Success",
			order:          dto.Order{OrderID: 123},
			requestId:      "request-123",
			kafkaError:     nil,
			repoError:      nil,
			expectedStatus: statuses.Sent,
			expectedError:  nil,
		},
		{
			name:           "Kafka Error",
			order:          dto.Order{OrderID: 123},
			requestId:      "request-123",
			kafkaError:     errors.New("kafka error"),
			repoError:      nil,
			expectedStatus: statuses.Failed,
			expectedError:  nil,
		},
		{
			name:           "Repo Error",
			order:          dto.Order{OrderID: 123},
			requestId:      "request-123",
			kafkaError:     nil,
			repoError:      errors.New("repo error"),
			expectedStatus: statuses.Sent,
			expectedError:  errors.New("failed to save message tracker: repo error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockKafkaClient := new(mocks.ClientKafka)
			mockRepo := new(mocks.MessageTrackerRepository)

			message, _ := json.Marshal(tt.order)
			key := []byte(tt.requestId)

			mockKafkaClient.On("SendMessage", mock.Anything, key, message, topics.OrderTopic).Return(tt.kafkaError)

			status := tt.expectedStatus
			model := convertor.OrderDtoToTrackingModel(message, tt.requestId, status)

			mockRepo.On("SaveMessage", mock.Anything, model).Return(tt.repoError)

			orderService := NewOrderService(mockRepo, mockKafkaClient)

			err := orderService.SaveOrderMessage(context.Background(), tt.order, tt.requestId)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			mockKafkaClient.AssertExpectations(t)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateStatusOrderMessage(t *testing.T) {
	tests := []struct {
		name          string
		dto           dto.OrderMessageResponse
		repoError     error
		expectedError error
	}{
		{
			name:          "Success",
			dto:           dto.OrderMessageResponse{RequestId: "request-123", Status: statuses.Sent},
			repoError:     nil,
			expectedError: nil,
		},
		{
			name:          "Repo Error",
			dto:           dto.OrderMessageResponse{RequestId: "request-123", Status: statuses.Sent},
			repoError:     errors.New("repo error"),
			expectedError: errors.New("failed to update message status: repo error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MessageTrackerRepository)

			model := convertor.OrderMessageDtoToTrackingModel(tt.dto)

			mockRepo.On("UpdateMessageStatusByRequestId", mock.Anything, model).Return(tt.repoError)

			orderService := NewOrderService(mockRepo, nil)

			err := orderService.UpdateStatusOrderMessage(context.Background(), tt.dto)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
