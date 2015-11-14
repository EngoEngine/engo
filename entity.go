package engi

import (
	"reflect"
)

type Entity struct {
	id         string
	components map[string]Component
	requires   map[string]bool
	Pattern    string
}

func NewEntity(requires []string) *Entity {
	e := &Entity{
		id:         generateUUID(),
		requires:   make(map[string]bool),
		components: make(map[string]Component),
	}
	for _, req := range requires {
		e.requires[req] = true
	}
	return e
}

func (e *Entity) DoesRequire(name string) bool {
	return e.requires[name]
}

func (e *Entity) AddComponent(component Component) {
	e.components[reflect.TypeOf(component).String()] = component
}

func (e *Entity) RemoveComponent(component Component) {
	delete(e.components, reflect.TypeOf(component).String())
}

// GetComponent takes a double pointer to a Component,
// and populates it with the value of the right type.
func (e *Entity) Component(x interface{}) bool {
	v := reflect.ValueOf(x).Elem() // *T
	c, ok := e.components[v.Type().String()]
	if !ok {
		return false
	}
	v.Set(reflect.ValueOf(c))
	return true
}

func (e *Entity) ID() string {
	return e.id
}
