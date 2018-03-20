package engo

import "testing"

type btnState struct {
	down     bool
	justUp   bool
	justDown bool
}

// Button configuretion used when testing.
var btnSimpleCfg = [6]Button{
	Button{Triggers: []Key{KeyA, KeyC}, Name: "Button 1"},
	Button{Triggers: []Key{KeyB, KeyD}, Name: "Button 2"},
	Button{Triggers: []Key{KeyF2, KeyF5}, Name: "Button 3"},
	Button{Triggers: []Key{KeyF4, KeyF6}, Name: "Button 4"},
	Button{Triggers: []Key{KeyOne, KeyFour}, Name: "Button 5"},
	Button{Triggers: []Key{KeyTwo, KeyFive}, Name: "Button 6"},
}

// Expected button state @ pass 0
var btnPass0 = [6]btnState{
	btnState{down: false, justUp: false, justDown: false},
	btnState{down: false, justUp: false, justDown: false},
	btnState{down: false, justUp: false, justDown: false},
	btnState{down: false, justUp: false, justDown: false},
	btnState{down: false, justUp: false, justDown: false},
	btnState{down: false, justUp: false, justDown: false},
}

// Expected button state @ pass 1
var btnPass1 = [6]btnState{
	btnState{down: false, justUp: false, justDown: false},
	btnState{down: false, justUp: false, justDown: true},
	btnState{down: false, justUp: false, justDown: false},
	btnState{down: false, justUp: false, justDown: true},
	btnState{down: false, justUp: false, justDown: false},
	btnState{down: false, justUp: false, justDown: true},
}

// Expected button state @ pass 2
var btnPass2 = [6]btnState{
	btnState{down: false, justUp: false, justDown: false},
	btnState{down: true, justUp: false, justDown: false},
	btnState{down: false, justUp: false, justDown: false},
	btnState{down: true, justUp: false, justDown: false},
	btnState{down: false, justUp: false, justDown: false},
	btnState{down: true, justUp: false, justDown: false},
}

// Expected button state @ pass 3
var btnPass3 = [6]btnState{
	btnState{down: false, justUp: false, justDown: false},
	btnState{down: false, justUp: true, justDown: false},
	btnState{down: false, justUp: false, justDown: false},
	btnState{down: false, justUp: true, justDown: false},
	btnState{down: false, justUp: false, justDown: false},
	btnState{down: false, justUp: true, justDown: false},
}

// Expected button state @ pass 4
var btnPass4 = [6]btnState{
	btnState{down: false, justUp: false, justDown: false},
	btnState{down: false, justUp: false, justDown: false},
	btnState{down: false, justUp: false, justDown: false},
	btnState{down: false, justUp: false, justDown: false},
	btnState{down: false, justUp: false, justDown: false},
	btnState{down: false, justUp: false, justDown: false},
}

// Checks the state of all configured buttons against the expected state.
func runBtnChecks(msg string, t *testing.T, expect [6]btnState) {
	for i, cfg := range btnSimpleCfg {
		exp := expect[i]
		btn := Input.Button(cfg.Name)
		if exp.down != btn.Down() {
			t.Error(msg, " Invalid on: ", cfg.Name, " - Down")
		}
		if exp.justUp != btn.JustReleased() {
			t.Error(msg, " Invalid on: ", cfg.Name, " - Just Up")
		}
		if exp.justDown != btn.JustPressed() {
			t.Error(msg, " Invalid on: ", cfg.Name, " - Just Down")
		}
	}
}

// Test configured axes using a single key on one button.
func TestButtonSimple(t *testing.T) {
	Input = NewInputManager()

	for _, cfg := range btnSimpleCfg {
		Input.RegisterButton(cfg.Name, cfg.Triggers[0])
	}

	runBtnChecks("Init (0.0)", t, btnPass0)

	// Empty update pass0
	Input.update()
	runBtnChecks("Pass (0.1)", t, btnPass0)
	Input.update()
	runBtnChecks("Pass (0.2)", t, btnPass0)
	Input.update()
	runBtnChecks("Pass (0.3)", t, btnPass0)

	// Set even true pass1
	Input.update()
	Input.keys.Set(btnSimpleCfg[1].Triggers[0], true)
	Input.keys.Set(btnSimpleCfg[3].Triggers[0], true)
	Input.keys.Set(btnSimpleCfg[5].Triggers[0], true)

	// FixMe: this causes an error ? Because the the static
	// arrays get filled with values before engo rewrites them!
	//Input.keys.Set(F10, true)

	runBtnChecks("Pass (1.0)", t, btnPass1)

	// Keeps state on pass2
	Input.update()
	runBtnChecks("Pass (2.0)", t, btnPass2)
	Input.update()
	runBtnChecks("Pass (2.1)", t, btnPass2)
	Input.update()
	runBtnChecks("Pass (2.2)", t, btnPass2)
	Input.update()
	runBtnChecks("Pass (2.3)", t, btnPass2)

	// Set even true pass3
	Input.update()
	Input.keys.Set(btnSimpleCfg[1].Triggers[0], false)
	Input.keys.Set(btnSimpleCfg[3].Triggers[0], false)
	Input.keys.Set(btnSimpleCfg[5].Triggers[0], false)

	runBtnChecks("Pass (3.0)", t, btnPass3)

	// Keeps state on pass4
	Input.update()
	runBtnChecks("Pass (4.0)", t, btnPass4)
	Input.update()
	runBtnChecks("Pass (4.1)", t, btnPass4)
	Input.update()
	runBtnChecks("Pass (4.2)", t, btnPass4)
	Input.update()
	runBtnChecks("Pass (4.3)", t, btnPass4)
}

// Test configured axes using a multiple keys on one button.
func TestButtonComplex(t *testing.T) {
	Input = NewInputManager()

	for _, cfg := range btnSimpleCfg {
		Input.RegisterButton(cfg.Name, cfg.Triggers[0], cfg.Triggers[1])
	}

	runBtnChecks("Init (0.0)", t, btnPass0)

	// Empty update pass0
	Input.update()
	runBtnChecks("Pass (0.1)", t, btnPass0)
	Input.update()
	runBtnChecks("Pass (0.2)", t, btnPass0)
	Input.update()
	runBtnChecks("Pass (0.3)", t, btnPass0)

	// Set even true pass1
	Input.update()
	Input.keys.Set(btnSimpleCfg[1].Triggers[0], true)
	Input.keys.Set(btnSimpleCfg[3].Triggers[0], true)
	Input.keys.Set(btnSimpleCfg[5].Triggers[0], true)

	runBtnChecks("Pass (1.0)", t, btnPass1)

	// Keeps state on pass2
	Input.update()
	runBtnChecks("Pass (2.0)", t, btnPass2)
	Input.update()
	runBtnChecks("Pass (2.1)", t, btnPass2)
	Input.update()
	runBtnChecks("Pass (2.2)", t, btnPass2)
	Input.update()
	runBtnChecks("Pass (2.3)", t, btnPass2)

	// Set even true pass3
	Input.update()
	Input.keys.Set(btnSimpleCfg[1].Triggers[0], false)
	Input.keys.Set(btnSimpleCfg[3].Triggers[0], false)
	Input.keys.Set(btnSimpleCfg[5].Triggers[0], false)

	runBtnChecks("Pass (3.0)", t, btnPass3)

	// Keeps state on pass4
	Input.update()
	runBtnChecks("Pass (4.0)", t, btnPass4)
	Input.update()
	runBtnChecks("Pass (4.1)", t, btnPass4)
	Input.update()
	runBtnChecks("Pass (4.2)", t, btnPass4)
	Input.update()
	runBtnChecks("Pass (4.3)", t, btnPass4)

	// Set even true pass1 alt
	Input.update()
	Input.keys.Set(btnSimpleCfg[1].Triggers[1], true)
	Input.keys.Set(btnSimpleCfg[3].Triggers[1], true)
	Input.keys.Set(btnSimpleCfg[5].Triggers[1], true)

	runBtnChecks("Pass alt (1.0)", t, btnPass1)

	// Keeps state on pass2 alt
	Input.update()
	runBtnChecks("Pass alt (2.0)", t, btnPass2)
	Input.update()
	runBtnChecks("Pass alt (2.1)", t, btnPass2)
	Input.update()
	runBtnChecks("Pass alt (2.2)", t, btnPass2)
	Input.update()
	runBtnChecks("Pass alt (2.3)", t, btnPass2)

	// Set even true pass3 alt
	Input.update()
	Input.keys.Set(btnSimpleCfg[1].Triggers[1], false)
	Input.keys.Set(btnSimpleCfg[3].Triggers[1], false)
	Input.keys.Set(btnSimpleCfg[5].Triggers[1], false)

	runBtnChecks("Pass alt (3.0)", t, btnPass3)

	// Keeps state on pass4 alt
	Input.update()
	runBtnChecks("Pass alt (4.0)", t, btnPass4)
	Input.update()
	runBtnChecks("Pass alt (4.1)", t, btnPass4)
	Input.update()
	runBtnChecks("Pass alt (4.2)", t, btnPass4)
	Input.update()
	runBtnChecks("Pass alt (4.3)", t, btnPass4)
}

// Used to store results when benchmarking.
var btnResult [6]btnState

func checkBtnConfigOptimal(b *testing.B) {
	for i, cfg := range btnSimpleCfg {
		btn := Input.Button(cfg.Name)
		btnResult[i].down = btn.Down()
		btnResult[i].justUp = btn.JustReleased()
		btnResult[i].justDown = btn.JustPressed()
	}
}

// Benchmark sub-optimal state checks
func BenchmarkInputMgr_ButtonCleanState(b *testing.B) {
	Input = NewInputManager()

	for _, cfg := range btnSimpleCfg {
		Input.RegisterButton(cfg.Name, cfg.Triggers[0], cfg.Triggers[1])
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		checkBtnConfigOptimal(b)
	}
}

// Benchmark sub-optimal state checks
func BenchmarkInputMgr_ButtonFilledState(b *testing.B) {
	Input = NewInputManager()
	keyFillManager(Input.keys)

	for _, cfg := range btnSimpleCfg {
		Input.RegisterButton(cfg.Name, cfg.Triggers[0], cfg.Triggers[1])
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		checkBtnConfigOptimal(b)
	}
}

// Disabled but around when needed

//func checkBtnConfigSubOptimal(b *testing.B) {
//	for i, cfg := range btnSimpleCfg {
//		btnResult[i].down = Input.Button(cfg.Name).Down()
//		btnResult[i].justUp = Input.Button(cfg.Name).JustReleased()
//		btnResult[i].justDown = Input.Button(cfg.Name).JustPressed()
//	}
//}
// Benchmark sub-optimal state checks
//func BenchmarkInputMgr_ButtonSubOptimal(b *testing.B) {
//	Input = NewInputManager()
//
//	for _, cfg := range btnSimpleCfg {
//		Input.RegisterButton(cfg.Name, cfg.Triggers[0], cfg.Triggers[1])
//	}
//
//	b.ResetTimer()
//	for n := 0; n < b.N; n++ {
//		checkBtnConfigSubOptimal(b)
//	}
//}
