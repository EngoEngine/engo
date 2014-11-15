package engi

import (
	"log"
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
	log.Println(w.systems)
	for i, system := range w.systems {
		log.Println(i, system)
		system.AddEntity(entity)
		// w.systems[i] = system
	}
}

func (w *World) AddSystem(system Systemer) {
	system.New()
	w.systems = append(w.systems, system)
}

func (w *World) Entities() []*Entity {
	return w.entities
}

func (w *World) Systems() []Systemer {
	return w.systems
}

func (w *World) Update(dt float32) {
	for _, system := range w.Systems() {
		system.Pre()
		for _, entity := range system.Entities() {
			system.Update(entity, dt)
		}
		system.Post()
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
	New()
	Entities() []*Entity
	AddEntity(entity *Entity)
}

type System struct {
	entities []*Entity
}

func (s System) New()  {}
func (s System) Pre()  {}
func (s System) Post() {}

func (s System) Priority() int {
	return 0
}

func (s System) Entities() []*Entity {
	return s.entities
}

func (s *System) AddEntity(entity *Entity) {
	// log.Println(entity)
	s.entities = append(s.entities, entity)
}

type SpaceComponent struct {
	Position Point
	Width    float32
	Height   float32
}

func (sc SpaceComponent) Name() string {
	return "SpaceComponent"
}
