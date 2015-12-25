package main

import (
	"github.com/paked/engi"
	"github.com/paked/engi/ecs"
)

var (
	zoomSpeed   float32 = -0.125
	RunAction   *engi.AnimationAction
	WalkAction  *engi.AnimationAction
	StopAction  *engi.AnimationAction
	SkillAction *engi.AnimationAction
	DieAction   *engi.AnimationAction
	actions     []*engi.AnimationAction
)

type GameWorld struct{}

func (game *GameWorld) Preload() {
	engi.Files.Add("assets/hero.png")
	StopAction = &engi.AnimationAction{Name: "stop", Frames: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}}
	RunAction = &engi.AnimationAction{Name: "run", Frames: []int{16, 17, 18, 19, 20, 21}}
	WalkAction = &engi.AnimationAction{Name: "move", Frames: []int{11, 12, 13, 14, 15}}
	SkillAction = &engi.AnimationAction{Name: "skill", Frames: []int{44, 45, 46, 47, 48, 49, 50, 51, 52, 53}}
	DieAction = &engi.AnimationAction{Name: "die", Frames: []int{28, 29, 30}}
	actions = []*engi.AnimationAction{DieAction, StopAction, WalkAction, RunAction, SkillAction}
}

func (game *GameWorld) Setup(w *ecs.World) {
	engi.SetBg(0xFFFFFF)

	w.AddSystem(&engi.RenderSystem{})
	w.AddSystem(&engi.AnimationSystem{})
	w.AddSystem(&ControlSystem{})
	w.AddSystem(engi.NewMouseZoomer(zoomSpeed))

	spriteSheet := engi.NewSpritesheetFromFile("hero.png", 150, 150)

	w.AddEntity(game.CreateEntity(&engi.Point{0, 0}, spriteSheet, StopAction))
}

func (*GameWorld) Hide()        {}
func (*GameWorld) Show()        {}
func (*GameWorld) Type() string { return "GameWorld" }

func (game *GameWorld) CreateEntity(point *engi.Point, spriteSheet *engi.Spritesheet, action *engi.AnimationAction) *ecs.Entity {
	entity := ecs.NewEntity([]string{"AnimationSystem", "RenderSystem", "ControlSystem"})

	space := &engi.SpaceComponent{*point, 150, 150}
	render := engi.NewRenderComponent(spriteSheet.Cell(action.Frames[0]), engi.Point{3, 3}, "hero")
	animation := engi.NewAnimationComponent(spriteSheet.Drawables(), 0.1)
	animation.AddAnimationActions(actions)
	animation.SelectAnimationByAction(action)
	entity.AddComponent(render)
	entity.AddComponent(space)
	entity.AddComponent(animation)

	return entity
}

type ControlSystem struct {
	*ecs.System
}

func (ControlSystem) Type() string {
	return "ControlSystem"
}

func (c *ControlSystem) New(*ecs.World) {
	c.System = ecs.NewSystem()
}

func (c *ControlSystem) Update(entity *ecs.Entity, dt float32) {
	var a *engi.AnimationComponent

	if !entity.Component(&a) {
		return
	}

	if engi.Keys.Get(engi.ArrowRight).Down() {
		a.SelectAnimationByAction(WalkAction)
	} else if engi.Keys.Get(engi.Space).Down() {
		a.SelectAnimationByAction(SkillAction)
	} else {
		a.SelectAnimationByAction(StopAction)
	}
}

func main() {
	opts := engi.RunOptions{
		Title:  "Animation Demo",
		Width:  1024,
		Height: 640,
	}
	engi.Open(opts, &GameWorld{})
}
