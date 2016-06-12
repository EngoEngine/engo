package common

import (
	"fmt"
	"image/color"
	"log"
	"strings"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/gl"
	"github.com/engoengine/math"
)

// UnicodeCap is the amount of unicode characters the fonts will be able to use, starting from index 0.
var UnicodeCap = 200

const bufferSize = 10000

type Shader interface {
	Setup(*ecs.World) error
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

	camera        *CameraSystem
	cameraEnabled bool
}

func (s *basicShader) Setup(w *ecs.World) error {
	var err error
	s.program, err = LoadShader(`
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

	if err != nil {
		return err
	}

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
	s.indexVBO = engo.Gl.CreateBuffer()
	engo.Gl.BindBuffer(engo.Gl.ELEMENT_ARRAY_BUFFER, s.indexVBO)
	engo.Gl.BufferData(engo.Gl.ELEMENT_ARRAY_BUFFER, s.indices, engo.Gl.STATIC_DRAW)

	// Define things that should be read from the texture buffer
	s.inPosition = engo.Gl.GetAttribLocation(s.program, "in_Position")
	s.inTexCoords = engo.Gl.GetAttribLocation(s.program, "in_TexCoords")
	s.inColor = engo.Gl.GetAttribLocation(s.program, "in_Color")

	// Define things that should be set per draw
	s.matrixProjection = engo.Gl.GetUniformLocation(s.program, "matrixProjection")
	s.matrixView = engo.Gl.GetUniformLocation(s.program, "matrixView")
	s.matrixModel = engo.Gl.GetUniformLocation(s.program, "matrixModel")

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

	if s.cameraEnabled {
		for _, system := range w.Systems() {
			switch sys := system.(type) {
			case *CameraSystem:
				s.camera = sys
			}
		}
		if s.camera == nil {
			log.Println("WARNING: BasicShader has CameraEnabled, but CameraSystem was not found")
		}
	}

	return nil
}

func (s *basicShader) Pre() {
	engo.Gl.Enable(engo.Gl.BLEND)
	engo.Gl.BlendFunc(engo.Gl.SRC_ALPHA, engo.Gl.ONE_MINUS_SRC_ALPHA)

	// Enable shader and buffer, enable attributes in shader
	engo.Gl.UseProgram(s.program)
	engo.Gl.BindBuffer(engo.Gl.ELEMENT_ARRAY_BUFFER, s.indexVBO)
	engo.Gl.EnableVertexAttribArray(s.inPosition)
	engo.Gl.EnableVertexAttribArray(s.inTexCoords)
	engo.Gl.EnableVertexAttribArray(s.inColor)

	if engo.ScaleOnResize() {
		s.projectionMatrix[0] = 1 / (engo.GameWidth() / 2)
		s.projectionMatrix[4] = 1 / (-engo.GameHeight() / 2)
	} else {
		s.projectionMatrix[0] = 1 / (engo.CanvasWidth() / 2)
		s.projectionMatrix[4] = 1 / (-engo.CanvasHeight() / 2)
	}

	if s.cameraEnabled {
		s.viewMatrix[1], s.viewMatrix[0] = math.Sincos(s.camera.angle * math.Pi / 180)
		s.viewMatrix[3] = -s.viewMatrix[1]
		s.viewMatrix[4] = s.viewMatrix[0]
		s.viewMatrix[6] = -s.camera.x
		s.viewMatrix[7] = -s.camera.y
		s.viewMatrix[8] = s.camera.z
	} else {
		s.viewMatrix[6] = -1 / s.projectionMatrix[0]
		s.viewMatrix[7] = 1 / s.projectionMatrix[4]
	}

	engo.Gl.UniformMatrix3fv(s.matrixProjection, false, s.projectionMatrix)
	engo.Gl.UniformMatrix3fv(s.matrixView, false, s.viewMatrix)
}

func (s *basicShader) Draw(ren *RenderComponent, space *SpaceComponent) {
	if s.lastBuffer != ren.buffer || ren.buffer == nil {
		s.updateBuffer(ren, space)

		engo.Gl.BindBuffer(engo.Gl.ARRAY_BUFFER, ren.buffer)
		engo.Gl.VertexAttribPointer(s.inPosition, 2, engo.Gl.FLOAT, false, 20, 0)
		engo.Gl.VertexAttribPointer(s.inTexCoords, 2, engo.Gl.FLOAT, false, 20, 8)
		engo.Gl.VertexAttribPointer(s.inColor, 4, engo.Gl.UNSIGNED_BYTE, true, 20, 16)

		s.lastBuffer = ren.buffer
	}

	if s.lastTexture != ren.Drawable.Texture() {
		engo.Gl.BindTexture(engo.Gl.TEXTURE_2D, ren.Drawable.Texture())

		s.lastTexture = ren.Drawable.Texture()
	}

	if s.lastRepeating != ren.Repeat {
		var val int
		switch ren.Repeat {
		case ClampToEdge:
			val = engo.Gl.CLAMP_TO_EDGE
		case ClampToBorder:
			val = engo.Gl.CLAMP_TO_EDGE
		case Repeat:
			val = engo.Gl.REPEAT
		case MirroredRepeat:
			val = engo.Gl.MIRRORED_REPEAT
		}

		engo.Gl.TexParameteri(engo.Gl.TEXTURE_2D, engo.Gl.TEXTURE_WRAP_S, val)
		engo.Gl.TexParameteri(engo.Gl.TEXTURE_2D, engo.Gl.TEXTURE_WRAP_T, val)
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

	engo.Gl.UniformMatrix3fv(s.matrixModel, false, s.modelMatrix)

	engo.Gl.DrawElements(engo.Gl.TRIANGLES, 6, engo.Gl.UNSIGNED_SHORT, 0)
}

func (s *basicShader) Post() {
	s.lastTexture = nil
	s.lastBuffer = nil

	// Cleanup
	engo.Gl.DisableVertexAttribArray(s.inPosition)
	engo.Gl.DisableVertexAttribArray(s.inTexCoords)
	engo.Gl.DisableVertexAttribArray(s.inColor)

	engo.Gl.BindTexture(engo.Gl.TEXTURE_2D, nil)
	engo.Gl.BindBuffer(engo.Gl.ARRAY_BUFFER, nil)
	engo.Gl.BindBuffer(engo.Gl.ELEMENT_ARRAY_BUFFER, nil)

	engo.Gl.Disable(engo.Gl.BLEND)
}

func (s *basicShader) updateBuffer(ren *RenderComponent, space *SpaceComponent) {
	if len(ren.bufferContent) == 0 {
		ren.bufferContent = make([]float32, 20) // because we add 20 elements to it
	}

	if changed := s.generateBufferContent(ren, space, ren.bufferContent); !changed {
		return
	}

	if ren.buffer == nil {
		ren.buffer = engo.Gl.CreateBuffer()
	}
	engo.Gl.BindBuffer(engo.Gl.ARRAY_BUFFER, ren.buffer)
	engo.Gl.BufferData(engo.Gl.ARRAY_BUFFER, ren.bufferContent, engo.Gl.STATIC_DRAW)
}

func (s *basicShader) generateBufferContent(ren *RenderComponent, space *SpaceComponent, buffer []float32) bool {
	// We shouldn't use SpaceComponent to get width/height, because this usually already contains the Scale (which
	// is being added elsewhere, so we don't want to over-do it)
	w := ren.Drawable.Width()
	h := ren.Drawable.Height()

	tint := colorToFloat32(ren.Color)

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

	indicesRectangles    []uint16
	indicesRectanglesVBO *gl.Buffer

	inPosition int
	inColor    int

	matrixProjection *gl.UniformLocation
	matrixView       *gl.UniformLocation
	matrixModel      *gl.UniformLocation
	inRadius         *gl.UniformLocation
	inCenter         *gl.UniformLocation
	inViewport       *gl.UniformLocation
	inBorderWidth    *gl.UniformLocation
	inBorderColor    *gl.UniformLocation

	projectionMatrix []float32
	viewMatrix       []float32
	modelMatrix      []float32

	camera        *CameraSystem
	cameraEnabled bool

	lastBuffer *gl.Buffer
}

func (l *legacyShader) Setup(w *ecs.World) error {
	var err error
	l.program, err = LoadShader(`
attribute vec2 in_Position;
attribute vec4 in_Color;

uniform mat3 matrixProjection;
uniform mat3 matrixView;
uniform mat3 matrixModel;
uniform vec2 in_Radius;
uniform vec2 in_Center;
uniform vec2 in_Viewport;
uniform float in_BorderWidth;
uniform vec4 in_BorderColor;

varying vec4 var_Color;
varying vec2 var_Radius;
varying vec2 var_Center;
varying float var_BorderWidth;
varying vec4 var_BorderColor;

void main() {
  var_Color = in_Color;

  vec3 matr = matrixProjection * matrixView * matrixModel * vec3(in_Position, 1.0);
  gl_Position = vec4(matr.xy, 0, matr.z);

  if (in_Radius.x > 0.0 && in_Radius.y > 0.0)
  {
    var_Radius = in_Radius;
    var_BorderWidth = in_BorderWidth;
    var_BorderColor = in_BorderColor;

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
varying float var_BorderWidth;
varying vec4 var_BorderColor;

void main (void) {
  gl_FragColor = var_Color;

  float halfBorder = var_BorderWidth / 2.0;

  if (var_Radius.x > 0.0 && var_Radius.y > 0.0)
  {
    if (pow(gl_FragCoord.x - var_Center.x, 2.0) / pow(var_Radius.x - halfBorder, 2.0) + pow(gl_FragCoord.y - var_Center.y, 2.0) / pow(var_Radius.y - halfBorder, 2.0) > 1.0)
    {
      if (pow(gl_FragCoord.x - var_Center.x, 2.0) / pow(var_Radius.x + halfBorder, 2.0) + pow(gl_FragCoord.y - var_Center.y, 2.0) / pow(var_Radius.y + halfBorder, 2.0) > 1.0)
	  {
	    gl_FragColor.w = 0.0;
	  } else {
	    gl_FragColor = var_BorderColor;
	  }
    }
  }
}`)

	if err != nil {
		return err
	}

	// Create and populate indices buffer
	l.indicesRectangles = []uint16{0, 1, 2, 0, 2, 3}
	l.indicesRectanglesVBO = engo.Gl.CreateBuffer()
	engo.Gl.BindBuffer(engo.Gl.ELEMENT_ARRAY_BUFFER, l.indicesRectanglesVBO)
	engo.Gl.BufferData(engo.Gl.ELEMENT_ARRAY_BUFFER, l.indicesRectangles, engo.Gl.STATIC_DRAW)

	// Define things that should be read from the texture buffer
	l.inPosition = engo.Gl.GetAttribLocation(l.program, "in_Position")
	l.inColor = engo.Gl.GetAttribLocation(l.program, "in_Color")

	// Define things that should be set per draw
	l.matrixProjection = engo.Gl.GetUniformLocation(l.program, "matrixProjection")
	l.matrixView = engo.Gl.GetUniformLocation(l.program, "matrixView")
	l.matrixModel = engo.Gl.GetUniformLocation(l.program, "matrixModel")
	l.inRadius = engo.Gl.GetUniformLocation(l.program, "in_Radius")
	l.inCenter = engo.Gl.GetUniformLocation(l.program, "in_Center")
	l.inViewport = engo.Gl.GetUniformLocation(l.program, "in_Viewport")
	l.inBorderWidth = engo.Gl.GetUniformLocation(l.program, "in_BorderWidth")
	l.inBorderColor = engo.Gl.GetUniformLocation(l.program, "in_BorderColor")

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

	if l.cameraEnabled {
		for _, system := range w.Systems() {
			switch sys := system.(type) {
			case *CameraSystem:
				l.camera = sys
			}
		}
		if l.camera == nil {
			log.Println("WARNING: BasicShader has CameraEnabled, but CameraSystem was not found")
		}
	}

	return nil
}

func (l *legacyShader) Pre() {
	engo.Gl.Enable(engo.Gl.BLEND)
	engo.Gl.BlendFunc(engo.Gl.SRC_ALPHA, engo.Gl.ONE_MINUS_SRC_ALPHA)

	// Bind shader and buffer, enable attributes
	engo.Gl.UseProgram(l.program)
	engo.Gl.EnableVertexAttribArray(l.inPosition)
	engo.Gl.EnableVertexAttribArray(l.inColor)

	if engo.ScaleOnResize() {
		l.projectionMatrix[0] = 1 / (engo.GameWidth() / 2)
		l.projectionMatrix[4] = 1 / (-engo.GameHeight() / 2)
	} else {
		l.projectionMatrix[0] = 1 / (engo.CanvasWidth() / 2)
		l.projectionMatrix[4] = 1 / (-engo.CanvasHeight() / 2)
	}

	if l.cameraEnabled {
		l.viewMatrix[1], l.viewMatrix[0] = math.Sincos(l.camera.angle * math.Pi / 180)
		l.viewMatrix[3] = -l.viewMatrix[1]
		l.viewMatrix[4] = l.viewMatrix[0]
		l.viewMatrix[6] = -l.camera.x
		l.viewMatrix[7] = -l.camera.y
		l.viewMatrix[8] = l.camera.z
	} else {
		l.viewMatrix[6] = -1 / l.projectionMatrix[0]
		l.viewMatrix[7] = 1 / l.projectionMatrix[4]
	}

	engo.Gl.UniformMatrix3fv(l.matrixProjection, false, l.projectionMatrix)
	engo.Gl.UniformMatrix3fv(l.matrixView, false, l.viewMatrix)
	engo.Gl.Uniform2f(l.inViewport, engo.GameWidth(), engo.GameHeight()) // TODO: canvasWidth/Height
}

func (l *legacyShader) updateBuffer(ren *RenderComponent, space *SpaceComponent) {
	if len(ren.bufferContent) == 0 {
		ren.bufferContent = make([]float32, l.computeBufferSize(ren.Drawable)) // because we add at most this many elements to it
	}
	if changed := l.generateBufferContent(ren, space, ren.bufferContent); !changed {
		return
	}

	if ren.buffer == nil {
		ren.buffer = engo.Gl.CreateBuffer()
	}
	engo.Gl.BindBuffer(engo.Gl.ARRAY_BUFFER, ren.buffer)
	engo.Gl.BufferData(engo.Gl.ARRAY_BUFFER, ren.bufferContent, engo.Gl.STATIC_DRAW)
}

func (l *legacyShader) computeBufferSize(draw Drawable) int {
	switch shape := draw.(type) {
	case Triangle:
		return 18
	case Rectangle:
		return 24
	case Circle:
		return 12
	case ComplexTriangles:
		return len(shape.Points) * 6
	default:
		return 0
	}
}

func (l *legacyShader) generateBufferContent(ren *RenderComponent, space *SpaceComponent, buffer []float32) bool {
	w := space.Width
	h := space.Height

	var changed bool

	tint := colorToFloat32(ren.Color)

	switch shape := ren.Drawable.(type) {
	case Triangle:
		switch shape.TriangleType {
		case TriangleIsosceles:
			setBufferValue(buffer, 0, w/2, &changed)
			//setBufferValue(buffer, 1, 0, &changed)
			setBufferValue(buffer, 2, tint, &changed)

			setBufferValue(buffer, 3, w, &changed)
			setBufferValue(buffer, 4, h, &changed)
			setBufferValue(buffer, 5, tint, &changed)

			//setBufferValue(buffer, 6, 0, &changed)
			setBufferValue(buffer, 7, h, &changed)
			setBufferValue(buffer, 8, tint, &changed)

			if shape.BorderWidth > 0 {
				borderTint := colorToFloat32(shape.BorderColor)

				setBufferValue(buffer, 9, w/2, &changed)
				//setBufferValue(buffer, 10, 0, &changed)
				setBufferValue(buffer, 11, borderTint, &changed)

				setBufferValue(buffer, 12, w, &changed)
				setBufferValue(buffer, 13, h, &changed)
				setBufferValue(buffer, 14, borderTint, &changed)

				//setBufferValue(buffer, 15, 0, &changed)
				setBufferValue(buffer, 16, h, &changed)
				setBufferValue(buffer, 17, borderTint, &changed)
			}
		case TriangleRight:
			//setBufferValue(buffer, 0, 0, &changed)
			//setBufferValue(buffer, 1, 0, &changed)
			setBufferValue(buffer, 2, tint, &changed)

			setBufferValue(buffer, 3, w, &changed)
			setBufferValue(buffer, 4, h, &changed)
			setBufferValue(buffer, 5, tint, &changed)

			//setBufferValue(buffer, 6, 0, &changed)
			setBufferValue(buffer, 7, h, &changed)
			setBufferValue(buffer, 8, tint, &changed)

			if shape.BorderWidth > 0 {
				borderTint := colorToFloat32(shape.BorderColor)

				//setBufferValue(buffer, 9, 0, &changed)
				//setBufferValue(buffer, 10, 0, &changed)
				setBufferValue(buffer, 11, borderTint, &changed)

				setBufferValue(buffer, 12, w, &changed)
				setBufferValue(buffer, 13, h, &changed)
				setBufferValue(buffer, 14, borderTint, &changed)

				//setBufferValue(buffer, 15, 0, &changed)
				setBufferValue(buffer, 16, h, &changed)
				setBufferValue(buffer, 17, borderTint, &changed)
			}
		}

	case Circle:
		halfWidth := shape.BorderWidth / 2
		setBufferValue(buffer, 0, -halfWidth, &changed)
		setBufferValue(buffer, 1, -halfWidth, &changed)
		setBufferValue(buffer, 2, tint, &changed)

		setBufferValue(buffer, 3, w+halfWidth, &changed)
		setBufferValue(buffer, 4, -halfWidth, &changed)
		setBufferValue(buffer, 5, tint, &changed)

		setBufferValue(buffer, 6, w+halfWidth, &changed)
		setBufferValue(buffer, 7, h+halfWidth, &changed)
		setBufferValue(buffer, 8, tint, &changed)

		setBufferValue(buffer, 9, -halfWidth, &changed)
		setBufferValue(buffer, 10, h+halfWidth, &changed)
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

		if shape.BorderWidth > 0 {
			borderTint := colorToFloat32(shape.BorderColor)

			//setBufferValue(buffer, 12, 0, &changed)
			//setBufferValue(buffer, 13, 0, &changed)
			setBufferValue(buffer, 14, borderTint, &changed)

			setBufferValue(buffer, 15, w, &changed)
			//setBufferValue(buffer, 16, 0, &changed)
			setBufferValue(buffer, 17, borderTint, &changed)

			setBufferValue(buffer, 18, w, &changed)
			setBufferValue(buffer, 19, h, &changed)
			setBufferValue(buffer, 20, borderTint, &changed)

			//setBufferValue(buffer, 21, 0, &changed)
			setBufferValue(buffer, 22, h, &changed)
			setBufferValue(buffer, 23, borderTint, &changed)
		}

	case ComplexTriangles:
		var index int
		for _, point := range shape.Points {
			setBufferValue(buffer, index, point.X*w, &changed)
			setBufferValue(buffer, index+1, point.Y*h, &changed)
			setBufferValue(buffer, index+2, tint, &changed)
			index += 3
		}

		if shape.BorderWidth > 0 {
			borderTint := colorToFloat32(shape.BorderColor)

			for _, point := range shape.Points {
				setBufferValue(buffer, index, point.X*w, &changed)
				setBufferValue(buffer, index+1, point.Y*h, &changed)
				setBufferValue(buffer, index+2, borderTint, &changed)
				index += 3
			}
		}
	default:
		unsupportedType(ren.Drawable)
	}

	return changed
}

func (l *legacyShader) Draw(ren *RenderComponent, space *SpaceComponent) {
	if l.lastBuffer != ren.buffer || ren.buffer == nil {
		l.updateBuffer(ren, space)

		engo.Gl.BindBuffer(engo.Gl.ARRAY_BUFFER, ren.buffer)
		engo.Gl.VertexAttribPointer(l.inPosition, 2, engo.Gl.FLOAT, false, 12, 0)
		engo.Gl.VertexAttribPointer(l.inColor, 4, engo.Gl.UNSIGNED_BYTE, true, 12, 8)

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

	engo.Gl.UniformMatrix3fv(l.matrixModel, false, l.modelMatrix)

	switch shape := ren.Drawable.(type) {
	case Triangle:
		engo.Gl.Uniform2f(l.inRadius, 0, 0)
		engo.Gl.DrawArrays(engo.Gl.TRIANGLES, 0, 3)

		if shape.BorderWidth > 0 {
			borderWidth := shape.BorderWidth
			if l.cameraEnabled {
				borderWidth /= l.camera.z
			}
			engo.Gl.LineWidth(borderWidth)
			engo.Gl.DrawArrays(engo.Gl.LINE_LOOP, 3, 3)
		}
	case Rectangle:
		engo.Gl.Uniform2f(l.inRadius, 0, 0)
		engo.Gl.BindBuffer(engo.Gl.ELEMENT_ARRAY_BUFFER, l.indicesRectanglesVBO)
		engo.Gl.DrawElements(engo.Gl.TRIANGLES, 6, engo.Gl.UNSIGNED_SHORT, 0)

		if shape.BorderWidth > 0 {
			borderWidth := shape.BorderWidth
			if l.cameraEnabled {
				borderWidth /= l.camera.z
			}
			engo.Gl.LineWidth(borderWidth)
			engo.Gl.DrawArrays(engo.Gl.LINE_LOOP, 4, 4)
		}
	case Circle:
		engo.Gl.Uniform1f(l.inBorderWidth, shape.BorderWidth/l.camera.z)
		if shape.BorderWidth > 0 {
			r, g, b, a := shape.BorderColor.RGBA()
			engo.Gl.Uniform4f(l.inBorderColor, float32(r>>8), float32(g>>8), float32(b>>8), float32(a>>8))
		}
		engo.Gl.Uniform2f(l.inRadius, (space.Width/2)/l.camera.z, (space.Height/2)/l.camera.z)
		engo.Gl.Uniform2f(l.inCenter, space.Width/2, space.Height/2)
		engo.Gl.BindBuffer(engo.Gl.ELEMENT_ARRAY_BUFFER, l.indicesRectanglesVBO)
		engo.Gl.DrawElements(engo.Gl.TRIANGLES, 6, engo.Gl.UNSIGNED_SHORT, 0)
	case ComplexTriangles:
		engo.Gl.Uniform2f(l.inRadius, 0, 0)
		engo.Gl.DrawArrays(engo.Gl.TRIANGLES, 0, len(shape.Points))

		if shape.BorderWidth > 0 {
			borderWidth := shape.BorderWidth
			if l.cameraEnabled {
				borderWidth /= l.camera.z
			}
			engo.Gl.LineWidth(borderWidth)
			engo.Gl.DrawArrays(engo.Gl.LINE_LOOP, len(shape.Points), len(shape.Points))
		}
	default:
		unsupportedType(ren.Drawable)
	}
}

func (l *legacyShader) Post() {
	l.lastBuffer = nil

	// Cleanup
	engo.Gl.DisableVertexAttribArray(l.inPosition)
	engo.Gl.DisableVertexAttribArray(l.inColor)

	engo.Gl.BindBuffer(engo.Gl.ARRAY_BUFFER, nil)
	engo.Gl.BindBuffer(engo.Gl.ELEMENT_ARRAY_BUFFER, nil)

	engo.Gl.Disable(engo.Gl.BLEND)
}

type textShader struct {
	program *gl.Program

	indicesRectangles    []uint16
	indicesRectanglesVBO *gl.Buffer

	inPosition  int
	inTexCoords int
	inColor     int

	matrixProjection *gl.UniformLocation
	matrixView       *gl.UniformLocation
	matrixModel      *gl.UniformLocation

	projectionMatrix []float32
	viewMatrix       []float32
	modelMatrix      []float32

	camera        *CameraSystem
	cameraEnabled bool

	lastBuffer  *gl.Buffer
	lastTexture *gl.Texture
}

func (l *textShader) Setup(w *ecs.World) error {
	var err error
	l.program, err = LoadShader(`
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

	if err != nil {
		return err
	}

	// Create and populate indices buffer
	l.indicesRectangles = make([]uint16, 6*bufferSize)
	for i, j := 0, 0; i < bufferSize*6; i, j = i+6, j+4 {
		l.indicesRectangles[i+0] = uint16(j + 0)
		l.indicesRectangles[i+1] = uint16(j + 1)
		l.indicesRectangles[i+2] = uint16(j + 2)
		l.indicesRectangles[i+3] = uint16(j + 0)
		l.indicesRectangles[i+4] = uint16(j + 2)
		l.indicesRectangles[i+5] = uint16(j + 3)
	}
	l.indicesRectanglesVBO = engo.Gl.CreateBuffer()
	engo.Gl.BindBuffer(engo.Gl.ELEMENT_ARRAY_BUFFER, l.indicesRectanglesVBO)
	engo.Gl.BufferData(engo.Gl.ELEMENT_ARRAY_BUFFER, l.indicesRectangles, engo.Gl.STATIC_DRAW)

	// Define things that should be read from the texture buffer
	l.inPosition = engo.Gl.GetAttribLocation(l.program, "in_Position")
	l.inTexCoords = engo.Gl.GetAttribLocation(l.program, "in_TexCoords")
	l.inColor = engo.Gl.GetAttribLocation(l.program, "in_Color")

	// Define things that should be set per draw
	l.matrixProjection = engo.Gl.GetUniformLocation(l.program, "matrixProjection")
	l.matrixView = engo.Gl.GetUniformLocation(l.program, "matrixView")
	l.matrixModel = engo.Gl.GetUniformLocation(l.program, "matrixModel")

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

	if l.cameraEnabled {
		for _, system := range w.Systems() {
			switch sys := system.(type) {
			case *CameraSystem:
				l.camera = sys
			}
		}
		if l.camera == nil {
			log.Println("WARNING: TextShader has CameraEnabled, but CameraSystem was not found")
		}
	}

	return nil
}

func (l *textShader) Pre() {
	engo.Gl.Enable(engo.Gl.BLEND)
	engo.Gl.BlendFunc(engo.Gl.SRC_ALPHA, engo.Gl.ONE_MINUS_SRC_ALPHA)

	// Bind shader and buffer, enable attributes
	engo.Gl.UseProgram(l.program)
	engo.Gl.BindBuffer(engo.Gl.ELEMENT_ARRAY_BUFFER, l.indicesRectanglesVBO)
	engo.Gl.EnableVertexAttribArray(l.inPosition)
	engo.Gl.EnableVertexAttribArray(l.inTexCoords)
	engo.Gl.EnableVertexAttribArray(l.inColor)

	if engo.ScaleOnResize() {
		l.projectionMatrix[0] = 1 / (engo.GameWidth() / 2)
		l.projectionMatrix[4] = 1 / (-engo.GameHeight() / 2)
	} else {
		l.projectionMatrix[0] = 1 / (engo.CanvasWidth() / 2)
		l.projectionMatrix[4] = 1 / (-engo.CanvasHeight() / 2)
	}

	if l.cameraEnabled {
		l.viewMatrix[1], l.viewMatrix[0] = math.Sincos(l.camera.angle * math.Pi / 180)
		l.viewMatrix[3] = -l.viewMatrix[1]
		l.viewMatrix[4] = l.viewMatrix[0]
		l.viewMatrix[6] = -l.camera.x
		l.viewMatrix[7] = -l.camera.y
		l.viewMatrix[8] = l.camera.z
	} else {
		l.viewMatrix[6] = -1 / l.projectionMatrix[0]
		l.viewMatrix[7] = 1 / l.projectionMatrix[4]
	}

	engo.Gl.UniformMatrix3fv(l.matrixProjection, false, l.projectionMatrix)
	engo.Gl.UniformMatrix3fv(l.matrixView, false, l.viewMatrix)
}

func (l *textShader) updateBuffer(ren *RenderComponent, space *SpaceComponent) {
	txt, ok := ren.Drawable.(Text)
	if !ok {
		unsupportedType(ren.Drawable)
		return
	}

	if len(ren.bufferContent) < 20*len(txt.Text) {
		ren.bufferContent = make([]float32, 20*len(txt.Text)) // TODO: update this to actual value?
	}
	if changed := l.generateBufferContent(ren, space, ren.bufferContent); !changed {
		return
	}

	if ren.buffer == nil {
		ren.buffer = engo.Gl.CreateBuffer()
	}
	engo.Gl.BindBuffer(engo.Gl.ARRAY_BUFFER, ren.buffer)
	engo.Gl.BufferData(engo.Gl.ARRAY_BUFFER, ren.bufferContent, engo.Gl.STATIC_DRAW)
}

func (l *textShader) generateBufferContent(ren *RenderComponent, space *SpaceComponent, buffer []float32) bool {
	var changed bool

	tint := colorToFloat32(ren.Color)
	txt, ok := ren.Drawable.(Text)
	if !ok {
		unsupportedType(ren.Drawable)
		return false
	}

	atlas, ok := atlasCache[txt.Font.URL]
	if !ok {
		// Generate texture first
		atlas = txt.Font.generateFontAtlas(UnicodeCap)
		atlasCache[txt.Font.URL] = atlas
	}

	var currentX float32
	var currentY float32

	var modifier float32 = 1
	if txt.RightToLeft {
		modifier = -1
	}

	for index, char := range txt.Text {
		// TODO: this might not work for all characters
		switch {
		case char == '\n':
			currentX = 0
			currentY += atlas.Height[char] + txt.LineSpacing*atlas.Height[char]
			continue
		case char < 32: // all system stuff should be ignored
			continue
		}

		offset := 20 * index

		// These five are at 0, 0:
		setBufferValue(buffer, 0+offset, currentX, &changed)
		setBufferValue(buffer, 1+offset, currentY, &changed)
		setBufferValue(buffer, 2+offset, atlas.XLocation[char]/atlas.TotalWidth, &changed)
		setBufferValue(buffer, 3+offset, atlas.YLocation[char]/atlas.TotalHeight, &changed)
		setBufferValue(buffer, 4+offset, tint, &changed)

		// These five are at 1, 0:
		setBufferValue(buffer, 5+offset, currentX+atlas.Width[char], &changed)
		setBufferValue(buffer, 6+offset, currentY, &changed)
		setBufferValue(buffer, 7+offset, (atlas.XLocation[char]+atlas.Width[char])/atlas.TotalWidth, &changed)
		setBufferValue(buffer, 8+offset, atlas.YLocation[char]/atlas.TotalHeight, &changed)
		setBufferValue(buffer, 9+offset, tint, &changed)

		// These five are at 1, 1:
		setBufferValue(buffer, 10+offset, currentX+atlas.Width[char], &changed)
		setBufferValue(buffer, 11+offset, currentY+atlas.Height[char], &changed)
		setBufferValue(buffer, 12+offset, (atlas.XLocation[char]+atlas.Width[char])/atlas.TotalWidth, &changed)
		setBufferValue(buffer, 13+offset, (atlas.YLocation[char]+atlas.Height[char])/atlas.TotalHeight, &changed)
		setBufferValue(buffer, 14+offset, tint, &changed)

		// These five are at 0, 1:
		setBufferValue(buffer, 15+offset, currentX, &changed)
		setBufferValue(buffer, 16+offset, currentY+atlas.Height[char], &changed)
		setBufferValue(buffer, 17+offset, atlas.XLocation[char]/atlas.TotalWidth, &changed)
		setBufferValue(buffer, 18+offset, (atlas.YLocation[char]+atlas.Height[char])/atlas.TotalHeight, &changed)
		setBufferValue(buffer, 19+offset, tint, &changed)

		currentX += modifier * (atlas.Width[char] + float32(txt.Font.Size)*txt.LetterSpacing)
	}

	return changed
}

func (l *textShader) Draw(ren *RenderComponent, space *SpaceComponent) {
	if l.lastBuffer != ren.buffer || ren.buffer == nil {
		l.updateBuffer(ren, space)

		engo.Gl.BindBuffer(engo.Gl.ARRAY_BUFFER, ren.buffer)
		engo.Gl.VertexAttribPointer(l.inPosition, 2, engo.Gl.FLOAT, false, 20, 0)
		engo.Gl.VertexAttribPointer(l.inTexCoords, 2, engo.Gl.FLOAT, false, 20, 8)
		engo.Gl.VertexAttribPointer(l.inColor, 4, engo.Gl.UNSIGNED_BYTE, true, 20, 16)

		l.lastBuffer = ren.buffer
	}

	txt, ok := ren.Drawable.(Text)
	if !ok {
		unsupportedType(ren.Drawable)
	}

	atlas, ok := atlasCache[txt.Font.URL]
	if !ok {
		// Generate texture first
		atlas = txt.Font.generateFontAtlas(UnicodeCap)
		atlasCache[txt.Font.URL] = atlas
	}

	if atlas.Texture != l.lastTexture {
		engo.Gl.BindTexture(engo.Gl.TEXTURE_2D, atlas.Texture)
		l.lastTexture = atlas.Texture
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

	engo.Gl.UniformMatrix3fv(l.matrixModel, false, l.modelMatrix)
	engo.Gl.DrawElements(engo.Gl.TRIANGLES, 6*len(txt.Text), engo.Gl.UNSIGNED_SHORT, 0)
}

func (l *textShader) Post() {
	l.lastBuffer = nil
	l.lastTexture = nil

	// Cleanup
	engo.Gl.DisableVertexAttribArray(l.inPosition)
	engo.Gl.DisableVertexAttribArray(l.inTexCoords)
	engo.Gl.DisableVertexAttribArray(l.inColor)

	engo.Gl.BindTexture(engo.Gl.TEXTURE_2D, nil)
	engo.Gl.BindBuffer(engo.Gl.ARRAY_BUFFER, nil)
	engo.Gl.BindBuffer(engo.Gl.ELEMENT_ARRAY_BUFFER, nil)

	engo.Gl.Disable(engo.Gl.BLEND)
}

// colorToFloat32 returns the float32 representation of the given color
func colorToFloat32(c color.Color) float32 {
	colorR, colorG, colorB, colorA := c.RGBA()
	colorR >>= 8
	colorG >>= 8
	colorB >>= 8
	colorA >>= 8

	red := colorR
	green := colorG << 8
	blue := colorB << 16
	alpha := colorA << 24

	return math.Float32frombits((alpha | blue | green | red) & 0xfeffffff)
}

var (
	DefaultShader   = &basicShader{cameraEnabled: true}
	HUDShader       = &basicShader{cameraEnabled: false}
	LegacyShader    = &legacyShader{cameraEnabled: true}
	LegacyHUDShader = &legacyShader{cameraEnabled: false}
	TextShader      = &textShader{cameraEnabled: true}
	TextHUDShader   = &textShader{cameraEnabled: false}
	shadersSet      bool
	atlasCache      = make(map[string]FontAtlas)
)

func initShaders(w *ecs.World) error {
	if !shadersSet {
		shaders := []Shader{
			DefaultShader,
			HUDShader,
			LegacyShader,
			LegacyHUDShader,
			TextShader,
			TextHUDShader,
		}
		var err error

		for _, shader := range shaders {
			err = shader.Setup(w)
			if err != nil {
				return err
			}
		}

		shadersSet = true
	}
	return nil
}

// LoadShader takes a Vertex-shader and Fragment-shader, compiles them and attaches them to a newly created glProgram.
// It will log possible compilation errors
func LoadShader(vertSrc, fragSrc string) (*gl.Program, error) {
	vertShader := engo.Gl.CreateShader(engo.Gl.VERTEX_SHADER)
	engo.Gl.ShaderSource(vertShader, vertSrc)
	engo.Gl.CompileShader(vertShader)
	if !engo.Gl.GetShaderiv(vertShader, engo.Gl.COMPILE_STATUS) {
		errorLog := engo.Gl.GetShaderInfoLog(vertShader)
		return nil, VertexShaderCompilationError{errorLog}
	}
	defer engo.Gl.DeleteShader(vertShader)

	fragShader := engo.Gl.CreateShader(engo.Gl.FRAGMENT_SHADER)
	engo.Gl.ShaderSource(fragShader, fragSrc)
	engo.Gl.CompileShader(fragShader)
	if !engo.Gl.GetShaderiv(fragShader, engo.Gl.COMPILE_STATUS) {
		errorLog := engo.Gl.GetShaderInfoLog(fragShader)
		return nil, FragmentShaderCompilationError{errorLog}
	}
	defer engo.Gl.DeleteShader(fragShader)

	program := engo.Gl.CreateProgram()
	engo.Gl.AttachShader(program, vertShader)
	engo.Gl.AttachShader(program, fragShader)
	engo.Gl.LinkProgram(program)

	return program, nil
}

// VertexShaderCompilationError is returned whenever the `LoadShader` method was unable to compile your Vertex-shader (GLSL)
type VertexShaderCompilationError struct {
	OpenGLError string
}

func (v VertexShaderCompilationError) Error() string {
	return fmt.Sprintf("an error occured compiling the vertex shader: %s", strings.Trim(v.OpenGLError, "\r\n"))
}

// FragmentShaderCompilationError is returned whenever the `LoadShader` method was unable to compile your Fragment-shader (GLSL)
type FragmentShaderCompilationError struct {
	OpenGLError string
}

func (f FragmentShaderCompilationError) Error() string {
	return fmt.Sprintf("an error occured compiling the fragment shader: %s", strings.Trim(f.OpenGLError, "\r\n"))
}
