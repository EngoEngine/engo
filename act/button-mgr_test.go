package act

import "testing"

func TestButtonMgr(t *testing.T) {
	amgr := NewActMgr()
	bmgr := NewButtonMgr(amgr)

	abtn := bmgr.SetNamed("Button A", KeyA, KeyB)
	bbtn := bmgr.SetNamed("Button B", KeyF3, KeyF4)
	cbtn := bmgr.SetNamed("Button C", MouseLeft, MouseRight)

	if abtn != bmgr.Id("Button A") {
		t.Error("Failed to verify id")
	}

	if bmgr.SetId(99, KeyC, KeyD) {
		t.Error("Set codes on to an invalid id ?")
	}

	if !bmgr.SetId(bbtn, KeyF1, KeyF2) {
		t.Error("Failed to set codes on a valid id")
	}

	amgr.Clear()
	amgr.Update()

	runButtonCheck("Init (0.A)", t, bmgr, abtn, StateIdle)
	runButtonCheck("Init (0.B)", t, bmgr, bbtn, StateIdle)
	runButtonCheck("Init (0.C)", t, bmgr, cbtn, StateIdle)

	amgr.Clear()
	amgr.SetState(KeyA, true)
	amgr.SetState(KeyF1, false)
	amgr.Update()
	amgr.SetState(MouseLeft, true)

	runButtonCheck("Pass (1.A)", t, bmgr, abtn, StateJustActive)
	runButtonCheck("Pass (1.B)", t, bmgr, bbtn, StateIdle)
	runButtonCheck("Pass (1.C)", t, bmgr, cbtn, StateIdle)

	amgr.Clear()
	amgr.Update()

	runButtonCheck("Pass (2.A)", t, bmgr, abtn, StateActive)
	runButtonCheck("Pass (2.B)", t, bmgr, bbtn, StateIdle)
	runButtonCheck("Pass (2.C)", t, bmgr, cbtn, StateJustActive)

	amgr.Clear()
	amgr.SetState(KeyA, false)
	amgr.SetState(KeyF1, true)
	amgr.Update()
	amgr.SetState(MouseLeft, false)

	runButtonCheck("Pass (3.A)", t, bmgr, abtn, StateJustIdle)
	runButtonCheck("Pass (3.B)", t, bmgr, bbtn, StateJustActive)
	runButtonCheck("Pass (3.C)", t, bmgr, cbtn, StateActive)

	amgr.Clear()
	amgr.Update()

	runButtonCheck("Pass (4.A)", t, bmgr, abtn, StateIdle)
	runButtonCheck("Pass (4.B)", t, bmgr, bbtn, StateActive)
	runButtonCheck("Pass (4.C)", t, bmgr, cbtn, StateJustIdle)

	amgr.Clear()
	amgr.Update()

	runButtonCheck("Pass (5.A)", t, bmgr, abtn, StateIdle)
	runButtonCheck("Pass (5.B)", t, bmgr, bbtn, StateActive)
	runButtonCheck("Pass (5.C)", t, bmgr, cbtn, StateIdle)

	amgr.Clear()
	amgr.SetState(KeyF1, false)
	amgr.Update()

	runButtonCheck("Pass (6.A)", t, bmgr, abtn, StateIdle)
	runButtonCheck("Pass (6.B)", t, bmgr, bbtn, StateJustIdle)
	runButtonCheck("Pass (6.C)", t, bmgr, cbtn, StateIdle)

	amgr.Clear()
	amgr.Update()

	runButtonCheck("Pass (7.A)", t, bmgr, abtn, StateIdle)
	runButtonCheck("Pass (7.B)", t, bmgr, bbtn, StateIdle)
	runButtonCheck("Pass (7.C)", t, bmgr, cbtn, StateIdle)

	amgr.Clear()
	amgr.SetState(KeyB, true)
	amgr.SetState(KeyF2, true)
	amgr.SetState(MouseRight, true)
	amgr.Update()

	runButtonCheck("Pass (8.A)", t, bmgr, abtn, StateJustActive)
	runButtonCheck("Pass (8.B)", t, bmgr, bbtn, StateJustActive)
	runButtonCheck("Pass (8.C)", t, bmgr, cbtn, StateJustActive)

	amgr.Clear()
	amgr.SetState(KeyB, false)
	amgr.SetState(KeyF1, true)
	amgr.SetState(MouseLeft, true)
	amgr.Update()

	runButtonCheck("Pass (9.A)", t, bmgr, abtn, StateJustIdle)
	runButtonCheck("Pass (9.B)", t, bmgr, bbtn, StateActive)
	runButtonCheck("Pass (9.C)", t, bmgr, cbtn, StateActive)

	amgr.Clear()
	amgr.SetState(KeyA, true)
	amgr.SetState(KeyB, true)
	amgr.SetState(KeyF2, false)
	amgr.SetState(MouseLeft, false)
	amgr.Update()

	runButtonCheck("Pass (10.A)", t, bmgr, abtn, StateJustActive)
	runButtonCheck("Pass (10.B)", t, bmgr, bbtn, StateActive)
	runButtonCheck("Pass (10.C)", t, bmgr, cbtn, StateActive)

	amgr.Clear()
	amgr.SetState(KeyA, false)
	amgr.SetState(MouseRight, false)
	amgr.Update()

	runButtonCheck("Pass (11.A)", t, bmgr, abtn, StateActive)
	runButtonCheck("Pass (11.B)", t, bmgr, bbtn, StateActive)
	runButtonCheck("Pass (11.C)", t, bmgr, cbtn, StateJustIdle)

	amgr.Clear()
	amgr.SetState(KeyB, false)
	amgr.SetState(KeyF1, false)
	amgr.Update()

	runButtonCheck("Pass (12.A)", t, bmgr, abtn, StateJustIdle)
	runButtonCheck("Pass (12.B)", t, bmgr, bbtn, StateJustIdle)
	runButtonCheck("Pass (12.C)", t, bmgr, cbtn, StateIdle)

	amgr.Clear()
	amgr.Update()

	runButtonCheck("Pass (13.A)", t, bmgr, abtn, StateIdle)
	runButtonCheck("Pass (13.B)", t, bmgr, bbtn, StateIdle)
	runButtonCheck("Pass (13.C)", t, bmgr, cbtn, StateIdle)
}

func runButtonCheck(msg string, t *testing.T, mgr *ButtonMgr, id uintptr, exp State) {
	if (StateIdle == exp) != mgr.Idle(id) {
		t.Error(msg, " - Invalid on: Idle")
	}
	if (StateActive == exp) != mgr.Active(id) {
		t.Error(msg, " - Invalid on: Active")
	}
	if (StateJustIdle == exp) != mgr.JustIdle(id) {
		t.Error(msg, " - Invalid on: Just Idle")
	}
	if (StateJustActive == exp) != mgr.JustActive(id) {
		t.Error(msg, " - Invalid on: Just Active")
	}
}

////////////////

func BenchmarkButtonMgr_CleanSimulate(b *testing.B) {
	amgr := NewActMgr()
	bmgr := NewButtonMgr(amgr)

	btn := bmgr.SetNamed("Button A", KeyA, KeyB)

	amgr.Clear()
	amgr.Update()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		amgr.Clear()
		amgr.Update()

		bResult[0] = bmgr.Idle(btn)

		amgr.Clear()
		amgr.SetState(KeyA, true)
		amgr.Update()

		bResult[1] = bmgr.Active(btn)

		amgr.Clear()
		amgr.Update()

		bResult[2] = bmgr.JustIdle(btn)

		amgr.Clear()
		amgr.SetState(KeyA, false)
		amgr.Update()

		bResult[3] = bmgr.JustActive(btn)
	}
}

func BenchmarkButtonMgr_FilledSimulate(b *testing.B) {
	amgr := NewActMgr()
	fillActMgr(amgr)
	bmgr := NewButtonMgr(amgr)

	btn := bmgr.SetNamed("Button A", KeyA, KeyB)

	amgr.Clear()
	amgr.Update()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		amgr.Clear()
		amgr.Update()

		bResult[0] = bmgr.Idle(btn)

		amgr.Clear()
		amgr.SetState(KeyA, true)
		amgr.Update()

		bResult[1] = bmgr.Active(btn)

		amgr.Clear()
		amgr.Update()

		bResult[2] = bmgr.JustIdle(btn)

		amgr.Clear()
		amgr.SetState(KeyA, false)
		amgr.Update()

		bResult[3] = bmgr.JustActive(btn)
	}
}
