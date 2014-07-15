// Copyright 2014 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eng

import (
	"log"
	"math"
)

const size = 1000
const degToRad = math.Pi / 180

var batchVert = ` 
attribute vec4 in_Position;
attribute vec4 in_Color;
attribute vec2 in_TexCoords;

uniform vec2 uf_Projection;

varying vec4 var_Color;
varying vec2 var_TexCoords;

void main() {
  var_Color = in_Color;
  var_TexCoords = in_TexCoords;
	gl_Position = vec4(in_Position.x / uf_Projection.x - 1.0,
										 in_Position.y / -uf_Projection.y + 1.0,
										 0.0, 1.0);
}`

var batchFrag = `
varying vec4 var_Color;
varying vec2 var_TexCoords;

uniform sampler2D uf_Texture;

void main (void) {
  gl_FragColor = var_Color * texture2D(uf_Texture, var_TexCoords);
}`

// A Batch allows geometry to be efficiently rendered by buffering
// render calls and sending them all at once.
type Batch struct {
	drawing          bool
	lastTexture      *Texture
	color            float32
	vertices         []float32
	vertexVBO        uint32
	indices          []uint16
	indexVBO         uint32
	index            int
	shader           *Shader
	customShader     *Shader
	blendingDisabled bool
	blendSrcFunc     uint32
	blendDstFunc     uint32
	inPosition       uint32
	inColor          uint32
	inTexCoords      uint32
	ufProjection     int32
	projX            float32
	projY            float32
}

func NewBatch() *Batch {
	batch := new(Batch)

	batch.shader = NewShader(batchVert, batchFrag)
	batch.inPosition = batch.shader.GetAttrib("in_Position")
	batch.inColor = batch.shader.GetAttrib("in_Color")
	batch.inTexCoords = batch.shader.GetAttrib("in_TexCoords")
	batch.ufProjection = batch.shader.GetUniform("uf_Projection")

	batch.color = NewColorA(1, 1, 1, 1).FloatBits()

	batch.vertices = make([]float32, 20*size)
	batch.indices = make([]uint16, 6*size)

	for i, j := 0, 0; i < size*6; i, j = i+6, j+4 {
		batch.indices[i+0] = uint16(j + 0)
		batch.indices[i+1] = uint16(j + 1)
		batch.indices[i+2] = uint16(j + 2)
		batch.indices[i+3] = uint16(j + 2)
		batch.indices[i+4] = uint16(j + 1)
		batch.indices[i+5] = uint16(j + 3)
	}

	GL.GenBuffers(1, &batch.indexVBO)
	GL.GenBuffers(1, &batch.vertexVBO)

	GL.BindBuffer(GL.ELEMENT_ARRAY_BUFFER, batch.indexVBO)
	GL.BufferData(GL.ELEMENT_ARRAY_BUFFER, 6*2*size, batch.indices, GL.STATIC_DRAW)

	GL.BindBuffer(GL.ARRAY_BUFFER, batch.vertexVBO)
	GL.BufferData(GL.ARRAY_BUFFER, 20*4*size, batch.vertices, GL.DYNAMIC_DRAW)

	batch.projX = float32(Width()) / 2
	batch.projY = float32(Height()) / 2

	batch.blendingDisabled = false
	batch.blendSrcFunc = GL.SRC_ALPHA
	batch.blendDstFunc = GL.ONE_MINUS_SRC_ALPHA

	return batch
}

// Begin calculates the combined matrix and sets up rendering. This
// must be called before calling Draw.
func (b *Batch) Begin() {
	if b.drawing {
		log.Fatal("Batch.End() must be called first")
	}
	b.drawing = true
	shader := b.shader
	if b.customShader != nil {
		shader = b.customShader
	}

	shader.Bind()
}

// End finishes up rendering and flushes any remaining geometry to the
// gpu. This must be called after a called to Begin.
func (b *Batch) End() {
	if !b.drawing {
		log.Fatal("Batch.Begin() must be called first")
	}
	if b.index > 0 {
		b.flush()
	}
	if !b.blendingDisabled {
		GL.Disable(GL.BLEND)
	}
	b.drawing = false

	GL.BindBuffer(GL.ARRAY_BUFFER, 0)
	GL.UseProgram(0)

	b.lastTexture = nil
}

func (b *Batch) flush() {
	if b.lastTexture == nil {
		return
	}

	if b.blendingDisabled {
		GL.Disable(GL.BLEND)
	} else {
		GL.Enable(GL.BLEND)
		GL.BlendFunc(b.blendSrcFunc, b.blendDstFunc)
	}

	GL.Enable(GL.TEXTURE_2D)
	GL.ActiveTexture(GL.TEXTURE0)
	b.lastTexture.Bind()

	GL.Uniform2f(b.ufProjection, b.projX, b.projY)

	GL.BindBuffer(GL.ARRAY_BUFFER, b.vertexVBO)
	GL.BufferSubData(GL.ARRAY_BUFFER, 0, 20*4*b.index, b.vertices)

	GL.EnableVertexAttribArray(b.inPosition)
	GL.EnableVertexAttribArray(b.inTexCoords)
	GL.EnableVertexAttribArray(b.inColor)

	GL.VertexAttribPointer(b.inPosition, 2, GL.FLOAT, false, 20, 0)
	GL.VertexAttribPointer(b.inTexCoords, 2, GL.FLOAT, false, 20, 8)
	GL.VertexAttribPointer(b.inColor, 4, GL.UNSIGNED_BYTE, true, 20, 16)

	GL.BindBuffer(GL.ELEMENT_ARRAY_BUFFER, b.indexVBO)
	GL.DrawElements(GL.TRIANGLES, int32(6*b.index), GL.UNSIGNED_SHORT, 0)

	b.index = 0
}

// Draw renders a Region with its top left corner at x, y. Scaling and
// rotation will be with respect to the origin. If color is nil, the
// current batch color will be used. If the backing texture in the
// region is different than the last rendered region, any pending
// geometry will be flushed. Switching textures is a relatively
// expensive operation.
func (b *Batch) Draw(r *Region, x, y, originX, originY, scaleX, scaleY, rotation float32) {
	if !b.drawing {
		log.Fatal("Batch.Begin() must be called first")
	}

	if r.texture != b.lastTexture {
		b.flush()
		b.lastTexture = r.texture
	}

	worldOriginX := x + originX
	worldOriginY := y + originY
	fx := -originX
	fy := -originY
	fx2 := float32(r.width) - originX
	fy2 := float32(r.height) - originY

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
		rot := float64(rotation * degToRad)

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

	idx := b.index * 20

	b.vertices[idx+0] = x1
	b.vertices[idx+1] = y1
	b.vertices[idx+2] = r.u
	b.vertices[idx+3] = r.v
	b.vertices[idx+4] = b.color

	b.vertices[idx+5] = x4
	b.vertices[idx+6] = y4
	b.vertices[idx+7] = r.u2
	b.vertices[idx+8] = r.v
	b.vertices[idx+9] = b.color

	b.vertices[idx+10] = x2
	b.vertices[idx+11] = y2
	b.vertices[idx+12] = r.u
	b.vertices[idx+13] = r.v2
	b.vertices[idx+14] = b.color

	b.vertices[idx+15] = x3
	b.vertices[idx+16] = y3
	b.vertices[idx+17] = r.u2
	b.vertices[idx+18] = r.v2
	b.vertices[idx+19] = b.color

	b.index += 1

	if b.index >= size {
		b.flush()
	}
}

// SetBlending will toggle blending for rendering on the batch.
// Blending is a relatively expensive operation and should be disabled
// if your goemetry is opaque.
func (b *Batch) SetBlending(v bool) {
	if v != b.blendingDisabled {
		b.flush()
		b.blendingDisabled = !b.blendingDisabled
	}
}

// SetBlendFunc sets the opengl src and dst blending functions. The
// default is gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA which will render
// any alpha channel in your textures as blank. Calling this will
// flush any pending geometry to the gpu.
func (b *Batch) SetBlendFunc(src, dst uint32) {
	b.flush()
	b.blendSrcFunc = src
	b.blendDstFunc = dst
}

// SetColor changes the current batch rendering tint. This defaults to white.
func (b *Batch) SetColor(color *Color) {
	b.color = color.FloatBits()
}

// SetShader changes the shader used to rendering geometry. If the
// passed in shader == nil, the batch will go back to using its
// default shader. The shader should name the incoming vertex
// position, color, and texture coordinates to 'in_Position',
// 'in_Color', and 'in_TexCoords' respectively. The transform projection
// matrix will be passed in as 'uf_Matrix'.
func (b *Batch) SetShader(shader *Shader) {
	b.customShader = shader
}

// SetProjection allows for setting the projection matrix manually.
func (b *Batch) SetProjection(x, y float32) {
	b.projX = x
	b.projY = y
}
