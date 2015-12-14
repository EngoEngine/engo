// Copyright 2014 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package engi

import (
	"fmt"

	"github.com/paked/engi/ecs"
	"github.com/paked/webgl"
)

var (
	Time        *Clock
	Files       *Loader
	Gl          *webgl.Context
	WorldBounds AABB

	currentWorld *ecs.World
	currentScene Scene
	Mailbox      *MessageManager
	cam          *cameraSystem

	scaleOnResize   = false
	fpsLimit        = 120
	headless        = false
	resetLoopTicker = make(chan bool, 1)
)

func Open(title string, width, height int, fullscreen bool, defaultScene Scene) {
	keyStates = make(map[Key]bool)
	Time = NewClock()
	Files = NewLoader()

	run(defaultScene, title, width, height, fullscreen)
}

func OpenHeadless(defaultScene Scene) {
	keyStates = make(map[Key]bool)
	Time = NewClock()
	Files = NewLoader() // TODO: do we want files in Headless mode?

	headless = true

	runHeadless(defaultScene)
}

func OpenHeadlessNoRun() {
	Time = NewClock()
	Files = NewLoader() // TODO: do we want files in Headless mode?

	headless = true
}

func SetBg(color uint32) {
	if !headless {
		r := float32((color>>16)&0xFF) / 255.0
		g := float32((color>>8)&0xFF) / 255.0
		b := float32(color&0xFF) / 255.0
		Gl.ClearColor(r, g, b, 1.0)
	}
}

func SetScaleOnResize(b bool) {
	scaleOnResize = b
}

func SetFPSLimit(limit int) error {
	if limit <= 0 {
		return fmt.Errorf("FPS Limit out of bounds. Requires > 0")
	}
	fpsLimit = limit
	resetLoopTicker <- true
	return nil
}
