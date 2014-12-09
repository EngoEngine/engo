package engi

import (
	// "log"
	"strconv"
)

type World struct {
	Game
	entities []*Entity
	systems  []Systemer
	batch    *Batch
}

func (w *World) AddEntity(entity *Entity) {
	entity.id = strconv.Itoa(len(w.entities))
	w.entities = append(w.entities, entity)
	for _, system := range w.systems {
		if entity.DoesRequire(system.Name()) {
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
	Cam.Update(dt)
	for _, system := range w.Systems() {
		system.Pre()
		for _, message := range system.Messages() {
			system.Receive(message)
		}

		for len(system.Messages()) != 0 {
			system.Dismiss(0)
		}

		for _, entity := range system.Entities() {
			if entity.Exists {
				system.Update(entity, dt)
			}
		}
		system.Post()
	}

	if Keys.KEY_ESCAPE.JustPressed() {
		Exit()
	}
}

func (w *World) Batch() *Batch {
	return w.batch
}
