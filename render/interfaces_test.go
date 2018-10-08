//+build !windows,!netgo,!android

package render

import (
	"testing"

	"engo.io/ecs"
)

// SafeBasic is to provide a BasicEntity until ecs is updated to provide the GetBasicEntity method
type SafeBasic struct {
	ecs.BasicEntity
}

//Satisfy BasicFace interface
func (sb *SafeBasic) GetBasicEntity() *ecs.BasicEntity {
	return &sb.BasicEntity
}

func NewSafeBasic() SafeBasic {
	return SafeBasic{
		BasicEntity: ecs.NewBasic(),
	}
}

type EveryComp struct {
	SafeBasic
	AnimationComponent
	MouseComponent
	RenderComponent
	common2.SpaceComponent
}

// Test_Every Creates an Everything component and tries to add and then remove it from each system to each system using AddByInterface.
// I can't test adding things that don't work as the code won't compile
func Test_Every(t *testing.T) {
	e := &EveryComp{
		SafeBasic: NewSafeBasic(),
	}

	//Wanted to use a loop to do this, but each "AddByInterface" is actually different nmind

	//AnimationSystem
	as := &AnimationSystem{}
	as.AddByInterface(e)
	if len(as.entities) != 1 {
		t.Logf("AnimationSystem should have 1 entity, got %d", len(as.entities))
		t.Fail()
	}
	as.Remove(*e.GetBasicEntity())
	if len(as.entities) > 0 {
		t.Logf("AnimationSystem should now be empty")
		t.Fail()
	}

	//MouseSystem
	ms := &MouseSystem{}
	ms.AddByInterface(e)
	if len(ms.entities) != 1 {
		t.Logf("MouseSystem should have 1 entity, got %d", len(ms.entities))
		t.Fail()
	}
	ms.Remove(*e.GetBasicEntity())
	if len(ms.entities) > 0 {
		t.Logf("MouseSystem should now be empty")
		t.Fail()
	}

	//RenderSystem
	rs := &RenderSystem{}
	rs.AddByInterface(e)
	if len(rs.entities) != 1 {
		t.Logf("RenderSystem should have 1 entity, got %d", len(rs.entities))
		t.Fail()
	}
	rs.Remove(*e.GetBasicEntity())
	if len(rs.entities) > 0 {
		t.Logf("RenderSystem should now be empty")
		t.Fail()
	}

}
