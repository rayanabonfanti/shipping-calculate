package model

// CalculateShippingRequest represents the input for shipping calculation
type CalculateShippingRequest struct {
	OriginZipcode      string            `json:"origin_zipcode"`
	DestinationZipcode string            `json:"destination_zipcode"`
	Weight             float64           `json:"weight"`
	Dimensions         PackageDimensions `json:"dimensions"`
	IsExpress          bool              `json:"is_express"`
}

// PackageDimensions represents package dimensions in centimeters
type PackageDimensions struct {
	Length float64 `json:"length"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

// CalculateShippingResponse represents the output of shipping calculation
type CalculateShippingResponse struct {
	ShippingCost          float64          `json:"shipping_cost"`
	EstimatedDeliveryTime string           `json:"estimated_delivery_time"`
	AvailableServices     []string         `json:"available_services"`
	ShippingOptions       []ShippingOption `json:"shipping_options"`
}

// ShippingOption represents a shipping service option
type ShippingOption struct {
	Service string  `json:"service"`
	Cost    float64 `json:"cost"`
	Time    string  `json:"time"`
}

// ShippingCalculationDetails holds internal calculation details
type ShippingCalculationDetails struct {
	BaseCost         float64
	WeightSurcharge  float64
	VolumeSurcharge  float64
	ExpressSurcharge float64
	TotalCost        float64
	EstimatedDays    int
}
