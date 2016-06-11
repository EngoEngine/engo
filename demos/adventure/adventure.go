package main

import (
	"image/color"
	"log"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
)

var (
	WalkUpAction    *common.Animation
	WalkDownAction  *common.Animation
	WalkLeftAction  *common.Animation
	WalkRightAction *common.Animation
	StopUpAction    *common.Animation
	StopDownAction  *common.Animation
	StopLeftAction  *common.Animation
	StopRightAction *common.Animation
	SkillAction     *common.Animation
	actions         []*common.Animation

	upButton    = "up"
	downButton  = "down"
	leftButton  = "left"
	rightButton = "right"
	model       = "motw.png"
	width       = 52
	height      = 73
	levelWidth  float32
	levelHeight float32
)

type DefaultScene struct{}

type Hero struct {
	ecs.BasicEntity
	common.AnimationComponent
	common.RenderComponent
	common.SpaceComponent
	ControlComponent
}

type ControlComponent struct {
	SchemeVert  string
	SchemeHoriz string
}

type Tile struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
	common.CollisionComponent
}

func (*DefaultScene) Preload() {

	// Load character model
	engo.Files.Load(model)

	// Load TileMap
	if err := engo.Files.Load("example.tmx"); err != nil {
		panic(err)
	}

	StopUpAction = &common.Animation{
		Name:   "upstop",
		Frames: []int{37},
	}

	StopDownAction = &common.Animation{
		Name:   "downstop",
		Frames: []int{1},
	}

	StopLeftAction = &common.Animation{
		Name:   "leftstop",
		Frames: []int{13},
	}

	StopRightAction = &common.Animation{
		Name:   "rightstop",
		Frames: []int{25},
	}

	WalkUpAction = &common.Animation{
		Name:   "up",
		Frames: []int{36, 37, 38},
		Loop:   true,
	}

	WalkDownAction = &common.Animation{
		Name:   "down",
		Frames: []int{0, 1, 2},
		Loop:   true,
	}

	WalkLeftAction = &common.Animation{
		Name:   "left",
		Frames: []int{12, 13, 14},
		Loop:   true,
	}

	WalkRightAction = &common.Animation{
		Name:   "right",
		Frames: []int{24, 25, 26},
		Loop:   true,
	}

	actions = []*common.Animation{
		StopUpAction,
		StopDownAction,
		StopLeftAction,
		StopRightAction,
		WalkUpAction,
		WalkDownAction,
		WalkLeftAction,
		WalkRightAction,
	}

	engo.Input.RegisterButton(upButton, engo.W, engo.ArrowUp)
	engo.Input.RegisterButton(leftButton, engo.A, engo.ArrowLeft)
	engo.Input.RegisterButton(rightButton, engo.D, engo.ArrowRight)
	engo.Input.RegisterButton(downButton, engo.S, engo.ArrowDown)
}

func (scene *DefaultScene) Setup(w *ecs.World) {
	common.SetBackground(color.White)

	w.AddSystem(&common.RenderSystem{})
	w.AddSystem(&common.AnimationSystem{})
	w.AddSystem(&ControlSystem{})

	// Setup TileMap
	resource, err := engo.Files.Resource("example.tmx")
	if err != nil {
		panic(err)
	}
	tmxResource := resource.(common.TMXResource)
	levelData := tmxResource.Level

	// Extract Map Size
	levelWidth = levelData.Bounds().Max.X
	levelHeight = levelData.Bounds().Max.Y

	// Create Hero
	spriteSheet := common.NewSpritesheetFromFile(model, width, height)

	hero := scene.CreateHero(
		engo.Point{engo.CanvasWidth() / 2, engo.CanvasHeight() / 2},
		spriteSheet,
	)

	hero.ControlComponent = ControlComponent{
		SchemeHoriz: "horizontal",
		SchemeVert:  "vertical",
	}

	hero.RenderComponent.SetZIndex(1)

	// Add our hero to the appropriate systems
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(
				&hero.BasicEntity,
				&hero.RenderComponent,
				&hero.SpaceComponent,
			)

		case *common.AnimationSystem:
			sys.Add(
				&hero.BasicEntity,
				&hero.AnimationComponent,
				&hero.RenderComponent,
			)

		case *ControlSystem:
			sys.Add(
				&hero.BasicEntity,
				&hero.AnimationComponent,
				&hero.ControlComponent,
				&hero.SpaceComponent,
			)
		}
	}

	// Create render and space components for each of the tiles
	tileComponents := make([]*Tile, 0)

	for _, tileLayer := range levelData.TileLayers {
		for _, tileElement := range tileLayer.Tiles {

			if tileElement.Image != nil {
				tile := &Tile{BasicEntity: ecs.NewBasic()}
				tile.RenderComponent = common.RenderComponent{
					Drawable: tileElement,
					Scale:    engo.Point{1, 1},
				}
				tile.SpaceComponent = common.SpaceComponent{
					Position: tileElement.Point,
					Width:    0,
					Height:   0,
				}

				if tileLayer.Name == "grass" {
					tile.RenderComponent.SetZIndex(0)
				}

				if tileLayer.Name == "trees" {
					tile.RenderComponent.SetZIndex(2)
				}

				tileComponents = append(tileComponents, tile)
			}
		}
	}

	for _, imageLayer := range levelData.ImageLayers {
		for _, imageElement := range imageLayer.Images {

			if imageElement.Image != nil {
				tile := &Tile{BasicEntity: ecs.NewBasic()}
				tile.RenderComponent = common.RenderComponent{
					Drawable: imageElement,
					Scale:    engo.Point{1, 1},
				}
				tile.SpaceComponent = common.SpaceComponent{
					Position: imageElement.Point,
					Width:    0,
					Height:   0,
				}

				if imageLayer.Name == "clouds" {
					tile.RenderComponent.SetZIndex(3)
				}

				tileComponents = append(tileComponents, tile)
			}
		}
	}

	// Add each of the tiles entities and its components to the render system
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			for _, v := range tileComponents {
				sys.Add(&v.BasicEntity, &v.RenderComponent, &v.SpaceComponent)
			}

		}
	}

	// Access Object Layers
	for _, objectLayer := range levelData.ObjectLayers {
		log.Println("This object layer is called " + objectLayer.Name)
		// Do something with every regular Object
		for _, object := range objectLayer.Objects {
			log.Println("This object is called " + object.Name)
		}

		// Do something with every polyline Object
		for _, polylineObject := range objectLayer.PolyObjects {
			log.Println("This object is called " + polylineObject.Name)
		}
	}

	// Setup character and movement
	engo.Input.RegisterAxis(
		"vertical",
		engo.AxisKeyPair{engo.ArrowUp, engo.ArrowDown},
	)

	engo.Input.RegisterAxis(
		"horizontal",
		engo.AxisKeyPair{engo.ArrowLeft, engo.ArrowRight},
	)

	// Add EntityScroller System
	w.AddSystem(&common.EntityScroller{
		SpaceComponent: &hero.SpaceComponent,
		TrackingBounds: levelData.Bounds(),
	})
}

func (*DefaultScene) Type() string { return "DefaultScene" }

func (*DefaultScene) CreateHero(point engo.Point, spriteSheet *common.Spritesheet) *Hero {
	hero := &Hero{BasicEntity: ecs.NewBasic()}

	hero.SpaceComponent = common.SpaceComponent{
		Position: point,
		Width:    float32(width),
		Height:   float32(height),
	}
	hero.RenderComponent = common.RenderComponent{
		Drawable: spriteSheet.Cell(0),
		Scale:    engo.Point{1, 1},
	}
	hero.AnimationComponent = common.NewAnimationComponent(spriteSheet.Drawables(), 0.1)

	hero.AnimationComponent.AddAnimations(actions)
	hero.AnimationComponent.SelectAnimationByName("downstop")

	return hero
}

type controlEntity struct {
	*ecs.BasicEntity
	*common.AnimationComponent
	*ControlComponent
	*common.SpaceComponent
}

type ControlSystem struct {
	entities []controlEntity
}

func (c *ControlSystem) Add(basic *ecs.BasicEntity, anim *common.AnimationComponent, control *ControlComponent, space *common.SpaceComponent) {
	c.entities = append(c.entities, controlEntity{basic, anim, control, space})
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

		// Add Character Movement Control
		if engo.Input.Button(upButton).JustPressed() {
			e.AnimationComponent.SelectAnimationByAction(WalkUpAction)
		} else if engo.Input.Button(downButton).JustPressed() {
			e.AnimationComponent.SelectAnimationByAction(WalkDownAction)
		} else if engo.Input.Button(leftButton).JustPressed() {
			e.AnimationComponent.SelectAnimationByAction(WalkLeftAction)
		} else if engo.Input.Button(rightButton).JustPressed() {
			e.AnimationComponent.SelectAnimationByAction(WalkRightAction)
		}

		if engo.Input.Button(upButton).JustReleased() {
			e.AnimationComponent.SelectAnimationByAction(StopUpAction)
		} else if engo.Input.Button(downButton).JustReleased() {
			e.AnimationComponent.SelectAnimationByAction(StopDownAction)
		} else if engo.Input.Button(leftButton).JustReleased() {
			e.AnimationComponent.SelectAnimationByAction(StopLeftAction)
		} else if engo.Input.Button(rightButton).JustReleased() {
			e.AnimationComponent.SelectAnimationByAction(StopRightAction)
		}

		speed := engo.GameWidth() * dt

		vert := engo.Input.Axis(e.ControlComponent.SchemeVert)
		e.SpaceComponent.Position.Y += speed * vert.Value()

		horiz := engo.Input.Axis(e.ControlComponent.SchemeHoriz)
		e.SpaceComponent.Position.X += speed * horiz.Value()

		// Add Game Border Limits
		var heightLimit float32 = levelHeight - e.SpaceComponent.Height
		if e.SpaceComponent.Position.Y < 0 {
			e.SpaceComponent.Position.Y = 0
		} else if e.SpaceComponent.Position.Y > heightLimit {
			e.SpaceComponent.Position.Y = heightLimit
		}

		var widthLimit float32 = levelWidth - e.SpaceComponent.Width
		if e.SpaceComponent.Position.X < 0 {
			e.SpaceComponent.Position.X = 0
		} else if e.SpaceComponent.Position.X > widthLimit {
			e.SpaceComponent.Position.X = widthLimit
		}
	}
}

func main() {
	opts := engo.RunOptions{
		Title:  "My Little Adventure",
		Width:  500,
		Height: 500,
	}
	engo.Run(opts, &DefaultScene{})
}
