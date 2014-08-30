// Copyright 2014 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package engi

type Responder interface {
	Render()
	Resize(width, height int)
	Preload()
	Setup()
	Close()
	Update(dt float32)
	Mouse(x, y float32, action Action)
	Scroll(amount float32)
	Key(key Key, modifier Modifier, action Action)
	Type(char rune)
}

type Game struct{}

func (g *Game) Preload()                          {}
func (g *Game) Setup()                            {}
func (g *Game) Close()                            {}
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
