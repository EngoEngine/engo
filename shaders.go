package engo

import (
	"engo.io/gl"
	"github.com/luxengine/math"
)

const bufferSize = 10000

type Shader interface {
	Initialize()
	Pre()
	Draw(texture *gl.Texture, buffer *gl.Buffer, x, y, scaleX, scaleY, rotation float32)
	Post()
}

type basicShader struct {
	indices  []uint16
	indexVBO *gl.Buffer
	program  *gl.Program

	lastTexture *gl.Texture

	inPosition  int
	inTexCoords int
	inColor     int

	matrixProjection *gl.UniformLocation
	matrixView       *gl.UniformLocation
	matrixModel      *gl.UniformLocation

	projectionMatrix []float32
	viewMatrix       []float32
	modelMatrix      []float32

	cameraEnabled bool
}

func (s *basicShader) Initialize() {
	s.program = LoadShader(`
#version 120

attribute vec2 in_Position;
attribute vec2 in_TexCoords;
attribute vec4 in_Color;

uniform mat3 matrixProjection;
uniform mat3 matrixView;
uniform mat3 matrixModel;

varying vec4 var_Color;
varying vec2 var_TexCoords;

void main() {
  var_Color = in_Color;
  var_TexCoords = in_TexCoords;

  vec3 matr = matrixProjection * matrixView * matrixModel * vec3(in_Position, 1.0);
  gl_Position = vec4(matr.xy, 0, matr.z);
}
`, `
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

	// Define things that should be read from the texture buffer
	s.inPosition = Gl.GetAttribLocation(s.program, "in_Position")
	s.inTexCoords = Gl.GetAttribLocation(s.program, "in_TexCoords")
	s.inColor = Gl.GetAttribLocation(s.program, "in_Color")

	// Define things that should be set per draw
	s.matrixProjection = Gl.GetUniformLocation(s.program, "matrixProjection")
	s.matrixView = Gl.GetUniformLocation(s.program, "matrixView")
	s.matrixModel = Gl.GetUniformLocation(s.program, "matrixModel")

	// Enable those things
	Gl.EnableVertexAttribArray(s.inPosition)
	Gl.EnableVertexAttribArray(s.inTexCoords)
	Gl.EnableVertexAttribArray(s.inColor)

	Gl.Enable(Gl.BLEND)
	Gl.BlendFunc(Gl.SRC_ALPHA, Gl.ONE_MINUS_SRC_ALPHA)

	s.projectionMatrix = make([]float32, 9)
	s.projectionMatrix[8] = 1

	s.viewMatrix = make([]float32, 9)
	s.viewMatrix[0] = 1
	s.viewMatrix[4] = 1
	s.viewMatrix[8] = 1

	s.modelMatrix = make([]float32, 9)
	s.modelMatrix[0] = 1
	s.modelMatrix[4] = 1
	s.modelMatrix[8] = 1
}

func (s *basicShader) Pre() {
	Gl.UseProgram(s.program)

	if scaleOnResize {
		s.projectionMatrix[0] = 1 / (gameWidth / 2)
		s.projectionMatrix[4] = 1 / (-gameHeight / 2)
	} else {
		s.projectionMatrix[0] = 1 / (windowWidth / 2)
		s.projectionMatrix[4] = 1 / (-windowHeight / 2)
	}

	if s.cameraEnabled {
		angle := cam.angle * math.Pi / 180

		s.viewMatrix[0] = math.Cos(angle)
		s.viewMatrix[1] = math.Sin(angle)
		s.viewMatrix[3] = -s.viewMatrix[1]
		s.viewMatrix[4] = s.viewMatrix[0]
		s.viewMatrix[6] = -cam.x
		s.viewMatrix[7] = -cam.y
		s.viewMatrix[8] = cam.z
	} else {
		s.viewMatrix[6] = -1 / s.projectionMatrix[0]
		s.viewMatrix[7] = 1 / s.projectionMatrix[4]
	}

	Gl.UniformMatrix3fv(s.matrixProjection, false, s.projectionMatrix)
	Gl.UniformMatrix3fv(s.matrixView, false, s.viewMatrix)
}

func (s *basicShader) Draw(texture *gl.Texture, buffer *gl.Buffer, x, y, scaleX, scaleY, rotation float32) {
	if s.lastTexture != texture {
		Gl.BindTexture(Gl.TEXTURE_2D, texture)
		Gl.BindBuffer(Gl.ARRAY_BUFFER, buffer)

		Gl.VertexAttribPointer(s.inPosition, 2, Gl.FLOAT, false, 20, 0)
		Gl.VertexAttribPointer(s.inTexCoords, 2, Gl.FLOAT, false, 20, 8)
		Gl.VertexAttribPointer(s.inColor, 4, Gl.UNSIGNED_BYTE, true, 20, 16)

		s.lastTexture = texture
	}

	angle := rotation * math.Pi / 180
	cos := math.Cos(angle)
	sin := math.Sin(angle)

	s.modelMatrix[0] = scaleX * cos
	s.modelMatrix[1] = scaleX * sin
	s.modelMatrix[3] = scaleY * -sin
	s.modelMatrix[4] = scaleY * cos
	s.modelMatrix[6] = x
	s.modelMatrix[7] = y

	Gl.UniformMatrix3fv(s.matrixModel, false, s.modelMatrix)

	Gl.DrawElements(Gl.TRIANGLES, 6, Gl.UNSIGNED_SHORT, 0)
}

func (s *basicShader) Post() {
	s.lastTexture = nil
}

var (
	DefaultShader = &basicShader{cameraEnabled: true}
	HUDShader     = &basicShader{cameraEnabled: false}
	shadersSet    bool
)

func initShaders() {
	if !shadersSet {
		DefaultShader.Initialize()
		HUDShader.Initialize()

		shadersSet = true
	}
}
