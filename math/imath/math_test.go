package imath_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/EngoEngine/engo/math/imath"
)

func TestSqrt(t *testing.T) {
	assert.Equal(t, 5, imath.Sqrt(25))
	assert.Equal(t, 4, imath.Sqrt(24))
	assert.Equal(t, 20, imath.Sqrt(400))
	assert.Equal(t, 19, imath.Sqrt(399))
}
