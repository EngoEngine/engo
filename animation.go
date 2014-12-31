// Copyright 2014 Harrison Shoebridge. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package engi

import (
	"log"
)

// Spritesheet is a class that stores a set of tiles from a file, used by tilemaps and animations
type Spritesheet struct {
	texture               *Texture        // The original texture
	CellWidth, CellHeight int             // The dimensions of the cells
	cache                 map[int]*Region // The cell cache cells
}

// Simple handler for creating a new spritesheet from a file
// textureName is the name of a texture already preloaded with engi.Files.Add
func NewSpritesheet(textureName string, cellWidth, cellHeight int) *Spritesheet {
	return &Spritesheet{texture: Files.Image(textureName), CellWidth: cellWidth, CellHeight: cellHeight, cache: make(map[int]*Region)}
}

// Get the region at the index i, updates and pulls from cache if need be
func (s *Spritesheet) Cell(i int) *Region {
	if r := s.cache[i]; r != nil {
		return r
	}
	s.cache[i] = getRegionOfSpriteSheet(s.texture, s.CellWidth, i)
	return s.cache[i]
}

// The amount of tiles on the x-axis of the spritesheet
func (s Spritesheet) Width() float32 {
	return s.texture.Width() / float32(s.CellWidth)
}

// The amount of tiles on the y-axis of the spritesheet
func (s Spritesheet) Height() float32 {
	return s.texture.Height() / float32(s.CellHeight)
}

// Component that controls animation in rendering entities
type AnimationComponent struct {
	Index            int              // What frame in the is being used
	_index           int              // The index of the tile that should currently be being displayed
	Rate             float32          // How often frames should increment, in seconds.
	change           float32          // The time since the last incrementation
	S                *Spritesheet     // Pointer to the source spritesheet
	Animations       map[string][]int // All possible animations
	CurrentAnimation []int            // The currently animation
}

// Create a new pointer to AnimationComponent
func NewAnimationComponent() *AnimationComponent {
	return &AnimationComponent{Animations: make(map[string][]int)}
}

func (ac *AnimationComponent) SelectAnimation(name string) {
	ac.CurrentAnimation = ac.Animations[name]
}

func (ac *AnimationComponent) AddAnimation(name string, indexes []int) {
	ac.Animations[name] = indexes
}

// Increment the frame
func (ac *AnimationComponent) Increment() {
	if len(ac.CurrentAnimation) == 0 {
		log.Println("No data for this animation")
		return
	}

	ac.Index += 1
	if ac.Index >= len(ac.CurrentAnimation) {
		ac.Index = 0
	}
	ac._index = ac.CurrentAnimation[ac.Index]
	ac.change = 0

}

//Bug(me) Don't need to use _index at all, use ac.CurrentAnimation[ac.Index] instead. Set ac.Index as a private member variable
func (ac *AnimationComponent) Cell() *Region {
	return ac.S.Cell(ac._index)
}

func (ac AnimationComponent) Name() string {
	return "AnimationComponent"
}

type AnimationSystem struct {
	*System
}

func (a *AnimationSystem) New() {
	a.System = &System{}
}

func (a AnimationSystem) Name() string {
	return "AnimationSystem"
}

func (a *AnimationSystem) Update(e *Entity, dt float32) {
	var (
		ac *AnimationComponent
		r  *RenderComponent
	)

	if !e.GetComponent(&ac) || !e.GetComponent(&r) {
		return
	}

	ac.change += dt
	if ac.change >= ac.Rate {
		ac.Increment()
		r.Display = ac.Cell()
	}
}
