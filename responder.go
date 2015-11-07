// Copyright 2014 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package engi

import (
	"testing"
)

type CustomGame interface {
	Preload()
	Setup(*World)
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

type inlineGame struct {
	preloadFunc func()
	setupFunc   func(*World)
}

func (m *inlineGame) Preload() {
	m.preloadFunc()
}

func (m *inlineGame) Setup(w *World) {
	m.setupFunc(w)
}

// NewGame allows you to create a `Responder` using two inline functions `preload` and `setup`.
func NewGame(preload func(), setup func(*World)) CustomGame {
	g := &inlineGame{preloadFunc: preload, setupFunc: setup}
	return g
}

// Bench is a helper-function to easily benchmark one frame, given a preload / setup function
func Bench(b *testing.B, preload func(), setup func(w *World)) {
	g := NewGame(preload, setup)

	OpenHeadlessNoRun()
	RunPreparation(g)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RunIteration()
	}
}
