package common

import (
	"fmt"
	"image/color"
	"log"
	"strings"
	"sync"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/math"
	"engo.io/gl"
)

// UnicodeCap is the amount of unicode characters the fonts will be able to use, starting from index 0.
var UnicodeCap = 200

const (
	// MaxSprites is the maximum number of sprites that can comprise a single batch.
	// 32767 is the max vertex index in OpenGL. Since each sprite has 4 vertices,
	// 32767 / 4 = 8191 max sprites.
	MaxSprites = 8191
	spriteSize = 20

	bufferSize = 10000

	defaultVertexShader = `
	attribute vec2 in_Position;
	attribute vec2 in_TexCoords;
	attribute vec4 in_Color;

	uniform mat3 matrixProjView;

	varying vec4 var_Color;
	varying vec2 var_TexCoords;

	void main() {
	  var_Color = in_Color;
	  var_TexCoords = in_TexCoords;

	  vec3 matr = matrixProjView * vec3(in_Position, 1.0);
	  gl_Position = vec4(matr.xy, 0, matr.z);
	}
`

	defaultFragmentShader = `
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
)

// Shader when implemented can be used in the RenderSystem as an OpenGl Shader.
//
// Setup holds the actual OpenGL shader data, and prepares any matrices and OpenGL calls for use.
//
// Pre is called just before the Draw step.
//
// Draw is the Draw step.
//
// Post is called just after the Draw step.
type Shader interface {
	Setup(*ecs.World) error
	Pre()
	Draw(*RenderComponent, *SpaceComponent)
	Post()
	SetCamera(*CameraSystem)
}

// CullingShader when implemented can be used in the RenderSystem to test if an entity should be rendered.
type CullingShader interface {
	Shader
	// PrepareCulling is called once per frame for the shader to initialize culling calculation.
	PrepareCulling()
	ShouldDraw(*RenderComponent, *SpaceComponent) bool
}

type basicShader struct {
	BatchSize int

	indices     []uint16
	indexBuffer *gl.Buffer
	program     *gl.Program

	vertices                     []float32
	vertexBuffer                 *gl.Buffer
	lastTexture                  *gl.Texture
	lastRepeating                TextureRepeating
	lastMagFilter, lastMinFilter ZoomFilter

	inPosition  int
	inTexCoords int
	inColor     int

	matrixProjView *gl.UniformLocation

	projectionMatrix *engo.Matrix
	viewMatrix       *engo.Matrix
	modelMatrix      *engo.Matrix
	cullingMatrix    *engo.Matrix

	camera        *CameraSystem
	cameraEnabled bool

	idx int
}

func (s *basicShader) Setup(w *ecs.World) error {
	if s.BatchSize > MaxSprites {
		return fmt.Errorf("%d is greater than the maximum batch size of %d", s.BatchSize, MaxSprites)
	}
	if s.BatchSize <= 0 {
		s.BatchSize = MaxSprites
	}
	// Create the vertex buffer for batching.
	s.vertices = make([]float32, s.BatchSize*spriteSize)
	s.vertexBuffer = engo.Gl.CreateBuffer()
	// Create and populate indices buffer. The size of the buffer depends on the batch size.
	// These should never change, so we can just initialize them once here and be done with it.
	numIndicies := s.BatchSize * 6
	s.indices = make([]uint16, numIndicies)
	for i, j := 0, 0; i < numIndicies; i, j = i+6, j+4 {
		s.indices[i+0] = uint16(j + 0)
		s.indices[i+1] = uint16(j + 1)
		s.indices[i+2] = uint16(j + 2)
		s.indices[i+3] = uint16(j + 0)
		s.indices[i+4] = uint16(j + 2)
		s.indices[i+5] = uint16(j + 3)
	}
	var err error
	s.program, err = LoadShader(defaultVertexShader, defaultFragmentShader)
	if err != nil {
		return err
	}
	s.indexBuffer = engo.Gl.CreateBuffer()
	engo.Gl.BindBuffer(engo.Gl.ELEMENT_ARRAY_BUFFER, s.indexBuffer)
	engo.Gl.BufferData(engo.Gl.ELEMENT_ARRAY_BUFFER, s.indices, engo.Gl.STATIC_DRAW)

	s.inPosition = engo.Gl.GetAttribLocation(s.program, "in_Position")
	s.inTexCoords = engo.Gl.GetAttribLocation(s.program, "in_TexCoords")
	s.inColor = engo.Gl.GetAttribLocation(s.program, "in_Color")

	s.matrixProjView = engo.Gl.GetUniformLocation(s.program, "matrixProjView")

	s.projectionMatrix = engo.IdentityMatrix()
	s.viewMatrix = engo.IdentityMatrix()
	s.modelMatrix = engo.IdentityMatrix()
	s.cullingMatrix = engo.IdentityMatrix()

	return nil
}

func (s *basicShader) Pre() {
	engo.Gl.Enable(engo.Gl.BLEND)
	engo.Gl.BlendFunc(engo.Gl.SRC_ALPHA, engo.Gl.ONE_MINUS_SRC_ALPHA)
	// Enable shader and buffer, enable attributes in shader
	engo.Gl.UseProgram(s.program)
	engo.Gl.BindBuffer(engo.Gl.ELEMENT_ARRAY_BUFFER, s.indexBuffer)
	engo.Gl.EnableVertexAttribArray(s.inPosition)
	engo.Gl.EnableVertexAttribArray(s.inTexCoords)
	engo.Gl.EnableVertexAttribArray(s.inColor)

	// The matrixProjView shader uniform is projection * view.
	// We do the multiplication on the CPU instead of sending each matrix to the shader and letting the GPU do the multiplication,
	// because it's likely faster to do the multiplication client side and send the result over the shader bus than to send two separate
	// buffers over the bus and then do the multiplication on the GPU.
	pv := s.projectionMatrix.Multiply(s.viewMatrix)
	engo.Gl.UniformMatrix3fv(s.matrixProjView, false, pv.Val[:])

	// Since we are batching client side, we only have one VBO, so we can just bind it now and use it for the entire frame.
	engo.Gl.BindBuffer(engo.Gl.ARRAY_BUFFER, s.vertexBuffer)
	engo.Gl.VertexAttribPointer(s.inPosition, 2, engo.Gl.FLOAT, false, 20, 0)
	engo.Gl.VertexAttribPointer(s.inTexCoords, 2, engo.Gl.FLOAT, false, 20, 8)
	engo.Gl.VertexAttribPointer(s.inColor, 4, engo.Gl.UNSIGNED_BYTE, true, 20, 16)
}

func (s *basicShader) PrepareCulling() {
	// (Re)initialize the projection matrix.
	s.projectionMatrix.Identity()
	if engo.ScaleOnResize() {
		s.projectionMatrix.Scale(1/(engo.GameWidth()/2), 1/(-engo.GameHeight()/2))
	} else {
		s.projectionMatrix.Scale(1/(engo.CanvasWidth()/(2*engo.CanvasScale())), 1/(-engo.CanvasHeight()/(2*engo.CanvasScale())))
	}
	// (Re)initialize the view matrix
	s.viewMatrix.Identity()
	if s.cameraEnabled {
		s.viewMatrix.Scale(1/s.camera.z, 1/s.camera.z)
		s.viewMatrix.Translate(-s.camera.x, -s.camera.y).Rotate(s.camera.angle)
	} else {
		scaleX, scaleY := s.projectionMatrix.ScaleComponent()
		s.viewMatrix.Translate(-1/scaleX, 1/scaleY)
	}
	s.cullingMatrix.Identity()
	s.cullingMatrix.Multiply(s.projectionMatrix).Multiply(s.viewMatrix)
	s.cullingMatrix.Scale(engo.GetGlobalScale().X, engo.GetGlobalScale().Y)
}

func (s *basicShader) ShouldDraw(rc *RenderComponent, sc *SpaceComponent) bool {
	tsc := SpaceComponent{
		Position: sc.Position,
		Width:    rc.Drawable.Width() * rc.Scale.X,
		Height:   rc.Drawable.Height() * rc.Scale.Y,
		Rotation: sc.Rotation,
	}

	c := tsc.Corners()
	c[0].MultiplyMatrixVector(s.cullingMatrix)
	c[1].MultiplyMatrixVector(s.cullingMatrix)
	c[2].MultiplyMatrixVector(s.cullingMatrix)
	c[3].MultiplyMatrixVector(s.cullingMatrix)

	return !((c[0].X < -1 && c[1].X < -1 && c[2].X < -1 && c[3].X < -1) || // All points left of the "viewport"
		(c[0].X > 1 && c[1].X > 1 && c[2].X > 1 && c[3].X > 1) || // All points right of the "viewport"
		(c[0].Y < -1 && c[1].Y < -1 && c[2].Y < -1 && c[3].Y < -1) || // All points above of the "viewport"
		(c[0].Y > 1 && c[1].Y > 1 && c[2].Y > 1 && c[3].Y > 1)) // All points below of the "viewport"
}

func (s *basicShader) Draw(ren *RenderComponent, space *SpaceComponent) {
	// If our texture (or any of its properties) has changed or we've reached the end of our buffer, flush before moving on.
	if s.lastTexture != ren.Drawable.Texture() {
		s.flush()
		engo.Gl.BindTexture(engo.Gl.TEXTURE_2D, ren.Drawable.Texture())
		s.lastTexture = ren.Drawable.Texture()
	} else if s.idx == len(s.vertices) {
		s.flush()
	}

	if s.lastRepeating != ren.Repeat {
		s.flush()
		var val int
		switch ren.Repeat {
		case NoRepeat:
			val = engo.Gl.CLAMP_TO_EDGE
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

	if s.lastMagFilter != ren.magFilter {
		s.flush()
		var val int
		switch ren.magFilter {
		case FilterNearest:
			val = engo.Gl.NEAREST
		case FilterLinear:
			val = engo.Gl.LINEAR
		}
		engo.Gl.TexParameteri(engo.Gl.TEXTURE_2D, engo.Gl.TEXTURE_MAG_FILTER, val)
	}

	if s.lastMinFilter != ren.minFilter {
		s.flush()
		var val int
		switch ren.minFilter {
		case FilterNearest:
			val = engo.Gl.NEAREST
		case FilterLinear:
			val = engo.Gl.LINEAR
		}
		engo.Gl.TexParameteri(engo.Gl.TEXTURE_2D, engo.Gl.TEXTURE_MIN_FILTER, val)
	}

	// Update the vertex buffer data.
	s.updateBuffer(ren, space)
	s.idx += 20
}

func (s *basicShader) Post() {
	s.flush()
	s.lastTexture = nil

	// Cleanup
	engo.Gl.DisableVertexAttribArray(s.inPosition)
	engo.Gl.DisableVertexAttribArray(s.inTexCoords)
	engo.Gl.DisableVertexAttribArray(s.inColor)

	engo.Gl.BindTexture(engo.Gl.TEXTURE_2D, nil)
	engo.Gl.BindBuffer(engo.Gl.ARRAY_BUFFER, nil)
	engo.Gl.BindBuffer(engo.Gl.ELEMENT_ARRAY_BUFFER, nil)

	engo.Gl.Disable(engo.Gl.BLEND)
}

func (s *basicShader) flush() {
	// If we haven't rendered anything yet, no point in flushing.
	if s.idx == 0 {
		return
	}
	engo.Gl.BufferData(engo.Gl.ARRAY_BUFFER, s.vertices, engo.Gl.STATIC_DRAW)
	// We only want to draw the indicies up to the number of sprites in the current batch.
	count := s.idx / 20 * 6
	engo.Gl.DrawElements(engo.Gl.TRIANGLES, count, engo.Gl.UNSIGNED_SHORT, 0)
	s.idx = 0
	// We need to reset the vertex buffer so that when we start drawing again, we don't accidentally use junk data.
	// The "simpler" way to do this would be to just create a new slice with make(), however that would cause the
	// previous slice to be marked for garbage collection and we'd prefer to keep the GC activity to a minimum.
	for i := range s.vertices {
		s.vertices[i] = 0
	}
}

func (s *basicShader) updateBuffer(ren *RenderComponent, space *SpaceComponent) {
	// For backwards compatibility, ren.Buffer is set to the VBO and ren.BufferContent
	// is set to the slice of the vertex buffer for the current sprite. This same slice is
	// populated with vertex data via generateBufferContent.
	ren.Buffer = s.vertexBuffer
	ren.BufferContent = s.vertices[s.idx : s.idx+20]
	s.generateBufferContent(ren, space, ren.BufferContent)
}

func (s *basicShader) makeModelMatrix(ren *RenderComponent, space *SpaceComponent) *engo.Matrix {
	// Instead of creating a new model matrix every time, we instead store a global one as a struct member
	// and just reset it for every sprite. This prevents us from allocating a bunch of new Matrix instances in memory
	// ultimately saving on GC activity.
	s.modelMatrix.Identity().Scale(engo.GetGlobalScale().X, engo.GetGlobalScale().Y).Translate(space.Position.X, space.Position.Y)
	if space.Rotation != 0 {
		s.modelMatrix.Rotate(space.Rotation)
	}
	s.modelMatrix.Scale(ren.Scale.X, ren.Scale.Y)
	return s.modelMatrix
}

func (s *basicShader) generateBufferContent(ren *RenderComponent, space *SpaceComponent, buffer []float32) bool {
	// We shouldn't use SpaceComponent to get width/height, because this usually already contains the Scale (which
	// is being added elsewhere, so we don't want to over-do it)
	w := ren.Drawable.Width()
	h := ren.Drawable.Height()

	tint := colorToFloat32(ren.Color)

	u, v, u2, v2 := ren.Drawable.View()

	if ren.Repeat != NoRepeat {
		u2 = space.Width / (ren.Drawable.Width() * ren.Scale.X)
		w *= u2
		v2 = space.Width / (ren.Drawable.Height() * ren.Scale.Y)
		h *= v2
	}

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

	// Since each sprite in the batch has a different transform, we can't just send the model matrix into
	// the shader and let the GPU take care of it. Instead, we need to multiply the current sprite's model matrix
	// with the position component for each vertex of the current sprite on the CPU, and send the transformed
	// positions to the shader directly.
	modelMatrix := s.makeModelMatrix(ren, space)
	s.multModel(modelMatrix, buffer[:2])
	s.multModel(modelMatrix, buffer[5:7])
	s.multModel(modelMatrix, buffer[10:12])
	s.multModel(modelMatrix, buffer[15:17])
	return changed
}

func (s *basicShader) multModel(m *engo.Matrix, v []float32) {
	tmp := engo.MultiplyMatrixVector(m, v)
	v[0] = tmp[0]
	v[1] = tmp[1]
}

func (s *basicShader) SetCamera(c *CameraSystem) {
	if s.cameraEnabled {
		s.camera = c
		s.viewMatrix.Identity().Translate(-s.camera.x, -s.camera.y).Rotate(s.camera.angle)
	} else {
		scaleX, scaleY := s.projectionMatrix.ScaleComponent()
		s.viewMatrix.Translate(-1/scaleX, 1/scaleY)
	}
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

varying vec4 var_Color;

void main() {
  var_Color = in_Color;

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

void main (void) {
  gl_FragColor = var_Color;
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
		l.projectionMatrix[0] = 1 / (engo.CanvasWidth() / (2 * engo.CanvasScale()))
		l.projectionMatrix[4] = 1 / (-engo.CanvasHeight() / (2 * engo.CanvasScale()))
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

func (l *legacyShader) updateBuffer(ren *RenderComponent, space *SpaceComponent) {
	if len(ren.BufferContent) == 0 {
		ren.BufferContent = make([]float32, l.computeBufferSize(ren.Drawable)) // because we add at most this many elements to it
	}
	if changed := l.generateBufferContent(ren, space, ren.BufferContent); !changed {
		return
	}

	if ren.Buffer == nil {
		ren.Buffer = engo.Gl.CreateBuffer()
	}
	engo.Gl.BindBuffer(engo.Gl.ARRAY_BUFFER, ren.Buffer)
	engo.Gl.BufferData(engo.Gl.ARRAY_BUFFER, ren.BufferContent, engo.Gl.STATIC_DRAW)
}

func (l *legacyShader) computeBufferSize(draw Drawable) int {
	switch shape := draw.(type) {
	case Triangle:
		return 65
	case Rectangle:
		return 90
	case Circle:
		return 1800
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
				b := shape.BorderWidth
				s, c := math.Sincos(math.Atan(2 * h / w))

				pts := [][]float32{
					//Left
					{w / 2, 0},
					{0, h},
					{b, h},
					{b, h},
					{(w / 2) + b*c, b * s},
					{w / 2, 0},
					//Right
					{w / 2, 0},
					{w, h},
					{w - b, h},
					{w - b, h},
					{(w / 2) - b*c, b * s},
					{w / 2, 0},
					//Bottom
					{0, h},
					{w, h},
					{b * c, h - b*s},
					{b * c, h - b*s},
					{w - b*c, h - b*s},
					{w, h},
				}

				for i, p := range pts {
					setBufferValue(buffer, 9+3*i, p[0], &changed)
					setBufferValue(buffer, 10+3*i, p[1], &changed)
					setBufferValue(buffer, 11+3*i, borderTint, &changed)
				}
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
				b := shape.BorderWidth

				pts := [][]float32{
					//Left
					{0, 0},
					{0, h},
					{b, h},
					{b, h},
					{b, b * h / w},
					{0, 0},
					//Right
					{0, 0},
					{w, h},
					{w - b, h},
					{w - b, h},
					{0, b},
					{0, 0},
					//Bottom
					{0, h},
					{w, h},
					{w - b*w/h, h - b},
					{w - b*w/h, h - b},
					{0, h - b},
					{0, h},
				}

				for i, p := range pts {
					setBufferValue(buffer, 9+3*i, p[0], &changed)
					setBufferValue(buffer, 10+3*i, p[1], &changed)
					setBufferValue(buffer, 11+3*i, borderTint, &changed)
				}
			}
		}

	case Circle:
		theta := float32(2.0 * math.Pi / 300.0)
		s, c := math.Sincos(theta)
		x := w / 2
		cx := w / 2
		bx := shape.BorderWidth
		y := float32(0.0)
		cy := h / 2
		by := shape.BorderWidth
		var borderTint float32
		hasBorder := shape.BorderWidth > 0
		if hasBorder {
			borderTint = colorToFloat32(shape.BorderColor)
		}
		for i := 0; i < 300; i++ {
			setBufferValue(buffer, i*3, x+cx-bx/2, &changed)
			setBufferValue(buffer, i*3+1, y+cy-by/2, &changed)
			setBufferValue(buffer, i*3+2, tint, &changed)
			if hasBorder {
				setBufferValue(buffer, i*3+900, x+cx, &changed)
				setBufferValue(buffer, i*3+901, y+cy, &changed)
				setBufferValue(buffer, i*3+902, borderTint, &changed)
			}
			t := x
			bt := bx
			x = c*x - s*y
			bx = c*bx - s*by
			y = s*t + c*y
			by = s*bt + c*by
		}

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

		setBufferValue(buffer, 9, w, &changed)
		setBufferValue(buffer, 10, h, &changed)
		setBufferValue(buffer, 11, tint, &changed)

		//setBufferValue(buffer, 12, 0, &changed)
		setBufferValue(buffer, 13, h, &changed)
		setBufferValue(buffer, 14, tint, &changed)

		//setBufferValue(buffer, 15, 0, &changed)
		//setBufferValue(buffer, 16, 0, &changed)
		setBufferValue(buffer, 17, tint, &changed)

		if shape.BorderWidth > 0 {
			borderTint := colorToFloat32(shape.BorderColor)
			b := shape.BorderWidth
			pts := [][]float32{
				//Top
				{0, 0},
				{w, 0},
				{w, b},
				{w, b},
				{0, b},
				{0, 0},
				//Right
				{w - b, b},
				{w, b},
				{w, h - b},
				{w, h - b},
				{w - b, h - b},
				{w - b, b},
				//Bottom
				{w, h - b},
				{w, h},
				{0, h},
				{0, h},
				{0, h - b},
				{w, h - b},
				//Left
				{0, b},
				{b, b},
				{b, h - b},
				{b, h - b},
				{0, h - b},
				{0, b},
			}

			for i, p := range pts {
				setBufferValue(buffer, 18+3*i, p[0], &changed)
				setBufferValue(buffer, 19+3*i, p[1], &changed)
				setBufferValue(buffer, 20+3*i, borderTint, &changed)
			}
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
	if l.lastBuffer != ren.Buffer || ren.Buffer == nil {
		l.updateBuffer(ren, space)

		engo.Gl.BindBuffer(engo.Gl.ARRAY_BUFFER, ren.Buffer)
		engo.Gl.VertexAttribPointer(l.inPosition, 2, engo.Gl.FLOAT, false, 12, 0)
		engo.Gl.VertexAttribPointer(l.inColor, 4, engo.Gl.UNSIGNED_BYTE, true, 12, 8)

		l.lastBuffer = ren.Buffer
	}

	if space.Rotation != 0 {
		sin, cos := math.Sincos(space.Rotation * math.Pi / 180)

		l.modelMatrix[0] = ren.Scale.X * engo.GetGlobalScale().X * cos
		l.modelMatrix[1] = ren.Scale.X * engo.GetGlobalScale().X * sin
		l.modelMatrix[3] = ren.Scale.Y * engo.GetGlobalScale().Y * -sin
		l.modelMatrix[4] = ren.Scale.Y * engo.GetGlobalScale().Y * cos
	} else {
		l.modelMatrix[0] = ren.Scale.X * engo.GetGlobalScale().X
		l.modelMatrix[1] = 0
		l.modelMatrix[3] = 0
		l.modelMatrix[4] = ren.Scale.Y * engo.GetGlobalScale().Y
	}

	l.modelMatrix[6] = space.Position.X * engo.GetGlobalScale().X
	l.modelMatrix[7] = space.Position.Y * engo.GetGlobalScale().Y

	engo.Gl.UniformMatrix3fv(l.matrixModel, false, l.modelMatrix)

	switch shape := ren.Drawable.(type) {
	case Triangle:
		num := 3
		if shape.BorderWidth > 0 {
			num = 21
		}
		engo.Gl.DrawArrays(engo.Gl.TRIANGLES, 0, num)
	case Rectangle:
		num := 6
		if shape.BorderWidth > 0 {
			num = 30
		}
		engo.Gl.DrawArrays(engo.Gl.TRIANGLES, 0, num)
	case Circle:
		// Circle stuff!
		if shape.BorderWidth > 0 {
			engo.Gl.DrawArrays(engo.Gl.TRIANGLE_FAN, 300, 300)
		}
		engo.Gl.DrawArrays(engo.Gl.TRIANGLE_FAN, 0, 300)
	case ComplexTriangles:
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

func (l *legacyShader) SetCamera(c *CameraSystem) {
	if l.cameraEnabled {
		l.camera = c
	}
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
		l.projectionMatrix[0] = 1 / (engo.CanvasWidth() / (2 * engo.CanvasScale()))
		l.projectionMatrix[4] = 1 / (-engo.CanvasHeight() / (2 * engo.CanvasScale()))
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

	if len(ren.BufferContent) < 20*len(txt.Text) {
		ren.BufferContent = make([]float32, 20*len(txt.Text)) // TODO: update this to actual value?
	}
	if changed := l.generateBufferContent(ren, space, ren.BufferContent); !changed {
		return
	}

	if ren.Buffer == nil {
		ren.Buffer = engo.Gl.CreateBuffer()
	}
	engo.Gl.BindBuffer(engo.Gl.ARRAY_BUFFER, ren.Buffer)
	engo.Gl.BufferData(engo.Gl.ARRAY_BUFFER, ren.BufferContent, engo.Gl.STATIC_DRAW)
}

func (l *textShader) generateBufferContent(ren *RenderComponent, space *SpaceComponent, buffer []float32) bool {
	var changed bool

	tint := colorToFloat32(ren.Color)
	txt, ok := ren.Drawable.(Text)
	if !ok {
		unsupportedType(ren.Drawable)
		return false
	}

	atlas, ok := atlasCache[*txt.Font]
	if !ok {
		// Generate texture first
		atlas = txt.Font.generateFontAtlas(UnicodeCap)
		atlasCache[*txt.Font] = atlas
	}

	var currentX float32
	var currentY float32

	var modifier float32 = 1
	if txt.RightToLeft {
		modifier = -1
	}

	letterSpace := float32(txt.Font.Size) * txt.LetterSpacing
	lineSpace := txt.LineSpacing * atlas.Height['X']

	for index, char := range txt.Text {
		// TODO: this might not work for all characters
		switch {
		case char == '\n':
			currentX = 0
			currentY += atlas.Height['X'] + lineSpace
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
		setBufferValue(buffer, 5+offset, currentX+atlas.Width[char]+letterSpace, &changed)
		setBufferValue(buffer, 6+offset, currentY, &changed)
		setBufferValue(buffer, 7+offset, (atlas.XLocation[char]+atlas.Width[char])/atlas.TotalWidth, &changed)
		setBufferValue(buffer, 8+offset, atlas.YLocation[char]/atlas.TotalHeight, &changed)
		setBufferValue(buffer, 9+offset, tint, &changed)

		// These five are at 1, 1:
		setBufferValue(buffer, 10+offset, currentX+atlas.Width[char]+letterSpace, &changed)
		setBufferValue(buffer, 11+offset, currentY+atlas.Height[char]+lineSpace, &changed)
		setBufferValue(buffer, 12+offset, (atlas.XLocation[char]+atlas.Width[char])/atlas.TotalWidth, &changed)
		setBufferValue(buffer, 13+offset, (atlas.YLocation[char]+atlas.Height[char])/atlas.TotalHeight, &changed)
		setBufferValue(buffer, 14+offset, tint, &changed)

		// These five are at 0, 1:
		setBufferValue(buffer, 15+offset, currentX, &changed)
		setBufferValue(buffer, 16+offset, currentY+atlas.Height[char]+lineSpace, &changed)
		setBufferValue(buffer, 17+offset, atlas.XLocation[char]/atlas.TotalWidth, &changed)
		setBufferValue(buffer, 18+offset, (atlas.YLocation[char]+atlas.Height[char])/atlas.TotalHeight, &changed)
		setBufferValue(buffer, 19+offset, tint, &changed)

		currentX += modifier * (atlas.Width[char] + letterSpace)
	}

	return changed
}

func (l *textShader) Draw(ren *RenderComponent, space *SpaceComponent) {
	if l.lastBuffer != ren.Buffer || ren.Buffer == nil {
		l.updateBuffer(ren, space)

		engo.Gl.BindBuffer(engo.Gl.ARRAY_BUFFER, ren.Buffer)
		engo.Gl.VertexAttribPointer(l.inPosition, 2, engo.Gl.FLOAT, false, 20, 0)
		engo.Gl.VertexAttribPointer(l.inTexCoords, 2, engo.Gl.FLOAT, false, 20, 8)
		engo.Gl.VertexAttribPointer(l.inColor, 4, engo.Gl.UNSIGNED_BYTE, true, 20, 16)

		l.lastBuffer = ren.Buffer
	}

	txt, ok := ren.Drawable.(Text)
	if !ok {
		unsupportedType(ren.Drawable)
	}

	atlas, ok := atlasCache[*txt.Font]
	if !ok {
		// Generate texture first
		atlas = txt.Font.generateFontAtlas(UnicodeCap)
		atlasCache[*txt.Font] = atlas
	}

	if atlas.Texture != l.lastTexture {
		engo.Gl.BindTexture(engo.Gl.TEXTURE_2D, atlas.Texture)
		l.lastTexture = atlas.Texture
	}

	engo.Gl.TexParameteri(engo.Gl.TEXTURE_2D, engo.Gl.TEXTURE_WRAP_S, engo.Gl.CLAMP_TO_EDGE)
	engo.Gl.TexParameteri(engo.Gl.TEXTURE_2D, engo.Gl.TEXTURE_WRAP_T, engo.Gl.CLAMP_TO_EDGE)

	if space.Rotation != 0 {
		sin, cos := math.Sincos(space.Rotation * math.Pi / 180)

		l.modelMatrix[0] = ren.Scale.X * engo.GetGlobalScale().X * cos
		l.modelMatrix[1] = ren.Scale.X * engo.GetGlobalScale().X * sin
		l.modelMatrix[3] = ren.Scale.Y * engo.GetGlobalScale().Y * -sin
		l.modelMatrix[4] = ren.Scale.Y * engo.GetGlobalScale().Y * cos
	} else {
		l.modelMatrix[0] = ren.Scale.X * engo.GetGlobalScale().X
		l.modelMatrix[1] = 0
		l.modelMatrix[3] = 0
		l.modelMatrix[4] = ren.Scale.Y * engo.GetGlobalScale().Y
	}

	l.modelMatrix[6] = space.Position.X * engo.GetGlobalScale().X
	l.modelMatrix[7] = space.Position.Y * engo.GetGlobalScale().Y

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

func (l *textShader) SetCamera(c *CameraSystem) {
	if l.cameraEnabled {
		l.camera = c
	}
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
	// DefaultShader is the shader picked when no other shader is used.
	DefaultShader = &basicShader{cameraEnabled: true}
	// HUDShader is the shader used for HUD elements.
	HUDShader = &basicShader{cameraEnabled: false}
	// LegacyShader is the shader used for drawing shapes.
	LegacyShader = &legacyShader{cameraEnabled: true}
	// LegacyHUDShader is the shader used for drawing shapes on the HUD.
	LegacyHUDShader = &legacyShader{cameraEnabled: false}
	// TextShader is the shader used to draw fonts from a FontAtlas
	TextShader = &textShader{cameraEnabled: true}
	// TextHUDShader is the shader used to draw fonts from a FontAtlas on the HUD.
	TextHUDShader = &textShader{cameraEnabled: false}

	BlendmapShader = &blendmapShader{cameraEnabled: true}
	shadersSet     bool
	atlasCache     = make(map[Font]FontAtlas)
	shaders        = []Shader{
		DefaultShader,
		HUDShader,
		LegacyShader,
		LegacyHUDShader,
		TextShader,
		TextHUDShader,
		BlendmapShader,
	}
)

// AddShader adds a shader to the list of shaders for initalization. They should
// be added before the Rendersystem is added, such as in the scene's Preload.
func AddShader(s Shader) {
	shaders = append(shaders, s)
}

var shaderInitMutex sync.Mutex

func initShaders(w *ecs.World) error {
	shaderInitMutex.Lock()
	defer shaderInitMutex.Unlock()

	if !shadersSet {
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

func newCamera(w *ecs.World) {
	shaderInitMutex.Lock()
	defer shaderInitMutex.Unlock()
	var cam *CameraSystem
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *CameraSystem:
			cam = sys
		}
	}
	if cam == nil {
		log.Println("Camera system was not found when changing scene!")
		return
	}
	for _, shader := range shaders {
		shader.SetCamera(cam)
	}
}

// VertexShaderCompilationError is returned whenever the `LoadShader` method was unable to compile your Vertex-shader (GLSL)
type VertexShaderCompilationError struct {
	OpenGLError string
}

// Error implements the error interface.
func (v VertexShaderCompilationError) Error() string {
	return fmt.Sprintf("an error occurred compiling the vertex shader: %s", strings.Trim(v.OpenGLError, "\r\n"))
}

// FragmentShaderCompilationError is returned whenever the `LoadShader` method was unable to compile your Fragment-shader (GLSL)
type FragmentShaderCompilationError struct {
	OpenGLError string
}

// Error implements the error interface.
func (f FragmentShaderCompilationError) Error() string {
	return fmt.Sprintf("an error occurred compiling the fragment shader: %s", strings.Trim(f.OpenGLError, "\r\n"))
}
