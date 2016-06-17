package act

import "testing"

func TestActMgr(t *testing.T) {
	mgr := NewActManager()

	mgr.Clear()
	mgr.Update()

	runActCheck("Init (0.A)", t, mgr, KeyA, StateIdle)
	runActCheck("Init (0.F12)", t, mgr, KeyF12, StateIdle)
	runActCheck("Init (0.Left)", t, mgr, MouseLeft, StateIdle)

	mgr.Clear()
	mgr.SetState(KeyA, true)
	mgr.SetState(KeyF12, false)
	mgr.Update()
	mgr.SetState(MouseLeft, true)

	runActCheck("Pass (1.A)", t, mgr, KeyA, StateJustActive)
	runActCheck("Pass (1.F12)", t, mgr, KeyF12, StateIdle)
	runActCheck("Pass (1.Left)", t, mgr, MouseLeft, StateIdle)

	mgr.Clear()
	mgr.Update()

	runActCheck("Pass (2.A)", t, mgr, KeyA, StateActive)
	runActCheck("Pass (2.F12)", t, mgr, KeyF12, StateIdle)
	runActCheck("Pass (2.Left)", t, mgr, MouseLeft, StateJustActive)

	mgr.Clear()
	mgr.SetState(KeyA, false)
	mgr.SetState(KeyF12, true)
	mgr.Update()
	mgr.SetState(MouseLeft, false)

	runActCheck("Pass (3.A)", t, mgr, KeyA, StateJustIdle)
	runActCheck("Pass (3.F12)", t, mgr, KeyF12, StateJustActive)
	runActCheck("Pass (3.Left)", t, mgr, MouseLeft, StateActive)

	mgr.Clear()
	mgr.Update()

	runActCheck("Pass (4.A)", t, mgr, KeyA, StateIdle)
	runActCheck("Pass (4.F12)", t, mgr, KeyF12, StateActive)
	runActCheck("Pass (4.Left)", t, mgr, MouseLeft, StateJustIdle)

	mgr.Clear()
	mgr.Update()

	runActCheck("Pass (5.A)", t, mgr, KeyA, StateIdle)
	runActCheck("Pass (5.F12)", t, mgr, KeyF12, StateActive)
	runActCheck("Pass (5.Left)", t, mgr, MouseLeft, StateIdle)

	mgr.Clear()
	mgr.SetState(KeyF12, false)
	mgr.Update()

	runActCheck("Pass (6.A)", t, mgr, KeyA, StateIdle)
	runActCheck("Pass (6.F12)", t, mgr, KeyF12, StateJustIdle)
	runActCheck("Pass (6.Left)", t, mgr, MouseLeft, StateIdle)

	mgr.Clear()
	mgr.Update()

	runActCheck("Pass (7.A)", t, mgr, KeyA, StateIdle)
	runActCheck("Pass (7.F12)", t, mgr, KeyF12, StateIdle)
	runActCheck("Pass (7.Left)", t, mgr, MouseLeft, StateIdle)
}

func runActCheck(msg string, t *testing.T, mgr *ActManager, act Code, exp State) {
	if exp != mgr.State(act) {
		t.Error(msg, " Invalid on: State")
	}
	if (StateIdle == exp) != mgr.Idle(act) {
		t.Error(msg, " - Invalid on: Idle")
	}
	if (StateActive == exp) != mgr.Active(act) {
		t.Error(msg, " - Invalid on: Active")
	}
	if (StateJustIdle == exp) != mgr.JustIdle(act) {
		t.Error(msg, " - Invalid on: Just Idle")
	}
	if (StateJustActive == exp) != mgr.JustActive(act) {
		t.Error(msg, " - Invalid on: Just Active")
	}
}

var bState [4]State
var bResult [4]bool

func BenchmarkActMgr_CleanSimulate(b *testing.B) {
	mgr := NewActManager()

	mgr.Clear()
	mgr.Update()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		mgr.Clear()
		mgr.Update()

		bState[0] = mgr.State(KeyA)

		mgr.Clear()
		mgr.SetState(KeyA, true)
		mgr.Update()

		bState[1] = mgr.State(KeyA)

		mgr.Clear()
		mgr.Update()

		bState[2] = mgr.State(KeyA)

		mgr.Clear()
		mgr.SetState(KeyA, false)
		mgr.Update()

		bState[3] = mgr.State(KeyA)
	}
}

func BenchmarkActMgr_FilledSimulate(b *testing.B) {
	mgr := NewActManager()
	fillActMgr(mgr)

	mgr.Clear()
	mgr.Update()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		mgr.Clear()
		mgr.Update()

		bState[0] = mgr.State(KeyA)

		mgr.Clear()
		mgr.SetState(KeyA, true)
		mgr.Update()

		bState[1] = mgr.State(KeyA)

		mgr.Clear()
		mgr.Update()

		bState[2] = mgr.State(KeyA)

		mgr.Clear()
		mgr.SetState(KeyA, false)
		mgr.Update()

		bState[3] = mgr.State(KeyA)
	}
}

// Disabled for CI - Use to diganose
//func BenchmarkActMgr_CleanUpdate(b *testing.B) {
//	mgr := NewActMgr()
//
//	mgr.Clear()
//	mgr.Update()
//
//	b.ResetTimer()
//	for n := 0; n < b.N; n++ {
//		mgr.Clear()
//		mgr.Update()
//	}
//}

// Disabled for CI - Use to diganose
//func BenchmarkActMgr_FilledUpdate(b *testing.B) {
//	mgr := NewActMgr()
//	fillActMgr(mgr)
//
//	mgr.Clear()
//	mgr.Update()
//
//	b.ResetTimer()
//	for n := 0; n < b.N; n++ {
//		mgr.Clear()
//		mgr.Update()
//	}
//}

// Disabled for CI - Use to diganose
//func BenchmarkActMgr_CleanUpdateSet(b *testing.B) {
//	mgr := NewActMgr()
//
//	mgr.Clear()
//	mgr.Update()
//
//	b.ResetTimer()
//	for n := 0; n < b.N; n++ {
//		mgr.Clear()
//		mgr.SetState(KeyA, true)
//		mgr.Update()
//	}
//}

// Disabled for CI - Usedto diganose
//func BenchmarkActMgr_FilledUpdateSet(b *testing.B) {
//	mgr := NewActMgr()
//	fillActMgr(mgr)
//
//	mgr.Clear()
//	mgr.Update()
//
//	b.ResetTimer()
//	for n := 0; n < b.N; n++ {
//		mgr.Clear()
//		mgr.SetState(KeyA, true)
//		mgr.Update()
//	}
//}

// Disabled for CI - Use to diganose
//func BenchmarkActMgr_CleanUpdateState(b *testing.B) {
//	mgr := NewActMgr()
//
//	mgr.Clear()
//	mgr.Update()
//
//	b.ResetTimer()
//	for n := 0; n < b.N; n++ {
//		mgr.Clear()
//		mgr.Update()
//		bState[0] = mgr.State(KeyA)
//		bResult[0] = mgr.Idle(KeyA)
//		bResult[1] = mgr.Active(KeyA)
//		bResult[2] = mgr.JustIdle(KeyA)
//		bResult[3] = mgr.JustActive(KeyA)
//	}
//}

// Disabled for CI - Use to diganose
//func BenchmarkActMgr_FilledUpdateState(b *testing.B) {
//	mgr := NewActMgr()
//	fillActMgr(mgr)
//
//	mgr.Clear()
//	mgr.Update()
//
//	b.ResetTimer()
//	for n := 0; n < b.N; n++ {
//		mgr.Clear()
//		mgr.Update()
//		bState[0] = mgr.State(KeyA)
//		bResult[0] = mgr.Idle(KeyA)
//		bResult[1] = mgr.Active(KeyA)
//		bResult[2] = mgr.JustIdle(KeyA)
//		bResult[3] = mgr.JustActive(KeyA)
//	}
//}

// Utility function that fills the ActMgr with code states.
func fillActMgr(mgr *ActManager) {
	mgr.SetState(KeyD, false)
	mgr.SetState(KeyE, false)
	mgr.SetState(KeyF, false)
	mgr.SetState(KeyG, false)
	mgr.SetState(KeyH, false)
	mgr.SetState(KeyI, false)
	mgr.SetState(KeyJ, false)
	mgr.SetState(KeyK, false)
	mgr.SetState(KeyL, false)
	mgr.SetState(KeyM, false)
	mgr.SetState(KeyN, false)
	mgr.SetState(KeyO, false)
	mgr.SetState(KeyP, false)
	mgr.SetState(KeyQ, false)
	mgr.SetState(KeyR, false)
	mgr.SetState(KeyS, false)
	mgr.SetState(KeyT, false)
	mgr.SetState(KeyU, false)
	mgr.SetState(KeyV, false)
	mgr.SetState(KeyW, false)
	mgr.SetState(KeyX, false)
	mgr.SetState(KeyY, false)
	mgr.SetState(KeyZ, false)

	mgr.SetState(Key0, false)
	mgr.SetState(Key1, false)
	mgr.SetState(Key2, false)
	mgr.SetState(Key3, false)
	mgr.SetState(Key4, false)
	mgr.SetState(Key5, false)
	mgr.SetState(Key6, false)
	mgr.SetState(Key7, false)
	mgr.SetState(Key8, false)
	mgr.SetState(Key9, false)

	mgr.SetState(KeyF5, false)
	mgr.SetState(KeyF6, false)
	mgr.SetState(KeyF7, false)
	mgr.SetState(KeyF8, false)
	mgr.SetState(KeyF9, false)
	mgr.SetState(KeyF10, false)
	mgr.SetState(KeyF11, false)
	mgr.SetState(KeyF12, false)

	mgr.SetState(KeyPad0, false)
	mgr.SetState(KeyPad1, false)
	mgr.SetState(KeyPad2, false)
	mgr.SetState(KeyPad3, false)
	mgr.SetState(KeyPad4, false)
	mgr.SetState(KeyPad5, false)
	mgr.SetState(KeyPad6, false)
	mgr.SetState(KeyPad7, false)
	mgr.SetState(KeyPad8, false)
	mgr.SetState(KeyPad9, false)
}
