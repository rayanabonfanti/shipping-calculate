package service

import (
	"context"
	"testing"

	"github.com/rbonfanti/shipping-calculator/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestNewShippingService(t *testing.T) {
	// Arrange
	// (no setup needed)

	// Act
	service := NewShippingService()

	// Assert
	assert.NotNil(t, service)
}

func TestCalculateShipping_ValidRequest_Standard(t *testing.T) {
	// Arrange
	ctx := context.Background()
	service := NewShippingService()
	req := &model.CalculateShippingRequest{
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

	// Act
	response, err := service.CalculateShipping(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Greater(t, response.ShippingCost, 0.0)
	assert.Equal(t, "2 dias", response.EstimatedDeliveryTime)
	assert.Equal(t, []string{"standard", "express"}, response.AvailableServices)
	assert.Len(t, response.ShippingOptions, 2)
}

func TestCalculateShipping_ValidRequest_Express(t *testing.T) {
	// Arrange
	ctx := context.Background()
	service := NewShippingService()
	req := &model.CalculateShippingRequest{
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

	// Act
	response, err := service.CalculateShipping(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Greater(t, response.ShippingCost, 0.0)
	assert.Equal(t, "1 dia", response.EstimatedDeliveryTime)
	assert.Equal(t, []string{"standard", "express"}, response.AvailableServices)
	assert.Len(t, response.ShippingOptions, 2)
}

func TestCalculateShipping_InvalidOriginZipcode(t *testing.T) {
	// Arrange
	ctx := context.Background()
	service := NewShippingService()
	req := &model.CalculateShippingRequest{
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

	// Act
	response, err := service.CalculateShipping(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "invalid origin_zipcode")
}

func TestCalculateShipping_InvalidDestinationZipcode(t *testing.T) {
	// Arrange
	ctx := context.Background()
	service := NewShippingService()
	req := &model.CalculateShippingRequest{
		OriginZipcode:      "12345678",
		DestinationZipcode: "",
		Weight:             1.0,
		Dimensions: model.PackageDimensions{
			Length: 10.0,
			Width:  10.0,
			Height: 10.0,
		},
		IsExpress: false,
	}

	// Act
	response, err := service.CalculateShipping(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "invalid destination_zipcode")
}

func TestCalculateShipping_InvalidWeight(t *testing.T) {
	// Arrange
	ctx := context.Background()
	service := NewShippingService()
	req := &model.CalculateShippingRequest{
		OriginZipcode:      "12345678",
		DestinationZipcode: "87654321",
		Weight:             0.0,
		Dimensions: model.PackageDimensions{
			Length: 10.0,
			Width:  10.0,
			Height: 10.0,
		},
		IsExpress: false,
	}

	// Act
	response, err := service.CalculateShipping(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "invalid weight")
}

func TestCalculateShipping_InvalidDimensions(t *testing.T) {
	// Arrange
	ctx := context.Background()
	service := NewShippingService()
	req := &model.CalculateShippingRequest{
		OriginZipcode:      "12345678",
		DestinationZipcode: "87654321",
		Weight:             1.0,
		Dimensions: model.PackageDimensions{
			Length: 0.0,
			Width:  10.0,
			Height: 10.0,
		},
		IsExpress: false,
	}

	// Act
	response, err := service.CalculateShipping(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "invalid dimensions")
}

func TestCalculateShipping_VolumeExceedsLimit(t *testing.T) {
	// Arrange
	ctx := context.Background()
	service := NewShippingService()
	req := &model.CalculateShippingRequest{
		OriginZipcode:      "12345678",
		DestinationZipcode: "87654321",
		Weight:             1.0,
		Dimensions: model.PackageDimensions{
			Length: 30.0,
			Width:  30.0,
			Height: 20.0,
		},
		IsExpress: false,
	}

	// Act
	response, err := service.CalculateShipping(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "invalid dimensions")
}

func TestCalculateShippingDetails_StandardShipping(t *testing.T) {
	// Arrange
	service := NewShippingService()
	baseCost := 1000.0
	weight := 1.0
	volume := 1000.0
	isExpress := false

	// Act
	details := service.calculateShippingDetails(baseCost, weight, volume, isExpress)

	// Assert
	assert.NotNil(t, details)
	assert.Equal(t, 1000.0, details.BaseCost)
	assert.Greater(t, details.WeightSurcharge, 0.0)
	assert.Greater(t, details.VolumeSurcharge, 0.0)
	assert.Equal(t, 0.0, details.ExpressSurcharge)
	assert.Greater(t, details.TotalCost, details.BaseCost)
	assert.Equal(t, 2, details.EstimatedDays)
}

func TestCalculateShippingDetails_ExpressShipping(t *testing.T) {
	// Arrange
	service := NewShippingService()
	baseCost := 1000.0
	weight := 1.0
	volume := 1000.0
	isExpress := true

	// Act
	details := service.calculateShippingDetails(baseCost, weight, volume, isExpress)

	// Assert
	assert.NotNil(t, details)
	assert.Equal(t, 1000.0, details.BaseCost)
	assert.Greater(t, details.WeightSurcharge, 0.0)
	assert.Greater(t, details.VolumeSurcharge, 0.0)
	assert.Greater(t, details.ExpressSurcharge, 0.0)
	assert.Greater(t, details.TotalCost, details.BaseCost)
	assert.Equal(t, 1, details.EstimatedDays)
}

func TestCalculateShippingDetails_HeavyPackage(t *testing.T) {
	// Arrange
	service := NewShippingService()
	baseCost := 1000.0
	weight := 5.0
	volume := 1000.0
	isExpress := false

	// Act
	details := service.calculateShippingDetails(baseCost, weight, volume, isExpress)

	// Assert
	assert.NotNil(t, details)
	assert.Equal(t, 1000.0, details.BaseCost)
	assert.Greater(t, details.WeightSurcharge, 0.0)
	assert.Greater(t, details.VolumeSurcharge, 0.0)
	assert.Equal(t, 0.0, details.ExpressSurcharge)
	assert.Greater(t, details.TotalCost, details.BaseCost)
	assert.Equal(t, 2, details.EstimatedDays)
}

func TestCalculateShippingDetails_LargeVolume(t *testing.T) {
	// Arrange
	service := NewShippingService()
	baseCost := 1000.0
	weight := 1.0
	volume := 5000.0
	isExpress := false

	// Act
	details := service.calculateShippingDetails(baseCost, weight, volume, isExpress)

	// Assert
	assert.NotNil(t, details)
	assert.Equal(t, 1000.0, details.BaseCost)
	assert.Greater(t, details.WeightSurcharge, 0.0)
	assert.Greater(t, details.VolumeSurcharge, 0.0)
	assert.Equal(t, 0.0, details.ExpressSurcharge)
	assert.Greater(t, details.TotalCost, details.BaseCost)
	assert.Equal(t, 2, details.EstimatedDays)
}

func TestBuildResponse_StandardShipping(t *testing.T) {
	// Arrange
	service := NewShippingService()
	details := &model.ShippingCalculationDetails{
		BaseCost:         1000.0,
		WeightSurcharge:  200.0,
		VolumeSurcharge:  50.0,
		ExpressSurcharge: 0.0,
		TotalCost:        1250.0,
		EstimatedDays:    2,
	}
	isExpress := false

	// Act
	response := service.buildResponse(details, isExpress)

	// Assert
	assert.NotNil(t, response)
	assert.Equal(t, 1250.0, response.ShippingCost)
	assert.Equal(t, "2 dias", response.EstimatedDeliveryTime)
	assert.Equal(t, []string{"standard", "express"}, response.AvailableServices)
	assert.Len(t, response.ShippingOptions, 2)
	assert.Equal(t, "standard", response.ShippingOptions[0].Service)
	assert.Equal(t, 1250.0, response.ShippingOptions[0].Cost)
	assert.Equal(t, "2 dias", response.ShippingOptions[0].Time)
	assert.Equal(t, "express", response.ShippingOptions[1].Service)
	assert.Greater(t, response.ShippingOptions[1].Cost, response.ShippingOptions[0].Cost)
	assert.Equal(t, "1 dia", response.ShippingOptions[1].Time)
}

func TestBuildResponse_ExpressShipping(t *testing.T) {
	// Arrange
	service := NewShippingService()
	details := &model.ShippingCalculationDetails{
		BaseCost:         1000.0,
		WeightSurcharge:  200.0,
		VolumeSurcharge:  50.0,
		ExpressSurcharge: 625.0,
		TotalCost:        1875.0,
		EstimatedDays:    1,
	}
	isExpress := true

	// Act
	response := service.buildResponse(details, isExpress)

	// Assert
	assert.NotNil(t, response)
	expectedExpressCost := 1250.0 * (1 + 0.50) // 50% surcharge
	assert.Equal(t, expectedExpressCost, response.ShippingCost)
	assert.Equal(t, "1 dia", response.EstimatedDeliveryTime)
	assert.Equal(t, []string{"standard", "express"}, response.AvailableServices)
	assert.Len(t, response.ShippingOptions, 2)
	assert.Equal(t, "standard", response.ShippingOptions[0].Service)
	assert.Equal(t, 1250.0, response.ShippingOptions[0].Cost)
	assert.Equal(t, "2 dias", response.ShippingOptions[0].Time)
	assert.Equal(t, "express", response.ShippingOptions[1].Service)
	assert.Equal(t, expectedExpressCost, response.ShippingOptions[1].Cost)
	assert.Equal(t, "1 dia", response.ShippingOptions[1].Time)
}

func TestCalculateShippingDetails_ZeroWeight(t *testing.T) {
	// Arrange
	service := NewShippingService()
	baseCost := 1000.0
	weight := 0.1
	volume := 1000.0
	isExpress := false

	// Act
	details := service.calculateShippingDetails(baseCost, weight, volume, isExpress)

	// Assert
	assert.NotNil(t, details)
	assert.Equal(t, 1000.0, details.BaseCost)
	assert.GreaterOrEqual(t, details.WeightSurcharge, 0.0)
	assert.Greater(t, details.VolumeSurcharge, 0.0)
	assert.Equal(t, 0.0, details.ExpressSurcharge)
	assert.Greater(t, details.TotalCost, details.BaseCost)
	assert.Equal(t, 2, details.EstimatedDays)
}

func TestCalculateShippingDetails_ZeroVolume(t *testing.T) {
	// Arrange
	service := NewShippingService()
	baseCost := 1000.0
	weight := 1.0
	volume := 100.0
	isExpress := false

	// Act
	details := service.calculateShippingDetails(baseCost, weight, volume, isExpress)

	// Assert
	assert.NotNil(t, details)
	assert.Equal(t, 1000.0, details.BaseCost)
	assert.Greater(t, details.WeightSurcharge, 0.0)
	assert.GreaterOrEqual(t, details.VolumeSurcharge, 0.0)
	assert.Equal(t, 0.0, details.ExpressSurcharge)
	assert.Greater(t, details.TotalCost, details.BaseCost)
	assert.Equal(t, 2, details.EstimatedDays)
}

func TestCalculateShippingDetails_ExpressWithHeavyPackage(t *testing.T) {
	// Arrange
	service := NewShippingService()
	baseCost := 1000.0
	weight := 10.0
	volume := 5000.0
	isExpress := true

	// Act
	details := service.calculateShippingDetails(baseCost, weight, volume, isExpress)

	// Assert
	assert.NotNil(t, details)
	assert.Equal(t, 1000.0, details.BaseCost)
	assert.Greater(t, details.WeightSurcharge, 0.0)
	assert.Greater(t, details.VolumeSurcharge, 0.0)
	assert.Greater(t, details.ExpressSurcharge, 0.0)
	assert.Greater(t, details.TotalCost, details.BaseCost)
	assert.Equal(t, 1, details.EstimatedDays)
}

func TestCalculateShipping_WithDifferentZipcodeFormats(t *testing.T) {
	// Arrange
	ctx := context.Background()
	service := NewShippingService()
	req := &model.CalculateShippingRequest{
		OriginZipcode:      "12345-678",
		DestinationZipcode: "87654-321",
		Weight:             1.0,
		Dimensions: model.PackageDimensions{
			Length: 10.0,
			Width:  10.0,
			Height: 10.0,
		},
		IsExpress: false,
	}

	// Act
	response, err := service.CalculateShipping(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Greater(t, response.ShippingCost, 0.0)
}

func TestCalculateShipping_WithSpacesInZipcode(t *testing.T) {
	// Arrange
	ctx := context.Background()
	service := NewShippingService()
	req := &model.CalculateShippingRequest{
		OriginZipcode:      "12345 678",
		DestinationZipcode: "87654 321",
		Weight:             1.0,
		Dimensions: model.PackageDimensions{
			Length: 10.0,
			Width:  10.0,
			Height: 10.0,
		},
		IsExpress: false,
	}

	// Act
	response, err := service.CalculateShipping(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Greater(t, response.ShippingCost, 0.0)
}

func TestCalculateShipping_ExampleCase(t *testing.T) {
	// Arrange
	ctx := context.Background()
	service := NewShippingService()
	req := &model.CalculateShippingRequest{
		OriginZipcode:      "1414",
		DestinationZipcode: "1428",
		Weight:             0.5,
		Dimensions: model.PackageDimensions{
			Length: 15.0,
			Width:  8.0,
			Height: 2.0,
		},
		IsExpress: false,
	}

	// Act
	response, err := service.CalculateShipping(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, response)
	// Base cost: 1000 (distance < 1000)
	// Weight surcharge: 0.5kg / 0.5 = 1.0 → 1000 * 0.10 * 1.0 = 100
	// Volume surcharge: 240cm³ / 1000 = 0.24 → 1000 * 0.05 * 0.24 = 12
	// Standard cost: 1000 + 100 + 12 = 1112
	expectedStandardCost := 1112.0
	assert.Equal(t, expectedStandardCost, response.ShippingCost)
	assert.Equal(t, "2 dias", response.EstimatedDeliveryTime)
	assert.Equal(t, []string{"standard", "express"}, response.AvailableServices)
	assert.Len(t, response.ShippingOptions, 2)
	assert.Equal(t, "standard", response.ShippingOptions[0].Service)
	assert.Equal(t, expectedStandardCost, response.ShippingOptions[0].Cost)
	assert.Equal(t, "2 dias", response.ShippingOptions[0].Time)
	assert.Equal(t, "express", response.ShippingOptions[1].Service)
	expectedExpressCost := expectedStandardCost * 1.5 // 50% surcharge
	assert.Equal(t, expectedExpressCost, response.ShippingOptions[1].Cost)
	assert.Equal(t, "1 dia", response.ShippingOptions[1].Time)
}

func TestCalculateBaseCost_SameRegion(t *testing.T) {
	// Arrange
	service := NewShippingService()
	originZipcode := "1414"
	destinationZipcode := "1428"

	// Act
	baseCost := service.calculateBaseCost(originZipcode, destinationZipcode)

	// Assert
	// Distance is 14 (< 1000), so should return base cost
	assert.Equal(t, 1000.0, baseCost)
}

func TestCalculateBaseCost_DifferentRegions(t *testing.T) {
	// Arrange
	service := NewShippingService()
	originZipcode := "01000-000"
	destinationZipcode := "20000-000"

	// Act
	baseCost := service.calculateBaseCost(originZipcode, destinationZipcode)

	// Assert
	// Distance is 10000, so should have increased base cost
	assert.Greater(t, baseCost, 1000.0)
}

func TestCalculateShippingDetails_WeightSurcharge_10PercentPerHalfKg(t *testing.T) {
	// Arrange
	service := NewShippingService()
	baseCost := 1000.0
	weight := 1.0 // 1 kg = 2 units of 0.5 kg
	volume := 1000.0
	isExpress := false

	// Act
	details := service.calculateShippingDetails(baseCost, weight, volume, isExpress)

	// Assert
	// Weight multiplier: 1.0 / 0.5 = 2.0
	// Weight surcharge: 1000 * 0.10 * 2.0 = 200
	expectedWeightSurcharge := 200.0
	assert.Equal(t, expectedWeightSurcharge, details.WeightSurcharge)
}

func TestCalculateShippingDetails_WeightSurcharge_MultipleHalfKgs(t *testing.T) {
	// Arrange
	service := NewShippingService()
	baseCost := 1000.0
	weight := 2.5 // 2.5 kg = 5 units of 0.5 kg
	volume := 1000.0
	isExpress := false

	// Act
	details := service.calculateShippingDetails(baseCost, weight, volume, isExpress)

	// Assert
	// Weight multiplier: 2.5 / 0.5 = 5.0
	// Weight surcharge: 1000 * 0.10 * 5.0 = 500
	expectedWeightSurcharge := 500.0
	assert.Equal(t, expectedWeightSurcharge, details.WeightSurcharge)
}

func TestCalculateShippingDetails_VolumeSurcharge_5PercentPer1000Cm3(t *testing.T) {
	// Arrange
	service := NewShippingService()
	baseCost := 1000.0
	weight := 1.0
	volume := 2000.0 // 2000 cm³ = 2 units of 1000 cm³
	isExpress := false

	// Act
	details := service.calculateShippingDetails(baseCost, weight, volume, isExpress)

	// Assert
	// Volume multiplier: 2000 / 1000 = 2.0
	// Volume surcharge: 1000 * 0.05 * 2.0 = 100
	expectedVolumeSurcharge := 100.0
	assert.Equal(t, expectedVolumeSurcharge, details.VolumeSurcharge)
}

func TestCalculateShippingDetails_VolumeSurcharge_Multiple1000Cm3(t *testing.T) {
	// Arrange
	service := NewShippingService()
	baseCost := 1000.0
	weight := 1.0
	volume := 5000.0 // 5000 cm³ = 5 units of 1000 cm³
	isExpress := false

	// Act
	details := service.calculateShippingDetails(baseCost, weight, volume, isExpress)

	// Assert
	// Volume multiplier: 5000 / 1000 = 5.0
	// Volume surcharge: 1000 * 0.05 * 5.0 = 250
	expectedVolumeSurcharge := 250.0
	assert.Equal(t, expectedVolumeSurcharge, details.VolumeSurcharge)
}

func TestCalculateShippingDetails_ExpressSurcharge_50Percent(t *testing.T) {
	// Arrange
	service := NewShippingService()
	baseCost := 1000.0
	weight := 1.0
	volume := 1000.0
	isExpress := true

	// Act
	details := service.calculateShippingDetails(baseCost, weight, volume, isExpress)

	// Assert
	// Weight surcharge: 1000 * 0.10 * 2.0 = 200
	// Volume surcharge: 1000 * 0.05 * 1.0 = 50
	// Subtotal: 1000 + 200 + 50 = 1250
	// Express surcharge: 1250 * 0.50 = 625
	expectedSubtotal := 1250.0
	expectedExpressSurcharge := expectedSubtotal * 0.50
	assert.Equal(t, expectedExpressSurcharge, details.ExpressSurcharge)
}

func TestCalculateShipping_CompleteCalculation_Standard(t *testing.T) {
	// Arrange
	ctx := context.Background()
	service := NewShippingService()
	req := &model.CalculateShippingRequest{
		OriginZipcode:      "01000-000",
		DestinationZipcode: "02000-000",
		Weight:             1.5, // 1.5 kg = 3 units of 0.5 kg
		Dimensions: model.PackageDimensions{
			Length: 20.0,
			Width:  15.0,
			Height: 10.0, // Volume = 3000 cm³ = 3 units of 1000 cm³
		},
		IsExpress: false,
	}

	// Act
	response, err := service.CalculateShipping(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, response)
	// Base cost: depends on distance (10000 difference)
	// Weight surcharge: baseCost * 0.10 * 3.0
	// Volume surcharge: baseCost * 0.05 * 3.0
	// Standard cost should be greater than base cost
	assert.Greater(t, response.ShippingCost, 1000.0)
	assert.Equal(t, "2 dias", response.EstimatedDeliveryTime)
}

func TestCalculateShipping_CompleteCalculation_Express(t *testing.T) {
	// Arrange
	ctx := context.Background()
	service := NewShippingService()
	req := &model.CalculateShippingRequest{
		OriginZipcode:      "01000-000",
		DestinationZipcode: "02000-000",
		Weight:             1.0,
		Dimensions: model.PackageDimensions{
			Length: 10.0,
			Width:  10.0,
			Height: 10.0, // Volume = 1000 cm³
		},
		IsExpress: true,
	}

	// Act
	response, err := service.CalculateShipping(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, response)
	// Express cost should be 50% more than standard
	standardCost := response.ShippingOptions[0].Cost
	expressCost := response.ShippingOptions[1].Cost
	expectedExpressCost := standardCost * 1.5
	assert.Equal(t, expectedExpressCost, expressCost)
	assert.Equal(t, expectedExpressCost, response.ShippingCost)
	assert.Equal(t, "1 dia", response.EstimatedDeliveryTime)
}

func TestCalculateBaseCost_InvalidZipcode_NonNumeric(t *testing.T) {
	// Arrange
	service := NewShippingService()
	originZipcode := "abc"
	destinationZipcode := "def"

	// Act
	baseCost := service.calculateBaseCost(originZipcode, destinationZipcode)

	// Assert
	// Should return default base cost when conversion fails
	assert.Equal(t, 1000.0, baseCost)
}

func TestCalculateBaseCost_InvalidZipcode_Empty(t *testing.T) {
	// Arrange
	service := NewShippingService()
	originZipcode := ""
	destinationZipcode := ""

	// Act
	baseCost := service.calculateBaseCost(originZipcode, destinationZipcode)

	// Assert
	// Should return default base cost when conversion fails
	assert.Equal(t, 1000.0, baseCost)
}

func TestCalculateBaseCost_NegativeDistance(t *testing.T) {
	// Arrange
	service := NewShippingService()
	originZipcode := "20000"
	destinationZipcode := "10000"

	// Act
	baseCost := service.calculateBaseCost(originZipcode, destinationZipcode)

	// Assert
	// Distance is 10000, should have increased base cost
	assert.Greater(t, baseCost, 1000.0)
}

func TestCalculateBaseCost_WithHyphensAndSpaces(t *testing.T) {
	// Arrange
	service := NewShippingService()
	originZipcode := "12 345-678"
	destinationZipcode := "87-654 321"

	// Act
	baseCost := service.calculateBaseCost(originZipcode, destinationZipcode)

	// Assert
	// Should normalize and calculate correctly
	assert.Greater(t, baseCost, 0.0)
}

func TestCalculateShippingDetails_WeightSurcharge_ExactHalfKg(t *testing.T) {
	// Arrange
	service := NewShippingService()
	baseCost := 1000.0
	weight := 0.5 // Exactly 0.5 kg = 1 unit
	volume := 1000.0
	isExpress := false

	// Act
	details := service.calculateShippingDetails(baseCost, weight, volume, isExpress)

	// Assert
	// Weight multiplier: 0.5 / 0.5 = 1.0
	// Weight surcharge: 1000 * 0.10 * 1.0 = 100
	expectedWeightSurcharge := 100.0
	assert.Equal(t, expectedWeightSurcharge, details.WeightSurcharge)
}

func TestCalculateShippingDetails_WeightSurcharge_LessThanHalfKg(t *testing.T) {
	// Arrange
	service := NewShippingService()
	baseCost := 1000.0
	weight := 0.25 // Less than 0.5 kg
	volume := 1000.0
	isExpress := false

	// Act
	details := service.calculateShippingDetails(baseCost, weight, volume, isExpress)

	// Assert
	// Weight multiplier: 0.25 / 0.5 = 0.5
	// Weight surcharge: 1000 * 0.10 * 0.5 = 50
	expectedWeightSurcharge := 50.0
	assert.Equal(t, expectedWeightSurcharge, details.WeightSurcharge)
}

func TestCalculateShippingDetails_VolumeSurcharge_Exact1000Cm3(t *testing.T) {
	// Arrange
	service := NewShippingService()
	baseCost := 1000.0
	weight := 1.0
	volume := 1000.0 // Exactly 1000 cm³ = 1 unit
	isExpress := false

	// Act
	details := service.calculateShippingDetails(baseCost, weight, volume, isExpress)

	// Assert
	// Volume multiplier: 1000 / 1000 = 1.0
	// Volume surcharge: 1000 * 0.05 * 1.0 = 50
	expectedVolumeSurcharge := 50.0
	assert.Equal(t, expectedVolumeSurcharge, details.VolumeSurcharge)
}

func TestCalculateShippingDetails_VolumeSurcharge_LessThan1000Cm3(t *testing.T) {
	// Arrange
	service := NewShippingService()
	baseCost := 1000.0
	weight := 1.0
	volume := 500.0 // Less than 1000 cm³
	isExpress := false

	// Act
	details := service.calculateShippingDetails(baseCost, weight, volume, isExpress)

	// Assert
	// Volume multiplier: 500 / 1000 = 0.5
	// Volume surcharge: 1000 * 0.05 * 0.5 = 25
	expectedVolumeSurcharge := 25.0
	assert.Equal(t, expectedVolumeSurcharge, details.VolumeSurcharge)
}

func TestCalculateShippingDetails_ExpressSurcharge_ZeroSubtotal(t *testing.T) {
	// Arrange
	service := NewShippingService()
	baseCost := 0.0
	weight := 0.0
	volume := 0.0
	isExpress := true

	// Act
	details := service.calculateShippingDetails(baseCost, weight, volume, isExpress)

	// Assert
	assert.Equal(t, 0.0, details.BaseCost)
	assert.Equal(t, 0.0, details.WeightSurcharge)
	assert.Equal(t, 0.0, details.VolumeSurcharge)
	assert.Equal(t, 0.0, details.ExpressSurcharge)
	assert.Equal(t, 0.0, details.TotalCost)
	assert.Equal(t, 1, details.EstimatedDays)
}

func TestBuildResponse_StandardShipping_PluralDays(t *testing.T) {
	// Arrange
	service := NewShippingService()
	details := &model.ShippingCalculationDetails{
		BaseCost:         1000.0,
		WeightSurcharge:  200.0,
		VolumeSurcharge:  50.0,
		ExpressSurcharge: 0.0,
		TotalCost:        1250.0,
		EstimatedDays:    2,
	}
	isExpress := false

	// Act
	response := service.buildResponse(details, isExpress)

	// Assert
	assert.NotNil(t, response)
	assert.Equal(t, "2 dias", response.EstimatedDeliveryTime)
	assert.Equal(t, "2 dias", response.ShippingOptions[0].Time)
}

func TestBuildResponse_ExpressShipping_SingularDay(t *testing.T) {
	// Arrange
	service := NewShippingService()
	details := &model.ShippingCalculationDetails{
		BaseCost:         1000.0,
		WeightSurcharge:  200.0,
		VolumeSurcharge:  50.0,
		ExpressSurcharge: 625.0,
		TotalCost:        1875.0,
		EstimatedDays:    1,
	}
	isExpress := true

	// Act
	response := service.buildResponse(details, isExpress)

	// Assert
	assert.NotNil(t, response)
	assert.Equal(t, "1 dia", response.EstimatedDeliveryTime)
	assert.Equal(t, "1 dia", response.ShippingOptions[1].Time)
}

func TestCalculateShipping_InvalidWeight_Negative(t *testing.T) {
	// Arrange
	ctx := context.Background()
	service := NewShippingService()
	req := &model.CalculateShippingRequest{
		OriginZipcode:      "12345678",
		DestinationZipcode: "87654321",
		Weight:             -1.0,
		Dimensions: model.PackageDimensions{
			Length: 10.0,
			Width:  10.0,
			Height: 10.0,
		},
		IsExpress: false,
	}

	// Act
	response, err := service.CalculateShipping(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "invalid weight")
}

func TestCalculateShipping_InvalidDimensions_NegativeLength(t *testing.T) {
	// Arrange
	ctx := context.Background()
	service := NewShippingService()
	req := &model.CalculateShippingRequest{
		OriginZipcode:      "12345678",
		DestinationZipcode: "87654321",
		Weight:             1.0,
		Dimensions: model.PackageDimensions{
			Length: -1.0,
			Width:  10.0,
			Height: 10.0,
		},
		IsExpress: false,
	}

	// Act
	response, err := service.CalculateShipping(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "invalid dimensions")
}

func TestCalculateShipping_InvalidDimensions_NegativeWidth(t *testing.T) {
	// Arrange
	ctx := context.Background()
	service := NewShippingService()
	req := &model.CalculateShippingRequest{
		OriginZipcode:      "12345678",
		DestinationZipcode: "87654321",
		Weight:             1.0,
		Dimensions: model.PackageDimensions{
			Length: 10.0,
			Width:  -1.0,
			Height: 10.0,
		},
		IsExpress: false,
	}

	// Act
	response, err := service.CalculateShipping(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "invalid dimensions")
}

func TestCalculateShipping_InvalidDimensions_NegativeHeight(t *testing.T) {
	// Arrange
	ctx := context.Background()
	service := NewShippingService()
	req := &model.CalculateShippingRequest{
		OriginZipcode:      "12345678",
		DestinationZipcode: "87654321",
		Weight:             1.0,
		Dimensions: model.PackageDimensions{
			Length: 10.0,
			Width:  10.0,
			Height: -1.0,
		},
		IsExpress: false,
	}

	// Act
	response, err := service.CalculateShipping(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "invalid dimensions")
}
