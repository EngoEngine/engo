// Copyright 2014 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eng

type Responder interface {
	Load()
	Setup()
	Draw()
	Resize(width, height int)
	Update(delta float32)
	Mouse(x, y float32, action Action)
	Scroll(amount float32)
	Key(key Key, modifier Modifier, action Action)
	Type(char rune)
}

type Game struct{}

func (g *Game) Load()                             {}
func (g *Game) Setup()                            {}
func (g *Game) Draw()                             {}
func (g *Game) Resize(width, height int)          {}
func (g *Game) Update(dt float32)                 {}
func (g *Game) Mouse(x, y float32, action Action) {}
func (g *Game) Scroll(amount float32)             {}
func (g *Game) Type(char rune)                    {}
func (g *Game) Key(key Key, mod Modifier, act Action) {
	if key == Escape {
		Exit()
	}
}
