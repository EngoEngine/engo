package engo

import (
	"engo.io/engo/math"
)

const (
	// Epsilon is some tiny value that determines how precisely equal we want our
	// floats to be.
	Epsilon float32 = 1e-3
	// MinNormal is the smallest normal value possible.
	MinNormal = float32(1.1754943508222875e-38) // 1 / 2**(127 - 1)
)

// AABB describes two points of a rectangle: the upper-left corner and the lower-right corner. It should always hold that
// `Min.X <= Max.X` and `Min.Y <= Max.Y`.
type AABB struct {
	Min, Max Point
}

// A Container is a 2D closed shape which contains a set of points.
type Container interface {
	// Contains reports whether the container contains the given point.
	Contains(p Point) bool
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
	Line
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
	return FloatEqual(p.X, p2.X) && FloatEqual(p.Y, p2.Y)
}

// PointDistance returns the euclidean distance between p and p2
func (p *Point) PointDistance(p2 Point) float32 {
	return math.Sqrt(p.PointDistanceSquared(p2))
}

// PointDistanceSquared returns the squared euclidean distance between p and p2
func (p *Point) PointDistanceSquared(p2 Point) float32 {
	return (p.X-p2.X)*(p.X-p2.X) + (p.Y-p2.Y)*(p.Y-p2.Y)
}

// ProjectOnto returns the vector produced by projecting p on to p2
// returns an empty Point if they can't project onto one another
func (p *Point) ProjectOnto(p2 Point) Point {
	dot := p.X*p2.X + p.Y*p2.Y
	denom := p2.X*p2.X + p2.Y*p2.Y
	if FloatEqual(denom, 0) {
		return Point{}
	}
	return Point{
		dot / denom * p2.X,
		dot / denom * p2.Y,
	}
}

// Normalize returns the unit vector from p, and its magnitude.
// if you try to normalize the null vector, the return value will be null values
func (p *Point) Normalize() (Point, float32) {
	if p.X == 0 && p.Y == 0 {
		return *p, 0
	}

	mag := math.Sqrt(p.X*p.X + p.Y*p.Y)
	unit := Point{p.X / mag, p.Y / mag}

	return unit, mag
}

// Within reports whether the point is contained within the given container.
func (p Point) Within(c Container) bool {
	return c.Contains(p)
}

// PointSide returns which side of the line l the point p sits on
// true means the point is below/left of the line
// false means the point is above/right of the line or touching the line
func (l *Line) PointSide(point Point) bool {
	one := (point.X - l.P1.X) * (l.P2.Y - l.P1.Y)
	two := (point.Y - l.P1.Y) * (l.P2.X - l.P1.X)

	return math.Signbit(one - two)
}

// Angle returns the euclidean angle of l relative to X = 0
// The return angle is in radians and goes counter-clockwise and returns [-pi, pi]
func (l *Line) Angle() float32 {
	return math.Atan2(l.P1.X-l.P2.X, l.P1.Y-l.P2.Y)
}

// PointDistance Returns the euclidean distance from the point p to the
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
// intersect and true if there is intersection, nil and false when line
// segments one and two do not intersect
func LineIntersection(one, two Line) (Point, bool) {
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

	t := qmpcs / rcs
	u := qmpcr / rcs
	// if rcs == 0 then it's either collinear or parallel. It'll be +/- inf, so it'll skip this statement and return at the end
	if t >= 0 && t <= 1 && u >= 0 && u <= 1 {
		// the two line segments meet at the point p + t r = q + u s.
		return Point{p.X + t*r.X, p.Y + t*r.Y}, true
	}

	return Point{}, false
}

// LineTraceFraction returns the trace fraction of tracer through boundary
// 1 means no intersection
// 0 means tracer's origin lies on the boundary line
func LineTraceFraction(tracer, boundary Line) float32 {
	pt, intersect := LineIntersection(tracer, boundary)
	if !intersect {
		return 1
	}

	traceMag := tracer.P1.PointDistance(pt)
	lineMag := tracer.P1.PointDistance(tracer.P2)

	return traceMag / lineMag
}

// LineTrace runs a series of line traces from tracer to each boundary line
// and returns the nearest trace values
func LineTrace(tracer Line, boundaries []Line) Trace {
	var t Trace
	t.Fraction = math.Inf(1)

	for _, cl := range boundaries {
		fraction := LineTraceFraction(tracer, cl)

		if fraction < t.Fraction {
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

// FloatEqual is a safe utility function to compare floats.
// It's Taken from http://floating-point-gui.de/errors/comparison/
//
// It is slightly altered to not call Abs when not needed.
func FloatEqual(a, b float32) bool {
	return FloatEqualThreshold(a, b, Epsilon)
}

// FloatEqualThreshold is a utility function to compare floats.
// It's Taken from http://floating-point-gui.de/errors/comparison/
//
// It is slightly altered to not call Abs when not needed.
//
// This differs from FloatEqual in that it lets you pass in your comparison threshold, so that you can adjust the comparison value to your specific needs
func FloatEqualThreshold(a, b, epsilon float32) bool {
	if a == b { // Handles the case of inf or shortcuts the loop when no significant error has accumulated
		return true
	}

	if math.IsNaN(a) || math.IsNaN(b) {
		return false // Can't be equal if NaN
	}

	diff := math.Abs(a - b)
	if a*b == 0 || diff < MinNormal { // If a or b are 0 or both are extremely close to it
		return diff < epsilon*epsilon
	}

	// Else compare difference
	return diff/(math.Abs(a)+math.Abs(b)) < epsilon
}
