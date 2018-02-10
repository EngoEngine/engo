package common

import (
	"testing"

	"engo.io/ecs"
	"github.com/stretchr/testify/assert"
)

//Test_MinDistance tests, that the calculated step off distance from
// HitBox.MinStepOffD - returns a valid distance pair.
func Test_MinDistance(t *testing.T) {
	a := HitBox{20, 20, 20, 20}
	b := HitBox{8, 9, 20, 20}
	c := HitBox{9, 8, 20, 20}

	//right
	dx, dy := a.MinimumStepOffD(b)

	assert.True(t, dx == 8, "right: Dx wrong = %f", dx)
	assert.True(t, dy == 0, "right: Dy wrony = %f", dy)

	//left
	dx, dy = b.MinimumStepOffD(a)

	assert.True(t, dx == -8, "left wrong = %f", dx)
	assert.True(t, dy == 0, "left wrony = %f", dy)

	//top
	dx, dy = c.MinimumStepOffD(a)
	assert.True(t, dx == 0, "top wrong = %f", dx)
	assert.True(t, dy == -8, "top wrony = %f", dy)
	//bottom
	dx, dy = a.MinimumStepOffD(c)
	assert.True(t, dx == 0, "bottom wrong = %f", dx)
	assert.True(t, dy == 8, "bottom wrong = %f", dy)
}

//hitme is a super simplified hitable object.
//For use of the system, users should combine Components to achieve the required interface
type hitme struct {
	ecs.BasicEntity
	box         HitBox
	main, group HitGroup
}

//Shunt moves the hitbox
//In a full system the SpaceComponent should provide this method
func (hm *hitme) Shunt(dx, dy float32) {
	hm.box.x += dx
	hm.box.y += dy
}

//Simply returns the stored hitbox
//In a full system, this method can be provided by the SpaceComponent,
//But also overridden, by the containing entity
func (hm *hitme) GetHitBox() HitBox {
	return hm.box
}

//HitGroups, Possibly provided by a CollisionComponent at some point.
func (hm *hitme) HitGroups() (HitGroup, HitGroup) {
	return hm.main, hm.group
}

//Test_Update creates a system, sends it objects, calls update once,
//then sees if it has done what it should to the objects
func Test_Update(t *testing.T) {
	hent := func(x, y, w, h float32, gm, gg HitGroup) *hitme {
		nb := ecs.NewBasic()
		return &hitme{
			BasicEntity: nb,
			box: HitBox{
				x: x, y: y, w: w, h: h},
			main:  gm,
			group: gg,
		}
	}

	ts := []*hitme{
		hent(20, 20, 20, 20, 1, 0),
		hent(7, 20, 20, 20, 0, 1),
	}

	sys := HitSystem{Solid: 1}

	for _, v := range ts {
		sys.Add(v)
	}
	sys.Update(0.001)

	hb := ts[0].GetHitBox()
	if hb.x != 27 {
		t.Errorf("No solid collision happened")
	}

}
