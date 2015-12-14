package ecs

import (
	"runtime"
	"sort"
)

// World contains a bunch of Entitys, and a bunch of Systems. It is
// the recommended way to run ecs
type World struct {
	entities map[string]*Entity
	systems  Systemers

	isSetup bool
	serial  bool
}

// New initialises the World
func (w *World) New() {
	if w.isSetup {
		return
	}

	w.entities = make(map[string]*Entity)

	/*
		// Default WorldBounds values
		WorldBounds.Max = Point{Width(), Height()}
	*/

	// Short-circuit bypass if there's only 1 core
	if runtime.NumCPU() == 1 {
		w.serial = true
	} else {
		w.serial = false
	}

	w.isSetup = true
}

// AddEntity adds a new Entity to the World, and its required Systems
func (w *World) AddEntity(entity *Entity) {
	w.entities[entity.ID()] = entity

	for _, system := range w.systems {
		if entity.DoesRequire(system.Type()) {
			system.AddEntity(entity)
		}
	}
}

// RemoveEntity removes an Entity from the World and its required Systems
func (w *World) RemoveEntity(entity *Entity) {
	for _, system := range w.systems {
		if entity.DoesRequire(system.Type()) {
			system.RemoveEntity(entity)
		}
	}

	delete(w.entities, entity.ID())
}

// AddSystem adds a new System to the World, and then sorts them based on Priority
func (w *World) AddSystem(system Systemer) {
	system.New(w)
	w.systems = append(w.systems, system)
	sort.Sort(w.systems)
}

// Entities returns the list of Entities
func (w *World) Entities() []*Entity {
	entities := make([]*Entity, len(w.entities))
	i := 0
	for _, v := range w.entities {
		entities[i] = v
		i++
	}

	return entities
}

// Systems returns a list of Systems
func (w *World) Systems() []Systemer {
	return w.systems
}

func (w *World) HasSystem(systemType string) bool {
	for _, s := range w.systems {
		if s.Type() == systemType {
			return true
		}
	}
	return false
}

// Update is called on each frame, with dt being the time difference in seconds since the last Update call
func (w *World) Update(dt float32) {
	complChan := make(chan struct{})
	for _, system := range w.Systems() {
		system.Pre()

		entities := system.Entities()
		count := len(entities)

		// Calling them serial / in parallel, depending on the settings
		if w.serial || !system.RunInParallel() {
			for _, entity := range entities {
				system.Update(entity, dt)
			}
		} else {
			for _, entity := range entities {
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
}
