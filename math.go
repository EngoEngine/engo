package engi

import (
	"math"
)

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

func (p *Point) AddScalar(s float32) {
	p.X += s
	p.Y += s
}

func (p *Point) SubtractScalar(s float32) {
	p.X -= s
	p.Y -= s
}

func (p *Point) MultiplyScalar(s float32) {
	p.X *= s
	p.Y *= s
}

func (p *Point) Add(p2 Point) {
	p.X += p2.X
	p.Y += p2.Y
}

func (p *Point) Subtract(p2 Point) {
	p.X -= p2.X
	p.Y -= p2.Y
}

func (p *Point) Multiply(p2 Point) {
	p.X *= p2.X
	p.Y *= p2.Y
}

func (p *Point) DistanceFrom(p2 Point) float32 {
	return float32(math.Sqrt(float64(p.DistanceFromSquared(p2))))
}

func (p *Point) DistanceFromSquared(p2 Point) float32 {
	return float32(math.Pow(float64(p.X-p2.X), 2) + math.Pow(float64(p.Y-p2.Y), 2))
}
