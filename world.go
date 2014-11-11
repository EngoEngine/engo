package engi

import (
	"strconv"
)

type World struct {
	entities []*Entity
	systems  []*System
}

func (w *World) AddEntity(entity *Entity) {
	entity.id = strconv.Itoa(len(w.entities))
	w.entities = append(w.entities, entity)
}

func (w *World) AddSystem() {

}

func (w *World) Entities() []*Entity {
	return w.entities
}

func (w *World) Systems() []*System {
	return w.systems
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

type System struct {
}
