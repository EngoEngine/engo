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
	for i, system := range w.systems {
		if entity.DoesRequire(system.Name()) {
			log.Println(i, system)
			system.AddEntity(entity)
		}
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

	if Keys.KEY_ESCAPE.JustPressed() {
		Exit()
	}
}
