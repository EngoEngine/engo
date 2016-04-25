package main

import (
	"image/color"

	"engo.io/ecs"
	"engo.io/engo"
)

type DefaultScene struct{}

type Guy struct {
	ecs.BasicEntity
	engo.MouseComponent
	engo.RenderComponent
	engo.SpaceComponent
}

<<<<<<< HEAD
// generateBackground creates a background of green tiles - might not be the most efficient way to do this
func generateBackground() *engo.RenderComponent {
	rect := image.Rect(0, 0, int(boxWidth), int(boxHeight))
	img := image.NewNRGBA(rect)
	c1 := color.RGBA{102, 153, 0, 255}
	for i := rect.Min.X; i < rect.Max.X; i++ {
		for j := rect.Min.Y; j < rect.Max.Y; j++ {
			img.Set(i, j, c1)
		}
	}
	bgTexture := engo.NewImageObject(img)
	fieldRender := engo.NewRenderComponent(engo.NewTexture(bgTexture), engo.Point{1, 1})
	return fieldRender
=======
func (*DefaultScene) Preload() {
	// Load all files from the data directory. `false` means: do not do it recursively.
	engo.Files.AddFromDir("data", false)
>>>>>>> 28393c45ef7ce198babe3c6854931398faaba25c
}

func (*DefaultScene) Setup(w *ecs.World) {
	engo.SetBackground(color.White)

	w.AddSystem(&engo.MouseSystem{})
	w.AddSystem(&engo.RenderSystem{})
	w.AddSystem(&ControlSystem{})
	w.AddSystem(&engo.MouseZoomer{-0.125})

	// Retrieve a texture
	texture := engo.Files.Image("icon.png")

<<<<<<< HEAD
func (*GameWorld) Hide()        {}
func (*GameWorld) Show()        {}
func (*GameWorld) Exit()        {}
func (*GameWorld) Type() string { return "GameWorld" }
=======
	// Create an entity
	guy := Guy{BasicEntity: ecs.NewBasic()}
>>>>>>> 28393c45ef7ce198babe3c6854931398faaba25c

	// Initialize the components, set scale to 8x
	guy.RenderComponent = engo.NewRenderComponent(texture, engo.Point{8, 8})
	guy.SpaceComponent = engo.SpaceComponent{
		Position: engo.Point{0, 0},
		Width:    texture.Width() * guy.RenderComponent.Scale.X,
		Height:   texture.Height() * guy.RenderComponent.Scale.Y,
	}
	// guy.MouseComponent doesn't have to be set, because its default values will do

	// Add our guy to appropriate systems
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *engo.RenderSystem:
			sys.Add(&guy.BasicEntity, &guy.RenderComponent, &guy.SpaceComponent)
		case *engo.MouseSystem:
			sys.Add(&guy.BasicEntity, &guy.MouseComponent, &guy.SpaceComponent, &guy.RenderComponent)
		case *ControlSystem:
			sys.Add(&guy.BasicEntity, &guy.MouseComponent)
		}
	}
}

<<<<<<< HEAD
	entity.AddComponent(generateBackground())
	entity.AddComponent(&engo.MouseComponent{})
	entity.AddComponent(&engo.SpaceComponent{
		Position: engo.Point{0, 0},
		Width:    boxWidth,
		Height:   boxHeight,
	})
=======
func (*DefaultScene) Type() string { return "GameWorld" }
>>>>>>> 28393c45ef7ce198babe3c6854931398faaba25c

type controlEntity struct {
	*ecs.BasicEntity
	*engo.MouseComponent
}

type ControlSystem struct {
	entities []controlEntity
}

func (c *ControlSystem) Add(basic *ecs.BasicEntity, mouse *engo.MouseComponent) {
	c.entities = append(c.entities, controlEntity{basic, mouse})
}

func (c *ControlSystem) Remove(basic ecs.BasicEntity) {
	delete := -1
	for index, e := range c.entities {
		if e.BasicEntity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		c.entities = append(c.entities[:delete], c.entities[delete+1:]...)
	}
}

func (c *ControlSystem) Update(dt float32) {
	for _, e := range c.entities {
		if e.MouseComponent.Enter {
			engo.SetCursor(engo.Hand)
		} else if e.MouseComponent.Leave {
			engo.SetCursor(nil)
		}
	}
}

func main() {
	opts := engo.RunOptions{
		Title:  "Mouse Demo",
		Width:  1024,
		Height: 640,
	}
	engo.Run(opts, &DefaultScene{})
}
