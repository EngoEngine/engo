//+build demo

package main

import (
	"image/color"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
)

var (
	upButton    = "up"
	downButton  = "down"
	leftButton  = "left"
	rightButton = "right"

	levelData *common.Level
)

const (
	SPEED_MESSAGE = "SpeedMessage"
	SPEED_SCALE   = 64
)

type DefaultScene struct{}

type Car struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
	ControlComponent
	SpeedComponent
}

type ControlComponent struct {
	SchemeVert  string
	SchemeHoriz string
	TileLayer   *common.TileLayer
}

type Tile struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

func (*DefaultScene) Preload() {
	// Load character model
	engo.Files.Load("orange_vehicles.png")

	// Load TileMap
	if err := engo.Files.Load("map.tmx"); err != nil {
		panic(err)
	}

	engo.Input.RegisterButton(upButton, engo.KeyW, engo.KeyArrowUp)
	engo.Input.RegisterButton(leftButton, engo.KeyA, engo.KeyArrowLeft)
	engo.Input.RegisterButton(rightButton, engo.KeyD, engo.KeyArrowRight)
	engo.Input.RegisterButton(downButton, engo.KeyS, engo.KeyArrowDown)
}

func (scene *DefaultScene) Setup(u engo.Updater) {
	w, _ := u.(*ecs.World)

	common.SetBackground(color.White)

	w.AddSystem(&common.RenderSystem{})
	w.AddSystem(&SpeedSystem{})
	w.AddSystem(&ControlSystem{})

	// Setup TileMap
	resource, err := engo.Files.Resource("map.tmx")
	if err != nil {
		panic(err)
	}
	tmxResource := resource.(common.TMXResource)
	levelData = tmxResource.Level

	// Create Hero
	spriteSheet := common.NewSpritesheetFromFile("orange_vehicles.png", 32, 32)

	car := &Car{BasicEntity: ecs.NewBasic()}

	car.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{X: 0, Y: 0},
		Width:    float32(32) * 4,
		Height:   float32(32) * 4,
	}
	car.RenderComponent = common.RenderComponent{
		Drawable: spriteSheet.Cell(0),
		Scale:    engo.Point{4, 4},
	}

	car.SpeedComponent = SpeedComponent{}
	car.ControlComponent = ControlComponent{
		SchemeHoriz: "horizontal",
		SchemeVert:  "vertical",
		TileLayer:   levelData.TileLayers[0],
	}

	car.RenderComponent.SetZIndex(1)

	// Add our hero to the appropriate systems
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(
				&car.BasicEntity,
				&car.RenderComponent,
				&car.SpaceComponent,
			)

		case *ControlSystem:
			sys.Add(
				&car.BasicEntity,
				&car.ControlComponent,
				&car.SpaceComponent,
			)

		case *SpeedSystem:
			sys.Add(
				&car.BasicEntity,
				&car.SpeedComponent,
				&car.SpaceComponent,
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
				tile.SetZIndex(0)
				tile.SpaceComponent = common.SpaceComponent{
					Position: tileElement.Point,
					Width:    0,
					Height:   0,
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

	// Setup character and movement
	engo.Input.RegisterAxis(
		"vertical",
		engo.AxisKeyPair{engo.KeyArrowUp, engo.KeyArrowDown},
		engo.AxisKeyPair{engo.KeyW, engo.KeyS},
	)

	engo.Input.RegisterAxis(
		"horizontal",
		engo.AxisKeyPair{engo.KeyArrowLeft, engo.KeyArrowRight},
		engo.AxisKeyPair{engo.KeyA, engo.KeyD},
	)

	// Add EntityScroller System
	w.AddSystem(&common.EntityScroller{
		SpaceComponent: &car.SpaceComponent,
		TrackingBounds: levelData.Bounds(),
	})
}

func (*DefaultScene) Type() string { return "DefaultScene" }

type SpeedMessage struct {
	*ecs.BasicEntity
	engo.Point
}

func (SpeedMessage) Type() string {
	return SPEED_MESSAGE
}

type SpeedComponent struct {
	engo.Point
}

type speedEntity struct {
	*ecs.BasicEntity
	*SpeedComponent
	*common.SpaceComponent
}

type SpeedSystem struct {
	entities []speedEntity
}

func (s *SpeedSystem) New(*ecs.World) {
	engo.Mailbox.Listen(SPEED_MESSAGE, func(message engo.Message) {
		speed, isSpeed := message.(SpeedMessage)
		if isSpeed {
			for _, e := range s.entities {
				if e.ID() == speed.BasicEntity.ID() {
					e.SpeedComponent.Point = speed.Point
				}
			}
		}
	})
}

func (s *SpeedSystem) Add(basic *ecs.BasicEntity, speed *SpeedComponent, space *common.SpaceComponent) {
	s.entities = append(s.entities, speedEntity{basic, speed, space})
}

func (s *SpeedSystem) Remove(basic ecs.BasicEntity) {
	delete := -1
	for index, e := range s.entities {
		if e.BasicEntity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		s.entities = append(s.entities[:delete], s.entities[delete+1:]...)
	}
}

func (s *SpeedSystem) Update(dt float32) {
	for _, e := range s.entities {
		speed := engo.GameWidth() * dt
		prev := e.SpaceComponent.Position
		e.SpaceComponent.Position.X = e.SpaceComponent.Position.X + speed*e.SpeedComponent.Point.X
		e.SpaceComponent.Position.Y = e.SpaceComponent.Position.Y + speed*e.SpeedComponent.Point.Y

		var t *common.Tile
		if e.SpaceComponent.Position.X >= 0 {
			t = levelData.GetTile(e.SpaceComponent.Center())
		} else {
			t = levelData.GetTile(engo.Point{
				X: e.SpaceComponent.Position.X - float32(levelData.TileWidth),
				Y: e.SpaceComponent.Position.Y + float32(levelData.TileHeight),
			})
		}

		if t == nil {
			e.SpaceComponent.Position = prev
		}
	}
}

type controlEntity struct {
	*ecs.BasicEntity
	*ControlComponent
	*common.SpaceComponent
}

type ControlSystem struct {
	entities []controlEntity
}

func (c *ControlSystem) Add(basic *ecs.BasicEntity, control *ControlComponent, space *common.SpaceComponent) {
	c.entities = append(c.entities, controlEntity{basic, control, space})
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

func getSpeed(e controlEntity) (p engo.Point, changed bool) {
	p.X = engo.Input.Axis(e.ControlComponent.SchemeHoriz).Value()
	p.Y = engo.Input.Axis(e.ControlComponent.SchemeVert).Value()
	origX, origY := p.X, p.Y

	if engo.Input.Button(upButton).JustPressed() {
		p.Y = -1
	} else if engo.Input.Button(downButton).JustPressed() {
		p.Y = 1
	}
	if engo.Input.Button(leftButton).JustPressed() {
		p.X = -1
	} else if engo.Input.Button(rightButton).JustPressed() {
		p.X = 1
	}

	if engo.Input.Button(upButton).JustReleased() || engo.Input.Button(downButton).JustReleased() {
		p.Y = 0
		changed = true
		if engo.Input.Button(upButton).Down() {
			p.Y = -1
		} else if engo.Input.Button(downButton).Down() {
			p.Y = 1
		} else if engo.Input.Button(leftButton).Down() {
			p.X = -1
		} else if engo.Input.Button(rightButton).Down() {
			p.X = 1
		}
	}
	if engo.Input.Button(leftButton).JustReleased() || engo.Input.Button(rightButton).JustReleased() {
		p.X = 0
		changed = true
		if engo.Input.Button(leftButton).Down() {
			p.X = -1
		} else if engo.Input.Button(rightButton).Down() {
			p.X = 1
		} else if engo.Input.Button(upButton).Down() {
			p.Y = -1
		} else if engo.Input.Button(downButton).Down() {
			p.Y = 1
		}
	}
	changed = changed || p.X != origX || p.Y != origY
	return
}

func (c *ControlSystem) Update(dt float32) {
	for _, e := range c.entities {
		if vector, changed := getSpeed(e); changed {
			speed := dt * SPEED_SCALE
			vector, _ = vector.Normalize()
			vector.MultiplyScalar(speed)
			engo.Mailbox.Dispatch(SpeedMessage{e.BasicEntity, vector})
		}
	}
}

func main() {
	opts := engo.RunOptions{
		Title:  "My Little Isometric Adventure",
		Width:  500,
		Height: 500,
	}
	engo.Run(opts, &DefaultScene{})
}
