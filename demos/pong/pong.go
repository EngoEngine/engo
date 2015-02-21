package main

import (
	"fmt"
	"github.com/paked/engi"
	"log"
	"math/rand"
)

type PongGame struct {
	engi.World
}

func (pong PongGame) Preload() {
	engi.Files.Add(engi.NewResource("ball", "assets/ball.png"),
		engi.NewResource("paddle", "assets/paddle.png"),
		engi.NewResource("font", "assets/font.png"))
}

func (pong *PongGame) Setup() {
	engi.SetBg(0x2d3739)
	pong.AddSystem(&engi.RenderSystem{})
	pong.AddSystem(&engi.CollisionSystem{})
	pong.AddSystem(&SpeedSystem{})
	pong.AddSystem(&ControlSystem{})
	pong.AddSystem(&BallSystem{})
	pong.AddSystem(&ScoreSystem{})

	ball := engi.NewEntity([]string{"RenderSystem", "CollisionSystem", "SpeedSystem", "BallSystem"})
	ballTexture := engi.Files.Image("ball")
	ballRender := engi.NewRenderComponent(ballTexture, engi.Point{2, 2}, "ball")
	ballSpace := engi.SpaceComponent{engi.Point{(engi.Width() - ballTexture.Width()) / 2, (engi.Height() - ballTexture.Height()) / 2}, ballTexture.Width() * ballRender.Scale.X, ballTexture.Height() * ballRender.Scale.Y}
	ballCollision := engi.CollisionComponent{Main: true, Solid: true}
	ballSpeed := SpeedComponent{}
	ballSpeed.Point = engi.Point{300, 100}
	ball.AddComponent(&ballRender)
	ball.AddComponent(&ballSpace)
	ball.AddComponent(&ballCollision)
	ball.AddComponent(&ballSpeed)
	pong.AddEntity(ball)

	score := engi.NewEntity([]string{"RenderSystem", "ScoreSystem"})
	scoreTex := engi.NewGridFont(engi.Files.Image("font"), 20, 20)
	scoreRender := engi.NewRenderComponent(scoreTex, engi.Point{1, 1}, "YOLO <3")
	scoreSpace := engi.SpaceComponent{engi.Point{100, 100}, 10, 10}
	score.AddComponent(&scoreRender)
	score.AddComponent(&scoreSpace)
	pong.AddEntity(score)

	schemes := []string{"WASD", ""}
	for i := 0; i < 2; i++ {
		paddle := engi.NewEntity([]string{"RenderSystem", "CollisionSystem", "ControlSystem"})
		paddleTexture := engi.Files.Image("paddle")
		paddleRender := engi.NewRenderComponent(paddleTexture, engi.Point{2, 2}, "paddle")
		x := float32(0)
		if i != 0 {
			x = 800 - 16
		}
		paddleSpace := engi.SpaceComponent{engi.Point{x, (engi.Height() - paddleTexture.Height()) / 2}, paddleRender.Scale.X * paddleTexture.Width(), paddleRender.Scale.Y * paddleTexture.Height()}
		paddleControl := ControlComponent{schemes[i]}
		paddleCollision := engi.CollisionComponent{Main: false, Solid: true}
		paddle.AddComponent(&paddleRender)
		paddle.AddComponent(&paddleSpace)
		paddle.AddComponent(&paddleControl)
		paddle.AddComponent(&paddleCollision)
		pong.AddEntity(paddle)
	}
}

type SpeedSystem struct {
	*engi.System
}

func (ms *SpeedSystem) New() {
	ms.System = &engi.System{}
	engi.Mailbox.Listen("CollisionMessage", func(message interface{}) {
		log.Println("collision")
		collision, isCollision := message.(engi.CollisionMessage)
		if isCollision {
			var speed *SpeedComponent
			if !collision.Entity.GetComponent(&speed) {
				return
			}

			speed.X *= -1
		}
	})
}

func (ms SpeedSystem) Name() string {
	return "SpeedSystem"
}

func (ms SpeedSystem) Update(entity *engi.Entity, dt float32) {
	var speed *SpeedComponent
	var space *engi.SpaceComponent
	if !entity.GetComponent(&speed) || !entity.GetComponent(&space) {
		return
	}
	space.Position.X += speed.X * dt
	space.Position.Y += speed.Y * dt
}

func (ms SpeedSystem) Receive(message engi.Message) {
}

type SpeedComponent struct {
	engi.Point
}

func (speed SpeedComponent) Name() string {
	return "SpeedComponent"
}

type BallSystem struct {
	*engi.System
}

func (bs *BallSystem) New() {
	bs.System = &engi.System{}
}

func (bs BallSystem) Name() string {
	return "BallSystem"
}

func (bs *BallSystem) Update(entity *engi.Entity, dt float32) {
	var space *engi.SpaceComponent
	var speed *SpeedComponent
	if !entity.GetComponent(&space) || !entity.GetComponent(&speed) {
		return
	}

	if space.Position.X < 0 {
		engi.Mailbox.Dispatch("ScoreMesasge", ScoreMessage{1})

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
		engi.Mailbox.Dispatch("ScoreMessage", ScoreMessage{2})

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
	*engi.System
}

func (c *ControlSystem) Name() string {
	return "ControlSystem"
}
func (c *ControlSystem) New() {
	c.System = &engi.System{}
}

func (c *ControlSystem) Update(entity *engi.Entity, dt float32) {
	//Check scheme
	// -Move entity based on that
	var control *ControlComponent
	var space *engi.SpaceComponent

	if !entity.GetComponent(&space) || !entity.GetComponent(&control) {
		return
	}
	up := false
	down := false
	if control.Scheme == "WASD" {
		up = engi.Keys.KEY_W.Down()
		down = engi.Keys.KEY_S.Down()
	} else {
		up = engi.Keys.KEY_UP.Down()
		down = engi.Keys.KEY_DOWN.Down()
	}

	if up {
		space.Position.Y -= 800 * dt
	}

	if down {
		space.Position.Y += 800 * dt
	}

}

type ControlComponent struct {
	Scheme string
}

func (cs ControlComponent) Name() string {
	return "ControlComponent"
}

type ScoreSystem struct {
	*engi.System
	PlayerOneScore, PlayerTwoScore int
}

func (score *ScoreSystem) Name() string {
	return "ScoreSystem"
}

func (sc *ScoreSystem) New() {
	sc.System = &engi.System{}
	engi.Mailbox.Listen("ScoreMessage", func(message interface{}) {
		scoreMessage, isScore := message.(ScoreMessage)
		if !isScore {
			return
		}

		if scoreMessage.Player != 1 {
			sc.PlayerOneScore += 1
		} else {
			sc.PlayerTwoScore += 1
		}
	})
}

func (sc *ScoreSystem) Receive(message engi.Message) {

}

func (c *ScoreSystem) Update(entity *engi.Entity, dt float32) {
	var render *engi.RenderComponent
	var space *engi.SpaceComponent
	if !entity.GetComponent(&render) || !entity.GetComponent(&space) {
		return
	}

	render.Label = fmt.Sprintf("%v vs %v", c.PlayerOneScore, c.PlayerTwoScore)

	width := len(render.Label) * 20

	space.Position.X = float32(400 - (width / 2))

}

type ScoreMessage struct {
	Player int
}

func (score ScoreMessage) Type() string {
	return "ScoreMessage"
}

func main() {
	engi.Open("Pong", 800, 800, false, &PongGame{})
}
