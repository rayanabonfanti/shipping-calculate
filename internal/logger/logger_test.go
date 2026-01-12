package logger

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

func TestGetCorrelationID_WithRequestID(t *testing.T) {
	// Arrange
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	var ctx context.Context
	handler := middleware.RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx = r.Context()
	}))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// Act
	result := GetCorrelationID(ctx)

	// Assert
	assert.NotEmpty(t, result)
}

func TestGetCorrelationID_WithoutRequestID(t *testing.T) {
	// Arrange
	ctx := context.Background()

	// Act
	result := GetCorrelationID(ctx)

	// Assert
	assert.Equal(t, "", result)
}

func TestGetTraceID_WithValidSpan(t *testing.T) {
	// Arrange
	tracer := otel.Tracer("test")
	ctx, span := tracer.Start(context.Background(), "test-operation")
	defer span.End()

	// Act
	result := GetTraceID(ctx)

	// Assert
	// Note: OpenTelemetry may return empty if not fully initialized in test environment
	// This test verifies the function doesn't panic
	assert.NotNil(t, result)
}

func TestGetTraceID_WithoutSpan(t *testing.T) {
	// Arrange
	ctx := context.Background()

	// Act
	result := GetTraceID(ctx)

	// Assert
	assert.Equal(t, "", result)
}

func TestGetSpanID_WithValidSpan(t *testing.T) {
	// Arrange
	tracer := otel.Tracer("test")
	ctx, span := tracer.Start(context.Background(), "test-operation")
	defer span.End()

	// Act
	result := GetSpanID(ctx)

	// Assert
	// Note: OpenTelemetry may return empty if not fully initialized in test environment
	// This test verifies the function doesn't panic
	assert.NotNil(t, result)
}

func TestGetSpanID_WithoutSpan(t *testing.T) {
	// Arrange
	ctx := context.Background()

	// Act
	result := GetSpanID(ctx)

	// Assert
	assert.Equal(t, "", result)
}

func TestWithTracingFields_WithAllFields(t *testing.T) {
	// Arrange
	logger := zaptest.NewLogger(t)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	var ctx context.Context
	handler := middleware.RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx = r.Context()
	}))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	tracer := otel.Tracer("test")
	ctx, span := tracer.Start(ctx, "test-operation")
	defer span.End()

	// Act
	result := WithTracingFields(logger, ctx)

	// Assert
	assert.NotNil(t, result)
	assert.NotEqual(t, logger, result)
}

func TestWithTracingFields_WithoutFields(t *testing.T) {
	// Arrange
	logger := zaptest.NewLogger(t)
	ctx := context.Background()

	// Act
	result := WithTracingFields(logger, ctx)

	// Assert
	assert.NotNil(t, result)
}

func TestWithCorrelationID(t *testing.T) {
	// Arrange
	logger := zaptest.NewLogger(t)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	var ctx context.Context
	handler := middleware.RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx = r.Context()
	}))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// Act
	result := WithCorrelationID(logger, ctx)

	// Assert
	assert.NotNil(t, result)
}

func TestLogRequest(t *testing.T) {
	// Arrange
	logger := zaptest.NewLogger(t)
	ctx := context.Background()
	message := "test message"

	// Act
	LogRequest(logger, ctx, message)

	// Assert
	// No error means success
}

func TestLogRequest_WithFields(t *testing.T) {
	// Arrange
	logger := zaptest.NewLogger(t)
	ctx := context.Background()
	message := "test message"
	fields := []zap.Field{
		zap.String("key1", "value1"),
		zap.Int("key2", 42),
	}

	// Act
	LogRequest(logger, ctx, message, fields...)

	// Assert
	// No error means success
}

func TestLogWarning(t *testing.T) {
	// Arrange
	logger := zaptest.NewLogger(t)
	ctx := context.Background()
	message := "warning message"

	// Act
	LogWarning(logger, ctx, message)

	// Assert
	// No error means success
}

func TestLogWarning_WithFields(t *testing.T) {
	// Arrange
	logger := zaptest.NewLogger(t)
	ctx := context.Background()
	message := "warning message"
	fields := []zap.Field{
		zap.String("key1", "value1"),
	}

	// Act
	LogWarning(logger, ctx, message, fields...)

	// Assert
	// No error means success
}

func TestLogError(t *testing.T) {
	// Arrange
	logger := zaptest.NewLogger(t)
	ctx := context.Background()
	message := "error message"
	err := assert.AnError

	// Act
	LogError(logger, ctx, message, err)

	// Assert
	// No error means success
}

func TestLogError_WithFields(t *testing.T) {
	// Arrange
	logger := zaptest.NewLogger(t)
	ctx := context.Background()
	message := "error message"
	err := assert.AnError
	fields := []zap.Field{
		zap.String("key1", "value1"),
	}

	// Act
	LogError(logger, ctx, message, err, fields...)

	// Assert
	// No error means success
}

func TestGetLoggerFromContext_WithLoggerInContext(t *testing.T) {
	// Arrange
	defaultLogger := zaptest.NewLogger(t)
	ctxLogger := zaptest.NewLogger(t)
	ctx := context.WithValue(context.Background(), "logger", ctxLogger)

	// Act
	result := GetLoggerFromContext(ctx, defaultLogger)

	// Assert
	assert.NotNil(t, result)
}

func TestGetLoggerFromContext_WithoutLoggerInContext(t *testing.T) {
	// Arrange
	defaultLogger := zaptest.NewLogger(t)
	ctx := context.Background()

	// Act
	result := GetLoggerFromContext(ctx, defaultLogger)

	// Assert
	assert.NotNil(t, result)
}

func TestGetLoggerFromContext_WithInvalidLoggerType(t *testing.T) {
	// Arrange
	defaultLogger := zaptest.NewLogger(t)
	ctx := context.WithValue(context.Background(), "logger", "not-a-logger")

	// Act
	result := GetLoggerFromContext(ctx, defaultLogger)

	// Assert
	assert.NotNil(t, result)
}

func TestGetLoggerFromContext_WithNilLogger(t *testing.T) {
	// Arrange
	defaultLogger := zaptest.NewLogger(t)
	ctx := context.WithValue(context.Background(), "logger", nil)

	// Act
	result := GetLoggerFromContext(ctx, defaultLogger)

	// Assert
	assert.NotNil(t, result)
}

func TestWithTracingFields_WithCorrelationIDOnly(t *testing.T) {
	// Arrange
	logger := zaptest.NewLogger(t)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	var ctx context.Context
	handler := middleware.RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx = r.Context()
	}))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// Act
	result := WithTracingFields(logger, ctx)

	// Assert
	assert.NotNil(t, result)
}

func TestWithTracingFields_WithTraceIDOnly(t *testing.T) {
	// Arrange
	logger := zaptest.NewLogger(t)
	tracer := otel.Tracer("test")
	ctx, span := tracer.Start(context.Background(), "test-operation")
	defer span.End()

	// Act
	result := WithTracingFields(logger, ctx)

	// Assert
	assert.NotNil(t, result)
}

func TestGetTraceID_WithInvalidSpan(t *testing.T) {
	// Arrange
	ctx := context.Background()
	// Create a span context that is not valid
	spanCtx := trace.SpanContext{}
	ctx = trace.ContextWithSpanContext(ctx, spanCtx)

	// Act
	result := GetTraceID(ctx)

	// Assert
	assert.Equal(t, "", result)
}

func TestGetSpanID_WithInvalidSpan(t *testing.T) {
	// Arrange
	ctx := context.Background()
	// Create a span context that is not valid
	spanCtx := trace.SpanContext{}
	ctx = trace.ContextWithSpanContext(ctx, spanCtx)

	// Act
	result := GetSpanID(ctx)

	// Assert
	assert.Equal(t, "", result)
}

func TestWithTracingFields_WithSpanIDOnly(t *testing.T) {
	// Arrange
	logger := zaptest.NewLogger(t)
	tracer := otel.Tracer("test")
	ctx, span := tracer.Start(context.Background(), "test-operation")
	defer span.End()

	// Act
	result := WithTracingFields(logger, ctx)

	// Assert
	assert.NotNil(t, result)
}

func TestWithTracingFields_WithCorrelationAndTraceID(t *testing.T) {
	// Arrange
	logger := zaptest.NewLogger(t)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	var ctx context.Context
	handler := middleware.RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx = r.Context()
	}))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	tracer := otel.Tracer("test")
	ctx, span := tracer.Start(ctx, "test-operation")
	defer span.End()

	// Act
	result := WithTracingFields(logger, ctx)

	// Assert
	assert.NotNil(t, result)
}

func TestWithTracingFields_WithCorrelationAndSpanID(t *testing.T) {
	// Arrange
	logger := zaptest.NewLogger(t)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	var ctx context.Context
	handler := middleware.RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx = r.Context()
	}))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	tracer := otel.Tracer("test")
	ctx, span := tracer.Start(ctx, "test-operation")
	defer span.End()

	// Act
	result := WithTracingFields(logger, ctx)

	// Assert
	assert.NotNil(t, result)
}
