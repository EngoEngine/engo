package engo

import "testing"

type keyState struct {
	state    int
	up       bool
	down     bool
	justUp   bool
	justDown bool
}

// Keys used when testing
var keySimpleCfg = [12]Key{
	A, B, C, D,
	F1, F2, F3, F4,
	Slash, NumTwo,
	Enter, ArrowLeft,
}

// Expected key state @ pass 0
var initPass0 = [12]keyState{
	keyState{state: KeyStateUp},
	keyState{state: KeyStateUp},
	keyState{state: KeyStateUp},
	keyState{state: KeyStateUp},
	keyState{state: KeyStateUp},
	keyState{state: KeyStateUp},
	keyState{state: KeyStateUp},
	keyState{state: KeyStateUp},
	keyState{state: KeyStateUp},
	keyState{state: KeyStateUp},
	keyState{state: KeyStateUp},
	keyState{state: KeyStateUp},
}

// Expected key state @ pass 1
var initPass1 = [12]keyState{
	keyState{state: KeyStateUp},
	keyState{state: KeyStateJustDown},
	keyState{state: KeyStateUp},
	keyState{state: KeyStateJustDown},
	keyState{state: KeyStateUp},
	keyState{state: KeyStateJustDown},
	keyState{state: KeyStateUp},
	keyState{state: KeyStateJustDown},
	keyState{state: KeyStateUp},
	keyState{state: KeyStateJustDown},
	keyState{state: KeyStateUp},
	keyState{state: KeyStateJustDown},
}

// Expected key state @ pass 2
var initPass2 = [12]keyState{
	keyState{state: KeyStateJustDown},
	keyState{state: KeyStateDown},
	keyState{state: KeyStateJustDown},
	keyState{state: KeyStateDown},
	keyState{state: KeyStateJustDown},
	keyState{state: KeyStateDown},
	keyState{state: KeyStateJustDown},
	keyState{state: KeyStateDown},
	keyState{state: KeyStateJustDown},
	keyState{state: KeyStateDown},
	keyState{state: KeyStateJustDown},
	keyState{state: KeyStateDown},
}

// Expected key state @ pass 3
var initPass3 = [12]keyState{
	keyState{state: KeyStateDown},
	keyState{state: KeyStateDown},
	keyState{state: KeyStateDown},
	keyState{state: KeyStateDown},
	keyState{state: KeyStateDown},
	keyState{state: KeyStateDown},
	keyState{state: KeyStateDown},
	keyState{state: KeyStateDown},
	keyState{state: KeyStateDown},
	keyState{state: KeyStateDown},
	keyState{state: KeyStateDown},
	keyState{state: KeyStateDown},
}

// Expected key state @ pass 4
var initPass4 = [12]keyState{
	keyState{state: KeyStateDown},
	keyState{state: KeyStateJustUp},
	keyState{state: KeyStateDown},
	keyState{state: KeyStateJustUp},
	keyState{state: KeyStateDown},
	keyState{state: KeyStateJustUp},
	keyState{state: KeyStateDown},
	keyState{state: KeyStateJustUp},
	keyState{state: KeyStateDown},
	keyState{state: KeyStateJustUp},
	keyState{state: KeyStateDown},
	keyState{state: KeyStateJustUp},
}

// Expected key state @ pass 5
var initPass5 = [12]keyState{
	keyState{state: KeyStateJustUp},
	keyState{state: KeyStateUp},
	keyState{state: KeyStateJustUp},
	keyState{state: KeyStateUp},
	keyState{state: KeyStateJustUp},
	keyState{state: KeyStateUp},
	keyState{state: KeyStateJustUp},
	keyState{state: KeyStateUp},
	keyState{state: KeyStateJustUp},
	keyState{state: KeyStateUp},
	keyState{state: KeyStateJustUp},
	keyState{state: KeyStateUp},
}

// Expected key state @ pass 6
var initPass6 = [12]keyState{
	keyState{state: KeyStateUp},
	keyState{state: KeyStateUp},
	keyState{state: KeyStateUp},
	keyState{state: KeyStateUp},
	keyState{state: KeyStateUp},
	keyState{state: KeyStateUp},
	keyState{state: KeyStateUp},
	keyState{state: KeyStateUp},
	keyState{state: KeyStateUp},
	keyState{state: KeyStateUp},
	keyState{state: KeyStateUp},
	keyState{state: KeyStateUp},
}

// Checks the state of keys in the configuration against the expected state.
func runKeyChecks(msg string, t *testing.T, mgr *KeyManager, expect [12]keyState) {
	for i, cd := range keySimpleCfg {
		exp := expect[i]
		sto := mgr.Get(cd)

		if exp.state != sto.State() {
			t.Error(msg, " Invalid on: ", cd, " - State")
		}
		if (KeyStateUp == exp.state) != sto.Up() {
			t.Error(msg, " Invalid on: ", cd, " - Up")
		}
		if (KeyStateDown == exp.state) != sto.Down() {
			t.Error(msg, " Invalid on: ", cd, " - Down")
		}
		if (KeyStateJustUp == exp.state) != sto.JustReleased() {
			t.Error(msg, " Invalid on: ", cd, " - Just Up")
		}
		if (KeyStateJustDown == exp.state) != sto.JustPressed() {
			t.Error(msg, " Invalid on: ", cd, " - Just Down")
		}
	}
}

// This test sets up a key manager and changes the states a few times, every
// time the state changes the results are checked against the expected outcome.
func TestKeyManager(t *testing.T) {
	mgr := NewKeyManager()
	runKeyChecks("Init (0.0)", t, mgr, initPass0)

	// Empty update pass0
	mgr.update()
	runKeyChecks("Pass (0.1)", t, mgr, initPass0)
	mgr.update()
	runKeyChecks("Pass (0.2)", t, mgr, initPass0)
	mgr.update()
	runKeyChecks("Pass (0.3)", t, mgr, initPass0)

	// Set uneven true pass1
	mgr.update()

	mgr.Set(keySimpleCfg[1], true)
	mgr.Set(keySimpleCfg[3], true)
	mgr.Set(keySimpleCfg[5], true)
	mgr.Set(keySimpleCfg[7], true)
	mgr.Set(keySimpleCfg[9], true)
	mgr.Set(keySimpleCfg[11], true)

	runKeyChecks("Pass (1.0)", t, mgr, initPass1)

	// Set even true pass2
	mgr.update()

	mgr.Set(keySimpleCfg[0], true)
	mgr.Set(keySimpleCfg[2], true)
	mgr.Set(keySimpleCfg[4], true)
	mgr.Set(keySimpleCfg[6], true)
	mgr.Set(keySimpleCfg[8], true)
	mgr.Set(keySimpleCfg[10], true)

	runKeyChecks("Pass (2.0)", t, mgr, initPass2)

	// Keeps state
	mgr.update()
	runKeyChecks("Pass (3.0)", t, mgr, initPass3)
	mgr.update()
	runKeyChecks("Pass (3.1)", t, mgr, initPass3)
	mgr.update()
	runKeyChecks("Pass (3.2)", t, mgr, initPass3)

	// Set uneven to false
	mgr.update()

	mgr.Set(keySimpleCfg[1], false)
	mgr.Set(keySimpleCfg[3], false)
	mgr.Set(keySimpleCfg[5], false)
	mgr.Set(keySimpleCfg[7], false)
	mgr.Set(keySimpleCfg[9], false)
	mgr.Set(keySimpleCfg[11], false)

	runKeyChecks("Pass (4.0)", t, mgr, initPass4)

	// Set uneven to false
	mgr.update()

	mgr.Set(keySimpleCfg[0], false)
	mgr.Set(keySimpleCfg[2], false)
	mgr.Set(keySimpleCfg[4], false)
	mgr.Set(keySimpleCfg[6], false)
	mgr.Set(keySimpleCfg[8], false)
	mgr.Set(keySimpleCfg[10], false)

	runKeyChecks("Pass (5.0)", t, mgr, initPass5)

	// Final keeps state
	mgr.update()
	runKeyChecks("Pass (6.0)", t, mgr, initPass6)
	mgr.update()
	runKeyChecks("Pass (6.1)", t, mgr, initPass6)
	mgr.update()
	runKeyChecks("Pass (6.2)", t, mgr, initPass6)
}

// Used to store results when benchmarking.
var keyResult [12]keyState

// Fast way to grab current key state, checks and store them externaly.
func checkKeyConfigOptimal(b *testing.B, mgr *KeyManager) {
	for i, cd := range keySimpleCfg {
		st := mgr.Get(cd)
		keyResult[i].up = st.Up()
		keyResult[i].down = st.Down()
		keyResult[i].justUp = st.JustReleased()
		keyResult[i].justDown = st.JustPressed()
	}
}

// Benchmark state checks with a clean manager.
func BenchmarkKeyMgr_CleanState(b *testing.B) {
	mgr := NewKeyManager()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		checkKeyConfigOptimal(b, mgr)
	}
}

// Benchmark state checks with a full key manager.
func BenchmarkKeyMgr_FilledState(b *testing.B) {
	mgr := NewKeyManager()
	keyFillManager(mgr)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		checkKeyConfigOptimal(b, mgr)
	}
}

// Benchmark update with a clean manager and state checks.
func BenchmarkKeyMgr_CleanUpdate(b *testing.B) {
	mgr := NewKeyManager()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		mgr.update()
		checkKeyConfigOptimal(b, mgr)
	}
}

// Benchmark update with a full key manager and state checks.
func BenchmarkKeyMgr_FilledUpdate(b *testing.B) {
	mgr := NewKeyManager()
	keyFillManager(mgr)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		mgr.update()
		checkKeyConfigOptimal(b, mgr)
	}
}

// Slow way to check key state, checks and store them externaly.
func checkKeyConfigSubOptimal(b *testing.B, mgr *KeyManager) {
	for i, cd := range keySimpleCfg {
		keyResult[i].up = mgr.Get(cd).Up()
		keyResult[i].down = mgr.Get(cd).Down()
		keyResult[i].justUp = mgr.Get(cd).JustReleased()
		keyResult[i].justDown = mgr.Get(cd).JustPressed()
	}
}

// Benchmark sub-optimal state checks on a clean key manager.
func BenchmarkKeyMgr_CleanSubOptimal(b *testing.B) {
	mgr := NewKeyManager()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		checkKeyConfigSubOptimal(b, mgr)
	}
}

// Benchmark sub-optimal state checks on a filled key manager.
func BenchmarkKeyMgr_FilledSubOptimal(b *testing.B) {
	mgr := NewKeyManager()
	keyFillManager(mgr)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		checkKeyConfigSubOptimal(b, mgr)
	}
}

// Utility function that fills the KeyManager with key states.
func keyFillManager(mgr *KeyManager) {
	mgr.Set(keySimpleCfg[0], false)
	mgr.Set(keySimpleCfg[1], false)
	mgr.Set(keySimpleCfg[2], false)
	mgr.Set(keySimpleCfg[3], false)
	mgr.Set(keySimpleCfg[4], false)
	mgr.Set(keySimpleCfg[5], false)
	mgr.Set(keySimpleCfg[6], false)
	mgr.Set(keySimpleCfg[7], false)
	mgr.Set(keySimpleCfg[8], false)
	mgr.Set(keySimpleCfg[9], false)
	mgr.Set(keySimpleCfg[10], false)
	mgr.Set(keySimpleCfg[11], false)

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
