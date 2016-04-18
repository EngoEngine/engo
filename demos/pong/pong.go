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

	ball := ecs.NewEntity("RenderSystem", "CollisionSystem", "SpeedSystem", "BallSystem")
	ballTexture := engo.Files.Image("ball.png")
	ballRender := engo.NewRenderComponent(ballTexture, engo.Point{2, 2}, "ball")
	ballSpace := &engo.SpaceComponent{engo.Point{(engo.Width() - ballTexture.Width()) / 2, (engo.Height() - ballTexture.Height()) / 2}, ballTexture.Width() * ballRender.Scale().X, ballTexture.Height() * ballRender.Scale().Y}
	ballCollision := &engo.CollisionComponent{Main: true, Solid: true}
	ballSpeed := &SpeedComponent{}
	ballSpeed.Point = engo.Point{300, 100}

	ball.AddComponent(ballRender)
	ball.AddComponent(ballSpace)
	ball.AddComponent(ballCollision)
	ball.AddComponent(ballSpeed)
	err := w.AddEntity(ball)
	if err != nil {
		log.Println(err)
	}

	score := ecs.NewEntity("RenderSystem", "ScoreSystem")

	scoreRender := engo.NewRenderComponent(basicFont.Render(" "), engo.Point{1, 1}, "YOLO <3")

	scoreSpace := &engo.SpaceComponent{engo.Point{100, 100}, 100, 100}
	score.AddComponent(scoreRender)
	score.AddComponent(scoreSpace)
	err = w.AddEntity(score)
	if err != nil {
		log.Println(err)
	}

	schemes := []string{"WASD", ""}
	for i := 0; i < 2; i++ {
		paddle := ecs.NewEntity("RenderSystem", "CollisionSystem", "ControlSystem")
		paddleTexture := engo.Files.Image("paddle.png")
		paddleRender := engo.NewRenderComponent(paddleTexture, engo.Point{2, 2}, "paddle")
		x := float32(0)
		if i != 0 {
			x = 800 - 16
		}

		paddleSpace := &engo.SpaceComponent{engo.Point{x, (engo.Height() - paddleTexture.Height()) / 2}, paddleRender.Scale().X * paddleTexture.Width(), paddleRender.Scale().Y * paddleTexture.Height()}
		paddleControl := &ControlComponent{schemes[i]}
		paddleCollision := &engo.CollisionComponent{Main: false, Solid: true}
		paddle.AddComponent(paddleRender)
		paddle.AddComponent(paddleSpace)
		paddle.AddComponent(paddleControl)
		paddle.AddComponent(paddleCollision)
		err = w.AddEntity(paddle)
		if err != nil {
			log.Println(err)
		}
	}
}

func (*PongGame) Hide()        {}
func (*PongGame) Show()        {}
func (*PongGame) Exit()        {}
func (*PongGame) Type() string { return "PongGame" }

type SpeedSystem struct {
	ecs.LinearSystem
}

func (*SpeedSystem) Type() string { return "SpeedSystem" }

func (ms *SpeedSystem) New(*ecs.World) {
	engo.Mailbox.Listen("CollisionMessage", func(message engo.Message) {
		log.Println("collision")
		collision, isCollision := message.(engo.CollisionMessage)
		if isCollision {
			var speed *SpeedComponent
			if !collision.Entity.Component(&speed) {
				return
			}

			speed.X *= -1
		}
	})
}

func (ms *SpeedSystem) UpdateEntity(entity *ecs.Entity, dt float32) {
	var speed *SpeedComponent
	var space *engo.SpaceComponent
	if !entity.Component(&speed) || !entity.Component(&space) {
		return
	}
	space.Position.X += speed.X * dt
	space.Position.Y += speed.Y * dt
}

func (ms *SpeedSystem) Receive(message engo.Message) {}

type SpeedComponent struct {
	engo.Point
}

func (*SpeedComponent) Type() string {
	return "SpeedComponent"
}

type BallSystem struct {
	ecs.LinearSystem
}

func (*BallSystem) Type() string { return "BallSystem" }

func (bs *BallSystem) New(*ecs.World) {}

func (bs *BallSystem) UpdateEntity(entity *ecs.Entity, dt float32) {
	var space *engo.SpaceComponent
	var speed *SpeedComponent
	if !entity.Component(&space) || !entity.Component(&speed) {
		return
	}

	if space.Position.X < 0 {
		engo.Mailbox.Dispatch(ScoreMessage{1})

		space.Position.X = 400 - 16
		space.Position.Y = 400 - 16
		speed.X = 800 * rand.Float32()
		speed.Y = 800 * rand.Float32()
	}

	if space.Position.Y < 0 {
		space.Position.Y = 0
		speed.Y *= -1
	}

	if space.Position.X > (800 - 16) {
		engo.Mailbox.Dispatch(ScoreMessage{2})

		space.Position.X = 400 - 16
		space.Position.Y = 400 - 16
		speed.X = 800 * rand.Float32()
		speed.Y = 800 * rand.Float32()
	}

	if space.Position.Y > (800 - 16) {
		space.Position.Y = 800 - 16
		speed.Y *= -1
	}
}

type ControlSystem struct {
	ecs.LinearSystem
}

func (*ControlSystem) Type() string { return "ControlSystem" }

func (c *ControlSystem) New(*ecs.World) {}

func (c *ControlSystem) UpdateEntity(entity *ecs.Entity, dt float32) {
	//Check scheme
	// -Move entity based on that
	var control *ControlComponent
	var space *engo.SpaceComponent

	if !entity.Component(&space) || !entity.Component(&control) {
		return
	}
	up := false
	down := false
	if control.Scheme == "WASD" {
		up = engo.Keys.Get(engo.W).Down()
		down = engo.Keys.Get(engo.S).Down()
	} else {
		up = engo.Keys.Get(engo.ArrowUp).Down()
		down = engo.Keys.Get(engo.ArrowDown).Down()
	}

	if up {
		if space.Position.Y > 0 {
			space.Position.Y -= 800 * dt
		}
	}

	if down {
		if (space.Height + space.Position.Y) < 800 {
			space.Position.Y += 800 * dt
		}
	}

}

type ControlComponent struct {
	Scheme string
}

func (*ControlComponent) Type() string {
	return "ControlComponent"
}

type ScoreSystem struct {
	ecs.LinearSystem
	PlayerOneScore, PlayerTwoScore int
	upToDate                       bool
	scoreLock                      sync.RWMutex
}

func (*ScoreSystem) Type() string { return "ScoreSystem" }

func (sc *ScoreSystem) New(*ecs.World) {
	sc.upToDate = true
	engo.Mailbox.Listen("ScoreMessage", func(message engo.Message) {
		scoreMessage, isScore := message.(ScoreMessage)
		if !isScore {
			return
		}

		sc.scoreLock.Lock()
		if scoreMessage.Player != 1 {
			sc.PlayerOneScore += 1
		} else {
			sc.PlayerTwoScore += 1
		}
		log.Println("The score is now", sc.PlayerOneScore, "vs", sc.PlayerTwoScore)
		sc.upToDate = false
		sc.scoreLock.Unlock()
	})
}

func (c *ScoreSystem) UpdateEntity(entity *ecs.Entity, dt float32) {
	var render *engo.RenderComponent
	var space *engo.SpaceComponent

	if !entity.Component(&render) || !entity.Component(&space) {
		return
	}

	if !c.upToDate {
		c.scoreLock.RLock()
		label := fmt.Sprintf("%v vs %v", c.PlayerOneScore, c.PlayerTwoScore)
		c.upToDate = true
		c.scoreLock.RUnlock()

		render.SetDrawable(basicFont.Render(label))
		width := len(label) * 20

		space.Position.X = float32(400 - (width / 2))
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
