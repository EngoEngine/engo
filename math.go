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

	// RadToDeg is multiplied with a radian value to get the equivalant value in degrees.
	RadToDeg = 180 / math.Pi
	// DegToRad is multiplied with a degree value to get the equivalent value in radians.
	DegToRad = math.Pi / 180
	// Matrix row/column indexes
	m00 = 0
	m01 = 3
	m02 = 6
	m10 = 1
	m11 = 4
	m12 = 7
	m20 = 2
	m21 = 5
	m22 = 8
)

// AABB describes two points of a rectangle: the upper-left corner and the lower-right corner. It should always hold that
// `Min.X <= Max.X` and `Min.Y <= Max.Y`.
type AABB struct {
	Min, Max Point
}

// AABBer is an interface for everything that provides information about its axis aligned bounding box.
type AABBer interface {
	// AABB returns the axis aligned bounding box.
	AABB() AABB
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

// Matrix describes a 3x3 column-major matrix useful for 2D transformations.
type Matrix struct {
	Val [9]float32
	tmp [9]float32
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
func (p *Point) Set(x, y float32) *Point {
	p.X = x
	p.Y = y
	return p
}

// AddScalar adds s to each component of p
func (p *Point) AddScalar(s float32) *Point {
	p.X += s
	p.Y += s
	return p
}

// SubtractScalar subtracts s from each component of p
func (p *Point) SubtractScalar(s float32) *Point {
	p.AddScalar(-s)
	return p
}

// MultiplyScalar multiplies each component of p by s
func (p *Point) MultiplyScalar(s float32) *Point {
	p.X *= s
	p.Y *= s
	return p
}

// Add sets the components of p to the pointwise summation of p + p2
func (p *Point) Add(p2 Point) *Point {
	p.X += p2.X
	p.Y += p2.Y
	return p
}

// Subtract sets the components of p to the pointwise difference of p - p2
func (p *Point) Subtract(p2 Point) *Point {
	p.X -= p2.X
	p.Y -= p2.Y
	return p
}

// Multiply sets the components of p to the pointwise product of p * p2
func (p *Point) Multiply(p2 Point) *Point {
	p.X *= p2.X
	p.Y *= p2.Y
	return p
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

// Set sets the matrix to the given float slice and returns the matrix. The float
// slice must have at least 9 elements. If the float slie contains more than 9 elements,
// only the first 9 will be copied.
func (m *Matrix) Set(val []float32) *Matrix {
	copy(m.Val[:], val)
	return m
}

// Identity sets the matrix to the Identity matrix and returns the matrix.
func (m *Matrix) Identity() *Matrix {
	m.Val[m00] = 1
	m.Val[m10] = 0
	m.Val[m20] = 0
	m.Val[m01] = 0
	m.Val[m11] = 1
	m.Val[m21] = 0
	m.Val[m02] = 0
	m.Val[m12] = 0
	m.Val[m22] = 1
	return m
}

// Multiply postmultiplies m matrix with m2 and stores the result in m, returning m.
// Multiplaction is the result of m2 times m.
func (m *Matrix) Multiply(m2 *Matrix) *Matrix {
	multiplyMatricies(m.Val[:], m2.Val[:])
	return m
}

// Translate translates m by the point (x, y).
func (m *Matrix) Translate(x, y float32) *Matrix {
	m.tmp[m00] = 1
	m.tmp[m10] = 0
	m.tmp[m20] = 0

	m.tmp[m01] = 0
	m.tmp[m11] = 1
	m.tmp[m21] = 0

	m.tmp[m02] = x
	m.tmp[m12] = y
	m.tmp[m22] = 1

	multiplyMatricies(m.Val[:], m.tmp[:])
	return m
}

// TranslatePoint translates m by the point p.
func (m *Matrix) TranslatePoint(p Point) *Matrix {
	return m.Translate(p.X, p.Y)
}

// Scale scales m by x and y.
func (m *Matrix) Scale(x, y float32) *Matrix {
	m.tmp[m00] = x
	m.tmp[m10] = 0
	m.tmp[m20] = 0

	m.tmp[m01] = 0
	m.tmp[m11] = y
	m.tmp[m21] = 0

	m.tmp[m02] = 0
	m.tmp[m12] = 0
	m.tmp[m22] = 1
	multiplyMatricies(m.Val[:], m.tmp[:])
	return m
}

// ScaleComponent returns the current scale component of m.
// This assumes uniform scaling.
func (m *Matrix) ScaleComponent() (x, y float32) {
	return m.Val[m00], m.Val[m11]
}

// TranslationComponent returns the current translation component of m.
// This assumes uniform scaling.
func (m *Matrix) TranslationComponent() (x, y float32) {
	return m.Val[m02], m.Val[m12]
}

// RotationComponent returns the current rotation component of m in degrees.
// This assumes uniform scaling.
func (m *Matrix) RotationComponent() float32 {
	return m.RotationComponentRad() * RadToDeg
}

// RotationComponentRad returns the current rotation component of m in radians.
// This assumes uniform scaling.
func (m *Matrix) RotationComponentRad() float32 {
	return math.Atan2(m.Val[m10], m.Val[m00])
}

// Rotate rorates m counter-clockwise by deg degrees.
func (m *Matrix) Rotate(deg float32) *Matrix {
	return m.RotateRad(deg * DegToRad)
}

// RotateRad rotates m counter-clockwise by rad radians.
func (m *Matrix) RotateRad(rad float32) *Matrix {
	if rad == 0 {
		return m
	}
	sin, cos := math.Sincos(rad)
	m.tmp[m00] = cos
	m.tmp[m10] = sin
	m.tmp[m20] = 0

	m.tmp[m01] = -sin
	m.tmp[m11] = cos
	m.tmp[m21] = 0

	m.tmp[m02] = 0
	m.tmp[m12] = 0
	m.tmp[m22] = 1
	multiplyMatricies(m.Val[:], m.tmp[:])
	return m
}

// PointSide returns which side of the line l the point p sits on
// true means the point is below/left of the line
// false means the point is above/right of the line or touching the line
func (l *Line) PointSide(point Point) bool {
	one := (point.X - l.P1.X) * (l.P2.Y - l.P1.Y)
	two := (point.Y - l.P1.Y) * (l.P2.X - l.P1.X)

	return math.Signbit(one - two)
}

// Angle returns the euclidean angle of l in radians relative to a vertical line, going
// positive as you head towards the positive x-axis (clockwise) and negative
// as you head towards the negative x-axis. Values returned are [-pi, pi].
func (l *Line) Angle() float32 {
	return math.Atan2(l.P1.X-l.P2.X, l.P1.Y-l.P2.Y)
}

// AngleDeg returns the euclidean angle of l in degrees relative to a vertical line, going
// positive as you head towards the positive x-axis (clockwise) and negative
// as you head towards the negative x-axis. Values returned are [-180, 180].
func (l *Line) AngleDeg() float32 {
	x := l.P2.X - l.P1.X
	y := l.P2.Y - l.P1.Y
	if x == 0 {
		if y > 0 {
			return 180
		}
		return 0
	}

	deg := math.Atan(x/y) * 180 / math.Pi
	if x > 0 && y < 0 {
		deg = -deg
	} else if x < 0 && y < 0 {
		deg = 360 - deg
	} else {
		deg = 180 - deg
	}
	return deg
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

// IdentityMatrix returns a new identity matrix.
func IdentityMatrix() *Matrix {
	return new(Matrix).Identity()
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

func multiplyMatricies(m1, m2 []float32) {
	v00 := m1[m00]*m2[m00] + m1[m01]*m2[m10] + m1[m02]*m2[m20]
	v01 := m1[m00]*m2[m01] + m1[m01]*m2[m11] + m1[m02]*m2[m21]
	v02 := m1[m00]*m2[m02] + m1[m01]*m2[m12] + m1[m02]*m2[m22]

	v10 := m1[m10]*m2[m00] + m1[m11]*m2[m10] + m1[m12]*m2[m20]
	v11 := m1[m10]*m2[m01] + m1[m11]*m2[m11] + m1[m12]*m2[m21]
	v12 := m1[m10]*m2[m02] + m1[m11]*m2[m12] + m1[m12]*m2[m22]

	v20 := m1[m20]*m2[m00] + m1[m21]*m2[m10] + m1[m22]*m2[m20]
	v21 := m1[m20]*m2[m01] + m1[m21]*m2[m11] + m1[m22]*m2[m21]
	v22 := m1[m20]*m2[m02] + m1[m21]*m2[m12] + m1[m22]*m2[m22]
	m1[m00] = v00
	m1[m10] = v10
	m1[m20] = v20
	m1[m01] = v01
	m1[m11] = v11
	m1[m21] = v21
	m1[m02] = v02
	m1[m12] = v12
	m1[m22] = v22
}

// MultiplyMatrixVector multiplies the matrix m with the float32 vector v and returns the result.
// The size of vector v MUST be 2 or 3. If v is size 2, a 3rd component is automatically added with
// value of 1.0.
func MultiplyMatrixVector(m *Matrix, v []float32) []float32 {
	if len(v) == 2 {
		v = []float32{v[0], v[1], 1}
	}
	v00 := m.Val[m00]*v[m00] + m.Val[m01]*v[m10] + m.Val[m02]*v[m20]
	v10 := m.Val[m10]*v[m00] + m.Val[m11]*v[m10] + m.Val[m12]*v[m20]
	v20 := m.Val[m20]*v[m00] + m.Val[m21]*v[m10] + m.Val[m22]*v[m20]
	return []float32{v00, v10, v20}
}

// MultiplyMatrixVector multiplies the matrix m with the point and returns the result.
func (p *Point) MultiplyMatrixVector(m *Matrix) *Point {
	x := m.Val[m00]*p.X + m.Val[m01]*p.Y + m.Val[m02]
	y := m.Val[m10]*p.X + m.Val[m11]*p.Y + m.Val[m12]
	p.X, p.Y = x, y
	return p
}
