package engo

import (
	"fmt"
	"image/color"
	"sort"

	"engo.io/ecs"
	"engo.io/gl"
	"github.com/luxengine/math"
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
}

type RenderComponent struct {
	// Hidden is used to prevent drawing by OpenGL
	Hidden bool

	// transparency is the level of transparency that is used to draw the texture
	transparency float32

	Scale  Point
	Color  color.Color
	shader Shader
	zIndex float32

	drawable      Drawable
	buffer        *gl.Buffer
	bufferContent []float32
}

func NewRenderComponent(d Drawable, scale Point) RenderComponent {
	rc := RenderComponent{
		transparency: 1,
		Color:        color.White,
		Scale:        scale,
	}
	rc.SetDrawable(d)

	return rc
}

func (r *RenderComponent) SetTransparency(t float32) {
	r.transparency = t
	// regen buffer content as we just changed transparency for the whole
	// sprite
	r.generateBufferContent()
}

func (r *RenderComponent) SetDrawable(d Drawable) {
	r.drawable = d
	r.preloadTexture()
}

func (r *RenderComponent) Drawable() Drawable {
	return r.drawable
}

func (r *RenderComponent) SetShader(s Shader) {
	r.shader = s
	Mailbox.Dispatch(&renderChangeMessage{})
}

func (r *RenderComponent) SetZIndex(index float32) {
	r.zIndex = index
	Mailbox.Dispatch(&renderChangeMessage{})
}

// Init is called to initialize the RenderElement
func (ren *RenderComponent) preloadTexture() {
	if ren.drawable == nil || headless {
		return
	}

	ren.bufferContent = ren.generateBufferContent()

	ren.buffer = Gl.CreateBuffer()
	Gl.BindBuffer(Gl.ARRAY_BUFFER, ren.buffer)
	Gl.BufferData(Gl.ARRAY_BUFFER, ren.bufferContent, Gl.STATIC_DRAW)
}

// generateBufferContent computes information about the 4 vertices needed to draw the texture, which should
// be stored in the buffer
func (ren *RenderComponent) generateBufferContent() []float32 {
	w := ren.drawable.Width()
	h := ren.drawable.Height()

	colorR, colorG, colorB, _ := ren.Color.RGBA()

	red := colorR
	green := colorG << 8
	blue := colorB << 16
	alpha := uint32(ren.transparency*255.0) << 24

	tint := math.Float32frombits((alpha | blue | green | red) & 0xfeffffff)

	u, v, u2, v2 := ren.drawable.View()

	return []float32{0, 0, u, v, tint, w, 0, u2, v, tint, w, h, u2, v2, tint, 0, h, u, v2, tint}
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

	if !headless {
		initShaders()
	}

	Mailbox.Listen("renderChangeMessage", func(Message) {
		rs.sortingNeeded = true
	})
}

func (rs *RenderSystem) Add(basic *ecs.BasicEntity, render *RenderComponent, space *SpaceComponent) {
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
	if headless {
		return
	}

	if rs.sortingNeeded {
		sort.Sort(rs.entities)
		rs.sortingNeeded = false
	}

	Gl.Clear(Gl.COLOR_BUFFER_BIT)

	// TODO: it's linear for now, but that might very well be a bad idea
	for _, e := range rs.entities {
		if e.RenderComponent.Hidden {
			continue // with other entities
		}

		// Retrieve a shader, may be the default one -- then use it if we aren't already using it
		shader := e.RenderComponent.shader
		if shader == nil {
			shader = DefaultShader
		}

		// Change Shader if we have to
		if shader != rs.currentShader {
			if rs.currentShader != nil {
				rs.currentShader.Post()
			}
			shader.Pre()
			rs.currentShader = shader
		}

		rs.currentShader.Draw(e.RenderComponent.drawable.Texture(), e.RenderComponent.buffer,
			e.SpaceComponent.Position.X, e.SpaceComponent.Position.Y,
			e.RenderComponent.Scale.X, e.RenderComponent.Scale.Y,
			e.SpaceComponent.Rotation)
	}

	if rs.currentShader != nil {
		rs.currentShader.Post()
		rs.currentShader = nil
	}
}
