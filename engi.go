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
	Cam         *Camera
	Wo          Responder
	WorldBounds AABB

	fpsLimit        = 120
	resetLoopTicker = make(chan bool, 1)
)

func Open(title string, width, height int, fullscreen bool, r Responder) {
	states = make(map[Key]bool)
	responder = r
	Time = NewClock()
	Files = NewLoader()
	SetCamera(&Camera{})
	Wo = r
	run(title, width, height, fullscreen)
}

func SetCamera(c *Camera) {
	Cam = c
}

func SetBg(color uint32) {
	r := float32((color>>16)&0xFF) / 255.0
	g := float32((color>>8)&0xFF) / 255.0
	b := float32(color&0xFF) / 255.0
	Gl.ClearColor(r, g, b, 1.0)
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
