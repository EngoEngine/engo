package engi

import (
	"strconv"
)

type World struct {
	Game
	entities []*Entity
	systems  []Systemer

	defaultBatch *Batch
	hudBatch     *Batch

	isSetup bool
}

func (w *World) New() {
	if !w.isSetup {

		w.defaultBatch = NewBatch(Width(), Height(), batchVert, batchFrag)
		w.hudBatch = NewBatch(Width(), Height(), hudVert, hudFrag)

		w.isSetup = true
	}
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
}

func (w *World) Post() {}

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

func (w *World) Batch(prio PriorityLevel) *Batch {
	if prio >= HUDGround {
		return w.hudBatch
	} else {
		return w.defaultBatch
	}
}
