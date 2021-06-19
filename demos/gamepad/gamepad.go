package main

import (
	"bytes"
	"image/color"
	"strconv"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"

	"golang.org/x/image/font/gofont/gomonobold"
)

var font *common.Font

type DefaultScene struct{}

type panel struct {
	ecs.BasicEntity

	common.RenderComponent
	common.SpaceComponent

	InputComponent
}

type label struct {
	ecs.BasicEntity

	common.RenderComponent
	common.SpaceComponent
}

type GamepadInputType uint8

const (
	GamepadInputA GamepadInputType = iota
	GamepadInputB
	GamepadInputX
	GamepadInputY
	GamepadInputStart
	GamepadInputBack
	GamepadInputGuide
	GamepadInputLeftBumper
	GamepadInputRightBumper
	GamepadInputLeftThumb
	GamepadInputRightThumb
	GamepadInputDpadUp
	GamepadInputDpadRight
	GamepadInputDpadDown
	GamepadInputDpadLeft
	GamepadInputLeftTrigger
	GamepadInputRightTrigger
	GamepadInputLeftAxisX
	GamepadInputLeftAxisY
	GamepadInputRightAxisX
	GamepadInputRightAxisY
)

type InputComponent struct {
	Trigger GamepadInputType
}

type InputEntity struct {
	*ecs.BasicEntity

	*common.RenderComponent

	*InputComponent
}

type InputSystem struct {
	entities []InputEntity
}

func (s *InputSystem) Add(basic *ecs.BasicEntity, rend *common.RenderComponent, input *InputComponent) {
	s.entities = append(s.entities, InputEntity{basic, rend, input})
}

func (s *InputSystem) Remove(basic ecs.BasicEntity) {
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

func (s *InputSystem) Update(float32) {
	// Retrieve the Gamepad
	gamepad := engo.Input.Gamepad("Player1")
	if gamepad == nil {
		println("No gamepad found for Player1.")
		return
	}
	for _, entity := range s.entities {
		switch entity.Trigger {
		case GamepadInputA:
			if gamepad.A.Up() {
				entity.Color = color.White
			} else if gamepad.A.JustPressed() {
				entity.Color = color.RGBA{0, 255, 0, 255}
			} else if gamepad.A.Down() {
				entity.Color = color.RGBA{255, 0, 0, 255}
			}
		case GamepadInputB:
			if gamepad.B.Up() {
				entity.Color = color.White
			} else if gamepad.B.JustPressed() {
				entity.Color = color.RGBA{0, 255, 0, 255}
			} else if gamepad.B.Down() {
				entity.Color = color.RGBA{255, 0, 0, 255}
			}
		case GamepadInputX:
			if gamepad.X.Up() {
				entity.Color = color.White
			} else if gamepad.X.JustPressed() {
				entity.Color = color.RGBA{0, 255, 0, 255}
			} else if gamepad.X.Down() {
				entity.Color = color.RGBA{255, 0, 0, 255}
			}
		case GamepadInputY:
			if gamepad.Y.Up() {
				entity.Color = color.White
			} else if gamepad.Y.JustPressed() {
				entity.Color = color.RGBA{0, 255, 0, 255}
			} else if gamepad.Y.Down() {
				entity.Color = color.RGBA{255, 0, 0, 255}
			}
		case GamepadInputStart:
			if gamepad.Start.Up() {
				entity.Color = color.White
			} else if gamepad.Start.JustPressed() {
				entity.Color = color.RGBA{0, 255, 0, 255}
			} else if gamepad.Start.Down() {
				entity.Color = color.RGBA{255, 0, 0, 255}
			}
		case GamepadInputBack:
			if gamepad.Back.Up() {
				entity.Color = color.White
			} else if gamepad.Back.JustPressed() {
				entity.Color = color.RGBA{0, 255, 0, 255}
			} else if gamepad.Back.Down() {
				entity.Color = color.RGBA{255, 0, 0, 255}
			}
		case GamepadInputGuide:
			if gamepad.Guide.Up() {
				entity.Color = color.White
			} else if gamepad.Guide.JustPressed() {
				entity.Color = color.RGBA{0, 255, 0, 255}
			} else if gamepad.Guide.Down() {
				entity.Color = color.RGBA{255, 0, 0, 255}
			}
		case GamepadInputLeftBumper:
			if gamepad.LeftBumper.Up() {
				entity.Color = color.White
			} else if gamepad.LeftBumper.JustPressed() {
				entity.Color = color.RGBA{0, 255, 0, 255}
			} else if gamepad.LeftBumper.Down() {
				entity.Color = color.RGBA{255, 0, 0, 255}
			}
		case GamepadInputRightBumper:
			if gamepad.RightBumper.Up() {
				entity.Color = color.White
			} else if gamepad.RightBumper.JustPressed() {
				entity.Color = color.RGBA{0, 255, 0, 255}
			} else if gamepad.RightBumper.Down() {
				entity.Color = color.RGBA{255, 0, 0, 255}
			}
		case GamepadInputLeftThumb:
			if gamepad.LeftThumb.Up() {
				entity.Color = color.White
			} else if gamepad.LeftThumb.JustPressed() {
				entity.Color = color.RGBA{0, 255, 0, 255}
			} else if gamepad.LeftThumb.Down() {
				entity.Color = color.RGBA{255, 0, 0, 255}
			}
		case GamepadInputRightThumb:
			if gamepad.RightThumb.Up() {
				entity.Color = color.White
			} else if gamepad.RightThumb.JustPressed() {
				entity.Color = color.RGBA{0, 255, 0, 255}
			} else if gamepad.RightThumb.Down() {
				entity.Color = color.RGBA{255, 0, 0, 255}
			}
		case GamepadInputDpadUp:
			if gamepad.DpadUp.Up() {
				entity.Color = color.White
			} else if gamepad.DpadUp.JustPressed() {
				entity.Color = color.RGBA{0, 255, 0, 255}
			} else if gamepad.DpadUp.Down() {
				entity.Color = color.RGBA{255, 0, 0, 255}
			}
		case GamepadInputDpadRight:
			if gamepad.DpadRight.Up() {
				entity.Color = color.White
			} else if gamepad.DpadRight.JustPressed() {
				entity.Color = color.RGBA{0, 255, 0, 255}
			} else if gamepad.DpadRight.Down() {
				entity.Color = color.RGBA{255, 0, 0, 255}
			}
		case GamepadInputDpadDown:
			if gamepad.DpadDown.Up() {
				entity.Color = color.White
			} else if gamepad.DpadDown.JustPressed() {
				entity.Color = color.RGBA{0, 255, 0, 255}
			} else if gamepad.DpadDown.Down() {
				entity.Color = color.RGBA{255, 0, 0, 255}
			}
		case GamepadInputDpadLeft:
			if gamepad.DpadLeft.Up() {
				entity.Color = color.White
			} else if gamepad.DpadLeft.JustPressed() {
				entity.Color = color.RGBA{0, 255, 0, 255}
			} else if gamepad.DpadLeft.Down() {
				entity.Color = color.RGBA{255, 0, 0, 255}
			}
		case GamepadInputLeftAxisX:
			entity.Drawable = common.Text{
				Font: font,
				Text: strconv.FormatFloat(float64(gamepad.LeftX.Value()), 'f', 2, 32),
			}
		}
	}
}

func (*DefaultScene) Preload() {
	if err := engo.Files.LoadReaderData("gomonobold.ttf", bytes.NewReader(gomonobold.TTF)); err != nil {
		panic("unable to load gomonobold.ttf! Error was: " + err.Error())
	}

	// Register the gamepad
	err := engo.Input.RegisterGamepad("Player1")
	if err != nil {
		println("Unable to find suitable Gamepad. Error was: ", err.Error())
	}
}

func (d *DefaultScene) Setup(u engo.Updater) {
	w, _ := u.(*ecs.World)

	common.SetBackground(color.White)

	w.AddSystem(&common.RenderSystem{})
	w.AddSystem(&common.FPSSystem{Display: true})
	w.AddSystem(&InputSystem{})

	font = &common.Font{
		URL:  "gomonobold.ttf",
		FG:   color.Black,
		Size: 48,
	}
	if err := font.CreatePreloaded(); err != nil {
		panic("unable to create gomonobold.ttf! Error was: " + err.Error())
	}

	// Buttons
	//  A
	//    Label
	lA := label{BasicEntity: ecs.NewBasic()}
	lA.Drawable = common.Text{
		Font: font,
		Text: "A",
	}
	lA.Position = engo.Point{X: 835, Y: 465}
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&lA.BasicEntity, &lA.RenderComponent, &lA.SpaceComponent)
		}
	}
	//    Panel
	pA := panel{BasicEntity: ecs.NewBasic()}
	pA.Drawable = common.Circle{BorderWidth: 5, BorderColor: color.Black}
	pA.Color = color.White
	pA.Position = engo.Point{X: 800, Y: 365}
	pA.Height = 100
	pA.Width = 100
	pA.Trigger = GamepadInputA
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&pA.BasicEntity, &pA.RenderComponent, &pA.SpaceComponent)
		case *InputSystem:
			sys.Add(&pA.BasicEntity, &pA.RenderComponent, &pA.InputComponent)
		}
	}
	//  B
	//    Label
	lB := label{BasicEntity: ecs.NewBasic()}
	lB.Drawable = common.Text{
		Font: font,
		Text: "B",
	}
	lB.Position = engo.Point{X: 980, Y: 275}
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&lB.BasicEntity, &lB.RenderComponent, &lB.SpaceComponent)
		}
	}
	//    Panel
	pB := panel{BasicEntity: ecs.NewBasic()}
	pB.Drawable = common.Circle{BorderWidth: 5, BorderColor: color.Black}
	pB.Color = color.White
	pB.Position = engo.Point{X: 880, Y: 265}
	pB.Height = 100
	pB.Width = 100
	pB.Trigger = GamepadInputB
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&pB.BasicEntity, &pB.RenderComponent, &pB.SpaceComponent)
		case *InputSystem:
			sys.Add(&pB.BasicEntity, &pB.RenderComponent, &pB.InputComponent)
		}
	}
	//  X
	//    Label
	lX := label{BasicEntity: ecs.NewBasic()}
	lX.Drawable = common.Text{
		Font: font,
		Text: "X",
	}
	lX.Position = engo.Point{X: 675, Y: 275}
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&lX.BasicEntity, &lX.RenderComponent, &lX.SpaceComponent)
		}
	}
	//    Panel
	pX := panel{BasicEntity: ecs.NewBasic()}
	pX.Drawable = common.Circle{BorderWidth: 5, BorderColor: color.Black}
	pX.Color = color.White
	pX.Position = engo.Point{X: 725, Y: 265}
	pX.Height = 100
	pX.Width = 100
	pX.Trigger = GamepadInputX
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&pX.BasicEntity, &pX.RenderComponent, &pX.SpaceComponent)
		case *InputSystem:
			sys.Add(&pX.BasicEntity, &pX.RenderComponent, &pX.InputComponent)
		}
	}
	//  Y
	//    Label
	lY := label{BasicEntity: ecs.NewBasic()}
	lY.Drawable = common.Text{
		Font: font,
		Text: "Y",
	}
	lY.Position = engo.Point{X: 835, Y: 100}
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&lY.BasicEntity, &lY.RenderComponent, &lY.SpaceComponent)
		}
	}
	//    Panel
	pY := panel{BasicEntity: ecs.NewBasic()}
	pY.Drawable = common.Circle{BorderWidth: 5, BorderColor: color.Black}
	pY.Color = color.White
	pY.Position = engo.Point{X: 800, Y: 165}
	pY.Height = 100
	pY.Width = 100
	pY.Trigger = GamepadInputY
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&pY.BasicEntity, &pY.RenderComponent, &pY.SpaceComponent)
		case *InputSystem:
			sys.Add(&pY.BasicEntity, &pY.RenderComponent, &pY.InputComponent)
		}
	}
	//  Back
	//    Label
	lBack := label{BasicEntity: ecs.NewBasic()}
	lBack.Drawable = common.Text{
		Font: font,
		Text: "Back",
	}
	lBack.Position = engo.Point{X: 285, Y: 200}
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&lBack.BasicEntity, &lBack.RenderComponent, &lBack.SpaceComponent)
		}
	}
	//    Panel
	pBack := panel{BasicEntity: ecs.NewBasic()}
	pBack.Drawable = common.Rectangle{BorderWidth: 5, BorderColor: color.Black}
	pBack.Color = color.White
	pBack.Position = engo.Point{X: 300, Y: 275}
	pBack.Height = 75
	pBack.Width = 125
	pBack.Trigger = GamepadInputBack
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&pBack.BasicEntity, &pBack.RenderComponent, &pBack.SpaceComponent)
		case *InputSystem:
			sys.Add(&pBack.BasicEntity, &pBack.RenderComponent, &pBack.InputComponent)
		}
	}
	//  Start
	//    Label
	lStart := label{BasicEntity: ecs.NewBasic()}
	lStart.Drawable = common.Text{
		Font: font,
		Text: "Start",
	}
	lStart.Position = engo.Point{X: 485, Y: 200}
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&lStart.BasicEntity, &lStart.RenderComponent, &lStart.SpaceComponent)
		}
	}
	//    Panel
	pStart := panel{BasicEntity: ecs.NewBasic()}
	pStart.Drawable = common.Rectangle{BorderWidth: 5, BorderColor: color.Black}
	pStart.Color = color.White
	pStart.Position = engo.Point{X: 500, Y: 275}
	pStart.Height = 75
	pStart.Width = 125
	pStart.Trigger = GamepadInputStart
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&pStart.BasicEntity, &pStart.RenderComponent, &pStart.SpaceComponent)
		case *InputSystem:
			sys.Add(&pStart.BasicEntity, &pStart.RenderComponent, &pStart.InputComponent)
		}
	}
	//  Guide
	//    Label
	lGuide := label{BasicEntity: ecs.NewBasic()}
	lGuide.Drawable = common.Text{
		Font: font,
		Text: "Guide",
	}
	lGuide.Position = engo.Point{X: 390, Y: 450}
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&lGuide.BasicEntity, &lGuide.RenderComponent, &lGuide.SpaceComponent)
		}
	}
	//    Panel
	pGuide := panel{BasicEntity: ecs.NewBasic()}
	pGuide.Drawable = common.Circle{BorderWidth: 5, BorderColor: color.Black}
	pGuide.Color = color.White
	pGuide.Position = engo.Point{X: 400, Y: 375}
	pGuide.Height = 75
	pGuide.Width = 125
	pGuide.Trigger = GamepadInputGuide
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&pGuide.BasicEntity, &pGuide.RenderComponent, &pGuide.SpaceComponent)
		case *InputSystem:
			sys.Add(&pGuide.BasicEntity, &pGuide.RenderComponent, &pGuide.InputComponent)
		}
	}
	//  LeftBumper
	//    Label
	lLeftBumper := label{BasicEntity: ecs.NewBasic()}
	lLeftBumper.Drawable = common.Text{
		Font: font,
		Text: "LeftBumper",
	}
	lLeftBumper.Position = engo.Point{X: 15, Y: 50}
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&lLeftBumper.BasicEntity, &lLeftBumper.RenderComponent, &lLeftBumper.SpaceComponent)
		}
	}
	//    Panel
	pLeftBumper := panel{BasicEntity: ecs.NewBasic()}
	pLeftBumper.Drawable = common.Rectangle{BorderWidth: 5, BorderColor: color.Black}
	pLeftBumper.Color = color.White
	pLeftBumper.Position = engo.Point{X: 50, Y: 125}
	pLeftBumper.Height = 50
	pLeftBumper.Width = 200
	pLeftBumper.Trigger = GamepadInputLeftBumper
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&pLeftBumper.BasicEntity, &pLeftBumper.RenderComponent, &pLeftBumper.SpaceComponent)
		case *InputSystem:
			sys.Add(&pLeftBumper.BasicEntity, &pLeftBumper.RenderComponent, &pLeftBumper.InputComponent)
		}
	}
	//  RightBumper
	//    Label
	lRightBumper := label{BasicEntity: ecs.NewBasic()}
	lRightBumper.Drawable = common.Text{
		Font: font,
		Text: "RightBumper",
	}
	lRightBumper.Position = engo.Point{X: 515, Y: 50}
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&lRightBumper.BasicEntity, &lRightBumper.RenderComponent, &lRightBumper.SpaceComponent)
		}
	}
	//    Panel
	pRightBumper := panel{BasicEntity: ecs.NewBasic()}
	pRightBumper.Drawable = common.Rectangle{BorderWidth: 5, BorderColor: color.Black}
	pRightBumper.Color = color.White
	pRightBumper.Position = engo.Point{X: 550, Y: 125}
	pRightBumper.Height = 50
	pRightBumper.Width = 200
	pRightBumper.Trigger = GamepadInputRightBumper
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&pRightBumper.BasicEntity, &pRightBumper.RenderComponent, &pRightBumper.SpaceComponent)
		case *InputSystem:
			sys.Add(&pRightBumper.BasicEntity, &pRightBumper.RenderComponent, &pRightBumper.InputComponent)
		}
	}
	//  LeftThumb
	//    Panel
	pLeftThumb := panel{BasicEntity: ecs.NewBasic()}
	pLeftThumb.Drawable = common.Circle{BorderWidth: 5, BorderColor: color.Black}
	pLeftThumb.Color = color.White
	pLeftThumb.Position = engo.Point{X: 225, Y: 425}
	pLeftThumb.Height = 150
	pLeftThumb.Width = 150
	pLeftThumb.Trigger = GamepadInputLeftThumb
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&pLeftThumb.BasicEntity, &pLeftThumb.RenderComponent, &pLeftThumb.SpaceComponent)
		case *InputSystem:
			sys.Add(&pLeftThumb.BasicEntity, &pLeftThumb.RenderComponent, &pLeftThumb.InputComponent)
		}
	}
	//  RightThumb
	//    Panel
	pRightThumb := panel{BasicEntity: ecs.NewBasic()}
	pRightThumb.Drawable = common.Circle{BorderWidth: 5, BorderColor: color.Black}
	pRightThumb.Color = color.White
	pRightThumb.Position = engo.Point{X: 550, Y: 425}
	pRightThumb.Height = 150
	pRightThumb.Width = 150
	pRightThumb.Trigger = GamepadInputRightThumb
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&pRightThumb.BasicEntity, &pRightThumb.RenderComponent, &pRightThumb.SpaceComponent)
		case *InputSystem:
			sys.Add(&pRightThumb.BasicEntity, &pRightThumb.RenderComponent, &pRightThumb.InputComponent)
		}
	}
	//  DpadUp
	pDpadUp := panel{BasicEntity: ecs.NewBasic()}
	pDpadUp.Drawable = common.Triangle{BorderWidth: 5, BorderColor: color.Black}
	pDpadUp.Color = color.White
	pDpadUp.Position = engo.Point{X: 100, Y: 200}
	pDpadUp.Height = 50
	pDpadUp.Width = 75
	pDpadUp.Trigger = GamepadInputDpadUp
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&pDpadUp.BasicEntity, &pDpadUp.RenderComponent, &pDpadUp.SpaceComponent)
		case *InputSystem:
			sys.Add(&pDpadUp.BasicEntity, &pDpadUp.RenderComponent, &pDpadUp.InputComponent)
		}
	}
	//  DpadRight
	//    Panel
	pDpadRight := panel{BasicEntity: ecs.NewBasic()}
	pDpadRight.Drawable = common.Triangle{BorderWidth: 5, BorderColor: color.Black}
	pDpadRight.Color = color.White
	pDpadRight.Position = engo.Point{X: 250, Y: 275}
	pDpadRight.Height = 50
	pDpadRight.Width = 75
	pDpadRight.Rotation = 90
	pDpadRight.Trigger = GamepadInputDpadRight
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&pDpadRight.BasicEntity, &pDpadRight.RenderComponent, &pDpadRight.SpaceComponent)
		case *InputSystem:
			sys.Add(&pDpadRight.BasicEntity, &pDpadRight.RenderComponent, &pDpadRight.InputComponent)
		}
	}
	//  DpadDown
	//    Panel
	pDpadDown := panel{BasicEntity: ecs.NewBasic()}
	pDpadDown.Drawable = common.Triangle{BorderWidth: 5, BorderColor: color.Black}
	pDpadDown.Color = color.White
	pDpadDown.Position = engo.Point{X: 170, Y: 425}
	pDpadDown.Height = 50
	pDpadDown.Width = 75
	pDpadDown.Rotation = 180
	pDpadDown.Trigger = GamepadInputDpadDown
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&pDpadDown.BasicEntity, &pDpadDown.RenderComponent, &pDpadDown.SpaceComponent)
		case *InputSystem:
			sys.Add(&pDpadDown.BasicEntity, &pDpadDown.RenderComponent, &pDpadDown.InputComponent)
		}
	}
	//  DpadLeft
	//    Panel
	pDpadLeft := panel{BasicEntity: ecs.NewBasic()}
	pDpadLeft.Drawable = common.Triangle{BorderWidth: 5, BorderColor: color.Black}
	pDpadLeft.Color = color.White
	pDpadLeft.Position = engo.Point{X: 25, Y: 350}
	pDpadLeft.Height = 50
	pDpadLeft.Width = 75
	pDpadLeft.Rotation = 270
	pDpadLeft.Trigger = GamepadInputDpadLeft
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&pDpadLeft.BasicEntity, &pDpadLeft.RenderComponent, &pDpadLeft.SpaceComponent)
		case *InputSystem:
			sys.Add(&pDpadLeft.BasicEntity, &pDpadLeft.RenderComponent, &pDpadLeft.InputComponent)
		}
	}
	// Axes
	//  LeftX
	//    Label
	lLeftX := label{BasicEntity: ecs.NewBasic()}
	lLeftX.Drawable = common.Text{
		Font: font,
		Text: "X:",
	}
	lLeftX.Position = engo.Point{X: 25, Y: 525}
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&lLeftX.BasicEntity, &lLeftX.RenderComponent, &lLeftX.SpaceComponent)
		}
	}
	//    AxisNumber
	nLeftX := panel{BasicEntity: ecs.NewBasic()}
	nLeftX.Drawable = common.Text{
		Font: font,
		Text: "0.00",
	}
	nLeftX.Position = engo.Point{X: 45, Y: 525}
	nLeftX.Trigger = GamepadInputLeftAxisX
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&nLeftX.BasicEntity, &nLeftX.RenderComponent, &nLeftX.SpaceComponent)
		case *InputSystem:
			sys.Add(&nLeftX.BasicEntity, &nLeftX.RenderComponent, &nLeftX.InputComponent)
		}
	}
	//  LeftY
	//    Label
	lLeftY := label{BasicEntity: ecs.NewBasic()}
	lLeftY.Drawable = common.Text{
		Font: font,
		Text: "Y:",
	}
	lLeftY.Position = engo.Point{X: 15, Y: 575}
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&lLeftY.BasicEntity, &lLeftY.RenderComponent, &lLeftY.SpaceComponent)
		}
	}
	//    AxisNumber
	//  RightX
	//    Label
	lRightX := label{BasicEntity: ecs.NewBasic()}
	lRightX.Drawable = common.Text{
		Font: font,
		Text: "X:",
	}
	lRightX.Position = engo.Point{X: 715, Y: 525}
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&lRightX.BasicEntity, &lRightX.RenderComponent, &lRightX.SpaceComponent)
		}
	}
	//    AxisNumber
	//  RightY
	//    Label
	lRightY := label{BasicEntity: ecs.NewBasic()}
	lRightY.Drawable = common.Text{
		Font: font,
		Text: "Y:",
	}
	lRightY.Position = engo.Point{X: 715, Y: 575}
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&lRightY.BasicEntity, &lRightY.RenderComponent, &lRightY.SpaceComponent)
		}
	}
	//    AxisNumber
	//  LeftTrigger
	//  RightTrigger
}

func (*DefaultScene) Type() string { return "GameWorld" }

func main() {
	opts := engo.RunOptions{
		Title:    "Gamepad Demo",
		Width:    1024,
		Height:   640,
		FPSLimit: 200,
	}
	engo.Run(opts, &DefaultScene{})
}
