// Here we are testing that add by interface allows all the various components to be added simply to the various systems.

// Simply compiling should show all the bits are pointing in the right direction

package common

import (
	"engo.io/ecs"
	"engo.io/engo"
	"testing"
)

//Until ecs is updated with GetBasicEntity(), this allows us to give BasicEntity the Interface it Needs
type SafeBasic struct {
	ecs.BasicEntity
}

func (sb *SafeBasic) GetBasicEntity() *ecs.BasicEntity {
	return &sb.BasicEntity
}

type Collidable struct {
	*SafeBasic
	*SpaceComponent
	*CollisionComponent
}

func Test_Collision(t *testing.T) {

	ob1 := Collidable{
		SafeBasic:      &SafeBasic{ecs.NewBasic()},
		SpaceComponent: &SpaceComponent{Width: 50, Height: 50, Position: engo.Point{50, 50}},
		CollisionComponent: &CollisionComponent{
			Solid: true,
			Main:  true,
		},
	}

	ob2 := Collidable{
		SafeBasic:      &SafeBasic{ecs.NewBasic()},
		SpaceComponent: &SpaceComponent{Width: 50, Height: 50, Position: engo.Point{25, 25}},
		CollisionComponent: &CollisionComponent{
			Solid: true,
			Main:  false,
		},
	}

	cs := CollisionSystem{}

	cs.AddByInterface(ob1)
	cs.AddByInterface(ob2)

	if len(cs.entities) < 2 {
		t.Log("Collision system should have 2 entites, got %d", len(cs.entities))
		t.Fail()

	}

}
