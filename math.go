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

type Matrix struct {
	A, B, C, D, TX, TY float32
}

func NewMatrix() *Matrix {
	return &Matrix{1, 0, 0, 1, 0, 0}
}
