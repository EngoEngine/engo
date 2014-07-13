// Copyright 2013 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package eng provides functionality for creating 2d games.
package eng

import (
	gl "github.com/chsc/gogl/gl32"
	glfw "github.com/go-gl/glfw3"
	"image"
	"image/draw"
	"log"
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
	Escape                     = glfw.KeyEsc
	F1                         = glfw.KeyF1
	F2                         = glfw.KeyF2
	F3                         = glfw.KeyF3
	F4                         = glfw.KeyF4
	F5                         = glfw.KeyF5
	F6                         = glfw.KeyF6
	F7                         = glfw.KeyF7
	F8                         = glfw.KeyF8
	F9                         = glfw.KeyF9
	F10                        = glfw.KeyF10
	F11                        = glfw.KeyF11
	F12                        = glfw.KeyF12
	F13                        = glfw.KeyF13
	F14                        = glfw.KeyF14
	F15                        = glfw.KeyF15
	F16                        = glfw.KeyF16
	F17                        = glfw.KeyF17
	F18                        = glfw.KeyF18
	F19                        = glfw.KeyF19
	F20                        = glfw.KeyF20
	F21                        = glfw.KeyF21
	F22                        = glfw.KeyF22
	F23                        = glfw.KeyF23
	F24                        = glfw.KeyF24
	F25                        = glfw.KeyF25
	Up                         = glfw.KeyUp
	Down                       = glfw.KeyDown
	Left                       = glfw.KeyLeft
	Right                      = glfw.KeyRight
	Lshift                     = glfw.KeyLshift
	Rshift                     = glfw.KeyRshift
	Lctrl                      = glfw.KeyLctrl
	Rctrl                      = glfw.KeyRctrl
	Lalt                       = glfw.KeyLalt
	Ralt                       = glfw.KeyRalt
	Tab                        = glfw.KeyTab
	Space                      = glfw.KeySpace
	Enter                      = glfw.KeyEnter
	Backspace                  = glfw.KeyBackspace
	Insert                     = glfw.KeyInsert
	Del                        = glfw.KeyDel
	Pageup                     = glfw.KeyPageup
	Pagedown                   = glfw.KeyPagedown
	Home                       = glfw.KeyHome
	End                        = glfw.KeyEnd
	KP0                        = glfw.KeyKP0
	KP1                        = glfw.KeyKP1
	KP2                        = glfw.KeyKP2
	KP3                        = glfw.KeyKP3
	KP4                        = glfw.KeyKP4
	KP5                        = glfw.KeyKP5
	KP6                        = glfw.KeyKP6
	KP7                        = glfw.KeyKP7
	KP8                        = glfw.KeyKP8
	KP9                        = glfw.KeyKP9
	KPDivide                   = glfw.KeyKPDivide
	KPMultiply                 = glfw.KeyKPMultiply
	KPSubtract                 = glfw.KeyKPSubtract
	KPAdd                      = glfw.KeyKPAdd
	KPDecimal                  = glfw.KeyKPDecimal
	KPEqual                    = glfw.KeyKPEqual
	KPEnter                    = glfw.KeyKPEnter
	KPNumlock                  = glfw.KeyKPNumlock
	Capslock                   = glfw.KeyCapslock
	Scrolllock                 = glfw.KeyScrolllock
	Pause                      = glfw.KeyPause
	Lsuper                     = glfw.KeyLsuper
	Rsuper                     = glfw.KeyRsuper
	Menu                       = glfw.KeyMenu
	Mouse1                     = glfw.Mouse1
	Mouse2                     = glfw.Mouse2
	Mouse3                     = glfw.Mouse3
	Mouse4                     = glfw.Mouse4
	Mouse5                     = glfw.Mouse5
	Mouse6                     = glfw.Mouse6
	Mouse7                     = glfw.Mouse7
	Mouse8                     = glfw.Mouse8
	MouseLeft                  = glfw.MouseLeft
	MouseRight                 = glfw.MouseRight
	MouseMiddle                = glfw.MouseMiddle
)

var (
	responder    Responder
	config       *Config
	timing       *stats
	defaultFont  *Font
	bgColor      *Color
	blankTexture *Texture
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

// A Responder describes an interface for application events.
//
// Open is called after the opengl context and window have been
// created. You should load assets and create eng objects in this method.
type Responder interface {
	Open()
	Close()
	Update(delta float32)
	Draw()
	MouseMove(x, y float32)
	MouseDown(x, y float32, button int)
	MouseUp(x, y float32, button int)
	MouseScroll(x, y float32, amount int)
	KeyType(key rune)
	KeyDown(key int)
	KeyUp(key int)
	Resize(width, height int)
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

	if err := glfw.Init(); err != nil {
		log.Fatal(err)
	}
	defer glfw.Terminate()

	if !config.Resizable {
		glfw.OpenWindowHint(glfw.WindowNoResize, 1)
	}
	glfw.OpenWindowHint(glfw.FsaaSamples, config.Fsaa)

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

	if err := gl.Init(); err != nil {
		log.Fatal(err)
	}

	config.Width, config.Height = glfw.WindowSize()

	bgColor = NewColor(0, 0, 0)

	glfw.SetWindowTitle(config.Title)

	if !config.Fullscreen {
		glfw.SetWindowPos((mode.W-width)/2, (mode.H-height)/2)
	}

	if config.Vsync {
		glfw.SetSwapInterval(1)
	}

	responder.Open()
	defer responder.Close()

	glfw.SetWindowSizeCallback(func(w, h int) {
		config.Width, config.Height = glfw.WindowSize()
		responder.Resize(w, h)
	})

	glfw.SetMousePosCallback(func(x, y int) {
		responder.MouseMove(float32(x), float32(y))
	})

	glfw.SetMouseButtonCallback(func(b, s int) {
		x, y := glfw.MousePos()
		if s == glfw.KeyPress {
			responder.MouseDown(float32(x), float32(y), b)
		} else {
			responder.MouseUp(float32(x), float32(y), b)
		}
	})

	var lastWheel int
	glfw.SetMouseWheelCallback(func(pos int) {
		if lastWheel-pos != 0 {
			x, y := glfw.MousePos()
			responder.MouseScroll(float32(x), float32(y), lastWheel-pos)
			lastWheel = pos
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
			responder.KeyType(rune(k))
		}
	})

	timing = NewStats(config.LogFPS)
	timing.Update()

	for glfw.WindowParam(glfw.Opened) == 1 {
		responder.Update(float32(timing.Dt))
		gl.ClearColor(gl.Float(bgColor.R), gl.Float(bgColor.G), gl.Float(bgColor.B), gl.Float(bgColor.A))
		gl.Clear(gl.COLOR_BUFFER_BIT)
		responder.Draw()
		glfw.SwapBuffers()
		timing.Update()
	}
}

func Log(l ...interface{}) {
	log.Println(l...)
}

// Clear manually clears with a given color. Mostly used with a Canvas.
func Clear(color *Color) {
	gl.ClearColor(gl.Float(color.R), gl.Float(color.G), gl.Float(color.B), gl.Float(color.A))
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
	glfw.SetWindowSize(w, h)
}

// Focused indicates if the game window is currently focused.
func Focused() bool {
	return glfw.WindowParam(glfw.Active) == gl.TRUE
}

// Exit closes the window and breaks out of the game loop.
func Exit() {
	glfw.CloseWindow()
}

// MouseX returns the cursor's horizontal position within the window.
func MouseX() int {
	x, _ := glfw.MousePos()
	return x
}

// MouseY returns the cursor's vetical position within the window.
func MouseY() int {
	_, y := glfw.MousePos()
	return y
}

// MousePos returns the cursor's X and Y position within the window.
func MousePos() (int, int) {
	return glfw.MousePos()
}

// SetMousePos sets the cursor's X and Y position within the window.
func SetMousePos(x, y int) {
	glfw.SetMousePos(x, y)
}

// SetMouseCursor shows or hides the cursor.
func SetMouseCursor(on bool) {
	if on {
		glfw.Enable(glfw.MouseCursor)
	} else {
		glfw.Disable(glfw.MouseCursor)
	}
}

// MousePressed takes a mouse button constant and indicates if it is
// currently pressed.
func MousePressed(b int) bool {
	return glfw.MouseButton(b) == glfw.KeyPress
}

// KeyPressed takes a key constant and indicates if it is currently pressed.
func KeyPressed(k int) bool {
	return glfw.Key(k) == glfw.KeyPress
}

// SetKeyRepeat toggles key repeat either on or off.
func SetKeyRepeat(repeat bool) {
	if repeat {
		glfw.Enable(glfw.KeyRepeat)
	} else {
		glfw.Disable(glfw.KeyRepeat)
	}
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

// Returns a white 1x1 pixel texture.
func BlankTexture() *Texture {
	if blankTexture == nil {
		img := image.NewRGBA(image.Rect(0, 0, 1, 1))
		draw.Draw(img, img.Bounds(), &image.Uniform{NewColor(1, 1, 1)}, image.ZP, draw.Src)
		blankTexture = NewTexture(img)
	}
	return blankTexture
}
