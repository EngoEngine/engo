package common

import (
	"bytes"
	"log"
	"strings"
	"testing"

	"engo.io/ecs"
	"engo.io/engo"
	"github.com/stretchr/testify/assert"
)

var cam *CameraSystem

func initialize() {
	engo.Mailbox = &engo.MessageManager{}
	CameraBounds = engo.AABB{Min: engo.Point{X: 0, Y: 0}, Max: engo.Point{X: 300, Y: 300}}
	engo.SetGlobalScale(engo.Point{X: 1, Y: 1})
	w := &ecs.World{}

	cam = &CameraSystem{}
	cam.New(w)
}

type CameraTestScene struct{}

func (*CameraTestScene) Preload()         {}
func (*CameraTestScene) Setup(*ecs.World) {}
func (*CameraTestScene) Type() string     { return "CameraTestScene" }

func TestCameraMoveX(t *testing.T) {
	initialize()

	currentX := cam.X()

	cam.moveY(0.50)
	assert.Equal(t, cam.X(), currentX, "Moving in the Y direction shouldn't change the X location")

	cam.moveY(-1)
	assert.Equal(t, cam.X(), currentX, "Moving in the Y direction shouldn't change the X location")

	cam.zoom(0.50)
	assert.Equal(t, cam.X(), currentX, "Zooming in should not change the X location")

	cam.zoom(-1)
	assert.Equal(t, cam.X(), currentX, "Zooming out should not change the X location")

	cam.moveX(10)
	assert.Equal(t, cam.X(), currentX+10, "Moving by 10 units, should have moved the camera by 10 units")

	cam.moveX(-10)
	assert.Equal(t, cam.X(), currentX, "Moving by -10 units, should have moved the camera back by 10 units")

	cam.moveX(305)
	assert.Equal(t, cam.X(), CameraBounds.Max.X, "Moving too many unit, should have moved the camera to the maximum")

	cam.moveX(-305)
	assert.Equal(t, cam.X(), CameraBounds.Min.X, "Moving too many units back, should have moved the camera to the minimum")
}

func TestCameraMoveY(t *testing.T) {
	initialize()

	currentY := cam.Y()

	cam.moveX(0.50)
	assert.Equal(t, cam.Y(), currentY, "Moving in the X direction shouldn't change the Y location")

	cam.moveX(-1)
	assert.Equal(t, cam.Y(), currentY, "Moving in the X direction shouldn't change the Y location")

	cam.zoom(0.50)
	assert.Equal(t, cam.Y(), currentY, "Zooming in should not change the Y location")

	cam.zoom(-1)
	assert.Equal(t, cam.Y(), currentY, "Zooming out should not change the Y location")

	cam.moveY(10)
	assert.Equal(t, cam.Y(), currentY+10, "Moving by 10 units, should have moved the camera by 10 units")

	cam.moveY(-10)
	assert.Equal(t, cam.Y(), currentY, "Moving by -10 units, should have moved the camera back by 10 units")

	cam.moveY(305)
	assert.Equal(t, cam.Y(), CameraBounds.Max.Y, "Moving too many unit, should have moved the camera to the maximum")

	cam.moveY(-305)
	assert.Equal(t, cam.Y(), CameraBounds.Min.Y, "Moving too many units back, should have moved the camera to the minimum")
}

func TestCameraZoom(t *testing.T) {
	initialize()

	currentZ := cam.Z()

	cam.moveX(0.5)
	assert.Equal(t, cam.Z(), currentZ, "Moving in the X direction shouldn't change the zoom level")

	cam.moveX(-1)
	assert.Equal(t, cam.Z(), currentZ, "Moving in the X direction shouldn't change the zoom level")

	cam.moveY(0.5)
	assert.Equal(t, cam.Z(), currentZ, "Moving in the Y direction shouldn't change the zoom level")

	cam.moveY(-1)
	assert.Equal(t, cam.Z(), currentZ, "Moving in the Y direction shouldn't change the zoom level")

	cam.zoom(0.5)
	assert.Equal(t, cam.Z(), currentZ+0.5, "Should be zoomed in an additional 0.5")

	cam.zoom(-0.5)
	assert.Equal(t, cam.Z(), currentZ, "Should be zoomed out to the starting zoom level")

	cam.zoom(-1000)
	assert.Equal(t, cam.Z(), MinZoom, "Should be zoomed out to the minimum zoom level")

	cam.zoom(1000)
	assert.Equal(t, cam.Z(), MaxZoom, "Should be zoomed in to the maximum zoom level")
}

func TestCameraMoveToX(t *testing.T) {
	initialize()

	currentX := cam.X()

	cam.moveToX(currentX + 5)
	assert.Equal(t, cam.X(), currentX+5, "Moving to current + 5 should get us to current + 5")

	cam.moveToX(600)
	assert.Equal(t, cam.X(), CameraBounds.Max.X, "Moving to a location out of bounds, should get us to the maximum")

	cam.moveToX(-10)
	assert.Equal(t, cam.X(), CameraBounds.Min.X, "Moving to a location out of bounds, should get us to the minimum")
}

func TestCameraMoveToY(t *testing.T) {
	initialize()

	currentY := cam.Y()

	cam.moveToY(currentY + 5)
	assert.Equal(t, cam.Y(), currentY+5, "Moving to current + 5 should get us to current + 5")

	cam.moveToY(600)
	assert.Equal(t, cam.Y(), CameraBounds.Max.Y, "Moving to a location out of bounds, should get us to the maximum")

	cam.moveToY(-10)
	assert.Equal(t, cam.Y(), CameraBounds.Min.Y, "Moving to a location out of bounds, should get us to the minimum")
}

func TestCameraZoomTo(t *testing.T) {
	initialize()

	currentZ := cam.Z()

	cam.zoomTo(currentZ + 1)
	assert.Equal(t, cam.Z(), currentZ+1, "Zooming to current + 1 should get us to current + 1")

	cam.zoomTo(600)
	assert.Equal(t, cam.Z(), MaxZoom, "Zooming too close, should get us to the minimum distance")

	cam.zoomTo(-10)
	assert.Equal(t, cam.Z(), MinZoom, "Zooming too far, should get us to the maximum distance")
}

func TestCameraAddOnlyOne(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)

	engo.Mailbox = &engo.MessageManager{}
	CameraBounds = engo.AABB{Min: engo.Point{X: 0, Y: 0}, Max: engo.Point{X: 300, Y: 300}}
	w := &ecs.World{}

	w.AddSystem(&CameraSystem{})
	w.AddSystem(&CameraSystem{})

	expected := "More than one CameraSystem was added to the World. The RenderSystem adds a CameraSystem if none exist when it's added.\n"
	if !strings.HasSuffix(buf.String(), expected) {
		t.Error("adding more than one CameraSystem did not write expected output to log")
	}
}
