package engo

import (
    "testing"
    
    "github.com/github.com/stretchr/testify/assert"
)

func TestSetScaleOnResize(t *testing.T) {
    SetScaleOnResize(true)
    assert.Equal(t, opts.ScaleOnResize, true)
    assert.Equal(t, ScaleOnResize(), true)
    SetScaleOnResize(false)
    assert.Equal(t, opts.ScaleOnResize, false)
    assert.Equal(t, ScaleOnResize(), false)
}

func TestSetOverrideCloseAction(t *testing.T) {
    SetOverrideCloseAction(true)
    assert.Equal(t, opts.OverrideCloseAction, true)
    SetOverrideCloseAction(false)
    assert.Equal(t, opts.OverrideCloseAction, false)
}

func TestSetFPSLimit(t *testing.T) {
    err := SetFPSLimit(60)
    assert.Nil(t, err)
    assert.Equal(t, opts.FPSLimit, 60)
    err = SetFPSLimit(-1)
    assert.NotNil(t, err)
}

func TestHeadless(t *testing.T) {
    SetHeadless(false)
    assert.Equal(t, opts.HeadlessMode, false)
    assert.Equal(t, Headless(), false)
    SetHeadless(true)
    assert.Equal(t, opts.HeadlessMode, true)
    assert.Equal(t, Headless)
}


func TestExit(t *testing.T) {
    assert.Equal(t, closeGame, false)
    Exit()
    assert.Equal(t, closeGame, true)
}