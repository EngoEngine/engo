package eng

import (
	gl "github.com/chsc/gogl/gl21"
)

const (
	M00 = 0
	M01 = 4
	M02 = 8
	M03 = 12
	M10 = 1
	M11 = 5
	M12 = 9
	M13 = 13
	M20 = 2
	M21 = 6
	M22 = 10
	M23 = 14
	M30 = 3
	M31 = 7
	M32 = 11
	M33 = 15
)

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
	*m = *o
	return m
}

func (m *Matrix) SetOrtho(x, y, width, height float32) *Matrix {
	m.setToOrtho(x, x+width, y, y+height, 0, 1)
	return m
}

func (m *Matrix) setToOrtho(left, right, bottom, top, near, far float32) *Matrix {
	//this.idt();

	x_orth := 2 / (right - left)
	y_orth := 2 / (top - bottom)
	z_orth := -2 / (far - near)

	tx := -(right + left) / (right - left)
	ty := -(top + bottom) / (top - bottom)
	tz := -(far + near) / (far - near)

	m[M00] = gl.Float(x_orth)
	m[M10] = 0
	m[M20] = 0
	m[M30] = 0
	m[M01] = 0
	m[M11] = gl.Float(y_orth)
	m[M21] = 0
	m[M31] = 0
	m[M02] = 0
	m[M12] = 0
	m[M22] = gl.Float(z_orth)
	m[M32] = 0
	m[M03] = gl.Float(tx)
	m[M13] = gl.Float(ty)
	m[M23] = gl.Float(tz)
	m[M33] = 1

	return m
}

func (m *Matrix) Mul(o *Matrix) *Matrix {
	m00 := m[M00]
	m01 := m[M01]
	m02 := m[M02]
	m03 := m[M03]
	m10 := m[M10]
	m11 := m[M11]
	m12 := m[M12]
	m13 := m[M13]
	m20 := m[M20]
	m21 := m[M21]
	m22 := m[M22]
	m23 := m[M23]
	m30 := m[M30]
	m31 := m[M31]
	m32 := m[M32]
	m33 := m[M33]

	m[M00] = m00*o[M00] + m01*o[M10] + m02*o[M20] + m03*o[M30]
	m[M01] = m00*o[M01] + m01*o[M11] + m02*o[M21] + m03*o[M31]
	m[M02] = m00*o[M02] + m01*o[M12] + m02*o[M22] + m03*o[M32]
	m[M03] = m00*o[M03] + m01*o[M13] + m02*o[M23] + m03*o[M33]
	m[M10] = m10*o[M00] + m11*o[M10] + m12*o[M20] + m13*o[M30]
	m[M11] = m10*o[M01] + m11*o[M11] + m12*o[M21] + m13*o[M31]
	m[M12] = m10*o[M02] + m11*o[M12] + m12*o[M22] + m13*o[M32]
	m[M13] = m10*o[M03] + m11*o[M13] + m12*o[M23] + m13*o[M33]
	m[M20] = m20*o[M00] + m21*o[M10] + m22*o[M20] + m23*o[M30]
	m[M21] = m20*o[M01] + m21*o[M11] + m22*o[M21] + m23*o[M31]
	m[M22] = m20*o[M02] + m21*o[M12] + m22*o[M22] + m23*o[M32]
	m[M23] = m20*o[M03] + m21*o[M13] + m22*o[M23] + m23*o[M33]
	m[M30] = m30*o[M00] + m31*o[M10] + m32*o[M20] + m33*o[M30]
	m[M31] = m30*o[M01] + m31*o[M11] + m32*o[M21] + m33*o[M31]
	m[M32] = m30*o[M02] + m31*o[M12] + m32*o[M22] + m33*o[M32]
	m[M33] = m30*o[M03] + m31*o[M13] + m32*o[M23] + m33*o[M33]

	return m
}
