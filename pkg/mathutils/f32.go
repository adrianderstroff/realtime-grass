// Package mathutils provides utility functions for scalar and vectorial math.
package mathutils

import "math"

// MaxF32 is a float32 wrapper for the float64 function math.Max.
func MaxF32(vala, valb float32) float32 {
	return float32(math.Max(float64(vala), float64(valb)))
}

// AbsF32 is a float32 wrapper for the float64 function math.Abs.
func AbsF32(val float32) float32 {
	return float32(math.Abs(float64(val)))
}

// RoundF32 is a float32 wrapper for the float64 function math.Round.
func RoundF32(val float32) float32 {
	return float32(math.Round(float64(val)))
}

// CeilF32 is a float32 wrapper for the float64 function math.Ceil.
func CeilF32(val float32) float32 {
	return float32(math.Ceil(float64(val)))
}

// SqrtF32 is a float32 wrapper for the float64 function math.Sqrt.
func SqrtF32(val float32) float32 {
	return float32(math.Sqrt(float64(val)))
}
