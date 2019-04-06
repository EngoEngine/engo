package math

import (
	engomath "github.com/EngoEngine/math"
)

// Inf returns positive infinity if sign >= 0, negative infinity if sign < 0.
func Inf(sign int) float32 {
	return engomath.Inf(sign)
}

// NaN returns an IEEE 754 ``not-a-number'' value.
func NaN() float32 {
	return engomath.NaN()
}

// IsNaN reports whether f is an IEEE 754 ``not-a-number'' value.
func IsNaN(f float32) bool {
	return engomath.IsNaN(f)
}

// IsInf reports whether f is an infinity, according to sign.
// If sign > 0, IsInf reports whether f is positive infinity.
// If sign < 0, IsInf reports whether f is negative infinity.
// If sign == 0, IsInf reports whether f is either infinity.
func IsInf(f float32, sign int) bool {
	return engomath.IsInf(f, sign)
}
