package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"sync"

	"engo.io/ecs"
	"engo.io/engo"
)

type PongGame struct{}

var (
	basicFont *engo.Font
)

type Ball struct {
	ecs.BasicEntity
	engo.RenderComponent
	engo.SpaceComponent
	engo.CollisionComponent
	SpeedComponent
}

type Score struct {
	ecs.BasicEntity
	engo.RenderComponent
	engo.SpaceComponent
}

type Paddle struct {
	ecs.BasicEntity
	ControlComponent
	engo.CollisionComponent
	engo.RenderComponent
	engo.SpaceComponent
}

func (pong *PongGame) Preload() {
	engo.Files.AddFromDir("assets", true)
}

func (pong *PongGame) Setup(w *ecs.World) {
	engo.SetBackground(color.Black)
	w.AddSystem(&engo.RenderSystem{})
	w.AddSystem(&engo.CollisionSystem{})
	w.AddSystem(&SpeedSystem{})
	w.AddSystem(&ControlSystem{})
	w.AddSystem(&BallSystem{})
	w.AddSystem(&ScoreSystem{})

	basicFont = (&engo.Font{URL: "Roboto-Regular.ttf", Size: 32, FG: color.NRGBA{255, 255, 255, 255}})
	if err := basicFont.CreatePreloaded(); err != nil {
		log.Fatalln("Could not load font:", err)
	}

	ballTexture := engo.Files.Image("ball.png")

	ball := Ball{BasicEntity: ecs.NewBasic()}
	ball.RenderComponent = engo.NewRenderComponent(ballTexture, engo.Point{2, 2}, "ball")
	ball.SpaceComponent = engo.SpaceComponent{
		Position: engo.Point{(engo.Width() - ballTexture.Width()) / 2, (engo.Height() - ballTexture.Height()) / 2},
		Width:    ballTexture.Width() * ball.RenderComponent.Scale().X,
		Height:   ballTexture.Height() * ball.RenderComponent.Scale().Y,
	}
	ball.CollisionComponent = engo.CollisionComponent{Main: true, Solid: true}
	ball.SpeedComponent = SpeedComponent{Point: engo.Point{300, 1000}}

	// Add our entity to the appropriate systems
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *engo.RenderSystem:
			sys.Add(&ball.BasicEntity, &ball.RenderComponent, &ball.SpaceComponent)
		case *engo.CollisionSystem:
			sys.Add(&ball.BasicEntity, &ball.CollisionComponent, &ball.SpaceComponent)
		case *SpeedSystem:
			sys.Add(&ball.BasicEntity, &ball.SpeedComponent, &ball.SpaceComponent)
		case *BallSystem:
			sys.Add(&ball.BasicEntity, &ball.SpeedComponent, &ball.SpaceComponent)
		}
	}

	score := Score{BasicEntity: ecs.NewBasic()}
	score.RenderComponent = engo.NewRenderComponent(basicFont.Render(" "), engo.Point{1, 1}, "YOLO <3")
	score.SpaceComponent = engo.SpaceComponent{engo.Point{100, 100}, 100, 100}

	// Add our entity to the appropriate systems
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *engo.RenderSystem:
			sys.Add(&score.BasicEntity, &score.RenderComponent, &score.SpaceComponent)
		case *ScoreSystem:
			sys.Add(&score.BasicEntity, &score.RenderComponent, &score.SpaceComponent)
		}
	}

	schemes := []string{"WASD", ""}
	paddleTexture := engo.Files.Image("paddle.png")

	for i := 0; i < 2; i++ {
		paddle := Paddle{BasicEntity: ecs.NewBasic()}
		paddle.RenderComponent = engo.NewRenderComponent(paddleTexture, engo.Point{2, 2}, "paddle")

		x := float32(0)
		if i != 0 {
			x = 800 - 16
		}

		paddle.SpaceComponent = engo.SpaceComponent{
			Position: engo.Point{x, (engo.Height() - paddleTexture.Height()) / 2},
			Width:    paddle.RenderComponent.Scale().X * paddleTexture.Width(),
			Height:   paddle.RenderComponent.Scale().Y * paddleTexture.Height(),
		}
		paddle.ControlComponent = ControlComponent{schemes[i]}
		paddle.CollisionComponent = engo.CollisionComponent{Main: false, Solid: true}

		// Add our entity to the appropriate systems
		for _, system := range w.Systems() {
			switch sys := system.(type) {
			case *engo.RenderSystem:
				sys.Add(&paddle.BasicEntity, &paddle.RenderComponent, &paddle.SpaceComponent)
			case *engo.CollisionSystem:
				sys.Add(&paddle.BasicEntity, &paddle.CollisionComponent, &paddle.SpaceComponent)
			case *ControlSystem:
				sys.Add(&paddle.BasicEntity, &paddle.ControlComponent, &paddle.SpaceComponent)
			}
		}
	}
}

func (*PongGame) Hide()        {}
func (*PongGame) Show()        {}
func (*PongGame) Exit()        {}
func (*PongGame) Type() string { return "PongGame" }

type SpeedComponent struct {
	engo.Point
}

type ControlComponent struct {
	Scheme string
}

type speedEntity struct {
	*ecs.BasicEntity
	*SpeedComponent
	*engo.SpaceComponent
}

type SpeedSystem struct {
	entities []speedEntity
}

func (s *SpeedSystem) New(*ecs.World) {
	engo.Mailbox.Listen("CollisionMessage", func(message engo.Message) {
		log.Println("collision")

		collision, isCollision := message.(engo.CollisionMessage)
		if isCollision {
			// See if we also have that Entity, and if so, change the speed
			for _, e := range s.entities {
				if e.ID() == collision.Entity.BasicEntity.ID() {
					e.SpeedComponent.X *= -1
				}
			}
		}
	})
}

func (s *SpeedSystem) Add(basic *ecs.BasicEntity, speed *SpeedComponent, space *engo.SpaceComponent) {
	s.entities = append(s.entities, speedEntity{basic, speed, space})
}

func (s *SpeedSystem) Remove(basic ecs.BasicEntity) {
	delete := -1
	for index, e := range s.entities {
		if e.BasicEntity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		s.entities = append(s.entities[:delete], s.entities[delete+1:]...)
	}
}

func (s *SpeedSystem) Update(dt float32) {
	for _, e := range s.entities {
		e.SpaceComponent.Position.X += e.SpeedComponent.X * dt
		e.SpaceComponent.Position.Y += e.SpeedComponent.Y * dt
	}
}

type ballEntity struct {
	*ecs.BasicEntity
	*SpeedComponent
	*engo.SpaceComponent
}

type BallSystem struct {
	entities []ballEntity
}

func (b *BallSystem) Add(basic *ecs.BasicEntity, speed *SpeedComponent, space *engo.SpaceComponent) {
	b.entities = append(b.entities, ballEntity{basic, speed, space})
}

func (b *BallSystem) Remove(basic ecs.BasicEntity) {
	delete := -1
	for index, e := range b.entities {
		if e.BasicEntity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		b.entities = append(b.entities[:delete], b.entities[delete+1:]...)
	}
}

func (b *BallSystem) Update(dt float32) {
	for _, e := range b.entities {
		if e.SpaceComponent.Position.X < 0 {
			engo.Mailbox.Dispatch(ScoreMessage{1})

			e.SpaceComponent.Position.X = 400 - 16
			e.SpaceComponent.Position.Y = 400 - 16
			e.SpeedComponent.X = 800 * rand.Float32()
			e.SpeedComponent.Y = 800 * rand.Float32()
		}

		if e.SpaceComponent.Position.Y < 0 {
			e.SpaceComponent.Position.Y = 0
			e.SpeedComponent.Y *= -1
		}

		if e.SpaceComponent.Position.X > (800 - 16) {
			engo.Mailbox.Dispatch(ScoreMessage{2})

			e.SpaceComponent.Position.X = 400 - 16
			e.SpaceComponent.Position.Y = 400 - 16
			e.SpeedComponent.X = 800 * rand.Float32()
			e.SpeedComponent.Y = 800 * rand.Float32()
		}

		if e.SpaceComponent.Position.Y > (800 - 16) {
			e.SpaceComponent.Position.Y = 800 - 16
			e.SpeedComponent.Y *= -1
		}
	}
}

type controlEntity struct {
	*ecs.BasicEntity
	*ControlComponent
	*engo.SpaceComponent
}

type ControlSystem struct {
	entities []controlEntity
}

func (c *ControlSystem) Add(basic *ecs.BasicEntity, control *ControlComponent, space *engo.SpaceComponent) {
	c.entities = append(c.entities, controlEntity{basic, control, space})
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
		up := false
		down := false
		if e.ControlComponent.Scheme == "WASD" {
			up = engo.Keys.Get(engo.W).Down()
			down = engo.Keys.Get(engo.S).Down()
		} else {
			up = engo.Keys.Get(engo.ArrowUp).Down()
			down = engo.Keys.Get(engo.ArrowDown).Down()
		}

		if up {
			if e.SpaceComponent.Position.Y > 0 {
				e.SpaceComponent.Position.Y -= 800 * dt
			}
		}

		if down {
			if (e.SpaceComponent.Height + e.SpaceComponent.Position.Y) < 800 {
				e.SpaceComponent.Position.Y += 800 * dt
			}
		}
	}
}

type scoreEntity struct {
	*ecs.BasicEntity
	*engo.RenderComponent
	*engo.SpaceComponent
}

type ScoreSystem struct {
	entities []scoreEntity

	PlayerOneScore, PlayerTwoScore int
	upToDate                       bool
	scoreLock                      sync.RWMutex
}

func (s *ScoreSystem) New(*ecs.World) {
	s.upToDate = true
	engo.Mailbox.Listen("ScoreMessage", func(message engo.Message) {
		scoreMessage, isScore := message.(ScoreMessage)
		if !isScore {
			return
		}

		s.scoreLock.Lock()
		if scoreMessage.Player != 1 {
			s.PlayerOneScore += 1
		} else {
			s.PlayerTwoScore += 1
		}
		log.Println("The score is now", s.PlayerOneScore, "vs", s.PlayerTwoScore)
		s.upToDate = false
		s.scoreLock.Unlock()
	})
}

func (c *ScoreSystem) Add(basic *ecs.BasicEntity, render *engo.RenderComponent, space *engo.SpaceComponent) {
	c.entities = append(c.entities, scoreEntity{basic, render, space})
}

func (s *ScoreSystem) Remove(basic ecs.BasicEntity) {
	delete := -1
	for index, e := range s.entities {
		if e.BasicEntity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		s.entities = append(s.entities[:delete], s.entities[delete+1:]...)
	}
}

func (s *ScoreSystem) Update(dt float32) {
	for _, e := range s.entities {
		if !s.upToDate {
			s.scoreLock.RLock()
			label := fmt.Sprintf("%v vs %v", s.PlayerOneScore, s.PlayerTwoScore)
			s.upToDate = true
			s.scoreLock.RUnlock()

			e.RenderComponent.SetDrawable(basicFont.Render(label))
			width := len(label) * 20

			e.SpaceComponent.Position.X = float32(400 - (width / 2))
		}
	}
}

type ScoreMessage struct {
	Player int
}

func (ScoreMessage) Type() string {
	return "ScoreMessage"
}

func main() {
	opts := engo.RunOptions{
		Title:         "Pong Demo",
		Width:         800,
		Height:        800,
		ScaleOnResize: true,
	}
	engo.Run(opts, &PongGame{})
}
