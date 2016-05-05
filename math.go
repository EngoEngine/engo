package engo

import (
	"github.com/luxengine/math"
)

type Point struct {
	X, Y float32
}

type Line struct {
	P1 Point
	P2 Point
}

type Trace struct {
	Fraction float32
	EndPosition Point
	*Line
}

func (p *Point) Set(x, y float32) {
	p.X = x
	p.Y = y

}

func (p *Point) Dot(other Point) float32 {
	return p.X*other.X + p.Y*other.Y
}

// 2D cross product is magnitude of 3D cross product
func (p *Point) Cross(other Point) float32 {
	return p.X*other.Y - p.Y*other.X
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

func (p *Point) Divide(p2 Point) {
	p.X /= p2.X
	p.Y /= p2.Y
}

func (p *Point) PointDistance(p2 Point) float32 {
	return math.Sqrt(p.PointDistanceSquared(p2))
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
	mag := math.Sqrt(a.X*a.X + a.Y*a.Y)
	unit := Point{a.X / mag, a.Y / mag}

	return unit, mag
}

// Returns which side of the line the point is on
// This is useful if you have a point of reference
func (l *Line) PointSide(point Point) bool {
	one := (point.X - l.P1.X) * (l.P2.Y - l.P1.Y)
	two := (point.Y - l.P1.Y) * (l.P2.X - l.P1.X)

	return math.Signbit(one - two)
}

// Returns the line's angle relative to Y = 0
func (l *Line) Angle() float32 {
	return math.Atan2(l.P1.X-l.P2.X, l.P1.Y-l.P2.Y)
}

// Returns the squared euclidean distance from a point to a line *segment*
func (l *Line) PointDistance(point Point) float32 {
	return math.Sqrt(l.PointDistanceSquared(point))
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

// Left Hand Normal
func (l *Line) Normal() Point {
	dx := l.P2.X - l.P1.X
	dy := l.P2.Y - l.P1.Y
	inverse := Point{dy, -dx}
	unit, _ := inverse.Normalize()

	return unit
}


// Returns the point where the two line *segments* intersect
func LineIntersection(one, two *Line) Point {
	p := one.P1
	q := two.P1

	r := one.P2
	r.Subtract(p)
	s := two.P2
	s.Subtract(q)

	// t = (q − p) × s / (r × s)
	// u = (q − p) × r / (r × s)
	// So then we define
	// qmp = (q - p)
	// rcs = (r × s)
	// and we get simply:
	// t = qmp × s / rcs
	// u = qmp × r / rcs
	qmp := q
	qmp.Subtract(p)
	qmpcs := qmp.Cross(s)
	qmpcr := qmp.Cross(r)
	rcs := r.Cross(s)

	// Collinear
	if rcs == 0 && qmpcr == 0 {
		return Point{-1, -1}
	}

	// Parallel
	if rcs == 0 && qmpcr != 0 {
		return Point{-1, -1}
	}

	t := qmpcs / rcs
	u := qmpcr / rcs
	// rcs != 0 at this point
	if t >= 0 && t <= 1 && u >= 0 && u <= 1 {
		// the two line segments meet at the point p + t r = q + u s.
		return Point{p.X + t*r.X, p.Y + t*r.Y}
	}

	return Point{-1, -1}
}

// Returns the trace fraction of tracer through boundary
// 1 means no intersection
// 0 means tracer's origin lies on the boundary line
func LineTraceFraction(tracer, boundary *Line) float32 {

	pt := LineIntersection(tracer, boundary)
	if pt.X == -1 && pt.Y == -1 {
		return 1
	}

	traceMag := tracer.P1.PointDistance(pt)
	lineMag := tracer.P1.PointDistance(tracer.P2)

	if traceMag > lineMag {
		return 1
	}

	if lineMag == 0 {
		return 0
	}

	return traceMag / lineMag
}

// Runs a series of line tracers from tracer to each boundary line
// and returns the nearest trace values
func LineTrace(tracer *Line, boundaries []*Line) Trace {
	var t Trace

	for _, cl := range boundaries {
		//TODO why are some lines nil here?
		//fmt.Println("Line:", cl)
		if cl == nil {
			continue
		}

		fraction := LineTraceFraction(tracer, cl)

		if t.Line == nil || fraction < t.Fraction {
			t.Fraction = fraction
			t.Line = cl

			moveVector := tracer.P2
			moveVector.Subtract(tracer.P1)
			moveVector.MultiplyScalar(t.Fraction)
			t.EndPosition = tracer.P1
			t.EndPosition.Add(moveVector)
		}
	}

	return t
}
