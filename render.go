package engo

import (
	"fmt"
	"image/color"
	"math"
	"sort"

	"engo.io/ecs"
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

func NewRenderComponent(d Drawable, scale Point) *RenderComponent {
	rc := &RenderComponent{
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

func (*RenderComponent) Type() string {
	return "RenderComponent"
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

	// TODO: ask why this doesn't work
	// ren.bufferContent = make([]float32, 0)
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

type renderEntityList []*ecs.Entity

func (r renderEntityList) Len() int {
	return len(r)
}

func (r renderEntityList) Less(i, j int) bool {
	var (
		rc1 *RenderComponent
		rc2 *RenderComponent
		ok  bool
	)
	if rc1, ok = r[i].ComponentFast(rc1).(*RenderComponent); !ok {
		return false // those without render component go last
	}
	if rc2, ok = r[i].ComponentFast(rc1).(*RenderComponent); !ok {
		return true // those without render component go last
	}

	// Sort by shader-pointer if they have the same zIndex
	if rc1.zIndex == rc2.zIndex {
		// TODO: optimize this for performance
		return fmt.Sprintf("%p", rc1.shader) < fmt.Sprintf("%p", rc2.shader)
	}

	return rc1.zIndex < rc2.zIndex
}

func (r renderEntityList) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

type RenderSystem struct {
	renders renderEntityList
	world   *ecs.World

	sortingNeeded bool
	currentShader Shader
}

func (rs *RenderSystem) New(w *ecs.World) {
	rs.world = w

	if !headless {
		initShaders()
	}

	Mailbox.Listen("renderChangeMessage", func(Message) {
		rs.sortingNeeded = true
	})
}

func (rs *RenderSystem) AddEntity(e *ecs.Entity) {
	rs.renders = append(rs.renders, e)
	rs.sortingNeeded = true
}

func (rs *RenderSystem) RemoveEntity(e *ecs.Entity) {
	var removeIndex int = -1
	for index, entity := range rs.renders {
		if entity.ID() == e.ID() {
			removeIndex = index
			break
		}
	}
	if removeIndex >= 0 {
		rs.renders = append(rs.renders[:removeIndex], rs.renders[removeIndex+1:]...) // TODO: test for edge cases
		rs.sortingNeeded = true
	}
}

func (rs *RenderSystem) Update(dt float32) {
	if headless {
		return
	}

	if rs.sortingNeeded {
		sort.Sort(rs.renders)
		rs.sortingNeeded = false
	}

	Gl.Clear(Gl.COLOR_BUFFER_BIT)

	// TODO: it's linear for now, but that might very well be a bad idea
	for _, entity := range rs.renders {
		var (
			render *RenderComponent
			space  *SpaceComponent
			ok     bool
		)

		if render, ok = entity.ComponentFast(render).(*RenderComponent); !ok {
			continue // with other entities
		}

		if render.Hidden {
			continue // with other entities
		}

		if space, ok = entity.ComponentFast(space).(*SpaceComponent); !ok {
			continue // with other entities
		}

		// Retrieve a shader, may be the default one -- then use it if we aren't already using it
		shader := render.shader
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

		rs.currentShader.Draw(render.drawable.Texture(), render.buffer, space.Position.X, space.Position.Y, render.Scale.X, render.Scale.Y, space.Rotation)
	}

	if rs.currentShader != nil {
		rs.currentShader.Post()
		rs.currentShader = nil
	}
}

func (*RenderSystem) Type() string {
	return "RenderSystem"
}

func (*RenderSystem) Priority() int {
	return RenderSystemPriority
}
