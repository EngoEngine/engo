// +build darwin,!arm,!arm64 linux windows
// +build !ios,!android,!netgo

package engo

import "github.com/go-gl/glfw/v3.1/glfw"

const (
	// KeyGrave represents the '`' keyboard key
	KeyGrave Key = Key(glfw.KeyGraveAccent)
	// KeyDash represents the '-' keyboard key
	KeyDash Key = Key(glfw.KeyMinus)
	// KeyApostrophe represents the `'` keyboard key
	KeyApostrophe Key = Key(glfw.KeyApostrophe)
	// KeySemicolon represents the ';' keyboard key
	KeySemicolon Key = Key(glfw.KeySemicolon)
	// KeyEquals reprsents the '=' keyboard key
	KeyEquals Key = Key(glfw.KeyEqual)
	// KeyComma represents the ',' keyboard key
	KeyComma Key = Key(glfw.KeyComma)
	// KeyPeriod represents the '.' keyboard key
	KeyPeriod Key = Key(glfw.KeyPeriod)
	// KeySlash represents the '/' keyboard key
	KeySlash Key = Key(glfw.KeySlash)
	// KeyBackslash represents the '\' keyboard key
	KeyBackslash Key = Key(glfw.KeyBackslash)
	// KeyBackspace represents the backspace keyboard key
	KeyBackspace Key = Key(glfw.KeyBackspace)
	// KeyTab represents the tab keyboard key
	KeyTab Key = Key(glfw.KeyTab)
	// KeyCapsLock represents the caps lock keyboard key
	KeyCapsLock Key = Key(glfw.KeyCapsLock)
	// KeySpace represents the space keyboard key
	KeySpace Key = Key(glfw.KeySpace)
	// KeyEnter represents the enter keyboard key
	KeyEnter Key = Key(glfw.KeyEnter)
	// KeyEscape represents the escape keyboard key
	KeyEscape Key = Key(glfw.KeyEscape)
	// KeyInsert represents the insert keyboard key
	KeyInsert Key = Key(glfw.KeyInsert)
	// KeyPrintScreen represents the print screen keyboard key often
	// represented by 'Prt Scrn', 'Prt Scn', or 'Print Screen'
	KeyPrintScreen Key = Key(glfw.KeyPrintScreen)
	// KeyDelete represents the delete keyboard key
	KeyDelete Key = Key(glfw.KeyDelete)
	// KeyPageUp represents the page up keyboard key
	KeyPageUp Key = Key(glfw.KeyPageUp)
	// KeyPageDown represents the page down keyboard key
	KeyPageDown Key = Key(glfw.KeyPageDown)
	// KeyHome represents the home keyboard key
	KeyHome Key = Key(glfw.KeyHome)
	// KeyEnd represents the end keyboard key
	KeyEnd Key = Key(glfw.KeyEnd)
	// KeyPause represents the pause keyboard key
	KeyPause Key = Key(glfw.KeyPause)
	// KeyScrollLock represents the scroll lock keyboard key
	KeyScrollLock Key = Key(glfw.KeyScrollLock)
	// KeyArrowLeft represents the arrow left keyboard key
	KeyArrowLeft Key = Key(glfw.KeyLeft)
	// KeyArrowRight represents the arrow right keyboard key
	KeyArrowRight Key = Key(glfw.KeyRight)
	// KeyArrowDown represents the down arrow keyboard key
	KeyArrowDown Key = Key(glfw.KeyDown)
	// KeyArrowUp represents the up arrow keyboard key
	KeyArrowUp Key = Key(glfw.KeyUp)
	// KeyLeftBracket represents the '[' keyboard key
	KeyLeftBracket Key = Key(glfw.KeyLeftBracket)
	// KeyLeftShift represents the left shift keyboard key
	KeyLeftShift Key = Key(glfw.KeyLeftShift)
	// KeyLeftControl represents the left control keyboard key
	KeyLeftControl Key = Key(glfw.KeyLeftControl)
	// KeyLeftSuper represents the left super keyboard key
	// (Windows key on Microsoft Windows, Command key on Apple OSX, and varies on Linux)
	KeyLeftSuper Key = Key(glfw.KeyLeftSuper)
	// KeyLeftAlt represents the left alt keyboard key
	KeyLeftAlt Key = Key(glfw.KeyLeftAlt)
	// KeyRightBracket represents the ']' keyboard key
	KeyRightBracket Key = Key(glfw.KeyRightBracket)
	// KeyRightShift represents the right shift keyboard key
	KeyRightShift Key = Key(glfw.KeyRightShift)
	// KeyRightControl represents the right control keyboard key
	KeyRightControl Key = Key(glfw.KeyRightControl)
	// KeyRightSuper represents the right super keyboard key
	// (Windows key on Microsoft Windows, Command key on Apple OSX, and varies on Linux)
	KeyRightSuper Key = Key(glfw.KeyRightSuper)
	// KeyRightAlt represents the left alt keyboard key
	KeyRightAlt Key = Key(glfw.KeyRightAlt)
	// KeyZero represents the '0' keyboard key
	KeyZero Key = Key(glfw.Key0)
	// KeyOne represents the '1' keyboard key
	KeyOne Key = Key(glfw.Key1)
	// KeyTwo represents the '2' keyboard key
	KeyTwo Key = Key(glfw.Key2)
	// KeyThree represents the '3' keyboard key
	KeyThree Key = Key(glfw.Key3)
	// KeyFour represents the '4' keyboard key
	KeyFour Key = Key(glfw.Key4)
	// KeyFive represents the '5' keyboard key
	KeyFive Key = Key(glfw.Key5)
	// KeySix represents the '6' keyboard key
	KeySix Key = Key(glfw.Key6)
	// KeySeven represents the '7' keyboard key
	KeySeven Key = Key(glfw.Key7)
	// KeyEight represents the '8' keyboard key
	KeyEight Key = Key(glfw.Key8)
	// KeyNine represents the  '9' keyboard key
	KeyNine Key = Key(glfw.Key9)
	// KeyF1 represents the 'F1' keyboard key
	KeyF1 Key = Key(glfw.KeyF1)
	// KeyF2 represents the 'F2' keyboard key
	KeyF2 Key = Key(glfw.KeyF2)
	// KeyF3 represents the 'F3' keyboard key
	KeyF3 Key = Key(glfw.KeyF3)
	// KeyF4 represents the 'F4' keyboard key
	KeyF4 Key = Key(glfw.KeyF4)
	// KeyF5 represents the 'F5' keyboard key
	KeyF5 Key = Key(glfw.KeyF5)
	// KeyF6 represents the 'F6' keyboard key
	KeyF6 Key = Key(glfw.KeyF6)
	// KeyF7 represents the 'F7' keyboard key
	KeyF7 Key = Key(glfw.KeyF7)
	// KeyF8 represents the 'F8' keyboard key
	KeyF8 Key = Key(glfw.KeyF8)
	// KeyF9 represents the 'F9' keyboard key
	KeyF9 Key = Key(glfw.KeyF9)
	// KeyF10 represents the 'F10' keyboard key
	KeyF10 Key = Key(glfw.KeyF10)
	// KeyF11 represents the 'F11' keyboard key
	KeyF11 Key = Key(glfw.KeyF11)
	// KeyF12 represents the 'F12' keyboard key
	KeyF12 Key = Key(glfw.KeyF12)
	// KeyA represents the 'A' keyboard key
	KeyA Key = Key(glfw.KeyA)
	// KeyB represents the 'B' keyboard key
	KeyB Key = Key(glfw.KeyB)
	// KeyC represents the 'C' keyboard key
	KeyC Key = Key(glfw.KeyC)
	// KeyD represents the 'D' keyboard key '
	KeyD Key = Key(glfw.KeyD)
	// KeyE represents the 'E' keyboard key
	KeyE Key = Key(glfw.KeyE)
	// KeyF represents the 'F' keyboard key
	KeyF Key = Key(glfw.KeyF)
	// KeyG represents the 'G' keyboard key
	KeyG Key = Key(glfw.KeyG)
	// KeyH represents the 'H' keyboard key
	KeyH Key = Key(glfw.KeyH)
	// KeyI represents the 'I' keyboard key
	KeyI Key = Key(glfw.KeyI)
	// KeyJ represents the 'J' keyboard key
	KeyJ Key = Key(glfw.KeyJ)
	// KeyK represents the 'K' keyboard key
	KeyK Key = Key(glfw.KeyK)
	// KeyL represents the 'L' keyboard key
	KeyL Key = Key(glfw.KeyL)
	// KeyM represents the 'M' keyboard key
	KeyM Key = Key(glfw.KeyM)
	// KeyN represents the 'N' keyboard key
	KeyN Key = Key(glfw.KeyN)
	// KeyO represents the 'O' keyboard key
	KeyO Key = Key(glfw.KeyO)
	// KeyP represents the 'P' keyboard key
	KeyP Key = Key(glfw.KeyP)
	// KeyQ represents the 'Q' keyboard key
	KeyQ Key = Key(glfw.KeyQ)
	// KeyR represents the 'R' keyboard key
	KeyR Key = Key(glfw.KeyR)
	// KeyS represents the 'S' keyboard key
	KeyS Key = Key(glfw.KeyS)
	// KeyT represents the 'T' keyboard key
	KeyT Key = Key(glfw.KeyT)
	// KeyU represents the 'U' keyboard key
	KeyU Key = Key(glfw.KeyU)
	// KeyV represents the 'V' keyboard key
	KeyV Key = Key(glfw.KeyV)
	// KeyW represents the 'W' keyboard key
	KeyW Key = Key(glfw.KeyW)
	// KeyX represents the 'X' keyboard key
	KeyX Key = Key(glfw.KeyX)
	// KeyY represents the 'Y' keyboard key
	KeyY Key = Key(glfw.KeyY)
	// KeyZ represents the 'Z' keyboard key
	KeyZ Key = Key(glfw.KeyZ)
	// KeyNumLock represents the NumLock keyboard key on the numpad
	KeyNumLock Key = Key(glfw.KeyNumLock)
	// KeyNumMultiply represents the NumMultiply keyboard key on the numpad
	KeyNumMultiply Key = Key(glfw.KeyKPMultiply)
	// KeyNumDivide represents the NumDivide keyboard key on the numpad
	KeyNumDivide Key = Key(glfw.KeyKPDivide)
	// KeyNumAdd represents the NumAdd keyboard key on the numpad
	KeyNumAdd Key = Key(glfw.KeyKPAdd)
	// KeyNumSubtract represents the NumSubtract keyboard key on the numpad
	KeyNumSubtract Key = Key(glfw.KeyKPSubtract)
	// KeyNumZero represents the NumZero keyboard key on the numpad
	KeyNumZero Key = Key(glfw.KeyKP0)
	// KeyNumOne represents the NumOne keyboard key on the numpad
	KeyNumOne Key = Key(glfw.KeyKP1)
	// KeyNumTwo represents the NumTwo keyboard key on the numpad
	KeyNumTwo Key = Key(glfw.KeyKP2)
	// KeyNumThree represents the NumThree keyboard key on the numpad
	KeyNumThree Key = Key(glfw.KeyKP3)
	// KeyNumFour represents the NumFour keyboard key on the numpad
	KeyNumFour Key = Key(glfw.KeyKP4)
	// KeyNumFive represents the NumFive keyboard key on the numpad
	KeyNumFive Key = Key(glfw.KeyKP5)
	// KeyNumSix represents the NumSix keyboard key on the numpad
	KeyNumSix Key = Key(glfw.KeyKP6)
	// KeyNumSeven represents the NumSeven keyboard key on the numpad
	KeyNumSeven Key = Key(glfw.KeyKP7)
	// KeyNumEight represents the NumEight keyboard key on the numpad
	KeyNumEight Key = Key(glfw.KeyKP8)
	// KeyNumNine represents the NumNine keyboard key on the numpad
	KeyNumNine Key = Key(glfw.KeyKP9)
	// KeyNumDecimal represents the NumDecimal keyboard key on the numpad
	KeyNumDecimal Key = Key(glfw.KeyKPDecimal)
	// KeyNumEnter represents the NumEnter keyboard key on the numpad
	KeyNumEnter Key = Key(glfw.KeyKPEnter)
)
