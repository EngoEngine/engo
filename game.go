// Copyright 2013 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eng

// A Game provides a default implementation for a responder so you can
// simply override those methods you might want to respond to.
type Game struct {
}

func (g *Game) Open() {
}

func (g *Game) Close() {
}

func (g *Game) Update(dt float32) {
}

func (g *Game) Draw() {
}

func (g *Game) MouseMove(x, y int) {
}

func (g *Game) MouseDown(x, y int, b Button) {
}

func (g *Game) MouseUp(x, y int, b Button) {
}

func (g *Game) MouseScroll(x, y, a int) {
}

func (g *Game) KeyType(k rune) {
}

func (g *Game) KeyDown(k Key) {
}

func (g *Game) KeyUp(k Key) {
}

func (g *Game) Resize(w, h int) {
}
