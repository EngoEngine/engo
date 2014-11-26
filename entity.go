package engi

import (
	"reflect"
)

type Entity struct {
	id         string
	components map[reflect.Type]Component
	requires   map[string]bool
	Exists     bool
}

func NewEntity(requires []string) *Entity {
	e := &Entity{requires: make(map[string]bool), components: make(map[reflect.Type]Component)}
	for _, req := range requires {
		e.requires[req] = true
	}
	e.Exists = true
	return e
}

func (e *Entity) DoesRequire(name string) bool { return e.requires[name] }

func (e *Entity) AddComponent(component Component) {
	e.components[reflect.TypeOf(component)] = component
}

// GetComponent takes a double pointer to a Component,
// and populates it with the value of the right type.
func (e *Entity) GetComponent(x interface{}) bool {
	v := reflect.ValueOf(x).Elem() // *T
	c, ok := e.components[v.Type()]
	if !ok {
		return false
	}
	v.Set(reflect.ValueOf(c))
	return true
}

func (e *Entity) ID() string {
	return e.id
}
