package common

import (
	"log"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/math"
)

// Cursor is a reference to a GLFW-cursor - to be used with the `SetCursor` method.
type Cursor uint8

const (
	//CursorNone is no visible cursor.
	CursorNone = iota
	//CursorArrow is an arrow cursor.
	CursorArrow
	//CursorCrosshair is a crosshair cursor.
	CursorCrosshair
	//CursorHand is a hand cursor.
	CursorHand
	//CursorIBeam is a text-editor I cursor.
	CursorIBeam
	//CursorHResize is the horizontal resize window cursor.
	CursorHResize
	//CursorVResize is the vertical resize window cursor.
	CursorVResize
)

// MouseSystemPriority is the priority of the MouseSystem
const MouseSystemPriority = 100

// Mouse is the representation of the physical mouse
type Mouse struct {
	// X is the current x position of the mouse in the game
	X float32
	// Y is the current y position of the mouse in the game
	Y float32
	// ScrollX is the current scrolled position on the x component
	ScrollX float32
	// ScrollY is the current scrolled position on the y component
	ScrollY float32
	// Action is the currently active Action
	Action engo.Action
	// Button is which button is being pressed on the mouse
	Button engo.MouseButton
	// Modifier is whether any modifier mouse buttons are being pressed
	Modifer engo.Modifier
}

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
	// Dragged is true whenever the entity space was left-clicked,
	// and then the mouse started moving (while holding)
	Dragged bool
	// RightClicked is true whenever the entity space was right-clicked
	// in this frame
	RightClicked bool
	// RightDragged is true whenever the entity space was right-clicked,
	// and then the mouse started moving (while holding)
	RightDragged bool
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
	Modifier engo.Modifier

	// startedDragging is used internally to see if *this* is the object that is being dragged
	startedDragging bool
	// startedRightDragging is used internally to see if *this* is the object that is being right-dragged
	rightStartedDragging bool
}

type mouseEntity struct {
	*ecs.BasicEntity
	*MouseComponent
	*SpaceComponent
	*RenderComponent
}

// MouseSystem listens for mouse events, and changes value for MouseComponent accordingly
type MouseSystem struct {
	entities []mouseEntity
	world    *ecs.World
	camera   *CameraSystem

	mouseX    float32
	mouseY    float32
	mouseDown bool
}

// Priority returns a priority higher than most, to ensure that this System runs before all others
func (m *MouseSystem) Priority() int { return MouseSystemPriority }

// New initializes the MouseSystem. It is run before any updates.
func (m *MouseSystem) New(w *ecs.World) {
	m.world = w

	// First check if the CameraSystem is available
	for _, system := range m.world.Systems() {
		switch sys := system.(type) {
		case *CameraSystem:
			m.camera = sys
		}
	}

	if m.camera == nil {
		log.Println("ERROR: CameraSystem not found - have you added the `RenderSystem` before the `MouseSystem`?")
		return
	}
}

// Add adds a new entity to the MouseSystem.
// * RenderComponent is only required if you're using the HUDShader on this Entity.
// * SpaceComponent is required whenever you want to know specific mouse-events on this Entity (like hover,
//   click, etc.). If you don't need those, then you can omit the SpaceComponent.
// * MouseComponent is always required.
// * BasicEntity is always required.
func (m *MouseSystem) Add(basic *ecs.BasicEntity, mouse *MouseComponent, space *SpaceComponent, render *RenderComponent) {
	m.entities = append(m.entities, mouseEntity{basic, mouse, space, render})
}

// AddByInterface adds the Entity to the system as long as it satisfies, Mouseable.  Any Entity containing a BasicEntity,MouseComponent, and RenderComponent, automatically does this.
func (m *MouseSystem) AddByInterface(i ecs.Identifier) {
	o, _ := i.(Mouseable)
	m.Add(o.GetBasicEntity(), o.GetMouseComponent(), o.GetSpaceComponent(), o.GetRenderComponent())
}

// Remove removes an entity from the MouseSystem.
func (m *MouseSystem) Remove(basic ecs.BasicEntity) {
	var delete = -1
	for index, entity := range m.entities {
		if entity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		m.entities = append(m.entities[:delete], m.entities[delete+1:]...)
	}
}

// Update updates all the entities in the MouseSystem.
func (m *MouseSystem) Update(dt float32) {
	// Translate Mouse.X and Mouse.Y into "game coordinates"
	switch engo.CurrentBackEnd {
	case engo.BackEndGLFW:
		m.mouseX = engo.Input.Mouse.X*m.camera.Z() + (m.camera.X()-(engo.GameWidth()/2)*m.camera.Z())/engo.GetGlobalScale().X
		m.mouseY = engo.Input.Mouse.Y*m.camera.Z() + (m.camera.Y()-(engo.GameHeight()/2)*m.camera.Z())/engo.GetGlobalScale().Y
	case engo.BackEndMobile:
		m.mouseX = engo.Input.Mouse.X*m.camera.Z() + (m.camera.X()-(engo.GameWidth()/2)*m.camera.Z()+(engo.ResizeXOffset/2))/engo.GetGlobalScale().X
		m.mouseY = engo.Input.Mouse.Y*m.camera.Z() + (m.camera.Y()-(engo.GameHeight()/2)*m.camera.Z()+(engo.ResizeYOffset/2))/engo.GetGlobalScale().Y
	case engo.BackEndWeb:
		m.mouseX = engo.Input.Mouse.X*m.camera.Z() + (m.camera.X()-(engo.GameWidth()/2)*m.camera.Z()+(engo.ResizeXOffset/2))/engo.GetGlobalScale().X
		m.mouseY = engo.Input.Mouse.Y*m.camera.Z() + (m.camera.Y()-(engo.GameHeight()/2)*m.camera.Z()+(engo.ResizeYOffset/2))/engo.GetGlobalScale().X
	}

	// Rotate if needed
	if m.camera.angle != 0 {
		sin, cos := math.Sincos(m.camera.angle * math.Pi / 180)
		m.mouseX, m.mouseY = m.mouseX*cos+m.mouseY*sin, m.mouseY*cos-m.mouseX*sin
	}

	for _, e := range m.entities {
		// Reset all values except these
		*e.MouseComponent = MouseComponent{
			Track:                e.MouseComponent.Track,
			Hovered:              e.MouseComponent.Hovered,
			startedDragging:      e.MouseComponent.startedDragging,
			rightStartedDragging: e.MouseComponent.rightStartedDragging,
		}

		if e.MouseComponent.Track {
			// track mouse position so that systems that need to stay on the mouse
			// position can do it (think an RTS when placing a new building and
			// you get a ghost building following your mouse until you click to
			// place it somewhere in your world.
			e.MouseComponent.MouseX = m.mouseX
			e.MouseComponent.MouseY = m.mouseY
		}

		mx := m.mouseX
		my := m.mouseY

		if e.SpaceComponent == nil {
			continue // with other entities
		}

		if e.RenderComponent != nil {
			// Hardcoded special case for the HUD | TODO: make generic instead of hardcoding
			if e.RenderComponent.shader == HUDShader || e.RenderComponent.shader == LegacyHUDShader {
				mx = engo.Input.Mouse.X
				my = engo.Input.Mouse.Y
			}

			if e.RenderComponent.Hidden {
				continue // skip hidden components
			}
		}

		// If the Mouse component is a tracker we always update it
		// Check if the X-value is within range
		// and if the Y-value is within range
		if e.MouseComponent.Track || e.MouseComponent.startedDragging ||
			e.SpaceComponent.Contains(engo.Point{X: mx, Y: my}) {

			e.MouseComponent.Enter = !e.MouseComponent.Hovered
			e.MouseComponent.Hovered = true
			e.MouseComponent.Released = false

			if !e.MouseComponent.Track {
				// If we're tracking, we've already set these
				e.MouseComponent.MouseX = mx
				e.MouseComponent.MouseY = my
			}

			switch engo.Input.Mouse.Action {
			case engo.Press:
				switch engo.Input.Mouse.Button {
				case engo.MouseButtonLeft:
					e.MouseComponent.Clicked = true
					e.MouseComponent.startedDragging = true
				case engo.MouseButtonRight:
					e.MouseComponent.RightClicked = true
					e.MouseComponent.rightStartedDragging = true
				}

				m.mouseDown = true
			case engo.Release:
				switch engo.Input.Mouse.Button {
				case engo.MouseButtonLeft:
					e.MouseComponent.Released = true
				case engo.MouseButtonRight:
					e.MouseComponent.RightReleased = true
				}
			case engo.Move:
				if m.mouseDown && e.MouseComponent.startedDragging {
					e.MouseComponent.Dragged = true
				}
				if m.mouseDown && e.MouseComponent.rightStartedDragging {
					e.MouseComponent.RightDragged = true
				}
			}
		} else {
			if e.MouseComponent.Hovered {
				e.MouseComponent.Leave = true
			}

			e.MouseComponent.Hovered = false
		}

		if engo.Input.Mouse.Action == engo.Release {
			// dragging stops as soon as one of the currently pressed buttons
			// is released
			e.MouseComponent.Dragged = false
			e.MouseComponent.startedDragging = false
			// TODO maybe separate out the release into left-button release and right-button release
			e.MouseComponent.rightStartedDragging = false
			// mouseDown goes false as soon as one of the pressed buttons is
			// released. Effectively ending any dragging
			m.mouseDown = false
		}

		// propagate the modifiers to the mouse component so that game
		// implementers can take different decisions based on those
		e.MouseComponent.Modifier = engo.Input.Mouse.Modifer
	}
}
