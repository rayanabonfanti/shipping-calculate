package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rbonfanti/shipping-calculator/internal/model"
	"github.com/rbonfanti/shipping-calculator/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zaptest"
)

// MockShippingService is a mock implementation of ShippingServiceInterface
type MockShippingService struct {
	mock.Mock
}

func (m *MockShippingService) CalculateShipping(ctx context.Context, req *model.CalculateShippingRequest) (*model.CalculateShippingResponse, error) {
	args := m.Called(ctx, req)
	resp := args.Get(0)
	err := args.Error(1)
	if resp == nil {
		return nil, err
	}
	return resp.(*model.CalculateShippingResponse), err
}

// addRequestID adds RequestID to the request context for testing
// This simulates the chi middleware.RequestID behavior by running the request through the middleware
func addRequestID(req *http.Request) *http.Request {
	var updatedReq *http.Request
	// Use chi middleware to add RequestID
	handler := middleware.RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Capture the request with updated context
		updatedReq = r
	}))

	// Run request through middleware to get RequestID in context
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// Return the updated request if middleware processed it, otherwise return original
	if updatedReq != nil {
		return updatedReq
	}
	return req
}

func TestNewShippingHandler(t *testing.T) {
	// Arrange
	shippingService := service.NewShippingService()
	logger := zaptest.NewLogger(t)

	// Act
	handler := NewShippingHandler(shippingService, logger)

	// Assert
	assert.NotNil(t, handler)
	assert.NotNil(t, handler.service)
	assert.NotNil(t, handler.logger)
}

func TestCalculateShipping_ValidRequest_Standard(t *testing.T) {
	// Arrange
	mockService := new(MockShippingService)
	logger := zaptest.NewLogger(t)
	handler := NewShippingHandler(mockService, logger)

	reqBody := model.CalculateShippingRequest{
		OriginZipcode:      "12345678",
		DestinationZipcode: "87654321",
		Weight:             1.0,
		Dimensions: model.PackageDimensions{
			Length: 10.0,
			Width:  10.0,
			Height: 10.0,
		},
		IsExpress: false,
	}
	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/calculate", bytes.NewReader(bodyBytes))
	req = addRequestID(req)
	w := httptest.NewRecorder()

	expectedResponse := &model.CalculateShippingResponse{
		ShippingCost:          1250.0,
		EstimatedDeliveryTime: "2 dias",
		AvailableServices:     []string{"standard", "express"},
		ShippingOptions: []model.ShippingOption{
			{Service: "standard", Cost: 1250.0, Time: "2 dias"},
			{Service: "express", Cost: 1875.0, Time: "1 dia"},
		},
	}

	mockService.On("CalculateShipping", mock.Anything, mock.MatchedBy(func(req *model.CalculateShippingRequest) bool {
		return req.OriginZipcode == "12345678" &&
			req.DestinationZipcode == "87654321" &&
			req.Weight == 1.0 &&
			req.IsExpress == false
	})).Return(expectedResponse, nil).Once()

	// Act
	handler.CalculateShipping(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)

	var response model.CalculateShippingResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.ShippingCost, response.ShippingCost)
	assert.Equal(t, expectedResponse.EstimatedDeliveryTime, response.EstimatedDeliveryTime)
}

func TestCalculateShipping_InvalidJSON(t *testing.T) {
	// Arrange
	mockService := new(MockShippingService)
	logger := zaptest.NewLogger(t)
	handler := NewShippingHandler(mockService, logger)

	req := httptest.NewRequest(http.MethodPost, "/calculate", bytes.NewReader([]byte("invalid json")))
	req = addRequestID(req)
	w := httptest.NewRecorder()

	// Act
	handler.CalculateShipping(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertNumberOfCalls(t, "CalculateShipping", 0)

	var errorResponse map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, "invalid request body", errorResponse["error"])
}

func TestCalculateShipping_EmptyBody(t *testing.T) {
	// Arrange
	mockService := new(MockShippingService)
	logger := zaptest.NewLogger(t)
	handler := NewShippingHandler(mockService, logger)

	req := httptest.NewRequest(http.MethodPost, "/calculate", bytes.NewReader([]byte("")))
	req = addRequestID(req)
	w := httptest.NewRecorder()

	// Act
	handler.CalculateShipping(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertNumberOfCalls(t, "CalculateShipping", 0)

	var errorResponse map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, "invalid request body", errorResponse["error"])
}

func TestCalculateShipping_ServiceError(t *testing.T) {
	// Arrange
	mockService := new(MockShippingService)
	logger := zaptest.NewLogger(t)
	handler := NewShippingHandler(mockService, logger)

	reqBody := model.CalculateShippingRequest{
		OriginZipcode:      "12345678",
		DestinationZipcode: "87654321",
		Weight:             1.0,
		Dimensions: model.PackageDimensions{
			Length: 10.0,
			Width:  10.0,
			Height: 10.0,
		},
		IsExpress: false,
	}
	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/calculate", bytes.NewReader(bodyBytes))
	req = addRequestID(req)
	w := httptest.NewRecorder()

	expectedError := errors.New("invalid origin_zipcode: origin_zipcode is required")
	mockService.On("CalculateShipping", mock.Anything, mock.Anything).Return(nil, expectedError).Once()

	// Act
	handler.CalculateShipping(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertExpectations(t)

	var errorResponse map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, expectedError.Error(), errorResponse["error"])
}

func TestCalculateShipping_ValidationError(t *testing.T) {
	// Arrange
	mockService := new(MockShippingService)
	logger := zaptest.NewLogger(t)
	handler := NewShippingHandler(mockService, logger)

	reqBody := model.CalculateShippingRequest{
		OriginZipcode:      "",
		DestinationZipcode: "87654321",
		Weight:             1.0,
		Dimensions: model.PackageDimensions{
			Length: 10.0,
			Width:  10.0,
			Height: 10.0,
		},
		IsExpress: false,
	}
	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/calculate", bytes.NewReader(bodyBytes))
	req = addRequestID(req)
	w := httptest.NewRecorder()

	expectedError := errors.New("invalid origin_zipcode: origin_zipcode is required")
	mockService.On("CalculateShipping", mock.Anything, mock.Anything).Return(nil, expectedError).Once()

	// Act
	handler.CalculateShipping(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertExpectations(t)

	var errorResponse map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, expectedError.Error(), errorResponse["error"])
}

func TestCalculateShipping_ExpressShipping(t *testing.T) {
	// Arrange
	mockService := new(MockShippingService)
	logger := zaptest.NewLogger(t)
	handler := NewShippingHandler(mockService, logger)

	reqBody := model.CalculateShippingRequest{
		OriginZipcode:      "12345678",
		DestinationZipcode: "87654321",
		Weight:             1.0,
		Dimensions: model.PackageDimensions{
			Length: 10.0,
			Width:  10.0,
			Height: 10.0,
		},
		IsExpress: true,
	}
	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/calculate", bytes.NewReader(bodyBytes))
	req = addRequestID(req)
	w := httptest.NewRecorder()

	expectedResponse := &model.CalculateShippingResponse{
		ShippingCost:          1875.0,
		EstimatedDeliveryTime: "1 dia",
		AvailableServices:     []string{"standard", "express"},
		ShippingOptions: []model.ShippingOption{
			{Service: "standard", Cost: 1250.0, Time: "2 dias"},
			{Service: "express", Cost: 1875.0, Time: "1 dia"},
		},
	}

	mockService.On("CalculateShipping", mock.Anything, mock.MatchedBy(func(req *model.CalculateShippingRequest) bool {
		return req.IsExpress == true
	})).Return(expectedResponse, nil).Once()

	// Act
	handler.CalculateShipping(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)

	var response model.CalculateShippingResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.ShippingCost, response.ShippingCost)
	assert.Equal(t, expectedResponse.EstimatedDeliveryTime, response.EstimatedDeliveryTime)
}

func TestCalculateShipping_NilBody(t *testing.T) {
	// Arrange
	mockService := new(MockShippingService)
	logger := zaptest.NewLogger(t)
	handler := NewShippingHandler(mockService, logger)

	req := httptest.NewRequest(http.MethodPost, "/calculate", nil)
	req = addRequestID(req)
	w := httptest.NewRecorder()

	// Act
	handler.CalculateShipping(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertNumberOfCalls(t, "CalculateShipping", 0)

	var errorResponse map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, "invalid request body", errorResponse["error"])
}

func TestWriteJSON_ErrorEncoding(t *testing.T) {
	// Arrange
	mockService := new(MockShippingService)
	logger := zaptest.NewLogger(t)
	handler := NewShippingHandler(mockService, logger)
	ctx := context.Background()
	w := httptest.NewRecorder()
	invalidData := make(chan int)

	// Act
	handler.writeJSON(ctx, w, http.StatusOK, invalidData)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}
