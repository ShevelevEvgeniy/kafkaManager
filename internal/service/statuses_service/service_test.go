package statuses_service

import (
	"context"
	"errors"
	"testing"

	"github.com/ShevelevEvgeniy/kafkaManager/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockery --name MessageTrackerRepository --dir internal/postgres/repository/repository_interfaces --output internal/service/mocks --outpkg mocks

func TestGetStatus(t *testing.T) {
	tests := []struct {
		name           string
		requestId      string
		mockReturn     string
		mockError      error
		expectedStatus string
		expectedError  error
	}{
		{
			name:           "Success",
			requestId:      "request-123",
			mockReturn:     "status-success",
			mockError:      nil,
			expectedStatus: "status-success",
			expectedError:  nil,
		},
		{
			name:           "Repository Error",
			requestId:      "request-123",
			mockReturn:     "",
			mockError:      errors.New("repository error"),
			expectedStatus: "",
			expectedError:  errors.New("failed to get message by request id: repository error"),
		},
		{
			name:           "No Status Found",
			requestId:      "request-123",
			mockReturn:     "",
			mockError:      nil,
			expectedStatus: "",
			expectedError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MessageTrackerRepository)
			mockRepo.On("GetMessageByRequestId", mock.Anything, tt.requestId).Return(tt.mockReturn, tt.mockError)

			statusService := NewStatusService(mockRepo)
			status, err := statusService.GetStatus(context.Background(), tt.requestId)

			assert.Equal(t, tt.expectedStatus, status)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
