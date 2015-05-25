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

func (p *Point) PointDistance(p2 Point) float32 {
	return float32(math.Sqrt(float64(p.PointDistanceSquared(p2))))
}

func (p *Point) PointDistanceSquared(p2 Point) float32 {
	return (p.X-p2.X)*(p.X-p2.X) + (p.Y-p2.Y)*(p.Y-p2.Y)
}

// Returns the vector produced by projecting a on to b
func (a *Point) ProjectOnto(b Point) Point {
	dot := a.X*b.X + a.Y*b.Y
	proj := Point{
		dot / (b.X*b.X + b.Y*b.Y) * b.X,
		dot / (b.X*b.X + b.Y*b.Y) * b.Y,
	}
	return proj
}

// Returns the unit vector from a, and it's magnitude
func (a *Point) Normalize() (Point, float32) {
	mag := float32(math.Sqrt(float64(a.X*a.X + a.Y*a.Y)))
	unit := Point{a.X / mag, a.Y / mag}

	return unit, mag
}

type Line struct {
	P1 Point
	P2 Point
}

// Returns which side of the line the point is on
// This is useful if you have a point of reference
func (l *Line) PointSide(point Point) bool {
	one := (point.X - l.P1.X) * (l.P2.Y - l.P1.Y)
	two := (point.Y - l.P1.Y) * (l.P2.X - l.P1.X)

	return math.Signbit(float64(one - two))
}

// Returns the line's angle relative to Y = 0
func (l *Line) Angle() float32 {
	return float32(math.Atan2(float64(l.P1.X-l.P2.X), float64(l.P1.Y-l.P2.Y)))
}

// Returns the squared euclidean distance from a point to a line *segment*
func (l *Line) PointDistance(point Point) float32 {
	return float32(math.Sqrt(float64(l.PointDistanceSquared(point))))
}

// Returns the squared euclidean distance from a point to a line *segment*
func (l *Line) PointDistanceSquared(point Point) float32 {
	p1 := l.P1
	p2 := l.P2

	x0 := point.X
	y0 := point.Y
	x1 := p1.X
	y1 := p1.Y
	x2 := p2.X
	y2 := p2.Y

	l2 := (y2-y1)*(y2-y1) + (x2-x1)*(x2-x1)
	if l2 == 0 {
		return (y0-y1)*(y0-y1) + (x0-x1)*(x0-x1)
	}

	t := ((x0-x1)*(x2-x1) + (y0-y1)*(y2-y1)) / l2

	if t < 0 {
		return (y0-y1)*(y0-y1) + (x0-x1)*(x0-x1)
	} else if t > 1 {
		return (y0-y2)*(y0-y2) + (x0-x2)*(x0-x2)
	}

	return (x0-(x1+t*(x2-x1)))*(x0-(x1+t*(x2-x1))) +
		(y0-(y1+t*(y2-y1)))*(y0-(y1+t*(y2-y1)))
}

// Returns the point where the two lines intersect
func (l *Line) LineIntersection(l2 Line) Point {
	x1 := l.P1.X
	x2 := l.P2.X
	x3 := l2.P1.X
	x4 := l2.P2.X

	y1 := l.P1.Y
	y2 := l.P2.Y
	y3 := l2.P1.Y
	y4 := l2.P2.Y

	denom := ((x1-x2)*(y3-y4) - (y1-y2)*(x3-x4))
	if denom == 0 {
		return Point{-1, -1}
	}

	px := ((x1*y2-y1*x2)*(x3-x4) - (x1-x2)*(x3*y4-y3*x4)) / denom
	py := ((x1*y2-y1*x2)*(y3-y4) - (y1-y2)*(x3*y4-y3*x4)) / denom

	return Point{px, py}
}
