package common

import (
	"fmt"
	"testing"

	"engo.io/ecs"
	"engo.io/engo"
	"github.com/stretchr/testify/assert"
)

func TestSpaceComponent_Contains(t *testing.T) {
	space := SpaceComponent{Width: 100, Height: 100}
	pass := []engo.Point{
		engo.Point{X: 10, Y: 10},
		engo.Point{X: 50, Y: 50},
		engo.Point{X: 10, Y: 50},
		engo.Point{X: 99, Y: 99},
	}
	fail := []engo.Point{
		// Totally not within:
		engo.Point{X: -10, Y: -10},
		engo.Point{X: 120, Y: 120},

		// Only one axis within:
		engo.Point{X: 50, Y: 120},
		engo.Point{X: 120, Y: 50},

		// On the edge:
		engo.Point{X: 0, Y: 0},
		engo.Point{X: 0, Y: 50},
		engo.Point{X: 0, Y: 100},
		engo.Point{X: 50, Y: 0},
		engo.Point{X: 50, Y: 100},
		engo.Point{X: 100, Y: 0},
		engo.Point{X: 100, Y: 50},
		engo.Point{X: 100, Y: 100},
	}

	for _, p := range pass {
		assert.True(t, space.Contains(p), fmt.Sprintf("point %v should be within area", p))
	}

	for _, f := range fail {
		assert.False(t, space.Contains(f), fmt.Sprintf("point %v should not be within area", f))
	}
}

func TestSpaceComponent_Corners(t *testing.T) {
	space1 := SpaceComponent{Width: 1, Height: 1}
	exp1 := [4]engo.Point{engo.Point{X: 0, Y: 0}, engo.Point{X: 1, Y: 0}, engo.Point{X: 0, Y: 1}, engo.Point{X: 1, Y: 1}}
	act1 := space1.Corners()
	for i := 0; i < 4; i++ {
		assert.True(t, exp1[i].Equal(act1[i]), fmt.Sprintf("corner %d did not match for rotation %f (got %v expected %v)", i, space1.Rotation, act1[i], exp1[i]))
	}

	space2 := SpaceComponent{Width: 1, Height: 1, Rotation: 90}
	exp2 := [4]engo.Point{engo.Point{X: 0, Y: 0}, engo.Point{X: 0, Y: 1}, engo.Point{X: -1, Y: 0}, engo.Point{X: -1, Y: 1}}
	act2 := space2.Corners()
	for i := 0; i < 4; i++ {
		assert.True(t, exp2[i].Equal(act2[i]), fmt.Sprintf("corner %d did not match for rotation %f (got %v expected %v)", i, space2.Rotation, act2[i], exp2[i]))
	}
}

const (
	Ball = 1 << iota
	Bat
)

//Test GroupSolid working
func Test_GroupSolid(t *testing.T) {
	//All items in same place, have to collide
	CE := func(m, g CollisionGroup) collisionEntity {
		nb := ecs.NewBasic()
		return collisionEntity{
			BasicEntity: &nb,
			CollisionComponent: &CollisionComponent{
				Main:  m,
				Group: g,
			},
			//All objects in same position
			SpaceComponent: &SpaceComponent{engo.Point{X: 10, Y: 10}, 50, 50, 0},
		}
	}
	ents := []collisionEntity{
		CE(Ball, 0),     //The Ball
		CE(Bat, Ball),   //The Batt
		CE(0, Ball|Bat), //The Wall
		CE(0, 0),        //Ghost Should not collide with anything
	}
	sys := CollisionSystem{
		entities: ents,
		Solids:   Ball, //Only the ball should move as Solid
	}
	sys.Update(0.01)

	if ents[0].Position == ents[1].Position {
		t.Log("Ball should collide Solid")
		t.Fail()
	}

	if ents[3].Collides != 0 {
		t.Log("Ghost should not collide with anything")
		t.Fail()
	}
	for i := 0; i < 2; i++ { //Ball and Bat
		if ents[i].Collides == 0 {
			t.Logf("object %d should collides", i)
			t.Fail()
		}
	}
}

func TestSpaceComponent_Center(t *testing.T) {
	components := []SpaceComponent{
		SpaceComponent{Width: 0, Height: 0},
		SpaceComponent{Width: 100, Height: 100},
		SpaceComponent{Width: 100, Height: 200},
	}
	points := []engo.Point{
		engo.Point{X: 10, Y: 10},
		engo.Point{X: 50, Y: 50},
		engo.Point{X: 10, Y: 50},
		engo.Point{X: 99, Y: 99},
	}

	for _, sc := range components {
		for _, p := range points {
			sc.SetCenter(p)
			c := sc.Center()
			assert.True(t, c.Equal(p), fmt.Sprintf("center %v should be equal to point %v", c, p))
		}
	}
}
