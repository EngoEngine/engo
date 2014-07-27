// Copyright 2014 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package engi

import (
	"encoding/json"
	"log"
	"math"
	"strings"
)

const size = 10000
const degToRad = math.Pi / 180

var batchVert = ` 
attribute vec4 in_Position;
attribute vec4 in_Color;
attribute vec2 in_TexCoords;

uniform vec2 uf_Projection;

varying vec4 var_Color;
varying vec2 var_TexCoords;

const vec2 center = vec2(-1.0, 1.0);

void main() {
  var_Color = in_Color;
  var_TexCoords = in_TexCoords;
	gl_Position = vec4(in_Position.x / uf_Projection.x + center.x,
										 in_Position.y / -uf_Projection.y + center.y,
										 0.0, 1.0);
}`

var batchFrag = `
#ifdef GL_ES
#define LOWP lowp
precision mediump float;
#else
#define LOWP
#endif

varying vec4 var_Color;
varying vec2 var_TexCoords;

uniform sampler2D uf_Texture;

void main (void) {
  gl_FragColor = var_Color * texture2D(uf_Texture, var_TexCoords);
}`

// A Shader abstracts the loading, compiling, and linking of shader
// programs, which can directly modify the rendering of vertices and pixels.
type Shader struct {
	id *ProgramObject
}

// NewShader takes the source of a vertex and fragment shader and
// returns a compiled and linked shader program.
func NewShader(vertSrc, fragSrc string) *Shader {
	vertShader := GL.CreateShader(GL.VERTEX_SHADER)
	GL.ShaderSource(vertShader, vertSrc)
	GL.CompileShader(vertShader)
	defer GL.DeleteShader(vertShader)

	fragShader := GL.CreateShader(GL.FRAGMENT_SHADER)
	GL.ShaderSource(fragShader, fragSrc)
	GL.CompileShader(fragShader)
	defer GL.DeleteShader(fragShader)

	program := GL.CreateProgram()
	GL.AttachShader(program, vertShader)
	GL.AttachShader(program, fragShader)

	GL.LinkProgram(program)

	return &Shader{program}
}

// Bind turns the shader on to be used during rendering.
func (s *Shader) Bind() {
	GL.UseProgram(s.id)
}

// GetUniform returns the location of the named uniform.
func (s *Shader) GetUniform(uniform string) *UniformObject {
	return GL.GetUniformLocation(s.id, uniform)
}

// GetAttrib returns the location of the named attribute.
func (s *Shader) GetAttrib(attrib string) int {
	return GL.GetAttribLocation(s.id, attrib)
}

// A region represents a portion of a texture that can be rendered
// using a Batch.
type Region struct {
	texture       *Texture
	u, v          float32
	u2, v2        float32
	width, height float32
}

// NewRegion constructs an image from the rectangle x, y, w, h on the
// given texture.
func NewRegion(texture *Texture, x, y, w, h int) *Region {
	invTexWidth := 1.0 / float32(texture.Width())
	invTexHeight := 1.0 / float32(texture.Height())

	u := float32(x) * invTexWidth
	v := float32(y) * invTexHeight
	u2 := float32(x+w) * invTexWidth
	v2 := float32(y+h) * invTexHeight
	width := float32(math.Abs(float64(w)))
	height := float32(math.Abs(float64(h)))

	return &Region{texture, u, v, u2, v2, width, height}
}

// NewRegionFull returns a region that covers the entire texture.
func NewRegionFull(texture *Texture) *Region {
	return NewRegion(texture, 0, 0, int(texture.Width()), int(texture.Height()))
}

// Flip will swap the region's image on the x and/or y axes.
func (r *Region) Flip(x, y bool) {
	if x {
		tmp := r.u
		r.u = r.u2
		r.u2 = tmp
	}
	if y {
		tmp := r.v
		r.v = r.v2
		r.v2 = tmp
	}
}

func (r *Region) Width() float32 {
	return float32(r.width)
}

func (r *Region) Height() float32 {
	return float32(r.height)
}

// A Texture wraps an opengl texture and is mostly used for loading
// images and constructing Regions.
type Texture struct {
	id        *TextureObject
	width     int
	height    int
	minFilter int
	maxFilter int
	uWrap     int
	vWrap     int
}

func NewTexture(img Image) *Texture {
	id := GL.CreateTexture()

	GL.BindTexture(GL.TEXTURE_2D, id)

	GL.TexParameteri(GL.TEXTURE_2D, GL.TEXTURE_WRAP_S, GL.CLAMP_TO_EDGE)
	GL.TexParameteri(GL.TEXTURE_2D, GL.TEXTURE_WRAP_T, GL.CLAMP_TO_EDGE)
	GL.TexParameteri(GL.TEXTURE_2D, GL.TEXTURE_MIN_FILTER, GL.LINEAR)
	GL.TexParameteri(GL.TEXTURE_2D, GL.TEXTURE_MAG_FILTER, GL.NEAREST)

	if img.Data() == nil {
		panic("Texture image data is nil.")
	}

	GL.TexImage2D(GL.TEXTURE_2D, 0, GL.RGBA, img.Width(), img.Height(), 0, GL.RGBA, GL.UNSIGNED_BYTE, img.Data())

	return &Texture{id, img.Width(), img.Height(), GL.LINEAR, GL.LINEAR, GL.CLAMP_TO_EDGE, GL.CLAMP_TO_EDGE}
}

// Split creates Regions from every width, height rect going from left
// to right, then down. This is useful for simple images with uniform cells.
func (t *Texture) Split(w, h int) []*Region {
	x := 0
	y := 0
	width := int(t.Width())
	height := int(t.Height())

	rows := height / h
	cols := width / w

	startX := x
	tiles := make([]*Region, 0)
	for row := 0; row < rows; row++ {
		x = startX
		for col := 0; col < cols; col++ {
			tiles = append(tiles, NewRegion(t, x, y, w, h))
			x += w
		}
		y += h
	}

	return tiles
}

func (t *Texture) Unpack(path string) map[string]*Region {
	regions := make(map[string]*Region)

	var data interface{}
	err := json.Unmarshal([]byte(path), &data)
	if err != nil {
		log.Fatal(err)
	}

	root := data.(map[string]interface{})
	frames := root["frames"].([]interface{})
	for _, frameData := range frames {
		frame := frameData.(map[string]interface{})
		name := strings.Split(frame["filename"].(string), ".")[0]
		rect := frame["frame"].(map[string]interface{})
		x := int(rect["x"].(float64))
		y := int(rect["y"].(float64))
		w := int(rect["w"].(float64))
		h := int(rect["h"].(float64))
		regions[name] = NewRegion(t, x, y, w, h)
	}

	return regions
}

// Delete will dispose of the texture.
func (t *Texture) Delete() {
	GL.DeleteTexture(t.id)
}

// Bind will bind the texture.
func (t *Texture) Bind() {
	GL.BindTexture(GL.TEXTURE_2D, t.id)
}

// Unbind will unbind all textures.
func (t *Texture) Unbind() {
	GL.BindTexture(GL.TEXTURE_2D, nil)
}

// Width returns the width of the texture.
func (t *Texture) Width() int {
	return t.width
}

// Height returns the height of the texture.
func (t *Texture) Height() int {
	return t.height
}

// SetFilter sets the filter type used when scaling a texture up or
// down. The default is nearest which will not doing any interpolation
// between pixels.
func (t *Texture) SetFilter(min, max int) {
	t.minFilter = min
	t.maxFilter = max
	t.Bind()
	GL.TexParameteri(GL.TEXTURE_2D, GL.TEXTURE_MIN_FILTER, min)
	GL.TexParameteri(GL.TEXTURE_2D, GL.TEXTURE_MAG_FILTER, max)
}

// Returns the current min and max filters used.
func (t *Texture) Filter() (int, int) {
	return t.minFilter, t.maxFilter
}

func (t *Texture) SetWrap(u, v int) {
	t.uWrap = u
	t.vWrap = v
	t.Bind()
	GL.TexParameteri(GL.TEXTURE_2D, GL.TEXTURE_WRAP_S, u)
	GL.TexParameteri(GL.TEXTURE_2D, GL.TEXTURE_WRAP_T, v)
}

func (t *Texture) Wrap() (int, int) {
	return t.uWrap, t.vWrap
}

// A Batch allows geometry to be efficiently rendered by buffering
// render calls and sending them all at once.
type Batch struct {
	drawing      bool
	lastTexture  *Texture
	vertices     []float32
	vertexVBO    *BufferObject
	indices      []uint16
	indexVBO     *BufferObject
	index        int
	shader       *Shader
	inPosition   int
	inColor      int
	inTexCoords  int
	ufProjection *UniformObject
	projX        float32
	projY        float32
}

func NewBatch(width, height float32) *Batch {
	batch := new(Batch)

	batch.shader = NewShader(batchVert, batchFrag)
	batch.inPosition = batch.shader.GetAttrib("in_Position")
	batch.inColor = batch.shader.GetAttrib("in_Color")
	batch.inTexCoords = batch.shader.GetAttrib("in_TexCoords")
	batch.ufProjection = batch.shader.GetUniform("uf_Projection")

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

	batch.indexVBO = GL.CreateBuffer()
	batch.vertexVBO = GL.CreateBuffer()

	GL.BindBuffer(GL.ELEMENT_ARRAY_BUFFER, batch.indexVBO)
	GL.BufferData(GL.ELEMENT_ARRAY_BUFFER, batch.indices, GL.STATIC_DRAW)

	GL.BindBuffer(GL.ARRAY_BUFFER, batch.vertexVBO)
	GL.BufferData(GL.ARRAY_BUFFER, batch.vertices, GL.DYNAMIC_DRAW)

	GL.EnableVertexAttribArray(batch.inPosition)
	GL.EnableVertexAttribArray(batch.inTexCoords)
	GL.EnableVertexAttribArray(batch.inColor)

	GL.VertexAttribPointer(batch.inPosition, 2, GL.FLOAT, false, 20, 0)
	GL.VertexAttribPointer(batch.inTexCoords, 2, GL.FLOAT, false, 20, 8)
	GL.VertexAttribPointer(batch.inColor, 4, GL.UNSIGNED_BYTE, true, 20, 16)

	batch.projX = width / 2
	batch.projY = height / 2

	GL.Enable(GL.BLEND)
	GL.BlendFunc(GL.SRC_ALPHA, GL.ONE_MINUS_SRC_ALPHA)

	return batch
}

// Begin calculates the combined matrix and sets up rendering. This
// must be called before calling Draw.
func (b *Batch) Begin() {
	if b.drawing {
		log.Fatal("Batch.End() must be called first")
	}
	b.drawing = true
	b.shader.Bind()
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
	b.drawing = false

	b.lastTexture = nil
}

func (b *Batch) flush() {
	if b.lastTexture == nil {
		return
	}

	b.lastTexture.Bind()

	GL.Uniform2f(b.ufProjection, b.projX, b.projY)

	GL.BufferSubData(GL.ARRAY_BUFFER, 0, 20*4*b.index, b.vertices)
	GL.DrawElements(GL.TRIANGLES, 6*b.index, GL.UNSIGNED_SHORT, 0)

	b.index = 0
}

func (batch *Batch) Render(s *Sprite) {
	if s == nil || s.region == nil {
		return
	}

	if !batch.drawing {
		log.Fatal("Batch.Begin() must be called first")
	}

	r := s.region
	if r.texture != batch.lastTexture {
		if batch.lastTexture != nil {
			batch.flush()
		}
		batch.lastTexture = r.texture
	}

	color := s.color
	vertices := batch.vertices

	aX := s.Anchor.X
	aY := s.Anchor.Y

	w0 := r.width * (1 - aX)
	w1 := r.width * -aX

	h0 := r.height * (1 - aY)
	h1 := r.height * -aY

	transform := s.transform
	a := transform.A
	b := transform.C
	c := transform.B
	d := transform.D
	tx := transform.TX
	ty := transform.TY

	idx := batch.index * 20

	vertices[idx+0] = a*w1 + c*h1 + tx
	vertices[idx+1] = d*h1 + b*w1 + ty
	vertices[idx+2] = r.u
	vertices[idx+3] = r.v
	vertices[idx+4] = color

	vertices[idx+5] = a*w0 + c*h1 + tx
	vertices[idx+6] = d*h1 + b*w0 + ty
	vertices[idx+7] = r.u2
	vertices[idx+8] = r.v
	vertices[idx+9] = color

	vertices[idx+10] = a*w0 + c*h0 + tx
	vertices[idx+11] = d*h0 + b*w0 + ty
	vertices[idx+12] = r.u2
	vertices[idx+13] = r.v2
	vertices[idx+14] = color

	vertices[idx+15] = a*w1 + c*h0 + tx
	vertices[idx+16] = d*h0 + b*w1 + ty
	vertices[idx+17] = r.u
	vertices[idx+18] = r.v2
	vertices[idx+19] = color

	batch.index += 1

	if batch.index >= size {
		batch.flush()
	}
}

func (b *Batch) Draw(r *Region, x, y, originX, originY, scaleX, scaleY, rotation, color float32) {
	if !b.drawing {
		log.Fatal("Batch.Begin() must be called first")
	}

	if r.texture != b.lastTexture {
		if b.lastTexture != nil {
			b.flush()
		}
		b.lastTexture = r.texture
	}

	x -= originX * r.width
	y -= originY * r.height

	originX = r.width * originX
	originY = r.height * originY

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
	b.vertices[idx+4] = color

	b.vertices[idx+5] = x4
	b.vertices[idx+6] = y4
	b.vertices[idx+7] = r.u2
	b.vertices[idx+8] = r.v
	b.vertices[idx+9] = color

	b.vertices[idx+10] = x2
	b.vertices[idx+11] = y2
	b.vertices[idx+12] = r.u
	b.vertices[idx+13] = r.v2
	b.vertices[idx+14] = color

	b.vertices[idx+15] = x3
	b.vertices[idx+16] = y3
	b.vertices[idx+17] = r.u2
	b.vertices[idx+18] = r.v2
	b.vertices[idx+19] = color

	b.index += 1

	if b.index >= size {
		b.flush()
	}
}

func (b *Batch) SetProjection(width, height float32) {
	b.projX = width / 2
	b.projY = height / 2
}
