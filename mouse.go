package engi

import (
	"github.com/paked/engi/ecs"
)

// MouseComponent is the location for the MouseSystem to store its results;
// to be used / viewed by other Systems
type MouseComponent struct {
	// Clicked is true whenever the Mouse was clicked over
	// the entity space in this frame
	Clicked bool
	// Released is true whenever the left mouse button is released over the
	// entity space in this frame
	Released bool
	// Hovered is true whenever the Mouse is hovering
	// the entity space in this frame. This does not necessarily imply that
	// the mouse button was pressed down in your entity space.
	Hovered bool
	// Dragged is true whenever the entity space was clicked,
	// and then the mouse started moving (while holding)
	Dragged bool
	// RightClicked is true whenever the entity space was right-clicked
	// in this frame
	RightClicked bool
	// RightReleased is true whenever the right mouse button is released over
	// the entity space in this frame. This does not necessarily imply that
	// the mouse button was pressed down in your entity space.
	RightReleased bool
	// Enter is true whenever the Mouse entered the entity space in that frame,
	// but wasn't in that space during the previous frame
	Enter bool
	// Leave is true whenever the Mouse was in the space on the previous frame,
	// but now isn't
	Leave bool
	// Position of the mouse at any moment this is generally used
	// in conjunction with Track = true
	MouseX float32
	MouseY float32
	// Set manually this to true and your mouse component will track the mouse
	// and your entity will always be able to receive an updated mouse
	// component even if its space is not under the mouse cursor
	// WARNING: you MUST know why you want to use this because it will
	// have serious performance impacts if you have many entities with
	// a MouseComponent in tracking mode.
	// This is ideally used for a really small number of entities
	// that must really be aware of the mouse details event when the
	// mouse is not hovering them
	Track bool
	// Modifier is used to store the eventual modifiers that were pressed during
	// the same time the different click events occurred
	Modifier Modifier
}

// Type returns the string representation of the MouseComponent type
func (*MouseComponent) Type() string {
	return "MouseComponent"
}

// MouseSystem listens for mouse events, and changes value for MouseComponent accordingly
type MouseSystem struct {
	ecs.LinearSystem

	mouseX    float32
	mouseY    float32
	mouseDown bool
}

// Type returns the string representation of the MouseSystem type
func (*MouseSystem) Type() string { return "MouseSystem" }

// New initializes the MouseSystem
func (m *MouseSystem) New(*ecs.World) {}

// Priority returns a priority of 10 (higher than most) to ensure that this System runs before all others
func (m *MouseSystem) Priority() int {
	return 10
}

// Pre is called before all Update calls, and is used to compute internal values
func (m *MouseSystem) Pre() {
	// Translate Mouse.X and Mouse.Y into "game coordinates"
	m.mouseX = Mouse.X*Cam.z*(gameWidth/windowWidth) + Cam.x - (gameWidth/2)*Cam.z
	m.mouseY = Mouse.Y*Cam.z*(gameHeight/windowHeight) + Cam.y - (gameHeight/2)*Cam.z
}

// Post is called after all Update calls, and is used to compute internal values
// NOTE: we do not reset modifiers here because we always set them to meaningful
// values when a mouse event is propagated to the mouse components
func (m *MouseSystem) Post() {
	// reset mouse.Action value to something meaningless to avoid
	// catching the same "signal" twice
	Mouse.Action = NEUTRAL
}

// EUpdate sets the MouseComponent values for each Entity
func (m *MouseSystem) UpdateEntity(entity *ecs.Entity, dt float32) {
	var (
		mc     *MouseComponent
		space  *SpaceComponent
		render *RenderComponent
		ok     bool
	)

	// We need MouseComponent to save our findings
	if mc, ok = entity.ComponentFast(mc).(*MouseComponent); !ok {
		return
	}

	// We need SpaceComponent for the location
	if space, ok = entity.ComponentFast(space).(*SpaceComponent); !ok {
		return
	}

	// We need RenderComponent for the Priority
	if render, ok = entity.ComponentFast(render).(*RenderComponent); !ok {
		return
	}

	// Reset some values
	mc.Leave = false

	mx := m.mouseX
	my := m.mouseY

	// Special case: HUD
	if render.priority >= HUDGround {
		mx = Mouse.X
		my = Mouse.Y
	}

	// if the Mouse component is a tracker we always update it
	// Check if the X-value is within range
	// and if the Y-value is within range
	if mc.Track || mx > space.Position.X && mx < (space.Position.X+space.Width) &&
		my > space.Position.Y && my < (space.Position.Y+space.Height) {

		mc.Enter = !mc.Hovered
		mc.Hovered = true

		mc.Released = false
		// track mouse position so that systems that need to stay on the mouse
		// position can do it (think an RTS when placing a new building and
		// you get a ghost building following your mouse until you click to
		// place it somewhere in your world.
		mc.MouseX = mx
		mc.MouseY = my

		// propagate the modifiers to the mouse component so that game
		// implementers can take different decisions based on those
		mc.Modifier = Mouse.Modifer

		switch Mouse.Action {
		case PRESS:
			switch Mouse.Button {
			case MouseButtonLeft:
				mc.Clicked = true
			case MouseButtonRight:
				mc.RightClicked = true
			}
			m.mouseDown = true
		case RELEASE:
			switch Mouse.Button {
			case MouseButtonLeft:
				mc.Released = true
			case MouseButtonRight:
				mc.RightReleased = true
			}
			// dragging stops as soon as one of the currently pressed buttons
			// is released
			mc.Dragged = false
			// mouseDown goes false as soon as one of the pressed buttons is
			// released. Effectively ending any dragging
			m.mouseDown = false
		case MOVE:
			if m.mouseDown {
				mc.Dragged = true
			}
		}
	} else {
		if mc.Hovered {
			mc.Leave = true
			// propagate the modifiers to the mouse component so that game
			// implementers can take different decisions based on those
			mc.Modifier = Mouse.Modifer
		}
		mc.Hovered = false
	}
}
