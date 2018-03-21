package engo

import "testing"

type axState struct {
	value float32
}

type axKeyCfg struct {
	Name  string
	Pairs []AxisKeyPair
}

// Axis configuretion used when testing.
var axSimpleCfg = [6]axKeyCfg{
	axKeyCfg{
		Name: "Axis 1",
		Pairs: []AxisKeyPair{
			AxisKeyPair{Min: KeyA, Max: KeyB},
			AxisKeyPair{Min: KeyC, Max: KeyD},
		},
	},
	axKeyCfg{
		Name: "Axis 2",
		Pairs: []AxisKeyPair{
			AxisKeyPair{Min: KeyE, Max: KeyF},
			AxisKeyPair{Min: KeyG, Max: KeyH},
		},
	},
	axKeyCfg{
		Name: "Axis 3",
		Pairs: []AxisKeyPair{
			AxisKeyPair{Min: KeyF1, Max: KeyF2},
			AxisKeyPair{Min: KeyF3, Max: KeyF4},
		},
	},
	axKeyCfg{
		Name: "Axis 4",
		Pairs: []AxisKeyPair{
			AxisKeyPair{Min: KeyF5, Max: KeyF6},
			AxisKeyPair{Min: KeyF7, Max: KeyF8},
		},
	},
	axKeyCfg{
		Name: "Axis 5",
		Pairs: []AxisKeyPair{
			AxisKeyPair{Min: KeyOne, Max: KeyTwo},
			AxisKeyPair{Min: KeyThree, Max: KeyFour},
		},
	},
	axKeyCfg{
		Name: "Axis 6",
		Pairs: []AxisKeyPair{
			AxisKeyPair{Min: KeyArrowUp, Max: KeyArrowDown},
			AxisKeyPair{Min: KeyArrowLeft, Max: KeyArrowRight},
		},
	},
}

// Expected axes values @ pass 0
var axPass0 = [6]axState{
	axState{value: 0.0},
	axState{value: 0.0},
	axState{value: 0.0},
	axState{value: 0.0},
	axState{value: 0.0},
	axState{value: 0.0},
}

// Expected axes values @ pass 1
var axPass1 = [6]axState{
	axState{value: 0.0},
	axState{value: 0.0},
	axState{value: 0.0},
	axState{value: 0.0},
	axState{value: 0.0},
	axState{value: 0.0},
}

// Expected axes values @ pass 2
var axPass2 = [6]axState{
	axState{value: 0.0},
	axState{value: 1.0},
	axState{value: 0.0},
	axState{value: 1.0},
	axState{value: 0.0},
	axState{value: 1.0},
}

// Expected axes values @ pass 3
var axPass3 = [6]axState{
	axState{value: 0.0},
	axState{value: 0.0},
	axState{value: 0.0},
	axState{value: 0.0},
	axState{value: 0.0},
	axState{value: 0.0},
}

// Expected axes values @ pass 4
var axPass4 = [6]axState{
	axState{value: 0.0},
	axState{value: 0.0},
	axState{value: 0.0},
	axState{value: 0.0},
	axState{value: 0.0},
	axState{value: 0.0},
}

// Expected axes values @ pass 5
var axPass5 = [6]axState{
	axState{value: 0.0},
	axState{value: 0.0},
	axState{value: 0.0},
	axState{value: 0.0},
	axState{value: 0.0},
	axState{value: 0.0},
}

// Expected axes values @ pass 6
var axPass6 = [6]axState{
	axState{value: 0.0},
	axState{value: -1.0},
	axState{value: 0.0},
	axState{value: -1.0},
	axState{value: 0.0},
	axState{value: -1.0},
}

// Expected axes values @ pass 7
var axPass7 = [6]axState{
	axState{value: 0.0},
	axState{value: 0.0},
	axState{value: 0.0},
	axState{value: 0.0},
	axState{value: 0.0},
	axState{value: 0.0},
}

// Expected axes values @ pass 8
var axPass8 = [6]axState{
	axState{value: 0.0},
	axState{value: 0.0},
	axState{value: 0.0},
	axState{value: 0.0},
	axState{value: 0.0},
	axState{value: 0.0},
}

// Checks the value of all configured axis against the expected values.
func runAxisChecks(msg string, t *testing.T, expect [6]axState) {
	for i, cfg := range axSimpleCfg {
		exp := expect[i]
		axi := Input.Axis(cfg.Name)
		if exp.value != axi.Value() {
			t.Error(msg, " Invalid on: ", cfg.Name, " - Value")
		}
	}
}

// Test configured axes using a single pair on one axis.
func TestAxisSimple(t *testing.T) {
	Input = NewInputManager()

	for _, cfg := range axSimpleCfg {
		Input.RegisterAxis(cfg.Name, cfg.Pairs[0])
	}

	runAxisChecks("Init (0.0)", t, axPass0)

	// Empty update pass0
	Input.update()
	runAxisChecks("Pass (0.1)", t, axPass0)
	Input.update()
	runAxisChecks("Pass (0.2)", t, axPass0)
	Input.update()
	runAxisChecks("Pass (0.3)", t, axPass0)

	// Set even true pass1
	Input.update()
	Input.keys.Set(axSimpleCfg[1].Pairs[0].Max, true)
	Input.keys.Set(axSimpleCfg[3].Pairs[0].Max, true)
	Input.keys.Set(axSimpleCfg[5].Pairs[0].Max, true)

	runAxisChecks("Pass (1.0)", t, axPass1)

	// Keeps state on pass2
	Input.update()
	runAxisChecks("Pass (2.0)", t, axPass2)
	Input.update()
	runAxisChecks("Pass (2.1)", t, axPass2)
	Input.update()
	runAxisChecks("Pass (2.2)", t, axPass2)
	Input.update()
	runAxisChecks("Pass (2.3)", t, axPass2)

	// Set even true pass3
	Input.update()
	Input.keys.Set(axSimpleCfg[1].Pairs[0].Max, false)
	Input.keys.Set(axSimpleCfg[3].Pairs[0].Max, false)
	Input.keys.Set(axSimpleCfg[5].Pairs[0].Max, false)

	runAxisChecks("Pass (3.0)", t, axPass3)

	// Keeps state on pass4
	Input.update()
	runAxisChecks("Pass (4.0)", t, axPass4)
	Input.update()
	runAxisChecks("Pass (4.1)", t, axPass4)
	Input.update()
	runAxisChecks("Pass (4.2)", t, axPass4)
	Input.update()
	runAxisChecks("Pass (4.3)", t, axPass4)

	// Set even true pass5
	Input.update()
	Input.keys.Set(axSimpleCfg[1].Pairs[0].Min, true)
	Input.keys.Set(axSimpleCfg[3].Pairs[0].Min, true)
	Input.keys.Set(axSimpleCfg[5].Pairs[0].Min, true)

	runAxisChecks("Pass (5.0)", t, axPass5)

	// Keeps state on pass6
	Input.update()
	runAxisChecks("Pass (6.0)", t, axPass6)
	Input.update()
	runAxisChecks("Pass (6.1)", t, axPass6)
	Input.update()
	runAxisChecks("Pass (6.2)", t, axPass6)
	Input.update()
	runAxisChecks("Pass (6.3)", t, axPass6)

	// Set even true pass7
	Input.update()
	Input.keys.Set(axSimpleCfg[1].Pairs[0].Min, false)
	Input.keys.Set(axSimpleCfg[3].Pairs[0].Min, false)
	Input.keys.Set(axSimpleCfg[5].Pairs[0].Min, false)

	runAxisChecks("Pass (7.0)", t, axPass7)

	// Keeps state on pass8
	Input.update()
	runAxisChecks("Pass (8.0)", t, axPass8)
	Input.update()
	runAxisChecks("Pass (8.1)", t, axPass8)
	Input.update()
	runAxisChecks("Pass (8.2)", t, axPass8)
	Input.update()
	runAxisChecks("Pass (8.3)", t, axPass8)
}

// Test configured axes using multiple pairs on one axis.
func TestAxisComplex(t *testing.T) {
	Input = NewInputManager()

	for _, cfg := range axSimpleCfg {
		Input.RegisterAxis(cfg.Name, cfg.Pairs[0], cfg.Pairs[1])
	}

	runAxisChecks("Init (0.0)", t, axPass0)

	// Empty update pass0
	Input.update()
	runAxisChecks("Pass (0.1)", t, axPass0)
	Input.update()
	runAxisChecks("Pass (0.2)", t, axPass0)
	Input.update()
	runAxisChecks("Pass (0.3)", t, axPass0)

	// Set even true pass1
	Input.update()
	Input.keys.Set(axSimpleCfg[1].Pairs[0].Max, true)
	Input.keys.Set(axSimpleCfg[3].Pairs[0].Max, true)
	Input.keys.Set(axSimpleCfg[5].Pairs[0].Max, true)

	runAxisChecks("Pass (1.0)", t, axPass1)

	// Keeps state on pass2
	Input.update()
	runAxisChecks("Pass (2.0)", t, axPass2)
	Input.update()
	runAxisChecks("Pass (2.1)", t, axPass2)
	Input.update()
	runAxisChecks("Pass (2.2)", t, axPass2)
	Input.update()
	runAxisChecks("Pass (2.3)", t, axPass2)

	// Set even true pass3
	Input.update()
	Input.keys.Set(axSimpleCfg[1].Pairs[0].Max, false)
	Input.keys.Set(axSimpleCfg[3].Pairs[0].Max, false)
	Input.keys.Set(axSimpleCfg[5].Pairs[0].Max, false)

	runAxisChecks("Pass (3.0)", t, axPass3)

	// Keeps state on pass4
	Input.update()
	runAxisChecks("Pass (4.0)", t, axPass4)
	Input.update()
	runAxisChecks("Pass (4.1)", t, axPass4)
	Input.update()
	runAxisChecks("Pass (4.2)", t, axPass4)
	Input.update()
	runAxisChecks("Pass (4.3)", t, axPass4)

	// Set even true pass5
	Input.update()
	Input.keys.Set(axSimpleCfg[1].Pairs[0].Min, true)
	Input.keys.Set(axSimpleCfg[3].Pairs[0].Min, true)
	Input.keys.Set(axSimpleCfg[5].Pairs[0].Min, true)

	runAxisChecks("Pass (5.0)", t, axPass5)

	// Keeps state on pass6
	Input.update()
	runAxisChecks("Pass (6.0)", t, axPass6)
	Input.update()
	runAxisChecks("Pass (6.1)", t, axPass6)
	Input.update()
	runAxisChecks("Pass (6.2)", t, axPass6)
	Input.update()
	runAxisChecks("Pass (6.3)", t, axPass6)

	// Set even true pass7
	Input.update()
	Input.keys.Set(axSimpleCfg[1].Pairs[0].Min, false)
	Input.keys.Set(axSimpleCfg[3].Pairs[0].Min, false)
	Input.keys.Set(axSimpleCfg[5].Pairs[0].Min, false)

	runAxisChecks("Pass (7.0)", t, axPass7)

	// Keeps state on pass8
	Input.update()
	runAxisChecks("Pass (8.0)", t, axPass8)
	Input.update()
	runAxisChecks("Pass (8.1)", t, axPass8)
	Input.update()
	runAxisChecks("Pass (8.2)", t, axPass8)
	Input.update()
	runAxisChecks("Pass (8.3)", t, axPass8)

	// Set even true pass1 alt
	Input.update()
	Input.keys.Set(axSimpleCfg[1].Pairs[1].Max, true)
	Input.keys.Set(axSimpleCfg[3].Pairs[1].Max, true)
	Input.keys.Set(axSimpleCfg[5].Pairs[1].Max, true)

	runAxisChecks("Pass alt (1.0)", t, axPass1)

	// Keeps state on pass2 alt
	Input.update()
	runAxisChecks("Pass alt (2.0)", t, axPass2)
	Input.update()
	runAxisChecks("Pass alt (2.1)", t, axPass2)
	Input.update()
	runAxisChecks("Pass alt (2.2)", t, axPass2)
	Input.update()
	runAxisChecks("Pass alt (2.3)", t, axPass2)

	// Set even true pass3 alt
	Input.update()
	Input.keys.Set(axSimpleCfg[1].Pairs[1].Max, false)
	Input.keys.Set(axSimpleCfg[3].Pairs[1].Max, false)
	Input.keys.Set(axSimpleCfg[5].Pairs[1].Max, false)

	runAxisChecks("Pass alt (3.0)", t, axPass3)

	// Keeps state on pass4 alt
	Input.update()
	runAxisChecks("Pass alt (4.0)", t, axPass4)
	Input.update()
	runAxisChecks("Pass alt (4.1)", t, axPass4)
	Input.update()
	runAxisChecks("Pass alt (4.2)", t, axPass4)
	Input.update()
	runAxisChecks("Pass alt (4.3)", t, axPass4)

	// Set even true pass5 alt
	Input.update()
	Input.keys.Set(axSimpleCfg[1].Pairs[1].Min, true)
	Input.keys.Set(axSimpleCfg[3].Pairs[1].Min, true)
	Input.keys.Set(axSimpleCfg[5].Pairs[1].Min, true)

	runAxisChecks("Pass alt (5.0)", t, axPass5)

	// Keeps state on pass6 alt
	Input.update()
	runAxisChecks("Pass alt (6.0)", t, axPass6)
	Input.update()
	runAxisChecks("Pass alt (6.1)", t, axPass6)
	Input.update()
	runAxisChecks("Pass alt (6.2)", t, axPass6)
	Input.update()
	runAxisChecks("Pass alt (6.3)", t, axPass6)

	// Set even true pass7
	Input.update()
	Input.keys.Set(axSimpleCfg[1].Pairs[1].Min, false)
	Input.keys.Set(axSimpleCfg[3].Pairs[1].Min, false)
	Input.keys.Set(axSimpleCfg[5].Pairs[1].Min, false)

	runAxisChecks("Pass alt (7.0)", t, axPass7)

	// Keeps state on pass8 alt
	Input.update()
	runAxisChecks("Pass alt (8.0)", t, axPass8)
	Input.update()
	runAxisChecks("Pass alt (8.1)", t, axPass8)
	Input.update()
	runAxisChecks("Pass alt (8.2)", t, axPass8)
	Input.update()
	runAxisChecks("Pass alt (8.3)", t, axPass8)
}

// Checks the state of the two mouse axes against provided values.
func runAxisMouse(msg string, t *testing.T, x float32, y float32) {

	if x != Input.Axis("mouse x").Value() {
		t.Error(msg, "Invalid x value: ", x, "!=", Input.Axis("mouse x").Value())
	}
	if y != Input.Axis("mouse y").Value() {
		t.Error(msg, "Invalid y value: ", y, "!=", Input.Axis("mouse y").Value())
	}
}

// Test some state changes on the two mouse axes.
func TestAxisMouse(t *testing.T) {
	Input = NewInputManager()

	SetGlobalScale(Point{1, 1})

	Input.RegisterAxis("mouse x", NewAxisMouse(AxisMouseHori))
	Input.RegisterAxis("mouse y", NewAxisMouse(AxisMouseVert))

	runAxisMouse("Pass 0", t, 0.0, 0.0)

	Input.Mouse.X = 1.0
	Input.Mouse.Y = 1.0
	runAxisMouse("Pass 1", t, 1.0, 1.0)

	// Resets state to 0.0
	runAxisMouse("Pass 2", t, 0.0, 0.0)
	runAxisMouse("Pass 3", t, 0.0, 0.0)

	Input.Mouse.Y = -1.0
	Input.Mouse.X = -1.0
	runAxisMouse("Pass 4", t, -2.0, -2.0)

	Input.Mouse.X = 0.0
	Input.Mouse.Y = 0.0
	runAxisMouse("Pass 4", t, 1.0, 1.0)

	runAxisMouse("Pass 6", t, 0.0, 0.0)
	runAxisMouse("Pass 7", t, 0.0, 0.0)
}

// Used to store results when benchmarking.
var axResult [6]axState

// Fast check for all configured axes, stores the value externaly.
func checkAxisConfigValue(b *testing.B) {
	for i, cfg := range axSimpleCfg {
		axResult[i].value = Input.Axis(cfg.Name).Value()
	}
}

// Benchmark values checks with a clean key manager.
func BenchmarkInputMgr_AxisCleanState(b *testing.B) {
	Input = NewInputManager()

	for _, cfg := range axSimpleCfg {
		Input.RegisterAxis(cfg.Name, cfg.Pairs[0], cfg.Pairs[1])
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		checkAxisConfigValue(b)
	}
}

// Benchmark values checks with a full key manager.
func BenchmarkInputMgr_AxisFilledState(b *testing.B) {
	Input = NewInputManager()
	keyFillManager(Input.keys)

	for _, cfg := range axSimpleCfg {
		Input.RegisterAxis(cfg.Name, cfg.Pairs[0], cfg.Pairs[1])
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		checkAxisConfigValue(b)
	}
}

// Disabled but around when needed

//func checkAxisMouseValue(b *testing.B) {
//	axResult[0].value = Input.Axis("mouse x").Value()
//	axResult[1].value = Input.Axis("mouse y").Value()
//}
// Benchmark sub-optimal state checks
//func BenchmarkInputMgr_CleanMouseAxisValues(b *testing.B) {
//	Input = NewInputManager()
//
//	Input.RegisterAxis("mouse x", NewAxisMouse(AxisMouseHori))
//	Input.RegisterAxis("mouse y", NewAxisMouse(AxisMouseVert))
//
//	b.ResetTimer()
//	for n := 0; n < b.N; n++ {
//		checkAxisMouseValue(b)
//	}
//}
