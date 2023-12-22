package utils

import "fmt"

// GetFormatString Determine the number of digits needed based on the length
func GetFormatString(length int) string {
	numDigits := 1
	for length >= 10 {
		length /= 10
		numDigits++
	}
	return fmt.Sprintf("%%0%dd", numDigits)
}
