//+build netgo

package math

import (
	"math"
)

// Cos returns the cosine of the radian argument x.
//
// Special cases are:
//	Cos(±Inf) = NaN
//	Cos(NaN) = NaN
func Cos(x float32) float32 {
	return float32(math.Cos(float64(x)))
}

// Sin returns the sine of the radian argument x.
//
// Special cases are:
//	Sin(±0) = ±0
//	Sin(±Inf) = NaN
//	Sin(NaN) = NaN
func Sin(x float32) float32 {
	return float32(math.Sin(float64(x)))
}
