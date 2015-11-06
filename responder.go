// Copyright 2014 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package engi

import (
	"testing"
)

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
	AddEntity(e *Entity)
	Batch(PriorityLevel) *Batch
	New()
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
	World
	preloadFunc func(*World)
	setupFunc   func(*World)
}

func (m *inlineGame) Preload() {
	m.preloadFunc(&m.World)
}

func (m *inlineGame) Setup() {
	m.setupFunc(&m.World)
}

// NewGame allows you to create a `Responder` using two inline functions `preload` and `setup`.
func NewGame(preload, setup func(*World)) Responder {
	g := &inlineGame{preloadFunc: preload, setupFunc: setup}
	return g
}

// Bench is a helper-function to easily benchmark one frame, given a preload / setup function
func Bench(b *testing.B, preload, setup func(w *World)) {
	g := NewGame(preload, setup)

	OpenHeadlessNoRun(g)
	RunPreparation()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RunIteration()
	}
}
