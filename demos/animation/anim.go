package main

import (
	"image/color"

	"engo.io/ecs"
	"engo.io/engo"
)

var (
	WalkAction  *engo.Animation
	StopAction  *engo.Animation
	SkillAction *engo.Animation
	actions     []*engo.Animation
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
	StopAction = &engo.Animation{Name: "stop", Frames: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}}
	WalkAction = &engo.Animation{Name: "move", Frames: []int{11, 12, 13, 14, 15}, Loop: true}
	SkillAction = &engo.Animation{Name: "skill", Frames: []int{44, 45, 46, 47, 48, 49, 50, 51, 52, 53}}
	actions = []*engo.Animation{StopAction, WalkAction, SkillAction}
}

func (scene *DefaultScene) Setup(w *ecs.World) {
	engo.SetBackground(color.White)

	w.AddSystem(&engo.RenderSystem{})
	w.AddSystem(&engo.AnimationSystem{})
	w.AddSystem(&ControlSystem{})

	spriteSheet := engo.NewSpritesheetFromFile("hero.png", 150, 150)

	hero := scene.CreateEntity(&engo.Point{0, 0}, spriteSheet)

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

func (*DefaultScene) Type() string { return "GameWorld" }

func (*DefaultScene) CreateEntity(point *engo.Point, spriteSheet *engo.Spritesheet) *Animation {
	entity := &Animation{BasicEntity: ecs.NewBasic()}

	entity.SpaceComponent = engo.SpaceComponent{*point, 150, 150}
	entity.RenderComponent = engo.NewRenderComponent(spriteSheet.Cell(0), engo.Point{3, 3})
	entity.AnimationComponent = engo.NewAnimationComponent(spriteSheet.Drawables(), 0.1)
	entity.AnimationComponent.AddAnimations(actions)

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
		} /* else {
			e.AnimationComponent.SelectAnimationByAction(StopAction)
		}*/
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
