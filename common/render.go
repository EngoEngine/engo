package common

import (
	"fmt"
	"image/color"
	"sort"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/gl"
)

const (
	RenderSystemPriority = -1000
)

type renderChangeMessage struct{}

func (renderChangeMessage) Type() string {
	return "renderChangeMessage"
}

type Drawable interface {
	Texture() *gl.Texture
	Width() float32
	Height() float32
	View() (float32, float32, float32, float32)
	Close()
}

type TextureRepeating uint8

const (
	ClampToEdge TextureRepeating = iota
	ClampToBorder
	Repeat
	MirroredRepeat
)

type RenderComponent struct {
	// Hidden is used to prevent drawing by OpenGL
	Hidden bool
	// Scale is the scale at which to render, in the X and Y axis. Not defining Scale, will default to engo.Point{1, 1}
	Scale engo.Point
	// Color defines how much of the color-components of the texture get used
	Color color.Color
	// Drawable refers to the Texture that should be drawn
	Drawable Drawable
	// Repeat defines how to repeat the Texture if the viewport of the texture is larger than the texture itself
	Repeat TextureRepeating

	shader Shader
	zIndex float32

	buffer        *gl.Buffer
	bufferContent []float32
}

func (r *RenderComponent) SetShader(s Shader) {
	r.shader = s
	engo.Mailbox.Dispatch(&renderChangeMessage{})
}

func (r *RenderComponent) SetZIndex(index float32) {
	r.zIndex = index
	engo.Mailbox.Dispatch(&renderChangeMessage{})
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
		return fmt.Sprintf("%p", r[i].RenderComponent.shader) < fmt.Sprintf("%p", r[j].RenderComponent.shader)
	}

	return r[i].RenderComponent.zIndex < r[j].RenderComponent.zIndex
}

func (r renderEntityList) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

type RenderSystem struct {
	entities renderEntityList
	world    *ecs.World

	sortingNeeded bool
	currentShader Shader
}

func (*RenderSystem) Priority() int { return RenderSystemPriority }

func (rs *RenderSystem) New(w *ecs.World) {
	rs.world = w

	w.AddSystem(&CameraSystem{})

	if !engo.Headless() {
		initShaders(w)
		engo.Gl.Enable(engo.Gl.MULTISAMPLE)
	}

	engo.Mailbox.Listen("renderChangeMessage", func(engo.Message) {
		rs.sortingNeeded = true
	})
}

func (rs *RenderSystem) Add(basic *ecs.BasicEntity, render *RenderComponent, space *SpaceComponent) {
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

func (rs *RenderSystem) Remove(basic ecs.BasicEntity) {
	var delete int = -1
	for index, entity := range rs.entities {
		if entity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		rs.entities = append(rs.entities[:delete], rs.entities[delete+1:]...)
		rs.sortingNeeded = true
	}
}

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
			e.RenderComponent.Scale = engo.Point{1, 1}
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

func SetBackground(c color.Color) {
	if !engo.Headless() {
		r, g, b, a := c.RGBA()

		engo.Gl.ClearColor(float32(r)/0xffff, float32(g)/0xffff, float32(b)/0xffff, float32(a)/0xffff)
	}
}
