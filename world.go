package engi

import (
	"strconv"
)

type World struct {
	Game
	entities []*Entity
	systems  []Systemer
}

func (w *World) AddEntity(entity *Entity) {
	entity.id = strconv.Itoa(len(w.entities))
	w.entities = append(w.entities, entity)
}

func (w *World) AddSystem(system Systemer) {
	w.systems = append(w.systems, system)
}

func (w *World) Entities() []*Entity {
	return w.entities
}

func (w *World) Systems() []Systemer {
	return w.systems
}

func (w *World) Update(dt float32) {
	for _, entity := range w.Entities() {
		for _, system := range w.Systems() {
			if entity.DoesRequire(system.Name()) {
				system.Pre()
				system.Update(entity, dt)
				system.Post()
			}
		}
	}
}

type Entity struct {
	id         string
	components []Component
	requires   []string
}

func NewEntity(requires []string) *Entity {
	return &Entity{requires: requires}
}

func (e *Entity) DoesRequire(name string) bool {
	for _, requirement := range e.requires {
		if requirement == name {
			return true
		}
	}

	return false
}

func (e *Entity) AddComponent(component Component) {
	e.components = append(e.components, component)
}

func (e *Entity) GetComponent(name string) Component {
	for _, component := range e.components {
		if component.Name() == name {
			return component
		}
	}
	return nil
}

func (e *Entity) ID() string {
	return e.id
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

type Systemer interface {
	Update(entity *Entity, dt float32)
	Name() string
	Priority() int
	Pre()
	Post()
}

type TestSystem struct{}

func (ts TestSystem) Pre()  {}
func (ts TestSystem) Post() {}
func (ts TestSystem) Update(entity *Entity, dt float32) {
	print(entity.ID() + "YOLO\n")
}

func (ts TestSystem) Name() string {
	return "TestSystem"
}
func (ts TestSystem) Priority() int {
	return 0
}
