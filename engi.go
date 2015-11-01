// Copyright 2014 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package engi

import (
	"fmt"
	"github.com/paked/webgl"
)

var (
	responder   Responder
	Time        *Clock
	Files       *Loader
	Gl          *webgl.Context
	Mailbox     MessageManager
	cam         *cameraSystem
	Wo          Responder
	WorldBounds AABB

	fpsLimit        = 120
	headless        bool
	resetLoopTicker = make(chan bool, 1)
)

func Open(title string, width, height int, fullscreen bool, r Responder) {
	keyStates = make(map[Key]bool)
	responder = r
	Time = NewClock()
	Files = NewLoader()
	Wo = r
	run(title, width, height, fullscreen)
}

func OpenHeadless(r Responder) {
	Time = NewClock()
	Files = NewLoader() // TODO: do we want files in Headless mode?

	// TODO: change these (#35)
	responder = r
	Wo = r
	headless = true

	runHeadless()
}

func SetBg(color uint32) {
	if !headless {
		r := float32((color>>16)&0xFF) / 255.0
		g := float32((color>>8)&0xFF) / 255.0
		b := float32(color&0xFF) / 255.0
		Gl.ClearColor(r, g, b, 1.0)
	}
}

func SetFPSLimit(limit int) error {
	if limit <= 0 {
		return fmt.Errorf("FPS Limit out of bounds. Requires > 0")
	}
	fpsLimit = limit
	resetLoopTicker <- true
	return nil
}

func Width() float32 {
	return width()
}

func Height() float32 {
	return height()
}

func Exit() {
	exit()
}
