package engi

import (
	"strconv"
)

type World struct {
	Game
	entities []*Entity
	systems  []System
}

func (w *World) AddEntity(entity *Entity) {
	entity.id = strconv.Itoa(len(w.entities))
	w.entities = append(w.entities, entity)
}

func (w *World) AddSystem(system System) {
	w.systems = append(w.systems, system)
}

func (w *World) Entities() []*Entity {
	return w.entities
}

func (w *World) Systems() []System {
	return w.systems
}

func (w *World) Update(dt float32) {
	for _, system := range w.Systems() {
		system.Update(dt)
	}

	for i, entity := range w.Entities() {
		if entity.id == "" {
			w.entities[i] = nil
		}
	}
}

type Entity struct {
	id         string
	components []Component
}

func (e *Entity) AddComponent(component Component) {
	e.components = append(e.components, component)
}

type Component interface {
	Name() string
}

type PositionComponent struct {
	X, Y int
}

func (pc PositionComponent) Name() string {
	return "Position"
}

type System interface {
	Update(dt float32)
	Name() string
	Priority() int
}

type TestSystem struct{}

func (ts TestSystem) Update(dt float32) {}

func (ts TestSystem) Name() string {
	return "TestSystem"
}
func (ts TestSystem) Priority() int {
	return 0
}
