package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rbonfanti/shipping-calculator/internal/handler"
	"github.com/rbonfanti/shipping-calculator/internal/service"
	"github.com/rbonfanti/shipping-calculator/telemetry"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	// Initialize OpenTelemetry
	shutdown, err := telemetry.InitOpenTelemetry(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize OpenTelemetry: %v", err)
	}
	defer func() {
		if err := shutdown(); err != nil {
			log.Printf("Error shutting down OpenTelemetry: %v", err)
		}
	}()

	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// Initialize services
	shippingService := service.NewShippingService()

	// Initialize handlers
	shippingHandler := handler.NewShippingHandler(shippingService, logger)

	// Setup router
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(otelMiddleware)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Register routes
	r.Post("/calculate", shippingHandler.CalculateShipping)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Graceful shutdown
	go func() {
		logger.Info("Server starting", zap.String("port", port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Server shutting down")
	if err := server.Close(); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	}

	// Shutdown OpenTelemetry
	if err := shutdown(); err != nil {
		logger.Error("Error shutting down OpenTelemetry", zap.Error(err))
	}
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// otelMiddleware creates OpenTelemetry spans for HTTP requests
func otelMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Extract trace context from headers
		ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(r.Header))

		// Get tracer
		tracer := otel.Tracer("shipping-calculator")

		// Start span
		ctx, span := tracer.Start(ctx, r.Method+" "+r.URL.Path,
			trace.WithAttributes(
				semconv.HTTPMethod(r.Method),
				semconv.HTTPURL(r.URL.String()),
				semconv.HTTPRoute(r.URL.Path),
			),
		)
		defer span.End()

		// Add span to request context
		r = r.WithContext(ctx)

		// Wrap ResponseWriter to capture status code
		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK, // default status
		}

		// Call next handler
		next.ServeHTTP(wrapped, r)

		// Set span status based on response
		span.SetAttributes(semconv.HTTPStatusCodeKey.Int(wrapped.statusCode))
		if wrapped.statusCode >= 400 {
			span.RecordError(nil)
		}
	})
}
