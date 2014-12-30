package engi

import (
	"strconv"
)

type World struct {
	Game
	entities []*Entity
	systems  []Systemer
	batch    *Batch
}

func (w *World) New() {
	w.batch = NewBatch(Width(), Height())
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

func (w *World) Pre() {
	Gl.Clear(Gl.COLOR_BUFFER_BIT)
	w.batch.Begin()
}

func (w *World) Post() {
	w.batch.End()
}

func (w *World) Update(dt float32) {
	w.Pre()
	Cam.Update(dt)
	for _, system := range w.Systems() {
		system.Pre()
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

	w.Post()
}

func (w *World) Batch() *Batch {
	return w.batch
}
