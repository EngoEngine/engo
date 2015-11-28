package engi

import (
	"runtime"
	"sort"
)

type World struct {
	entities map[string]*Entity
	systems  Systemers

	isSetup bool
	paused  bool
	serial  bool
}

func (w *World) new() {
	if w.isSetup {
		return
	}
	w.entities = make(map[string]*Entity)

	// Default WorldBounds values
	WorldBounds.Max = Point{Width(), Height()}

	// Initialize cameraSystem
	cam = &cameraSystem{}
	w.AddSystem(cam)

	// Short-circuit bypass if there's only 1 core
	if runtime.NumCPU() == 1 {
		w.serial = true
	} else {
		w.serial = false
	}
	w.isSetup = true
}

func (w *World) AddEntity(entity *Entity) {
	w.entities[entity.ID()] = entity

	for _, system := range w.systems {
		if entity.DoesRequire(system.Type()) {
			system.AddEntity(entity)
		}
	}
}

func (w *World) RemoveEntity(entity *Entity) {
	for _, system := range w.systems {
		if entity.DoesRequire(system.Type()) {
			system.RemoveEntity(entity)
		}
	}
	delete(w.entities, entity.ID())
}

func (w *World) AddSystem(system Systemer) {
	system.New(w)
	w.systems = append(w.systems, system)
	sort.Sort(w.systems)
}

func (w *World) Entities() []*Entity {
	entities := make([]*Entity, 0, len(w.entities))
	for _, v := range w.entities {
		entities = append(entities, v)
	}
	return entities
}

func (w *World) Systems() []Systemer {
	return w.systems
}

func (w *World) pre() {
	if !headless {
		Gl.Clear(Gl.COLOR_BUFFER_BIT)
	}
}

func (w *World) post() {}

func (w *World) update(dt float32) {
	w.pre()

	var unp *UnpauseComponent

	complChan := make(chan struct{})
	for _, system := range w.Systems() {
		if headless && system.SkipOnHeadless() {
			continue // so skip it
		}

		system.Pre()

		entities := system.Entities()
		count := len(entities)

		// Concurrency performance maximized at 20+ entities
		// Performance tuning should be conducted for entity updates
		if w.serial || count < 20 || !system.RunInParallel() {
			for _, entity := range entities {
				if w.paused && !entity.Component(&unp) {
					continue // with other entities
				}
				system.Update(entity, dt)
			}
		} else {
			for _, entity := range entities {
				if w.paused && !entity.Component(&unp) {
					count--
					continue // with other entities
				}
				go func(entity *Entity) {
					system.Update(entity, dt)
					complChan <- struct{}{}
				}(entity)
			}
			for ; count > 0; count-- {
				<-complChan
			}
		}
		system.Post()
	}
	close(complChan)
	w.post()
}
