package engo

import "testing"

type ExpState struct {
	up       bool
	down     bool
	justUp   bool
	justDown bool
}

var result [12]ExpState

////////////////

var initList = [12]Key{
	A, B, C, D,
	F1, F2, F3, F4,
	Slash, NumTwo,
	Enter, ArrowLeft,
}

var initPass0 = [12]ExpState{
	ExpState{up: true, down: false, justUp: false, justDown: false},
	ExpState{up: true, down: false, justUp: false, justDown: false},
	ExpState{up: true, down: false, justUp: false, justDown: false},
	ExpState{up: true, down: false, justUp: false, justDown: false},
	ExpState{up: true, down: false, justUp: false, justDown: false},
	ExpState{up: true, down: false, justUp: false, justDown: false},
	ExpState{up: true, down: false, justUp: false, justDown: false},
	ExpState{up: true, down: false, justUp: false, justDown: false},
	ExpState{up: true, down: false, justUp: false, justDown: false},
	ExpState{up: true, down: false, justUp: false, justDown: false},
	ExpState{up: true, down: false, justUp: false, justDown: false},
	ExpState{up: true, down: false, justUp: false, justDown: false},
}

var initPass1 = [12]ExpState{
	ExpState{up: true, down: false, justUp: false, justDown: false},
	ExpState{up: false, down: false, justUp: false, justDown: true},
	ExpState{up: true, down: false, justUp: false, justDown: false},
	ExpState{up: false, down: false, justUp: false, justDown: true},
	ExpState{up: true, down: false, justUp: false, justDown: false},
	ExpState{up: false, down: false, justUp: false, justDown: true},
	ExpState{up: true, down: false, justUp: false, justDown: false},
	ExpState{up: false, down: false, justUp: false, justDown: true},
	ExpState{up: true, down: false, justUp: false, justDown: false},
	ExpState{up: false, down: false, justUp: false, justDown: true},
	ExpState{up: true, down: false, justUp: false, justDown: false},
	ExpState{up: false, down: false, justUp: false, justDown: true},
}

var initPass2 = [12]ExpState{
	ExpState{up: false, down: false, justUp: false, justDown: true},
	ExpState{up: false, down: true, justUp: false, justDown: false},
	ExpState{up: false, down: false, justUp: false, justDown: true},
	ExpState{up: false, down: true, justUp: false, justDown: false},
	ExpState{up: false, down: false, justUp: false, justDown: true},
	ExpState{up: false, down: true, justUp: false, justDown: false},
	ExpState{up: false, down: false, justUp: false, justDown: true},
	ExpState{up: false, down: true, justUp: false, justDown: false},
	ExpState{up: false, down: false, justUp: false, justDown: true},
	ExpState{up: false, down: true, justUp: false, justDown: false},
	ExpState{up: false, down: false, justUp: false, justDown: true},
	ExpState{up: false, down: true, justUp: false, justDown: false},
}

var initPass3 = [12]ExpState{
	ExpState{up: false, down: true, justUp: false, justDown: false},
	ExpState{up: false, down: true, justUp: false, justDown: false},
	ExpState{up: false, down: true, justUp: false, justDown: false},
	ExpState{up: false, down: true, justUp: false, justDown: false},
	ExpState{up: false, down: true, justUp: false, justDown: false},
	ExpState{up: false, down: true, justUp: false, justDown: false},
	ExpState{up: false, down: true, justUp: false, justDown: false},
	ExpState{up: false, down: true, justUp: false, justDown: false},
	ExpState{up: false, down: true, justUp: false, justDown: false},
	ExpState{up: false, down: true, justUp: false, justDown: false},
	ExpState{up: false, down: true, justUp: false, justDown: false},
	ExpState{up: false, down: true, justUp: false, justDown: false},
}

var initPass4 = [12]ExpState{
	ExpState{up: false, down: true, justUp: false, justDown: false},
	ExpState{up: false, down: false, justUp: true, justDown: false},
	ExpState{up: false, down: true, justUp: false, justDown: false},
	ExpState{up: false, down: false, justUp: true, justDown: false},
	ExpState{up: false, down: true, justUp: false, justDown: false},
	ExpState{up: false, down: false, justUp: true, justDown: false},
	ExpState{up: false, down: true, justUp: false, justDown: false},
	ExpState{up: false, down: false, justUp: true, justDown: false},
	ExpState{up: false, down: true, justUp: false, justDown: false},
	ExpState{up: false, down: false, justUp: true, justDown: false},
	ExpState{up: false, down: true, justUp: false, justDown: false},
	ExpState{up: false, down: false, justUp: true, justDown: false},
}

var initPass5 = [12]ExpState{
	ExpState{up: false, down: false, justUp: true, justDown: false},
	ExpState{up: true, down: false, justUp: false, justDown: false},
	ExpState{up: false, down: false, justUp: true, justDown: false},
	ExpState{up: true, down: false, justUp: false, justDown: false},
	ExpState{up: false, down: false, justUp: true, justDown: false},
	ExpState{up: true, down: false, justUp: false, justDown: false},
	ExpState{up: false, down: false, justUp: true, justDown: false},
	ExpState{up: true, down: false, justUp: false, justDown: false},
	ExpState{up: false, down: false, justUp: true, justDown: false},
	ExpState{up: true, down: false, justUp: false, justDown: false},
	ExpState{up: false, down: false, justUp: true, justDown: false},
	ExpState{up: true, down: false, justUp: false, justDown: false},
}

var initPass6 = [12]ExpState{
	ExpState{up: true, down: false, justUp: false, justDown: false},
	ExpState{up: true, down: false, justUp: false, justDown: false},
	ExpState{up: true, down: false, justUp: false, justDown: false},
	ExpState{up: true, down: false, justUp: false, justDown: false},
	ExpState{up: true, down: false, justUp: false, justDown: false},
	ExpState{up: true, down: false, justUp: false, justDown: false},
	ExpState{up: true, down: false, justUp: false, justDown: false},
	ExpState{up: true, down: false, justUp: false, justDown: false},
	ExpState{up: true, down: false, justUp: false, justDown: false},
	ExpState{up: true, down: false, justUp: false, justDown: false},
	ExpState{up: true, down: false, justUp: false, justDown: false},
	ExpState{up: true, down: false, justUp: false, justDown: false},
}

////////////////

func runInitList(msg string, t *testing.T, mgr *KeyManager, expect [12]ExpState) {
	for i, cd := range initList {
		exp := expect[i]
		ste := mgr.Get(cd)
		if exp.up != ste.Up() {
			t.Error(msg, " Invalid on: ", cd, " - Up")
		}
		if exp.down != ste.Down() {
			t.Error(msg, " Invalid on: ", cd, " - Down")
		}
		if exp.justUp != ste.JustReleased() {
			t.Error(msg, " Invalid on: ", cd, " - Just Up")
		}
		if exp.justDown != ste.JustPressed() {
			t.Error(msg, " Invalid on: ", cd, " - Just Down")
		}
	}
}

func TestKeyManager(t *testing.T) {
	mgr := NewKeyManager()
	runInitList("Init (0.0)", t, mgr, initPass0)

	// Empty update pass0
	mgr.update()
	runInitList("Pass (0.1)", t, mgr, initPass0)
	mgr.update()
	runInitList("Pass (0.2)", t, mgr, initPass0)
	mgr.update()
	runInitList("Pass (0.3)", t, mgr, initPass0)

	// Set uneven true pass1
	mgr.update()

	mgr.Set(initList[1], true)
	mgr.Set(initList[3], true)
	mgr.Set(initList[5], true)
	mgr.Set(initList[7], true)
	mgr.Set(initList[9], true)
	mgr.Set(initList[11], true)

	runInitList("Pass (1.0)", t, mgr, initPass1)

	// Set even true pass2
	mgr.update()

	mgr.Set(initList[0], true)
	mgr.Set(initList[2], true)
	mgr.Set(initList[4], true)
	mgr.Set(initList[6], true)
	mgr.Set(initList[8], true)
	mgr.Set(initList[10], true)

	runInitList("Pass (2.0)", t, mgr, initPass2)

	// Keeps state
	mgr.update()
	runInitList("Pass (3.0)", t, mgr, initPass3)
	mgr.update()
	runInitList("Pass (3.1)", t, mgr, initPass3)
	mgr.update()
	runInitList("Pass (3.2)", t, mgr, initPass3)
	mgr.update()
	runInitList("Pass (3.3)", t, mgr, initPass3)

	// Set uneven to false
	mgr.update()

	mgr.Set(initList[1], false)
	mgr.Set(initList[3], false)
	mgr.Set(initList[5], false)
	mgr.Set(initList[7], false)
	mgr.Set(initList[9], false)
	mgr.Set(initList[11], false)

	runInitList("Pass (4.0)", t, mgr, initPass4)

	// Set uneven to false
	mgr.update()

	mgr.Set(initList[0], false)
	mgr.Set(initList[2], false)
	mgr.Set(initList[4], false)
	mgr.Set(initList[6], false)
	mgr.Set(initList[8], false)
	mgr.Set(initList[10], false)

	runInitList("Pass (5.0)", t, mgr, initPass5)

	// Final keeps state
	mgr.update()
	runInitList("Pass (6.0)", t, mgr, initPass6)
	mgr.update()
	runInitList("Pass (6.1)", t, mgr, initPass6)
	mgr.update()
	runInitList("Pass (6.2)", t, mgr, initPass6)
	mgr.update()
	runInitList("Pass (6.3)", t, mgr, initPass6)
}

////////////////

func fillManager(mgr *KeyManager) {
	mgr.Set(initList[0], false)
	mgr.Set(initList[1], false)
	mgr.Set(initList[2], false)
	mgr.Set(initList[3], false)
	mgr.Set(initList[4], false)
	mgr.Set(initList[5], false)
	mgr.Set(initList[6], false)
	mgr.Set(initList[7], false)
	mgr.Set(initList[8], false)
	mgr.Set(initList[9], false)
	mgr.Set(initList[10], false)
	mgr.Set(initList[11], false)

	mgr.Set(D, false)
	mgr.Set(E, false)
	mgr.Set(F, false)
	mgr.Set(G, false)
	mgr.Set(H, false)
	mgr.Set(I, false)
	mgr.Set(J, false)
	mgr.Set(K, false)
	mgr.Set(L, false)
	mgr.Set(M, false)
	mgr.Set(N, false)
	mgr.Set(O, false)
	mgr.Set(P, false)
	mgr.Set(Q, false)
	mgr.Set(R, false)
	mgr.Set(S, false)

	mgr.Set(T, false)
	mgr.Set(U, false)
	mgr.Set(V, false)
	mgr.Set(W, false)
	mgr.Set(X, false)
	mgr.Set(Y, false)
	mgr.Set(Z, false)

	mgr.Set(Zero, false)
	mgr.Set(One, false)
	mgr.Set(Two, false)
	mgr.Set(Three, false)
	mgr.Set(Four, false)
	mgr.Set(Five, false)
	mgr.Set(Six, false)
	mgr.Set(Seven, false)
	mgr.Set(Eight, false)
	mgr.Set(Nine, false)

	mgr.Set(F5, false)
	mgr.Set(F6, false)
	mgr.Set(F7, false)
	mgr.Set(F8, false)
	mgr.Set(F9, false)
	mgr.Set(F10, false)
	mgr.Set(F11, false)
	mgr.Set(F12, false)

	mgr.Set(NumZero, false)
	mgr.Set(NumOne, false)
	mgr.Set(NumTwo, false)
	mgr.Set(NumThree, false)
	mgr.Set(NumFour, false)
	mgr.Set(NumFive, false)
	mgr.Set(NumSix, false)
	mgr.Set(NumSeven, false)
	mgr.Set(NumEight, false)
	mgr.Set(NumNine, false)
}

func checkInitListOptimal(b *testing.B, mgr *KeyManager) {
	for i, cd := range initList {
		st := mgr.Get(cd)
		result[i].up = st.Up()
		result[i].down = st.Down()
		result[i].justUp = st.JustReleased()
		result[i].justDown = st.JustPressed()
	}
}

// Benchmark optimal state checks
func BenchmarkKeyMgrCleanState(b *testing.B) {
	mgr := NewKeyManager()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		checkInitListOptimal(b, mgr)
	}
}

// Benchmark update with state checks to avoide optimizing.
func BenchmarkKeyMgrCleanUpdate(b *testing.B) {
	mgr := NewKeyManager()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		mgr.update()
		checkInitListOptimal(b, mgr)
	}
}

// Benchmark optimal state checks with keys inside the map.
func BenchmarkKeyMgrFilledState(b *testing.B) {
	mgr := NewKeyManager()
	fillManager(mgr)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		checkInitListOptimal(b, mgr)
	}
}

// Benchmark update with keys inside the map and state checks to avoide optimizing.
func BenchmarkKeyMgrFilledUpdate(b *testing.B) {
	mgr := NewKeyManager()
	fillManager(mgr)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		mgr.update()
		checkInitListOptimal(b, mgr)
	}
}

////////////////

func checkInitListSubOptimal(b *testing.B, mgr *KeyManager) {
	for i, cd := range initList {
		result[i].up = mgr.Get(cd).Up()
		result[i].down = mgr.Get(cd).Down()
		result[i].justUp = mgr.Get(cd).JustReleased()
		result[i].justDown = mgr.Get(cd).JustPressed()
	}
}

// Benchmark sub-optimal state checks
func BenchmarkKeyMgrCleanStateSubOptimal(b *testing.B) {
	mgr := NewKeyManager()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		checkInitListSubOptimal(b, mgr)
	}
}
