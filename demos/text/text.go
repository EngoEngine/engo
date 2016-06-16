package main

import (
	"image/color"
	"log"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/act"
	"engo.io/engo/common"
)

type DefaultScene struct{}

var (
	zoomSpeed   float32 = -0.125
	scrollSpeed float32 = 700

	worldWidth  int = 800
	worldHeight int = 800
)

type MyLabel struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

func (*DefaultScene) Preload() {
	err := engo.Files.Load("Roboto-Regular.ttf")
	if err != nil {
		panic(err)
	}
}

// Setup is called before the main loop is started
func (*DefaultScene) Setup(w *ecs.World) {
	common.SetBackground(color.White)
	w.AddSystem(&common.RenderSystem{})

	// Get the default axis id's
	vertAxis := engo.Axes.Id(engo.DefaultVerticalAxis)
	horiAxis := engo.Axes.Id(engo.DefaultHorizontalAxis)
	if 0 == vertAxis {
		log.Println("Default vertical axis not found, setting up fall back!")
		vertAxis = engo.Axes.SetNamed(engo.DefaultVerticalAxis, act.AxisPair{act.KeyW, act.KeyS}, act.AxisPair{act.KeyUp, act.KeyDown})
	}
	if 0 == horiAxis {
		log.Println("Default horizontal axis not found, setting up fall back!")
		horiAxis = engo.Axes.SetNamed(engo.DefaultHorizontalAxis, act.AxisPair{act.KeyA, act.KeyD}, act.AxisPair{act.KeyLeft, act.KeyRight})
	}

	// Adding KeyboardScroller so we can actually see the difference between the HUD and non-HUD text
	w.AddSystem(common.NewKeyboardScroller(scrollSpeed, horiAxis, vertAxis))
	w.AddSystem(&common.MouseZoomer{zoomSpeed})

	fnt := &common.Font{
		URL:  "Roboto-Regular.ttf",
		FG:   color.Black,
		Size: 64,
	}
	err := fnt.CreatePreloaded()
	if err != nil {
		panic(err)
	}

	label1 := MyLabel{BasicEntity: ecs.NewBasic()}
	label1.RenderComponent.Drawable = common.Text{
		Font: fnt,
		Text: "Hello world !",
	}
	label1.SetShader(common.HUDShader)

	label2 := MyLabel{BasicEntity: ecs.NewBasic()}
	label2.RenderComponent.Drawable = common.Text{
		Font:          fnt,
		Text:          "This may also be text\nwhich includes a newline. ",
		LineSpacing:   0.5,
		LetterSpacing: 0.15,
	}

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&label1.BasicEntity, &label1.RenderComponent, &label1.SpaceComponent)
			sys.Add(&label2.BasicEntity, &label2.RenderComponent, &label2.SpaceComponent)
		}
	}
}

func (*DefaultScene) Type() string { return "Game" }

func main() {
	opts := engo.RunOptions{
		Title:          "Text Demo",
		Width:          worldWidth,
		Height:         worldHeight,
		StandardInputs: true,
	}
	engo.Run(opts, &DefaultScene{})
}
