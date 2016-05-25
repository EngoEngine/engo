package engo

import (
	"testing"

	"github.com/github.com/stretchr/testify/assert"
)

func TestSetScaleOnResize(t *testing.T) {
	SetScaleOnResize(true)
	assert.True(t, opts.ScaleOnResize)
	assert.True(t, ScaleOnResize())
	SetScaleOnResize(false)
	assert.False(t, opts.ScaleOnResize)
	assert.False(t, ScaleOnResize())
}

func TestSetOverrideCloseAction(t *testing.T) {
	SetOverrideCloseAction(true)
	assert.True(t, opts.OverrideCloseAction)
	SetOverrideCloseAction(false)
	assert.True(t, opts.OverrideCloseAction)
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
	assert.False(t, opts.HeadlessMode)
	assert.False(t, Headless())
	SetHeadless(true)
	assert.True(t, opts.HeadlessMode)
	assert.Equal(t, Headless)
}

func TestExit(t *testing.T) {
	assert.False(t, closeGame)
	Exit()
	assert.True(t, closeGame)
}
