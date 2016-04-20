package demoutils

import (
	"image"
	"image/color"

	"engo.io/ecs"
	"engo.io/engo"
)

type Background struct {
	ecs.BasicEntity
	engo.RenderComponent
	engo.SpaceComponent
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

	bgTexture := engo.NewImageObject(img)

	bg := &Background{BasicEntity: ecs.NewBasic()}
	bg.RenderComponent = engo.NewRenderComponent(engo.NewTexture(bgTexture), engo.Point{1, 1}, "Background")
	bg.SpaceComponent = engo.SpaceComponent{engo.Point{0, 0}, float32(width), float32(height)}

	for _, system := range world.Systems() {
		switch sys := system.(type) {
		case *engo.RenderSystem:
			sys.Add(&bg.BasicEntity, &bg.RenderComponent, &bg.SpaceComponent)
		}
	}

	return bg
}
