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

func NewSpritesheetFromTexture(texture *Texture, cellWidth, cellHeight int) *Spritesheet {
	return &Spritesheet{texture: texture, CellWidth: cellWidth, CellHeight: cellHeight, cache: make(map[int]*Region)}
}

// Simple handler for creating a new spritesheet from a file
// textureName is the name of a texture already preloaded with engi.Files.Add
func NewSpritesheetFromFile(textureName string, cellWidth, cellHeight int) *Spritesheet {
	return NewSpritesheetFromTexture(Files.Image(textureName), cellWidth, cellHeight)
}

// Get the region at the index i, updates and pulls from cache if need be
func (s *Spritesheet) Cell(index int) *Region {
	if r := s.cache[index]; r != nil {
		return r
	}
	s.cache[index] = regionFromSheet(s.texture, s.CellWidth, s.CellHeight, index)

	return s.cache[index]
}

func (s *Spritesheet) Renderable(index int) Renderable {
	return s.Cell(index)
}

func (s *Spritesheet) Renderables() []Renderable {
	renderables := make([]Renderable, s.CellCount())

	for i := 0; i < s.CellCount(); i++ {
		renderables[i] = s.Renderable(i)
	}

	return renderables
}

func (s *Spritesheet) CellCount() int {
	return int(s.Width()) * int(s.Height())
}

func (s *Spritesheet) Cells() []*Region {
	cellsNo := s.CellCount()
	cells := make([]*Region, cellsNo)
	for i := 0; i < cellsNo; i++ {
		cells[i] = s.Cell(i)
	}

	return cells
}

// The amount of tiles on the x-axis of the spritesheet
func (s Spritesheet) Width() float32 {
	return s.texture.Width() / float32(s.CellWidth)
}

// The amount of tiles on the y-axis of the spritesheet
func (s Spritesheet) Height() float32 {
	return s.texture.Height() / float32(s.CellHeight)
}

type AnimationAction struct {
	Name   string
	Frames []int
}

// Component that controls animation in rendering entities
type AnimationComponent struct {
	index            int              // What frame in the is being used
	Rate             float32          // How often frames should increment, in seconds.
	change           float32          // The time since the last incrementation
	Renderables      []Renderable     // Renderables
	Animations       map[string][]int // All possible animations
	CurrentAnimation []int            // The current animation
}

// Create a new pointer to AnimationComponent
func NewAnimationComponent(renderables []Renderable, rate float32) *AnimationComponent {
	return &AnimationComponent{
		Animations:  make(map[string][]int),
		Renderables: renderables,
		Rate:        rate,
	}
}

func (ac *AnimationComponent) SelectAnimationByName(name string) {
	ac.CurrentAnimation = ac.Animations[name]
}

func (ac *AnimationComponent) SelectAnimationByAction(action *AnimationAction) {
	ac.CurrentAnimation = ac.Animations[action.Name]
}

func (ac *AnimationComponent) AddAnimationAction(action *AnimationAction) {
	ac.Animations[action.Name] = action.Frames
}

func (ac *AnimationComponent) AddAnimationActions(actions []*AnimationAction) {
	for _, action := range actions {
		ac.Animations[action.Name] = action.Frames
	}
}

func (ac *AnimationComponent) Cell() Renderable {
	idx := ac.CurrentAnimation[ac.index]

	return ac.Renderables[idx]
}

func (*AnimationComponent) Type() string {
	return "AnimationComponent"
}

type AnimationSystem struct {
	*System
}

func (a *AnimationSystem) New() {
	a.System = &System{}
}

func (AnimationSystem) Type() string {
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
		a.NextFrame(ac)
		r.Display = ac.Cell()
	}
}

func (a *AnimationSystem) NextFrame(ac *AnimationComponent) {
	if len(ac.CurrentAnimation) == 0 {
		log.Println("No data for this animation")
		return
	}

	ac.index += 1
	if ac.index >= len(ac.CurrentAnimation) {
		ac.index = 0
	}
	ac.change = 0
}
