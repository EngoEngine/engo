package engi

import (
	"testing"

	"github.com/paked/engi/ecs"
	"github.com/stretchr/testify/assert"
)

func initialize() {
	Mailbox = &MessageManager{}
	WorldBounds = AABB{Point{0, 0}, Point{300, 300}}
	currentWorld = &ecs.World{}

	Cam = &cameraSystem{}
	Cam.New(currentWorld)
}

type CameraTestScene struct{}

func (*CameraTestScene) Preload()         {}
func (*CameraTestScene) Setup(*ecs.World) {}
func (*CameraTestScene) Show()            {}
func (*CameraTestScene) Hide()            {}
func (*CameraTestScene) Type() string     { return "CameraTestScene" }

func TestCameraMoveX(t *testing.T) {
	initialize()

	currentX := Cam.X()

	Cam.moveY(0.50)
	assert.Equal(t, Cam.X(), currentX, "Moving in the Y direction shouldn't change the X location")

	Cam.moveY(-1)
	assert.Equal(t, Cam.X(), currentX, "Moving in the Y direction shouldn't change the X location")

	Cam.zoom(0.50)
	assert.Equal(t, Cam.X(), currentX, "Zooming in should not change the X location")

	Cam.zoom(-1)
	assert.Equal(t, Cam.X(), currentX, "Zooming out should not change the X location")

	Cam.moveX(10)
	assert.Equal(t, Cam.X(), currentX+10, "Moving by 10 units, should have moved the camera by 10 units")

	Cam.moveX(-10)
	assert.Equal(t, Cam.X(), currentX, "Moving by -10 units, should have moved the camera back by 10 units")

	Cam.moveX(305)
	assert.Equal(t, Cam.X(), WorldBounds.Max.X, "Moving too many unit, should have moved the camera to the maximum")

	Cam.moveX(-305)
	assert.Equal(t, Cam.X(), WorldBounds.Min.X, "Moving too many units back, should have moved the camera to the minimum")
}

func TestCameraMoveY(t *testing.T) {
	initialize()

	currentY := Cam.Y()

	Cam.moveX(0.50)
	assert.Equal(t, Cam.Y(), currentY, "Moving in the X direction shouldn't change the Y location")

	Cam.moveX(-1)
	assert.Equal(t, Cam.Y(), currentY, "Moving in the X direction shouldn't change the Y location")

	Cam.zoom(0.50)
	assert.Equal(t, Cam.Y(), currentY, "Zooming in should not change the Y location")

	Cam.zoom(-1)
	assert.Equal(t, Cam.Y(), currentY, "Zooming out should not change the Y location")

	Cam.moveY(10)
	assert.Equal(t, Cam.Y(), currentY+10, "Moving by 10 units, should have moved the camera by 10 units")

	Cam.moveY(-10)
	assert.Equal(t, Cam.Y(), currentY, "Moving by -10 units, should have moved the camera back by 10 units")

	Cam.moveY(305)
	assert.Equal(t, Cam.Y(), WorldBounds.Max.Y, "Moving too many unit, should have moved the camera to the maximum")

	Cam.moveY(-305)
	assert.Equal(t, Cam.Y(), WorldBounds.Min.Y, "Moving too many units back, should have moved the camera to the minimum")
}

func TestCameraZoom(t *testing.T) {
	initialize()

	currentZ := Cam.Z()

	Cam.moveX(0.5)
	assert.Equal(t, Cam.Z(), currentZ, "Moving in the X direction shouldn't change the zoom level")

	Cam.moveX(-1)
	assert.Equal(t, Cam.Z(), currentZ, "Moving in the X direction shouldn't change the zoom level")

	Cam.moveY(0.5)
	assert.Equal(t, Cam.Z(), currentZ, "Moving in the Y direction shouldn't change the zoom level")

	Cam.moveY(-1)
	assert.Equal(t, Cam.Z(), currentZ, "Moving in the Y direction shouldn't change the zoom level")

	Cam.zoom(0.5)
	assert.Equal(t, Cam.Z(), currentZ+0.5, "Should be zoomed in an additional 0.5")

	Cam.zoom(-0.5)
	assert.Equal(t, Cam.Z(), currentZ, "Should be zoomed out to the starting zoom level")

	Cam.zoom(-1000)
	assert.Equal(t, Cam.Z(), MinZoom, "Should be zoomed out to the minimum zoom level")

	Cam.zoom(1000)
	assert.Equal(t, Cam.Z(), MaxZoom, "Should be zoomed in to the maximum zoom level")
}

func TestCameraMoveToX(t *testing.T) {
	initialize()

	currentX := Cam.X()

	Cam.moveToX(currentX + 5)
	assert.Equal(t, Cam.X(), currentX+5, "Moving to current + 5 should get us to current + 5")

	Cam.moveToX(600)
	assert.Equal(t, Cam.X(), WorldBounds.Max.X, "Moving to a location out of bounds, should get us to the maximum")

	Cam.moveToX(-10)
	assert.Equal(t, Cam.X(), WorldBounds.Min.X, "Moving to a location out of bounds, should get us to the minimum")
}

func TestCameraMoveToY(t *testing.T) {
	initialize()

	currentY := Cam.Y()

	Cam.moveToY(currentY + 5)
	assert.Equal(t, Cam.Y(), currentY+5, "Moving to current + 5 should get us to current + 5")

	Cam.moveToY(600)
	assert.Equal(t, Cam.Y(), WorldBounds.Max.Y, "Moving to a location out of bounds, should get us to the maximum")

	Cam.moveToY(-10)
	assert.Equal(t, Cam.Y(), WorldBounds.Min.Y, "Moving to a location out of bounds, should get us to the minimum")
}

func TestCameraZoomTo(t *testing.T) {
	initialize()

	currentZ := Cam.Z()

	Cam.zoomTo(currentZ + 1)
	assert.Equal(t, Cam.Z(), currentZ+1, "Zooming to current + 1 should get us to current + 1")

	Cam.zoomTo(600)
	assert.Equal(t, Cam.Z(), MaxZoom, "Zooming too close, should get us to the minimum distance")

	Cam.zoomTo(-10)
	assert.Equal(t, Cam.Z(), MinZoom, "Zooming too far, should get us to the maximum distance")
}
