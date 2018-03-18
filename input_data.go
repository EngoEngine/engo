package engo

// Action corresponds to a control action such as move, press, release
type Action int

// Key correspends to a keyboard key
type Key int

// Modifier represents a special key pressed along with another key
type Modifier int

var (
	// Move is an action representing mouse movement
	Move = Action(0)
	// Press is an action representing a mouse press/click
	Press = Action(1)
	// Release is an action representing a mouse a release
	Release = Action(2)
	// Neutral represents a neutral action
	Neutral = Action(99)
	// Shift represents the shift modifier.
	// It is triggered when the shift key is pressed simultaneously with another key
	Shift = Modifier(0x0001)
	// Control represents the control modifier
	// It is triggered when the ctrl key is pressed simultaneously with another key
	Control = Modifier(0x0002)
	// Alt represents the alt modifier
	// It is triggered when the alt key is pressed simultaneously with another key
	Alt = Modifier(0x0004)
	// Super represents the super modifier
	// (Windows key on Microsoft Windows, Command key on Apple OSX, and varies on Linux)
	// It is triggered when the super key is pressed simultaneously with another key
	Super = Modifier(0x0008)
)

// MouseButton corresponds to a mouse button.
type MouseButton int

// Mouse buttons
const (
	// MouseButtonLeft represents the left mouse button
	MouseButtonLeft MouseButton = 0
	// MouseButtonRight represents the right mouse button
	MouseButtonRight MouseButton = 1
	// MouseButtonMiddle represent the middle mosue button
	MouseButtonMiddle MouseButton = 2
	// MouseButton4 represents the 4th mouse button
	MouseButton4 MouseButton = 3
	// MouseButton5 represents the 5th mouse button
	MouseButton5 MouseButton = 4
	// MouseButton6 represents the 6th mouse button
	MouseButton6 MouseButton = 5
	// MouseButton7 represents the 7th mouse button
	MouseButton7 MouseButton = 6
	// MouseButton4 represents the last mouse button
	MouseButtonLast MouseButton = 7
)

// MouseState represents the current state of the Mouse (or latest Touch-events).
type MouseState struct {
	// X and Y are the coordinates of the Mouse, relative to the `Canvas`.
	X, Y float32
	// ScrollX and ScrollY are the amount of scrolling the user has done with his mouse wheel in the respective directions.
	ScrollX, ScrollY float32
	// Action indicates what the gamer currently has done with his mouse.
	Action Action
	// Button indicates which button is being pressed by the gamer (if any).
	Button MouseButton
	// Modifier indicates which modifier (shift, alt, etc.) has been pressed during the Action.
	Modifier Modifier
}
