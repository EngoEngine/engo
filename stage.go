// Copyright 2014 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package engi

const whiteBits = -1.7014117331926443e+38

type Stage struct {
	*Sprite
	batch *Batch
}

func NewStage() *Stage {
	return &Stage{NewGroup(), nil}
}

func (s *Stage) SetBg(color uint32) {
	r := float32((color>>16)&0xFF) / 255.0
	g := float32((color>>8)&0xFF) / 255.0
	b := float32(color&0xFF) / 255.0
	GL.ClearColor(r, g, b, 1.0)
}

func (s *Stage) Load(name, path string) {
	Files.Add(name, path)
}

func (s *Stage) NewSprite(region *Region, x, y float32) *Sprite {
	return s.AddChild(NewSprite(region, x, y))
}

// Width returns the current window width.
func (s *Stage) Width() float32 {
	return float32(config.Width)
}

// Width returns the current window width.
func (s *Stage) Height() float32 {
	return float32(config.Height)
}

func (s *Stage) Delta() float32 {
	return float32(timing.Dt)
}

func (s *Stage) Time() float32 {
	return float32(timing.Start.Sub(timing.Then).Seconds())
}

func (s *Stage) Fps() float32 {
	return float32(timing.Fps)
}

func (s *Stage) Exit() {
	exit()
}

func (s *Stage) Preload() {}
func (s *Stage) init() {
	s.batch = NewBatch(s.Width(), s.Height())
}
func (s *Stage) Setup() {}
func (s *Stage) draw() {
	s.batch.Begin()
	for _, child := range s.children {
		child.render(s.batch)
	}
	s.batch.End()
}
func (s *Stage) resize(width, height int) {
	s.batch.SetProjection(s.Width()/2, s.Height()/2)
}
func (s *Stage) Update()                           {}
func (s *Stage) Mouse(x, y float32, action Action) {}
func (s *Stage) Scroll(amount float32)             {}
func (s *Stage) Type(char rune)                    {}
func (s *Stage) Key(key Key, mod Modifier, act Action) {
	if key == Escape {
		s.Exit()
	}
}

/*
type alignment float32

const (
	LEFT   = alignment(0)
	CENTER = alignment(0.5)
	RIGHT  = alignment(1)
)

type Text struct {
	*Sprite
	font  *Font
	text  string
	width float32
	align int
}

func NewText(font *Font, x, y float32, text string) *Text {
	container := NewGroup()
	container.Position.Set(x, y)
	t := &Text{Sprite: container, font: font, text: text}
	t.Build()
	return t
}

func (t *Text) SetTint(tint uint32) {
	for _, child := range t.children {
		child.SetTint(tint)
	}
	t.Sprite.SetTint(tint)
}

func (t *Text) SetAlign(align alignment) {
	t.align = align
	t.Build()
}

func (t *Text) Build() {
	t.width = 0
	for _, v := range t.text {
		i, ok := t.font.mapRune(v)
		if ok {
			offset := t.font.offsets[i]
			t.width += offset.xadvance * t.Scale.X
		}
	}

	px := -(t.width * t.align)
	py := -(float32(0)
	for _, v := range t.text {
		i, exists := t.font.mapRune(v)
		if exists {
			region := t.font.regions[i]
			offset := t.font.offsets[i]
			x := px + offset.xoffset
			y := py + offset.yoffset - (region.height * t.Anchor.Y)
			sprite := NewSprite(region, x, y)
			sprite.color = t.color
			t.AddChild(sprite)
			px += offset.xadvance * sx
		}
	}
}
*/
