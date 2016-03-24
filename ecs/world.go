package ecs

import (
	"log"
	"runtime"
	"sort"
)

// UnmetRequirement is an error that can be raise by ecs.(*World).AddEntity()
// when said entity requires a system this world does not know of.
type UnmetRequirement struct {
	Msg string
}

// Error just implements the error interface for our missing requirement
func (e UnmetRequirement) Error() string {
	return "world cannot add entity with unmet requirement: " + e.Msg
}

// World contains a bunch of Entities, and a bunch of Systems.
// It is the recommended way to run ecs
type World struct {
	entities  map[string]*Entity
	systemMap map[string]System // tracks presence of a system in the world
	systems   Systems

	isSetup bool
	serial  bool
}

// New initialises the World
func (w *World) New() {
	if w.isSetup {
		return
	}

	w.entities = make(map[string]*Entity)
	w.systemMap = make(map[string]System)

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
// In case the entity you are trying to add is requiring an System that this
// world does not know this method will return an UnmetRequirement error
func (w *World) AddEntity(entity *Entity) error {
	w.entities[entity.ID()] = entity
	reqs := entity.Requirements()
	for _, req := range reqs {
		_, ok := w.systemMap[req]
		if !ok {
			// return an error with info about the first missing req
			// WARNING this will not compile a list of all your missing requirements
			// and will just fail on the first missing one.
			return UnmetRequirement{Msg: req}
		}
	}

	for _, system := range w.systems {
		if entity.DoesRequire(system.Type()) {
			system.AddEntity(entity)
		}
	}
	return nil
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
func (w *World) AddSystem(system System) {
	// Special checks for LinearSystem
	if linSys, ok := (system).(linearSystemSetter); ok {
		var pre LinearSystemPre
		var post LinearSystemPost

		if preFunc, ok := (system).(LinearSystemPre); ok {
			pre = preFunc
		}

		if postFunc, ok := (system).(LinearSystemPost); ok {
			post = postFunc
		}

		if update, ok := (system).(LinearSystemUpdate); ok {
			linSys.setLinearSystemFunctions(pre, update, post)
		} else {
			log.Println("Warning: Linear System", system.Type(), "does not implement ecs.LinearSystemUpdate")
		}
	}

	system.New(w)
	w.systems = append(w.systems, system)
	// update system map so that we can quickly test if a system is present
	// in the world
	w.systemMap[system.Type()] = system
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
func (w *World) Systems() []System {
	return w.systems
}

// HasSystem tests if a given system is present in this world
func (w *World) HasSystem(systemType string) bool {
	_, ok := w.systemMap[systemType]
	return ok
}

// Update is called on each frame, with `dt` being the time difference in seconds since the last `Update` call
func (w *World) Update(dt float32) {
	for _, system := range w.Systems() {
		system.Update(dt)
	}
}
