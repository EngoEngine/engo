package core

import (
	"testing"

	"engo.io/ecs"
	"engo.io/engo"
	"github.com/stretchr/testify/assert"
)

func TestShaders(t *testing.T) {
	engo.SetHeadless(true)
	engo.CreateWindow("", 100, 100, false, 1)
	w := &ecs.World{}
	w.AddSystem(&cameraSystem{})
	err := initShaders(w)
	assert.NoError(t, err)
}
