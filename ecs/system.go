// Copyright 2014-2015 Harrison Shoebridge and other contributors.
// All rights reserved. Use of this source code is governed by a
// BSD-style license that can be found in the LICENSE file.
package ecs

// Component is a piece of data which belongs to an Entity
type Component interface {
	Type() string
}

// Systemer is an interface which implements a System. A System
// iterates over the Entitys it is required by, and can process
// their Components
type Systemer interface {
	// Type returns an individual string identifier, usually the struct name
	// eg. "RenderSystem", "CollisionSystem"...
	Type() string
	// Priority is used to create the order in which Systemers are processed
	Priority() int
	// RunInParallel checks whether the System can run in parallel. This is ran every update,
	// so it is possible to run serial when the system has < 10 entities, and then run in parallel
	// when not
	RunInParallel() bool

	// New is the initialisation of the System
	New(*World)
	// Pre is ran just before the System updates, every frame
	Pre()
	// Post is ran just after the System updates, every frame
	Post()
	// Update is ran every frame, and for each Entity which belongs to the
	// System
	Update(entity *Entity, dt float32)

	// Entities returns a slice of all Entities
	Entities() []*Entity
	// AddEntity adds a new Entity to the System
	AddEntity(entity *Entity)
	// AddEntity removes an Entity from the System
	RemoveEntity(entity *Entity)
}

// System is the default implementation of the Systemer interface.
type System struct {
	EntityMap            map[string]*Entity
	ShouldSkipOnHeadless bool
}

// NewSystem returns a new default System
func NewSystem() *System {
	s := &System{}
	s.EntityMap = make(map[string]*Entity)
	return s
}

func (s System) New()  {}
func (s System) Pre()  {}
func (s System) Post() {}

func (s System) Priority() int {
	return 0
}
func (s System) RunInParallel() bool { return false }

func (s System) Entities() []*Entity {
	list := make([]*Entity, len(s.EntityMap))
	i := 0
	for _, ent := range s.EntityMap {
		list[i] = ent
		i++
	}
	return list
}

func (s *System) AddEntity(entity *Entity) {
	s.EntityMap[entity.ID()] = entity
}

func (s *System) RemoveEntity(entity *Entity) {
	delete(s.EntityMap, entity.ID())
}

// Systemers implements a sortable list of System. It is indexed on System.Priority().
type Systemers []Systemer

func (s Systemers) Len() int {
	return len(s)
}

func (s Systemers) Less(i, j int) bool {
	return s[i].Priority() < s[j].Priority()
}

func (s Systemers) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
