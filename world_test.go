package engi

import (
	"testing"
)

type TestSystem struct {
	*System
}

func (ts *TestSystem) New() {
	ts.System = &System{}
}

func (*TestSystem) Type() string {
	return "TestSystem"
}

func (ts *TestSystem) Update(e *Entity, dt float32) {}

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

func TestAddSystem(t *testing.T) {
	world := World{}
	system := &TestSystem{}
	world.AddSystem(system)

	if len(world.Systems()) != 1 {
		t.Fail()
	}
}

func TestAddComponent(t *testing.T) {
	world := World{}
	world.AddSystem(&TestSystem{})
	entity := NewEntity([]string{"TestSystem"})
	world.AddEntity(entity)
	component := SpaceComponent{Position: Point{0, 10}, Width: 100, Height: 100}
	entity.AddComponent(component)
	if len(entity.components) != 1 {
		t.Fail()
	}
}
