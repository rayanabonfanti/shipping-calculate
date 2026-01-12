package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/rbonfanti/shipping-calculator/internal/logger"
	"github.com/rbonfanti/shipping-calculator/internal/model"
	"github.com/rbonfanti/shipping-calculator/internal/validator"
	"go.uber.org/zap"
)

const (
	// Base shipping cost in cents (10.00 BRL = 1000 cents)
	baseCostCents = 1000.0

	// Weight surcharge: 10% of base cost per 0.5 kg
	weightSurchargeRate = 0.10
	weightUnit          = 0.5 // kg

	// Volume surcharge: 5% of base cost per 1.000 cm³
	volumeSurchargeRate = 0.05
	volumeUnit          = 1000.0 // cm³

	// Express shipping surcharge: 50% of subtotal
	expressSurchargeRate = 0.50

	// Estimated delivery days
	standardDeliveryDays = 2
	expressDeliveryDays  = 1
)

// ShippingServiceInterface defines the contract for shipping calculation service
type ShippingServiceInterface interface {
	CalculateShipping(ctx context.Context, req *model.CalculateShippingRequest) (*model.CalculateShippingResponse, error)
}

// ShippingService handles shipping calculation business logic
type ShippingService struct{}

// NewShippingService creates a new shipping service instance
func NewShippingService() *ShippingService {
	return &ShippingService{}
}

// CalculateShipping calculates shipping cost and delivery time based on package details
func (s *ShippingService) CalculateShipping(ctx context.Context, req *model.CalculateShippingRequest) (*model.CalculateShippingResponse, error) {
	// Get logger from context with correlation_id
	zapLogger := logger.GetLoggerFromContext(ctx, zap.L())

	// Validate request
	if err := validator.ValidateZipcode(req.OriginZipcode, "origin_zipcode"); err != nil {
		logger.LogWarning(zapLogger, ctx, "Solicitação com parâmetros inválidos",
			zap.String("param", "origin_zipcode"),
			zap.String("valor", req.OriginZipcode),
			zap.Error(err),
		)
		return nil, fmt.Errorf("invalid origin_zipcode: %w", err)
	}

	if err := validator.ValidateZipcode(req.DestinationZipcode, "destination_zipcode"); err != nil {
		logger.LogWarning(zapLogger, ctx, "Solicitação com parâmetros inválidos",
			zap.String("param", "destination_zipcode"),
			zap.String("valor", req.DestinationZipcode),
			zap.Error(err),
		)
		return nil, fmt.Errorf("invalid destination_zipcode: %w", err)
	}

	if err := validator.ValidateWeight(req.Weight); err != nil {
		logger.LogWarning(zapLogger, ctx, "Solicitação com parâmetros inválidos",
			zap.String("param", "weight"),
			zap.Float64("valor", req.Weight),
			zap.Error(err),
		)
		return nil, fmt.Errorf("invalid weight: %w", err)
	}

	volume := validator.CalculateVolume(req.Dimensions.Length, req.Dimensions.Width, req.Dimensions.Height)
	if err := validator.ValidateDimensions(req.Dimensions.Length, req.Dimensions.Width, req.Dimensions.Height); err != nil {
		logger.LogWarning(zapLogger, ctx, "Solicitação com parâmetros inválidos",
			zap.String("param", "dimensions"),
			zap.Float64("volume", volume),
			zap.Error(err),
		)
		return nil, fmt.Errorf("invalid dimensions: %w", err)
	}

	// Calculate base cost based on distance between zipcodes
	baseCost := s.calculateBaseCost(req.OriginZipcode, req.DestinationZipcode)

	// Calculate shipping cost
	details := s.calculateShippingDetails(baseCost, req.Weight, volume, req.IsExpress)

	// Log calculation details with structured fields
	logger.LogRequest(zapLogger, ctx, "Detalhes do cálculo",
		zap.Float64("custo_base", details.BaseCost),
		zap.Float64("acréscimo_peso", details.WeightSurcharge),
		zap.Float64("acréscimo_volume", details.VolumeSurcharge),
	)

	// Build response
	response := s.buildResponse(details, req.IsExpress)

	// Log result with structured fields
	logger.LogRequest(zapLogger, ctx, "Resultado do cálculo",
		zap.Float64("custo_envio", response.ShippingCost),
		zap.String("tempo_estimado", response.EstimatedDeliveryTime),
	)

	return response, nil
}

// calculateBaseCost calculates the base shipping cost based on distance between zipcodes
func (s *ShippingService) calculateBaseCost(originZipcode, destinationZipcode string) float64 {
	// Normalize zipcodes (remove hyphens and spaces)
	originNormalized := strings.ReplaceAll(strings.ReplaceAll(originZipcode, "-", ""), " ", "")
	destNormalized := strings.ReplaceAll(strings.ReplaceAll(destinationZipcode, "-", ""), " ", "")

	// Convert to numbers (use first 4-8 digits)
	originNum, err1 := strconv.ParseFloat(originNormalized, 64)
	destNum, err2 := strconv.ParseFloat(destNormalized, 64)

	// If conversion fails, use default base cost
	if err1 != nil || err2 != nil {
		return baseCostCents
	}

	// Calculate distance as absolute difference
	distance := originNum - destNum
	if distance < 0 {
		distance = -distance
	}

	// Base cost increases with distance
	// For same region (distance < 1000): base cost
	// For different regions: base cost * (1 + distance/10000)
	// This provides a simple distance-based pricing model
	if distance < 1000 {
		return baseCostCents
	}

	// Scale factor: 1% increase per 1000 units of distance difference
	distanceFactor := 1.0 + (distance / 10000.0)
	return baseCostCents * distanceFactor
}

// calculateShippingDetails performs the actual shipping cost calculation
func (s *ShippingService) calculateShippingDetails(baseCost, weight, volume float64, isExpress bool) *model.ShippingCalculationDetails {

	// Weight surcharge: 10% of base cost per 0.5 kg
	weightMultiplier := weight / weightUnit
	weightSurcharge := baseCost * weightSurchargeRate * weightMultiplier

	// Volume surcharge: 5% of base cost per 1000 cm³
	volumeMultiplier := volume / volumeUnit
	volumeSurcharge := baseCost * volumeSurchargeRate * volumeMultiplier

	// Subtotal before express surcharge
	subtotal := baseCost + weightSurcharge + volumeSurcharge

	// Express surcharge: 50% of subtotal if express
	var expressSurcharge float64
	if isExpress {
		expressSurcharge = subtotal * expressSurchargeRate
	}

	// Total cost
	totalCost := subtotal + expressSurcharge

	// Estimated delivery days
	estimatedDays := standardDeliveryDays
	if isExpress {
		estimatedDays = expressDeliveryDays
	}

	return &model.ShippingCalculationDetails{
		BaseCost:         baseCost,
		WeightSurcharge:  weightSurcharge,
		VolumeSurcharge:  volumeSurcharge,
		ExpressSurcharge: expressSurcharge,
		TotalCost:        totalCost,
		EstimatedDays:    estimatedDays,
	}
}

// buildResponse constructs the response with all shipping options
func (s *ShippingService) buildResponse(details *model.ShippingCalculationDetails, isExpress bool) *model.CalculateShippingResponse {
	// Calculate standard shipping cost (without express surcharge)
	standardCost := details.BaseCost + details.WeightSurcharge + details.VolumeSurcharge

	// Calculate express shipping cost (with express surcharge)
	expressCost := standardCost * (1 + expressSurchargeRate)

	// Determine which cost to return based on request
	var shippingCost float64
	var estimatedTime string
	if isExpress {
		shippingCost = expressCost
		estimatedTime = fmt.Sprintf("%d dia", expressDeliveryDays)
		if expressDeliveryDays > 1 {
			estimatedTime = fmt.Sprintf("%d dias", expressDeliveryDays)
		}
	} else {
		shippingCost = standardCost
		estimatedTime = fmt.Sprintf("%d dias", standardDeliveryDays)
	}

	// Build shipping options
	shippingOptions := []model.ShippingOption{
		{
			Service: "standard",
			Cost:    standardCost,
			Time:    fmt.Sprintf("%d dias", standardDeliveryDays),
		},
		{
			Service: "express",
			Cost:    expressCost,
			Time:    fmt.Sprintf("%d dia", expressDeliveryDays),
		},
	}

	return &model.CalculateShippingResponse{
		ShippingCost:          shippingCost,
		EstimatedDeliveryTime: estimatedTime,
		AvailableServices:     []string{"standard", "express"},
		ShippingOptions:       shippingOptions,
	}
}
