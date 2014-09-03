// Copyright 2014 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package engi

import (
	"log"
	"math"

	"github.com/ajhager/webgl"
)

const size = 10000

type Drawable interface {
	Texture() *webgl.Texture
	Width() float32
	Height() float32
	View() (float32, float32, float32, float32)
}

type Batch struct {
	drawing      bool
	lastTexture  *webgl.Texture
	vertices     []float32
	vertexVBO    *webgl.Buffer
	indices      []uint16
	indexVBO     *webgl.Buffer
	index        int
	shader       *webgl.Program
	inPosition   int
	inColor      int
	inTexCoords  int
	ufProjection *webgl.UniformLocation
	projX        float32
	projY        float32
}

func NewBatch(width, height float32) *Batch {
	batch := new(Batch)

	batch.shader = LoadShader(batchVert, batchFrag)
	batch.inPosition = gl.GetAttribLocation(batch.shader, "in_Position")
	batch.inColor = gl.GetAttribLocation(batch.shader, "in_Color")
	batch.inTexCoords = gl.GetAttribLocation(batch.shader, "in_TexCoords")
	batch.ufProjection = gl.GetUniformLocation(batch.shader, "uf_Projection")

	batch.vertices = make([]float32, 20*size)
	batch.indices = make([]uint16, 6*size)

	for i, j := 0, 0; i < size*6; i, j = i+6, j+4 {
		batch.indices[i+0] = uint16(j + 0)
		batch.indices[i+1] = uint16(j + 1)
		batch.indices[i+2] = uint16(j + 2)
		batch.indices[i+3] = uint16(j + 0)
		batch.indices[i+4] = uint16(j + 2)
		batch.indices[i+5] = uint16(j + 3)
	}

	batch.indexVBO = gl.CreateBuffer()
	batch.vertexVBO = gl.CreateBuffer()

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, batch.indexVBO)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, batch.indices, gl.STATIC_DRAW)

	gl.BindBuffer(gl.ARRAY_BUFFER, batch.vertexVBO)
	gl.BufferData(gl.ARRAY_BUFFER, batch.vertices, gl.DYNAMIC_DRAW)

	gl.EnableVertexAttribArray(batch.inPosition)
	gl.EnableVertexAttribArray(batch.inTexCoords)
	gl.EnableVertexAttribArray(batch.inColor)

	gl.VertexAttribPointer(batch.inPosition, 2, gl.FLOAT, false, 20, 0)
	gl.VertexAttribPointer(batch.inTexCoords, 2, gl.FLOAT, false, 20, 8)
	gl.VertexAttribPointer(batch.inColor, 4, gl.UNSIGNED_BYTE, true, 20, 16)

	batch.projX = width / 2
	batch.projY = height / 2

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	return batch
}

func (b *Batch) Begin() {
	if b.drawing {
		log.Fatal("Batch.End() must be called first")
	}
	b.drawing = true
	gl.UseProgram(b.shader)
}

func (b *Batch) End() {
	if !b.drawing {
		log.Fatal("Batch.Begin() must be called first")
	}
	if b.index > 0 {
		b.flush()
	}
	b.drawing = false

	b.lastTexture = nil
}

func (b *Batch) flush() {
	if b.lastTexture == nil {
		return
	}

	gl.BindTexture(gl.TEXTURE_2D, b.lastTexture)

	gl.Uniform2f(b.ufProjection, b.projX, b.projY)

	gl.BufferSubData(gl.ARRAY_BUFFER, 0, b.vertices)
	gl.DrawElements(gl.TRIANGLES, 6*b.index, gl.UNSIGNED_SHORT, 0)

	b.index = 0
}

func (b *Batch) SetProjection(width, height float32) {
	b.projX = width / 2
	b.projY = height / 2
}

func (b *Batch) Draw(r Drawable, x, y, originX, originY, scaleX, scaleY, rotation float32, color uint32, transparency float32) {
	if !b.drawing {
		log.Fatal("Batch.Begin() must be called first")
	}

	if r.Texture() != b.lastTexture {
		if b.lastTexture != nil {
			b.flush()
		}
		b.lastTexture = r.Texture()
	}

	x -= originX * r.Width()
	y -= originY * r.Height()

	originX = r.Width() * originX
	originY = r.Height() * originY

	worldOriginX := x + originX
	worldOriginY := y + originY
	fx := -originX
	fy := -originY
	fx2 := float32(r.Width()) - originX
	fy2 := float32(r.Height()) - originY

	if scaleX != 1 || scaleY != 1 {
		fx *= scaleX
		fy *= scaleY
		fx2 *= scaleX
		fy2 *= scaleY
	}

	p1x := fx
	p1y := fy
	p2x := fx
	p2y := fy2
	p3x := fx2
	p3y := fy2
	p4x := fx2
	p4y := fy

	var x1 float32
	var y1 float32
	var x2 float32
	var y2 float32
	var x3 float32
	var y3 float32
	var x4 float32
	var y4 float32

	if rotation != 0 {
		rot := float64(rotation * (math.Pi / 180.0))

		cos := float32(math.Cos(rot))
		sin := float32(math.Sin(rot))

		x1 = cos*p1x - sin*p1y
		y1 = sin*p1x + cos*p1y

		x2 = cos*p2x - sin*p2y
		y2 = sin*p2x + cos*p2y

		x3 = cos*p3x - sin*p3y
		y3 = sin*p3x + cos*p3y

		x4 = x1 + (x3 - x2)
		y4 = y3 - (y2 - y1)
	} else {
		x1 = p1x
		y1 = p1y

		x2 = p2x
		y2 = p2y

		x3 = p3x
		y3 = p3y

		x4 = p4x
		y4 = p4y
	}

	x1 += worldOriginX
	y1 += worldOriginY
	x2 += worldOriginX
	y2 += worldOriginY
	x3 += worldOriginX
	y3 += worldOriginY
	x4 += worldOriginX
	y4 += worldOriginY

	red := (color >> 16) & 0xFF
	green := ((color >> 8) & 0xFF) << 8
	blue := (color & 0xFF) << 16
	alpha := uint32(transparency*255.0) << 24
	tint := math.Float32frombits((alpha | blue | green | red) & 0xfeffffff)

	idx := b.index * 20

	u, v, u2, v2 := r.View()

	b.vertices[idx+0] = x1
	b.vertices[idx+1] = y1
	b.vertices[idx+2] = u
	b.vertices[idx+3] = v
	b.vertices[idx+4] = tint

	b.vertices[idx+5] = x4
	b.vertices[idx+6] = y4
	b.vertices[idx+7] = u2
	b.vertices[idx+8] = v
	b.vertices[idx+9] = tint

	b.vertices[idx+10] = x3
	b.vertices[idx+11] = y3
	b.vertices[idx+12] = u2
	b.vertices[idx+13] = v2
	b.vertices[idx+14] = tint

	b.vertices[idx+15] = x2
	b.vertices[idx+16] = y2
	b.vertices[idx+17] = u
	b.vertices[idx+18] = v2
	b.vertices[idx+19] = tint

	b.index += 1

	if b.index >= size {
		b.flush()
	}
}
