package engi

type Entity struct {
	id         string
	components []Component
	requires   map[string]bool
}

func NewEntity(requires []string) *Entity {
	e := &Entity{requires: make(map[string]bool)}
	for _, req := range requires {
		e.requires[req] = true
	}
	return e
}

func (e *Entity) DoesRequire(name string) bool { return e.requires[name] }

func (e *Entity) AddComponent(component Component) {
	e.components = append(e.components, component)
}

func (e *Entity) GetComponent(name string) Component {
	for _, component := range e.components {
		if component.Name() == name {
			return component
		}
	}
	return nil
}

func (e *Entity) ID() string {
	return e.id
}
