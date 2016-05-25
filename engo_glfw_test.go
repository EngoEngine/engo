//+build !netgo,!android

package engo

import (
    "testing"
    
    "github.com/github.com/stretchr/testify/assert"
)

func TestGameSize(t *testing.T) {
    assert.Equal(t, GameWidth(), gameWidth)
    assert.Equal(t, GameHeight(), gameHeight)    
}

func TestCursorPos(t *testing.T) {
    assert.Equal(t, CursorPos(), window.GetCursorPos())
}

func TestWindowSize(t *testing.T) {
    assert.Equal(t, WindowSize(), window.GetSize())
    assert.Equal(t, WindowWidth(), windowWidth)
    assert.Equal(t, WindowHeight(), windowHeight)
}

func TestCanvasSize(t *testing.T) {
    assert.Equal(t, CanvasWidth(), canvasWidth)
    assert.Equal(t, CanvasHeight(), canvasHeight)
}