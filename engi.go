// Copyright 2014 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package engi

import (
	"log"
	"math"
	"time"
)

var (
	responder Responder
	width     float32
	height    float32
	timing    *stats
	Files     *Loader
	GL        *gl2
)

type Responder interface {
	Render()
	Resize(width, height int)
	Preload()
	Setup()
	Update(dt float32)
	Mouse(x, y float32, action Action)
	Scroll(amount float32)
	Key(key Key, modifier Modifier, action Action)
	Type(char rune)
}

type Game struct{}

func (g *Game) Preload()                          {}
func (g *Game) Setup()                            {}
func (g *Game) Update(dt float32)                 {}
func (g *Game) Render()                           {}
func (g *Game) Resize(w, h int)                   {}
func (g *Game) Mouse(x, y float32, action Action) {}
func (g *Game) Scroll(amount float32)             {}
func (g *Game) Key(key Key, modifier Modifier, action Action) {
	if key == Escape {
		Exit()
	}
}
func (g *Game) Type(char rune) {}

type stats struct {
	Elapsed, Dt, Fps, Frames, Period float64
	Start                            time.Time
	Then                             time.Time
	show                             bool
}

func NewStats(show bool) *stats {
	st := new(stats)
	st.Start = time.Now()
	st.Period = 1
	st.Update()
	st.show = show
	return st
}

func (t *stats) Update() {
	now := time.Now()
	t.Frames += 1
	t.Dt = now.Sub(t.Then).Seconds()
	t.Elapsed += t.Dt
	t.Then = now

	if t.Elapsed >= t.Period {
		t.Fps = t.Frames / t.Period
		t.Elapsed = math.Mod(t.Elapsed, t.Period)
		t.Frames = 0
		if t.show {
			log.Println(t.Fps)
		}
	}
}

func Open(title string, width, height int, fullscreen bool, r Responder) {
	responder = r
	Files = NewLoader()
	run(title, width, height, fullscreen)
}

func SetBg(color uint32) {
	r := float32((color>>16)&0xFF) / 255.0
	g := float32((color>>8)&0xFF) / 255.0
	b := float32(color&0xFF) / 255.0
	GL.ClearColor(r, g, b, 1.0)
}

func Width() float32 {
	return float32(width)
}

func Height() float32 {
	return float32(height)
}

func Delta() float32 {
	return float32(timing.Dt)
}

func Fps() float32 {
	return float32(timing.Fps)
}

func Exit() {
	exit()
}
