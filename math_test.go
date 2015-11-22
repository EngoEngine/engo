package engi

import "testing"

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
