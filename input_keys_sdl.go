// +build sdl

package engo

import "github.com/veandco/go-sdl2/sdl"

const (
	// KeyGrave represents the '`' keyboard key
	KeyGrave Key = Key(sdl.K_BACKQUOTE)
	// KeyDash represents the '-' keyboard key
	KeyDash Key = Key(sdl.K_MINUS)
	// KeyApostrophe represents the `'` keyboard key
	KeyApostrophe Key = Key(sdl.K_QUOTE)
	// KeySemicolon represents the ';' keyboard key
	KeySemicolon Key = Key(sdl.K_SEMICOLON)
	// KeyEquals reprsents the '=' keyboard key
	KeyEquals Key = Key(sdl.K_EQUALS)
	// KeyComma represents the ',' keyboard key
	KeyComma Key = Key(sdl.K_COMMA)
	// KeyPeriod represents the '.' keyboard key
	KeyPeriod Key = Key(sdl.K_PERIOD)
	// KeySlash represents the '/' keyboard key
	KeySlash Key = Key(sdl.K_SLASH)
	// KeyBackslash represents the '\' keyboard key
	KeyBackslash Key = Key(sdl.K_BACKSLASH)
	// KeyBackspace represents the backspace keyboard key
	KeyBackspace Key = Key(sdl.K_BACKSPACE)
	// KeyTab represents the tab keyboard key
	KeyTab Key = Key(sdl.K_TAB)
	// KeyCapsLock represents the caps lock keyboard key
	KeyCapsLock Key = Key(sdl.K_CAPSLOCK)
	// KeySpace represents the space keyboard key
	KeySpace Key = Key(sdl.K_SPACE)
	// KeyEnter represents the enter keyboard key
	KeyEnter Key = Key(sdl.K_RETURN)
	// KeyEscape represents the escape keyboard key
	KeyEscape Key = Key(sdl.K_ESCAPE)
	// KeyInsert represents the insert keyboard key
	KeyInsert Key = Key(sdl.K_INSERT)
	// KeyPrintScreen represents the print screen keyboard key often
	// represented by 'Prt Scrn', 'Prt Scn', or 'Print Screen'
	KeyPrintScreen Key = Key(sdl.K_PRINTSCREEN)
	// KeyDelete represents the delete keyboard key
	KeyDelete Key = Key(sdl.K_DELETE)
	// KeyPageUp represents the page up keyboard key
	KeyPageUp Key = Key(sdl.K_PAGEUP)
	// KeyPageDown represents the page down keyboard key
	KeyPageDown Key = Key(sdl.K_PAGEDOWN)
	// KeyHome represents the home keyboard key
	KeyHome Key = Key(sdl.K_HOME)
	// KeyEnd represents the end keyboard key
	KeyEnd Key = Key(sdl.K_END)
	// KeyPause represents the pause keyboard key
	KeyPause Key = Key(sdl.K_PAUSE)
	// KeyScrollLock represents the scroll lock keyboard key
	KeyScrollLock Key = Key(sdl.K_SCROLLLOCK)
	// KeyArrowLeft represents the arrow left keyboard key
	KeyArrowLeft Key = Key(sdl.K_LEFT)
	// KeyArrowRight represents the arrow right keyboard key
	KeyArrowRight Key = Key(sdl.K_RIGHT)
	// KeyArrowDown represents the down arrow keyboard key
	KeyArrowDown Key = Key(sdl.K_DOWN)
	// KeyArrowUp represents the up arrow keyboard key
	KeyArrowUp Key = Key(sdl.K_UP)
	// KeyLeftBracket represents the '[' keyboard key
	KeyLeftBracket Key = Key(sdl.K_LEFTBRACKET)
	// KeyLeftShift represents the left shift keyboard key
	KeyLeftShift Key = Key(sdl.K_LSHIFT)
	// KeyLeftControl represents the left control keyboard key
	KeyLeftControl Key = Key(sdl.K_LCTRL)
	// KeyLeftSuper represents the left super keyboard key
	// (Windows key on Microsoft Windows, Command key on Apple OSX, and varies on Linux)
	KeyLeftSuper Key = Key(sdl.K_LGUI)
	// KeyLeftAlt represents the left alt keyboard key
	KeyLeftAlt Key = Key(sdl.K_LALT)
	// KeyRightBracket represents the ']' keyboard key
	KeyRightBracket Key = Key(sdl.K_RIGHTBRACKET)
	// KeyRightShift represents the right shift keyboard key
	KeyRightShift Key = Key(sdl.K_RSHIFT)
	// KeyRightControl represents the right control keyboard key
	KeyRightControl Key = Key(sdl.K_RCTRL)
	// KeyRightSuper represents the right super keyboard key
	// (Windows key on Microsoft Windows, Command key on Apple OSX, and varies on Linux)
	KeyRightSuper Key = Key(sdl.K_RGUI)
	// KeyRightAlt represents the left alt keyboard key
	KeyRightAlt Key = Key(sdl.K_RALT)
	// KeyZero represents the '0' keyboard key
	KeyZero Key = Key(sdl.K_0)
	// KeyOne represents the '1' keyboard key
	KeyOne Key = Key(sdl.K_1)
	// KeyTwo represents the '2' keyboard key
	KeyTwo Key = Key(sdl.K_2)
	// KeyThree represents the '3' keyboard key
	KeyThree Key = Key(sdl.K_3)
	// KeyFour represents the '4' keyboard key
	KeyFour Key = Key(sdl.K_4)
	// KeyFive represents the '5' keyboard key
	KeyFive Key = Key(sdl.K_5)
	// KeySix represents the '6' keyboard key
	KeySix Key = Key(sdl.K_6)
	// KeySeven represents the '7' keyboard key
	KeySeven Key = Key(sdl.K_7)
	// KeyEight represents the '8' keyboard key
	KeyEight Key = Key(sdl.K_8)
	// KeyNine represents the  '9' keyboard key
	KeyNine Key = Key(sdl.K_9)
	// KeyF1 represents the 'F1' keyboard key
	KeyF1 Key = Key(sdl.K_F1)
	// KeyF2 represents the 'F2' keyboard key
	KeyF2 Key = Key(sdl.K_F2)
	// KeyF3 represents the 'F3' keyboard key
	KeyF3 Key = Key(sdl.K_F3)
	// KeyF4 represents the 'F4' keyboard key
	KeyF4 Key = Key(sdl.K_F4)
	// KeyF5 represents the 'F5' keyboard key
	KeyF5 Key = Key(sdl.K_F5)
	// KeyF6 represents the 'F6' keyboard key
	KeyF6 Key = Key(sdl.K_F6)
	// KeyF7 represents the 'F7' keyboard key
	KeyF7 Key = Key(sdl.K_F7)
	// KeyF8 represents the 'F8' keyboard key
	KeyF8 Key = Key(sdl.K_F8)
	// KeyF9 represents the 'F9' keyboard key
	KeyF9 Key = Key(sdl.K_F9)
	// KeyF10 represents the 'F10' keyboard key
	KeyF10 Key = Key(sdl.K_F10)
	// KeyF11 represents the 'F11' keyboard key
	KeyF11 Key = Key(sdl.K_F11)
	// KeyF12 represents the 'F12' keyboard key
	KeyF12 Key = Key(sdl.K_F12)
	// KeyA represents the 'A' keyboard key
	KeyA Key = Key(sdl.K_a)
	// KeyB represents the 'B' keyboard key
	KeyB Key = Key(sdl.K_b)
	// KeyC represents the 'C' keyboard key
	KeyC Key = Key(sdl.K_c)
	// KeyD represents the 'D' keyboard key '
	KeyD Key = Key(sdl.K_d)
	// KeyE represents the 'E' keyboard key
	KeyE Key = Key(sdl.K_e)
	// KeyF represents the 'F' keyboard key
	KeyF Key = Key(sdl.K_f)
	// KeyG represents the 'G' keyboard key
	KeyG Key = Key(sdl.K_g)
	// KeyH represents the 'H' keyboard key
	KeyH Key = Key(sdl.K_h)
	// KeyI represents the 'I' keyboard key
	KeyI Key = Key(sdl.K_i)
	// KeyJ represents the 'J' keyboard key
	KeyJ Key = Key(sdl.K_j)
	// KeyK represents the 'K' keyboard key
	KeyK Key = Key(sdl.K_k)
	// KeyL represents the 'L' keyboard key
	KeyL Key = Key(sdl.K_l)
	// KeyM represents the 'M' keyboard key
	KeyM Key = Key(sdl.K_m)
	// KeyN represents the 'N' keyboard key
	KeyN Key = Key(sdl.K_n)
	// KeyO represents the 'O' keyboard key
	KeyO Key = Key(sdl.K_o)
	// KeyP represents the 'P' keyboard key
	KeyP Key = Key(sdl.K_p)
	// KeyQ represents the 'Q' keyboard key
	KeyQ Key = Key(sdl.K_q)
	// KeyR represents the 'R' keyboard key
	KeyR Key = Key(sdl.K_r)
	// KeyS represents the 'S' keyboard key
	KeyS Key = Key(sdl.K_s)
	// KeyT represents the 'T' keyboard key
	KeyT Key = Key(sdl.K_t)
	// KeyU represents the 'U' keyboard key
	KeyU Key = Key(sdl.K_u)
	// KeyV represents the 'V' keyboard key
	KeyV Key = Key(sdl.K_v)
	// KeyW represents the 'W' keyboard key
	KeyW Key = Key(sdl.K_w)
	// KeyX represents the 'X' keyboard key
	KeyX Key = Key(sdl.K_x)
	// KeyY represents the 'Y' keyboard key
	KeyY Key = Key(sdl.K_y)
	// KeyZ represents the 'Z' keyboard key
	KeyZ Key = Key(sdl.K_z)
	// KeyNumLock represents the NumLock keyboard key on the numpad
	KeyNumLock Key = Key(sdl.K_NUMLOCKCLEAR)
	// KeyNumMultiply represents the NumMultiply keyboard key on the numpad
	KeyNumMultiply Key = Key(sdl.K_KP_MULTIPLY)
	// KeyNumDivide represents the NumDivide keyboard key on the numpad
	KeyNumDivide Key = Key(sdl.K_KP_DIVIDE)
	// KeyNumAdd represents the NumAdd keyboard key on the numpad
	KeyNumAdd Key = Key(sdl.K_KP_PLUS)
	// KeyNumSubtract represents the NumSubtract keyboard key on the numpad
	KeyNumSubtract Key = Key(sdl.K_KP_MINUS)
	// KeyNumZero represents the NumZero keyboard key on the numpad
	KeyNumZero Key = Key(sdl.K_KP_0)
	// KeyNumOne represents the NumOne keyboard key on the numpad
	KeyNumOne Key = Key(sdl.K_KP_1)
	// KeyNumTwo represents the NumTwo keyboard key on the numpad
	KeyNumTwo Key = Key(sdl.K_KP_2)
	// KeyNumThree represents the NumThree keyboard key on the numpad
	KeyNumThree Key = Key(sdl.K_KP_3)
	// KeyNumFour represents the NumFour keyboard key on the numpad
	KeyNumFour Key = Key(sdl.K_KP_4)
	// KeyNumFive represents the NumFive keyboard key on the numpad
	KeyNumFive Key = Key(sdl.K_KP_5)
	// KeyNumSix represents the NumSix keyboard key on the numpad
	KeyNumSix Key = Key(sdl.K_KP_6)
	// KeyNumSeven represents the NumSeven keyboard key on the numpad
	KeyNumSeven Key = Key(sdl.K_KP_7)
	// KeyNumEight represents the NumEight keyboard key on the numpad
	KeyNumEight Key = Key(sdl.K_KP_8)
	// KeyNumNine represents the NumNine keyboard key on the numpad
	KeyNumNine Key = Key(sdl.K_KP_9)
	// KeyNumDecimal represents the NumDecimal keyboard key on the numpad
	KeyNumDecimal Key = Key(sdl.K_KP_DECIMAL)
	// KeyNumEnter represents the NumEnter keyboard key on the numpad
	KeyNumEnter Key = Key(sdl.K_KP_ENTER)
)
