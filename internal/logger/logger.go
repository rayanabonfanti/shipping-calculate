package logger

import (
	"context"

	"github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// GetCorrelationID extracts correlation_id from context (from chi middleware.RequestID)
func GetCorrelationID(ctx context.Context) string {
	if reqID := middleware.GetReqID(ctx); reqID != "" {
		return reqID
	}
	return ""
}

// GetTraceID extracts trace_id from OpenTelemetry span context
func GetTraceID(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		return span.SpanContext().TraceID().String()
	}
	return ""
}

// GetSpanID extracts span_id from OpenTelemetry span context
func GetSpanID(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		return span.SpanContext().SpanID().String()
	}
	return ""
}

// WithTracingFields adds correlation_id, trace_id, and span_id to logger fields
func WithTracingFields(logger *zap.Logger, ctx context.Context) *zap.Logger {
	fields := []zap.Field{}

	// Add correlation_id
	if correlationID := GetCorrelationID(ctx); correlationID != "" {
		fields = append(fields, zap.String("correlation_id", correlationID))
	}

	// Add trace_id
	if traceID := GetTraceID(ctx); traceID != "" {
		fields = append(fields, zap.String("trace_id", traceID))
	}

	// Add span_id
	if spanID := GetSpanID(ctx); spanID != "" {
		fields = append(fields, zap.String("span_id", spanID))
	}

	if len(fields) > 0 {
		return logger.With(fields...)
	}
	return logger
}

// WithCorrelationID adds correlation_id to logger fields (deprecated, use WithTracingFields)
func WithCorrelationID(logger *zap.Logger, ctx context.Context) *zap.Logger {
	return WithTracingFields(logger, ctx)
}

// LogRequest logs a request with structured fields including trace_id and span_id
func LogRequest(logger *zap.Logger, ctx context.Context, message string, fields ...zap.Field) {
	logger = WithTracingFields(logger, ctx)
	logger.Info(message, fields...)
}

// LogWarning logs a warning with structured fields including trace_id and span_id
func LogWarning(logger *zap.Logger, ctx context.Context, message string, fields ...zap.Field) {
	logger = WithTracingFields(logger, ctx)
	logger.Warn(message, fields...)
}

// LogError logs an error with structured fields including trace_id and span_id
func LogError(logger *zap.Logger, ctx context.Context, message string, err error, fields ...zap.Field) {
	logger = WithTracingFields(logger, ctx)
	allFields := append(fields, zap.Error(err))
	logger.Error(message, allFields...)
}

// GetLoggerFromContext extracts logger from context or returns the default logger
// The returned logger includes trace_id and span_id from the context
func GetLoggerFromContext(ctx context.Context, defaultLogger *zap.Logger) *zap.Logger {
	if ctxLogger := ctx.Value("logger"); ctxLogger != nil {
		if l, ok := ctxLogger.(*zap.Logger); ok {
			return WithTracingFields(l, ctx)
		}
	}
	return WithTracingFields(defaultLogger, ctx)
}
