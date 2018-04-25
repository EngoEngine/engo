//+build netgo android ios darwin,arm darwin,arm64

package engo

const (
	// KeyGrave represents the '`' keyboard key
	KeyGrave Key = 192
	// KeyDash represents the '-' keyboard key
	KeyDash Key = 189
	// KeyApostrophe represents the `'` keyboard key
	KeyApostrophe Key = 222
	// KeySemicolon represents the ';' keyboard key
	KeySemicolon Key = 186
	// KeyEquals reprsents the '=' keyboard key
	KeyEquals Key = 187
	// KeyComma represents the ',' keyboard key
	KeyComma Key = 188
	// KeyPeriod represents the '.' keyboard key
	KeyPeriod Key = 190
	// KeySlash represents the '/' keyboard key
	KeySlash Key = 191
	// KeyBackslash represents the '\' keyboard key
	KeyBackslash Key = 220
	// KeyBackspace represents the backspace keyboard key
	KeyBackspace Key = 8
	// KeyTab represents the tab keyboard key
	KeyTab Key = 9
	// KeyCapsLock represents the caps lock keyboard key
	KeyCapsLock Key = 20
	// KeySpace represents the space keyboard key
	KeySpace Key = 32
	// KeyEnter represents the enter keyboard key
	KeyEnter Key = 13
	// KeyEscape represents the escape keyboard key
	KeyEscape Key = 27
	// KeyInsert represents the insert keyboard key
	KeyInsert Key = 45
	// KeyPrintScreen represents the print screen keyboard key often
	// represented by 'Prt Scrn', 'Prt Scn', or 'Print Screen'
	KeyPrintScreen Key = 42
	// KeyDelete represents the delete keyboard key
	KeyDelete Key = 46
	// KeyPageUp represents the page up keyboard key
	KeyPageUp Key = 33
	// KeyPageDown represents the page down keyboard key
	KeyPageDown Key = 34
	// KeyHome represents the home keyboard key
	KeyHome Key = 36
	// KeyEnd represents the end keyboard key
	KeyEnd Key = 35
	// KeyPause represents the pause keyboard key
	KeyPause Key = 19
	// KeyScrollLock represents the scroll lock keyboard key
	KeyScrollLock Key = 145
	// KeyAllowLeft represents the arrow left keyboard key
	KeyArrowLeft Key = 37
	// KeyArrowRight represents the arrow right keyboard key
	KeyArrowRight Key = 39
	// KeyArrowDown represents the down arrow keyboard key
	KeyArrowDown Key = 40
	// KeyArrowUp represents the up arrow keyboard key
	KeyArrowUp Key = 38
	// KeyLeftBracket represents the '[' keyboard key
	KeyLeftBracket Key = 219
	// KeyLeftShift represents the left shift keyboard key
	KeyLeftShift Key = 16
	// KeyLeftControl represents the left control keyboard key
	KeyLeftControl Key = 17
	// KeyLeftSuper represents the left super keyboard key
	// (Windows key on Microsoft Windows, Command key on Apple OSX, and varies on Linux)
	KeyLeftSuper Key = 73
	// KeyLeftAlt represents the left alt keyboard key
	KeyLeftAlt Key = 18
	// KeyRightBracket represents the ']' keyboard key
	KeyRightBracket Key = 221
	// KeyRightShift represents the right shift keyboard key
	KeyRightShift Key = 16
	// KeyRightControl represents the right control keyboard key
	KeyRightControl Key = 17
	// KeyRightSuper represents the right super keyboard key
	// (Windows key on Microsoft Windows, Command key on Apple OSX, and varies on Linux)
	KeyRightSuper Key = 73
	// KeyRightAlt represents the left alt keyboard key
	KeyRightAlt Key = 18
	// KeyZero represents the '0' keyboard key
	KeyZero Key = 48
	// KeyOne represents the '1' keyboard key
	KeyOne Key = 49
	// KeyTwo represents the '2' keyboard key
	KeyTwo Key = 50
	// KeyThree represents the '3' keyboard key
	KeyThree Key = 51
	// KeyFour represents the '4' keyboard key
	KeyFour Key = 52
	// KeyFive represents the '5' keyboard key
	KeyFive Key = 53
	// KeySix represents the '6' keyboard key
	KeySix Key = 54
	// KeySeven represents the '7' keyboard key
	KeySeven Key = 55
	// KeyEight represents the '8' keyboard key
	KeyEight Key = 56
	// KeyNine represents the  '9' keyboard key
	KeyNine Key = 57
	// KeyF1 represents the 'F1' keyboard key
	KeyF1 Key = 112
	// KeyF2 represents the 'F2' keyboard key
	KeyF2 Key = 113
	// KeyF3 represents the 'F3' keyboard key
	KeyF3 Key = 114
	// KeyF4 represents the 'F4' keyboard key
	KeyF4 Key = 115
	// KeyF5 represents the 'F5' keyboard key
	KeyF5 Key = 116
	// KeyF6 represents the 'F6' keyboard key
	KeyF6 Key = 117
	// KeyF7 represents the 'F7' keyboard key
	KeyF7 Key = 118
	// KeyF8 represents the 'F8' keyboard key
	KeyF8 Key = 119
	// KeyF9 represents the 'F9' keyboard key
	KeyF9 Key = 120
	// KeyF10 represents the 'F10' keyboard key
	KeyF10 Key = 121
	// KeyF11 represents the 'F11' keyboard key
	KeyF11 Key = 122
	// KeyF12 represents the 'F12' keyboard key
	KeyF12 Key = 123
	// KeyA represents the 'A' keyboard key
	KeyA Key = 65
	// KeyB represents the 'B' keyboard key
	KeyB Key = 66
	// KeyC represents the 'C' keyboard key
	KeyC Key = 67
	// KeyD represents the 'D' keyboard key '
	KeyD Key = 68
	// KeyE represents the 'E' keyboard key
	KeyE Key = 69
	// KeyF represents the 'F' keyboard key
	KeyF Key = 70
	// KeyG represents the 'G' keyboard key
	KeyG Key = 71
	// KeyH represents the 'H' keyboard key
	KeyH Key = 72
	// KeyI represents the 'I' keyboard key
	KeyI Key = 73
	// KeyJ represents the 'J' keyboard key
	KeyJ Key = 74
	// KeyK represents the 'K' keyboard key
	KeyK Key = 75
	// KeyL represents the 'L' keyboard key
	KeyL Key = 76
	// KeyM represents the 'M' keyboard key
	KeyM Key = 77
	// KeyN represents the 'N' keyboard key
	KeyN Key = 78
	// KeyO represents the 'O' keyboard key
	KeyO Key = 79
	// KeyP represents the 'P' keyboard key
	KeyP Key = 80
	// KeyQ represents the 'Q' keyboard key
	KeyQ Key = 81
	// KeyR represents the 'R' keyboard key
	KeyR Key = 82
	// KeyS represents the 'S' keyboard key
	KeyS Key = 83
	// KeyT represents the 'T' keyboard key
	KeyT Key = 84
	// KeyU represents the 'U' keyboard key
	KeyU Key = 85
	// KeyV represents the 'V' keyboard key
	KeyV Key = 86
	// KeyW represents the 'W' keyboard key
	KeyW Key = 87
	// KeyX represents the 'X' keyboard key
	KeyX Key = 88
	// KeyY represents the 'Y' keyboard key
	KeyY Key = 89
	// KeyZ represents the 'Z' keyboard key
	KeyZ Key = 90
	// KeyNumLock represents the NumLock keyboard key on the numpad
	KeyNumLock Key = 144
	// KeyNumMultiply represents the NumMultiply keyboard key on the numpad
	KeyNumMultiply Key = 106
	// KeyNumDivide represents the NumDivide keyboard key on the numpad
	KeyNumDivide Key = 111
	// KeyNumAdd represents the NumAdd keyboard key on the numpad
	KeyNumAdd Key = 107
	// KeyNumSubtract represents the NumSubtract keyboard key on the numpad
	KeyNumSubtract Key = 109
	// KeyNumZero represents the NumZero keyboard key on the numpad
	KeyNumZero Key = 96
	// KeyNumOne represents the NumOne keyboard key on the numpad
	KeyNumOne Key = 97
	// KeyNumTwo represents the NumTwo keyboard key on the numpad
	KeyNumTwo Key = 98
	// KeyNumThree represents the NumThree keyboard key on the numpad
	KeyNumThree Key = 99
	// KeyNumFour represents the NumFour keyboard key on the numpad
	KeyNumFour Key = 100
	// KeyNumFive represents the NumFive keyboard key on the numpad
	KeyNumFive Key = 101
	// KeyNumSix represents the NumSix keyboard key on the numpad
	KeyNumSix Key = 102
	// KeyNumSeven represents the NumSeven keyboard key on the numpad
	KeyNumSeven Key = 103
	// KeyNumEight represents the NumEight keyboard key on the numpad
	KeyNumEight Key = 104
	// KeyNumNine represents the NumNine keyboard key on the numpad
	KeyNumNine Key = 105
	// KeyNumDecimal represents the NumDecimal keyboard key on the numpad
	KeyNumDecimal Key = 110
	// KeyNumEnter represents the NumEnter keyboard key on the numpad
	KeyNumEnter Key = 13
)
