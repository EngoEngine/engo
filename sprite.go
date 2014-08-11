package engi

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

type Sprite struct {
	Position *Point
	Scale    *Point
	Anchor   *Point
	Rotation float32
	Color    uint32
	Alpha    float32
	Region   *Region
}

func NewSprite(region *Region, x, y float32) *Sprite {
	return &Sprite{
		Position: &Point{x, y},
		Scale:    &Point{1, 1},
		Anchor:   &Point{0, 0},
		Rotation: 0,
		Color:    0xffffff,
		Alpha:    1,
		Region:   region,
	}
}

func (s *Sprite) Render(batch *Batch) {
	batch.Draw(s.Region, s.Position.X, s.Position.Y, s.Anchor.X, s.Anchor.Y, s.Scale.X, s.Scale.Y, s.Rotation, s.Color, s.Alpha)
}
