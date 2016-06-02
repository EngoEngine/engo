package engo

// Action corresponds to a control action
type Action int

// Key correspends to a keyboard key
type Key int

//Modifier represents a special key pressed along with another key
type Modifier int

var (
	// Move represents a move action
	Move = Action(0)
	// Press represents a press action
	Press = Action(1)
	// Release represents a release action
	Release = Action(2)
	// Neutral represents a neutral action
	Neutral = Action(99)
	// Shift represents the shift modifier
	Shift = Modifier(0x0001)
	// Control represents the control modifier
	Control = Modifier(0x0002)
	// Alt represents the alt modifier
	Alt = Modifier(0x0004)
	// Super represents the super modifier
	// (Windows key on Microsoft Windows, Command key on Apple OSX, and varies on Linux)
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
	MouseButton4      MouseButton = 3
	MouseButton5      MouseButton = 4
	MouseButton6      MouseButton = 5
	MouseButton7      MouseButton = 6
	MouseButton8      MouseButton = 7
	MouseButtonLast   MouseButton = 7
)

// those are default values for engo_js defined here because some of them are shared
// with engo_glfw.
// engo_glfw redefines the variables it needs to other values during init() so
var (
	Grave Key = 192
	// Dash represents the dash keyboard key '-'
	Dash Key = 189
	// Apostrophe represents the apostrophe keyboard key `'`
	Apostrophe Key = 222
	// Semicolon represents the semicolon keyboard key ';'
	Semicolon Key = 186
	// Equals reprsents the equals keyboard key '='
	Equals Key = 187
	// Comma represents the comma keyboard key ','
	Comma Key = 188
	// Period represents the period keyboard key '.'
	Period Key = 190
	// Slash represents the slash keyboard key '/'
	Slash Key = 191
	// Backslash represents the backslash keyboard key '\'
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
	// PrintScreen represents the print screen keyboard key 'prntscr'
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
	// LeftBracket represents the left bracket keyboard key '['
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
	// RightBracket represents the right bracket keyboard key ']'
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
	// Zero represents the zero keyboard key '0'
	Zero Key = 48
	// One represents the one keyboard key '1'
	One Key = 49
	// Two represents the two keyboard key '2'
	Two Key = 50
	// Three represents the three keyboard key '3'
	Three Key = 51
	// Four represents the four keyboard key '4'
	Four Key = 52
	// Five represents the five keyboard key '5'
	Five Key = 53
	// Six represents the six keyboard key '6'
	Six Key = 54
	// Seven represents the seven keyboard key '7'
	Seven Key = 55
	// Eight represents the eight keyboard key '8'
	Eight Key = 56
	// Nine represents the nine keyboard key '9'
	Nine Key = 57
	// F1 represents the f1 keyboard key
	F1 Key = 112
	// F2 represents the f2 keyboard key
	F2 Key = 113
	// F3 represents the f3 keyboard key
	F3 Key = 114
	// F4 represents the f4 keyboard key
	F4 Key = 115
	// F5 represents the f5 keyboard key
	F5 Key = 116
	// F6 represents the f6 keyboard key
	F6 Key = 117
	// F7 represents the f7 keyboard key
	F7 Key = 118
	// F8 represents the f8 keyboard key
	F8 Key = 119
	// F9 represents the f9 keyboard key
	F9 Key = 120
	// F10 represents the f10 keyboard key
	F10 Key = 121
	// F11 represents the f11 keyboard key
	F11 Key = 122
	// F12 represents the f12 keyboard key
	F12 Key = 123
	// A represents the a keyboard key 'A'
	A Key = 65
	// B represents the b keyboard key 'B'
	B Key = 66
	// C represents the c keyboard key 'C'
	C Key = 67
	// D represents the d keyboard key 'D'
	D Key = 68
	// E represents the e keyboard key 'E'
	E Key = 69
	// F represents the f keyboard key 'F'
	F Key = 70
	// G represents the g keyboard key 'G'
	G Key = 71
	// H represents the h keyboard key 'H'
	H Key = 72
	// I represents the i keyboard key 'I'
	I Key = 73
	// J represents the j keyboard key 'J'
	J Key = 74
	// K represents the k keyboard key 'K'
	K Key = 75
	// L represents the l keyboard key 'L'
	L Key = 76
	// M represents the m keyboard key 'M'
	M Key = 77
	// N represents the n keyboard key 'N'
	N Key = 78
	// O represents the o keyboard key 'O'
	O Key = 79
	// P represents the p keyboard key 'P'
	P Key = 80
	// Q represents the q keyboard key 'Q'
	Q Key = 81
	// R represents the r keyboard key 'R'
	R Key = 82
	// S represents the s keyboard key 'S'
	S Key = 83
	// T represents the t keyboard key 'T'
	T Key = 84
	// U represents the u keyboard key 'U'
	U Key = 85
	// V represents the v keyboard key 'V'
	V Key = 86
	// W represents the w keyboard key 'W'
	W Key = 87
	// X represents the x keyboard key 'X'
	X Key = 88
	// Y represents the y keyboard key 'Y'
	Y Key = 89
	// Z represents the z keyboard key 'Z'
	Z Key = 90
	// NumLock represents the NumLock keyboard key
	NumLock Key = 144
	// NumMultiply represents the NumLock keyboard key
	NumMultiply Key = 106
	// NumDivide represents the NumDivide keyboard key
	NumDivide Key = 111
	// NumAdd represents the NumAdd keyboard key
	NumAdd Key = 107
	// NumSubtract represents the NumSubtract keyboard key
	NumSubtract Key = 109
	// NumZero represents the NumZero keyboard key
	NumZero Key = 96
	// NumOne represents the NumOne keyboard key
	NumOne Key = 97
	// NumTwo represents the NumTwo keyboard key
	NumTwo Key = 98
	// NumThree represents the NumThree keyboard key
	NumThree Key = 99
	// NumFour represents the NumFour keyboard key
	NumFour Key = 100
	// NumFive represents the NumFive keyboard key
	NumFive Key = 101
	// NumSiz represents the NumSix keyboard key
	NumSix Key = 102
	// NumSeven represents the NumSeven keyboard key
	NumSeven Key = 103
	// NumEight represents the NumEight keyboard key
	NumEight Key = 104
	// NumNine represents the NumNine keyboard key
	NumNine Key = 105
	// NumDecimal represents the NumDecimal keyboard key
	NumDecimal Key = 110
	// NumEnter represents the NumEnter keyboard key
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
