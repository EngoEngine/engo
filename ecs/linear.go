package ecs

import (
	"log"
	"reflect"
	"runtime"
)

// LinearSystem is the default implementation of the System interface, which handles Entities in a linear fashion
// Implement `LinearSystemFunctions` and inherit this, in order to use it
type LinearSystem struct {
	// Entities holds the Entity-references as given by the World
	Entities []*Entity
	// RunInParallel indicates whether or not the UpdateEntity function should be called in parallel
	RunInParallel bool
	// Prio allows the `World` to order `System`s: the lower the priority-value, the
	// sooner it will be processed by the `World`.
	Prio int

	functions LinearSystemFunctions
}

type LinearSystemFunctions interface {
	Pre()
	UpdateEntity(entity *Entity, dt float32)
	Post()
}

// Some functions that should be overriden:
func (s *LinearSystem) New(*World) {}
func (LinearSystem) Type() string  { return "generic LinearSystem" }

// Update is called by the `World`
func (s *LinearSystem) Update(dt float32) {
	s.functions.Pre()

	count := len(s.Entities)

	// Calling them serial / in parallel, depending on the settings
	if processSystemInSerial || !s.RunInParallel {
		for _, entity := range s.Entities {
			s.functions.UpdateEntity(entity, dt)
		}
	} else {
		complChan := make(chan struct{})
		for _, entity := range s.Entities {
			go func(entity *Entity) {
				s.functions.UpdateEntity(entity, dt)
				complChan <- struct{}{}
			}(entity)
		}
		for ; count > 0; count-- {
			<-complChan
		}
		close(complChan)
	}
	s.functions.Post()
}

func (s *LinearSystem) AddEntity(entity *Entity) {
	s.Entities = append(s.Entities, entity)
}

// TODO: not sure if one can still override this?
func (s *LinearSystem) RemoveEntity(entity *Entity) {
	for index, e := range s.Entities {
		if e.ID() == entity.ID() {
			// Found it, now let's remove it - TODO: make sure this works for the edge case (index+1 being out-of-range?)
			s.Entities = append(s.Entities[:index], s.Entities[index+1:]...)
		}
	}
}

func (s *LinearSystem) Priority() int {
	return s.Prio
}

func (s *LinearSystem) setLinearSystemFunctions(funcs LinearSystemFunctions) {
	if s != nil {
		s.functions = funcs
	} else {
		log.Println("Warning:", reflect.TypeOf(funcs).String(), "has a pointer-reference to LinearSystem (*ecs.LinearSystem). ")
	}
}

type basicSystemSetter interface {
	setLinearSystemFunctions(LinearSystemFunctions)
}

var (
	processSystemInSerial bool
)

func init() {
	// Run in serial if there's only 1 core
	processSystemInSerial = runtime.NumCPU() == 1
}
