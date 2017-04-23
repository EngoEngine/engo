package imath

import (
	"testing"
)

func TestSqrt(t *testing.T) {
	t.Log(Sqrt(25))
	t.Log(Sqrt(24))
	t.Log(Sqrt(400))
	t.Log(Sqrt(399))
}
