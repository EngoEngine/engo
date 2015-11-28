// Copyright 2014 Harrison Shoebridge. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ecs

type Component interface {
	Type() string
}

type Systemer interface {
	Type() string
	Priority() int
	RunInParallel() bool

	New(*World)
	Pre()
	Post()
	Update(entity *Entity, dt float32)

	Entities() []*Entity
	AddEntity(entity *Entity)
	RemoveEntity(entity *Entity)
}

type System struct {
	EntityMap            map[string]*Entity
	ShouldSkipOnHeadless bool
}

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
