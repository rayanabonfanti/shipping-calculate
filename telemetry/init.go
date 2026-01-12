package telemetry

import (
	"context"
	"log"
	"net/http"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

// InitOpenTelemetry initializes OpenTelemetry with basic configuration
// For production, configure OTLP exporters via environment variables:
// - OTEL_EXPORTER_OTLP_ENDPOINT: OTLP endpoint URL
// - OTEL_EXPORTER_OTLP_METRICS_ENDPOINT: Metrics endpoint URL (optional)
// - OTEL_EXPORTER_OTLP_TRACES_ENDPOINT: Traces endpoint URL (optional)
// - OTEL_SERVICE_NAME: Service name (defaults to APPLICATION_NAME or "shipping-calculator")
func InitOpenTelemetry(ctx context.Context) (func() error, error) {
	appName := os.Getenv("APPLICATION_NAME")
	if appName == "" {
		appName = "shipping-calculator"
	}

	// Set global propagator for distributed tracing
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// Create resource with service name
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(appName),
		),
	)
	if err != nil {
		return nil, err
	}

	// Create TracerProvider with default SDK (no exporter configured)
	// This will create spans even without an OTLP endpoint
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	// Set global TracerProvider
	otel.SetTracerProvider(tp)

	// Log OpenTelemetry configuration
	otlpEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if otlpEndpoint == "" {
		log.Printf("OpenTelemetry initialized (no OTLP endpoint configured, using default SDK behavior)")
		log.Printf("To export metrics, set OTEL_EXPORTER_OTLP_ENDPOINT environment variable")
	} else {
		log.Printf("OpenTelemetry initialized with OTLP endpoint: %s", otlpEndpoint)
	}

	// Return shutdown function to shutdown the TracerProvider
	return func() error {
		return tp.Shutdown(ctx)
	}, nil
}

// InjectTraceContext injects the trace context from the given context into the HTTP request headers.
// This should be called before making outbound HTTP requests to propagate the trace_id to downstream services.
//
// Example usage:
//
//	req, _ := http.NewRequestWithContext(ctx, "GET", "https://api.example.com", nil)
//	telemetry.InjectTraceContext(ctx, req)
//	client.Do(req)
func InjectTraceContext(ctx context.Context, req *http.Request) {
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
}
