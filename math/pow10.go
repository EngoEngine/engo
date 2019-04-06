package math

import (
	engomath "github.com/EngoEngine/math"
)

// Pow10 returns 10**e, the base-10 exponential of e.
//
// Special cases are:
//	Pow10(e) = +Inf for e > 309
//	Pow10(e) = 0 for e < -324
func Pow10(e int) float32 {
	return engomath.Pow10(e)
}
