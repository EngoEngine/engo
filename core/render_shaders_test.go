package core

import (
	"testing"

	"engo.io/ecs"
	"engo.io/engo"
	"github.com/stretchr/testify/assert"
)

type defaultScene struct{}

func (defaultScene) Type() string { return "defaultScene" }
func (defaultScene) Preload()     {}
func (defaultScene) Setup(w *ecs.World) {
	w.AddSystem(&RenderSystem{})
}

func TestShaders(t *testing.T) {
	engo.CreateWindow("", 100, 100, false, 1)

	w := &ecs.World{}
	w.AddSystem(&cameraSystem{})
	err := initShaders(w)
	assert.NoError(t, err)
}
