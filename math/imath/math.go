package imath

import (
	engoimath "github.com/EngoEngine/math/imath"
)

// Integer limit values.
const (
	MaxInt8   = 1<<7 - 1
	MinInt8   = -1 << 7
	MaxInt16  = 1<<15 - 1
	MinInt16  = -1 << 15
	MaxInt32  = 1<<31 - 1
	MinInt32  = -1 << 31
	MaxInt64  = 1<<63 - 1
	MinInt64  = -1 << 63
	MaxUint8  = 1<<8 - 1
	MaxUint16 = 1<<16 - 1
	MaxUint32 = 1<<32 - 1
	MaxUint64 = 1<<64 - 1
)

// Abs returns the absolute value of x.
func Abs(x int) int {
	return engoimath.Abs(x)
}

// Cbrt returns the cube root of x.
func Cbrt(x int) int {
	return engoimath.Cbrt(x)
}

// Copysign returns a value with the magnitude of x and the sign of y.
func Copysign(x, y int) int {
	return engoimath.Copysign(x, y)
}

// Dim returns the maximum of x-y or 0.
func Dim(x, y int) int {
	return engoimath.Dim(x, y)
}

// Exp2 returns 2**x, the base-2 exponential of x.
func Exp2(x int) int {
	return engoimath.Exp2(x)
}

// Intbits return the binary representation of i.
func Intbits(i int) uint {
	return engoimath.Intbits(i)
}

// Intfrombits returns the int represented from b.
func Intfrombits(b uint) int {
	return engoimath.Intfrombits(b)
}

// Hypot returns Sqrt(p*p + q*q), taking care to avoid unnecessary overflow and
// underflow.
func Hypot(p, q int) int {
	return engoimath.Hypot(p, q)
}

// Log returns the natural logarithm of x.
func Log(x int) int {
	return engoimath.Log(x)
}

// Max returns the larger of x or y.
func Max(x, y int) int {
	return engoimath.Max(x, y)
}

// Min returns the smaller of x or y.
func Min(x, y int) int {
	return engoimath.Min(x, y)
}

// Mod returns the x%y.
func Mod(x, y int) int {
	return engoimath.Mod(x, y)
}

// Nextafter returns the next representable int value after x towards y.
func Nextafter(x, y int) (r int) {
	return engoimath.Nextafter(x, y)
}

// Pow returns x**y, the base-x exponential of y.
func Pow(x, y int) int {
	return engoimath.Pow(x, y)
}

// Pow10 returns 10**e, the base-10 exponential of e.
func Pow10(e int) int {
	return engoimath.Pow10(e)
}

// Signbit returns true if x is negative or negative zero.
func Signbit(x int) bool {
	return engoimath.Signbit(x)
}

// Sqrt returns the square root of x.
func Sqrt(x int) int {
	return engoimath.Sqrt(x)
}
