package engi

import (
	"log"
	"testing"
)

type TestSystem struct {
	*System
}

func (ts *TestSystem) New() {
	ts.System = NewSystem()
}

func (*TestSystem) Type() string {
	return "TestSystem"
}

func (ts *TestSystem) Update(e *Entity, dt float32) {}

func TestAddEntity(t *testing.T) {
	headless = true
	world := World{}
	world.new()
	entityOne := NewEntity(nil)
	world.AddEntity(entityOne)
	entityTwo := NewEntity(nil)
	world.AddEntity(entityTwo)
	if len(world.Entities()) != 2 {
		log.Printf("Entities not added.  %d != 2: %+v\n", len(world.Entities()), world.Entities())
		t.Fail()
	}
}

func TestAddSystem(t *testing.T) {
	headless = true
	world := World{}
	world.new()

	before := len(world.Systems())

	system := &TestSystem{}
	world.AddSystem(system)

	if len(world.Systems()) != before+1 {
		t.Fail()
	}
}

func TestAddComponent(t *testing.T) {
	headless = true
	world := World{}
	world.new()
	world.AddSystem(&TestSystem{})
	entity := NewEntity([]string{"TestSystem"})
	world.AddEntity(entity)
	component := &SpaceComponent{Position: Point{0, 10}, Width: 100, Height: 100}
	entity.AddComponent(component)
	if len(entity.components) != 1 {
		t.Fail()
	}
}
