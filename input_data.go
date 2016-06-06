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

// those are default values for engo_js defined here because some of them are shared
// with engo_glfw.
// engo_glfw redefines the variables it needs to other values during init() so
var (
	// Grave represents the '`' keyboard key
	Grave Key = 192
	// Dash represents the '-' keyboard key
	Dash Key = 189
	// Apostrophe represents the `'` keyboard key
	Apostrophe Key = 222
	// Semicolon represents the ';' keyboard key
	Semicolon Key = 186
	// Equals reprsents the '=' keyboard key
	Equals Key = 187
	// Comma represents the ',' keyboard key
	Comma Key = 188
	// Period represents the '.' keyboard key
	Period Key = 190
	// Slash represents the '/' keyboard key
	Slash Key = 191
	// Backslash represents the '\' keyboard key
	Backslash Key = 220
	//Backspace represents the backspace keyboard key
	Backspace Key = 8
	// Tab represents the tab keyboard key
	Tab Key = 9
	// CapsLock represents the caps lock keyboard key
	CapsLock Key = 20
	// Space represents the space keyboard key
	Space Key = 32
	// Enter represents the enter keyboard key
	Enter Key = 13
	// Escape represents the escape keyboard key
	Escape Key = 27
	// Insert represents the insert keyboard key
	Insert Key = 45
	// PrintScreen represents the print screen keyboard key often
	// represented by 'Prt Scrn', 'Prt Scn', or 'Print Screen'
	PrintScreen Key = 42
	// Delete represents the delete keyboard key
	Delete Key = 46
	// PageUp represents the page up keyboard key
	PageUp Key = 33
	// PageDown represents the page down keyboard key
	PageDown Key = 34
	// Home represents the home keyboard key
	Home Key = 36
	// End represents the end keyboard key
	End Key = 35
	// Pause represents the pause keyboard key
	Pause Key = 19
	// ScrollLock represents the scroll lock keyboard key
	ScrollLock Key = 145
	// AllowLeft represents the arrow left keyboard key
	ArrowLeft Key = 37
	// ArrowRight represents the arrow right keyboard key
	ArrowRight Key = 39
	// ArrowDown represents the down arrow keyboard key
	ArrowDown Key = 40
	// ArrowUp represents the up arrow keyboard key
	ArrowUp Key = 38
	// LeftBracket represents the '[' keyboard key
	LeftBracket Key = 219
	// LeftShift represents the left shift keyboard key
	LeftShift Key = 16
	// LeftControl represents the left control keyboard key
	LeftControl Key = 17
	// LeftSuper represents the left super keyboard key
	// (Windows key on Microsoft Windows, Command key on Apple OSX, and varies on Linux)
	LeftSuper Key = 73
	// LeftAlt represents the left alt keyboard key
	LeftAlt Key = 18
	// RightBracket represents the ']' keyboard key
	RightBracket Key = 221
	// RightShift represents the right shift keyboard key
	RightShift Key = 16
	// RightControl represents the right control keyboard key
	RightControl Key = 17
	// RightSuper represents the right super keyboard key
	// (Windows key on Microsoft Windows, Command key on Apple OSX, and varies on Linux)
	RightSuper Key = 73
	// RightAlt represents the left alt keyboard key
	RightAlt Key = 18
	// Zero represents the '0' keyboard key
	Zero Key = 48
	// One represents the '1' keyboard key
	One Key = 49
	// Two represents the '2' keyboard key
	Two Key = 50
	// Three represents the '3' keyboard key
	Three Key = 51
	// Four represents the '4' keyboard key
	Four Key = 52
	// Five represents the '5' keyboard key
	Five Key = 53
	// Six represents the '6' keyboard key
	Six Key = 54
	// Seven represents the '7' keyboard key
	Seven Key = 55
	// Eight represents the '8' keyboard key
	Eight Key = 56
	// Nine represents the  '9' keyboard key
	Nine Key = 57
	// F1 represents the 'F1' keyboard key
	F1 Key = 112
	// F2 represents the 'F2' keyboard key
	F2 Key = 113
	// F3 represents the 'F3' keyboard key
	F3 Key = 114
	// F4 represents the 'F4' keyboard key
	F4 Key = 115
	// F5 represents the 'F5' keyboard key
	F5 Key = 116
	// F6 represents the 'F6' keyboard key
	F6 Key = 117
	// F7 represents the 'F7' keyboard key
	F7 Key = 118
	// F8 represents the 'F8' keyboard key
	F8 Key = 119
	// F9 represents the 'F9' keyboard key
	F9 Key = 120
	// F10 represents the 'F10' keyboard key
	F10 Key = 121
	// F11 represents the 'F11' keyboard key
	F11 Key = 122
	// F12 represents the 'F12' keyboard key
	F12 Key = 123
	// A represents the 'A' keyboard key
	A Key = 65
	// B represents the 'B' keyboard key
	B Key = 66
	// C represents the 'C' keyboard key
	C Key = 67
	// D represents the 'D' keyboard key '
	D Key = 68
	// E represents the 'E' keyboard key
	E Key = 69
	// F represents the 'F' keyboard key
	F Key = 70
	// G represents the 'G' keyboard key
	G Key = 71
	// H represents the 'H' keyboard key
	H Key = 72
	// I represents the 'I' keyboard key
	I Key = 73
	// J represents the 'J' keyboard key
	J Key = 74
	// K represents the 'K' keyboard key
	K Key = 75
	// L represents the 'L' keyboard key
	L Key = 76
	// M represents the 'M' keyboard key
	M Key = 77
	// N represents the 'N' keyboard key
	N Key = 78
	// O represents the 'O' keyboard key
	O Key = 79
	// P represents the 'P' keyboard key
	P Key = 80
	// Q represents the 'Q' keyboard key
	Q Key = 81
	// R represents the 'R' keyboard key
	R Key = 82
	// S represents the 'S' keyboard key
	S Key = 83
	// T represents the 'T' keyboard key
	T Key = 84
	// U represents the 'U' keyboard key
	U Key = 85
	// V represents the 'V' keyboard key
	V Key = 86
	// W represents the 'W' keyboard key
	W Key = 87
	// X represents the 'X' keyboard key
	X Key = 88
	// Y represents the 'Y' keyboard key
	Y Key = 89
	// Z represents the 'Z' keyboard key
	Z Key = 90
	// NumLock represents the NumLock keyboard key on the numpad
	NumLock Key = 144
	// NumMultiply represents the NumMultiply keyboard key on the numpad
	NumMultiply Key = 106
	// NumDivide represents the NumDivide keyboard key on the numpad
	NumDivide Key = 111
	// NumAdd represents the NumAdd keyboard key on the numpad
	NumAdd Key = 107
	// NumSubtract represents the NumSubtract keyboard key on the numpad
	NumSubtract Key = 109
	// NumZero represents the NumZero keyboard key on the numpad
	NumZero Key = 96
	// NumOne represents the NumOne keyboard key on the numpad
	NumOne Key = 97
	// NumTwo represents the NumTwo keyboard key on the numpad
	NumTwo Key = 98
	// NumThree represents the NumThree keyboard key on the numpad
	NumThree Key = 99
	// NumFour represents the NumFour keyboard key on the numpad
	NumFour Key = 100
	// NumFive represents the NumFive keyboard key on the numpad
	NumFive Key = 101
	// NumSiz represents the NumSix keyboard key on the numpad
	NumSix Key = 102
	// NumSeven represents the NumSeven keyboard key on the numpad
	NumSeven Key = 103
	// NumEight represents the NumEight keyboard key on the numpad
	NumEight Key = 104
	// NumNine represents the NumNine keyboard key on the numpad
	NumNine Key = 105
	// NumDecimal represents the NumDecimal keyboard key on the numpad
	NumDecimal Key = 110
	// NumEnter represents the NumEnter keyboard key on the numpad
	NumEnter Key = 13
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
