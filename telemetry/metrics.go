package telemetry

import (
	"context"
	"log"
	"os"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

var (
	once     sync.Once
	instance *instruments
)

// Package telemetry provides helper functions to record service metrics
// using the OpenTelemetry Go SDK. The following code is an example of how a healthy metric helper should be implemented
// Users may feel free to erase or change it to fit their needs
//
// Quick reference
//
//  1. Instruments should be created only once. Metric names should be attached to instruments and appear in this file exclusively. In the example:
//     - latencyOperationA   (Int64Histogram): internal operation latency in ms
//     - memoryServer      (Int64Gauge): memory usage of the server
//     - httpRequestHandled (Int64Counter)  : total HTTP requests processed
//
//  2. Usage example:
//     start := time.Now()
//     // ... business logic ...
//     telemetry.RecordLatencyOperationA(ctx, time.Since(start).Milliseconds(), "generate_invoice")
//
//     telemetry.IncrementHttpRequestHandled(ctx, r.Method, http.StatusOK)
//
//  3. Conventions:
//     • Instrument names are snake_case and describe *what* is measured.
//     • Attributes/labels use snake.case keys and whenever possible follow the OTel semantic conventions (as shown in the examples)
//     • Always propagate the incoming context when recording.
//
// To extend:
//   - Add a new instrument to the instruments struct.
//   - Instantiate it inside getInstance() with its corresponding metric name and description.
//   - Expose a helper function that records or adds values following the patterns in the examples.
//   - Metric attributes should be passed as primitive arguments and then converted to OTel attributes inside the helper function (as shown in the examples)
//
// More info:
//
// OTel: https://opentelemetry.io/docs/specs/semconv/general/metrics
type instruments struct {
	latencyOperationA                 metric.Int64Histogram
	memoryServer                      metric.Int64Gauge
	httpRequestHandled                metric.Int64Counter
	shipmentCalculate                 metric.Int64Counter
	shipmentCalculateTime             metric.Int64Histogram
	shipmentCalculateCostDistribution metric.Float64Histogram
	shipmentCalculateError            metric.Int64Counter
}

func getInstance() *instruments {
	once.Do(func() {
		appName := os.Getenv("APPLICATION_NAME")
		if appName == "" {
			appName = "shipping-calculator"
		}

		metricPrefix := "shipping.calculate"

		meter := otel.Meter(appName)

		latencyOperationA, err := meter.Int64Histogram("latency_operation",
			metric.WithDescription("The latency of the processed operation"))
		if err != nil {
			log.Fatalf("Failed to create instrument histogram: %v", err)
		}

		memoryServer, err := meter.Int64Gauge("memory_server",
			metric.WithDescription("The current memory server used"))
		if err != nil {
			log.Fatalf("Failed to create instrument gauge: %v", err)
		}

		httpRequestHandled, err := meter.Int64Counter("http_requests_total",
			metric.WithDescription("The total number of HTTP requests"))
		if err != nil {
			log.Fatalf("Failed to create instrument counter: %v", err)
		}

		shipmentCalculate, err := meter.Int64Counter(metricPrefix,
			metric.WithDescription("Contador de cálculos solicitados"))
		if err != nil {
			log.Fatalf("Failed to create instrument counter: %v", err)
		}

		shipmentCalculateTime, err := meter.Int64Histogram(metricPrefix+".time",
			metric.WithDescription("Tempo de resposta"))
		if err != nil {
			log.Fatalf("Failed to create instrument histogram: %v", err)
		}

		shipmentCalculateCostDistribution, err := meter.Float64Histogram(metricPrefix+".cost.distribution",
			metric.WithDescription("Distribuição dos custos calculados"))
		if err != nil {
			log.Fatalf("Failed to create instrument histogram: %v", err)
		}

		shipmentCalculateError, err := meter.Int64Counter(metricPrefix+".error",
			metric.WithDescription("Contador de erros"))
		if err != nil {
			log.Fatalf("Failed to create instrument counter: %v", err)
		}

		instance = &instruments{
			latencyOperationA:                 latencyOperationA,
			memoryServer:                      memoryServer,
			httpRequestHandled:                httpRequestHandled,
			shipmentCalculate:                 shipmentCalculate,
			shipmentCalculateTime:             shipmentCalculateTime,
			shipmentCalculateCostDistribution: shipmentCalculateCostDistribution,
			shipmentCalculateError:            shipmentCalculateError,
		}
	})

	return instance
}

func RecordLatencyOperationA(ctx context.Context, latency int64, resource string) {
	getInstance().latencyOperationA.Record(ctx, latency,
		metric.WithAttributes(attribute.String("resource", resource)))
}

func RecordMemoryHeapServer(ctx context.Context, amount int64) {
	getInstance().memoryServer.Record(ctx, amount, metric.WithAttributes(
		semconv.TypeHeap,
		semconv.TelemetrySDKLanguageGo))
}

func RecordMemoryNoHeapServer(ctx context.Context, amount int64) {
	getInstance().memoryServer.Record(ctx, amount, metric.WithAttributes(
		semconv.TypeNonHeap,
		semconv.TelemetrySDKLanguageGo))
}

func IncrementHttpRequestHandled(ctx context.Context, httpMethod string, status int) {
	getInstance().httpRequestHandled.Add(ctx, 1, metric.WithAttributes(
		attribute.String("http.resource", "request"),
		semconv.HTTPMethod(httpMethod),
		semconv.HTTPStatusCodeKey.Int(status)))
}

// IncrementShipmentCalculate increments the shipment calculation counter
func IncrementShipmentCalculate(ctx context.Context) {
	getInstance().shipmentCalculate.Add(ctx, 1)
}

// RecordShipmentCalculateTime records the time taken to calculate shipment
func RecordShipmentCalculateTime(ctx context.Context, timeMs int64) {
	getInstance().shipmentCalculateTime.Record(ctx, timeMs)
}

// RecordShipmentCalculateCostDistribution records the shipping cost distribution
func RecordShipmentCalculateCostDistribution(ctx context.Context, cost float64) {
	getInstance().shipmentCalculateCostDistribution.Record(ctx, cost)
}

// IncrementShipmentCalculateError increments the shipment calculation error counter
func IncrementShipmentCalculateError(ctx context.Context) {
	getInstance().shipmentCalculateError.Add(ctx, 1)
}
