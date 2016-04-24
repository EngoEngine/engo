// +build !netgo

package engo

type Shader interface {
	Initialize(width, height float32)
	Pre()
	Draw(texture *gl.Texture, buffer *gl.Buffer, x, y, rotation float32)
	Post()
}

type defaultShader struct {
	indices  []uint16
	indexVBO *gl.Buffer
	program  *gl.Program

	projX float32
	projY float32

	lastTexture *gl.Texture

	inPosition   int
	inTexCoords  int
	inColor      int
	ufCamera     *gl.UniformLocation
	ufPosition   *gl.UniformLocation
	ufProjection *gl.UniformLocation
}

type hudShader struct {
	indices  []uint16
	indexVBO *gl.Buffer
	program  *gl.Program

	projX float32
	projY float32

	lastTexture *gl.Texture

	inPosition   int
	inTexCoords  int
	inColor      int
	ufPosition   *gl.UniformLocation
	ufProjection *gl.UniformLocation
}

func (s *defaultShader) Draw(texture *gl.Texture, buffer *gl.Buffer, x, y, rotation float32) {
	if s.lastTexture != texture {
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

func (s *hudShader) Draw(texture *gl.Texture, buffer *gl.Buffer, x, y, rotation float32) {
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
