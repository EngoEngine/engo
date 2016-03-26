package engi

import (
	"log"

	"github.com/paked/engi/ecs"
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

// Creates a new animation component
func NewAnimationComponent(drawables []Drawable, rate float32) *AnimationComponent {
	return &AnimationComponent{
		Animations: make(map[string][]int),
		Drawables:  drawables,
		Rate:       rate,
	}
}

// Sets the current animation by name
func (ac *AnimationComponent) SelectAnimationByName(name string) {
	ac.CurrentAnimation = ac.Animations[name]
}

// Sets the current animation by action
func (ac *AnimationComponent) SelectAnimationByAction(action *AnimationAction) {
	ac.CurrentAnimation = ac.Animations[action.Name]
}

// Adds an action to the component
func (ac *AnimationComponent) AddAnimationAction(action *AnimationAction) {
	ac.Animations[action.Name] = action.Frames
}

// Add multiplate actions to the component
func (ac *AnimationComponent) AddAnimationActions(actions []*AnimationAction) {
	for _, action := range actions {
		ac.Animations[action.Name] = action.Frames
	}
}

func (ac *AnimationComponent) Cell() Drawable {
	idx := ac.CurrentAnimation[ac.index]

	return ac.Drawables[idx]
}

// Switches to the next frame in the current animation
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

// Returns the type of the component (Animation)
func (*AnimationComponent) Type() string {
	return "AnimationComponent"
}

type AnimationSystem struct {
	*ecs.System
}

// Creates a new Animation System
func (a *AnimationSystem) New(*ecs.World) {
	a.System = ecs.NewSystem()
}

// Returns the type of the system (Animation)
func (AnimationSystem) Type() string {
	return "AnimationSystem"
}

// Animation system updated, called according to the FPS
func (a *AnimationSystem) Update(e *ecs.Entity, dt float32) {
	var (
		ac *AnimationComponent
		r  *RenderComponent
		ok bool
	)

	if ac, ok = e.ComponentFast(ac).(*AnimationComponent); !ok {
		return
	}
	if r, ok = e.ComponentFast(r).(*RenderComponent); !ok {
		return
	}

	ac.change += dt
	if ac.change >= ac.Rate {
		ac.NextFrame()
		r.SetDrawable(ac.Cell())
	}
}
