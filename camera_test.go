package engi

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCameraMoveX(t *testing.T) {
	cam := &cameraSystem{}
	cam.New()
	WorldBounds = AABB{Point{0, 0}, Point{300, 300}}

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
	assert.Equal(t, cam.X(), WorldBounds.Max.X, "Moving too many unit, should have moved the camera to the maximum")

	cam.moveX(-305)
	assert.Equal(t, cam.X(), WorldBounds.Min.X, "Moving too many units back, should have moved the camera to the minimum")
}

func TestCameraMoveY(t *testing.T) {
	cam := &cameraSystem{}
	cam.New()
	WorldBounds = AABB{Point{0, 0}, Point{300, 300}}

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
	assert.Equal(t, cam.Y(), WorldBounds.Max.Y, "Moving too many unit, should have moved the camera to the maximum")

	cam.moveY(-305)
	assert.Equal(t, cam.Y(), WorldBounds.Min.Y, "Moving too many units back, should have moved the camera to the minimum")
}

func TestCameraZoom(t *testing.T) {
	cam := &cameraSystem{}
	cam.New()

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
	cam := &cameraSystem{}
	cam.New()
	WorldBounds = AABB{Point{0, 0}, Point{300, 300}}

	currentX := cam.X()

	cam.moveToX(currentX + 5)
	assert.Equal(t, cam.X(), currentX+5, "Moving to current + 5 should get us to current + 5")

	cam.moveToX(600)
	assert.Equal(t, cam.X(), WorldBounds.Max.X, "Moving to a location out of bounds, should get us to the maximum")

	cam.moveToX(-10)
	assert.Equal(t, cam.X(), WorldBounds.Min.X, "Moving to a location out of bounds, should get us to the minimum")
}

func TestCameraMoveToY(t *testing.T) {
	cam := &cameraSystem{}
	cam.New()
	WorldBounds = AABB{Point{0, 0}, Point{300, 300}}

	currentY := cam.Y()

	cam.moveToY(currentY + 5)
	assert.Equal(t, cam.Y(), currentY+5, "Moving to current + 5 should get us to current + 5")

	cam.moveToY(600)
	assert.Equal(t, cam.Y(), WorldBounds.Max.Y, "Moving to a location out of bounds, should get us to the maximum")

	cam.moveToY(-10)
	assert.Equal(t, cam.Y(), WorldBounds.Min.Y, "Moving to a location out of bounds, should get us to the minimum")
}

func TestCameraZoomTo(t *testing.T) {
	cam := &cameraSystem{}
	cam.New()

	currentZ := cam.Z()

	cam.zoomTo(currentZ + 1)
	assert.Equal(t, cam.Z(), currentZ+1, "Zooming to current + 1 should get us to current + 1")

	cam.zoomTo(600)
	assert.Equal(t, cam.Z(), MaxZoom, "Zooming too close, should get us to the minimum distance")

	cam.zoomTo(-10)
	assert.Equal(t, cam.Z(), MinZoom, "Zooming too far, should get us to the maximum distance")
}
