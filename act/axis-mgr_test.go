package act

import "testing"

func TestAxisMgr(t *testing.T) {
	amgr := NewActMgr()
	bmgr := NewAxisMgr(amgr)

	aaxi := bmgr.SetNamed("Axis A", AxisPair{Min: KeyA, Max: KeyB})
	baxi := bmgr.SetNamed("Axis B", AxisPair{Min: KeyF3, Max: KeyF4})
	caxi := bmgr.SetNamed("Axis C",
		AxisPair{Min: KeyPad0, Max: KeyPad1},
		AxisPair{Min: MouseLeft, Max: MouseRight},
	)

	if aaxi != bmgr.Id("Axis A") {
		t.Error("Failed to verify id")
	}

	if bmgr.SetId(99, AxisPair{Min: KeyC, Max: KeyD}) {
		t.Error("Set codes on to an invalid id ?")
	}

	if !bmgr.SetId(baxi, AxisPair{Min: KeyF1, Max: KeyF2}) {
		t.Error("Failed to set codes on a valid id")
	}

	amgr.Clear()
	amgr.Update()

	runAxisCheck("Init (0.A)", t, bmgr, aaxi, StateIdle, StateIdle, StateIdle, 0.0)
	runAxisCheck("Init (0.B)", t, bmgr, baxi, StateIdle, StateIdle, StateIdle, 0.0)
	runAxisCheck("Init (0.C)", t, bmgr, caxi, StateIdle, StateIdle, StateIdle, 0.0)

	amgr.Clear()
	amgr.SetState(KeyA, true)
	amgr.SetState(KeyF1, false)
	amgr.Update()
	amgr.SetState(MouseRight, true)

	runAxisCheck("Pass (1.A)", t, bmgr, aaxi, StateJustActive, StateIdle, StateJustActive, 0.0)
	runAxisCheck("Pass (1.B)", t, bmgr, baxi, StateIdle, StateIdle, StateIdle, 0.0)
	runAxisCheck("Pass (1.C)", t, bmgr, caxi, StateIdle, StateIdle, StateIdle, 0.0)

	amgr.Clear()
	amgr.Update()

	runAxisCheck("Pass (2.A)", t, bmgr, aaxi, StateActive, StateIdle, StateActive, -1.0)
	runAxisCheck("Pass (2.B)", t, bmgr, baxi, StateIdle, StateIdle, StateIdle, 0.0)
	runAxisCheck("Pass (2.C)", t, bmgr, caxi, StateIdle, StateJustActive, StateJustActive, 0.0)

	amgr.Clear()
	amgr.SetState(KeyA, false)
	amgr.SetState(KeyF1, true)
	amgr.Update()
	amgr.SetState(MouseRight, false)

	runAxisCheck("Pass (3.A)", t, bmgr, aaxi, StateJustIdle, StateIdle, StateJustIdle, 0.0)
	runAxisCheck("Pass (3.B)", t, bmgr, baxi, StateJustActive, StateIdle, StateJustActive, 0.0)
	runAxisCheck("Pass (3.C)", t, bmgr, caxi, StateIdle, StateActive, StateActive, 1.0)

	amgr.Clear()
	amgr.Update()

	runAxisCheck("Pass (4.A)", t, bmgr, aaxi, StateIdle, StateIdle, StateIdle, 0.0)
	runAxisCheck("Pass (4.B)", t, bmgr, baxi, StateActive, StateIdle, StateActive, -1.0)
	runAxisCheck("Pass (4.C)", t, bmgr, caxi, StateIdle, StateJustIdle, StateJustIdle, 0.0)

	amgr.Clear()
	amgr.Update()

	runAxisCheck("Pass (5.A)", t, bmgr, aaxi, StateIdle, StateIdle, StateIdle, 0.0)
	runAxisCheck("Pass (5.B)", t, bmgr, baxi, StateActive, StateIdle, StateActive, -1.0)
	runAxisCheck("Pass (5.C)", t, bmgr, caxi, StateIdle, StateIdle, StateIdle, 0.0)

	amgr.Clear()
	amgr.SetState(KeyA, true)
	amgr.SetState(KeyB, true)
	amgr.SetState(KeyF2, true)
	amgr.SetState(KeyPad1, true)
	amgr.Update()

	runAxisCheck("Pass (6.A)", t, bmgr, aaxi, StateJustActive, StateJustActive, StateJustActive, 0.0)
	runAxisCheck("Pass (6.B)", t, bmgr, baxi, StateActive, StateJustActive, StateActive, -1.0)
	runAxisCheck("Pass (6.C)", t, bmgr, caxi, StateIdle, StateJustActive, StateJustActive, 0.0)

	amgr.Clear()
	amgr.SetState(MouseLeft, true)
	amgr.Update()

	runAxisCheck("Pass (7.A)", t, bmgr, aaxi, StateActive, StateActive, StateActive, 0.0)
	runAxisCheck("Pass (7.B)", t, bmgr, baxi, StateActive, StateActive, StateActive, 0.0)
	runAxisCheck("Pass (7.C)", t, bmgr, caxi, StateJustActive, StateActive, StateActive, 1.0)

	amgr.Clear()
	amgr.SetState(KeyA, false)
	amgr.SetState(KeyF2, false)
	amgr.Update()

	runAxisCheck("Pass (8.A)", t, bmgr, aaxi, StateJustIdle, StateActive, StateActive, 1.0)
	runAxisCheck("Pass (8.B)", t, bmgr, baxi, StateActive, StateJustIdle, StateActive, -1.0)
	runAxisCheck("Pass (8.C)", t, bmgr, caxi, StateActive, StateActive, StateActive, 0.0)

	amgr.Clear()
	amgr.SetState(KeyB, false)
	amgr.SetState(KeyF1, false)
	amgr.SetState(KeyPad1, false)
	amgr.SetState(MouseLeft, false)
	amgr.Update()

	runAxisCheck("Pass (9.A)", t, bmgr, aaxi, StateIdle, StateJustIdle, StateJustIdle, 0.0)
	runAxisCheck("Pass (9.B)", t, bmgr, baxi, StateJustIdle, StateIdle, StateJustIdle, 0.0)
	runAxisCheck("Pass (9.C)", t, bmgr, caxi, StateJustIdle, StateJustIdle, StateJustIdle, 0.0)

	amgr.Clear()
	amgr.Update()

	runAxisCheck("Pass (10.A)", t, bmgr, aaxi, StateIdle, StateIdle, StateIdle, 0.0)
	runAxisCheck("Pass (10.B)", t, bmgr, baxi, StateIdle, StateIdle, StateIdle, 0.0)
	runAxisCheck("Pass (10.C)", t, bmgr, caxi, StateIdle, StateIdle, StateIdle, 0.0)
}

func runAxisCheck(msg string, t *testing.T, mgr *AxisMgr, id uintptr, min State, max State, exp State, val float32) {
	runMinAxisCheck(msg, t, mgr, id, min)
	runMaxAxisCheck(msg, t, mgr, id, max)
	runValueAxisCheck(msg, t, mgr, id, exp, val)
}

func runMinAxisCheck(msg string, t *testing.T, mgr *AxisMgr, id uintptr, exp State) {
	if (StateIdle == exp) != mgr.MinIdle(id) {
		t.Error(msg, " - Min - Invalid on: Idle")
	}
	if (StateActive == exp) != mgr.MinActive(id) {
		t.Error(msg, " - Min - Invalid on: Active")
	}
	if (StateJustIdle == exp) != mgr.MinJustIdle(id) {
		t.Error(msg, " - Min - Invalid on: Just Idle")
	}
	if (StateJustActive == exp) != mgr.MinJustActive(id) {
		t.Error(msg, " - Min - Invalid on: Just Active")
	}
}

func runMaxAxisCheck(msg string, t *testing.T, mgr *AxisMgr, id uintptr, exp State) {
	if (StateIdle == exp) != mgr.MaxIdle(id) {
		t.Error(msg, " - Max - Invalid on: Idle")
	}
	if (StateActive == exp) != mgr.MaxActive(id) {
		t.Error(msg, " - Max - Invalid on: Active")
	}
	if (StateJustIdle == exp) != mgr.MaxJustIdle(id) {
		t.Error(msg, " - Max - Invalid on: Just Idle")
	}
	if (StateJustActive == exp) != mgr.MaxJustActive(id) {
		t.Error(msg, " - Max - Invalid on: Just Active")
	}
}

func runValueAxisCheck(msg string, t *testing.T, mgr *AxisMgr, id uintptr, exp State, val float32) {
	if val != mgr.Value(id) {
		t.Error(msg, " - Val - Invalid value")
	}
	if (StateIdle == exp) != mgr.Idle(id) {
		t.Error(msg, " - Val - Invalid on: Idle")
	}
	if (StateActive == exp) != mgr.Active(id) {
		t.Error(msg, " - Val - Invalid on: Active")
	}
	if (StateJustIdle == exp) != mgr.JustIdle(id) {
		t.Error(msg, " - Val - Invalid on: Just Idle")
	}
	if (StateJustActive == exp) != mgr.JustActive(id) {
		t.Error(msg, " - Val-  Invalid on: Just Active")
	}
}

////////////////

func BenchmarkAxisMgr_CleanSimulate(b *testing.B) {
	amgr := NewActMgr()
	bmgr := NewAxisMgr(amgr)

	axi := bmgr.SetNamed("Axis A", AxisPair{Min: KeyA, Max: KeyB})

	amgr.Clear()
	amgr.Update()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		amgr.Clear()
		amgr.Update()

		bResult[0] = bmgr.Idle(axi)

		amgr.Clear()
		amgr.SetState(KeyA, true)
		amgr.Update()

		bResult[1] = bmgr.Active(axi)

		amgr.Clear()
		amgr.Update()

		bResult[2] = bmgr.JustIdle(axi)

		amgr.Clear()
		amgr.SetState(KeyA, false)
		amgr.Update()

		bResult[3] = bmgr.JustActive(axi)
	}
}

func BenchmarkAxisMgr_FilledSimulate(b *testing.B) {
	amgr := NewActMgr()
	fillActMgr(amgr)
	bmgr := NewAxisMgr(amgr)

	axi := bmgr.SetNamed("Axis A", AxisPair{Min: KeyA, Max: KeyB})

	amgr.Clear()
	amgr.Update()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		amgr.Clear()
		amgr.Update()

		bResult[0] = bmgr.Idle(axi)

		amgr.Clear()
		amgr.SetState(KeyA, true)
		amgr.Update()

		bResult[1] = bmgr.Active(axi)

		amgr.Clear()
		amgr.Update()

		bResult[2] = bmgr.JustIdle(axi)

		amgr.Clear()
		amgr.SetState(KeyA, false)
		amgr.Update()

		bResult[3] = bmgr.JustActive(axi)
	}
}
