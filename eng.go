package eng

import (
	"bytes"
	gl "github.com/chsc/gogl/gl21"
	"github.com/go-gl/glfw"
	"log"
)

const (
	NOKEY       = -1
	Esc         = glfw.KeyEsc
	F1          = glfw.KeyF1
	F2          = glfw.KeyF2
	F3          = glfw.KeyF3
	F4          = glfw.KeyF4
	F5          = glfw.KeyF5
	F6          = glfw.KeyF6
	F7          = glfw.KeyF7
	F8          = glfw.KeyF8
	F9          = glfw.KeyF9
	F10         = glfw.KeyF10
	F11         = glfw.KeyF11
	F12         = glfw.KeyF12
	F13         = glfw.KeyF13
	F14         = glfw.KeyF14
	F15         = glfw.KeyF15
	F16         = glfw.KeyF16
	F17         = glfw.KeyF17
	F18         = glfw.KeyF18
	F19         = glfw.KeyF19
	F20         = glfw.KeyF20
	F21         = glfw.KeyF21
	F22         = glfw.KeyF22
	F23         = glfw.KeyF23
	F24         = glfw.KeyF24
	F25         = glfw.KeyF25
	Up          = glfw.KeyUp
	Down        = glfw.KeyDown
	Left        = glfw.KeyLeft
	Right       = glfw.KeyRight
	Lshift      = glfw.KeyLshift
	Rshift      = glfw.KeyRshift
	Lctrl       = glfw.KeyLctrl
	Rctrl       = glfw.KeyRctrl
	Lalt        = glfw.KeyLalt
	Ralt        = glfw.KeyRalt
	Tab         = glfw.KeyTab
	Enter       = glfw.KeyEnter
	Backspace   = glfw.KeyBackspace
	Insert      = glfw.KeyInsert
	Del         = glfw.KeyDel
	Pageup      = glfw.KeyPageup
	Pagedown    = glfw.KeyPagedown
	Home        = glfw.KeyHome
	End         = glfw.KeyEnd
	KP0         = glfw.KeyKP0
	KP1         = glfw.KeyKP1
	KP2         = glfw.KeyKP2
	KP3         = glfw.KeyKP3
	KP4         = glfw.KeyKP4
	KP5         = glfw.KeyKP5
	KP6         = glfw.KeyKP6
	KP7         = glfw.KeyKP7
	KP8         = glfw.KeyKP8
	KP9         = glfw.KeyKP9
	KPDivide    = glfw.KeyKPDivide
	KPMultiply  = glfw.KeyKPMultiply
	KPSubtract  = glfw.KeyKPSubtract
	KPAdd       = glfw.KeyKPAdd
	KPDecimal   = glfw.KeyKPDecimal
	KPEqual     = glfw.KeyKPEqual
	KPEnter     = glfw.KeyKPEnter
	KPNumlock   = glfw.KeyKPNumlock
	Capslock    = glfw.KeyCapslock
	Scrolllock  = glfw.KeyScrolllock
	Pause       = glfw.KeyPause
	Lsuper      = glfw.KeyLsuper
	Rsuper      = glfw.KeyRsuper
	Menu        = glfw.KeyMenu
	MouseLeft   = glfw.MouseLeft
	MouseRight  = glfw.MouseRight
	MouseMiddle = glfw.MouseMiddle
)

var (
	responder   Responder
	config      *Config
	timing      *stats
	batch       *Batch
	DefaultFont *Font
)

type Config struct {
	Title      string
	Width      int
	Height     int
	Fullscreen bool
	Vsync      bool
	Resizable  bool
	Fssa       int
}

type Responder interface {
	Init(s *Config)
	Open()
	Close()
	Update(dt float32)
	Draw()
	MouseMove(x, y int)
	MouseDown(x, y, b int)
	MouseUp(x, y, b int)
	KeyType(k int)
	KeyDown(k int)
	KeyUp(k int)
	Resize(w, h int)
}

func Run(r Responder) {
	responder = r

	if err := glfw.Init(); err != nil {
		log.Fatal(err)
	}
	defer glfw.Terminate()

	config = &Config{"Untitled", 1024, 640, false, true, false, 1}
	responder.Init(config)

	glfw.OpenWindowHint(glfw.OpenGLVersionMajor, 2)
	glfw.OpenWindowHint(glfw.OpenGLVersionMinor, 1)
	if !config.Resizable {
		glfw.OpenWindowHint(glfw.WindowNoResize, 1)
	}
	glfw.OpenWindowHint(glfw.FsaaSamples, config.Fssa)

	width := config.Width
	height := config.Height
	mode := glfw.DesktopMode()
	flag := glfw.Windowed

	if config.Fullscreen {
		flag = glfw.Fullscreen
		width = mode.W
		height = mode.H
	}

	if err := glfw.OpenWindow(width, height, 0, 0, 0, 0, 0, 0, flag); err != nil {
		log.Fatal(err)
	}
	defer glfw.CloseWindow()

	config.Width, config.Height = glfw.WindowSize()

	if err := gl.Init(); err != nil {
		log.Fatal(err)
	}

	batch = NewBatch()

	DefaultFont = NewFont(NewTexture(bytes.NewBuffer(Terminal())), 16, 16, "!\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~ ✵웃世界¢¥¤§©¨«¬£ª±²³´¶·¸¹º»¼½¾¿☐☑═║╔╗╚╝╠╣╦╩╬░▒▓☺☻☼♀♂▀▁▂▃▄▅▆▇█ÐÑÒÓÔÕÖÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏØÙÚÛÜÝàáâãäåèéêëìíîïðñòóôõö÷ùúûüýÿ♥♦♣♠♪♬æçø←↑→↓↔↕®‼ꀥ")

	responder.Open()

	if !config.Fullscreen {
		glfw.SetWindowPos((mode.W-width)/2, (mode.H-height)/2)
	}

	if config.Vsync {
		glfw.SetSwapInterval(1)
	}

	glfw.SetWindowTitle(config.Title)

	glfw.SetWindowSizeCallback(func(w, h int) {
		config.Width, config.Height = glfw.WindowSize()
		batch.Resize()
		responder.Resize(w, h)
	})

	glfw.SetMousePosCallback(func(x, y int) {
		responder.MouseMove(x, config.Height-y)
	})

	glfw.SetMouseButtonCallback(func(b, s int) {
		x, y := glfw.MousePos()
		if s == glfw.KeyPress {
			responder.MouseDown(x, config.Height-y, b)
		} else {
			responder.MouseUp(x, config.Height-y, b)
		}
	})

	glfw.SetKeyCallback(func(k, s int) {
		if s == glfw.KeyPress {
			responder.KeyDown(k)
		} else {
			responder.KeyUp(k)
		}
	})

	glfw.SetCharCallback(func(k, s int) {
		if s == glfw.KeyPress {
			responder.KeyType(k)
		}
	})

	timing = NewStats(true)
	timing.Update()
	for glfw.WindowParam(glfw.Opened) == 1 {
		responder.Update(float32(timing.Dt))
		gl.Clear(gl.COLOR_BUFFER_BIT)
		responder.Draw()
		glfw.SwapBuffers()
		timing.Update()
	}
	responder.Close()
}

func Log(l interface{}) {
	log.Println(l)
}

func Width() int {
	return config.Width
}

func Height() int {
	return config.Height
}

func Size() (int, int) {
	return config.Width, config.Height
}

func SetSize(w, h int) {
	glfw.SetWindowSize(w, h)
}

func Focused() bool {
	return glfw.WindowParam(glfw.Active) == gl.TRUE
}

func Exit() {
	glfw.CloseWindow()
}

func MouseX() int {
	x, _ := glfw.MousePos()
	return x
}

func MouseY() int {
	_, y := glfw.MousePos()
	return config.Height - y
}

func MousePos() (int, int) {
	return glfw.MousePos()
}

func SetMousePos(x, y int) {
	glfw.SetMousePos(x, y)
}

func SetMouseCursor(on bool) {
	if on {
		glfw.Enable(glfw.MouseCursor)
	} else {
		glfw.Disable(glfw.MouseCursor)
	}
}

func MousePressed(b int) bool {
	return glfw.MouseButton(b) == glfw.KeyPress
}

func KeyPressed(k int) bool {
	return glfw.Key(k) == glfw.KeyPress
}

func SetKeyRepeat(repeat bool) {
	if repeat {
		glfw.Enable(glfw.KeyRepeat)
	} else {
		glfw.Disable(glfw.KeyRepeat)
	}
}

func SetBgColor(c *Color) {
	gl.ClearColor(gl.Float(c.R), gl.Float(c.G), gl.Float(c.B), gl.Float(c.A))
}

func SetColor(c *Color) {
	batch.SetColor(c)
}
