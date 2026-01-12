package validator

import (
	"fmt"
	"strings"
	"unicode"
)

const (
	maxVolumeCm3     = 15000.0
	minWeight        = 0.0
	zipcodeLength    = 8
	minZipcodeLength = 4
)

// ValidateZipcode validates Brazilian zipcode format without using regex to avoid ReDoS vulnerabilities
func ValidateZipcode(zipcode, fieldName string) error {
	if zipcode == "" {
		return fmt.Errorf("%s is required", fieldName)
	}

	// Normalize zipcode (remove hyphens and spaces)
	normalized := strings.ReplaceAll(strings.ReplaceAll(zipcode, "-", ""), " ", "")

	// Validate length (must be at least 4 digits and at most 8 digits)
	if len(normalized) < minZipcodeLength || len(normalized) > zipcodeLength {
		return fmt.Errorf("%s must be a valid zipcode format (4-8 digits)", fieldName)
	}

	// Validate that all characters are digits (manual check to avoid regex backtracking)
	for _, char := range normalized {
		if !unicode.IsDigit(char) {
			return fmt.Errorf("%s must be a valid zipcode format (4-8 digits)", fieldName)
		}
	}

	return nil
}

// ValidateWeight validates that weight is positive
func ValidateWeight(weight float64) error {
	if weight <= minWeight {
		return fmt.Errorf("weight must be greater than 0")
	}
	return nil
}

// ValidateDimensions validates that dimensions are positive and volume doesn't exceed limit
func ValidateDimensions(length, width, height float64) error {
	if length <= 0 {
		return fmt.Errorf("dimensions.length must be positive")
	}
	if width <= 0 {
		return fmt.Errorf("dimensions.width must be positive")
	}
	if height <= 0 {
		return fmt.Errorf("dimensions.height must be positive")
	}

	volume := length * width * height
	if volume > maxVolumeCm3 {
		return fmt.Errorf("package volume (%.2f cm³) exceeds maximum allowed volume (%.2f cm³)", volume, maxVolumeCm3)
	}

	return nil
}

// CalculateVolume calculates the volume in cm³ from dimensions
func CalculateVolume(length, width, height float64) float64 {
	return length * width * height
}
