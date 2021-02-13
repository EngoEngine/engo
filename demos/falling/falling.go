package main

import (
	"image/color"
	"log"
	"math/rand"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type Guy struct {
	ecs.BasicEntity
	common.CollisionComponent
	common.RenderComponent
	common.SpaceComponent
}

type Rock struct {
	ecs.BasicEntity
	common.CollisionComponent
	common.RenderComponent
	common.SpaceComponent
}

type DefaultScene struct{}

func (*DefaultScene) Preload() {
	err := engo.Files.Load("icon.png", "rock.png", "rock2.png")
	if err != nil {
		log.Println(err)
	}
}

func (*DefaultScene) Setup(u engo.Updater) {
	w, _ := u.(*ecs.World)

	common.SetBackground(color.White)

	// Add all of the systems
	w.AddSystem(&common.RenderSystem{})
	w.AddSystem(&common.CollisionSystem{Solids: 1})
	w.AddSystem(&DeathSystem{})
	w.AddSystem(&FallingSystem{})
	w.AddSystem(&ControlSystem{})
	w.AddSystem(&RockSpawnSystem{})

	texture, err := common.LoadedSprite("icon.png")
	if err != nil {
		log.Println(err)
	}

	// Create an entity
	guy := Guy{BasicEntity: ecs.NewBasic()}

	// Initialize the components, set scale to 4x
	guy.RenderComponent = common.RenderComponent{
		Drawable: texture,
		Scale:    engo.Point{4, 4},
	}
	guy.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{0, 0},
		Width:    texture.Width() * guy.RenderComponent.Scale.X,
		Height:   texture.Height() * guy.RenderComponent.Scale.Y,
	}
	guy.CollisionComponent = common.CollisionComponent{
		Main: 1,
	}

	// Add it to appropriate systems
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&guy.BasicEntity, &guy.RenderComponent, &guy.SpaceComponent)
		case *common.CollisionSystem:
			sys.Add(&guy.BasicEntity, &guy.CollisionComponent, &guy.SpaceComponent)
		case *ControlSystem:
			sys.Add(&guy.BasicEntity, &guy.SpaceComponent)
		}
	}
}

func (*DefaultScene) Type() string { return "Game" }

type controlEntity struct {
	*ecs.BasicEntity
	*common.SpaceComponent
}

type ControlSystem struct {
	entities []controlEntity
}

func (c *ControlSystem) Add(basic *ecs.BasicEntity, space *common.SpaceComponent) {
	c.entities = append(c.entities, controlEntity{basic, space})
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
	speed := 400 * dt

	for _, e := range c.entities {
		hori := engo.Input.Axis(engo.DefaultHorizontalAxis)
		e.SpaceComponent.Position.X += speed * hori.Value()

		vert := engo.Input.Axis(engo.DefaultVerticalAxis)
		e.SpaceComponent.Position.Y += speed * vert.Value()

		if e.SpaceComponent.Position.Y > engo.GameHeight() {
			c.Remove(*e.BasicEntity)
		}
	}
}

type RockSpawnSystem struct {
	world *ecs.World
}

func (rock *RockSpawnSystem) New(w *ecs.World) {
	rock.world = w
}

func (*RockSpawnSystem) Remove(ecs.BasicEntity) {}

func (rock *RockSpawnSystem) Update(dt float32) {
	// 4% change of spawning a rock each frame
	if rand.Float32() < .96 {
		return
	}

	position := engo.Point{
		X: rand.Float32() * engo.GameWidth(),
		Y: -32,
	}
	NewRock(rock.world, position)
}

func NewRock(world *ecs.World, position engo.Point) {
	rock := Rock{BasicEntity: ecs.NewBasic()}
	switch rand.Intn(3) {
	case 0:
		yscale := 1.0 + rand.Float32()
		rock.RenderComponent = common.RenderComponent{
			Drawable: common.Circle{},
			Color:    color.RGBA{0, 0, 0, 255},
		}
		rock.SpaceComponent = common.SpaceComponent{
			Position: position,
			Width:    16 * 4,
			Height:   16 * 4 * yscale,
			Rotation: 45 * rand.Float32(),
		}
		rock.AddShape(common.Shape{
			Ellipse: common.Ellipse{
				Rx: 32,
				Ry: 32 * yscale,
				Cx: 32,
				Cy: 32 * yscale,
			},
		})
	case 1:
		texture, _ := common.LoadedSprite("rock2.png")
		rock.RenderComponent = common.RenderComponent{
			Drawable: texture,
			Scale:    engo.Point{X: 4, Y: 4},
		}
		rock.SpaceComponent = common.SpaceComponent{
			Position: position,
			Width:    texture.Width() * rock.RenderComponent.Scale.X,
			Height:   texture.Height() * rock.RenderComponent.Scale.Y,
			Rotation: 45 * rand.Float32(),
		}
		pts := []float32{4, 0, 12, 0, 16, 4, 16, 13, 13, 13, 13, 16, 3, 16, 3, 13, 0, 13, 0, 4, 4, 0}
		lines := []engo.Line{}
		for i := 0; i < len(pts)-3; i += 2 {
			line := engo.Line{
				P1: engo.Point{
					X: pts[i] * 4,
					Y: pts[i+1] * 4,
				},
				P2: engo.Point{
					X: pts[i+2] * 4,
					Y: pts[i+3] * 4,
				},
			}
			lines = append(lines, line)
		}
		rock.AddShape(common.Shape{Lines: lines})
	default:
		texture, _ := common.LoadedSprite("rock.png")
		rock.RenderComponent = common.RenderComponent{
			Drawable: texture,
			Scale:    engo.Point{X: 4, Y: 4},
		}
		rock.SpaceComponent = common.SpaceComponent{
			Position: position,
			Width:    texture.Width() * rock.RenderComponent.Scale.X,
			Height:   texture.Height() * rock.RenderComponent.Scale.Y,
			Rotation: 45 * rand.Float32(),
		}
	}
	rock.CollisionComponent = common.CollisionComponent{Group: 1}

	for _, system := range world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&rock.BasicEntity, &rock.RenderComponent, &rock.SpaceComponent)
		case *common.CollisionSystem:
			sys.Add(&rock.BasicEntity, &rock.CollisionComponent, &rock.SpaceComponent)
		case *FallingSystem:
			sys.Add(&rock.BasicEntity, &rock.SpaceComponent)
		}
	}
}

type fallingEntity struct {
	*ecs.BasicEntity
	*common.SpaceComponent
}

type FallingSystem struct {
	entities []fallingEntity
	world    *ecs.World
}

func (f *FallingSystem) Add(basic *ecs.BasicEntity, space *common.SpaceComponent) {
	f.entities = append(f.entities, fallingEntity{basic, space})
}

func (f *FallingSystem) Remove(basic ecs.BasicEntity) {
	for i, e := range f.entities {
		if e.BasicEntity.ID() == basic.ID() {
			for _, system := range f.world.Systems() {
				switch system.(type) {
				case *FallingSystem:
				default:
					system.Remove(*e.BasicEntity)
				}
			}
			f.entities = append(f.entities[:i], f.entities[i+1:]...)
			break
		}
	}
}

func (f *FallingSystem) Update(dt float32) {
	speed := 400 * dt

	for _, e := range f.entities {
		e.SpaceComponent.Position.Y += speed

		if e.SpaceComponent.Position.Y > engo.GameHeight() {
			f.Remove(*e.BasicEntity)
		}
	}
}

func (f *FallingSystem) New(world *ecs.World) {
	f.world = world
}

type DeathSystem struct{}

func (*DeathSystem) New(*ecs.World) {
	// Subscribe to ScoreMessage
	engo.Mailbox.Listen("CollisionMessage", func(message engo.Message) {
		_, isCollision := message.(common.CollisionMessage)

		if isCollision {
			log.Println("DEAD")
		}
	})
}

func (*DeathSystem) Remove(ecs.BasicEntity) {}
func (*DeathSystem) Update(float32)         {}

func main() {
	opts := engo.RunOptions{
		Title:          "Falling Demo",
		Width:          1024,
		Height:         640,
		StandardInputs: true,
	}

	engo.Run(opts, &DefaultScene{})
}
