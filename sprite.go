package engi

import "math"

type Sprite struct {
	Position  *Point
	Scale     *Point
	Anchor    *Point
	Rotation  float32
	color     float32
	Alpha     float32
	region    *Region
	transform *Matrix
	parent    *Sprite
	children  []*Sprite
	alphaMult float32
}

func NewSprite(region *Region, x, y float32) *Sprite {
	return &Sprite{
		Position:  &Point{x, y},
		Scale:     &Point{1, 1},
		Anchor:    &Point{0, 0},
		Rotation:  0,
		color:     whiteBits,
		Alpha:     1,
		region:    region,
		transform: NewMatrix(),
		parent:    nil,
		children:  make([]*Sprite, 0),
		alphaMult: 1,
	}
}

func NewGroup() *Sprite {
	return NewSprite(nil, 0, 0)
}

func (s *Sprite) updateTransform() {
	rot := float64(s.Rotation * degToRad)

	sr := float32(math.Sin(float64(rot)))
	cr := float32(math.Cos(float64(rot)))

	parentTransform := s.parent.transform
	worldTransform := s.transform

	px := s.Anchor.X
	py := s.Anchor.Y

	a00 := cr * s.Scale.X
	a01 := -sr * s.Scale.Y
	a10 := sr * s.Scale.X
	a11 := cr * s.Scale.Y
	a02 := s.Position.X - a00*px - py*a01
	a12 := s.Position.Y - a11*py - px*a10
	b00 := parentTransform.A
	b01 := parentTransform.B
	b10 := parentTransform.C
	b11 := parentTransform.D

	worldTransform.A = b00*a00 + b01*a10
	worldTransform.B = b00*a01 + b01*a11
	worldTransform.TX = b00*a02 + b01*a12 + parentTransform.TX

	worldTransform.C = b10*a00 + b11*a10
	worldTransform.D = b10*a01 + b11*a11
	worldTransform.TY = b10*a02 + b11*a12 + parentTransform.TY
}

func (s *Sprite) AddChild(child *Sprite) *Sprite {
	return s.AddChildAt(child, len(s.children))
}

func (s *Sprite) AddChildAt(child *Sprite, index int) *Sprite {
	if index < 0 || index > len(s.children) || child == nil {
		return nil
	}

	if child.parent != nil {
		child.parent.RemoveChild(child)
	}

	child.parent = s

	s.children = append(s.children, nil)
	copy(s.children[index+1:], s.children[index:])
	s.children[index] = child

	return child
}

func (s *Sprite) SetTint(tint uint32) {
	r := uint32((tint >> 16) & 0xFF)
	g := uint32((tint >> 8) & 0xFF)
	b := uint32(tint & 0xFF)
	a := uint32(s.Alpha * 255.0)
	s.color = math.Float32frombits((a<<24 | b<<16 | g<<8 | r) & 0xfeffffff)
}

func (s *Sprite) indexOf(child *Sprite) int {
	for i, val := range s.children {
		if val == child {
			return i
		}
	}
	return -1
}

func (s *Sprite) RemoveChild(child *Sprite) *Sprite {
	return s.RemoveChildAt(s.indexOf(child))
}

func (s *Sprite) RemoveChildAt(index int) *Sprite {
	if index < 0 || index >= len(s.children) {
		return nil
	}

	child := s.children[index]
	child.parent = nil

	copy(s.children[index:], s.children[index+1:])
	s.children[len(s.children)-1] = nil
	s.children = s.children[:len(s.children)-1]

	return child
}

func (s *Sprite) render(batch *Batch) {
	s.updateTransform()
	batch.Render(s)
	for _, child := range s.children {
		child.render(batch)
	}
}
