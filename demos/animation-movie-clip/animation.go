//+build demo

package main

import (
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/EngoEngine/engo/format/mc"
)

const (
	Height = 640
	Width  = Height * 1.6
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	opts := engo.RunOptions{
		Title: "Animation Demo",

		FPSLimit: 30,
		Width:    Width,
		Height:   Height,

		GlobalScale: engo.Point{X: 1.5, Y: 1.5},

		StandardInputs: true,
	}
	engo.Run(opts, &DefaultScene{})
}

func NewHeroEntity(position engo.Point, mcr *mc.MovieClipResource) *HeroEntity {
	entity := &HeroEntity{BasicEntity: ecs.NewBasic()}

	entity.RenderComponent = common.RenderComponent{}
	entity.RenderComponent.Drawable = mcr.Drawable
	entity.RenderComponent.Scale = engo.GetGlobalScale()

	entity.SpaceComponent = common.SpaceComponent{}
	entity.SpaceComponent.Position = position
	entity.SpaceComponent.Width = mcr.Drawable.Width()
	entity.SpaceComponent.Height = mcr.Drawable.Height()

	entity.AnimationComponent = common.NewAnimationComponent(mcr.SpriteSheet.Drawables(), 0.0)
	entity.AnimationComponent.AddAnimations(mcr.Actions)
	entity.AnimationComponent.AddDefaultAnimation(mcr.DefaultAction)

	return entity
}

type HeroEntity struct {
	ecs.BasicEntity
	common.AnimationComponent
	common.RenderComponent
	common.SpaceComponent
}

type DefaultScene struct{}

func (*DefaultScene) Preload() {
	err := engo.Files.Load(
		"black.png",
		"sheep.mc.json",
	)
	if err != nil {
		log.Fatalln(err)
	}
}

func (scene *DefaultScene) Setup(u engo.Updater) {
	w, _ := u.(*ecs.World)

	common.SetBackground(color.Alpha16{A: 0x7575})

	w.AddSystemInterface(&common.RenderSystem{}, new(common.Renderable), nil)
	w.AddSystemInterface(&common.AnimationSystem{}, new(common.Animationable), nil)
	w.AddSystemInterface(&ControlSystem{}, new(Controllable), nil)

	texture, err := common.LoadedSprite("black.png")
	if err != nil {
		log.Fatalln(err)
	}
	w.AddEntity(NewFieldEntity(
		engo.Point{200, 250},
		400,
		texture,
	))

	mcr, err := mc.LoadResource("sheep.mc.json")
	if err != nil {
		log.Fatalln(err)
	}
	hero := NewHeroEntity(engo.Point{200, 50}, mcr)
	w.AddEntity(hero)
}

func (*DefaultScene) Type() string { return "GameWorld" }

type controlEntity struct {
	*ecs.BasicEntity
	*common.AnimationComponent
}

type Controllable interface {
	common.BasicFace
	common.AnimationFace
}

type ControlSystem struct {
	entities map[uint64]controlEntity
}

func (c *ControlSystem) AddByInterface(i ecs.Identifier) {
	o, _ := i.(Controllable)
	c.Add(o.GetBasicEntity(), o.GetAnimationComponent())
}

func (c *ControlSystem) Add(basic *ecs.BasicEntity, anim *common.AnimationComponent) {
	if c.entities == nil {
		c.entities = make(map[uint64]controlEntity)
	}
	c.entities[basic.ID()] = controlEntity{basic, anim}
}

func (c *ControlSystem) Remove(basic ecs.BasicEntity) {
	if c.entities != nil {
		delete(c.entities, basic.ID())
	}
}

func (c *ControlSystem) Update(dt float32) {
	for _, e := range c.entities {
		if engo.Input.Button("action").JustPressed() {
			c.randAction(e.GetAnimationComponent())
		}
	}
}

func (c *ControlSystem) randAction(anim *common.AnimationComponent) {
	animCount := len(anim.Animations)
	list := make([]string, 0, animCount)
	for name := range anim.Animations {
		list = append(list, name)
	}
	anim.SelectAnimationByName(list[rand.Intn(animCount)])
}

type FieldEntity struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

func NewFieldEntity(position engo.Point, size float32, texture *common.Texture) *FieldEntity {
	entity := &FieldEntity{BasicEntity: ecs.NewBasic()}

	entity.RenderComponent = common.RenderComponent{Drawable: texture}
	entity.RenderComponent.Repeat = common.Repeat
	entity.RenderComponent.Scale = engo.GetGlobalScale()

	entity.SpaceComponent = common.SpaceComponent{Position: position}
	entity.SpaceComponent.Width = size
	entity.SpaceComponent.Height = size

	return entity
}
