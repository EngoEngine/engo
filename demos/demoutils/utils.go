//+build demo

package demoutils

import (
	"image"
	"image/color"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
)

type Background struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
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

	bgTexture := common.NewImageObject(img)

	bg := &Background{BasicEntity: ecs.NewBasic()}
	bg.RenderComponent = common.RenderComponent{Drawable: common.NewTextureSingle(bgTexture)}
	bg.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{0, 0},
		Width:    float32(width),
		Height:   float32(height),
	}

	for _, system := range world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&bg.BasicEntity, &bg.RenderComponent, &bg.SpaceComponent)
		}
	}

	return bg
}
