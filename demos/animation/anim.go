package main

import (
	"image/color"
	"log"

	"engo.io/engo"
	"engo.io/ecs"
)

var (
	zoomSpeed   float32 = -0.125
	RunAction   *engo.AnimationAction
	WalkAction  *engo.AnimationAction
	StopAction  *engo.AnimationAction
	SkillAction *engo.AnimationAction
	DieAction   *engo.AnimationAction
	actions     []*engo.AnimationAction
)

type GameWorld struct{}

func (game *GameWorld) Preload() {
	engo.Files.Add("assets/hero.png")
	StopAction = &engo.AnimationAction{Name: "stop", Frames: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}}
	RunAction = &engo.AnimationAction{Name: "run", Frames: []int{16, 17, 18, 19, 20, 21}}
	WalkAction = &engo.AnimationAction{Name: "move", Frames: []int{11, 12, 13, 14, 15}}
	SkillAction = &engo.AnimationAction{Name: "skill", Frames: []int{44, 45, 46, 47, 48, 49, 50, 51, 52, 53}}
	DieAction = &engo.AnimationAction{Name: "die", Frames: []int{28, 29, 30}}
	actions = []*engo.AnimationAction{DieAction, StopAction, WalkAction, RunAction, SkillAction}
}

func (game *GameWorld) Setup(w *ecs.World) {
	engo.SetBackground(color.White)

	w.AddSystem(&engo.RenderSystem{})
	w.AddSystem(&engo.AnimationSystem{})
	w.AddSystem(&ControlSystem{})
	w.AddSystem(&engo.MouseZoomer{zoomSpeed})

	spriteSheet := engo.NewSpritesheetFromFile("hero.png", 150, 150)

	err := w.AddEntity(game.CreateEntity(&engo.Point{0, 0}, spriteSheet, StopAction))
	if err != nil {
		log.Println(err)
	}
}

func (*GameWorld) Hide()        {}
func (*GameWorld) Show()        {}
func (*GameWorld) Exit() 	 	{}
func (*GameWorld) Type() string { return "GameWorld" }

func (game *GameWorld) CreateEntity(point *engo.Point, spriteSheet *engo.Spritesheet, action *engo.AnimationAction) *ecs.Entity {
	entity := ecs.NewEntity("AnimationSystem", "RenderSystem", "ControlSystem")

	space := &engo.SpaceComponent{*point, 150, 150}
	render := engo.NewRenderComponent(spriteSheet.Cell(action.Frames[0]), engo.Point{3, 3}, "hero")
	animation := engo.NewAnimationComponent(spriteSheet.Drawables(), 0.1)
	animation.AddAnimationActions(actions)
	animation.SelectAnimationByAction(action)
	entity.AddComponent(render)
	entity.AddComponent(space)
	entity.AddComponent(animation)

	return entity
}

type ControlSystem struct {
	ecs.LinearSystem
}

func (*ControlSystem) Type() string { return "ControlSystem" }
func (*ControlSystem) Pre()         {}
func (*ControlSystem) Post()        {}

func (c *ControlSystem) New(*ecs.World) {}

func (c *ControlSystem) UpdateEntity(entity *ecs.Entity, dt float32) {
	var a *engo.AnimationComponent

	if !entity.Component(&a) {
		return
	}

	if engo.Keys.Get(engo.ArrowRight).Down() {
		a.SelectAnimationByAction(WalkAction)
	} else if engo.Keys.Get(engo.Space).Down() {
		a.SelectAnimationByAction(SkillAction)
	} else {
		a.SelectAnimationByAction(StopAction)
	}
}

func main() {
	opts := engo.RunOptions{
		Title:  "Animation Demo",
		Width:  1024,
		Height: 640,
		DefaultCloseAction: true,
	}
	engo.Run(opts, &GameWorld{})
}
