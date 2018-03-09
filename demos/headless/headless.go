//+build demo

package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"sync"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
)

type PongGame struct{}

var (
	basicFont *common.Font
)

type Ball struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
	common.CollisionComponent
	SpeedComponent
}

type Score struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

type Paddle struct {
	ecs.BasicEntity
	ControlComponent
	common.CollisionComponent
	common.RenderComponent
	common.SpaceComponent
}

func (pong *PongGame) Preload() {
	err := engo.Files.Load("ball.png", "paddle.png", "Roboto-Regular.ttf")
	if err != nil {
		log.Println(err)
	}
}

func (pong *PongGame) Setup(u engo.Updater) {
	w, _ := u.(*ecs.World)

	common.SetBackground(color.Black)

	w.AddSystem(&common.RenderSystem{})
	w.AddSystem(&common.CollisionSystem{Solids: 1})
	w.AddSystem(&SpeedSystem{})
	w.AddSystem(&ControlSystem{})
	w.AddSystem(&BallSystem{})
	w.AddSystem(&ScoreSystem{})

	basicFont = (&common.Font{URL: "Roboto-Regular.ttf", Size: 32, FG: color.NRGBA{255, 255, 255, 255}})
	if err := basicFont.CreatePreloaded(); err != nil {
		log.Println("Could not load font:", err)
	}

	ballTexture, err := common.LoadedSprite("ball.png")
	if err != nil {
		log.Println(err)
	}

	ball := Ball{BasicEntity: ecs.NewBasic()}
	ball.RenderComponent = common.RenderComponent{
		Drawable: ballTexture,
		Scale:    engo.Point{2, 2},
	}
	ball.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{(engo.GameWidth() - ballTexture.Width()) / 2, (engo.GameHeight() - ballTexture.Height()) / 2},
		Width:    ballTexture.Width() * ball.RenderComponent.Scale.X,
		Height:   ballTexture.Height() * ball.RenderComponent.Scale.Y,
	}
	ball.CollisionComponent = common.CollisionComponent{
		Main: 1,
	}
	ball.SpeedComponent = SpeedComponent{Point: engo.Point{300, 1000}}

	// Add our entity to the appropriate systems
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&ball.BasicEntity, &ball.RenderComponent, &ball.SpaceComponent)
		case *common.CollisionSystem:
			sys.Add(&ball.BasicEntity, &ball.CollisionComponent, &ball.SpaceComponent)
		case *SpeedSystem:
			sys.Add(&ball.BasicEntity, &ball.SpeedComponent, &ball.SpaceComponent)
		case *BallSystem:
			sys.Add(&ball.BasicEntity, &ball.SpeedComponent, &ball.SpaceComponent)
		}
	}

	score := Score{BasicEntity: ecs.NewBasic()}
	score.RenderComponent = common.RenderComponent{Drawable: basicFont.Render(" ")}
	score.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{100, 100},
		Width:    100,
		Height:   100,
	}

	// Add our entity to the appropriate systems
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&score.BasicEntity, &score.RenderComponent, &score.SpaceComponent)
		case *ScoreSystem:
			sys.Add(&score.BasicEntity, &score.RenderComponent, &score.SpaceComponent)
		}
	}

	schemes := []string{"WASD", ""}
	paddleTexture, err := common.LoadedSprite("paddle.png")
	if err != nil {
		log.Println(err)
	}

	for i := 0; i < 2; i++ {
		paddle := Paddle{BasicEntity: ecs.NewBasic()}
		paddle.RenderComponent = common.RenderComponent{
			Drawable: paddleTexture,
			Scale:    engo.Point{2, 2},
		}

		x := float32(0)
		if i != 0 {
			x = 800 - 16
		}

		paddle.SpaceComponent = common.SpaceComponent{
			Position: engo.Point{x, (engo.GameHeight() - paddleTexture.Height()) / 2},
			Width:    paddle.RenderComponent.Scale.X * paddleTexture.Width(),
			Height:   paddle.RenderComponent.Scale.Y * paddleTexture.Height(),
		}
		paddle.ControlComponent = ControlComponent{schemes[i]}
		paddle.CollisionComponent = common.CollisionComponent{
			Group: 1,
		}

		// Add our entity to the appropriate systems
		for _, system := range w.Systems() {
			switch sys := system.(type) {
			case *common.RenderSystem:
				sys.Add(&paddle.BasicEntity, &paddle.RenderComponent, &paddle.SpaceComponent)
			case *common.CollisionSystem:
				sys.Add(&paddle.BasicEntity, &paddle.CollisionComponent, &paddle.SpaceComponent)
			case *ControlSystem:
				sys.Add(&paddle.BasicEntity, &paddle.ControlComponent, &paddle.SpaceComponent)
			}
		}
	}
}

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
	*common.SpaceComponent
}

type SpeedSystem struct {
	entities []speedEntity
}

func (s *SpeedSystem) New(*ecs.World) {
	engo.Mailbox.Listen("CollisionMessage", func(message engo.Message) {
		log.Println("collision")

		collision, isCollision := message.(common.CollisionMessage)
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

func (s *SpeedSystem) Add(basic *ecs.BasicEntity, speed *SpeedComponent, space *common.SpaceComponent) {
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
	*common.SpaceComponent
}

type BallSystem struct {
	entities []ballEntity
}

func (b *BallSystem) Add(basic *ecs.BasicEntity, speed *SpeedComponent, space *common.SpaceComponent) {
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
	*common.SpaceComponent
}

type ControlSystem struct {
	entities []controlEntity
}

func (c *ControlSystem) Add(basic *ecs.BasicEntity, control *ControlComponent, space *common.SpaceComponent) {
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
		speed := 800 * dt

		vert := engo.Input.Axis(engo.DefaultVerticalAxis)
		e.SpaceComponent.Position.Y += speed * vert.Value()

		if (e.SpaceComponent.Height + e.SpaceComponent.Position.Y) > 800 {
			e.SpaceComponent.Position.Y = 800
		} else if e.SpaceComponent.Position.Y < 0 {
			e.SpaceComponent.Position.Y = 0
		}
	}
}

type scoreEntity struct {
	*ecs.BasicEntity
	*common.RenderComponent
	*common.SpaceComponent
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

func (c *ScoreSystem) Add(basic *ecs.BasicEntity, render *common.RenderComponent, space *common.SpaceComponent) {
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

			// Clean up old one to prevent data leakage
			e.RenderComponent.Drawable.Close()

			e.RenderComponent.Drawable = basicFont.Render(label)
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
		HeadlessMode:   true,
		StandardInputs: true,
	}

	engo.Run(opts, &PongGame{})
}
