package engi

type Entity struct {
	id         string
	components []Component
	requires   []string
}

func NewEntity(requires []string) *Entity {
	return &Entity{requires: requires}
}

func (e *Entity) DoesRequire(name string) bool {
	for _, requirement := range e.requires {
		if requirement == name {
			return true
		}
	}

	return false
}

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
