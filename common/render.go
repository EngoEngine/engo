package common

import (
	"fmt"
	"image/color"
	"sort"

	"sync"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/gl"
)

const (
	// RenderSystemPriority is the priority of the RenderSystem
	RenderSystemPriority = -1000
)

type renderChangeMessage struct{}

func (renderChangeMessage) Type() string {
	return "renderChangeMessage"
}

// Drawable is that which can be rendered to OpenGL.
type Drawable interface {
	Texture() *gl.Texture
	Width() float32
	Height() float32
	View() (float32, float32, float32, float32)
	Close()
}

// TextureRepeating is the method used to repeat a texture in OpenGL.
type TextureRepeating uint8

const (
	// NoRepeat does not repeat the texture.
	NoRepeat TextureRepeating = iota
	// ClampToEdge stretches the texture to the edge of the viewport.
	ClampToEdge
	// ClampToBorder stretches the texture to the border of the viewpport.
	ClampToBorder
	// Repeat repeats the texture until the border of the viewport.
	Repeat
	// MirroredRepeat repeats a mirror image of the texture until the border of the viewport.
	MirroredRepeat
)

// ZoomFilter is a filter used when zooming in or out of a texture.
type ZoomFilter uint8

const (
	// FilterNearest is a simple nearest neighbor algorithm
	FilterNearest ZoomFilter = iota
	// FilterLinear is a bilinear interpolation algorithm
	FilterLinear
)

// RenderComponent is the component needed to render an entity.
type RenderComponent struct {
	// Hidden is used to prevent drawing by OpenGL
	Hidden bool
	// Scale is the scale at which to render, in the X and Y axis. Not defining Scale, will default to engo.Point{1, 1}
	Scale engo.Point
	// Color defines how much of the color-components of the texture get used
	Color color.Color
	// Drawable refers to the Texture that should be drawn
	Drawable Drawable
	// Repeat defines how to repeat the Texture if the SpaceComponent of the entity
	// is larger than the texture itself, after applying scale. Defaults to NoRepeat
	// which allows the texture to draw entirely without regard to th SpaceComponent
	// Do not set to anything other than NoRepeat for textures in a sprite sheet.
	// This does not yet work with sprite sheets.
	Repeat TextureRepeating

	magFilter, minFilter               ZoomFilter
	magFilterChanged, minFilterChanged bool

	shader Shader
	zIndex float32

	buffer        *gl.Buffer
	bufferContent []float32
}

// SetShader sets the shader used by the RenderComponent.
func (r *RenderComponent) SetShader(s Shader) {
	r.shader = s
	engo.Mailbox.Dispatch(&renderChangeMessage{})
}

// SetZIndex sets the order that the RenderComponent is drawn to the screen. Higher z-indices are drawn on top of
// lower ones if they overlap.
func (r *RenderComponent) SetZIndex(index float32) {
	r.zIndex = index
	engo.Mailbox.Dispatch(&renderChangeMessage{})
}

// SetMinFilter sets the ZoomFilter used for minimizing the RenderComponent
func (r *RenderComponent) SetMinFilter(z ZoomFilter) {
	r.minFilter = z
	r.minFilterChanged = true
}

// SetMagFilter sets the ZoomFilter used for magnifying the RenderComponent
func (r *RenderComponent) SetMagFilter(z ZoomFilter) {
	r.magFilter = z
	r.magFilterChanged = true
}

type renderEntity struct {
	*ecs.BasicEntity
	*RenderComponent
	*SpaceComponent
}

type renderEntityList []renderEntity

func (r renderEntityList) Len() int {
	return len(r)
}

func (r renderEntityList) Less(i, j int) bool {
	// Sort by shader-pointer if they have the same zIndex
	if r[i].RenderComponent.zIndex == r[j].RenderComponent.zIndex {
		// TODO: optimize this for performance
		return fmt.Sprintf("%p", &r[i].RenderComponent.shader) < fmt.Sprintf("%p", &r[j].RenderComponent.shader)
	}

	return r[i].RenderComponent.zIndex < r[j].RenderComponent.zIndex
}

func (r renderEntityList) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

// RenderSystem is the system that draws entities on the OpenGL surface. It requires
// a CameraSystem to work. If a CameraSystem is not in the World when you add RenderSystem
// one is automatically added to the world.
type RenderSystem struct {
	entities renderEntityList
	world    *ecs.World

	sortingNeeded bool
	currentShader Shader
}

// Priority implements the ecs.Prioritizer interface.
func (*RenderSystem) Priority() int { return RenderSystemPriority }

// New initializes the RenderSystem
func (rs *RenderSystem) New(w *ecs.World) {
	rs.world = w

	addCameraSystemOnce(w)

	if !engo.Headless() {
		initShaders(w)
		engo.Gl.Enable(engo.Gl.MULTISAMPLE)
	}

	engo.Mailbox.Listen("renderChangeMessage", func(engo.Message) {
		rs.sortingNeeded = true
	})
}

var cameraInitMutex sync.Mutex

func addCameraSystemOnce(w *ecs.World) {
	cameraInitMutex.Lock()
	defer cameraInitMutex.Unlock()

	camSystemAlreadyAdded := false
	for _, system := range w.Systems() {
		switch system.(type) {
		case *CameraSystem:
			camSystemAlreadyAdded = true
		}
	}
	if !camSystemAlreadyAdded {
		w.AddSystem(&CameraSystem{})
	}
}

// Add adds an entity to the RenderSystem. The entity needs a basic, render, and space component to be added to the system.
func (rs *RenderSystem) Add(basic *ecs.BasicEntity, render *RenderComponent, space *SpaceComponent) {
	// Do nothing if entity already exists
	if rs.EntityExists(basic) >= 0 {
		return
	}

	// Setting default shader
	if render.shader == nil {
		switch render.Drawable.(type) {
		case Triangle:
			render.shader = LegacyShader
		case Circle:
			render.shader = LegacyShader
		case Rectangle:
			render.shader = LegacyShader
		case ComplexTriangles:
			render.shader = LegacyShader
		case Text:
			render.shader = TextShader
		default:
			render.shader = DefaultShader
		}
	}

	// This is to prevent users from using the wrong one
	if render.shader == HUDShader {
		switch render.Drawable.(type) {
		case Triangle:
			render.shader = LegacyHUDShader
		case Circle:
			render.shader = LegacyHUDShader
		case Rectangle:
			render.shader = LegacyHUDShader
		case ComplexTriangles:
			render.shader = LegacyHUDShader
		case Text:
			render.shader = TextHUDShader
		default:
			render.shader = HUDShader
		}
	}

	rs.entities = append(rs.entities, renderEntity{basic, render, space})
	rs.sortingNeeded = true
}

// EntityExists looks if the entity is already into the System's entities. It will return the index >= 0 of the object into de rs.entities or -1 if it could not be found.
func (rs *RenderSystem) EntityExists(basic *ecs.BasicEntity) int {
	for index, entity := range rs.entities {
		if entity.ID() == basic.ID() {
			return index
		}
	}

	return -1
}

// AddByInterface adds any Renderable to the render system. Any Entity containing a BasicEntity,RenderComponent, and SpaceComponent anonymously does this automatically
func (rs *RenderSystem) AddByInterface(i ecs.Identifier) {
	o, _ := i.(Renderable)
	rs.Add(o.GetBasicEntity(), o.GetRenderComponent(), o.GetSpaceComponent())
}

// Remove removes an entity from the RenderSystem
func (rs *RenderSystem) Remove(basic ecs.BasicEntity) {
	var delete = rs.EntityExists(&basic)
	if delete >= 0 {
		rs.entities = append(rs.entities[:delete], rs.entities[delete+1:]...)
		rs.sortingNeeded = true
	}
}

// Update draws the entities in the RenderSystem to the OpenGL Surface.
func (rs *RenderSystem) Update(dt float32) {
	if engo.Headless() {
		return
	}

	if rs.sortingNeeded {
		sort.Sort(rs.entities)
		rs.sortingNeeded = false
	}

	engo.Gl.Clear(engo.Gl.COLOR_BUFFER_BIT)

	// TODO: it's linear for now, but that might very well be a bad idea
	for _, e := range rs.entities {
		if e.RenderComponent.Hidden {
			continue // with other entities
		}

		// Retrieve a shader, may be the default one -- then use it if we aren't already using it
		shader := e.RenderComponent.shader

		// Change Shader if we have to
		if shader != rs.currentShader {
			if rs.currentShader != nil {
				rs.currentShader.Post()
			}
			shader.Pre()
			rs.currentShader = shader
		}

		// Setting default scale to 1
		if e.RenderComponent.Scale.X == 0 && e.RenderComponent.Scale.Y == 0 {
			e.RenderComponent.Scale = engo.Point{X: 1, Y: 1}
		}

		// Setting default to white
		if e.RenderComponent.Color == nil {
			e.RenderComponent.Color = color.White
		}

		rs.currentShader.Draw(e.RenderComponent, e.SpaceComponent)
	}

	if rs.currentShader != nil {
		rs.currentShader.Post()
		rs.currentShader = nil
	}
}

// SetBackground sets the OpenGL ClearColor to the provided color.
func SetBackground(c color.Color) {
	if !engo.Headless() {
		r, g, b, a := c.RGBA()

		engo.Gl.ClearColor(float32(r)/0xffff, float32(g)/0xffff, float32(b)/0xffff, float32(a)/0xffff)
	}
}
