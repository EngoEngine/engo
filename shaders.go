package engo

import (
	"engo.io/gl"
	"github.com/luxengine/math"
	"log"
)

const bufferSize = 10000

type Shader interface {
	Initialize()
	Pre()
	Draw(*RenderComponent, *SpaceComponent)
	Post()
}

type basicShader struct {
	indices  []uint16
	indexVBO *gl.Buffer
	program  *gl.Program

	lastTexture   *gl.Texture
	lastBuffer    *gl.Buffer
	lastRepeating TextureRepeating

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
	// Enable shader and buffer, enable attributes in shader
	Gl.UseProgram(s.program)
	Gl.BindBuffer(Gl.ELEMENT_ARRAY_BUFFER, s.indexVBO)
	Gl.EnableVertexAttribArray(s.inPosition)
	Gl.EnableVertexAttribArray(s.inTexCoords)
	Gl.EnableVertexAttribArray(s.inColor)

	if scaleOnResize {
		s.projectionMatrix[0] = 1 / (gameWidth / 2)
		s.projectionMatrix[4] = 1 / (-gameHeight / 2)
	} else {
		s.projectionMatrix[0] = 1 / (windowWidth / 2)
		s.projectionMatrix[4] = 1 / (-windowHeight / 2)
	}

	if s.cameraEnabled {
		s.viewMatrix[1], s.viewMatrix[0] = math.Sincos(cam.angle * math.Pi / 180)
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

func (s *basicShader) Draw(ren *RenderComponent, space *SpaceComponent) {
	if s.lastBuffer != ren.buffer || ren.buffer == nil {
		s.updateBuffer(ren)

		Gl.BindBuffer(Gl.ARRAY_BUFFER, ren.buffer)
		Gl.VertexAttribPointer(s.inPosition, 2, Gl.FLOAT, false, 20, 0)
		Gl.VertexAttribPointer(s.inTexCoords, 2, Gl.FLOAT, false, 20, 8)
		Gl.VertexAttribPointer(s.inColor, 4, Gl.UNSIGNED_BYTE, true, 20, 16)

		s.lastBuffer = ren.buffer
	}

	if s.lastTexture != ren.Drawable.Texture() {
		Gl.BindTexture(Gl.TEXTURE_2D, ren.Drawable.Texture())

		s.lastTexture = ren.Drawable.Texture()
	}

	if s.lastRepeating != ren.Repeat {
		var val int
		switch ren.Repeat {
		case CLAMP_TO_EDGE:
			val = Gl.CLAMP_TO_EDGE
		case CLAMP_TO_BORDER:
			val = Gl.CLAMP_TO_EDGE
		case REPEAT:
			val = Gl.REPEAT
		case MIRRORED_REPEAT:
			val = Gl.MIRRORED_REPEAT
		}

		Gl.TexParameteri(Gl.TEXTURE_2D, Gl.TEXTURE_WRAP_S, val)
		Gl.TexParameteri(Gl.TEXTURE_2D, Gl.TEXTURE_WRAP_T, val)
	}

	if space.Rotation != 0 {
		sin, cos := math.Sincos(space.Rotation * math.Pi / 180)

		s.modelMatrix[0] = ren.Scale.X * cos
		s.modelMatrix[1] = ren.Scale.X * sin
		s.modelMatrix[3] = ren.Scale.Y * -sin
		s.modelMatrix[4] = ren.Scale.Y * cos
	} else {
		s.modelMatrix[0] = ren.Scale.X
		s.modelMatrix[1] = 0
		s.modelMatrix[3] = 0
		s.modelMatrix[4] = ren.Scale.Y
	}

	s.modelMatrix[6] = space.Position.X
	s.modelMatrix[7] = space.Position.Y

	Gl.UniformMatrix3fv(s.matrixModel, false, s.modelMatrix)

	Gl.DrawElements(Gl.TRIANGLES, 6, Gl.UNSIGNED_SHORT, 0)
}

func (s *basicShader) Post() {
	s.lastTexture = nil
	s.lastBuffer = nil

	// Cleanup
	Gl.BindTexture(Gl.TEXTURE_2D, nil)
	Gl.DisableVertexAttribArray(s.inPosition)
	Gl.DisableVertexAttribArray(s.inTexCoords)
	Gl.DisableVertexAttribArray(s.inColor)
}

func (s *basicShader) updateBuffer(ren *RenderComponent) {
	if len(ren.bufferContent) == 0 {
		ren.bufferContent = make([]float32, 20) // because we add 20 elements to it
	}

	if changed := s.generateBufferContent(ren, ren.bufferContent); !changed {
		return
	}

	if ren.buffer == nil {
		ren.buffer = Gl.CreateBuffer()
	}
	Gl.BindBuffer(Gl.ARRAY_BUFFER, ren.buffer)
	Gl.BufferData(Gl.ARRAY_BUFFER, ren.bufferContent, Gl.STATIC_DRAW)
}

func (s *basicShader) generateBufferContent(ren *RenderComponent, buffer []float32) bool {
	w := ren.Drawable.Width()
	h := ren.Drawable.Height()

	colorR, colorG, colorB, colorA := ren.Color.RGBA()
	colorR >>= 8
	colorG >>= 8
	colorB >>= 8
	colorA >>= 8

	red := colorR
	green := colorG << 8
	blue := colorB << 16
	alpha := colorA << 24

	tint := math.Float32frombits((alpha | blue | green | red) & 0xfeffffff)

	u, v, u2, v2 := ren.Drawable.View()

	var changed bool

	//setBufferValue(buffer, 0, 0, &changed)
	//setBufferValue(buffer, 1, 0, &changed)
	setBufferValue(buffer, 2, u, &changed)
	setBufferValue(buffer, 3, v, &changed)
	setBufferValue(buffer, 4, tint, &changed)

	setBufferValue(buffer, 5, w, &changed)
	//setBufferValue(buffer, 6, 0, &changed)
	setBufferValue(buffer, 7, u2, &changed)
	setBufferValue(buffer, 8, v, &changed)
	setBufferValue(buffer, 9, tint, &changed)

	setBufferValue(buffer, 10, w, &changed)
	setBufferValue(buffer, 11, h, &changed)
	setBufferValue(buffer, 12, u2, &changed)
	setBufferValue(buffer, 13, v2, &changed)
	setBufferValue(buffer, 14, tint, &changed)

	//setBufferValue(buffer, 15, 0, &changed)
	setBufferValue(buffer, 16, h, &changed)
	setBufferValue(buffer, 17, u, &changed)
	setBufferValue(buffer, 18, v2, &changed)
	setBufferValue(buffer, 19, tint, &changed)

	return changed
}

func setBufferValue(buffer []float32, index int, value float32, changed *bool) {
	if buffer[index] != value {
		buffer[index] = value
		*changed = true
	}
}

type legacyShader struct {
	program *gl.Program

	indicesTriangles     []uint16
	indicesRectangles    []uint16
	indicesTrianglesVBO  *gl.Buffer
	indicesRectanglesVBO *gl.Buffer

	inPosition int
	inColor    int

	matrixProjection *gl.UniformLocation
	matrixView       *gl.UniformLocation
	matrixModel      *gl.UniformLocation
	inRadius         *gl.UniformLocation
	inCenter         *gl.UniformLocation
	inViewport       *gl.UniformLocation

	projectionMatrix []float32
	viewMatrix       []float32
	modelMatrix      []float32

	cameraEnabled bool

	lastBuffer *gl.Buffer
}

func (l *legacyShader) Initialize() {
	l.program = LoadShader(`
attribute vec2 in_Position;
attribute vec4 in_Color;

uniform mat3 matrixProjection;
uniform mat3 matrixView;
uniform mat3 matrixModel;
uniform vec2 in_Radius;
uniform vec2 in_Center;
uniform vec2 in_Viewport;

varying vec4 var_Color;
varying vec2 var_Radius;
varying vec2 var_Center;

void main() {
  var_Color = in_Color;

  vec3 matr = matrixProjection * matrixView * matrixModel * vec3(in_Position, 1.0);
  gl_Position = vec4(matr.xy, 0, matr.z);

  if (in_Radius.x > 0.0 && in_Radius.y > 0.0)
  {
    var_Radius = in_Radius;

    vec3 vecCenter = (matrixProjection * matrixView * matrixModel * vec3(in_Center, 1.0));
    var_Center = (vecCenter.xy/vecCenter.z + vec2(1.0, 1.0)) * in_Viewport / vec2(2.0);
  } else {
    var_Radius = vec2(0.0, 0.0);
  }
}
`, `
#ifdef GL_ES
#define LOWP lowp
precision mediump float;
#else
#define LOWP
#endif

varying vec4 var_Color;
varying vec2 var_Radius;
varying vec2 var_Center;

void main (void) {
  gl_FragColor = var_Color;

  if (var_Radius.x > 0.0 && var_Radius.y > 0.0)
  {
    if (pow(gl_FragCoord.x - var_Center.x, 2.0) / pow(var_Radius.x, 2.0) + pow(gl_FragCoord.y - var_Center.y, 2.0) / pow(var_Radius.y, 2.0) > 1.0)
    {
      gl_FragColor.w = 0.0;
    }
  }
}`)

	// Create and populate indices buffer
	l.indicesTriangles = []uint16{0, 1, 2}
	l.indicesTrianglesVBO = Gl.CreateBuffer()
	Gl.BindBuffer(Gl.ELEMENT_ARRAY_BUFFER, l.indicesTrianglesVBO)
	Gl.BufferData(Gl.ELEMENT_ARRAY_BUFFER, l.indicesTriangles, Gl.STATIC_DRAW)

	l.indicesRectangles = []uint16{0, 1, 2, 0, 2, 3}
	l.indicesRectanglesVBO = Gl.CreateBuffer()
	Gl.BindBuffer(Gl.ELEMENT_ARRAY_BUFFER, l.indicesRectanglesVBO)
	Gl.BufferData(Gl.ELEMENT_ARRAY_BUFFER, l.indicesRectangles, Gl.STATIC_DRAW)

	// Define things that should be read from the texture buffer
	l.inPosition = Gl.GetAttribLocation(l.program, "in_Position")
	l.inColor = Gl.GetAttribLocation(l.program, "in_Color")

	// Define things that should be set per draw
	l.matrixProjection = Gl.GetUniformLocation(l.program, "matrixProjection")
	l.matrixView = Gl.GetUniformLocation(l.program, "matrixView")
	l.matrixModel = Gl.GetUniformLocation(l.program, "matrixModel")
	l.inRadius = Gl.GetUniformLocation(l.program, "in_Radius")
	l.inCenter = Gl.GetUniformLocation(l.program, "in_Center")
	l.inViewport = Gl.GetUniformLocation(l.program, "in_Viewport")

	// Not sure if they go here, or in Pre()
	Gl.Enable(Gl.BLEND)
	Gl.BlendFunc(Gl.SRC_ALPHA, Gl.ONE_MINUS_SRC_ALPHA)

	l.projectionMatrix = make([]float32, 9)
	l.projectionMatrix[8] = 1

	l.viewMatrix = make([]float32, 9)
	l.viewMatrix[0] = 1
	l.viewMatrix[4] = 1
	l.viewMatrix[8] = 1

	l.modelMatrix = make([]float32, 9)
	l.modelMatrix[0] = 1
	l.modelMatrix[4] = 1
	l.modelMatrix[8] = 1
}

func (l *legacyShader) Pre() {
	// Bind shader and buffer, enable attributes
	Gl.UseProgram(l.program)
	Gl.BindBuffer(Gl.ELEMENT_ARRAY_BUFFER, l.indicesTrianglesVBO)
	Gl.EnableVertexAttribArray(l.inPosition)
	Gl.EnableVertexAttribArray(l.inColor)

	if scaleOnResize {
		l.projectionMatrix[0] = 1 / (gameWidth / 2)
		l.projectionMatrix[4] = 1 / (-gameHeight / 2)
	} else {
		l.projectionMatrix[0] = 1 / (gameWidth / 2)   // TODO: canvasWidth
		l.projectionMatrix[4] = 1 / (-gameHeight / 2) // TODO: canvasHeight
	}

	if l.cameraEnabled {
		l.viewMatrix[1], l.viewMatrix[0] = math.Sincos(cam.angle * math.Pi / 180)
		l.viewMatrix[3] = -l.viewMatrix[1]
		l.viewMatrix[4] = l.viewMatrix[0]
		l.viewMatrix[6] = -cam.x
		l.viewMatrix[7] = -cam.y
		l.viewMatrix[8] = cam.z
	} else {
		l.viewMatrix[6] = -1 / l.projectionMatrix[0]
		l.viewMatrix[7] = 1 / l.projectionMatrix[4]
	}

	Gl.UniformMatrix3fv(l.matrixProjection, false, l.projectionMatrix)
	Gl.UniformMatrix3fv(l.matrixView, false, l.viewMatrix)
	Gl.Uniform2f(l.inViewport, gameWidth, gameHeight) // TODO: canvasWidth/Height
}

func (l *legacyShader) updateBuffer(ren *RenderComponent, space *SpaceComponent) {
	if len(ren.bufferContent) == 0 {
		ren.bufferContent = make([]float32, 12) // because we add at most this many elements to it
	}

	if changed := l.generateBufferContent(ren, space, ren.bufferContent); !changed {
		return
	}

	if ren.buffer == nil {
		ren.buffer = Gl.CreateBuffer()
	}
	Gl.BindBuffer(Gl.ARRAY_BUFFER, ren.buffer)
	Gl.BufferData(Gl.ARRAY_BUFFER, ren.bufferContent, Gl.STATIC_DRAW)
}

func (l *legacyShader) generateBufferContent(ren *RenderComponent, space *SpaceComponent, buffer []float32) bool {
	w := space.Width
	h := space.Height

	colorR, colorG, colorB, colorA := ren.Color.RGBA()
	colorR >>= 8
	colorG >>= 8
	colorB >>= 8
	colorA >>= 8

	red := colorR
	green := colorG << 8
	blue := colorB << 16
	alpha := colorA << 24

	tint := math.Float32frombits((alpha | blue | green | red) & 0xfeffffff)

	var changed bool

	switch ren.Drawable.(type) {
	case Triangle:
		setBufferValue(buffer, 0, w/2, &changed)
		//setBufferValue(buffer, 1, 0, &changed)
		setBufferValue(buffer, 2, tint, &changed)

		setBufferValue(buffer, 3, w, &changed)
		setBufferValue(buffer, 4, h, &changed)
		setBufferValue(buffer, 5, tint, &changed)

		//setBufferValue(buffer, 6, 0, &changed)
		setBufferValue(buffer, 7, h, &changed)
		setBufferValue(buffer, 8, tint, &changed)

	case Circle:
		//setBufferValue(buffer, 0, 0, &changed)
		//setBufferValue(buffer, 1, 0, &changed)
		setBufferValue(buffer, 2, tint, &changed)

		setBufferValue(buffer, 3, w, &changed)
		//setBufferValue(buffer, 4, 0, &changed)
		setBufferValue(buffer, 5, tint, &changed)

		setBufferValue(buffer, 6, w, &changed)
		setBufferValue(buffer, 7, h, &changed)
		setBufferValue(buffer, 8, tint, &changed)

		//setBufferValue(buffer, 9, 0, &changed)
		setBufferValue(buffer, 10, h, &changed)
		setBufferValue(buffer, 11, tint, &changed)

	case Rectangle:
		//setBufferValue(buffer, 0, 0, &changed)
		//setBufferValue(buffer, 1, 0, &changed)
		setBufferValue(buffer, 2, tint, &changed)

		setBufferValue(buffer, 3, w, &changed)
		//setBufferValue(buffer, 4, 0, &changed)
		setBufferValue(buffer, 5, tint, &changed)

		setBufferValue(buffer, 6, w, &changed)
		setBufferValue(buffer, 7, h, &changed)
		setBufferValue(buffer, 8, tint, &changed)

		//setBufferValue(buffer, 9, 0, &changed)
		setBufferValue(buffer, 10, h, &changed)
		setBufferValue(buffer, 11, tint, &changed)

	default:
		log.Println("Warning: type not supported")
	}

	return changed
}

func (l *legacyShader) Draw(ren *RenderComponent, space *SpaceComponent) {
	if l.lastBuffer != ren.buffer || ren.buffer == nil {
		l.updateBuffer(ren, space)

		Gl.BindBuffer(Gl.ARRAY_BUFFER, ren.buffer)
		Gl.VertexAttribPointer(l.inPosition, 2, Gl.FLOAT, false, 12, 0)
		Gl.VertexAttribPointer(l.inColor, 4, Gl.UNSIGNED_BYTE, true, 12, 8)

		l.lastBuffer = ren.buffer
	}

	if space.Rotation != 0 {
		sin, cos := math.Sincos(space.Rotation * math.Pi / 180)

		l.modelMatrix[0] = ren.Scale.X * cos
		l.modelMatrix[1] = ren.Scale.X * sin
		l.modelMatrix[3] = ren.Scale.Y * -sin
		l.modelMatrix[4] = ren.Scale.Y * cos
	} else {
		l.modelMatrix[0] = ren.Scale.X
		l.modelMatrix[1] = 0
		l.modelMatrix[3] = 0
		l.modelMatrix[4] = ren.Scale.Y
	}

	l.modelMatrix[6] = space.Position.X
	l.modelMatrix[7] = space.Position.Y

	Gl.UniformMatrix3fv(l.matrixModel, false, l.modelMatrix)

	switch ren.Drawable.(type) {
	case Triangle:
		Gl.Uniform2f(l.inRadius, 0, 0)
		Gl.BindBuffer(Gl.ELEMENT_ARRAY_BUFFER, l.indicesTrianglesVBO)
		Gl.DrawElements(Gl.TRIANGLES, 3, Gl.UNSIGNED_SHORT, 0)
	case Rectangle:
		Gl.Uniform2f(l.inRadius, 0, 0)
		Gl.BindBuffer(Gl.ELEMENT_ARRAY_BUFFER, l.indicesRectanglesVBO)
		Gl.DrawElements(Gl.TRIANGLES, 6, Gl.UNSIGNED_SHORT, 0)
	case Circle:
		Gl.Uniform2f(l.inRadius, (space.Width/2)/cam.z, (space.Height/2)/cam.z)
		Gl.Uniform2f(l.inCenter, space.Width/2, space.Height/2)
		Gl.BindBuffer(Gl.ELEMENT_ARRAY_BUFFER, l.indicesRectanglesVBO)
		Gl.DrawElements(Gl.TRIANGLES, 6, Gl.UNSIGNED_SHORT, 0)
	default:
		log.Println("Warning: type not supported")
	}
}

func (l *legacyShader) Post() {
	l.lastBuffer = nil

	// Cleanup
	Gl.DisableVertexAttribArray(l.inPosition)
	Gl.DisableVertexAttribArray(l.inColor)
}

var (
	DefaultShader   = &basicShader{cameraEnabled: true}
	HUDShader       = &basicShader{cameraEnabled: false}
	LegacyShader    = &legacyShader{cameraEnabled: true}
	LegacyHUDShader = &legacyShader{cameraEnabled: false}
	shadersSet      bool
)

func initShaders() {
	if !shadersSet {
		DefaultShader.Initialize()
		HUDShader.Initialize()
		LegacyShader.Initialize()
		LegacyHUDShader.Initialize()

		shadersSet = true
	}
}
