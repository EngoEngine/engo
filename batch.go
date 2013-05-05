package eng

import (
	gl "github.com/chsc/gogl/gl21"
	"unsafe"
)

const size = 1000

type Batch struct {
	drawing     bool
	lastTexture *Texture
	vertices    [size][2]gl.Float
	vertexVBO   gl.Uint
	colors      [size][4]gl.Float
	colorVBO    gl.Uint
	coords      [size][2]gl.Float
	coordVBO    gl.Uint
	index       gl.Sizei
	shader      *Shader
	combined    *Matrix
	projection  *Matrix
	transform   *Matrix
	r, g, b, a  gl.Float
}

func NewBatch() *Batch {
	batch := new(Batch)
	batch.shader = NewShader(vert, frag)

	batch.lastTexture = NewTexture("test.png")

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
	batch.projection = NewMatrix().SetOrtho(0, 0, float32(Width()), float32(Height()))
	batch.r = 1
	batch.g = 1
	batch.b = 1
	batch.a = 1

	return batch
}

func (b *Batch) Begin() {
	if b.drawing {
		panic("Batch.End() must be called first")
	}
	b.combined.Set(b.projection).Mul(b.transform)
	b.drawing = true
}

func (b *Batch) Draw(r *Region, x, y float32) {
	if r.texture != b.lastTexture {
		b.flush()
		b.lastTexture = r.texture
	}

	x1 := gl.Float(x)
	y1 := gl.Float(y)
	x2 := x1 + gl.Float(r.width)
	y2 := y1 + gl.Float(r.height)

	b.vertices[b.index+0][0] = x1
	b.vertices[b.index+0][1] = y1
	b.vertices[b.index+1][0] = x1
	b.vertices[b.index+1][1] = y2
	b.vertices[b.index+2][0] = x2
	b.vertices[b.index+2][1] = y2
	b.vertices[b.index+3][0] = x2
	b.vertices[b.index+3][1] = y1

	b.colors[b.index+0][0] = b.r
	b.colors[b.index+0][1] = b.g
	b.colors[b.index+0][2] = b.b
	b.colors[b.index+0][3] = b.a
	b.colors[b.index+1][0] = b.r
	b.colors[b.index+1][1] = b.g
	b.colors[b.index+1][2] = b.b
	b.colors[b.index+1][3] = b.a
	b.colors[b.index+2][0] = b.r
	b.colors[b.index+2][1] = b.g
	b.colors[b.index+2][2] = b.b
	b.colors[b.index+2][3] = b.a
	b.colors[b.index+3][0] = b.r
	b.colors[b.index+3][1] = b.g
	b.colors[b.index+3][2] = b.b
	b.colors[b.index+3][3] = b.a

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

func (bt *Batch) SetColor(r, g, b, a float32) {
	bt.r = gl.Float(r)
	bt.g = gl.Float(g)
	bt.b = gl.Float(b)
	bt.a = gl.Float(a)
}

func (b *Batch) Resize() {
	b.projection.SetOrtho(0, 0, float32(Width()), float32(Height()))
}

func (b *Batch) flush() {
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	gl.Enable(gl.TEXTURE_2D)
	gl.ActiveTexture(gl.TEXTURE0)
	b.lastTexture.Bind()

	b.shader.Bind()

	gl.UniformMatrix4fv(b.shader.UfMatrix, 1, gl.FALSE, &b.combined[0])
	gl.UniformMatrix4fv(b.shader.UfMatrix, 1, gl.FALSE, &b.combined[0])

	gl.BindBuffer(gl.ARRAY_BUFFER, b.vertexVBO)
	gl.BufferSubData(gl.ARRAY_BUFFER, gl.Intptr(0), gl.Sizeiptr(int(unsafe.Sizeof([2]gl.Float{}))*int(b.index)), gl.Pointer(&b.vertices[0]))
	gl.EnableVertexAttribArray(b.shader.InPosition)
	gl.VertexAttribPointer(b.shader.InPosition, 2, gl.FLOAT, gl.FALSE, 0, nil)

	gl.BindBuffer(gl.ARRAY_BUFFER, b.colorVBO)
	gl.BufferSubData(gl.ARRAY_BUFFER, gl.Intptr(0), gl.Sizeiptr(int(unsafe.Sizeof([4]gl.Float{}))*int(b.index)), gl.Pointer(&b.colors[0]))
	gl.EnableVertexAttribArray(b.shader.InColor)
	gl.VertexAttribPointer(b.shader.InColor, 4, gl.FLOAT, gl.FALSE, 0, nil)

	gl.BindBuffer(gl.ARRAY_BUFFER, b.coordVBO)
	gl.BufferSubData(gl.ARRAY_BUFFER, gl.Intptr(0), gl.Sizeiptr(int(unsafe.Sizeof([2]gl.Float{}))*int(b.index)), gl.Pointer(&b.coords[0]))
	gl.EnableVertexAttribArray(b.shader.InTexCoords)
	gl.VertexAttribPointer(b.shader.InTexCoords, 2, gl.FLOAT, gl.FALSE, 0, nil)

	gl.DrawArrays(gl.QUADS, 0, b.index)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	b.index = 0
}
