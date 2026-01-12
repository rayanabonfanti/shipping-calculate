package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateZipcode_ValidCases(t *testing.T) {
	tests := []struct {
		name      string
		zipcode   string
		fieldName string
	}{
		{
			name:      "valid zipcode with 8 digits",
			zipcode:   "12345678",
			fieldName: "origin_zipcode",
		},
		{
			name:      "valid zipcode with hyphen",
			zipcode:   "12345-678",
			fieldName: "destination_zipcode",
		},
		{
			name:      "valid zipcode with spaces",
			zipcode:   "12345 678",
			fieldName: "origin_zipcode",
		},
		{
			name:      "valid zipcode with 4 digits",
			zipcode:   "1414",
			fieldName: "origin_zipcode",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			// (no setup needed)

			// Act
			err := ValidateZipcode(tt.zipcode, tt.fieldName)

			// Assert
			assert.NoError(t, err)
		})
	}
}

func TestValidateZipcode_InvalidCases(t *testing.T) {
	tests := []struct {
		name        string
		zipcode     string
		fieldName   string
		expectedErr string
	}{
		{
			name:        "empty zipcode",
			zipcode:     "",
			fieldName:   "origin_zipcode",
			expectedErr: "origin_zipcode is required",
		},
		{
			name:        "zipcode with less than 4 digits",
			zipcode:     "123",
			fieldName:   "destination_zipcode",
			expectedErr: "destination_zipcode must be a valid zipcode format (4-8 digits)",
		},
		{
			name:        "zipcode with more than 8 digits",
			zipcode:     "123456789",
			fieldName:   "origin_zipcode",
			expectedErr: "origin_zipcode must be a valid zipcode format (4-8 digits)",
		},
		{
			name:        "zipcode with letters",
			zipcode:     "12345-abc",
			fieldName:   "destination_zipcode",
			expectedErr: "destination_zipcode must be a valid zipcode format (4-8 digits)",
		},
		{
			name:        "zipcode with special characters",
			zipcode:     "12345@678",
			fieldName:   "origin_zipcode",
			expectedErr: "origin_zipcode must be a valid zipcode format (4-8 digits)",
		},
		{
			name:        "zipcode with only letters",
			zipcode:     "abcdefgh",
			fieldName:   "origin_zipcode",
			expectedErr: "origin_zipcode must be a valid zipcode format (4-8 digits)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			// (no setup needed)

			// Act
			err := ValidateZipcode(tt.zipcode, tt.fieldName)

			// Assert
			assert.Error(t, err)
			assert.Equal(t, tt.expectedErr, err.Error())
		})
	}
}

func TestValidateWeight_ValidCases(t *testing.T) {
	tests := []struct {
		name   string
		weight float64
	}{
		{
			name:   "valid positive weight",
			weight: 1.0,
		},
		{
			name:   "valid small weight",
			weight: 0.1,
		},
		{
			name:   "valid large weight",
			weight: 100.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			// (no setup needed)

			// Act
			err := ValidateWeight(tt.weight)

			// Assert
			assert.NoError(t, err)
		})
	}
}

func TestValidateWeight_InvalidCases(t *testing.T) {
	tests := []struct {
		name        string
		weight      float64
		expectedErr string
	}{
		{
			name:        "zero weight",
			weight:      0.0,
			expectedErr: "weight must be greater than 0",
		},
		{
			name:        "negative weight",
			weight:      -1.0,
			expectedErr: "weight must be greater than 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			// (no setup needed)

			// Act
			err := ValidateWeight(tt.weight)

			// Assert
			assert.Error(t, err)
			assert.Equal(t, tt.expectedErr, err.Error())
		})
	}
}

func TestValidateDimensions_ValidCases(t *testing.T) {
	tests := []struct {
		name   string
		length float64
		width  float64
		height float64
	}{
		{
			name:   "valid dimensions",
			length: 10.0,
			width:  10.0,
			height: 10.0,
		},
		{
			name:   "valid large dimensions within limit",
			length: 25.0,
			width:  20.0,
			height: 30.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			// (no setup needed)

			// Act
			err := ValidateDimensions(tt.length, tt.width, tt.height)

			// Assert
			assert.NoError(t, err)
		})
	}
}

func TestValidateDimensions_InvalidCases(t *testing.T) {
	tests := []struct {
		name        string
		length      float64
		width       float64
		height      float64
		expectedErr string
	}{
		{
			name:        "zero length",
			length:      0.0,
			width:       10.0,
			height:      10.0,
			expectedErr: "dimensions.length must be positive",
		},
		{
			name:        "negative length",
			length:      -1.0,
			width:       10.0,
			height:      10.0,
			expectedErr: "dimensions.length must be positive",
		},
		{
			name:        "zero width",
			length:      10.0,
			width:       0.0,
			height:      10.0,
			expectedErr: "dimensions.width must be positive",
		},
		{
			name:        "negative width",
			length:      10.0,
			width:       -1.0,
			height:      10.0,
			expectedErr: "dimensions.width must be positive",
		},
		{
			name:        "zero height",
			length:      10.0,
			width:       10.0,
			height:      0.0,
			expectedErr: "dimensions.height must be positive",
		},
		{
			name:        "negative height",
			length:      10.0,
			width:       10.0,
			height:      -1.0,
			expectedErr: "dimensions.height must be positive",
		},
		{
			name:        "volume exceeds maximum",
			length:      30.0,
			width:       30.0,
			height:      20.0,
			expectedErr: "package volume (18000.00 cm³) exceeds maximum allowed volume (15000.00 cm³)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			// (no setup needed)

			// Act
			err := ValidateDimensions(tt.length, tt.width, tt.height)

			// Assert
			assert.Error(t, err)
			assert.Equal(t, tt.expectedErr, err.Error())
		})
	}
}

func TestCalculateVolume(t *testing.T) {
	tests := []struct {
		name     string
		length   float64
		width    float64
		height   float64
		expected float64
	}{
		{
			name:     "calculate volume for small package",
			length:   10.0,
			width:    10.0,
			height:   10.0,
			expected: 1000.0,
		},
		{
			name:     "calculate volume for medium package",
			length:   20.0,
			width:    15.0,
			height:   10.0,
			expected: 3000.0,
		},
		{
			name:     "calculate volume for large package",
			length:   25.0,
			width:    20.0,
			height:   15.0,
			expected: 7500.0,
		},
		{
			name:     "calculate volume with decimal dimensions",
			length:   10.5,
			width:    5.5,
			height:   2.5,
			expected: 144.375,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			// (no setup needed)

			// Act
			result := CalculateVolume(tt.length, tt.width, tt.height)

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidateZipcode_AllValidFormats(t *testing.T) {
	tests := []struct {
		name      string
		zipcode   string
		fieldName string
	}{
		{
			name:      "8 digits without separator",
			zipcode:   "12345678",
			fieldName: "origin_zipcode",
		},
		{
			name:      "8 digits with hyphen",
			zipcode:   "12345-678",
			fieldName: "destination_zipcode",
		},
		{
			name:      "8 digits with space",
			zipcode:   "12345 678",
			fieldName: "origin_zipcode",
		},
		{
			name:      "8 digits with hyphen and space",
			zipcode:   "12345 - 678",
			fieldName: "destination_zipcode",
		},
		{
			name:      "4 digits minimum",
			zipcode:   "1234",
			fieldName: "origin_zipcode",
		},
		{
			name:      "5 digits",
			zipcode:   "12345",
			fieldName: "destination_zipcode",
		},
		{
			name:      "6 digits",
			zipcode:   "123456",
			fieldName: "origin_zipcode",
		},
		{
			name:      "7 digits",
			zipcode:   "1234567",
			fieldName: "destination_zipcode",
		},
		{
			name:      "8 digits maximum",
			zipcode:   "12345678",
			fieldName: "origin_zipcode",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			// (no setup needed)

			// Act
			err := ValidateZipcode(tt.zipcode, tt.fieldName)

			// Assert
			assert.NoError(t, err)
		})
	}
}

func TestValidateZipcode_AllInvalidFormats(t *testing.T) {
	tests := []struct {
		name        string
		zipcode     string
		fieldName   string
		expectedErr string
	}{
		{
			name:        "3 digits (too short)",
			zipcode:     "123",
			fieldName:   "origin_zipcode",
			expectedErr: "origin_zipcode must be a valid zipcode format (4-8 digits)",
		},
		{
			name:        "9 digits (too long)",
			zipcode:     "123456789",
			fieldName:   "destination_zipcode",
			expectedErr: "destination_zipcode must be a valid zipcode format (4-8 digits)",
		},
		{
			name:        "letters mixed with digits",
			zipcode:     "1234abc",
			fieldName:   "origin_zipcode",
			expectedErr: "origin_zipcode must be a valid zipcode format (4-8 digits)",
		},
		{
			name:        "special characters",
			zipcode:     "12345@67",
			fieldName:   "destination_zipcode",
			expectedErr: "destination_zipcode must be a valid zipcode format (4-8 digits)",
		},
		{
			name:        "only special characters",
			zipcode:     "@#$%^&*",
			fieldName:   "origin_zipcode",
			expectedErr: "origin_zipcode must be a valid zipcode format (4-8 digits)",
		},
		{
			name:        "unicode characters",
			zipcode:     "12345á67",
			fieldName:   "destination_zipcode",
			expectedErr: "destination_zipcode must be a valid zipcode format (4-8 digits)",
		},
		{
			name:        "whitespace only",
			zipcode:     "   ",
			fieldName:   "origin_zipcode",
			expectedErr: "origin_zipcode must be a valid zipcode format (4-8 digits)",
		},
		{
			name:        "newline characters",
			zipcode:     "12345\n67",
			fieldName:   "destination_zipcode",
			expectedErr: "destination_zipcode must be a valid zipcode format (4-8 digits)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			// (no setup needed)

			// Act
			err := ValidateZipcode(tt.zipcode, tt.fieldName)

			// Assert
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func TestValidateWeight_EdgeCases_Valid(t *testing.T) {
	tests := []struct {
		name   string
		weight float64
	}{
		{
			name:   "very small positive weight",
			weight: 0.0001,
		},
		{
			name:   "very large weight",
			weight: 1000000.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			// (no setup needed)

			// Act
			err := ValidateWeight(tt.weight)

			// Assert
			assert.NoError(t, err)
		})
	}
}

func TestValidateWeight_EdgeCases_Invalid(t *testing.T) {
	tests := []struct {
		name        string
		weight      float64
		expectedErr string
	}{
		{
			name:        "exactly zero",
			weight:      0.0,
			expectedErr: "weight must be greater than 0",
		},
		{
			name:        "negative small",
			weight:      -0.0001,
			expectedErr: "weight must be greater than 0",
		},
		{
			name:        "negative large",
			weight:      -1000.0,
			expectedErr: "weight must be greater than 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			// (no setup needed)

			// Act
			err := ValidateWeight(tt.weight)

			// Assert
			assert.Error(t, err)
			assert.Equal(t, tt.expectedErr, err.Error())
		})
	}
}

func TestValidateDimensions_EdgeCases_Valid(t *testing.T) {
	tests := []struct {
		name   string
		length float64
		width  float64
		height float64
	}{
		{
			name:   "very small positive dimensions",
			length: 0.0001,
			width:  0.0001,
			height: 0.0001,
		},
		{
			name:   "dimensions at maximum volume limit",
			length: 25.0,
			width:  20.0,
			height: 30.0,
		},
		{
			name:   "dimensions just below maximum volume",
			length: 24.9,
			width:  20.0,
			height: 30.0,
		},
		{
			name:   "very large dimensions within limit",
			length: 100.0,
			width:  50.0,
			height: 3.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			// (no setup needed)

			// Act
			err := ValidateDimensions(tt.length, tt.width, tt.height)

			// Assert
			assert.NoError(t, err)
		})
	}
}

func TestValidateDimensions_EdgeCases_Invalid(t *testing.T) {
	tests := []struct {
		name        string
		length      float64
		width       float64
		height      float64
		expectedErr string
	}{
		{
			name:        "dimensions just above maximum volume",
			length:      25.1,
			width:       20.0,
			height:      30.0,
			expectedErr: "package volume",
		},
		{
			name:        "all dimensions zero",
			length:      0.0,
			width:       0.0,
			height:      0.0,
			expectedErr: "dimensions.length must be positive",
		},
		{
			name:        "all dimensions negative",
			length:      -1.0,
			width:       -1.0,
			height:      -1.0,
			expectedErr: "dimensions.length must be positive",
		},
		{
			name:        "very large dimensions exceeding limit",
			length:      100.0,
			width:       50.0,
			height:      4.0,
			expectedErr: "package volume",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			// (no setup needed)

			// Act
			err := ValidateDimensions(tt.length, tt.width, tt.height)

			// Assert
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func TestCalculateVolume_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		length   float64
		width    float64
		height   float64
		expected float64
	}{
		{
			name:     "very small dimensions",
			length:   0.1,
			width:    0.1,
			height:   0.1,
			expected: 0.001,
		},
		{
			name:     "very large dimensions",
			length:   100.0,
			width:    100.0,
			height:   100.0,
			expected: 1000000.0,
		},
		{
			name:     "decimal dimensions",
			length:   10.5,
			width:    5.5,
			height:   2.5,
			expected: 144.375,
		},
		{
			name:     "zero dimensions",
			length:   0.0,
			width:    0.0,
			height:   0.0,
			expected: 0.0,
		},
		{
			name:     "negative dimensions",
			length:   -10.0,
			width:    -5.0,
			height:   -2.0,
			expected: -100.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			// (no setup needed)

			// Act
			result := CalculateVolume(tt.length, tt.width, tt.height)

			// Assert
			assert.InDelta(t, tt.expected, result, 0.0001)
		})
	}
}
