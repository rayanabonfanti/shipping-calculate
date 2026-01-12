package telemetry

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func restoreEnvVars(appName string) {
	restoreEnvVar("APPLICATION_NAME", appName)
}

func restoreEnvVar(key, originalValue string) {
	os.Unsetenv(key)
	os.Setenv(key, originalValue)
}

func TestGetInstance(t *testing.T) {
	// Arrange
	originalAppName := os.Getenv("APPLICATION_NAME")
	defer restoreEnvVars(originalAppName)

	// Act
	instance := getInstance()

	// Assert
	assert.NotNil(t, instance)
}

func TestGetInstance_Singleton(t *testing.T) {
	// Arrange
	originalAppName := os.Getenv("APPLICATION_NAME")
	defer restoreEnvVars(originalAppName)

	// Act
	instance1 := getInstance()
	instance2 := getInstance()

	// Assert
	assert.Equal(t, instance1, instance2)
}

func TestRecordLatencyOperationA(t *testing.T) {
	// Arrange
	ctx := context.Background()
	latency := int64(100)
	resource := "test-resource"

	// Act
	RecordLatencyOperationA(ctx, latency, resource)

	// Assert
	// No error means success
}

func TestRecordMemoryHeapServer(t *testing.T) {
	// Arrange
	ctx := context.Background()
	amount := int64(1024)

	// Act
	RecordMemoryHeapServer(ctx, amount)

	// Assert
	// No error means success
}

func TestRecordMemoryNoHeapServer(t *testing.T) {
	// Arrange
	ctx := context.Background()
	amount := int64(2048)

	// Act
	RecordMemoryNoHeapServer(ctx, amount)

	// Assert
	// No error means success
}

func TestIncrementHttpRequestHandled(t *testing.T) {
	// Arrange
	ctx := context.Background()
	httpMethod := "GET"
	status := 200

	// Act
	IncrementHttpRequestHandled(ctx, httpMethod, status)

	// Assert
	// No error means success
}

func TestIncrementHttpRequestHandled_DifferentMethods(t *testing.T) {
	// Arrange
	ctx := context.Background()
	tests := []struct {
		method string
		status int
	}{
		{"GET", 200},
		{"POST", 201},
		{"PUT", 200},
		{"DELETE", 204},
		{"PATCH", 200},
	}

	for _, tt := range tests {
		// Act
		IncrementHttpRequestHandled(ctx, tt.method, tt.status)

		// Assert
		// No error means success
	}
}

func TestIncrementHttpRequestHandled_DifferentStatusCodes(t *testing.T) {
	// Arrange
	ctx := context.Background()
	tests := []struct {
		method string
		status int
	}{
		{"GET", 200},
		{"GET", 404},
		{"GET", 500},
		{"POST", 400},
		{"POST", 201},
	}

	for _, tt := range tests {
		// Act
		IncrementHttpRequestHandled(ctx, tt.method, tt.status)

		// Assert
		// No error means success
	}
}

func TestIncrementShipmentCalculate(t *testing.T) {
	// Arrange
	ctx := context.Background()

	// Act
	IncrementShipmentCalculate(ctx)

	// Assert
	// No error means success
}

func TestRecordShipmentCalculateTime(t *testing.T) {
	// Arrange
	ctx := context.Background()
	timeMs := int64(150)

	// Act
	RecordShipmentCalculateTime(ctx, timeMs)

	// Assert
	// No error means success
}

func TestRecordShipmentCalculateTime_DifferentValues(t *testing.T) {
	// Arrange
	ctx := context.Background()
	times := []int64{0, 50, 100, 200, 500, 1000, 5000}

	for _, timeMs := range times {
		// Act
		RecordShipmentCalculateTime(ctx, timeMs)

		// Assert
		// No error means success
	}
}

func TestRecordShipmentCalculateCostDistribution(t *testing.T) {
	// Arrange
	ctx := context.Background()
	cost := 1250.0

	// Act
	RecordShipmentCalculateCostDistribution(ctx, cost)

	// Assert
	// No error means success
}

func TestRecordShipmentCalculateCostDistribution_DifferentValues(t *testing.T) {
	// Arrange
	ctx := context.Background()
	costs := []float64{0.0, 100.0, 500.0, 1000.0, 2000.0, 5000.0, 10000.0}

	for _, cost := range costs {
		// Act
		RecordShipmentCalculateCostDistribution(ctx, cost)

		// Assert
		// No error means success
	}
}

func TestIncrementShipmentCalculateError(t *testing.T) {
	// Arrange
	ctx := context.Background()

	// Act
	IncrementShipmentCalculateError(ctx)

	// Assert
	// No error means success
}

func TestGetInstance_WithCustomAppName(t *testing.T) {
	// Arrange
	originalAppName := os.Getenv("APPLICATION_NAME")
	os.Setenv("APPLICATION_NAME", "custom-app")
	defer restoreEnvVars(originalAppName)

	// Act
	instance := getInstance()

	// Assert
	assert.NotNil(t, instance)
}

func TestAllMetrics_WithContext(t *testing.T) {
	// Arrange
	ctx := context.Background()

	// Act & Assert
	RecordLatencyOperationA(ctx, 100, "test")
	RecordMemoryHeapServer(ctx, 1024)
	RecordMemoryNoHeapServer(ctx, 2048)
	IncrementHttpRequestHandled(ctx, "GET", 200)
	IncrementShipmentCalculate(ctx)
	RecordShipmentCalculateTime(ctx, 150)
	RecordShipmentCalculateCostDistribution(ctx, 1250.0)
	IncrementShipmentCalculateError(ctx)

	// No error means success
}

func TestAllMetrics_WithNilContext(t *testing.T) {
	// Arrange
	var ctx context.Context

	// Act & Assert
	RecordLatencyOperationA(ctx, 100, "test")
	RecordMemoryHeapServer(ctx, 1024)
	RecordMemoryNoHeapServer(ctx, 2048)
	IncrementHttpRequestHandled(ctx, "GET", 200)
	IncrementShipmentCalculate(ctx)
	RecordShipmentCalculateTime(ctx, 150)
	RecordShipmentCalculateCostDistribution(ctx, 1250.0)
	IncrementShipmentCalculateError(ctx)

	// No error means success
}

func TestRecordLatencyOperationA_DifferentResources(t *testing.T) {
	// Arrange
	ctx := context.Background()
	resources := []string{"resource1", "resource2", "resource3", "test-resource"}

	for _, resource := range resources {
		// Act
		RecordLatencyOperationA(ctx, 100, resource)

		// Assert
		// No error means success
	}
}

func TestRecordLatencyOperationA_DifferentLatencies(t *testing.T) {
	// Arrange
	ctx := context.Background()
	latencies := []int64{0, 10, 50, 100, 500, 1000, 5000}

	for _, latency := range latencies {
		// Act
		RecordLatencyOperationA(ctx, latency, "test")

		// Assert
		// No error means success
	}
}

func TestRecordMemoryHeapServer_DifferentAmounts(t *testing.T) {
	// Arrange
	ctx := context.Background()
	amounts := []int64{0, 1024, 4096, 8192, 16384, 32768, 65536}

	for _, amount := range amounts {
		// Act
		RecordMemoryHeapServer(ctx, amount)

		// Assert
		// No error means success
	}
}

func TestRecordMemoryNoHeapServer_DifferentAmounts(t *testing.T) {
	// Arrange
	ctx := context.Background()
	amounts := []int64{0, 1024, 4096, 8192, 16384, 32768, 65536}

	for _, amount := range amounts {
		// Act
		RecordMemoryNoHeapServer(ctx, amount)

		// Assert
		// No error means success
	}
}
