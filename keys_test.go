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
	KeyA, KeyB, KeyC, KeyD,
	KeyF1, KeyF2, KeyF3, KeyF4,
	KeySlash, KeyNumTwo,
	KeyEnter, KeyArrowLeft,
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

	mgr.Set(KeyD, false)
	mgr.Set(KeyE, false)
	mgr.Set(KeyF, false)
	mgr.Set(KeyG, false)
	mgr.Set(KeyH, false)
	mgr.Set(KeyI, false)
	mgr.Set(KeyJ, false)
	mgr.Set(KeyK, false)
	mgr.Set(KeyL, false)
	mgr.Set(KeyM, false)
	mgr.Set(KeyN, false)
	mgr.Set(KeyO, false)
	mgr.Set(KeyP, false)
	mgr.Set(KeyQ, false)
	mgr.Set(KeyR, false)
	mgr.Set(KeyS, false)

	mgr.Set(KeyT, false)
	mgr.Set(KeyU, false)
	mgr.Set(KeyV, false)
	mgr.Set(KeyW, false)
	mgr.Set(KeyX, false)
	mgr.Set(KeyY, false)
	mgr.Set(KeyZ, false)

	mgr.Set(KeyZero, false)
	mgr.Set(KeyOne, false)
	mgr.Set(KeyTwo, false)
	mgr.Set(KeyThree, false)
	mgr.Set(KeyFour, false)
	mgr.Set(KeyFive, false)
	mgr.Set(KeySix, false)
	mgr.Set(KeySeven, false)
	mgr.Set(KeyEight, false)
	mgr.Set(KeyNine, false)

	mgr.Set(KeyF5, false)
	mgr.Set(KeyF6, false)
	mgr.Set(KeyF7, false)
	mgr.Set(KeyF8, false)
	mgr.Set(KeyF9, false)
	mgr.Set(KeyF10, false)
	mgr.Set(KeyF11, false)
	mgr.Set(KeyF12, false)

	mgr.Set(KeyNumZero, false)
	mgr.Set(KeyNumOne, false)
	mgr.Set(KeyNumTwo, false)
	mgr.Set(KeyNumThree, false)
	mgr.Set(KeyNumFour, false)
	mgr.Set(KeyNumFive, false)
	mgr.Set(KeyNumSix, false)
	mgr.Set(KeyNumSeven, false)
	mgr.Set(KeyNumEight, false)
	mgr.Set(KeyNumNine, false)
}
