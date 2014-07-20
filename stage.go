// Copyright 2014 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eng

var (
	whiteFloatBits = NewColor(255, 255, 255, 1).FloatBits()
)

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
	color    float32
}

func (e *Entity) SetColor(color *Color) {
	e.color = color.FloatBits()
}

func NewEntity(x, y float32) *Entity {
	return &Entity{
		Position: &Point{x, y},
		Scale:    &Point{1, 1},
		Pivot:    &Point{0.5, 0.5},
		Rotation: 0,
		color:    whiteFloatBits,
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
	batch.Draw(s.region, s.Position.X, s.Position.Y, s.Pivot.X, s.Pivot.Y, s.Scale.X, s.Scale.Y, s.Rotation, s.color)
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
				batch.Draw(region, x, y, 0, 0, sx, sy, 0, t.color)
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

func (s *Stage) SetBg(color *Color) {
	SetBgColor(color)
}

func (s *Stage) Add(object Displayer) {
	s.objects = append(s.objects, object)
}

func (s *Stage) Sprite(region *Region, x, y float32) *Sprite {
	sprite := NewSprite(region, x, y)
	s.Add(sprite)
	return sprite
}

func (s *Stage) Text(font *Font, x, y float32, content string) *Text {
	text := NewText(font, x, y, content)
	s.Add(text)
	return text
}

func (s *Stage) Load() {}
func (s *Stage) init() {
	s.batch = NewBatch()
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
func (s *Stage) Resize(width, height int)          {}
func (s *Stage) Update(dt float32)                 {}
func (s *Stage) Mouse(x, y float32, action Action) {}
func (s *Stage) Scroll(amount float32)             {}
func (s *Stage) Type(char rune)                    {}
func (s *Stage) Key(key Key, mod Modifier, act Action) {
	if key == Escape {
		Exit()
	}
}
