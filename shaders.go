package engo

import (
	"fmt"
)

const bufferSize = 10000
const vertSrc = `
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

}
`

var fragSrc = `
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
}
`

func (s *defaultShader) Initialize(width, height float32) {
	s.program = LoadShader(vertSrc, fragSrc)

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

	// We sometimes have to change our projection matrix
	Mailbox.Listen("WindowResizeMessage", func(m Message) {
		wrm, ok := m.(WindowResizeMessage)
		if !ok {
			return
		}

		if !scaleOnResize {
			s.SetProjection(float32(wrm.NewWidth), float32(wrm.NewHeight))
		}
	})
}

func (s *defaultShader) Pre() {
	Gl.UseProgram(s.program)
	Gl.Uniform2f(s.ufProjection, s.projX, s.projY)
	Gl.Uniform3f(s.ufCamera, cam.x, cam.y, cam.z)
}

func (s *defaultShader) Post() {
	s.lastTexture = nil
}

func (s *defaultShader) SetProjection(width, height float32) {
	s.projX = width / 2
	s.projY = height / 2
}

func (s *hudShader) Initialize(width, height float32) {
	s.program = LoadShader(vertSrc, fragSrc)

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
	s.ufPosition = Gl.GetUniformLocation(s.program, "uf_Position")
	s.ufProjection = Gl.GetUniformLocation(s.program, "uf_Projection")

	// Enable those things
	Gl.EnableVertexAttribArray(s.inPosition)
	Gl.EnableVertexAttribArray(s.inTexCoords)
	Gl.EnableVertexAttribArray(s.inColor)

	Gl.Enable(Gl.BLEND)
	Gl.BlendFunc(Gl.SRC_ALPHA, Gl.ONE_MINUS_SRC_ALPHA)

	// TODO: listen for Projection changes
}

func (s *hudShader) Pre() {
	Gl.UseProgram(s.program)
	Gl.Uniform2f(s.ufProjection, s.projX, s.projY)
}

func (s *hudShader) Post() {
	s.lastTexture = nil
}

func (s *hudShader) SetProjection(width, height float32) {
	s.projX = width / 2
	s.projY = height / 2
}

var (
	DefaultShader = &defaultShader{}
	HUDShader     = &hudShader{}
	shadersSet    bool
)

func initShaders(width, height float32) {
	if !shadersSet {
		fmt.Println("Initialized shaders", width, height)
		DefaultShader.Initialize(width, height)
		HUDShader.Initialize(width, height)

		shadersSet = true
	}
}
