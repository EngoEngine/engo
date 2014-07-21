// Copyright 2014 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package engi

import "math"

type Displayer interface {
	Display(*Batch)
}

type Point struct {
	X, Y float32
}

func (p *Point) Set(x, y float32) {
	p.X = x
	p.Y = y
}

func (p *Point) SetTo(v float32) {
	p.X = v
	p.Y = v
}

type Entity struct {
	Position *Point
	Scale    *Point
	Pivot    *Point
	Rotation float32
	Tint     uint32
	Alpha    float32
}

func (e *Entity) floatBits() float32 {
	r := uint32((e.Tint >> 16) & 0xFF)
	g := uint32((e.Tint >> 8) & 0xFF)
	b := uint32(e.Tint & 0xFF)
	a := uint32(e.Alpha * 255.0)
	return math.Float32frombits((a<<24 | b<<16 | g<<8 | r) & 0xfeffffff)
}

func NewEntity(x, y float32) *Entity {
	return &Entity{
		Position: &Point{x, y},
		Scale:    &Point{1, 1},
		Pivot:    &Point{0.5, 0.5},
		Rotation: 0,
		Tint:     0xffffff,
		Alpha:    1,
	}
}

type Sprite struct {
	*Entity
	region *Region
}

func NewSprite(region *Region, x, y float32) *Sprite {
	return &Sprite{Entity: NewEntity(x, y), region: region}
}

func (s *Sprite) Display(batch *Batch) {
	batch.Draw(s.region, s.Position.X, s.Position.Y, s.Pivot.X, s.Pivot.Y, s.Scale.X, s.Scale.Y, s.Rotation, s.floatBits())
}

type Text struct {
	*Entity
	font  *Font
	text  string
	width float32
}

func NewText(font *Font, x, y float32, text string) *Text {
	return &Text{Entity: NewEntity(x, y), font: font, text: text}
}

func (t *Text) Display(batch *Batch) {
	px := t.Position.X - t.width*t.Pivot.X
	py := t.Position.Y
	sx := t.Scale.X
	sy := t.Scale.Y
	for _, v := range t.text {
		i, ok := t.font.mapRune(v)
		if ok {
			region := t.font.regions[i]
			offset := t.font.offsets[i]
			x := px + offset.xoffset*sx
			y := py + offset.yoffset*sy - (region.height * sy * t.Pivot.Y)
			if t.width != 0 {
				batch.Draw(region, x, y, 0, 0, sx, sy, 0, t.floatBits())
			}
			px += offset.xadvance * sx
		}
	}
	if t.width == 0 {
		t.width = px - t.Position.X
	}
}

type Stage struct {
	batch   *Batch
	objects []Displayer
}

func NewStage() *Stage {
	return new(Stage)
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

func (s *Stage) Add(object Displayer) {
	s.objects = append(s.objects, object)
}

func (s *Stage) Sprite(name string, x, y float32) *Sprite {
	texture := NewTexture(Files.Image(name))
	region := NewRegion(texture, 0, 0, texture.Width(), texture.Height())
	sprite := NewSprite(region, x, y)
	s.Add(sprite)
	return sprite
}

func (s *Stage) Text(font *Font, x, y float32, content string) *Text {
	text := NewText(font, x, y, content)
	s.Add(text)
	return text
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
	s.objects = make([]Displayer, 0)
}
func (s *Stage) Setup() {}
func (s *Stage) draw() {
	s.batch.Begin()
	for _, object := range s.objects {
		object.Display(s.batch)
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
