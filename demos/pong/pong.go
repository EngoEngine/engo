package main

import (
	"github.com/paked/engi"
)

type PongGame struct {
	engi.World
}

func (pong PongGame) Preload() {
	engi.Files.Add("ball", "assets/ball.png")
	engi.Files.Add("paddle", "assets/paddle.png")
}

func (pong *PongGame) Setup() {
	engi.SetBg(0x2d3739)
	pong.AddSystem(&engi.RenderSystem{})
	pong.AddSystem(&engi.CollisionSystem{})
	pong.AddSystem(&MovementSystem{})

	ball := engi.NewEntity([]string{"RenderSystem", "CollisionSystem", "MovementSystem"})
	ballTexture := engi.Files.Image("ball")
	ballRender := engi.NewRenderComponent(ballTexture, engi.Point{2, 2}, "ball")
	ballSpace := engi.SpaceComponent{engi.Point{10, 10}, ballTexture.Width() * ballRender.Scale.X, ballTexture.Height() * ballRender.Scale.Y}
	ballCollisionMaster := engi.CollisionMasterComponent{}
	ballSpeed := SpeedComponent{}
	ballSpeed.Point = engi.Point{100, 0}
	ball.AddComponent(&ballRender)
	ball.AddComponent(&ballSpace)
	ball.AddComponent(&ballCollisionMaster)
	ball.AddComponent(&ballSpeed)
	pong.AddEntity(ball)

	for i := 0; i < 2; i++ {
		paddle := engi.NewEntity([]string{"RenderSystem", "CollisionSystem"})
		paddleTexture := engi.Files.Image("paddle")
		paddleRender := engi.NewRenderComponent(paddleTexture, engi.Point{2, 2}, "paddle")
		paddleSpace := engi.SpaceComponent{engi.Point{100 * float32(i), 10}, paddleRender.Scale.X * paddleTexture.Width(), paddleRender.Scale.Y * paddleTexture.Height()}
		paddle.AddComponent(&paddleRender)
		paddle.AddComponent(&paddleSpace)
		pong.AddEntity(paddle)
	}
}

type MovementSystem struct {
	*engi.System
}

func (ms *MovementSystem) New() {
	ms.System = &engi.System{}
}

func (ms MovementSystem) Name() string {
	return "MovementSystem"
}

func (ms MovementSystem) Update(entity *engi.Entity, dt float32) {
	speed, hasSpeed := entity.GetComponent("SpeedComponent").(*SpeedComponent)
	space, hasSpace := entity.GetComponent("SpaceComponent").(*engi.SpaceComponent)
	if hasSpeed && hasSpace {
		space.Position.X += speed.X * dt
		space.Position.Y += speed.Y * dt
	}
}

type SpeedComponent struct {
	engi.Point
}

func (speed SpeedComponent) Name() string {
	return "SpeedComponent"
}

func main() {
	engi.Open("Pong", 800, 800, false, &PongGame{})
}
