//+build demo

package main

import (
	"image/color"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
)

var (
	WalkAction  *common.Animation
	StopAction  *common.Animation
	SkillAction *common.Animation
	actions     []*common.Animation

	jumpButton   = "jump"
	actionButton = "action"
)

type DefaultScene struct{}

type Animation struct {
	ecs.BasicEntity
	common.AnimationComponent
	common.RenderComponent
	common.SpaceComponent
}

func (*DefaultScene) Preload() {
	engo.Files.Load("hero.png")

	StopAction = &common.Animation{Name: "stop", Frames: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}}
	WalkAction = &common.Animation{Name: "move", Frames: []int{11, 12, 13, 14, 15}, Loop: true}
	SkillAction = &common.Animation{Name: "skill", Frames: []int{44, 45, 46, 47, 48, 49, 50, 51, 52, 53}}
	actions = []*common.Animation{StopAction, WalkAction, SkillAction}

	engo.Input.RegisterButton(jumpButton, engo.KeySpace, engo.KeyX)
	engo.Input.RegisterButton(actionButton, engo.KeyD, engo.KeyArrowRight)
}

func (scene *DefaultScene) Setup(u engo.Updater) {
	w, _ := u.(*ecs.World)

	common.SetBackground(color.White)

	w.AddSystem(&common.RenderSystem{})
	w.AddSystem(&common.AnimationSystem{})
	w.AddSystem(&ControlSystem{})

	spriteSheet := common.NewSpritesheetFromFile("hero.png", 150, 150)

	hero := scene.CreateEntity(engo.Point{0, 0}, spriteSheet)

	// Add our hero to the appropriate systems
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&hero.BasicEntity, &hero.RenderComponent, &hero.SpaceComponent)
		case *common.AnimationSystem:
			sys.Add(&hero.BasicEntity, &hero.AnimationComponent, &hero.RenderComponent)
		case *ControlSystem:
			sys.Add(&hero.BasicEntity, &hero.AnimationComponent)
		}
	}
}

func (*DefaultScene) Type() string { return "GameWorld" }

func (*DefaultScene) CreateEntity(point engo.Point, spriteSheet *common.Spritesheet) *Animation {
	entity := &Animation{BasicEntity: ecs.NewBasic()}

	entity.SpaceComponent = common.SpaceComponent{
		Position: point,
		Width:    150,
		Height:   150,
	}
	entity.RenderComponent = common.RenderComponent{
		Drawable: spriteSheet.Cell(0),
		Scale:    engo.Point{3, 3},
	}
	entity.AnimationComponent = common.NewAnimationComponent(spriteSheet.Drawables(), 0.1)

	entity.AnimationComponent.AddAnimations(actions)
	entity.AnimationComponent.AddDefaultAnimation(StopAction)

	return entity
}

type controlEntity struct {
	*ecs.BasicEntity
	*common.AnimationComponent
}

type ControlSystem struct {
	entities []controlEntity
}

func (c *ControlSystem) Add(basic *ecs.BasicEntity, anim *common.AnimationComponent) {
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
		if engo.Input.Button(actionButton).JustPressed() {
			e.AnimationComponent.SelectAnimationByAction(WalkAction)
		} else if engo.Input.Button(jumpButton).JustPressed() {
			e.AnimationComponent.SelectAnimationByAction(SkillAction)
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
