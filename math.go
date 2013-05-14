package eng

import (
	gl "github.com/chsc/gogl/gl33"
	"math"
)

type Vector struct {
	X, Y, Z float32
}

func (v *Vector) Set(o *Vector) *Vector {
	v.X = o.X
	v.Y = o.Y
	v.Z = o.Z
	return v
}

func (v *Vector) Add(o *Vector) *Vector {
	v.X += o.X
	v.Y += o.Y
	v.Z += o.Z
	return v
}

func (v *Vector) Prj(m *Matrix) *Vector {
	w := 1 / (v.X*float32(m[3]) + v.Y*float32(m[7]) + v.Z*float32(m[11]) + float32(m[15]))
	v.X = w * (v.X*float32(m[0]) + v.Y*float32(m[4]) + v.Z*float32(m[8]) + float32(m[12]))
	v.Y = w * (v.X*float32(m[1]) + v.Y*float32(m[5]) + v.Z*float32(m[9]) + float32(m[13]))
	v.Z = w * (v.X*float32(m[2]) + v.Y*float32(m[6]) + v.Z*float32(m[10]) + float32(m[14]))
	return v
}

type Matrix [16]gl.Float

func NewMatrix() *Matrix {
	return &Matrix{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
}

func (m *Matrix) Set(o *Matrix) *Matrix {
	m[0] = o[0]
	m[1] = o[1]
	m[2] = o[2]
	m[3] = o[3]
	m[4] = o[4]
	m[5] = o[5]
	m[6] = o[6]
	m[7] = o[7]
	m[8] = o[8]
	m[9] = o[9]
	m[10] = o[10]
	m[11] = o[11]
	m[12] = o[12]
	m[13] = o[13]
	m[14] = o[14]
	m[15] = o[15]

	return m
}

func (m *Matrix) Identity() *Matrix {
	m[0] = 1
	m[1] = 0
	m[2] = 0
	m[3] = 0
	m[4] = 0
	m[5] = 1
	m[6] = 0
	m[7] = 0
	m[8] = 0
	m[9] = 0
	m[10] = 1
	m[11] = 0
	m[12] = 0
	m[13] = 0
	m[14] = 0
	m[15] = 1

	return m
}

func (m *Matrix) Mul(o *Matrix) *Matrix {
	a00 := m[0]
	a01 := m[1]
	a02 := m[2]
	a03 := m[3]
	a10 := m[4]
	a11 := m[5]
	a12 := m[6]
	a13 := m[7]
	a20 := m[8]
	a21 := m[9]
	a22 := m[10]
	a23 := m[11]
	a30 := m[12]
	a31 := m[13]
	a32 := m[14]
	a33 := m[15]

	b00 := o[0]
	b01 := o[1]
	b02 := o[2]
	b03 := o[3]
	b10 := o[4]
	b11 := o[5]
	b12 := o[6]
	b13 := o[7]
	b20 := o[8]
	b21 := o[9]
	b22 := o[10]
	b23 := o[11]
	b30 := o[12]
	b31 := o[13]
	b32 := o[14]
	b33 := o[15]

	m[0] = b00*a00 + b01*a10 + b02*a20 + b03*a30
	m[1] = b00*a01 + b01*a11 + b02*a21 + b03*a31
	m[2] = b00*a02 + b01*a12 + b02*a22 + b03*a32
	m[3] = b00*a03 + b01*a13 + b02*a23 + b03*a33
	m[4] = b10*a00 + b11*a10 + b12*a20 + b13*a30
	m[5] = b10*a01 + b11*a11 + b12*a21 + b13*a31
	m[6] = b10*a02 + b11*a12 + b12*a22 + b13*a32
	m[7] = b10*a03 + b11*a13 + b12*a23 + b13*a33
	m[8] = b20*a00 + b21*a10 + b22*a20 + b23*a30
	m[9] = b20*a01 + b21*a11 + b22*a21 + b23*a31
	m[10] = b20*a02 + b21*a12 + b22*a22 + b23*a32
	m[11] = b20*a03 + b21*a13 + b22*a23 + b23*a33
	m[12] = b30*a00 + b31*a10 + b32*a20 + b33*a30
	m[13] = b30*a01 + b31*a11 + b32*a21 + b33*a31
	m[14] = b30*a02 + b31*a12 + b32*a22 + b33*a32
	m[15] = b30*a03 + b31*a13 + b32*a23 + b33*a33

	return m
}

func cofactor(m0, m1, m2, m3, m4, m5, m6, m7, m8 gl.Float) gl.Float {
	return m0*(m4*m8-m5*m7) -
		m1*(m3*m8-m5*m6) +
		m2*(m3*m7-m4*m6)
}

func (m *Matrix) Inv() *Matrix {
	cf0 := cofactor(m[5], m[6], m[7], m[9], m[10], m[11], m[13], m[14], m[15])
	cf1 := cofactor(m[4], m[6], m[7], m[8], m[10], m[11], m[12], m[14], m[15])
	cf2 := cofactor(m[4], m[5], m[7], m[8], m[9], m[11], m[12], m[13], m[15])
	cf3 := cofactor(m[4], m[5], m[6], m[8], m[9], m[10], m[12], m[13], m[14])

	determinant := m[0]*cf0 - m[1]*cf1 + m[2]*cf2 - m[3]*cf3
	if math.Abs(float64(determinant)) <= 0.00001 {
		return m.Identity()
	}

	cf4 := cofactor(m[1], m[2], m[3], m[9], m[10], m[11], m[13], m[14], m[15])
	cf5 := cofactor(m[0], m[2], m[3], m[8], m[10], m[11], m[12], m[14], m[15])
	cf6 := cofactor(m[0], m[1], m[3], m[8], m[9], m[11], m[12], m[13], m[15])
	cf7 := cofactor(m[0], m[1], m[2], m[8], m[9], m[10], m[12], m[13], m[14])

	cf8 := cofactor(m[1], m[2], m[3], m[5], m[6], m[7], m[13], m[14], m[15])
	cf9 := cofactor(m[0], m[2], m[3], m[4], m[6], m[7], m[12], m[14], m[15])
	cf10 := cofactor(m[0], m[1], m[3], m[4], m[5], m[7], m[12], m[13], m[15])
	cf11 := cofactor(m[0], m[1], m[2], m[4], m[5], m[6], m[12], m[13], m[14])

	cf12 := cofactor(m[1], m[2], m[3], m[5], m[6], m[7], m[9], m[10], m[11])
	cf13 := cofactor(m[0], m[2], m[3], m[4], m[6], m[7], m[8], m[10], m[11])
	cf14 := cofactor(m[0], m[1], m[3], m[4], m[5], m[7], m[8], m[9], m[11])
	cf15 := cofactor(m[0], m[1], m[2], m[4], m[5], m[6], m[8], m[9], m[10])

	invDeterminant := 1.0 / determinant
	m[0] = invDeterminant * cf0
	m[1] = -invDeterminant * cf4
	m[2] = invDeterminant * cf8
	m[3] = -invDeterminant * cf12

	m[4] = -invDeterminant * cf1
	m[5] = invDeterminant * cf5
	m[6] = -invDeterminant * cf9
	m[7] = invDeterminant * cf13

	m[8] = invDeterminant * cf2
	m[9] = -invDeterminant * cf6
	m[10] = invDeterminant * cf10
	m[11] = -invDeterminant * cf14

	m[12] = -invDeterminant * cf3
	m[13] = invDeterminant * cf7
	m[14] = -invDeterminant * cf11
	m[15] = invDeterminant * cf15

	return m
}

func (m *Matrix) SetToOrtho(left, right, bottom, top, near, far float32) *Matrix {
	rl := right - left
	tb := top - bottom
	fn := far - near
	m[0] = gl.Float(2 / rl)
	m[1] = 0
	m[2] = 0
	m[3] = 0
	m[4] = 0
	m[5] = gl.Float(2 / tb)
	m[6] = 0
	m[7] = 0
	m[8] = 0
	m[9] = 0
	m[10] = gl.Float(-2 / fn)
	m[11] = 0
	m[12] = gl.Float(-(left + right) / rl)
	m[13] = gl.Float(-(top + bottom) / tb)
	m[14] = gl.Float(-(far + near) / fn)
	m[15] = 1
	return m
}

func (m *Matrix) SetToLookAt(eye, center, up *Vector) *Matrix {
	eyex := eye.X
	eyey := eye.Y
	eyez := eye.Z
	upx := up.X
	upy := up.Y
	upz := up.Z
	centerx := center.X
	centery := center.Y
	centerz := center.Z

	if eyex == centerx && eyey == centery && eyez == centerz {
		return m.Identity()
	}

	z0 := eyex - centerx
	z1 := eyey - centery
	z2 := eyez - centerz

	length := float32(1 / math.Sqrt(float64(z0*z0+z1*z1+z2*z2)))
	z0 *= length
	z1 *= length
	z2 *= length

	x0 := upy*z2 - upz*z1
	x1 := upz*z0 - upx*z2
	x2 := upx*z1 - upy*z0
	length = float32(math.Sqrt(float64(x0*x0 + x1*x1 + x2*x2)))
	if length == 0 {
		x0 = 0
		x1 = 0
		x2 = 0
	} else {
		length = 1 / length
		x0 *= length
		x1 *= length
		x2 *= length
	}

	y0 := z1*x2 - z2*x1
	y1 := z2*x0 - z0*x2
	y2 := z0*x1 - z1*x0

	length = float32(math.Sqrt(float64(y0*y0 + y1*y1 + y2*y2)))
	if length == 0 {
		y0 = 0
		y1 = 0
		y2 = 0
	} else {
		length = 1 / length
		y0 *= length
		y1 *= length
		y2 *= length
	}

	m[0] = gl.Float(x0)
	m[1] = gl.Float(y0)
	m[2] = gl.Float(z0)
	m[3] = 0
	m[4] = gl.Float(x1)
	m[5] = gl.Float(y1)
	m[6] = gl.Float(z1)
	m[7] = 0
	m[8] = gl.Float(x2)
	m[9] = gl.Float(y2)
	m[10] = gl.Float(z2)
	m[11] = 0
	m[12] = gl.Float(-(x0*eyex + x1*eyey + x2*eyez))
	m[13] = gl.Float(-(y0*eyex + y1*eyey + y2*eyez))
	m[14] = gl.Float(-(z0*eyex + z1*eyey + z2*eyez))
	m[15] = 1

	return m
}
