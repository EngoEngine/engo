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
				// Not triggering
				println(1)
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
	chippy.Exit()
}

const Escape = Key(keyboard.Escape)
