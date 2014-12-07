package engi

import (
	"log"
)

type Spritesheet struct {
	texture               *Texture
	CellWidth, CellHeight int
	cache                 map[int]*Region
}

func (s Spritesheet) Cell(i int) *Region {
	if r := s.cache[i]; r != nil {
		return r
	}
	s.cache[i] = getRegionOfSpriteSheet(s.texture, s.CellWidth, i)
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

func NewAnimationComponent() *AnimationComponent {
	return &AnimationComponent{Animations: make(map[string][]int)}
}

type AnimationComponent struct {
	Index            int
	_index           int
	Rate             float32
	Change           float32
	S                *Spritesheet
	Animations       map[string][]int
	CurrentAnimation []int
}

func (ac *AnimationComponent) SelectAnimation(name string) {
	ac.CurrentAnimation = ac.Animations[name]
}

func (ac *AnimationComponent) AddAnimation(name string, indexes []int) {
	ac.Animations[name] = indexes
}

func (ac *AnimationComponent) Increment() {
	if len(ac.CurrentAnimation) == 0 {
		log.Println("No data for this animation")
		return
	}

	// log.Println("Incrementing")
	ac.Index += 1
	if ac.Index >= len(ac.CurrentAnimation) {
		ac.Index = 0
	}
	ac._index = ac.CurrentAnimation[ac.Index]
	ac.Change = 0
}

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

	ac.Change += dt
	if ac.Change >= ac.Rate {
		ac.Increment()
		r.Display = ac.Cell()
	}
}
