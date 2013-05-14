// Copyright 2013 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eng

// A Stage wraps Batch and a Camera to provide a simple way of
// managing a viewport. This will eventually provide a 2d scenegraph.
type Stage struct {
	batch                     *Batch
	camera                    *Camera
	width, height             float32
	gutterWidth, gutterHeight float32
}

// NewStage stage with the given width and height, stretching width or
// height if need be if keepAspect is true.
func NewStage(width, height float32, keepAspect bool) *Stage {
	stage := new(Stage)

	if width == 0 {
		width = float32(Width())
	}
	if height == 0 {
		height = float32(Height())
	}

	stage.batch = NewBatch()
	stage.camera = NewCamera(stage.width, stage.height)

	stage.SetViewport(width, height, keepAspect)

	return stage
}

// SetViewport should be called when the window is resized or if you
// want to just resize the stage's view.
func (stage *Stage) SetViewport(width, height float32, keepAspect bool) {
	if keepAspect {
		screenWidth := float32(Width())
		screenHeight := float32(Height())
		if screenHeight/screenWidth < height/width {
			toScreenSpace := screenHeight / height
			toViewportSpace := height / screenHeight
			deviceWidth := width * toScreenSpace
			lengthen := (screenWidth - deviceWidth) * toViewportSpace
			stage.width = width + lengthen
			stage.height = height
			stage.gutterWidth = lengthen / 2
			stage.gutterHeight = 0
		} else {
			toScreenSpace := screenWidth / width
			toViewportSpace := width / screenWidth
			deviceHeight := height * toScreenSpace
			lengthen := (screenHeight - deviceHeight) * toViewportSpace
			stage.height = height + lengthen
			stage.width = width
			stage.gutterWidth = 0
			stage.gutterHeight = lengthen / 2
		}
	} else {
		stage.width = width
		stage.height = height
		stage.gutterWidth = 0
		stage.gutterHeight = 0
	}
	stage.camera.Position.X = stage.width / 2
	stage.camera.Position.Y = stage.height / 2
	stage.camera.ViewportWidth = stage.width
	stage.camera.ViewportHeight = stage.height
}

// Update should be called every time the underlying camera is modified.
func (s *Stage) Update() {
	s.camera.Update()
	s.batch.SetProjection(s.camera.Combined)
}

// Batch returns the stage's batch.
func (s *Stage) Batch() *Batch {
	return s.batch
}

// Camera returns the stage's camera.
func (s *Stage) Camera() *Camera {
	return s.camera
}

// ScreenToStage takes a point on the screen and returns the point in
// the position of that point with respect to the stage's view. The is
// often used to transform mouse clicks to stage coordinates.
func (s *Stage) ScreenToStage(x, y float32) (float32, float32) {
	tmp.X = x
	tmp.Y = y
	tmp.Z = 1
	s.camera.Unproject(tmp)
	return tmp.X, tmp.Y
}

// Width is the stage's width. This value should be used especially
// if the windows aspect ratio is used.
func (s *Stage) Width() float32 {
	return s.width
}

// GutterWidth is the extra width, if any, added to both the left and
// right sides of the stage when using the window's aspect ratio.
func (s *Stage) GutterWidth() float32 {
	return s.gutterWidth
}

// Height is the stage's height. This value should be used especially
// if the windows aspect ratio is used.
func (s *Stage) Height() float32 {
	return s.height
}

// GutterHeight is the extra height, if any, added to both the top and
// bottom sides of the stage when using the window's aspect ratio.
func (s *Stage) GutterHeight() float32 {
	return s.gutterHeight
}
