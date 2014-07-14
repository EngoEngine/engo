// Copyright 2013 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package eng provides functionality for creating 2d games.
package eng

import (
	"github.com/errcw/glow/gl/2.1/gl"
	glfw "github.com/go-gl/glfw3"
	"log"
)

type Action int
type Key int
type Modifier int

var (
	MOVE    = Action(0)
	PRESS   = Action(1)
	RELEASE = Action(2)
	SHIFT   = Modifier(0x0001)
	CONTROL = Modifier(0x0002)
	ALT     = Modifier(0x0004)
	SUPER   = Modifier(0x0008)
)

// Common OpenGL constants
const (
	BlendZero                  = gl.ZERO
	BlendOne                   = gl.ONE
	BlendSrcColor              = gl.SRC_COLOR
	BlendOneMinusSrcColor      = gl.ONE_MINUS_SRC_COLOR
	BlendDstColor              = gl.DST_COLOR
	BlendOneMinusDstColor      = gl.ONE_MINUS_DST_COLOR
	BlendSrcAlpha              = gl.SRC_ALPHA
	BlendOneMinusSrcAlpha      = gl.ONE_MINUS_SRC_ALPHA
	BlendDstAlpha              = gl.DST_ALPHA
	BlendOneMinusDstAlpha      = gl.ONE_MINUS_DST_ALPHA
	FilterNearest              = gl.NEAREST
	FilterLinear               = gl.LINEAR
	FilterMipMap               = gl.LINEAR_MIPMAP_LINEAR
	FilterLinearMipMapLinear   = gl.LINEAR_MIPMAP_LINEAR
	FilterNearestMipMapLinear  = gl.NEAREST_MIPMAP_LINEAR
	FilterLinearMipMapNearest  = gl.LINEAR_MIPMAP_NEAREST
	FilterNearestMipMapNearest = gl.NEAREST_MIPMAP_NEAREST
	WrapClampToEdge            = gl.CLAMP_TO_EDGE
	WrapRepeat                 = gl.REPEAT
	WrapMirroredRepeat         = gl.MIRRORED_REPEAT
	Escape                     = Key(glfw.KeyEscape)
	F1                         = Key(glfw.KeyF1)
	F2                         = Key(glfw.KeyF2)
	F3                         = Key(glfw.KeyF3)
	F4                         = Key(glfw.KeyF4)
	F5                         = Key(glfw.KeyF5)
	F6                         = Key(glfw.KeyF6)
	F7                         = Key(glfw.KeyF7)
	F8                         = Key(glfw.KeyF8)
	F9                         = Key(glfw.KeyF9)
	F10                        = Key(glfw.KeyF10)
	F11                        = Key(glfw.KeyF11)
	F12                        = Key(glfw.KeyF12)
	F13                        = Key(glfw.KeyF13)
	F14                        = Key(glfw.KeyF14)
	F15                        = Key(glfw.KeyF15)
	F16                        = Key(glfw.KeyF16)
	F17                        = Key(glfw.KeyF17)
	F18                        = Key(glfw.KeyF18)
	F19                        = Key(glfw.KeyF19)
	F20                        = Key(glfw.KeyF20)
	F21                        = Key(glfw.KeyF21)
	F22                        = Key(glfw.KeyF22)
	F23                        = Key(glfw.KeyF23)
	F24                        = Key(glfw.KeyF24)
	F25                        = Key(glfw.KeyF25)
	Up                         = Key(glfw.KeyUp)
	Down                       = Key(glfw.KeyDown)
	Left                       = Key(glfw.KeyLeft)
	Right                      = Key(glfw.KeyRight)
	LeftShift                  = Key(glfw.KeyLeftShift)
	RightShift                 = Key(glfw.KeyRightShift)
	LeftControl                = Key(glfw.KeyLeftControl)
	RightControl               = Key(glfw.KeyRightControl)
	LeftAlt                    = Key(glfw.KeyLeftAlt)
	RightAlt                   = Key(glfw.KeyRightAlt)
	Tab                        = Key(glfw.KeyTab)
	Space                      = Key(glfw.KeySpace)
	Enter                      = Key(glfw.KeyEnter)
	Backspace                  = Key(glfw.KeyBackspace)
	Insert                     = Key(glfw.KeyInsert)
	Delete                     = Key(glfw.KeyDelete)
	PageUp                     = Key(glfw.KeyPageUp)
	PageDown                   = Key(glfw.KeyPageDown)
	Home                       = Key(glfw.KeyHome)
	End                        = Key(glfw.KeyEnd)
	Kp0                        = Key(glfw.KeyKp0)
	Kp1                        = Key(glfw.KeyKp1)
	Kp2                        = Key(glfw.KeyKp2)
	Kp3                        = Key(glfw.KeyKp3)
	Kp4                        = Key(glfw.KeyKp4)
	Kp5                        = Key(glfw.KeyKp5)
	Kp6                        = Key(glfw.KeyKp6)
	Kp7                        = Key(glfw.KeyKp7)
	Kp8                        = Key(glfw.KeyKp8)
	Kp9                        = Key(glfw.KeyKp9)
	KpDivide                   = Key(glfw.KeyKpDivide)
	KpMultiply                 = Key(glfw.KeyKpMultiply)
	KpSubtract                 = Key(glfw.KeyKpSubtract)
	KpAdd                      = Key(glfw.KeyKpAdd)
	KpDecimal                  = Key(glfw.KeyKpDecimal)
	KpEqual                    = Key(glfw.KeyKpEqual)
	KpEnter                    = Key(glfw.KeyKpEnter)
	NumLock                    = Key(glfw.KeyNumLock)
	CapsLock                   = Key(glfw.KeyCapsLock)
	ScrollLock                 = Key(glfw.KeyScrollLock)
	Pause                      = Key(glfw.KeyPause)
	LeftSuper                  = Key(glfw.KeyLeftSuper)
	RightSuper                 = Key(glfw.KeyRightSuper)
	Menu                       = Key(glfw.KeyMenu)
)

var (
	responder   Responder
	window      *glfw.Window
	config      *Config
	timing      *stats
	defaultFont *Font
	bgColor     *Color
)

// A Config holds settings for your game's window and application.
type Config struct {
	// Title is the name of the created window.
	// Default: Untitled
	Title string

	// Width and Height are hints about the size of the window. You
	// may not end up with the indicated size, so you should always
	// query eng for the true width and height after initialization.
	// Default: 1024 x 640
	Width  int
	Height int

	// Fullscreen tells eng whether to open windowed or fullscreen.
	// Default: false
	Fullscreen bool

	// Vsync enables or disables vertical sync which will limit the
	// number of frames rendered per second to your monitor's refresh
	// rate. This may or may not be supported on certain platforms.
	// Default: true
	Vsync bool

	// Resizable tells eng if it should request a window that can be
	// resized by the user of your game.
	// Default: false
	Resizable bool

	// Fsaa indicates how many samples to use for the multisampling
	// buffer. Generally it will be 1, 2, 4, 8, or 16.
	// Default: 1
	Fsaa int

	// PrintFPS turns on a logging of the frames per second to the
	// console every second.
	// Default: false
	LogFPS bool
}

func NewConfig() *Config {
	return &Config{"Untitled", 800, 600, false, true, false, 1, false}
}

func Run(title string, width, height int, fullscreen bool, r Responder) {
	RunConfig(&Config{title, width, height, fullscreen, true, false, 1, false}, r)
}

// Run should be called with a type that satisfies the Responder
// interface. Windows will be setup using your Config and a runloop
// will start, blocking the main thread and calling methods on the
// given responder.
func RunConfig(c *Config, r Responder) {
	config = c
	responder = r

	glfw.SetErrorCallback(func(err glfw.ErrorCode, desc string) {
		log.Fatal("GLFW error %v: %v\n", err, desc)
	})

	if ok := glfw.Init(); ok {
		defer glfw.Terminate()
	}

	if !config.Resizable {
		glfw.WindowHint(glfw.Resizable, 0)
	}
	glfw.WindowHint(glfw.Samples, config.Fsaa)

	width := config.Width
	height := config.Height

	monitor, err := glfw.GetPrimaryMonitor()
	if err != nil {
		log.Fatal(err)
	}
	mode, err := monitor.GetVideoMode()
	if err != nil {
		log.Fatal(err)
	}

	if config.Fullscreen {
		width = mode.Width
		height = mode.Height
		glfw.WindowHint(glfw.Decorated, 0)
	} else {
		monitor = nil
	}

	title := config.Title

	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)

	window, err = glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer window.Destroy()
	window.MakeContextCurrent()

	config.Width, config.Height = window.GetSize()

	bgColor = NewColor(0, 0, 0)

	if !config.Fullscreen {
		window.SetPosition((mode.Width-width)/2, (mode.Height-height)/2)
	}

	if config.Vsync {
		glfw.SwapInterval(1)
	}

	if err := gl.Init(); err != nil {
		log.Fatal(err)
	}

	window.SetSizeCallback(func(window *glfw.Window, w, h int) {
		config.Width, config.Height = window.GetSize()
		responder.Resize(w, h)
	})

	window.SetCursorPositionCallback(func(window *glfw.Window, x, y float64) {
		responder.Mouse(float32(x), float32(y), MOVE)
	})

	window.SetMouseButtonCallback(func(window *glfw.Window, b glfw.MouseButton, a glfw.Action, m glfw.ModifierKey) {
		x, y := window.GetCursorPosition()
		if a == glfw.Press {
			responder.Mouse(float32(x), float32(y), PRESS)
		} else {
			responder.Mouse(float32(x), float32(y), RELEASE)
		}
	})

	window.SetScrollCallback(func(window *glfw.Window, xoff, yoff float64) {
		responder.Scroll(float32(yoff))
	})

	window.SetKeyCallback(func(window *glfw.Window, k glfw.Key, s int, a glfw.Action, m glfw.ModifierKey) {
		if a == glfw.Press {
			responder.Key(Key(k), Modifier(m), PRESS)
		} else {
			responder.Key(Key(k), Modifier(m), RELEASE)
		}
	})

	window.SetCharacterCallback(func(window *glfw.Window, char uint) {
		responder.Type(rune(char))
	})

	responder.Setup()

	timing = NewStats(config.LogFPS)
	timing.Update()

	for !window.ShouldClose() {
		responder.Update(float32(timing.Dt))
		Clear(bgColor)
		responder.Draw()
		window.SwapBuffers()
		glfw.PollEvents()
		timing.Update()
	}
}

func Log(l ...interface{}) {
	log.Println(l...)
}

// Clear manually clears with a given color. Mostly used with a Canvas.
func Clear(color *Color) {
	gl.ClearColor(color.R, color.G, color.B, color.A)
	gl.Clear(gl.COLOR_BUFFER_BIT)
}

// Width returns the current window width.
func Width() int {
	return config.Width
}

// Height returns the current window height.
func Height() int {
	return config.Height
}

// Size returns the current window width and height.
func Size() (int, int) {
	return config.Width, config.Height
}

// SetSize sets the window width and height.
func SetSize(w, h int) {
	window.SetSize(w, h)
}

// Focused indicates if the game window is currently focused.
func Focused() bool {
	return window.GetAttribute(glfw.Focused) == gl.TRUE
}

// Exit closes the window and breaks out of the game loop.
func Exit() {
	window.SetShouldClose(true)
}

// MouseX returns the cursor's horizontal position within the window.
func MouseX() float64 {
	x, _ := window.GetCursorPosition()
	return x
}

// MouseY returns the cursor's vetical position within the window.
func MouseY() float64 {
	_, y := window.GetCursorPosition()
	return y
}

// MousePos returns the cursor's X and Y position within the window.
func MousePos() (float64, float64) {
	return window.GetCursorPosition()
}

// SetMousePos sets the cursor's X and Y position within the window.
func SetMousePos(x, y float64) {
	window.SetCursorPosition(x, y)
}

// SetMouseCursor shows or hides the cursor.
func SetMouseCursor(on bool) {
	if on {
		window.SetInputMode(glfw.Cursor, glfw.CursorNormal)
	} else {
		window.SetInputMode(glfw.Cursor, glfw.CursorHidden)
	}
}

// MousePressed takes a mouse button constant and indicates if it is
// currently pressed.
func MousePressed(b glfw.MouseButton) bool {
	return window.GetMouseButton(b) == glfw.Press
}

// KeyPressed takes a key constant and indicates if it is currently pressed.
func KeyPressed(k glfw.Key) bool {
	return window.GetKey(k) == glfw.Press
}

// SetBgColor sets the default opengl clear color.
func SetBgColor(c *Color) {
	bgColor.R = c.R
	bgColor.G = c.G
	bgColor.B = c.B
	bgColor.A = c.A
}

// Dt returns the time since the last frame.
func Dt() float32 {
	return float32(timing.Dt)
}

// Fps returns the number of frames being rendered each second.
func Fps() float32 {
	return float32(timing.Fps)
}

// DefaultFont returns eng's built in font, creating it on first use.
func DefaultFont() *Font {
	if defaultFont == nil {
		defaultFont = NewBitmapFont(dfontimg(), dfonttxt)
	}
	return defaultFont
}
