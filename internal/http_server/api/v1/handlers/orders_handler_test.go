package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ShevelevEvgeniy/kafkaManager/internal/dto"
	"github.com/ShevelevEvgeniy/kafkaManager/internal/http_server/api/v1/handlers/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// mockery --name OrderService --dir internal/service/service_interfaces --output internal/http_server/api/v1/handlers/mocks --outpkg mocks

func TestCreateOrder(t *testing.T) {
	tests := []struct {
		name           string
		order          dto.Order
		requestID      string
		mockError      error
		expectedStatus int
		expectedBody   string
		setupMock      bool
	}{
		{
			name:           "Success",
			order:          getTestOrder(),
			requestID:      "test-request-id",
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":"200"}`,
			setupMock:      true,
		},
		{
			name:           "Bad Request - Invalid Body",
			order:          dto.Order{},
			requestID:      "test-request-id",
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"status":"400","error":"invalid request body"}`,
			setupMock:      false,
		},
		{
			name:           "Bad Request - Missing request_id",
			order:          getTestOrder(),
			requestID:      "",
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"status":"400","error":"request_id is required"}`,
			setupMock:      false,
		},
		{
			name:           "Internal Server Error",
			order:          getTestOrder(),
			requestID:      "test-request-id",
			mockError:      assert.AnError,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"status":"500","error":"failed to save orders"}`,
			setupMock:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(mocks.OrderService)
			if tt.setupMock {
				if tt.mockError != nil {
					mockService.On("SaveOrderMessage", mock.Anything, tt.order, tt.requestID).Return(tt.mockError)
				} else {
					mockService.On("SaveOrderMessage", mock.Anything, tt.order, tt.requestID).Return(nil)
				}
			}

			logger := zap.NewNop()
			v := validator.New()

			handler := NewOrdersHandler(logger, mockService, v)

			orderJSON, err := json.Marshal(tt.order)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest(http.MethodPost, "/api/v1/orders?request_id="+tt.requestID, bytes.NewBuffer(orderJSON))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			router := chi.NewRouter()
			router.Post("/api/v1/orders", handler.CreateOrder(context.Background()))

			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.JSONEq(t, tt.expectedBody, rr.Body.String())

			if tt.setupMock {
				mockService.AssertExpectations(t)
			}
		})
	}
}

func getTestOrder() dto.Order {
	return dto.Order{
		OrderID:     123456,
		OrderDate:   time.Date(2024, 7, 27, 14, 30, 0, 0, time.UTC),
		Status:      "pending",
		TotalAmount: 199.99,
		Products: []dto.Product{
			{
				ProductID: 98765,
				Name:      "Product A",
				Quantity:  2,
				Price:     50.00,
			},
			{
				ProductID: 87654,
				Name:      "Product B",
				Quantity:  1,
				Price:     99.99,
			},
		},
		Customer: dto.Customer{
			CustomerID: 54321,
			FirstName:  "John",
			LastName:   "Doe",
			Email:      "john.doe@example.com",
			Phone:      "+1234567890",
			Address:    "123 Elm Street, Springfield, IL",
		},
		PaymentMethod:  "credit_card",
		PaymentStatus:  "paid",
		DeliveryDate:   time.Date(2024, 8, 1, 14, 30, 0, 0, time.UTC),
		TrackingNumber: 112233,
		CreatedAt:      time.Date(2024, 7, 27, 14, 30, 0, 0, time.UTC),
		UpdatedAt:      time.Date(2024, 7, 27, 14, 30, 0, 0, time.UTC),
	}
}
