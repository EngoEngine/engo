package common

import (
	"log"

	"engo.io/ecs"
)

type Animation struct {
	Name   string
	Frames []int
	Loop   bool
}

// Component that controls animation in rendering entities
type AnimationComponent struct {
	Drawables        []Drawable            // Renderables
	Animations       map[string]*Animation // All possible animations
	CurrentAnimation *Animation            // The current animation
	Rate             float32               // How often frames should increment, in seconds.
	index            int                   // What frame in the is being used
	change           float32               // The time since the last incrementation
	def              string                // The default animation to play when nothing else is playing
}

func NewAnimationComponent(drawables []Drawable, rate float32) AnimationComponent {
	return AnimationComponent{
		Animations: make(map[string]*Animation),
		Drawables:  drawables,
		Rate:       rate,
	}
}

func (ac *AnimationComponent) SelectAnimationByName(name string) {
	ac.CurrentAnimation = ac.Animations[name]
	ac.index = 0
}

func (ac *AnimationComponent) SelectAnimationByAction(action *Animation) {
	ac.SelectAnimationByName(action.Name)
}

func (ac *AnimationComponent) AddDefaultAnimation(action *Animation) {
	ac.def = action.Name

	ac.AddAnimation(action)
}

func (ac *AnimationComponent) AddAnimation(action *Animation) {
	ac.Animations[action.Name] = action
}

func (ac *AnimationComponent) AddAnimations(actions []*Animation) {
	for _, action := range actions {
		ac.AddAnimation(action)
	}
}

func (ac *AnimationComponent) Cell() Drawable {
	idx := ac.CurrentAnimation.Frames[ac.index]

	return ac.Drawables[idx]
}

func (ac *AnimationComponent) NextFrame() {
	if len(ac.CurrentAnimation.Frames) == 0 {
		log.Println("No data for this animation")
		return
	}

	ac.index += 1
	ac.change = 0
	if ac.index >= len(ac.CurrentAnimation.Frames) {
		ac.index = 0

		if !ac.CurrentAnimation.Loop {
			ac.CurrentAnimation = nil
			return
		}
	}
}

type animationEntity struct {
	*ecs.BasicEntity
	*AnimationComponent
	*RenderComponent
}

type AnimationSystem struct {
	entities []animationEntity
}

func (a *AnimationSystem) Add(basic *ecs.BasicEntity, anim *AnimationComponent, render *RenderComponent) {
	a.entities = append(a.entities, animationEntity{basic, anim, render})
}

func (a *AnimationSystem) Remove(basic ecs.BasicEntity) {
	delete := -1
	for index, e := range a.entities {
		if e.BasicEntity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		a.entities = append(a.entities[:delete], a.entities[delete+1:]...)
	}
}

func (a *AnimationSystem) Update(dt float32) {
	for _, e := range a.entities {
		if e.AnimationComponent.CurrentAnimation == nil {
			if e.AnimationComponent.def == "" {
				return
			}

			e.AnimationComponent.SelectAnimationByName(e.AnimationComponent.def)
		}

		e.AnimationComponent.change += dt
		if e.AnimationComponent.change >= e.AnimationComponent.Rate {
			e.RenderComponent.Drawable = e.AnimationComponent.Cell()
			e.AnimationComponent.NextFrame()
		}
	}
}
