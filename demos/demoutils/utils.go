package demoutils

import (
	"image"
	"image/color"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/core"
)

type Background struct {
	ecs.BasicEntity
	core.RenderComponent
	core.SpaceComponent
}

// NewBackground creates a background of colored tiles - might not be the most efficient way to do this
// It gets added to the world as well, so we won't return anything.
func NewBackground(world *ecs.World, width, height int, colorA, colorB color.Color) *Background {
	rect := image.Rect(0, 0, width, height)

	img := image.NewNRGBA(rect)
	for i := rect.Min.X; i < rect.Max.X; i++ {
		for j := rect.Min.Y; j < rect.Max.Y; j++ {
			if i%40 > 20 {
				if j%40 > 20 {
					img.Set(i, j, colorA)
				} else {
					img.Set(i, j, colorB)
				}
			} else {
				if j%40 > 20 {
					img.Set(i, j, colorB)
				} else {
					img.Set(i, j, colorA)
				}
			}
		}
	}

	bgTexture := core.NewImageObject(img)

	bg := &Background{BasicEntity: ecs.NewBasic()}
	bg.RenderComponent = core.RenderComponent{Drawable: core.NewTextureSingle(bgTexture)}
	bg.SpaceComponent = core.SpaceComponent{
		Position: engo.Point{0, 0},
		Width:    float32(width),
		Height:   float32(height),
	}

	for _, system := range world.Systems() {
		switch sys := system.(type) {
		case *core.RenderSystem:
			sys.Add(&bg.BasicEntity, &bg.RenderComponent, &bg.SpaceComponent)
		}
	}

	return bg
}
