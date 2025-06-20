package udecimal

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShiftPointLeft(t *testing.T) {
	testcases := []struct {
		input    string
		shift    uint8
		expected string
	}{
		// Basic shifts
		{"123.45", 0, "123.45"},
		{"123.45", 1, "1234.5"},
		{"123.45", 2, "12345"},
		{"123.45", 3, "123450"},
		{"123.45", 4, "1234500"},

		// Shift with zero
		{"0", 1, "0"},
		{"0", 5, "0"},
		{"0.123", 3, "123"},
		{"0.123", 4, "1230"},

		// Negative numbers
		{"-123.45", 1, "-1234.5"},
		{"-123.45", 2, "-12345"},
		{"-0.123", 3, "-123"},

		// Edge cases with precision
		{"1.23456789012345678", 1, "12.3456789012345678"},
		{"1.23456789012345678", 5, "123456.789012345678"},

		// Large shifts beyond decimal places
		{"1.23", 5, "123000"},
		{"0.001", 3, "1"},
		{"0.001", 4, "10"},

		// Single digit cases
		{"1", 1, "10"},
		{"1", 3, "1000"},
		{"0.1", 1, "1"},
		{"0.01", 2, "1"},
	}

	for _, tc := range testcases {
		t.Run(tc.input+"_shift_"+string(rune(tc.shift+'0')), func(t *testing.T) {
			d := MustParse(tc.input)
			result := d.ShiftPointLeft(tc.shift)
			expected := MustParse(tc.expected)

			require.True(t, result.Equal(expected),
				"ShiftPointLeft(%d) on %s: got %s, expected %s",
				tc.shift, tc.input, result.String(), tc.expected)
		})
	}
}

func TestShiftPointRight(t *testing.T) {
	testcases := []struct {
		input    string
		shift    uint8
		expected string
	}{
		// Basic shifts
		{"123.45", 0, "123.45"},
		{"123.45", 1, "12.345"},
		{"123.45", 2, "1.2345"},
		{"123.45", 3, "0.12345"},
		{"123.45", 4, "0.012345"},
		{"123.45", 5, "0.0012345"},

		// Shift with zero
		{"0", 1, "0"},
		{"0", 5, "0"},

		// Negative numbers
		{"-123.45", 1, "-12.345"},
		{"-123.45", 2, "-1.2345"},
		{"-123.45", 3, "-0.12345"},

		// Integer shifts
		{"12345", 1, "1234.5"},
		{"12345", 2, "123.45"},
		{"12345", 5, "0.12345"},
		{"1", 1, "0.1"},
		{"1", 2, "0.01"},

		// Large numbers
		{"123456789", 3, "123456.789"},
		{"123456789", 6, "123.456789"},

		// Edge cases
		{"0.1", 1, "0.01"},
		{"0.01", 1, "0.001"},
	}

	for _, tc := range testcases {
		t.Run(tc.input+"_shift_"+string(rune(tc.shift+'0')), func(t *testing.T) {
			d := MustParse(tc.input)
			result := d.ShiftPointRight(tc.shift)
			expected := MustParse(tc.expected)

			require.True(t, result.Equal(expected),
				"ShiftPointRight(%d) on %s: got %s, expected %s",
				tc.shift, tc.input, result.String(), tc.expected)
		})
	}
}

func TestShiftPointRoundTrip(t *testing.T) {
	testcases := []struct {
		input  string
		shifts []uint8 // Only test shifts that won't hit precision limits
	}{
		{"123.45", []uint8{1, 2}},
		{"0.123", []uint8{1, 2, 3}},
		{"1234567890", []uint8{1, 2, 3, 5, 10}},
		{"0.000001", []uint8{1, 2, 3, 5}},
		{"-123.45", []uint8{1, 2}},
		{"1.23456789012345678", []uint8{1, 2}}, // Limited shifts for high precision numbers
	}

	for _, tc := range testcases {
		for _, shift := range tc.shifts {
			t.Run(tc.input+"_shift_"+string(rune(shift+'0')), func(t *testing.T) {
				original := MustParse(tc.input)

				// Only test round trip if it won't exceed precision limits
				if original.PrecUint()+shift <= maxPrec {
					// Shift right then left should return to original (within precision limits)
					shifted := original.ShiftPointRight(shift).ShiftPointLeft(shift)

					// For this test, we trim trailing zeros from both to compare values
					originalTrimmed := original.trimTrailingZeros()
					shiftedTrimmed := shifted.trimTrailingZeros()

					require.True(t, originalTrimmed.Equal(shiftedTrimmed),
						"Round trip failed for %s with shift %d: got %s",
						tc.input, shift, shifted.String())
				}
			})
		}
	}
}

func TestShiftPointPrecisionLimits(t *testing.T) {
	// Test that shifting doesn't exceed maxPrec
	d := MustParse("1.23")

	// Shift right by a large amount
	result := d.ShiftPointRight(25) // This should be clamped to maxPrec
	require.True(t, result.PrecUint() <= maxPrec,
		"Precision should not exceed maxPrec, got %d", result.PrecUint())

	// Test large left shifts
	largeShift := d.ShiftPointLeft(20)
	require.NotNil(t, largeShift, "Large left shift should not panic")
}

func TestShiftPointZero(t *testing.T) {
	zero := MustParse("0")

	// Shifting zero should always return zero
	leftShift := zero.ShiftPointLeft(5)
	rightShift := zero.ShiftPointRight(5)

	require.True(t, leftShift.IsZero(), "Left shifting zero should return zero")
	require.True(t, rightShift.IsZero(), "Right shifting zero should return zero")
}

func TestShiftPointExtremeValues(t *testing.T) {
	// Test with very small numbers
	small := MustParse("0.000000000000000001")
	shifted := small.ShiftPointLeft(18)
	expected := MustParse("1")
	require.True(t, shifted.Equal(expected),
		"Shifting small decimal: got %s, expected %s", shifted.String(), expected.String())

	// Test with numbers at precision limit
	precise := MustParse("1.2345678901234567890")
	shiftedPrecise := precise.ShiftPointRight(1)
	require.True(t, shiftedPrecise.PrecUint() <= maxPrec,
		"Precision should be within limits after shift")
}
