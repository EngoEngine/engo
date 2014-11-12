package engi

import (
	"testing"
)

func TestAddEntity(t *testing.T) {
	world := World{}
	entityOne := Entity{}
	world.AddEntity(&entityOne)
	entityTwo := Entity{}
	world.AddEntity(&entityTwo)
	if len(world.Entities()) != 2 {
		t.Fail()
	}
}

func TestAddComponent(t *testing.T) {
	world := World{}
	entity := Entity{}
	world.AddEntity(&entity)
	component := PositionComponent{0, 10}
	entity.AddComponent(component)
	if len(entity.components) != 1 {
		t.Fail()
	}
}

func TestAddSystem(t *testing.T) {
	world := World{}
	system := TestSystem{}
	world.AddSystem(system)

	if len(world.Systems()) != 1 {
		t.Fail()
	}
}
