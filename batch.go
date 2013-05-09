package eng

import (
	gl "github.com/chsc/gogl/gl21"
	"math"
	"unsafe"
)

const size = 1000
const degToRad = math.Pi / 180

type Batch struct {
	drawing      bool
	lastTexture  *Texture
	vertices     [size][2]gl.Float
	vertexVBO    gl.Uint
	colors       [size][4]gl.Float
	colorVBO     gl.Uint
	coords       [size][2]gl.Float
	coordVBO     gl.Uint
	index        gl.Sizei
	shader       *Shader
	customShader *Shader
	combined     *Matrix
	projection   *Matrix
	transform    *Matrix
	color        *Color
}

func NewBatch() *Batch {
	batch := new(Batch)
	batch.shader = NewShader(vert, frag)

	gl.GenBuffers(1, &batch.vertexVBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, batch.vertexVBO)
	gl.BufferData(gl.ARRAY_BUFFER, gl.Sizeiptr(int(unsafe.Sizeof([2]gl.Float{}))*size), gl.Pointer(&batch.vertices[0]), gl.DYNAMIC_DRAW)

	gl.GenBuffers(1, &batch.colorVBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, batch.colorVBO)
	gl.BufferData(gl.ARRAY_BUFFER, gl.Sizeiptr(int(unsafe.Sizeof([4]gl.Float{}))*size), gl.Pointer(&batch.colors[0]), gl.DYNAMIC_DRAW)

	gl.GenBuffers(1, &batch.coordVBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, batch.coordVBO)
	gl.BufferData(gl.ARRAY_BUFFER, gl.Sizeiptr(int(unsafe.Sizeof([2]gl.Float{}))*size), gl.Pointer(&batch.coords[0]), gl.DYNAMIC_DRAW)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	batch.combined = NewMatrix()
	batch.transform = NewMatrix()
	batch.projection = NewMatrix().SetToOrtho(0, float32(Width()), float32(Height()), 0, 0, 1)
	batch.color = NewColor(1, 1, 1, 1)

	return batch
}

func (b *Batch) Begin() {
	if b.drawing {
		panic("Batch.End() must be called first")
	}
	b.combined.Set(b.projection).Mul(b.transform)
	b.drawing = true
}

func (b *Batch) Draw(r *Region, x, y, originX, originY, scaleX, scaleY, rotation float32, color *Color) {
	if !b.drawing {
		panic("Batch.Begin() must be called first")
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

	b.vertices[b.index+0][0] = gl.Float(x1)
	b.vertices[b.index+0][1] = gl.Float(y1)
	b.vertices[b.index+1][0] = gl.Float(x2)
	b.vertices[b.index+1][1] = gl.Float(y2)
	b.vertices[b.index+2][0] = gl.Float(x3)
	b.vertices[b.index+2][1] = gl.Float(y3)
	b.vertices[b.index+3][0] = gl.Float(x4)
	b.vertices[b.index+3][1] = gl.Float(y4)

	b.colors[b.index+0][0] = gl.Float(c.R)
	b.colors[b.index+0][1] = gl.Float(c.G)
	b.colors[b.index+0][2] = gl.Float(c.B)
	b.colors[b.index+0][3] = gl.Float(c.A)
	b.colors[b.index+1][0] = gl.Float(c.R)
	b.colors[b.index+1][1] = gl.Float(c.G)
	b.colors[b.index+1][2] = gl.Float(c.B)
	b.colors[b.index+1][3] = gl.Float(c.A)
	b.colors[b.index+2][0] = gl.Float(c.R)
	b.colors[b.index+2][1] = gl.Float(c.G)
	b.colors[b.index+2][2] = gl.Float(c.B)
	b.colors[b.index+2][3] = gl.Float(c.A)
	b.colors[b.index+3][0] = gl.Float(c.R)
	b.colors[b.index+3][1] = gl.Float(c.G)
	b.colors[b.index+3][2] = gl.Float(c.B)
	b.colors[b.index+3][3] = gl.Float(c.A)

	b.coords[b.index+0][0] = r.u
	b.coords[b.index+0][1] = r.v
	b.coords[b.index+1][0] = r.u
	b.coords[b.index+1][1] = r.v2
	b.coords[b.index+2][0] = r.u2
	b.coords[b.index+2][1] = r.v2
	b.coords[b.index+3][0] = r.u2
	b.coords[b.index+3][1] = r.v

	b.index += 4

	if b.index >= size {
		b.flush()
	}
}

func (b *Batch) End() {
	if !b.drawing {
		panic("Batch.Begin() must be called first")
	}
	if b.index > 0 {
		b.flush()
	}
	b.drawing = false
}

func (b *Batch) SetColor(color *Color) {
	b.color.R = color.R
	b.color.G = color.G
	b.color.B = color.B
	b.color.A = color.A
}

func (b *Batch) SetShader(shader *Shader) {
	b.customShader = shader
}

func (b *Batch) SetProjection(m *Matrix) {
	b.projection.Set(m)
}

func (b *Batch) flush() {
	if b.lastTexture == nil {
		return
	}
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	gl.Enable(gl.TEXTURE_2D)
	gl.ActiveTexture(gl.TEXTURE0)
	b.lastTexture.Bind()

	shader := b.shader
	if b.customShader != nil {
		shader = b.customShader
	}

	shader.Bind()

	gl.UniformMatrix4fv(shader.UfMatrix, 1, gl.FALSE, &b.combined[0])

	gl.BindBuffer(gl.ARRAY_BUFFER, b.vertexVBO)
	gl.BufferSubData(gl.ARRAY_BUFFER, gl.Intptr(0), gl.Sizeiptr(int(unsafe.Sizeof([2]gl.Float{}))*int(b.index)), gl.Pointer(&b.vertices[0]))
	gl.EnableVertexAttribArray(shader.InPosition)
	gl.VertexAttribPointer(shader.InPosition, 2, gl.FLOAT, gl.FALSE, 0, nil)

	gl.BindBuffer(gl.ARRAY_BUFFER, b.colorVBO)
	gl.BufferSubData(gl.ARRAY_BUFFER, gl.Intptr(0), gl.Sizeiptr(int(unsafe.Sizeof([4]gl.Float{}))*int(b.index)), gl.Pointer(&b.colors[0]))
	gl.EnableVertexAttribArray(shader.InColor)
	gl.VertexAttribPointer(shader.InColor, 4, gl.FLOAT, gl.FALSE, 0, nil)

	gl.BindBuffer(gl.ARRAY_BUFFER, b.coordVBO)
	gl.BufferSubData(gl.ARRAY_BUFFER, gl.Intptr(0), gl.Sizeiptr(int(unsafe.Sizeof([2]gl.Float{}))*int(b.index)), gl.Pointer(&b.coords[0]))
	gl.EnableVertexAttribArray(shader.InTexCoords)
	gl.VertexAttribPointer(shader.InTexCoords, 2, gl.FLOAT, gl.FALSE, 0, nil)

	gl.DrawArrays(gl.QUADS, 0, b.index)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.UseProgram(0)

	b.index = 0
}
