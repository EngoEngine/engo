package main

import (
	"github.com/ajhager/eng"
	"github.com/ajhager/eng/spine"
)

var (
	skeleton *spine.Skeleton
	walk     *spine.Animation
	batch    *eng.Batch
	animTime float32
	dir      float32
)

type Game struct {
	*eng.Game
}

func (g *Game) Init(config *eng.Config) {
	config.Title = "Spine"
}

func (g *Game) Open() {
	batch = eng.NewBatch()
	skeleton = spine.NewSkeleton("data/spine", "spineboy.json")
	skeleton.X = 100
	skeleton.Y = 512
	skeleton.SetToSetupPose()
	walk = skeleton.Animation("walk")
	dir = 1
}

func (g *Game) Update(dt float32) {
	animTime += dt
	skeleton.X += dir * 200 * dt
	if skeleton.X < 100 {
		skeleton.FlipX = false
		dir = -dir
	}
	if skeleton.X > 924 {
		skeleton.FlipX = true
		dir = -dir
	}
	skeleton.Apply(walk, animTime, true)
	skeleton.UpdateWorldTransform()
}

func (g *Game) Draw() {
	batch.Begin()
	skeleton.Draw(batch)
	batch.End()
}

func main() {
	eng.Run(new(Game))
}
