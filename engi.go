// Copyright 2014 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package engi

var (
	responder Responder
	width     float32
	height    float32
	Time      *Clock
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

func Open(title string, width, height int, fullscreen bool, r Responder) {
	responder = r
	Time = NewClock()
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

func Exit() {
	exit()
}
