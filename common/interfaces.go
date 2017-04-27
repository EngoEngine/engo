// interfaces.go is intended to provide a simple means of adding components to each system
// Getters
// These are added functions to each class to allow them to meet the interfaces we use with AddByInterface methods on each system
// Faces
// The interfaces that end in "Face" are all met by a specific component, which can be composed into an Entity
// The word Get is used because, otherwise it would collide with the name of the object, when stored anonymously in a parent entity
// Ables
// The interfaces that end in "able" are those required by a specific system, and if an an object meets this interface it can be added to that system
// Note: *able* is used not *er* because they don't really do thing anything
// Note: The names have not been contracted for consistency, the interface is *Collisionable* not *Collidable*
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

// BasicFace is the means of accessing the ecs.BasicEntity class , it also has the ID method, to simplfy, finding an item within a system
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
	SpaceFace
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
