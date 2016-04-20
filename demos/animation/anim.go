package main

import (
	"image/color"

	"engo.io/ecs"
	"engo.io/engo"
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

type DefaultScene struct{}

type Animation struct {
	ecs.BasicEntity
	engo.AnimationComponent
	engo.RenderComponent
	engo.SpaceComponent
}

func (*DefaultScene) Preload() {
	engo.Files.Add("assets/hero.png")
	StopAction = &engo.AnimationAction{Name: "stop", Frames: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}}
	RunAction = &engo.AnimationAction{Name: "run", Frames: []int{16, 17, 18, 19, 20, 21}}
	WalkAction = &engo.AnimationAction{Name: "move", Frames: []int{11, 12, 13, 14, 15}}
	SkillAction = &engo.AnimationAction{Name: "skill", Frames: []int{44, 45, 46, 47, 48, 49, 50, 51, 52, 53}}
	DieAction = &engo.AnimationAction{Name: "die", Frames: []int{28, 29, 30}}
	actions = []*engo.AnimationAction{DieAction, StopAction, WalkAction, RunAction, SkillAction}
}

func (scene *DefaultScene) Setup(w *ecs.World) {
	engo.SetBackground(color.White)

	w.AddSystem(&engo.RenderSystem{})
	w.AddSystem(&engo.AnimationSystem{})
	w.AddSystem(&ControlSystem{})
	w.AddSystem(&engo.MouseZoomer{zoomSpeed})

	spriteSheet := engo.NewSpritesheetFromFile("hero.png", 150, 150)

	hero := scene.CreateEntity(&engo.Point{0, 0}, spriteSheet, StopAction)

	// Add our hero to the appropriate systems
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *engo.RenderSystem:
			sys.Add(&hero.BasicEntity, &hero.RenderComponent, &hero.SpaceComponent)
		case *engo.AnimationSystem:
			sys.Add(&hero.BasicEntity, &hero.AnimationComponent, &hero.RenderComponent)
		case *ControlSystem:
			sys.Add(&hero.BasicEntity, &hero.AnimationComponent)
		}
	}
}

func (*DefaultScene) Hide()        {}
func (*DefaultScene) Show()        {}
func (*DefaultScene) Exit()        {}
func (*DefaultScene) Type() string { return "GameWorld" }

func (*DefaultScene) CreateEntity(point *engo.Point, spriteSheet *engo.Spritesheet, action *engo.AnimationAction) *Animation {
	entity := &Animation{BasicEntity: ecs.NewBasic()}

	entity.SpaceComponent = engo.SpaceComponent{*point, 150, 150}
	entity.RenderComponent = engo.NewRenderComponent(spriteSheet.Cell(action.Frames[0]), engo.Point{3, 3}, "hero")
	entity.AnimationComponent = engo.NewAnimationComponent(spriteSheet.Drawables(), 0.1)
	entity.AnimationComponent.AddAnimationActions(actions)
	entity.AnimationComponent.SelectAnimationByAction(action)

	return entity
}

type controlEntity struct {
	*ecs.BasicEntity
	*engo.AnimationComponent
}

type ControlSystem struct {
	entities []controlEntity
}

func (c *ControlSystem) Add(basic *ecs.BasicEntity, anim *engo.AnimationComponent) {
	c.entities = append(c.entities, controlEntity{basic, anim})
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
		if engo.Keys.Get(engo.ArrowRight).Down() {
			e.AnimationComponent.SelectAnimationByAction(WalkAction)
		} else if engo.Keys.Get(engo.Space).Down() {
			e.AnimationComponent.SelectAnimationByAction(SkillAction)
		} else {
			e.AnimationComponent.SelectAnimationByAction(StopAction)
		}
	}
}

func main() {
	opts := engo.RunOptions{
		Title:  "Animation Demo",
		Width:  1024,
		Height: 640,
	}
	engo.Run(opts, &DefaultScene{})
}
