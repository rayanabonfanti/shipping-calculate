package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/rbonfanti/shipping-calculator/internal/logger"
	"github.com/rbonfanti/shipping-calculator/internal/model"
	"github.com/rbonfanti/shipping-calculator/internal/service"
	"github.com/rbonfanti/shipping-calculator/telemetry"
	"go.uber.org/zap"
)

// ShippingHandler handles HTTP requests for shipping calculations
type ShippingHandler struct {
	service service.ShippingServiceInterface
	logger  *zap.Logger
}

// NewShippingHandler creates a new shipping handler instance
func NewShippingHandler(shippingService service.ShippingServiceInterface, logger *zap.Logger) *ShippingHandler {
	return &ShippingHandler{
		service: shippingService,
		logger:  logger,
	}
}

// CalculateShipping handles POST /calculate requests
func (h *ShippingHandler) CalculateShipping(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	startTime := time.Now()

	// Record request metric
	telemetry.IncrementShipmentCalculate(ctx)

	// Decode request body
	var req model.CalculateShippingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		telemetry.IncrementShipmentCalculateError(ctx)
		logger.LogError(h.logger, ctx, "Erro no serviço de cálculo: falha ao decodificar requisição", err)
		h.writeJSON(ctx, w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	// Calculate volume for logging
	volume := req.Dimensions.Length * req.Dimensions.Width * req.Dimensions.Height

	// Log request with structured fields
	logger.LogRequest(h.logger, ctx, "Solicitação de cálculo de custos",
		zap.String("origem", req.OriginZipcode),
		zap.String("destino", req.DestinationZipcode),
		zap.Float64("peso", req.Weight),
		zap.Float64("volume", volume),
	)

	// Calculate shipping
	response, err := h.service.CalculateShipping(ctx, &req)
	if err != nil {
		telemetry.IncrementShipmentCalculateError(ctx)
		logger.LogError(h.logger, ctx, "Erro no serviço de cálculo", err)
		h.writeJSON(ctx, w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Record success metrics
	elapsed := time.Since(startTime)
	telemetry.RecordShipmentCalculateTime(ctx, elapsed.Milliseconds())
	telemetry.RecordShipmentCalculateCostDistribution(ctx, response.ShippingCost)

	// Return response
	h.writeJSON(ctx, w, http.StatusOK, response)
}

// writeJSON is a helper function to write JSON responses
func (h *ShippingHandler) writeJSON(ctx context.Context, w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.LogError(h.logger, ctx, "Erro ao codificar resposta JSON", err)
	}
}
