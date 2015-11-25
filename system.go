// Copyright 2014 Harrison Shoebridge. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package engi

type Component interface {
	Type() string
}

type Systemer interface {
	Update(entity *Entity, dt float32)
	Type() string
	Priority() int
	Pre()
	Post()
	New()
	Entities() []*Entity
	AddEntity(entity *Entity)
	RemoveEntity(entity *Entity)
	SkipOnHeadless() bool
	SetWorld(*World)
	RunInParallel() bool
}

type System struct {
	entities             map[string]*Entity
	messageQueue         []Message
	ShouldSkipOnHeadless bool
	World                *World
}

func NewSystem() *System {
	s := &System{}
	s.entities = make(map[string]*Entity)
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
	list := make([]*Entity, len(s.entities))
	i := 0
	for _, ent := range s.entities {
		list[i] = ent
		i++
	}
	return list
}

func (s *System) AddEntity(entity *Entity) {
	s.entities[entity.ID()] = entity
}

func (s *System) RemoveEntity(entity *Entity) {
	delete(s.entities, entity.ID())
}

func (s System) SkipOnHeadless() bool {
	return s.ShouldSkipOnHeadless
}

func (s *System) SetWorld(w *World) {
	s.World = w
}
