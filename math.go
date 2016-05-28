package engo

import (
	"github.com/engoengine/glm"
	"github.com/luxengine/math"
)

func init() {
	// This precision / error margin is required to work with float32 within `engo.Point` when checking for equality.
	glm.Epsilon = 1e-3
}

// AABB describes two points of a rectangle: the upper-left corner and the lower-right corner. It should always hold that
// `Min.X <= Max.X` and `Min.Y <= Max.Y`.
type AABB struct {
	Min, Max Point
}

// Point describes a coordinate in a 2 dimensional euclidean space
// it can also be thought of as a 2 dimensional vector from the origin
type Point struct {
	X, Y float32
}

// Line describes a line segment on a 2 dimensional euclidean space
// it can also be thought of as a 2 dimensional vector with an offset
type Line struct {
	P1 Point
	P2 Point
}

// Trace describes all the values computed from a line trace
type Trace struct {
	Fraction    float32
	EndPosition Point
	*Line
}

// Set sets the coordinates of p to x and y
func (p *Point) Set(x, y float32) {
	p.X = x
	p.Y = y
}

// AddScalar adds s to each component of p
func (p *Point) AddScalar(s float32) {
	p.X += s
	p.Y += s
}

// SubtractScalar subtracts s from each component of p
func (p *Point) SubtractScalar(s float32) {
	p.AddScalar(-s)
}

// MultiplyScalar multiplies each component of p by s
func (p *Point) MultiplyScalar(s float32) {
	p.X *= s
	p.Y *= s
}

// Add sets the components of p to the pointwise summation of p + p2
func (p *Point) Add(p2 Point) {
	p.X += p2.X
	p.Y += p2.Y
}

// Subtract sets the components of p to the pointwise difference of p - p2
func (p *Point) Subtract(p2 Point) {
	p.X -= p2.X
	p.Y -= p2.Y
}

// Multiply sets the components of p to the pointwise product of p * p2
func (p *Point) Multiply(p2 Point) {
	p.X *= p2.X
	p.Y *= p2.Y
}

// Equal indicates whether two points have the same value, avoiding issues with float precision
func (p *Point) Equal(p2 Point) bool {
	return glm.FloatEqual(p.X, p2.X) && glm.FloatEqual(p.Y, p2.Y)
}

// PointDistance returns the euclidean distance between p and p2
func (p *Point) PointDistance(p2 Point) float32 {
	return math.Sqrt(p.PointDistanceSquared(p2))
}

// PointDistanceSquared returns the squared euclidean distance between p and p2
func (p *Point) PointDistanceSquared(p2 Point) float32 {
	return (p.X-p2.X)*(p.X-p2.X) + (p.Y-p2.Y)*(p.Y-p2.Y)
}

// ProjectOnto returns the vector produced by projecting a on to b
func (a *Point) ProjectOnto(b Point) Point {
	dot := a.X*b.X + a.Y*b.Y
	proj := Point{
		dot / (b.X*b.X + b.Y*b.Y) * b.X,
		dot / (b.X*b.X + b.Y*b.Y) * b.Y,
	}
	return proj
}

// Normalize returns the unit vector from a, and its magnitude.
// if you try to normalize the null vector, the return value will be null values
func (a *Point) Normalize() (Point, float32) {
	if a.X == 0 && a.Y == 0 {
		return *a, 0
	}

	mag := math.Sqrt(a.X*a.X + a.Y*a.Y)
	unit := Point{a.X / mag, a.Y / mag}

	return unit, mag
}

// PointSide returns which side of the line l the point p sits on
func (l *Line) PointSide(point Point) bool {
	one := (point.X - l.P1.X) * (l.P2.Y - l.P1.Y)
	two := (point.Y - l.P1.Y) * (l.P2.X - l.P1.X)

	return math.Signbit(one - two)
}

// Angle returns the euclidean angle of l relative to Y = 0
func (l *Line) Angle() float32 {
	return math.Atan2(l.P1.X-l.P2.X, l.P1.Y-l.P2.Y)
}

// PointDistance Returns the squared euclidean distance from the point p to the
// line segment l
func (l *Line) PointDistance(point Point) float32 {
	return math.Sqrt(l.PointDistanceSquared(point))
}

// PointDistanceSquared returns the squared euclidean distance from the point p
// to the line segment l
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

// Normal returns the left hand normal of the line segment l
func (l *Line) Normal() Point {
	dx := l.P2.X - l.P1.X
	dy := l.P2.Y - l.P1.Y
	inverse := Point{dy, -dx}
	unit, _ := inverse.Normalize()

	return unit
}

// DotProduct returns the dot product between this and that
func DotProduct(this, that Point) float32 {
	return this.X*that.X + this.Y*that.Y
}

// CrossProduct returns the 2 dimensional cross product of this and that,
// which represents the magnitude of the three dimensional cross product
func CrossProduct(this, that Point) float32 {
	return this.X*that.Y - this.Y*that.X
}

// LineIntersection returns the point where the line segments one and two
// intersect
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
	qmpcs := CrossProduct(qmp, s)
	qmpcr := CrossProduct(qmp, r)
	rcs := CrossProduct(r, s)

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

// LineTraceFraction returns the trace fraction of tracer through boundary
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

// LineTrace runs a series of line traces from tracer to each boundary line
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
