package engo

import (
	"bytes"
	"log"
	"strings"
	"testing"
	"time"

	"engo.io/ecs"
)

type testScene struct{}

func (*testScene) Preload() {}

func (t *testScene) Setup(w *ecs.World) {}

func (*testScene) Type() string { return "testScene" }

// The tests for engo.go all have to use the headless option. Non-headless stuff is not
// testable via the cl only, and those are taken care of by building the demos via Travis CI
func TestRunHeadlessNoRunDefaults(t *testing.T) {
	Run(RunOptions{
		NoRun:        true,
		HeadlessMode: true,
	}, &testScene{})

	if opts.FPSLimit != 60 {
		t.Error("FPSLimit was not defaulted to 60")
	}

	if opts.MSAA != 1 {
		t.Error("MSAA was not defaulted to 1")
	}

	if opts.AssetsRoot != "assets" {
		t.Error("AssetsRoot was not defaulted to assets")
	}

	if opts.Width != 800 {
		t.Error("Width was not defaulted to 800")
	}

	if opts.Height != 800 {
		t.Error("Height was not defaulted to 800")
	}
}

func TestSetScaleOnResize(t *testing.T) {
	Run(RunOptions{
		HeadlessMode: true,
		NoRun:        true,
	}, &testScene{})
	if opts.ScaleOnResize {
		t.Error("ScaleOnResize didn't default to false.")
	}
	SetScaleOnResize(true)
	if !opts.ScaleOnResize {
		t.Error("SetScaleOnResize didn't set properly.")
	}
}

func TestSetOverrideCloseAction(t *testing.T) {
	Run(RunOptions{
		HeadlessMode: true,
		NoRun:        true,
	}, &testScene{})
	if opts.OverrideCloseAction {
		t.Error("OverrideCloseAction didn't default to false.")
	}
	SetOverrideCloseAction(true)
	if !opts.OverrideCloseAction {
		t.Error("SetOverrideCloseAction didn't set properly.")
	}
}

func TestHeadless(t *testing.T) {
	Run(RunOptions{
		HeadlessMode: true,
		NoRun:        true,
	}, &testScene{})
	if Headless() != opts.HeadlessMode {
		t.Error("Headless didn't return the proper value.")
	}
}

func TestScaleOnResize(t *testing.T) {
	Run(RunOptions{
		HeadlessMode: true,
		NoRun:        true,
	}, &testScene{})
	if ScaleOnResize() != opts.ScaleOnResize {
		t.Error("ScaleOnResize didn't return the proper value.")
	}
}

func TestGameWidthHeight(t *testing.T) {
	Run(RunOptions{
		HeadlessMode: true,
		NoRun:        true,
		Width:        100,
		Height:       50,
	}, &testScene{})
	if GameWidth() != 100 {
		t.Error("Width didn't return the proper value.")
	}
	if GameHeight() != 50 {
		t.Error("Height didn't return the proper value.")
	}
}

func TestSetFPSLimit(t *testing.T) {
	Run(RunOptions{
		HeadlessMode: true,
		NoRun:        true,
	}, &testScene{})
	SetFPSLimit(5)
	if opts.FPSLimit != 5 {
		t.Error("SetFPSLimit didn't set properly.")
	}

	expected := "FPS Limit out of bounds. Requires > 0"
	if err := SetFPSLimit(-5); err == nil {
		t.Error("Error wasn't recieved when SetFPSLimit was set to a negative number.")
	} else if err.Error() != expected {
		t.Errorf("Wrong error recieved when SetFPSLimit was set to a negative number. want %v, got %v", expected, err.Error())
	}
}

func TestRunNegativeMSAAPanic(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Error("No panic when MSAA was set to -5")
		}
		if r != "MSAA has to be greater or equal to 0" {
			t.Errorf("Wrong panic when MSAA was set to -5, got: %v", r)
		}
	}()
	Run(RunOptions{
		NoRun:        true,
		HeadlessMode: true,
		MSAA:         -5,
	}, &testScene{})
}

func TestRunStandardInputs(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)

	Run(RunOptions{
		NoRun:          true,
		HeadlessMode:   true,
		StandardInputs: true,
	}, &testScene{})

	expected := "Using standard inputs\n"
	if !strings.HasSuffix(buf.String(), expected) {
		t.Error("setting standard inputs did not write expected output to log")
	}
}

type testRunScene struct {
	updates int
}

func (*testRunScene) Preload() {}

func (t *testRunScene) Setup(w *ecs.World) {
	w.AddSystem(&testUpdate{updates: t.updates})
}

func (*testRunScene) Type() string { return "testRunScene" }

type testUpdate struct {
	updates, current int
}

func (*testUpdate) Remove(ecs.BasicEntity) {}

func (t *testUpdate) Update(float32) {
	t.current++
	if t.current >= t.updates {
		Exit()
	}
}

// This test tests running headless but also letting it go into the runloop
func TestRunHeadless(t *testing.T) {
	testChan := make(chan struct{})
	go func() {
		Run(RunOptions{
			HeadlessMode: true,
		}, &testRunScene{1})
		testChan <- struct{}{}
	}()
	select {
	case <-testChan:
		return
	case <-time.After(1 * time.Second):
		t.Error("Timed out while waiting for Headless Run to return from loop. Exit wasn't called within 1 second.")
	}
}

func TestOverrideCloseAction(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	opts.OverrideCloseAction = true
	expected := "[WARNING] default close action set to false, please make sure you manually handle this\n"

	closeEvent()

	if !strings.HasSuffix(buf.String(), expected) {
		t.Error("calling closeEvent with Override set did not write expected output to log")
	}
}
