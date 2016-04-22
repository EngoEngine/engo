package main

import (
	"image/color"
	"log"
	"math/rand"

	"engo.io/ecs"
	"engo.io/engo"
)

type Guy struct {
	ecs.BasicEntity
	engo.CollisionComponent
	engo.RenderComponent
	engo.SpaceComponent
}

type Rock struct {
	ecs.BasicEntity
	engo.CollisionComponent
	engo.RenderComponent
	engo.SpaceComponent
}

type DefaultScene struct{}

func (*DefaultScene) Preload() {
	// Add all the files in the data directory non recursively
	engo.Files.AddFromDir("data", false)
}

func (*DefaultScene) Setup(w *ecs.World) {
	engo.SetBackground(color.White)

	// Add all of the systems
	w.AddSystem(&engo.RenderSystem{})
	w.AddSystem(&engo.CollisionSystem{})
	w.AddSystem(&DeathSystem{})
	w.AddSystem(&FallingSystem{})
	w.AddSystem(&ControlSystem{})
	w.AddSystem(&RockSpawnSystem{})

	texture := engo.Files.Image("icon.png")

	// Create an entity
	guy := Guy{BasicEntity: ecs.NewBasic()}

	// Initialize the components, set scale to 4x
	guy.RenderComponent = engo.NewRenderComponent(texture, engo.Point{4, 4})
	guy.SpaceComponent = engo.SpaceComponent{
		Position: engo.Point{0, 0},
		Width:    texture.Width() * guy.RenderComponent.Scale().X,
		Height:   texture.Height() * guy.RenderComponent.Scale().Y,
	}
	guy.CollisionComponent = engo.CollisionComponent{Solid: true, Main: true}

	// Add it to appropriate systems
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *engo.RenderSystem:
			sys.Add(&guy.BasicEntity, &guy.RenderComponent, &guy.SpaceComponent)
		case *engo.CollisionSystem:
			sys.Add(&guy.BasicEntity, &guy.CollisionComponent, &guy.SpaceComponent)
		case *ControlSystem:
			sys.Add(&guy.BasicEntity, &guy.SpaceComponent)
		}
	}
}

func (*DefaultScene) Type() string { return "Game" }

type controlEntity struct {
	*ecs.BasicEntity
	*engo.SpaceComponent
}

type ControlSystem struct {
	entities []controlEntity
}

func (c *ControlSystem) Add(basic *ecs.BasicEntity, space *engo.SpaceComponent) {
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
		if engo.Keys.Get(engo.A).Down() {
			e.SpaceComponent.Position.X -= speed
		}

		if engo.Keys.Get(engo.D).Down() {
			e.SpaceComponent.Position.X += speed
		}

		if engo.Keys.Get(engo.W).Down() {
			e.SpaceComponent.Position.Y -= speed
		}

		if engo.Keys.Get(engo.S).Down() {
			e.SpaceComponent.Position.Y += speed
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
		X: rand.Float32() * engo.Width(),
		Y: -32,
	}
	NewRock(rock.world, position)
}

func NewRock(world *ecs.World, position engo.Point) {
	texture := engo.Files.Image("rock.png")

	rock := Rock{BasicEntity: ecs.NewBasic()}
	rock.RenderComponent = engo.NewRenderComponent(texture, engo.Point{4, 4})
	rock.SpaceComponent = engo.SpaceComponent{
		Position: position,
		Width:    texture.Width() * rock.RenderComponent.Scale().X,
		Height:   texture.Height() * rock.RenderComponent.Scale().Y,
	}
	rock.CollisionComponent = engo.CollisionComponent{Solid: true}

	for _, system := range world.Systems() {
		switch sys := system.(type) {
		case *engo.RenderSystem:
			sys.Add(&rock.BasicEntity, &rock.RenderComponent, &rock.SpaceComponent)
		case *engo.CollisionSystem:
			sys.Add(&rock.BasicEntity, &rock.CollisionComponent, &rock.SpaceComponent)
		case *FallingSystem:
			sys.Add(&rock.BasicEntity, &rock.SpaceComponent)
		}
	}
}

type fallingEntity struct {
	*ecs.BasicEntity
	*engo.SpaceComponent
}

type FallingSystem struct {
	entities []fallingEntity
}

func (f *FallingSystem) Add(basic *ecs.BasicEntity, space *engo.SpaceComponent) {
	f.entities = append(f.entities, fallingEntity{basic, space})
}

func (f *FallingSystem) Remove(basic ecs.BasicEntity) {
	delete := -1
	for index, e := range f.entities {
		if e.BasicEntity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		f.entities = append(f.entities[:delete], f.entities[delete+1:]...)
	}
}

func (f *FallingSystem) Update(dt float32) {
	for _, e := range f.entities {
		e.SpaceComponent.Position.Y += 200 * dt
	}
}

type DeathSystem struct{}

func (*DeathSystem) New(*ecs.World) {
	// Subscribe to ScoreMessage
	engo.Mailbox.Listen("CollisionMessage", func(message engo.Message) {
		_, isCollision := message.(engo.CollisionMessage)
		if isCollision {
			log.Println("DEAD")
		}
	})
}

func (*DeathSystem) Remove(ecs.BasicEntity) {}
func (*DeathSystem) Update(dt float32)      {}

func main() {
	opts := engo.RunOptions{
		Title:  "Falling Demo",
		Width:  1024,
		Height: 640,
	}
	engo.Run(opts, &DefaultScene{})
}
