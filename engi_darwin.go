// Copyright 2014 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !windows,!linux,!netgo,!android

package engi

import (
	glfw "github.com/go-gl/glfw3"
	"github.com/go-gl/glow/gl/2.1/gl"
	"log"
	"runtime"
)

var window *glfw.Window

func run() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

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

	if !config.Fullscreen {
		window.SetPosition((mode.Width-width)/2, (mode.Height-height)/2)
	}

	if config.Vsync {
		glfw.SwapInterval(1)
	}

	if err := gl.Init(); err != nil {
		log.Fatal(err)
	}
	GL = newgl2()

	GL.Viewport(0, 0, config.Width, config.Height)

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

	responder.Preload()

	timing = NewStats(config.LogFPS)
	timing.Update()

	Files.Load(func() {})

	responder.Setup()

	for !window.ShouldClose() {
		responder.Update(float32(timing.Dt))
		GL.Clear(gl.COLOR_BUFFER_BIT)
		responder.Render()
		window.SwapBuffers()
		glfw.PollEvents()
		timing.Update()
	}
}

func exit() {
	window.SetShouldClose(true)
}

const (
	Escape       = Key(glfw.KeyEscape)
	F1           = Key(glfw.KeyF1)
	F2           = Key(glfw.KeyF2)
	F3           = Key(glfw.KeyF3)
	F4           = Key(glfw.KeyF4)
	F5           = Key(glfw.KeyF5)
	F6           = Key(glfw.KeyF6)
	F7           = Key(glfw.KeyF7)
	F8           = Key(glfw.KeyF8)
	F9           = Key(glfw.KeyF9)
	F10          = Key(glfw.KeyF10)
	F11          = Key(glfw.KeyF11)
	F12          = Key(glfw.KeyF12)
	Up           = Key(glfw.KeyUp)
	Down         = Key(glfw.KeyDown)
	Left         = Key(glfw.KeyLeft)
	Right        = Key(glfw.KeyRight)
	LeftShift    = Key(glfw.KeyLeftShift)
	RightShift   = Key(glfw.KeyRightShift)
	LeftControl  = Key(glfw.KeyLeftControl)
	RightControl = Key(glfw.KeyRightControl)
	LeftAlt      = Key(glfw.KeyLeftAlt)
	RightAlt     = Key(glfw.KeyRightAlt)
	Tab          = Key(glfw.KeyTab)
	Space        = Key(glfw.KeySpace)
	Enter        = Key(glfw.KeyEnter)
	Backspace    = Key(glfw.KeyBackspace)
	Insert       = Key(glfw.KeyInsert)
	Delete       = Key(glfw.KeyDelete)
	PageUp       = Key(glfw.KeyPageUp)
	PageDown     = Key(glfw.KeyPageDown)
	Home         = Key(glfw.KeyHome)
	End          = Key(glfw.KeyEnd)
	Kp0          = Key(glfw.KeyKp0)
	Kp1          = Key(glfw.KeyKp1)
	Kp2          = Key(glfw.KeyKp2)
	Kp3          = Key(glfw.KeyKp3)
	Kp4          = Key(glfw.KeyKp4)
	Kp5          = Key(glfw.KeyKp5)
	Kp6          = Key(glfw.KeyKp6)
	Kp7          = Key(glfw.KeyKp7)
	Kp8          = Key(glfw.KeyKp8)
	Kp9          = Key(glfw.KeyKp9)
	KpDivide     = Key(glfw.KeyKpDivide)
	KpMultiply   = Key(glfw.KeyKpMultiply)
	KpSubtract   = Key(glfw.KeyKpSubtract)
	KpAdd        = Key(glfw.KeyKpAdd)
	KpDecimal    = Key(glfw.KeyKpDecimal)
	KpEqual      = Key(glfw.KeyKpEqual)
	KpEnter      = Key(glfw.KeyKpEnter)
	NumLock      = Key(glfw.KeyNumLock)
	CapsLock     = Key(glfw.KeyCapsLock)
	ScrollLock   = Key(glfw.KeyScrollLock)
	Pause        = Key(glfw.KeyPause)
	LeftSuper    = Key(glfw.KeyLeftSuper)
	RightSuper   = Key(glfw.KeyRightSuper)
	Menu         = Key(glfw.KeyMenu)
)
