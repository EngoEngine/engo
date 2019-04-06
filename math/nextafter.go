package math

import (
	engomath "github.com/EngoEngine/math"
)

// Nextafter returns the next representable float32 value after x towards y.
// Special cases:
//	Nextafter32(x, x)   = x
//      Nextafter32(NaN, y) = NaN
//      Nextafter32(x, NaN) = NaN
//
// Since this is a float32 math package the 32 bit version has no number and the
// 64 bit version has the number in the method name.
func Nextafter(x, y float32) float32 {
	return engomath.Nextafter(x, y)
}

// Nextafter64 returns the next representable float64 value after x towards y.
// Special cases:
//      Nextafter64(x, x)   = x
//      Nextafter64(NaN, y) = NaN
//      Nextafter64(x, NaN) = NaN
//
// Since this is a float32 math package the 32 bit version has no number and the
// 64 bit version has the number in the method name.
func Nextafter64(x, y float64) float64 {
	return engomath.Nextafter64(x, y)
}
