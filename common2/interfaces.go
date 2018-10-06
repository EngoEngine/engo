// interfaces.go is intended to provide a simple means of adding components to
// each system
//
// Getters
//
// These are added functions to each class to allow them to meet the interfaces
// we use with AddByInterface methods on each system.
//
// Faces
//
// The interfaces that end in "Face" are all met by a specific component, which
// can be composed into an Entity. The word Get is used because, otherwise it
// would collide with the name of the object, when stored anonymously in a
// parent entity.
//
// Ables
//
// The interfaces that end in "able" are those required by a specific system,
// and if an an object meets this interface it can be added to that system.
//
// Note: *able* is used not *er* because they don't really do thing anything.
//
// Note: The names have not been contracted for consistency, the interface is
// *Collisionable* not *Collidable*.
//
// Not-Ables
//
// The Not-Ables are interfaces of components used to flag entities to not add to the system,
// for use with the ecs.World.AddSystemInterface

package common2

import "engo.io/ecs"

// GetSpaceComponent Provides container classes ability to fulfil the interface and be accessed more simply by systems, eg in AddByInterface Methods
func (c *SpaceComponent) GetSpaceComponent() *SpaceComponent {
	return c
}

// GetCollisionComponent Provides container classes ability to fulfil the interface and be accessed more simply by systems, eg in AddByInterface Methods
func (c *CollisionComponent) GetCollisionComponent() *CollisionComponent {
	return c
}

// Faces

// BasicFace is the means of accessing the ecs.BasicEntity class , it also has the ID method, to simplify, finding an item within a system
type BasicFace interface {
	ID() uint64
	GetBasicEntity() *ecs.BasicEntity
}

// SpaceFace allows typesafe access to an anonymous SpaceComponent
type SpaceFace interface {
	GetSpaceComponent() *SpaceComponent
}

// CollisionFace allows typesafe access to an anonymous CollisionComponent
type CollisionFace interface {
	GetCollisionComponent() *CollisionComponent
}

// Combined for systems

// Collisionable is the required interface for the CollisionSystem.AddByInterface method
type Collisionable interface {
	BasicFace
	CollisionFace
	SpaceFace
}

// Not-Ables

// NotCollisionComponent is used to flag an entity as not in the CollisionSystem
// even if it has the proper components
type NotCollisionComponent struct{}

// GetNotCollisionComponent implements the NotCollisionable interface
func (n *NotCollisionComponent) GetNotCollisionComponent() *NotCollisionComponent {
	return n
}

// NotCollisionable is an interface used to flag an entity as not in the
// CollisionSystem even if it has the proper components
type NotCollisionable interface {
	GetNotCollisionComponent() *NotCollisionComponent
}
