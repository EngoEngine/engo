package math

import (
	engomath "github.com/EngoEngine/math"
)

// Tanh returns the hyperbolic tangent of x.
//
// Special cases are:
//	Tanh(±0) = ±0
//	Tanh(±Inf) = ±1
//	Tanh(NaN) = NaN
func Tanh(x float32) float32 {
	return engomath.Tanh(x)
}
