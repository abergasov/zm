package utils_test

import (
	"testing"
	"zm/internal/utils"

	"github.com/stretchr/testify/require"
)

func TestGetFormatString(t *testing.T) {
	tests := []struct {
		name   string
		length int
		want   string
	}{
		{"Single Digit", 5, "%01d"},
		{"Double Digits", 15, "%02d"},
		{"Triple Digits", 123, "%03d"},
		{"Four Digits", 9876, "%04d"},
		{"Zero", 0, "%01d"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, utils.GetFormatString(tt.length))
		})
	}
	t.Run("should return %01d for negative length", func(t *testing.T) {
		negativeLength := -5
		wantNegative := "%01d"
		gotNegative := utils.GetFormatString(negativeLength)
		require.Equal(t, gotNegative, wantNegative)
	})
	t.Run("should return %09d for very large length", func(t *testing.T) {
		veryLargeLength := 987654321
		wantVeryLarge := "%09d"
		gotVeryLarge := utils.GetFormatString(veryLargeLength)
		require.Equal(t, gotVeryLarge, wantVeryLarge)
	})
}
