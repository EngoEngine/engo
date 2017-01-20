// Here we are testing that add by interface allows all the various components to be added simply to the various systems.

// Simply compiling should show all the bits are pointing in the right direction

package common

import (
	"testing"
)

//Until ecs is updated, this allows us to give BasicEntity the Interface it Needs
type SafeBasic struct {
	*ecs.BasicEntity
}

func (sb *SafeBasic) GetBasicEntity() *BasicEntity {
	return sb.BasicEntity
}

type alltoall struct {
	*SafeBasic
	*AnimationComponent
	*AudioComponent
	*SpaceComponent
	*RenderComponent
	*MouseComponent
}

func Test_Adds(t *testing.T) {

	ob := alltoall{
		SafeBasic: SafeBasic{ecs.NewBasic()},
	}

}
