// Copyright 2013 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eng

import (
	"github.com/errcw/glow/gl/2.1/gl"
	"log"
	"math"
)

const size = 1000
const degToRad = math.Pi / 180

var batchVert = ` 
attribute vec4 in_Position;
attribute vec4 in_Color;
attribute vec2 in_TexCoords;

uniform mat4 uf_Matrix;

varying vec4 var_Color;
varying vec2 var_TexCoords;

void main() {
  var_Color = in_Color;
  var_TexCoords = in_TexCoords;
  gl_Position = uf_Matrix * in_Position;
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
	vertices         []float32
	vertexVBO        uint32
	colors           []float32
	colorVBO         uint32
	coords           []float32
	coordVBO         uint32
	index            int
	shader           *Shader
	customShader     *Shader
	combined         *Matrix
	projection       *Matrix
	transform        *Matrix
	color            *Color
	blendingDisabled bool
	blendSrcFunc     uint32
	blendDstFunc     uint32
	inPosition       uint32
	inColor          uint32
	inTexCoords      uint32
	ufMatrix         int32
}

func NewBatch() *Batch {
	batch := new(Batch)
	batch.vertices = make([]float32, 2*size)
	batch.colors = make([]float32, 4*size)
	batch.coords = make([]float32, 2*size)
	batch.shader = NewShader(batchVert, batchFrag)
	batch.inPosition = batch.shader.GetAttrib("in_Position")
	batch.inColor = batch.shader.GetAttrib("in_Color")
	batch.inTexCoords = batch.shader.GetAttrib("in_TexCoords")
	batch.ufMatrix = batch.shader.GetUniform("uf_Matrix")

	gl.GenBuffers(1, &batch.vertexVBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, batch.vertexVBO)
	gl.BufferData(gl.ARRAY_BUFFER, 2*4*4*size, gl.Ptr(batch.vertices), gl.DYNAMIC_DRAW)

	gl.GenBuffers(1, &batch.colorVBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, batch.colorVBO)
	gl.BufferData(gl.ARRAY_BUFFER, 4*4*size, gl.Ptr(batch.colors), gl.DYNAMIC_DRAW)

	gl.GenBuffers(1, &batch.coordVBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, batch.coordVBO)
	gl.BufferData(gl.ARRAY_BUFFER, 2*4*size, gl.Ptr(batch.coords), gl.DYNAMIC_DRAW)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	batch.combined = NewMatrix()
	batch.transform = NewMatrix()
	batch.projection = NewMatrix().SetToOrtho(0, float32(Width()), float32(Height()), 0, 0, 1)
	batch.color = NewColor(1, 1, 1)
	batch.blendingDisabled = false
	batch.blendSrcFunc = gl.SRC_ALPHA
	batch.blendDstFunc = gl.ONE_MINUS_SRC_ALPHA

	return batch
}

// Begin calculates the combined matrix and sets up rendering. This
// must be called before calling Draw.
func (b *Batch) Begin() {
	if b.drawing {
		log.Fatal("Batch.End() must be called first")
	}
	b.combined.Set(b.projection).Mul(b.transform)
	b.drawing = true
	shader := b.shader
	if b.customShader != nil {
		shader = b.customShader
	}

	shader.Bind()
}

func (b *Batch) DrawVerts(r *Region, verts [8]float32, color *Color) {
	if !b.drawing {
		log.Fatal("Batch.Begin() must be called first")
	}

	if r.texture != b.lastTexture {
		b.flush()
		b.lastTexture = r.texture
	}

	c := b.color
	if color != nil {
		c = color
	}

	vi := b.index * 2
	b.vertices[vi+0] = verts[0]
	b.vertices[vi+1] = verts[1]
	b.vertices[vi+2] = verts[2]
	b.vertices[vi+3] = verts[3]
	b.vertices[vi+4] = verts[4]
	b.vertices[vi+5] = verts[5]
	b.vertices[vi+6] = verts[6]
	b.vertices[vi+7] = verts[7]

	b.colors[b.index+0] = c.R
	b.colors[b.index+1] = c.G
	b.colors[b.index+2] = c.B
	b.colors[b.index+3] = c.A
	b.colors[b.index+4] = c.R
	b.colors[b.index+5] = c.G
	b.colors[b.index+6] = c.B
	b.colors[b.index+7] = c.A
	b.colors[b.index+8] = c.R
	b.colors[b.index+9] = c.G
	b.colors[b.index+10] = c.B
	b.colors[b.index+11] = c.A
	b.colors[b.index+12] = c.R
	b.colors[b.index+13] = c.G
	b.colors[b.index+14] = c.B
	b.colors[b.index+15] = c.A

	b.coords[b.index+0] = r.u
	b.coords[b.index+1] = r.v
	b.coords[b.index+2] = r.u
	b.coords[b.index+3] = r.v2
	b.coords[b.index+4] = r.u2
	b.coords[b.index+5] = r.v2
	b.coords[b.index+6] = r.u2
	b.coords[b.index+7] = r.v

	b.index += 4

	if b.index >= size {
		b.flush()
	}
}

// Draw renders a Region with its top left corner at x, y. Scaling and
// rotation will be with respect to the origin. If color is nil, the
// current batch color will be used. If the backing texture in the
// region is different than the last rendered region, any pending
// geometry will be flushed. Switching textures is a relatively
// expensive operation.
func (b *Batch) Draw(r *Region, x, y, originX, originY, scaleX, scaleY, rotation float32, color *Color) {
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

	c := b.color
	if color != nil {
		c = color
	}

	vi := b.index * 2
	b.vertices[vi+0] = x1
	b.vertices[vi+1] = y1
	b.vertices[vi+2] = x2
	b.vertices[vi+3] = y2
	b.vertices[vi+4] = x3
	b.vertices[vi+5] = y3
	b.vertices[vi+6] = x4
	b.vertices[vi+7] = y4

	ci := b.index * 4
	b.colors[ci+0] = c.R
	b.colors[ci+1] = c.G
	b.colors[ci+2] = c.B
	b.colors[ci+3] = c.A
	b.colors[ci+4] = c.R
	b.colors[ci+5] = c.G
	b.colors[ci+6] = c.B
	b.colors[ci+7] = c.A
	b.colors[ci+8] = c.R
	b.colors[ci+9] = c.G
	b.colors[ci+10] = c.B
	b.colors[ci+11] = c.A
	b.colors[ci+12] = c.R
	b.colors[ci+13] = c.G
	b.colors[ci+14] = c.B
	b.colors[ci+15] = c.A

	b.coords[vi+0] = r.u
	b.coords[vi+1] = r.v
	b.coords[vi+2] = r.u
	b.coords[vi+3] = r.v2
	b.coords[vi+4] = r.u2
	b.coords[vi+5] = r.v2
	b.coords[vi+6] = r.u2
	b.coords[vi+7] = r.v

	b.index += 4

	if b.index >= size {
		b.flush()
	}
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
		gl.Disable(gl.BLEND)
	}
	b.drawing = false

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.UseProgram(0)

	b.lastTexture = nil
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
	b.color.R = color.R
	b.color.G = color.G
	b.color.B = color.B
	b.color.A = color.A
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
// This is often used with a Camera.
func (b *Batch) SetProjection(m *Matrix) {
	b.projection.Set(m)
}

func (b *Batch) flush() {
	if b.lastTexture == nil {
		return
	}

	if b.blendingDisabled {
		gl.Disable(gl.BLEND)
	} else {
		gl.Enable(gl.BLEND)
		gl.BlendFunc(b.blendSrcFunc, b.blendDstFunc)
	}

	gl.Enable(gl.TEXTURE_2D)
	gl.ActiveTexture(gl.TEXTURE0)
	b.lastTexture.Bind()

	gl.UniformMatrix4fv(b.ufMatrix, 1, false, &b.combined[0])

	gl.BindBuffer(gl.ARRAY_BUFFER, b.vertexVBO)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, 4*2*b.index, gl.Ptr(b.vertices))
	gl.EnableVertexAttribArray(b.inPosition)
	gl.VertexAttribPointer(b.inPosition, 2, gl.FLOAT, false, 0, gl.PtrOffset(0))

	gl.BindBuffer(gl.ARRAY_BUFFER, b.colorVBO)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, 4*4*b.index, gl.Ptr(b.colors))
	gl.EnableVertexAttribArray(b.inColor)
	gl.VertexAttribPointer(b.inColor, 4, gl.FLOAT, false, 0, gl.PtrOffset(0))

	gl.BindBuffer(gl.ARRAY_BUFFER, b.coordVBO)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, 4*2*b.index, gl.Ptr(b.coords))
	gl.EnableVertexAttribArray(b.inTexCoords)
	gl.VertexAttribPointer(b.inTexCoords, 2, gl.FLOAT, false, 0, gl.PtrOffset(0))

	gl.DrawArrays(gl.QUADS, 0, int32(b.index))

	b.index = 0
}
