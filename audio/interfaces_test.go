//+build !windows,!netgo,!android

package audio

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
	AudioComponent
}

// Test_Every Creates an Everything component and tries to add and then remove it from each system to each system using AddByInterface.
// I can't test adding things that don't work as the code won't compile
func Test_Every(t *testing.T) {
	e := &EveryComp{
		SafeBasic: NewSafeBasic(),
	}

	//Wanted to use a loop to do this, but each "AddByInterface" is actually different nmind

	//AudioSystem
	aus := &AudioSystem{}
	aus.AddByInterface(e)
	if len(aus.entities) != 1 {
		t.Logf("AudioSystem should have 1 entity, got %d", len(aus.entities))
		t.Fail()
	}
	aus.Remove(*e.GetBasicEntity())
	if len(aus.entities) > 0 {
		t.Logf("AudioSystem should now be empty")
		t.Fail()
	}

}
