package engo

import (
	"log"

	"engo.io/ecs"
)

type AnimationAction struct {
	Name   string
	Frames []int
}

// Component that controls animation in rendering entities
type AnimationComponent struct {
	index            int              // What frame in the is being used
	Rate             float32          // How often frames should increment, in seconds.
	change           float32          // The time since the last incrementation
	Drawables        []Drawable       // Renderables
	Animations       map[string][]int // All possible animations
	CurrentAnimation []int            // The current animation
}

func NewAnimationComponent(drawables []Drawable, rate float32) AnimationComponent {
	return AnimationComponent{
		Animations: make(map[string][]int),
		Drawables:  drawables,
		Rate:       rate,
	}
}

func (ac *AnimationComponent) SelectAnimationByName(name string) {
	ac.CurrentAnimation = ac.Animations[name]
}

func (ac *AnimationComponent) SelectAnimationByAction(action *AnimationAction) {
	ac.CurrentAnimation = ac.Animations[action.Name]
}

func (ac *AnimationComponent) AddAnimationAction(action *AnimationAction) {
	ac.Animations[action.Name] = action.Frames
}

func (ac *AnimationComponent) AddAnimationActions(actions []*AnimationAction) {
	for _, action := range actions {
		ac.Animations[action.Name] = action.Frames
	}
}

func (ac *AnimationComponent) Cell() Drawable {
	idx := ac.CurrentAnimation[ac.index]

	return ac.Drawables[idx]
}

func (ac *AnimationComponent) NextFrame() {
	if len(ac.CurrentAnimation) == 0 {
		log.Println("No data for this animation")
		return
	}

	ac.index += 1
	if ac.index >= len(ac.CurrentAnimation) {
		ac.index = 0
	}
	ac.change = 0
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
		e.AnimationComponent.change += dt
		if e.AnimationComponent.change >= e.AnimationComponent.Rate {
			e.AnimationComponent.NextFrame()
			e.RenderComponent.SetDrawable(e.AnimationComponent.Cell())
		}
	}
}
