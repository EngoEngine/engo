package common

import (
	"fmt"
	"testing"

	"engo.io/engo"
	"github.com/stretchr/testify/assert"
)

func TestSpaceComponent_Within(t *testing.T) {
	space := SpaceComponent{Width: 100, Height: 100}
	pass := []engo.Point{
		engo.Point{10, 10},
		engo.Point{50, 50},
		engo.Point{10, 50},
		engo.Point{99, 99},
	}
	fail := []engo.Point{
		// Totally not within:
		engo.Point{-10, -10},
		engo.Point{120, 120},

		// Only one axis within:
		engo.Point{50, 120},
		engo.Point{120, 50},

		// On the edge:
		engo.Point{0, 0},
		engo.Point{0, 50},
		engo.Point{0, 100},
		engo.Point{50, 0},
		engo.Point{50, 100},
		engo.Point{100, 0},
		engo.Point{100, 50},
		engo.Point{100, 100},
	}

	for _, p := range pass {
		assert.True(t, space.Within(p), fmt.Sprintf("point %v should be within area", p))
	}

	for _, f := range fail {
		assert.False(t, space.Within(f), fmt.Sprintf("point %v should not be within area", f))
	}
}

func TestSpaceComponent_Corners(t *testing.T) {
	space1 := SpaceComponent{Width: 1, Height: 1}
	exp1 := [4]engo.Point{engo.Point{0, 0}, engo.Point{1, 0}, engo.Point{0, 1}, engo.Point{1, 1}}
	act1 := space1.Corners()
	for i := 0; i < 4; i++ {
		assert.True(t, exp1[i].Equal(act1[i]), fmt.Sprintf("corner %d did not match for rotation %f (got %v expected %v)", i, space1.Rotation, act1[i], exp1[i]))
	}

	space2 := SpaceComponent{Width: 1, Height: 1, Rotation: 90}
	exp2 := [4]engo.Point{engo.Point{0, 0}, engo.Point{0, 1}, engo.Point{-1, 0}, engo.Point{-1, 1}}
	act2 := space2.Corners()
	for i := 0; i < 4; i++ {
		assert.True(t, exp2[i].Equal(act2[i]), fmt.Sprintf("corner %d did not match for rotation %f (got %v expected %v)", i, space2.Rotation, act2[i], exp2[i]))
	}
}
