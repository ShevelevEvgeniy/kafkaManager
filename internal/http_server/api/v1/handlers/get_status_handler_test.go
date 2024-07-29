package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ShevelevEvgeniy/kafkaManager/internal/http_server/api/v1/handlers/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// mockery --name StatusService --dir internal/service/service_interfaces --output internal/http_server/api/v1/handlers/mocks --outpkg mocks

func TestGetStatusHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestID      string
		mockReturn     string
		mockError      error
		expectedStatus int
		expectedBody   string
		setupMock      bool
	}{
		{
			name:           "Success",
			requestID:      "test-request-id",
			mockReturn:     "status",
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":"200"}`,
			setupMock:      true,
		},
		{
			name:           "Bad Request",
			requestID:      "",
			mockReturn:     "",
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"status":"400","error":"request_id is required"}`,
			setupMock:      false,
		},
		{
			name:           "Internal Server Error",
			requestID:      "test-request-id",
			mockReturn:     "",
			mockError:      assert.AnError,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"status":"500","error":"internal server error"}`,
			setupMock:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(mocks.StatusService)
			if tt.setupMock {
				if tt.mockError != nil {
					mockService.On("GetStatus", mock.Anything, tt.requestID).Return(tt.mockReturn, tt.mockError)
				} else {
					mockService.On("GetStatus", mock.Anything, tt.requestID).Return(tt.mockReturn, nil)
				}
			}

			logger := zap.NewNop()

			handler := NewGetStatusHandler(logger, mockService)

			req, err := http.NewRequest(http.MethodGet, "/api/v1/get_status?request_id="+tt.requestID, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			router := chi.NewRouter()
			router.Get("/api/v1/get_status", handler.GetStatus(context.Background()))

			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.JSONEq(t, tt.expectedBody, rr.Body.String())

			if tt.setupMock {
				mockService.AssertExpectations(t)
			}
		})
	}
}
