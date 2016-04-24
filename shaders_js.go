//+build netgo

package engo

import (
	"github.com/gopherjs/gopherjs/js"
)

type Shader interface {
	Initialize(width, height float32)
	Pre()
	Draw(texture *js.Object, buffer *js.Object, x, y, rotation float32)
	Post()
}

type defaultShader struct {
	indices  []uint16
	indexVBO *js.Object
	program  *js.Object

	projX float32
	projY float32

	lastTexture *js.Object

	inPosition   int
	inTexCoords  int
	inColor      int
	ufCamera     *js.Object
	ufPosition   *js.Object
	ufProjection *js.Object
}

func (s *defaultShader) Draw(texture *js.Object, buffer *js.Object, x, y, rotation float32) {
	if s.lastTexture != texture {
		Gl.ActiveTexture(Gl.TEXTURE0)
		Gl.BindTexture(Gl.TEXTURE_2D, texture)
		Gl.BindBuffer(Gl.ARRAY_BUFFER, buffer)

		Gl.VertexAttribPointer(s.inPosition, 2, Gl.FLOAT, false, 20, 0)
		Gl.VertexAttribPointer(s.inTexCoords, 2, Gl.FLOAT, false, 20, 8)
		Gl.VertexAttribPointer(s.inColor, 4, Gl.UNSIGNED_BYTE, true, 20, 16)

		s.lastTexture = texture
	}

	// TODO: add rotation
	Gl.Uniform2f(s.ufPosition, x, y)
	Gl.DrawElements(Gl.TRIANGLES, 6, Gl.UNSIGNED_SHORT, 0)
}

type hudShader struct {
	indices  []uint16
	indexVBO *js.Object
	program  *js.Object

	projX float32
	projY float32

	lastTexture *js.Object

	inPosition   int
	inTexCoords  int
	inColor      int
	ufPosition   *js.Object
	ufProjection *js.Object
}

func (s *hudShader) Draw(texture *js.Object, buffer *js.Object, x, y, rotation float32) {
	if s.lastTexture != texture {
		Gl.BindTexture(Gl.TEXTURE_2D, texture)
		Gl.BindBuffer(Gl.ARRAY_BUFFER, buffer)

		Gl.VertexAttribPointer(s.inPosition, 2, Gl.FLOAT, false, 20, 0)
		Gl.VertexAttribPointer(s.inTexCoords, 2, Gl.FLOAT, false, 20, 8)
		Gl.VertexAttribPointer(s.inColor, 4, Gl.UNSIGNED_BYTE, true, 20, 16)

		s.lastTexture = texture
	}

	Gl.Uniform2f(s.ufPosition, x, y)
	Gl.DrawElements(Gl.TRIANGLES, 6, Gl.UNSIGNED_SHORT, 0)
}
