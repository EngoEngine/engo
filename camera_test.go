package engi_test

import (
	"github.com/paked/engi"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCameraMoveX(t *testing.T) {
	cam := &engi.Camera{}
	cam.Setup()
	engi.WorldBounds = engi.AABB{engi.Point{0, 0}, engi.Point{300, 300}}

	currentX := cam.X()

	cam.MoveY(0.50)
	assert.Equal(t, cam.X(), currentX, "Moving in the Y direction shouldn't change the X location")

	cam.MoveY(-1)
	assert.Equal(t, cam.X(), currentX, "Moving in the Y direction shouldn't change the X location")

	cam.Zoom(0.50)
	assert.Equal(t, cam.X(), currentX, "Zooming in should not change the X location")

	cam.Zoom(-1)
	assert.Equal(t, cam.X(), currentX, "Zooming out should not change the X location")

	cam.MoveX(10)
	assert.Equal(t, cam.X(), currentX+10, "Moving by 10 units, should have moved the camera by 10 units")

	cam.MoveX(-10)
	assert.Equal(t, cam.X(), currentX, "Moving by -10 units, should have moved the camera back by 10 units")

	cam.MoveX(305)
	assert.Equal(t, cam.X(), engi.WorldBounds.Max.X, "Moving too many unit, should have moved the camera to the maximum")

	cam.MoveX(-305)
	assert.Equal(t, cam.X(), engi.WorldBounds.Min.X, "Moving too many units back, should have moved the camera to the minimum")
}

func TestCameraMoveY(t *testing.T) {
	cam := &engi.Camera{}
	cam.Setup()
	engi.WorldBounds = engi.AABB{engi.Point{0, 0}, engi.Point{300, 300}}

	currentY := cam.Y()

	cam.MoveX(0.50)
	assert.Equal(t, cam.Y(), currentY, "Moving in the X direction shouldn't change the Y location")

	cam.MoveX(-1)
	assert.Equal(t, cam.Y(), currentY, "Moving in the X direction shouldn't change the Y location")

	cam.Zoom(0.50)
	assert.Equal(t, cam.Y(), currentY, "Zooming in should not change the Y location")

	cam.Zoom(-1)
	assert.Equal(t, cam.Y(), currentY, "Zooming out should not change the Y location")

	cam.MoveY(10)
	assert.Equal(t, cam.Y(), currentY+10, "Moving by 10 units, should have moved the camera by 10 units")

	cam.MoveY(-10)
	assert.Equal(t, cam.Y(), currentY, "Moving by -10 units, should have moved the camera back by 10 units")

	cam.MoveY(305)
	assert.Equal(t, cam.Y(), engi.WorldBounds.Max.Y, "Moving too many unit, should have moved the camera to the maximum")

	cam.MoveY(-305)
	assert.Equal(t, cam.Y(), engi.WorldBounds.Min.Y, "Moving too many units back, should have moved the camera to the minimum")
}

func TestCameraZoom(t *testing.T) {
	cam := &engi.Camera{}
	cam.Setup()

	currentZ := cam.Z()

	cam.MoveX(0.5)
	assert.Equal(t, cam.Z(), currentZ, "Moving in the X direction shouldn't change the zoom level")

	cam.MoveX(-1)
	assert.Equal(t, cam.Z(), currentZ, "Moving in the X direction shouldn't change the zoom level")

	cam.MoveY(0.5)
	assert.Equal(t, cam.Z(), currentZ, "Moving in the Y direction shouldn't change the zoom level")

	cam.MoveY(-1)
	assert.Equal(t, cam.Z(), currentZ, "Moving in the Y direction shouldn't change the zoom level")

	cam.Zoom(0.5)
	assert.Equal(t, cam.Z(), currentZ+0.5, "Should be zoomed in an additional 0.5")

	cam.Zoom(-0.5)
	assert.Equal(t, cam.Z(), currentZ, "Should be zoomed out to the starting zoom level")

	cam.Zoom(-1000)
	assert.Equal(t, cam.Z(), engi.MinZoom, "Should be zoomed out to the minimum zoom level")

	cam.Zoom(1000)
	assert.Equal(t, cam.Z(), engi.MaxZoom, "Should be zoomed in to the maximum zoom level")
}

func TestCameraMoveToX(t *testing.T) {
	cam := &engi.Camera{}
	cam.Setup()
	engi.WorldBounds = engi.AABB{engi.Point{0, 0}, engi.Point{300, 300}}

	currentX := cam.X()

	cam.MoveToX(currentX + 5)
	assert.Equal(t, cam.X(), currentX+5, "Moving to current + 5 should get us to current + 5")

	cam.MoveToX(600)
	assert.Equal(t, cam.X(), engi.WorldBounds.Max.X, "Moving to a location out of bounds, should get us to the maximum")

	cam.MoveToX(-10)
	assert.Equal(t, cam.X(), engi.WorldBounds.Min.X, "Moving to a location out of bounds, should get us to the minimum")
}

func TestCameraMoveToY(t *testing.T) {
	cam := &engi.Camera{}
	cam.Setup()
	engi.WorldBounds = engi.AABB{engi.Point{0, 0}, engi.Point{300, 300}}

	currentY := cam.Y()

	cam.MoveToY(currentY + 5)
	assert.Equal(t, cam.Y(), currentY+5, "Moving to current + 5 should get us to current + 5")

	cam.MoveToY(600)
	assert.Equal(t, cam.Y(), engi.WorldBounds.Max.Y, "Moving to a location out of bounds, should get us to the maximum")

	cam.MoveToY(-10)
	assert.Equal(t, cam.Y(), engi.WorldBounds.Min.Y, "Moving to a location out of bounds, should get us to the minimum")
}

func TestCameraZoomTo(t *testing.T) {
	cam := &engi.Camera{}
	cam.Setup()

	currentZ := cam.Z()

	cam.ZoomTo(currentZ + 5)
	assert.Equal(t, cam.Z(), currentZ+5, "Zooming to current + 5 should get us to current + 5")

	cam.ZoomTo(600)
	assert.Equal(t, cam.Z(), engi.MaxZoom, "Zooming too close, should get us to the minimum distance")

	cam.ZoomTo(-10)
	assert.Equal(t, cam.Z(), engi.MinZoom, "Zooming too far, should get us to the maximum distance")
}
