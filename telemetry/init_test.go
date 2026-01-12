package telemetry

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitOpenTelemetry_DefaultAppName(t *testing.T) {
	// Arrange
	ctx := context.Background()
	originalAppName := os.Getenv("APPLICATION_NAME")
	os.Unsetenv("APPLICATION_NAME")
	defer restoreEnvVar("APPLICATION_NAME", originalAppName)

	// Act
	shutdown, err := InitOpenTelemetry(ctx)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, shutdown)
	err = shutdown()
	assert.NoError(t, err)
}

func TestInitOpenTelemetry_CustomAppName(t *testing.T) {
	// Arrange
	ctx := context.Background()
	originalAppName := os.Getenv("APPLICATION_NAME")
	os.Setenv("APPLICATION_NAME", "test-app")
	defer restoreEnvVar("APPLICATION_NAME", originalAppName)

	// Act
	shutdown, err := InitOpenTelemetry(ctx)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, shutdown)
	err = shutdown()
	assert.NoError(t, err)
}

func TestInitOpenTelemetry_WithOTLPEndpoint(t *testing.T) {
	// Arrange
	ctx := context.Background()
	originalEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://localhost:4317")
	defer restoreEnvVar("OTEL_EXPORTER_OTLP_ENDPOINT", originalEndpoint)

	// Act
	shutdown, err := InitOpenTelemetry(ctx)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, shutdown)
	err = shutdown()
	assert.NoError(t, err)
}

func TestInitOpenTelemetry_ShutdownFunction(t *testing.T) {
	// Arrange
	ctx := context.Background()

	// Act
	shutdown, err := InitOpenTelemetry(ctx)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, shutdown)

	// Test shutdown function
	err = shutdown()
	assert.NoError(t, err)

	// Test multiple shutdown calls
	err = shutdown()
	assert.NoError(t, err)
}

func TestInitOpenTelemetry_ContextCancellation(t *testing.T) {
	// Arrange
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Act
	shutdown, err := InitOpenTelemetry(ctx)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, shutdown)
	err = shutdown()
	assert.NoError(t, err)
}
