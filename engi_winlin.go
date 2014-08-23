// Copyright 2014 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !darwin,!netgo,!android

package engi

import (
	"azul3d.org/chippy.v1"
	"azul3d.org/keyboard.v1"
	"azul3d.org/mouse.v1"
	"github.com/go-gl/glow/gl/2.1/gl"
	"log"
	"os"
	"runtime"
)

var window *chippy.Window

func program() {
	defer chippy.Exit()

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	screen := chippy.DefaultScreen()
	window = chippy.NewWindow()
	window.SetTitle(config.Title)
	window.SetSize(config.Width, config.Height)
	window.SetPositionCenter(screen)

	if config.Fullscreen {
		window.SetFullscreen(true)
		mode := screen.Mode()
		config.Width, config.Height = mode.Resolution()
	}

	err := window.Open(screen)
	if err != nil {
		log.Fatal(err)
	}

	configs := window.GLConfigs()
	bestConfig := chippy.GLChooseConfig(configs, chippy.GLWorstConfig, chippy.GLBestConfig)
	window.GLSetConfig(bestConfig)

	context, err := window.GLCreateContext(2, 1, chippy.GLCoreProfile, nil)
	if err != nil {
		log.Fatal(err)
	}
	window.GLMakeCurrent(context)

	if err := gl.Init(); err != nil {
		log.Fatal(err)
	}
	GL = newgl2()

	GL.Viewport(0, 0, config.Width, config.Height)

	responder.Preload()
	timing = NewStats(config.LogFPS)
	timing.Update()

	Files.Load(func() {})

	responder.Setup()

	events := window.Events()
	defer window.CloseEvents(events)

	for {
		responder.Update(float32(timing.Dt))
		GL.Clear(gl.COLOR_BUFFER_BIT)
		responder.Render()
		window.GLSwapBuffers()
		timing.Update()

		for i := 0; i < len(events); i++ {
			e := <-events
			switch ev := e.(type) {
			case chippy.ResizedEvent:
				config.Width = ev.Width
				config.Height = ev.Height
				GL.Viewport(0, 0, config.Width, config.Height)
				responder.Resize(config.Width, config.Height)
			case keyboard.StateEvent:
				switch ev.State {
				case keyboard.Up:
					responder.Key(Key(ev.Key), Modifier(0), RELEASE)
				case keyboard.Down:
					responder.Key(Key(ev.Key), Modifier(0), PRESS)
				}
			case keyboard.TypedEvent:
				responder.Type(ev.Rune)
			case mouse.Event:
				switch ev.State {
				case mouse.Up:
					responder.Mouse(float32(0), float32(0), RELEASE)
				case mouse.Down:
					responder.Mouse(float32(0), float32(0), PRESS)
				case mouse.ScrollForward:
					responder.Scroll(float32(1))
				case mouse.ScrollBack:
					responder.Scroll(float32(-1))
				}
			case chippy.CloseEvent, chippy.DestroyedEvent:
				return
			default:
			}
		}
	}
}

func run() {
	chippy.SetDebugOutput(os.Stdout)

	err := chippy.Init()
	if err != nil {
		log.Fatal(err)
	}

	go program()
	chippy.MainLoop()
}

func exit() {
	window.Destroy()
}

const (
	Dash         = Key(keyboard.Dash)
	Apostrophe   = Key(keyboard.Apostrophe)
	Semicolon    = Key(keyboard.Semicolon)
	Equals       = Key(keyboard.Equals)
	Comma        = Key(keyboard.Comma)
	Period       = Key(keyboard.Period)
	Slash        = Key(keyboard.ForwardSlash)
	Backslash    = Key(keyboard.BackSlash)
	Backspace    = Key(keyboard.Backspace)
	Tab          = Key(keyboard.Tab)
	CapsLock     = Key(keyboard.CapsLock)
	Space        = Key(keyboard.Space)
	Enter        = Key(keyboard.Enter)
	Escape       = Key(keyboard.Escape)
	Insert       = Key(keyboard.Insert)
	PrintScreen  = Key(keyboard.PrintScreen)
	Delete       = Key(keyboard.Delete)
	PageUp       = Key(keyboard.PageUp)
	PageDown     = Key(keyboard.PageDown)
	Home         = Key(keyboard.Home)
	End          = Key(keyboard.End)
	Pause        = Key(keyboard.Pause)
	ScrollLock   = Key(keyboard.ScrollLock)
	ArrowLeft    = Key(keyboard.ArrowLeft)
	ArrowRight   = Key(keyboard.ArrowRight)
	ArrowDown    = Key(keyboard.ArrowDown)
	ArrowUp      = Key(keyboard.ArrowUp)
	LeftBracket  = Key(keyboard.LeftBracket)
	LeftShift    = Key(keyboard.LeftShift)
	LeftControl  = Key(keyboard.LeftCtrl)
	LeftSuper    = Key(keyboard.LeftSuper)
	LeftAlt      = Key(keyboard.LeftAlt)
	RightBracket = Key(keyboard.RightBracket)
	RightShift   = Key(keyboard.RightShift)
	RightControl = Key(keyboard.RightCtrl)
	RightSuper   = Key(keyboard.RightSuper)
	RightAlt     = Key(keyboard.RightAlt)
	Zero         = Key(keyboard.Zero)
	One          = Key(keyboard.One)
	Two          = Key(keyboard.Two)
	Three        = Key(keyboard.Three)
	Four         = Key(keyboard.Four)
	Five         = Key(keyboard.Five)
	Six          = Key(keyboard.Six)
	Seven        = Key(keyboard.Seven)
	Eight        = Key(keyboard.Eight)
	Nine         = Key(keyboard.Nine)
	F1           = Key(keyboard.F1)
	F2           = Key(keyboard.F2)
	F3           = Key(keyboard.F3)
	F4           = Key(keyboard.F4)
	F5           = Key(keyboard.F5)
	F6           = Key(keyboard.F6)
	F7           = Key(keyboard.F7)
	F8           = Key(keyboard.F8)
	F9           = Key(keyboard.F9)
	F10          = Key(keyboard.F10)
	F11          = Key(keyboard.F11)
	F12          = Key(keyboard.F12)
	A            = Key(keyboard.A)
	B            = Key(keyboard.B)
	C            = Key(keyboard.C)
	D            = Key(keyboard.D)
	E            = Key(keyboard.E)
	F            = Key(keyboard.F)
	G            = Key(keyboard.G)
	H            = Key(keyboard.H)
	I            = Key(keyboard.I)
	J            = Key(keyboard.J)
	K            = Key(keyboard.K)
	L            = Key(keyboard.L)
	M            = Key(keyboard.M)
	N            = Key(keyboard.N)
	O            = Key(keyboard.O)
	P            = Key(keyboard.P)
	Q            = Key(keyboard.Q)
	R            = Key(keyboard.R)
	S            = Key(keyboard.S)
	T            = Key(keyboard.T)
	U            = Key(keyboard.U)
	V            = Key(keyboard.V)
	W            = Key(keyboard.W)
	X            = Key(keyboard.X)
	Y            = Key(keyboard.Y)
	Z            = Key(keyboard.Z)
	NumLock      = Key(keyboard.NumLock)
	NumMultiply  = Key(keyboard.NumMultiply)
	NumDivide    = Key(keyboard.NumDivide)
	NumAdd       = Key(keyboard.NumAdd)
	NumSubtract  = Key(keyboard.NumSubtract)
	NumZero      = Key(keyboard.NumZero)
	NumOne       = Key(keyboard.NumOne)
	NumTwo       = Key(keyboard.NumTwo)
	NumThree     = Key(keyboard.NumThree)
	NumFour      = Key(keyboard.NumFour)
	NumFive      = Key(keyboard.NumFive)
	NumSix       = Key(keyboard.NumSix)
	NumSeven     = Key(keyboard.NumSeven)
	NumEight     = Key(keyboard.NumEight)
	NumNine      = Key(keyboard.NumNine)
	NumDecimal   = Key(keyboard.NumDecimal)
	NumEnter     = Key(keyboard.NumEnter)
)
