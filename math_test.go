package engo

import (
	"testing"

	"engo.io/engo/math"
)

func TestFloatEqual(t *testing.T) {
	data := []struct {
		a, b float32
		exp  bool
	}{
		{
			a:   0,
			b:   0,
			exp: true,
		},
		{
			a:   5,
			b:   10,
			exp: false,
		},
		{
			a:   math.NaN(),
			b:   0,
			exp: false,
		},
		{
			a:   math.NaN(),
			b:   math.NaN(),
			exp: false,
		},
	}
	for _, d := range data {
		if actual := FloatEqual(d.a, d.b); actual != d.exp {
			t.Errorf("Test FloatEqual failed. a: %v, b: %v, wanted: %v, got: %v", d.a, d.b, d.exp, actual)
		}
	}
}

func TestPointPointAddition(t *testing.T) {
	a := Point{2, 5}
	b := Point{3, 1}
	a.Add(b)
	if a.X != 5 {
		t.Errorf("a.X should equal 5 not %v", a.X)
	}

	if a.Y != 6 {
		t.Errorf("a.Y should equal 6 not %v", a.Y)
	}
}

func TestPointPointSubtraction(t *testing.T) {
	a := Point{10, 15}
	b := Point{5, 2}
	a.Subtract(b)
	if a.X != 5 {
		t.Errorf("a.X should equal 5 not %v", a.X)
	}

	if a.Y != 13 {
		t.Errorf("a.Y should equal 13 not %v", a.Y)
	}
}

func TestPointPointMultiplication(t *testing.T) {
	a := Point{10, 2}
	b := Point{5, 6}
	a.Multiply(b)
	if a.X != 50 {
		t.Errorf("a.X should equal 50 not %v", a.X)
	}

	if a.Y != 12 {
		t.Errorf("a.Y should equal 12 not %v", a.Y)
	}
}

func TestPointScalarAddition(t *testing.T) {
	a := Point{2, 4}
	s := float32(1)
	a.AddScalar(s)

	if a.X != 3 {
		t.Errorf("a.X should equal 3 not %v", a.X)
	}

	if a.Y != 5 {
		t.Errorf("a.Y should equal 5 not %v", a.Y)
	}
}

func TestPointScalarSubtraction(t *testing.T) {
	a := Point{10, 20}
	s := float32(2)
	a.SubtractScalar(s)

	if a.X != 8 {
		t.Errorf("a.X should equal 8 not %v", a.X)
	}

	if a.Y != 18 {
		t.Errorf("a.Y should equal 18 not %v", a.Y)
	}
}

func TestPointScalarMultiplication(t *testing.T) {
	a := Point{5, 6}
	s := float32(3)

	a.MultiplyScalar(s)
	if a.X != 15 {
		t.Errorf("a.X should equal 15 not %v", a.X)
	}

	if a.Y != 18 {
		t.Errorf("a.X should equal 18 not %v", a.Y)
	}
}

func TestLineIntersection(t *testing.T) {
	//Parallel lines
	one := Line{
		Point{0, 0},
		Point{1, 1}}
	two := Line{
		Point{0, 1},
		Point{1, 2}}

	point, intersect := LineIntersection(one, two)

	if intersect {
		t.Errorf("Lines %v and %v should not intersect, they are parallel.  Intersection at: %v.", one, two, point)
	}

	//Collinear lines
	one = Line{
		Point{0, 0},
		Point{1, 1}}
	two = Line{
		Point{2, 2},
		Point{3, 3}}

	point, intersect = LineIntersection(one, two)

	if intersect {
		t.Errorf("Lines %v and %v should not intersect, they are collinear. Intersection at: %v", one, two, point)
	}

	//intersecting lines
	one = Line{
		Point{0, 0},
		Point{1, 1}}
	two = Line{
		Point{0, 1},
		Point{1, 1}}

	point, intersect = LineIntersection(one, two)

	if !intersect {
		t.Errorf("Lines %v and %v should intersect.", one, two)
	}

	if intersect && (point != Point{1, 1}) {
		t.Errorf("Lines %v and %v should intersect at point {1, 1}, but they are intersecting at %v.", one, two, point)
	}

}

func TestPointEqual(t *testing.T) {
	data := []struct {
		p1, p2 Point
		res    bool
	}{
		{
			p1:  Point{X: 0, Y: 0},
			p2:  Point{X: 0, Y: 0},
			res: true,
		},
		{
			p1:  Point{X: 3, Y: -8},
			p2:  Point{X: 3, Y: -8},
			res: true,
		},
		{
			p1:  Point{X: -3, Y: 3},
			p2:  Point{X: 3, Y: -3},
			res: false,
		},
		{
			p1:  Point{X: 3, Y: 6},
			p2:  Point{X: 2, Y: 5},
			res: false,
		},
	}

	for _, d := range data {
		if b := d.p1.Equal(d.p2); b != d.res {
			t.Errorf("Test Point.Equal failed p1: %v, p2: %v, res: %v, out: %v", d.p1, d.p2, d.res, b)
		}
	}
}

func TestPointSet(t *testing.T) {
	data := []Point{
		{X: 0, Y: 0},
		{X: 3, Y: 3},
		{X: 2, Y: 1},
		{X: 0, Y: 0},
		{X: -3, Y: 2},
		{X: 6, Y: -4},
		{X: -1, Y: -3},
	}
	p := Point{X: 0.5, Y: 0.5}
	for _, d := range data {
		if p.Set(d.X, d.Y); !p.Equal(d) {
			t.Errorf("Test Point.Set failed. wanted: %v, got:%v", d, p)
		}
	}
}

func TestPointPointDistance(t *testing.T) {
	data := []struct {
		p1, p2 Point
		exp    float32
	}{
		{p1: Point{X: 0, Y: 0}, p2: Point{X: 0, Y: 0}, exp: 0.0},
		{p1: Point{X: 0, Y: 0}, p2: Point{X: 3, Y: 4}, exp: 5.0},
		{p1: Point{X: -4, Y: -3}, p2: Point{X: 0, Y: 0}, exp: 5.0},
		{p1: Point{X: 2, Y: 4}, p2: Point{X: -1, Y: 0}, exp: 5.0},
		{p1: Point{X: -9, Y: -2}, p2: Point{X: -16, Y: -26}, exp: 25.0},
		{p1: Point{X: 161, Y: 240}, p2: Point{X: 322, Y: 480}, exp: 289.0},
	}
	for _, d := range data {
		if actual := d.p1.PointDistance(d.p2); actual != d.exp {
			t.Errorf("Test Point.PointDistance failed. p1: %v, p2: %v, wanted: %v, got: %v", d.p1, d.p2, d.exp, actual)
		}
	}
}

func TestPointProjectOnto(t *testing.T) {
	data := []struct {
		p1, p2, exp Point
	}{
		{p1: Point{X: 0, Y: 0}, p2: Point{X: 0, Y: 0}, exp: Point{X: 0, Y: 0}},
		{p1: Point{X: 5, Y: 5}, p2: Point{X: 5, Y: 5}, exp: Point{X: 5, Y: 5}},
		{p1: Point{X: 3, Y: 1}, p2: Point{X: 1, Y: 4}, exp: Point{X: 0.4117647, Y: 1.6470588}},
	}
	for _, d := range data {
		if actual := d.p1.ProjectOnto(d.p2); !actual.Equal(d.exp) {
			t.Errorf("Test Point.ProjectOnto failed. p1: %v, p2: %v, wanted: %v, got: %v", d.p1, d.p2, d.exp, actual)
		}
	}
}

func TestPointNormalize(t *testing.T) {
	data := []struct {
		p1, unit Point
		mag      float32
	}{
		{p1: Point{X: 0, Y: 0}, unit: Point{X: 0, Y: 0}, mag: 0.0},
		{p1: Point{X: 3, Y: 4}, unit: Point{X: 0.6, Y: 0.8}, mag: 5.0},
		{p1: Point{X: -3, Y: -4}, unit: Point{X: -0.6, Y: -0.8}, mag: 5.0},
	}
	for _, d := range data {
		if actualUnit, actualMag := d.p1.Normalize(); !actualUnit.Equal(d.unit) || actualMag != d.mag {
			t.Errorf("Test Point.Normalize failed. p1: %v, wanted: %v, %v, got: %v, %v", d.p1, d.unit, d.mag, actualUnit, actualMag)
		}
	}
}

type testContainer struct {
	bounds AABB
}

func (c testContainer) Contains(p Point) bool {
	if p.X <= c.bounds.Max.X && p.X >= c.bounds.Min.X && p.Y <= c.bounds.Max.Y && p.Y >= c.bounds.Min.Y {
		return true
	}
	return false
}

func TestPointWithin(t *testing.T) {
	data := []struct {
		box testContainer
		p   Point
		exp bool
	}{
		{
			box: testContainer{AABB{Min: Point{X: 0, Y: 0}, Max: Point{X: 0, Y: 0}}},
			p:   Point{X: 0, Y: 0},
			exp: true,
		},
		{
			box: testContainer{AABB{Min: Point{X: -5, Y: -5}, Max: Point{X: 5, Y: 5}}},
			p:   Point{X: 0, Y: 0},
			exp: true,
		},
		{
			box: testContainer{AABB{Min: Point{X: -5, Y: -5}, Max: Point{X: 5, Y: 5}}},
			p:   Point{X: 10, Y: 10},
			exp: false,
		},
	}
	for _, d := range data {
		if actual := d.p.Within(d.box); actual != d.exp {
			t.Errorf("Test Point.Within failed. box: %v, p: %v, want: %v, got: %v", d.box, d.p, d.exp, actual)
		}
	}
}

func TestLinePointSide(t *testing.T) {
	data := []struct {
		l   Line
		p   Point
		exp bool
	}{
		{
			l: Line{
				P1: Point{X: 0, Y: 0},
				P2: Point{X: 5, Y: 5},
			},
			p:   Point{X: -5, Y: 0},
			exp: true,
		},
		{
			l: Line{
				P1: Point{X: 0, Y: 0},
				P2: Point{X: 5, Y: 5},
			},
			p:   Point{X: 5, Y: 0},
			exp: false,
		},
		{
			l: Line{
				P1: Point{X: 0, Y: 0},
				P2: Point{X: 5, Y: 5},
			},
			p:   Point{X: 0, Y: 0},
			exp: false,
		},
	}
	for _, d := range data {
		if actual := d.l.PointSide(d.p); actual != d.exp {
			t.Errorf("Test Line.PointSide failed. l: %v, p: %v, want: %v, got: %v", d.l, d.p, d.exp, actual)
		}
	}
}

func TestLineAngle(t *testing.T) {
	data := []struct {
		l   Line
		exp float32
	}{
		{
			l: Line{
				P1: Point{X: 0, Y: 0},
				P2: Point{X: 0, Y: 0},
			},
			exp: 0,
		},
		{
			l: Line{
				P1: Point{X: 0, Y: 0},
				P2: Point{X: 3, Y: 4},
			},
			exp: -2.4980915,
		},
		{
			l: Line{
				P1: Point{X: 0, Y: 0},
				P2: Point{X: -3, Y: 4},
			},
			exp: 2.4980915,
		},
		{
			l: Line{
				P1: Point{X: 0, Y: 0},
				P2: Point{X: -3, Y: -4},
			},
			exp: 0.6435011,
		},
		{
			l: Line{
				P1: Point{X: 0, Y: 0},
				P2: Point{X: 3, Y: -4},
			},
			exp: -0.6435011,
		},
	}
	for _, d := range data {
		if res := d.l.Angle(); !FloatEqual(res, d.exp) {
			t.Errorf("Test Line.Angle failed. l: %v, wanted: %v, got: %v", d.l, d.exp, res)
		}
	}
}

func TestLineAngleDeg(t *testing.T) {
	data := []struct {
		l   Line
		exp float32
	}{
		{
			l: Line{
				P1: Point{X: 0, Y: 0},
				P2: Point{X: 0, Y: 0},
			},
			exp: 0,
		},
		{
			l: Line{
				P1: Point{X: 0, Y: 0},
				P2: Point{X: 0, Y: -1},
			},
			exp: 0,
		},
		{
			l: Line{
				P1: Point{X: 0, Y: 0},
				P2: Point{X: 1, Y: 0},
			},
			exp: 90,
		},
		{
			l: Line{
				P1: Point{X: 0, Y: 0},
				P2: Point{X: 0, Y: 1},
			},
			exp: 180,
		},
		{
			l: Line{
				P1: Point{X: 0, Y: 0},
				P2: Point{X: -1, Y: 0},
			},
			exp: 270,
		},
		{
			l: Line{
				P1: Point{X: 0, Y: 0},
				P2: Point{X: 2, Y: -2},
			},
			exp: 45,
		},
		{
			l: Line{
				P1: Point{X: 0, Y: 0},
				P2: Point{X: 2, Y: 2},
			},
			exp: 135,
		},
		{
			l: Line{
				P1: Point{X: 0, Y: 0},
				P2: Point{X: -0.7071068, Y: 0.7071067},
			},
			exp: 225,
		},
		{
			l: Line{
				P1: Point{X: 0, Y: 0},
				P2: Point{X: -2, Y: -2},
			},
			exp: 315,
		},
	}
	for _, d := range data {
		if res := d.l.AngleDeg(); !FloatEqual(res, d.exp) {
			t.Errorf("Test Line.AngleDeg failed. l: %v, wanted: %v, got: %v", d.l, d.exp, res)
		}
	}
}

func TestLinePointDistance(t *testing.T) {
	data := []struct {
		l   Line
		p   Point
		exp float32
	}{
		{
			l: Line{
				P1: Point{X: 0, Y: 0},
				P2: Point{X: 0, Y: 0},
			},
			p:   Point{X: 0, Y: 0},
			exp: 0,
		},
		{
			l: Line{
				P1: Point{X: 0, Y: 0},
				P2: Point{X: 8, Y: 8},
			},
			p:   Point{X: 3, Y: 4},
			exp: 0.70710677,
		},
		{
			l: Line{
				P1: Point{X: 0, Y: 0},
				P2: Point{X: 8, Y: 8},
			},
			p:   Point{X: 3, Y: 3},
			exp: 0,
		},
		{
			l: Line{
				P1: Point{X: 0, Y: 0},
				P2: Point{X: 8, Y: 8},
			},
			p:   Point{X: 13, Y: 10},
			exp: 5.3851647,
		},
		{
			l: Line{
				P1: Point{X: 0, Y: 0},
				P2: Point{X: 8, Y: 8},
			},
			p:   Point{X: -13, Y: -10},
			exp: 16.40122,
		},
	}
	for _, d := range data {
		if res := d.l.PointDistance(d.p); res != d.exp {
			t.Errorf("Test Line.PointDistance failed. l: %v, p: %v, wanted: %v, got: %v", d.l, d.p, d.exp, res)
		}
	}
}

func TestLineNormal(t *testing.T) {
	data := []struct {
		l   Line
		exp Point
	}{
		{
			l: Line{
				P1: Point{X: 0, Y: 0},
				P2: Point{X: 0, Y: 0},
			},
			exp: Point{X: 0, Y: 0},
		},
		{
			l: Line{
				P1: Point{X: 0, Y: 0},
				P2: Point{X: 0, Y: 8},
			},
			exp: Point{X: 1, Y: 0},
		},
		{
			l: Line{
				P1: Point{X: 0, Y: 0},
				P2: Point{X: 8, Y: 0},
			},
			exp: Point{X: 0, Y: -1},
		},
		{
			l: Line{
				P1: Point{X: 0, Y: 0},
				P2: Point{X: 8, Y: 8},
			},
			exp: Point{X: math.Sqrt(2) / 2, Y: -1 * math.Sqrt(2) / 2},
		},
	}
	for _, d := range data {
		if res := d.l.Normal(); !res.Equal(d.exp) {
			t.Errorf("Test Line.Normal failed. l: %v, wanted: %v, got: %v", d.l, d.exp, res)
		}
	}
}

func TestDotProduct(t *testing.T) {
	data := []struct {
		p1, p2 Point
		exp    float32
	}{
		{
			p1:  Point{X: 0, Y: 0},
			p2:  Point{X: 0, Y: 0},
			exp: 0,
		},
		{
			p1:  Point{X: 5, Y: 2},
			p2:  Point{X: 10, Y: 6},
			exp: 62,
		},
		{
			p1:  Point{X: 5, Y: 2},
			p2:  Point{X: -3, Y: 1},
			exp: -13,
		},
		{
			p1:  Point{X: -4, Y: -9},
			p2:  Point{X: -1, Y: 2},
			exp: -14,
		},
		{
			p1:  Point{X: 0, Y: 1},
			p2:  Point{X: 1, Y: 0},
			exp: 0,
		},
	}
	for _, d := range data {
		if res := DotProduct(d.p1, d.p2); !FloatEqual(res, d.exp) {
			t.Errorf("Test DotProduct failed. p1: %v, p2: %v, wanted: %v, got: %v", d.p1, d.p2, d.exp, res)
		}
	}
}

func TestLineTraceFraction(t *testing.T) {
	data := []struct {
		tracer, boundary Line
		exp              float32
	}{
		{
			tracer: Line{
				P1: Point{X: 0, Y: 0},
				P2: Point{X: 0, Y: 0},
			},
			boundary: Line{
				P1: Point{X: 0, Y: 0},
				P2: Point{X: 0, Y: 0},
			},
			exp: 1,
		},
		{
			tracer: Line{
				P1: Point{X: 0, Y: 0},
				P2: Point{X: 5, Y: 5},
			},
			boundary: Line{
				P1: Point{X: 1, Y: 0},
				P2: Point{X: 1, Y: 3},
			},
			exp: 0.2,
		},
		{
			tracer: Line{
				P1: Point{X: 1, Y: 1},
				P2: Point{X: 2, Y: 1},
			},
			boundary: Line{
				P1: Point{X: 0, Y: 0},
				P2: Point{X: 2, Y: 2},
			},
			exp: 0,
		},
	}
	for _, d := range data {
		if actual := LineTraceFraction(d.tracer, d.boundary); !FloatEqual(actual, d.exp) {
			t.Errorf("Test LineTraceFraction failed. tracer: %v, boundary: %v, wanted: %v, got: %v", d.tracer, d.boundary, d.exp, actual)
		}
	}
}

func TestLineTrace(t *testing.T) {
	data := []struct {
		tracer     Line
		boundaries []Line
		exp        Trace
	}{
		{
			tracer: Line{
				P1: Point{X: 0, Y: 0},
				P2: Point{X: 5, Y: 0},
			},
			boundaries: []Line{
				Line{
					P1: Point{X: 5, Y: 0},
					P2: Point{X: 5, Y: 5},
				},
				Line{
					P1: Point{X: 5, Y: 5},
					P2: Point{X: 0, Y: 5},
				},
				Line{
					P1: Point{X: 0, Y: 5},
					P2: Point{X: 0, Y: 0},
				},
				Line{
					P1: Point{X: -3, Y: -3},
					P2: Point{X: 3, Y: -3},
				},
			},
			exp: Trace{
				Fraction:    0,
				EndPosition: Point{X: 0, Y: 0},
				Line: Line{
					P1: Point{X: 0, Y: 5},
					P2: Point{X: 0, Y: 0},
				},
			},
		},
		{
			tracer: Line{
				P1: Point{X: 0, Y: 0},
				P2: Point{X: 5, Y: 0},
			},
			boundaries: []Line{
				Line{
					P1: Point{X: 5, Y: 5},
					P2: Point{X: 0, Y: 5},
				},
				Line{
					P1: Point{X: -3, Y: -3},
					P2: Point{X: 3, Y: -3},
				},
			},
			exp: Trace{
				Fraction:    1,
				EndPosition: Point{X: 5, Y: 0},
				Line: Line{
					P1: Point{X: 5, Y: 5},
					P2: Point{X: 0, Y: 5},
				},
			},
		},
		{
			tracer: Line{
				P1: Point{X: -3, Y: -3},
				P2: Point{X: 2, Y: 2},
			},
			boundaries: []Line{
				Line{
					P1: Point{X: 5, Y: 0},
					P2: Point{X: 5, Y: 5},
				},
				Line{
					P1: Point{X: 5, Y: 5},
					P2: Point{X: 0, Y: 5},
				},
				Line{
					P1: Point{X: 0, Y: 5},
					P2: Point{X: 0, Y: 0},
				},
				Line{
					P1: Point{X: -3, Y: -3},
					P2: Point{X: 3, Y: -3},
				},
			},
			exp: Trace{
				Fraction:    0,
				EndPosition: Point{X: -3, Y: -3},
				Line: Line{
					P1: Point{X: -3, Y: -3},
					P2: Point{X: 3, Y: -3},
				},
			},
		},
	}
	for _, d := range data {
		if actual := LineTrace(d.tracer, d.boundaries); !FloatEqual(actual.Fraction, d.exp.Fraction) || !actual.EndPosition.Equal(d.exp.EndPosition) || !actual.P1.Equal(d.exp.P1) || !actual.P2.Equal(d.exp.P2) {
			t.Errorf("Testing LineTrace failed. tracer: %v, boundaries: %v, wanted: %v, got: %v", d.tracer, d.boundaries, d.exp, actual)
		}
	}
}

func TestPointMethodChain(t *testing.T) {
	p1 := Point{X: 2, Y: 2}
	p2 := Point{X: 5, Y: 5}

	p1.Add(p2).Multiply(p2).Subtract(*p2.MultiplyScalar(2)).AddScalar(25).SubtractScalar(8)

	if p1.X != 42 {
		t.Errorf("Point method chain failed. p1.X should be 42, not %v", p1.X)
	}

	if p1.Y != 42 {
		t.Errorf("Point method chain failed. p1.Y should be 42, not %v", p1.Y)
	}
}

func TestMatrixSet(t *testing.T) {
	m := IdentityMatrix()
	for i := 0; i < 9; i++ {
		exp := float32(0)
		if i == 0 || i == 4 || i == 8 {
			exp = 1
		}
		if !FloatEqual(m.Val[i], exp) {
			t.Errorf("Identity was not the identity matrix. Index: %v\n Wanted: %v\n Got: %v\n", i, exp, m.Val[i])
		}
	}
	set := []float32{
		2, 3, 4,
		3, 4, 5,
		4, 5, 6,
	}
	m.Set(set)
	for i := 0; i < 9; i++ {
		if !FloatEqual(set[i], m.Val[i]) {
			t.Errorf("Set did not set properly. Index: %v\n Wanted: %v\n Got: %v\n", i, set[i], m.Val[i])
		}
	}
}

func TestMatrixMultiply(t *testing.T) {
	data := []struct {
		m1, m2, expected []float32
	}{
		{
			m1: []float32{
				0, 0, 0,
				0, 0, 0,
				0, 0, 0,
			},
			m2: []float32{
				0, 0, 0,
				0, 0, 0,
				0, 0, 0,
			},
			expected: []float32{
				0, 0, 0,
				0, 0, 0,
				0, 0, 0,
			},
		},
		{
			m1: []float32{
				1, 2, 3,
				4, 5, 6,
				7, 8, 9,
			},
			m2: []float32{
				1, 0, 0,
				0, 1, 0,
				0, 0, 1,
			},
			expected: []float32{
				1, 2, 3,
				4, 5, 6,
				7, 8, 9,
			},
		},
		{
			m1: []float32{
				0, 1, 2,
				2, 3, 4,
				3, 4, 2,
			},
			m2: []float32{
				1, 3, 2,
				2, 3, 2,
				2, 1, 1,
			},
			expected: []float32{
				12, 18, 18,
				12, 19, 20,
				5, 9, 10,
			},
		},
		{
			m1: []float32{
				1, 3, 2,
				2, 3, 2,
				2, 1, 1,
			},
			m2: []float32{
				0, 1, 2,
				2, 3, 4,
				3, 4, 2,
			},
			expected: []float32{
				6, 5, 4,
				16, 19, 14,
				15, 23, 16,
			},
		},
	}
	for _, d := range data {
		mat1, mat2 := IdentityMatrix(), IdentityMatrix()
		mat1.Set(d.m1)
		mat2.Set(d.m2)
		res := mat1.Multiply(mat2)
		for i := 0; i < 9; i++ {
			if !FloatEqual(res.Val[i], d.expected[i]) {
				t.Errorf("Multiplication was not correct.\n mat1: %v\n mat2: %v\n exp: %v\n res: %v\n index was: %v\ngot: %v\nwanted: %v\n", d.m1, d.m2, d.expected, res, i, res.Val[i], d.expected[i])
				return
			}
		}
	}
}

func TestMatrixTranslate(t *testing.T) {
	data := []struct {
		matrix, expected []float32
		transX, transY   float32
	}{
		{
			transX: 0,
			transY: 0,
			matrix: []float32{
				0, 1, 2,
				1, 2, 3,
				2, 3, 4,
			},
			expected: []float32{
				0, 1, 2,
				1, 2, 3,
				2, 3, 4,
			},
		},
		{
			transX: 15,
			transY: 15,
			matrix: []float32{
				0, 1, 2,
				1, 2, 3,
				2, 3, 4,
			},
			expected: []float32{
				0, 1, 2,
				1, 2, 3,
				17, 48, 79,
			},
		},
		{
			transX: -5,
			transY: -5,
			matrix: []float32{
				0, 1, 2,
				1, 2, 3,
				2, 3, 4,
			},
			expected: []float32{
				0, 1, 2,
				1, 2, 3,
				-3, -12, -21,
			},
		},
		{
			transX: 5,
			transY: -5,
			matrix: []float32{
				0, 1, 2,
				1, 2, 3,
				2, 3, 4,
			},
			expected: []float32{
				0, 1, 2,
				1, 2, 3,
				-3, -2, -1,
			},
		},
		{
			transX: -5,
			transY: 5,
			matrix: []float32{
				0, 1, 2,
				1, 2, 3,
				2, 3, 4,
			},
			expected: []float32{
				0, 1, 2,
				1, 2, 3,
				7, 8, 9,
			},
		},
	}
	for _, d := range data {
		mat1 := IdentityMatrix()
		mat1.Set(d.matrix)
		res := mat1.Translate(d.transX, d.transY)
		for i := 0; i < 9; i++ {
			if !FloatEqual(res.Val[i], d.expected[i]) {
				t.Errorf("Translation was not correct.\n matrix: %v\n exp: %v\n res: %v\n index was: %v\ngot: %v\nwanted: %v\n", d.matrix, d.expected, res, i, res.Val[i], d.expected[i])
				return
			}
		}
	}
}

func TestMatrixTranslatePoint(t *testing.T) {
	data := []struct {
		matrix, expected []float32
		point            Point
	}{
		{
			point: Point{
				X: 0,
				Y: 0,
			},
			matrix: []float32{
				0, 1, 2,
				1, 2, 3,
				2, 3, 4,
			},
			expected: []float32{
				0, 1, 2,
				1, 2, 3,
				2, 3, 4,
			},
		},
		{
			point: Point{
				X: 15,
				Y: 15,
			},
			matrix: []float32{
				0, 1, 2,
				1, 2, 3,
				2, 3, 4,
			},
			expected: []float32{
				0, 1, 2,
				1, 2, 3,
				17, 48, 79,
			},
		},
		{
			point: Point{
				X: -5,
				Y: -5,
			},
			matrix: []float32{
				0, 1, 2,
				1, 2, 3,
				2, 3, 4,
			},
			expected: []float32{
				0, 1, 2,
				1, 2, 3,
				-3, -12, -21,
			},
		},
		{
			point: Point{
				X: 5,
				Y: -5,
			},
			matrix: []float32{
				0, 1, 2,
				1, 2, 3,
				2, 3, 4,
			},
			expected: []float32{
				0, 1, 2,
				1, 2, 3,
				-3, -2, -1,
			},
		},
		{
			point: Point{
				X: -5,
				Y: 5,
			},
			matrix: []float32{
				0, 1, 2,
				1, 2, 3,
				2, 3, 4,
			},
			expected: []float32{
				0, 1, 2,
				1, 2, 3,
				7, 8, 9,
			},
		},
	}
	for _, d := range data {
		mat1 := IdentityMatrix()
		mat1.Set(d.matrix)
		res := mat1.TranslatePoint(d.point)
		for i := 0; i < 9; i++ {
			if !FloatEqual(res.Val[i], d.expected[i]) {
				t.Errorf("Translation was not correct.\n matrix: %v\n exp: %v\n res: %v\n index was: %v\ngot: %v\nwanted: %v\n", d.matrix, d.expected, res, i, res.Val[i], d.expected[i])
				return
			}
		}
	}
}

func TestMatrixScale(t *testing.T) {
	data := []struct {
		matrix, expected []float32
		scaleX, scaleY   float32
	}{
		{
			scaleX: 1,
			scaleY: 1,
			matrix: []float32{
				1, 2, 3,
				2, 2, 2,
				1, 3, 2,
			},
			expected: []float32{
				1, 2, 3,
				2, 2, 2,
				1, 3, 2,
			},
		},
		{
			scaleX: 0,
			scaleY: 0,
			matrix: []float32{
				1, 2, 3,
				2, 2, 2,
				1, 3, 2,
			},
			expected: []float32{
				0, 0, 0,
				0, 0, 0,
				1, 3, 2,
			},
		},
		{
			scaleX: 0,
			scaleY: 1,
			matrix: []float32{
				1, 2, 3,
				2, 2, 2,
				1, 3, 2,
			},
			expected: []float32{
				0, 0, 0,
				2, 2, 2,
				1, 3, 2,
			},
		},
		{
			scaleX: 1,
			scaleY: 0,
			matrix: []float32{
				1, 2, 3,
				2, 2, 2,
				1, 3, 2,
			},
			expected: []float32{
				1, 2, 3,
				0, 0, 0,
				1, 3, 2,
			},
		},
		{
			scaleX: 5,
			scaleY: 5,
			matrix: []float32{
				1, 2, 3,
				2, 2, 2,
				1, 3, 2,
			},
			expected: []float32{
				5, 10, 15,
				10, 10, 10,
				1, 3, 2,
			},
		},
		{
			scaleX: -5,
			scaleY: -5,
			matrix: []float32{
				1, 2, 3,
				2, 2, 2,
				1, 3, 2,
			},
			expected: []float32{
				-5, -10, -15,
				-10, -10, -10,
				1, 3, 2,
			},
		},
	}
	for _, d := range data {
		mat1 := IdentityMatrix()
		mat1.Set(d.matrix)
		res := mat1.Scale(d.scaleX, d.scaleY)
		for i := 0; i < 9; i++ {
			if !FloatEqual(res.Val[i], d.expected[i]) {
				t.Errorf("Scale was not correct.\n matrix: %v\n exp: %v\n res: %v\n index was: %v\ngot: %v\nwanted: %v\n", d.matrix, d.expected, res, i, res.Val[i], d.expected[i])
				return
			}
		}
	}
}

func TestMatrixScaleComponent(t *testing.T) {
	data := []struct {
		matrix               []float32
		expectedX, expectedY float32
	}{
		{
			matrix: []float32{
				0, 0, 0,
				0, 0, 0,
				0, 0, 0,
			},
			expectedX: 0,
			expectedY: 0,
		},
		{
			matrix: []float32{
				1, 2, 3,
				3, 1, 2,
				2, 3, 1,
			},
			expectedX: 1,
			expectedY: 1,
		},
		{
			matrix: []float32{
				5, 2, 3,
				6, -4, 3,
				1, 4, 2,
			},
			expectedX: 5,
			expectedY: -4,
		},
		{
			matrix: []float32{
				-5, 2, 3,
				6, -4, 3,
				1, 4, 2,
			},
			expectedX: -5,
			expectedY: -4,
		},
	}
	for _, d := range data {
		mat1 := IdentityMatrix()
		mat1.Set(d.matrix)
		if resX, resY := mat1.ScaleComponent(); !FloatEqual(resX, d.expectedX) || !FloatEqual(resY, d.expectedY) {
			t.Errorf("Scale Component did not match expected!\n ResX: %v\n expectedX: %v\n ResY: %v\n expectedY: %v\n", resX, d.expectedX, resY, d.expectedY)
			return
		}
	}
}

func TestMatrixTranslationComponent(t *testing.T) {
	data := []struct {
		matrix               []float32
		expectedX, expectedY float32
	}{
		{
			matrix: []float32{
				0, 0, 0,
				0, 0, 0,
				0, 0, 0,
			},
			expectedX: 0,
			expectedY: 0,
		},
		{
			matrix: []float32{
				1, 2, 3,
				4, 5, 6,
				7, 8, 9,
			},
			expectedX: 7,
			expectedY: 8,
		},
		{
			matrix: []float32{
				5, 2, 3,
				6, -4, 3,
				1, -4, 2,
			},
			expectedX: 1,
			expectedY: -4,
		},
		{
			matrix: []float32{
				-5, 2, 3,
				6, -4, 3,
				-1, 4, 2,
			},
			expectedX: -1,
			expectedY: 4,
		},
	}
	for _, d := range data {
		mat1 := IdentityMatrix()
		mat1.Set(d.matrix)
		if resX, resY := mat1.TranslationComponent(); !FloatEqual(resX, d.expectedX) || !FloatEqual(resY, d.expectedY) {
			t.Errorf("Translation Component did not match expected!\n Matrix: %v\n ResX: %v\n expectedX: %v\n ResY: %v\n expectedY: %v\n", d.matrix, resX, d.expectedX, resY, d.expectedY)
			return
		}
	}
}

func TestMatrixRotationComponent(t *testing.T) {
	data := []struct {
		matrix   []float32
		expected float32
	}{
		{
			matrix: []float32{
				0, 0, 0,
				0, 0, 0,
				0, 0, 0,
			},
			expected: 0,
		},
		{
			matrix: []float32{
				1, 2, 3,
				3, 1, 2,
				2, 3, 1,
			},
			expected: 63.43495,
		},
		{
			matrix: []float32{
				5, 2, 3,
				6, -4, 3,
				1, 4, 2,
			},
			expected: 21.801409,
		},
		{
			matrix: []float32{
				-5, 2, 3,
				6, -4, 3,
				1, 4, 2,
			},
			expected: 158.1986,
		},
	}
	for _, d := range data {
		mat1 := IdentityMatrix()
		mat1.Set(d.matrix)
		if res := mat1.RotationComponent(); !FloatEqual(res, d.expected) {
			t.Errorf("Rotation Component did not match expected!\nMatrix: %v\n Res: %v\n expected: %v\n", d.matrix, res, d.expected)
			return
		}
	}
}

func TestMatrixRotate(t *testing.T) {
	data := []struct {
		matrix, expected []float32
		rotation         float32
	}{
		{
			rotation: 0,
			matrix: []float32{
				1, 2, 3,
				6, 5, 4,
				7, 8, 9,
			},
			expected: []float32{
				1, 2, 3,
				6, 5, 4,
				7, 8, 9,
			},
		},
		{
			rotation: 35,
			matrix: []float32{
				1, 2, 3,
				6, 5, 4,
				7, 8, 9,
			},
			expected: []float32{
				4.2606106, 4.5061865, 4.751762,
				4.341336, 2.9486074, 1.5558789,
				7, 8, 9,
			},
		},
		{
			rotation: -35,
			matrix: []float32{
				1, 2, 3,
				6, 5, 4,
				7, 8, 9,
			},
			expected: []float32{
				-2.6223066, -1.2295781, 0.16315031,
				5.4884887, 5.2429132, 4.9973373,
				7, 8, 9,
			},
		},
		{
			rotation: 90,
			matrix: []float32{
				1, 2, 3,
				6, 5, 4,
				7, 8, 9,
			},
			expected: []float32{
				6, 5, 4,
				-1, -2, -3,
				7, 8, 9,
			},
		},
	}
	for _, d := range data {
		mat1 := IdentityMatrix()
		mat1.Set(d.matrix)
		res := mat1.Rotate(d.rotation)
		for i := 0; i < 9; i++ {
			if !FloatEqual(res.Val[i], d.expected[i]) {
				t.Errorf("Matrix Rotate did not return correct values.\n Matrix: %v\n Expected: %v\n Got: %v", d.matrix, d.expected, res.Val)
				return
			}
		}
	}
}

func TestMultiplyMatrixPoint(t *testing.T) {
	data := []struct {
		matrix []float32
		p, r   Point
	}{
		{
			matrix: []float32{
				1, 2, 3,
				4, 5, 6,
				7, 8, 0,
			},
			p: Point{
				X: 0,
				Y: 0,
			},
			r: Point{
				X: 7,
				Y: 8,
			},
		},
		{
			matrix: []float32{
				1, 2, 3,
				4, 5, 6,
				7, 8, 0,
			},
			p: Point{
				X: 5,
				Y: 5,
			},
			r: Point{
				X: 32,
				Y: 43,
			},
		},
		{
			matrix: []float32{
				1, 2, 3,
				4, 5, 6,
				7, 8, 0,
			},
			p: Point{
				X: -5,
				Y: -5,
			},
			r: Point{
				X: -18,
				Y: -27,
			},
		},
		{
			matrix: []float32{
				1, 2, 3,
				4, 5, 6,
				7, 8, 0,
			},
			p: Point{
				X: -5,
				Y: 5,
			},
			r: Point{
				X: 22,
				Y: 23,
			},
		},
		{
			matrix: []float32{
				1, 2, 3,
				4, 5, 6,
				7, 8, 0,
			},
			p: Point{
				X: 5,
				Y: -5,
			},
			r: Point{
				X: -8,
				Y: -7,
			},
		},
	}
	for _, d := range data {
		mat1 := IdentityMatrix()
		mat1.Set(d.matrix)
		if res := d.p.MultiplyMatrixVector(mat1); !res.Equal(d.r) {
			t.Errorf("MultiplyMatrixPoint did not return expected value. \nMatrix: %v\nPoint: %v\nWanted: %v\nGot: %v", d.matrix, d.p, d.r, res)
			return
		}
	}
}

func TestMultipleMatrixVector(t *testing.T) {
	data := []struct {
		matrix []float32
		p, r   []float32
	}{
		{
			matrix: []float32{
				1, 2, 3,
				4, 5, 6,
				7, 8, 0,
			},
			p: []float32{
				0,
				0,
			},
			r: []float32{
				7,
				8,
				0,
			},
		},
		{
			matrix: []float32{
				1, 2, 3,
				4, 5, 6,
				7, 8, 0,
			},
			p: []float32{
				5,
				5,
				0,
			},
			r: []float32{
				25,
				35,
				45,
			},
		},
		{
			matrix: []float32{
				1, 2, 3,
				4, 5, 6,
				7, 8, 0,
			},
			p: []float32{
				-5,
				-5,
				0,
			},
			r: []float32{
				-25,
				-35,
				-45,
			},
		},
		{
			matrix: []float32{
				1, 2, 3,
				4, 5, 6,
				7, 8, 0,
			},
			p: []float32{
				-5,
				5,
				0,
			},
			r: []float32{
				15,
				15,
				15,
			},
		},
		{
			matrix: []float32{
				1, 2, 3,
				4, 5, 6,
				7, 8, 0,
			},
			p: []float32{
				5,
				-5,
				0,
			},
			r: []float32{
				-15,
				-15,
				-15,
			},
		},
	}
	for _, d := range data {
		mat1 := IdentityMatrix()
		mat1.Set(d.matrix)
		res := MultiplyMatrixVector(mat1, d.p)
		for i := 0; i < 3; i++ {
			if !FloatEqual(res[i], d.r[i]) {
				t.Errorf("MultiplyMatrixVector did not return the correct value. \nMatrix %v\n With vector %v\nFor index %v\nWanted: %v\nGot: %v", mat1, d.p, i, res[i], d.r[i])
				return
			}
		}
	}
}
