package math

import (
	engomath "github.com/engoengine/math"
)

// Signbit returns true if x is negative or negative zero.
func Signbit(x float32) bool {
	return engomath.Signbit(x)
}
