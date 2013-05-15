package main

import (
	"github.com/ajhager/eng"
	"github.com/ajhager/eng/spine"
)

var skeleton *spine.Skeleton
var batch *eng.Batch

type Game struct {
	*eng.Game
}

func (g *Game) Init(config *eng.Config) {
	config.Title = "Spine"
}

func (g *Game) Open() {
	batch = eng.NewBatch()
	skeleton = spine.NewSkeleton("data/spine", "spineboy.json")
	skeleton.X = 512
	skeleton.Y = 512
	skeleton.FlipY = true
	skeleton.SetToSetupPose()
}

func (g *Game) Update(dt float32) {
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
