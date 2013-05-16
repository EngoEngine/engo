package main

import (
	"github.com/ajhager/eng"
	"github.com/ajhager/eng/spine"
	"math"
)

var (
	skeleton *spine.Skeleton
	walk     *spine.Animation
	jump     *spine.Animation
	batch    *eng.Batch
	animTime float32
	dir      float32
	ramp     float32
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
	jump = skeleton.Animation("jump")
	dir = 1
}

func (g *Game) Update(dt float32) {
	animTime += dt
	skeleton.X += dir * (200 * (1 - ramp)) * dt
	if skeleton.X < 100 {
		skeleton.FlipX = false
		dir = -dir
	}
	if skeleton.X > 924 {
		skeleton.FlipX = true
		dir = -dir
	}
	skeleton.Apply(walk, animTime, true)
	if eng.KeyPressed(eng.Space) && jump.Duration() > animTime {
		dir = float32(math.Copysign(2, float64(dir)))
		ramp += dt
	} else {
		dir = float32(math.Copysign(1, float64(dir)))
		ramp -= dt
	}
	if ramp > 1 {
		ramp = 1
	}
	if ramp < 0 {
		ramp = 0
	}
	skeleton.Mix(jump, animTime, false, ramp)
	skeleton.UpdateWorldTransform()
}

func (g *Game) Draw() {
	batch.Begin()
	skeleton.Draw(batch)
	batch.End()
}

func (g *Game) KeyDown(k eng.Key) {
	if k == eng.Space {
		animTime = 0
		ramp = 0
	}
}

func main() {
	eng.Run(new(Game))
}
