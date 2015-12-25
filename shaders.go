package engi

import (
	"github.com/paked/webgl"
)

const bufferSize = 10000

type Shader interface {
	Initialize(width, height float32)
	Pre()
	Draw(texture *webgl.Texture, buffer *webgl.Buffer, x, y, rotation float32)
	Post()
}

type DefaultShader struct {
	indices  []uint16
	indexVBO *webgl.Buffer
	program  *webgl.Program

	projX float32
	projY float32

	lastTexture *webgl.Texture

	inPosition   int
	inTexCoords  int
	inColor      int
	ufCamera     *webgl.UniformLocation
	ufPosition   *webgl.UniformLocation
	ufProjection *webgl.UniformLocation
}

func (s *DefaultShader) Initialize(width, height float32) {
	s.program = LoadShader(`
#version 120

attribute vec2 in_Position;
attribute vec2 in_TexCoords;
attribute vec4 in_Color;

uniform vec2 uf_Position;
uniform vec3 uf_Camera;
uniform vec2 uf_Projection;

varying vec4 var_Color;
varying vec2 var_TexCoords;

void main() {
  var_Color = in_Color;
  var_TexCoords = in_TexCoords;

  gl_Position = vec4((in_Position.x + uf_Position.x - uf_Camera.x)/  uf_Projection.x,
  					 (in_Position.y + uf_Position.y - uf_Camera.y)/ -uf_Projection.y,
  					 0.0, uf_Camera.z);

}`, `
/* Fragment Shader */
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
}`)

	// Create and populate indices buffer
	s.indices = make([]uint16, 6*bufferSize)
	for i, j := 0, 0; i < bufferSize*6; i, j = i+6, j+4 {
		s.indices[i+0] = uint16(j + 0)
		s.indices[i+1] = uint16(j + 1)
		s.indices[i+2] = uint16(j + 2)
		s.indices[i+3] = uint16(j + 0)
		s.indices[i+4] = uint16(j + 2)
		s.indices[i+5] = uint16(j + 3)
	}
	s.indexVBO = Gl.CreateBuffer()
	Gl.BindBuffer(Gl.ELEMENT_ARRAY_BUFFER, s.indexVBO)
	Gl.BufferData(Gl.ELEMENT_ARRAY_BUFFER, s.indices, Gl.STATIC_DRAW)

	s.SetProjection(width, height)

	// Define things that should be read from the texture buffer
	s.inPosition = Gl.GetAttribLocation(s.program, "in_Position")
	s.inTexCoords = Gl.GetAttribLocation(s.program, "in_TexCoords")
	s.inColor = Gl.GetAttribLocation(s.program, "in_Color")

	// Define things that should be set per draw
	s.ufCamera = Gl.GetUniformLocation(s.program, "uf_Camera")
	s.ufPosition = Gl.GetUniformLocation(s.program, "uf_Position")
	s.ufProjection = Gl.GetUniformLocation(s.program, "uf_Projection")

	// Enable those things
	Gl.EnableVertexAttribArray(s.inPosition)
	Gl.EnableVertexAttribArray(s.inTexCoords)
	Gl.EnableVertexAttribArray(s.inColor)

	Gl.Enable(Gl.BLEND)
	Gl.BlendFunc(Gl.SRC_ALPHA, Gl.ONE_MINUS_SRC_ALPHA)
}

func (s *DefaultShader) Pre() {
	Gl.UseProgram(s.program)
	Gl.Uniform2f(s.ufProjection, s.projX, s.projY)
	Gl.Uniform3f(s.ufCamera, cam.x, cam.y, cam.z)
}

func (s *DefaultShader) Draw(texture *webgl.Texture, buffer *webgl.Buffer, x, y, rotation float32) {
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

func (s *DefaultShader) Post() {
	s.lastTexture = nil
}

func (s *DefaultShader) SetProjection(width, height float32) {
	s.projX = width / 2
	s.projY = height / 2
}

type HUDShader struct {
	indices  []uint16
	indexVBO *webgl.Buffer
	program  *webgl.Program

	projX float32
	projY float32

	lastTexture        *webgl.Texture
	drawCount          int
	coordBuffer        *webgl.Buffer
	coordBufferContent []float32
	coordIndices       []uint16
	coordIndicesVBO    *webgl.Buffer

	inPosition    int
	inTexCoords   int
	inColor       int
	inCoordinates int
	ufPosition    *webgl.UniformLocation
	ufProjection  *webgl.UniformLocation
}

func (s *HUDShader) Initialize(width, height float32) {
	s.program = LoadShader(`
#version 120

attribute vec2 in_Position;
attribute vec2 in_TexCoords;
attribute vec4 in_Color;
attribute vec2 in_Coordinates;

uniform vec2 uf_Projection;

varying vec4 var_Color;
varying vec2 var_TexCoords;

void main() {
  var_Color = in_Color;
  var_TexCoords = in_TexCoords;

  gl_Position = vec4((in_Position.x + in_Coordinates.x)/  uf_Projection.x - 1.0,
  					 (in_Position.y + in_Coordinates.y)/ -uf_Projection.y + 1.0,
  					 0.0, 1.0);

}`, `
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
}`)

	// Create and populate indices buffer
	s.coordIndices = make([]uint16, 6*bufferSize)
	for i, j := 0, 0; i < bufferSize*6; i, j = i+6, j+4 {
		s.coordIndices[i+0] = uint16(j + 0)
		s.coordIndices[i+1] = uint16(j + 0)
		s.coordIndices[i+2] = uint16(j + 0)
		s.coordIndices[i+3] = uint16(j + 0)
		s.coordIndices[i+4] = uint16(j + 0)
		s.coordIndices[i+5] = uint16(j + 0)
	}
	s.coordIndicesVBO = Gl.CreateBuffer()
	Gl.BindBuffer(Gl.ELEMENT_ARRAY_BUFFER, s.coordIndicesVBO)
	Gl.BufferData(Gl.ELEMENT_ARRAY_BUFFER, s.coordIndices, Gl.STATIC_DRAW)

	s.indices = make([]uint16, 6*bufferSize)
	for i, j := 0, 0; i < bufferSize*6; i, j = i+6, j+4 {
		s.indices[i+0] = uint16(j + 0)
		s.indices[i+1] = uint16(j + 1)
		s.indices[i+2] = uint16(j + 2)
		s.indices[i+3] = uint16(j + 0)
		s.indices[i+4] = uint16(j + 2)
		s.indices[i+5] = uint16(j + 3)
	}
	s.indexVBO = Gl.CreateBuffer()
	Gl.BindBuffer(Gl.ELEMENT_ARRAY_BUFFER, s.indexVBO)
	Gl.BufferData(Gl.ELEMENT_ARRAY_BUFFER, s.indices, Gl.STATIC_DRAW)

	s.SetProjection(width, height)

	// Define things that should be read from the texture buffer
	s.inPosition = Gl.GetAttribLocation(s.program, "in_Position")
	s.inTexCoords = Gl.GetAttribLocation(s.program, "in_TexCoords")
	s.inColor = Gl.GetAttribLocation(s.program, "in_Color")
	s.inCoordinates = Gl.GetAttribLocation(s.program, "in_Coordinates")

	// Define things that should be set per draw
	s.ufPosition = Gl.GetUniformLocation(s.program, "uf_Position")
	s.ufProjection = Gl.GetUniformLocation(s.program, "uf_Projection")

	// Enable those things
	Gl.EnableVertexAttribArray(s.inPosition)
	Gl.EnableVertexAttribArray(s.inTexCoords)
	Gl.EnableVertexAttribArray(s.inColor)
	Gl.EnableVertexAttribArray(s.inCoordinates)

	Gl.Enable(Gl.BLEND)
	Gl.BlendFunc(Gl.SRC_ALPHA, Gl.ONE_MINUS_SRC_ALPHA)

	// TODO: listen for Projection changes
}

func (s *HUDShader) Pre() {
	Gl.UseProgram(s.program)
	Gl.BindBuffer(Gl.ELEMENT_ARRAY_BUFFER, s.indexVBO)

	Gl.Uniform2f(s.ufProjection, s.projX, s.projY)

	s.coordBufferContent = make([]float32, 0)
}

func (s *HUDShader) Draw(texture *webgl.Texture, buffer *webgl.Buffer, x, y, rotation float32) {
	if s.lastTexture != texture {
		Gl.BindTexture(Gl.TEXTURE_2D, texture)
		Gl.BindBuffer(Gl.ARRAY_BUFFER, buffer)

		Gl.VertexAttribPointer(s.inPosition, 2, Gl.FLOAT, false, 20, 0)
		Gl.VertexAttribPointer(s.inTexCoords, 2, Gl.FLOAT, false, 20, 8)
		Gl.VertexAttribPointer(s.inColor, 4, Gl.UNSIGNED_BYTE, true, 20, 16)

		s.lastTexture = texture
	}

	s.coordBufferContent = append(s.coordBufferContent, []float32{x, y, rotation}...)
	s.drawCount++

	if s.drawCount > bufferSize {
		s.flush()
	}
}

func (s *HUDShader) flush() {
	Gl.BindBuffer(Gl.ARRAY_BUFFER, s.coordBuffer)
	Gl.BufferData(Gl.ARRAY_BUFFER, s.coordBufferContent, Gl.DYNAMIC_DRAW)
	Gl.VertexAttribPointer(s.inCoordinates, 3, Gl.FLOAT, false, 12, 0)

	Gl.DrawElements(Gl.TRIANGLES, 6*s.drawCount, Gl.UNSIGNED_SHORT, 0)
	s.drawCount = 0
}

func (s *HUDShader) Post() {
	s.flush()
	s.lastTexture = nil
}

func (s *HUDShader) SetProjection(width, height float32) {
	s.projX = width / 2
	s.projY = height / 2
}

// ShadersLibrary is the manager for the Shaders
type ShadersLibrary struct {
	setup bool

	def     DefaultShader
	shaders []Shader
}

var Shaders = ShadersLibrary{
	shaders: make([]Shader, HighestGround+1),
}

// Registers the `Shader` for the given `PriorityLevel`; possibly overwriting previously registered Shaders
// It does no initialization whatsoever
func (s *ShadersLibrary) Register(prio PriorityLevel, sh Shader) {
	s.shaders[prio] = sh
}

// Get returns the `Shader` that should be used for the given `PriorityLevel`
func (s *ShadersLibrary) Get(prio PriorityLevel) Shader {
	if sh := s.shaders[prio]; sh != nil {
		return sh
	}
	return &s.def
}
