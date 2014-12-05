package engi

import (
// "log"
)

type Spritesheet struct {
	texture               *Texture
	CellWidth, CellHeight int
	cache                 map[int]*Region
}

func (s Spritesheet) Cell(i int) *Region {
	s.cache[i] = getRegionOfSpriteSheet(s.texture, 16, i)
	return s.cache[i]
}

func (s Spritesheet) Width() float32 {
	return s.texture.Width() / float32(s.CellWidth)
}

func (s Spritesheet) Height() float32 {
	return s.texture.Height() / float32(s.CellHeight)
}

func NewSpritesheet(filename string, cellsize int) *Spritesheet {
	return &Spritesheet{texture: Files.Image(filename), CellWidth: cellsize, CellHeight: cellsize, cache: make(map[int]*Region)}
}

type AnimationComponent struct {
	Index  int
	Rate   float32
	Change float32
	Tick   float32
	S      *Spritesheet
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

	ac.Change += dt
	if ac.Change >= ac.Rate {
		ac.Index += 1
		if ac.Index >= int(ac.S.Width()*ac.S.Height()) {
			ac.Index = 0
		}
		ac.Change = 0
		r.Display = ac.S.Cell(ac.Index)
	}
}
