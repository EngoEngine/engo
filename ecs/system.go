package ecs

// System is an interface which implements an ECS-System. A System
// should iterate over its Entities on `Update`, in any way suitable
// for the current implementation.
type System interface {
	// Type returns a unique string identifier, usually the struct name
	// eg. "RenderSystem", "CollisionSystem"...
	Type() string
	// Priority is used to create the order in which Systems (in the World) are processed
	Priority() int

	// New is the initialisation of the System
	New(*World)
	// Update is ran every frame, with `dt` being the time in seconds since the last frame
	Update(dt float32)

	// AddEntity adds a new Entity to the System
	AddEntity(entity *Entity)
	// RemoveEntity removes an Entity from the System
	RemoveEntity(entity *Entity)
}

// Systems implements a sortable list of `System`. It is indexed on `System.Priority()`.
type Systems []System

func (s Systems) Len() int {
	return len(s)
}

func (s Systems) Less(i, j int) bool {
	return s[i].Priority() > s[j].Priority()
}

func (s Systems) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
