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

package common

import "engo.io/ecs"

// Getters

// GetAnimationComponent Provides container classes ability to fulfil the interface and be accessed more simply by systems, eg in AddByInterface Methods
func (c *AnimationComponent) GetAnimationComponent() *AnimationComponent {
	return c
}

// GetMouseComponent Provides container classes ability to fulfil the interface and be accessed more simply by systems, eg in AddByInterface Methods
func (c *MouseComponent) GetMouseComponent() *MouseComponent {
	return c
}

// GetAudioComponent Provides container classes ability to fulfil the interface and be accessed more simply by systems, eg in AddByInterface Methods
func (c *AudioComponent) GetAudioComponent() *AudioComponent {
	return c
}

// GetRenderComponent Provides container classes ability to fulfil the interface and be accessed more simply by systems, eg in AddByInterface Methods
func (c *RenderComponent) GetRenderComponent() *RenderComponent {
	return c
}

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

// AnimationFace allows typesafe Access to an Annonymous child AnimationComponent
type AnimationFace interface {
	GetAnimationComponent() *AnimationComponent
}

// MouseFace allows typesafe access to an Anonymous child MouseComponent
type MouseFace interface {
	GetMouseComponent() *MouseComponent
}

// AudioFace allows typesafe access to an anonymouse child AudioComponent
type AudioFace interface {
	GetAudioComponent() *AudioComponent
}

// RenderFace allows typesafe access to an anonymous RenderComponent
type RenderFace interface {
	GetRenderComponent() *RenderComponent
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

// Animationable is the required interface for AnimationSystem.AddByInterface method
type Animationable interface {
	BasicFace
	AnimationFace
	RenderFace
}

// Mouseable is the required interface for the MouseSystem AddByInterface method
type Mouseable interface {
	BasicFace
	MouseFace
	SpaceFace
	RenderFace
}

// Audioable is the required interface for the AudioSystem.AddByInterface method
type Audioable interface {
	BasicFace
	AudioFace
}

// Renderable is the required interface for the RenderSystem.AddByInterface method
type Renderable interface {
	BasicFace
	RenderFace
	SpaceFace
}

// Collisionable is the required interface for the CollisionSystem.AddByInterface method
type Collisionable interface {
	BasicFace
	CollisionFace
	SpaceFace
}

// Not-Ables

// NotAnimationComponent is used to flag an entity as not in the AnimationSystem
// even if it has the proper components
type NotAnimationComponent struct{}

// GetNotAnimationComponent implements the NotAnimationable interface
func (n *NotAnimationComponent) GetNotAnimationComponent() *NotAnimationComponent {
	return n
}

// NotAnimationable is an interface used to flag an entity as not in the
// AnimationSystem even if it has the proper components
type NotAnimationable interface {
	GetNotAnimationComponent() *NotAnimationComponent
}

// NotMouseComponent is used to flag an entity as not in the MouseSystem even if
// it has the proper components
type NotMouseComponent struct{}

// GetNotMouseComponent implements the NotMouseable interface
func (n *NotMouseComponent) GetNotMouseComponent() *NotMouseComponent {
	return n
}

// NotMouseable is an interface used to flag an entity as not in the MouseSystem
// even if it has the proper components
type NotMouseable interface {
	GetNotMouseComponent() *NotMouseComponent
}

// NotAudioComponent is used to flag an entity as not in the AudioSystem even if
// it has the proper components
type NotAudioComponent struct{}

// GetNotAudioComponent implements the NotAudioable interface
func (n *NotAudioComponent) GetNotAudioComponent() *NotAudioComponent {
	return n
}

// NotAudioable is an interface used to flag an entity as not in the AudioSystem
// even if it has the proper components
type NotAudioable interface {
	GetNotAudioComponent() *NotAudioComponent
}

// NotRenderComponent is used to flag an entity as not in the RenderSystem even
// if it has the proper components
type NotRenderComponent struct{}

// GetNotRenderComponent implements the NotRenderable interface
func (n *NotRenderComponent) GetNotRenderComponent() *NotRenderComponent {
	return n
}

// NotRenderable is an interface used to flag an entity as not in the
// Rendersystem even if it has the proper components
type NotRenderable interface {
	GetNotRenderComponent() *NotRenderComponent
}

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
