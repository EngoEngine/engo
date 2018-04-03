package common

import (
	"log"

	"engo.io/ecs"
)

// Animation represents properties of an animation.
type Animation struct {
	Name   string
	Frames []int
	Loop   bool
}

// AnimationComponent tracks animations of an entity it is part of.
// This component should be created using NewAnimationComponent.
type AnimationComponent struct {
	Drawables        []Drawable            // Renderables
	Animations       map[string]*Animation // All possible animations
	CurrentAnimation *Animation            // The current animation
	Rate             float32               // How often frames should increment, in seconds.
	index            int                   // What frame in the is being used
	change           float32               // The time since the last incrementation
	def              *Animation            // The default animation to play when nothing else is playing
}

// NewAnimationComponent creates an AnimationComponent containing all given
// drawables. Animations will be played using the given rate.
func NewAnimationComponent(drawables []Drawable, rate float32) AnimationComponent {
	return AnimationComponent{
		Animations: make(map[string]*Animation),
		Drawables:  drawables,
		Rate:       rate,
	}
}

// SelectAnimationByName sets the current animation. The name must be
// registered.
func (ac *AnimationComponent) SelectAnimationByName(name string) {
	ac.CurrentAnimation = ac.Animations[name]
	ac.index = 0
}

// SelectAnimationByAction sets the current animation.
// An nil action value selects the default animation.
func (ac *AnimationComponent) SelectAnimationByAction(action *Animation) {
	ac.CurrentAnimation = action
	ac.index = 0
}

// AddDefaultAnimation adds an animation which is used when no other animation is playing.
func (ac *AnimationComponent) AddDefaultAnimation(action *Animation) {
	ac.AddAnimation(action)
	ac.def = action
}

// AddAnimation registers an animation under its name, making it available
// through SelectAnimationByName.
func (ac *AnimationComponent) AddAnimation(action *Animation) {
	ac.Animations[action.Name] = action
}

// AddAnimations registers all given animations.
func (ac *AnimationComponent) AddAnimations(actions []*Animation) {
	for _, action := range actions {
		ac.AddAnimation(action)
	}
}

// Cell returns the drawable for the current frame.
func (ac *AnimationComponent) Cell() Drawable {
	idx := ac.CurrentAnimation.Frames[ac.index]

	return ac.Drawables[idx]
}

// NextFrame advances the current animation by one frame.
func (ac *AnimationComponent) NextFrame() {
	if len(ac.CurrentAnimation.Frames) == 0 {
		log.Println("No data for this animation")
		return
	}

	ac.index++
	ac.change = 0
	if ac.index >= len(ac.CurrentAnimation.Frames) {
		ac.index = 0

		if !ac.CurrentAnimation.Loop {
			ac.CurrentAnimation = nil
			return
		}
	}
}

// AnimationSystem tracks AnimationComponents, advancing their current animation.
type AnimationSystem struct {
	entities map[ecs.BasicEntity]animationEntity
}

type animationEntity struct {
	*AnimationComponent
	*RenderComponent
}

// Add starts tracking the given entity.
func (a *AnimationSystem) Add(basic *ecs.BasicEntity, anim *AnimationComponent, render *RenderComponent) {
	if a.entities == nil {
		a.entities = make(map[ecs.BasicEntity]animationEntity)
	}
	a.entities[*basic] = animationEntity{anim, render}
}

// AddByInterface Allows an Entity to be added directly using the Animtionable interface. which every entity containing the BasicEntity,AnimationComponent,and RenderComponent anonymously, automatically satisfies.
func (a *AnimationSystem) AddByInterface(i ecs.Identifier) {
	o, _ := i.(Animationable)
	a.Add(o.GetBasicEntity(), o.GetAnimationComponent(), o.GetRenderComponent())
}

// Remove stops tracking the given entity.
func (a *AnimationSystem) Remove(basic ecs.BasicEntity) {
	if a.entities != nil {
		delete(a.entities, basic)
	}
}

// Update advances the animations of all tracked entities.
func (a *AnimationSystem) Update(dt float32) {
	for _, e := range a.entities {
		if e.AnimationComponent.CurrentAnimation == nil {
			if e.AnimationComponent.def == nil {
				return
			}
			e.AnimationComponent.SelectAnimationByAction(e.AnimationComponent.def)
		}

		e.AnimationComponent.change += dt
		if e.AnimationComponent.change >= e.AnimationComponent.Rate {
			e.RenderComponent.Drawable = e.AnimationComponent.Cell()
			e.AnimationComponent.NextFrame()
		}
	}
}
