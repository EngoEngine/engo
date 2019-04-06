package math

import (
	engomath "github.com/EngoEngine/math"
)

// Cos returns the cosine of the radian argument x.
//
// Special cases are:
//	Cos(±Inf) = NaN
//	Cos(NaN) = NaN
func Cos(x float32) float32 {
	return engomath.Cos(x)
}

// Sin returns the sine of the radian argument x.
//
// Special cases are:
//	Sin(±0) = ±0
//	Sin(±Inf) = NaN
//	Sin(NaN) = NaN
func Sin(x float32) float32 {
	return engomath.Sin(x)
}
