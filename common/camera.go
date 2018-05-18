package common

import (
	"log"
	"sync"
	"time"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/math"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	// MouseRotatorPriority is the priority for the MouseRotatorSystem.
	// Priorities determine the order in which the system is updated.
	MouseRotatorPriority = 100
	// MouseZoomerPriority is the priority for he MouseZoomerSystem.
	// Priorities determine the order in which the system is updated.
	MouseZoomerPriority = 110
	// EdgeScrollerPriority is the priority for the EdgeScrollerSystem.
	// Priorities determine the order in which the system is updated.
	EdgeScrollerPriority = 120
	// KeyboardScrollerPriority is the priority for the KeyboardScrollerSystem.
	// Priorities determine the order in which the system is updated.
	KeyboardScrollerPriority = 130
	// EntityScrollerPriority is the priority for the EntityScrollerSystem.
	// Priorities determine the order in which the system is updated.
	EntityScrollerPriority = 140
)

var (
	// MinZoom is the closest the camera position can be relative to the
	// rendered surface. Smaller numbers of MinZoom allows greater
	// perceived zooming "in".
	MinZoom float32 = 0.25
	// MaxZoom is the farthest the camera position can be relative to the
	// rendered surface. Larger numbers of MaxZoom allows greater
	// perceived zooming "out".
	MaxZoom float32 = 3

	// CameraBounds is the bounding box of the camera
	CameraBounds engo.AABB
)

type cameraEntity struct {
	*ecs.BasicEntity
	*SpaceComponent
}

// CameraSystem is a System that manages the state of the virtual camera. Only
// one CameraSystem can be in a World at a time. If more than one CameraSystem
// is added to the World, it will panic.
type CameraSystem struct {
	x, y, z       float32
	tracking      cameraEntity // The entity that is currently being followed
	trackRotation bool         // Rotate with the entity

	// angle is the angle of the camera, in degrees (not radians!)
	angle float32

	longTasks map[CameraAxis]*CameraMessage
}

// New initializes the CameraSystem.
func (cam *CameraSystem) New(w *ecs.World) {
	num := 0
	for _, sys := range w.Systems() {
		switch sys.(type) {
		case *CameraSystem:
			num++
		}
	}
	if num > 0 { //initalizer is called before added to w.systems
		warning("More than one CameraSystem was added to the World. The RenderSystem adds a CameraSystem if none exist when it's added.")
	}

	if CameraBounds.Max.X == 0 && CameraBounds.Max.Y == 0 {
		CameraBounds.Max = engo.Point{X: engo.GameWidth(), Y: engo.GameHeight()}
	}

	cam.x = CameraBounds.Max.X / 2
	cam.y = CameraBounds.Max.Y / 2
	cam.z = 1

	cam.longTasks = make(map[CameraAxis]*CameraMessage)

	engo.Mailbox.Listen("CameraMessage", func(msg engo.Message) {
		cammsg, ok := msg.(CameraMessage)
		if !ok {
			return
		}

		// Stop with whatever we're doing now
		if _, ok := cam.longTasks[cammsg.Axis]; ok {
			delete(cam.longTasks, cammsg.Axis)
		}

		if cammsg.Duration > time.Duration(0) {
			cam.longTasks[cammsg.Axis] = &cammsg
			return // because it's handled incrementally
		}

		if cammsg.Incremental {
			switch cammsg.Axis {
			case XAxis:
				cam.moveX(cammsg.Value)
			case YAxis:
				cam.moveY(cammsg.Value)
			case ZAxis:
				cam.zoom(cammsg.Value)
			case Angle:
				cam.rotate(cammsg.Value)
			}
		} else {
			switch cammsg.Axis {
			case XAxis:
				cam.moveToX(cammsg.Value)
			case YAxis:
				cam.moveToY(cammsg.Value)
			case ZAxis:
				cam.zoomTo(cammsg.Value)
			case Angle:
				cam.rotateTo(cammsg.Value)
			}
		}
	})
}

// Remove does nothing since the CameraSystem has only one entity, the camera itself.
// This is here to implement the ecs.System interface.
func (cam *CameraSystem) Remove(ecs.BasicEntity) {}

// Update updates the camera. lLong tasks are attempted to update incrementally in batches.
func (cam *CameraSystem) Update(dt float32) {
	for axis, longTask := range cam.longTasks {
		if !longTask.Incremental {
			longTask.Incremental = true

			switch axis {
			case XAxis:
				longTask.Value -= cam.x
			case YAxis:
				longTask.Value -= cam.y
			case ZAxis:
				longTask.Value -= cam.z
			case Angle:
				longTask.Value -= cam.angle
			}
		}

		// Set speed if needed
		if longTask.speed == 0 {
			longTask.speed = longTask.Value / float32(longTask.Duration.Seconds())
		}

		dAxis := longTask.speed * dt
		switch axis {
		case XAxis:
			cam.moveX(dAxis)
		case YAxis:
			cam.moveY(dAxis)
		case ZAxis:
			cam.zoom(dAxis)
		case Angle:
			cam.rotate(dAxis)
		}

		longTask.Duration -= time.Duration(dt)
		if longTask.Duration <= time.Duration(0) {
			delete(cam.longTasks, axis)
		}
	}

	if cam.tracking.BasicEntity == nil {
		return
	}

	if cam.tracking.SpaceComponent == nil {
		log.Println("Should be tracking", cam.tracking.BasicEntity.ID(), "but SpaceComponent is nil")
		cam.tracking.BasicEntity = nil
		return
	}

	cam.centerCam(cam.tracking.SpaceComponent.Position.X+cam.tracking.SpaceComponent.Width/2,
		cam.tracking.SpaceComponent.Position.Y+cam.tracking.SpaceComponent.Height/2,
		cam.z,
	)
	if cam.trackRotation {
		cam.rotateTo(cam.tracking.SpaceComponent.Rotation)
	}
}

// FollowEntity sets the camera to follow the entity with BasicEntity basic
// and SpaceComponent space.
func (cam *CameraSystem) FollowEntity(basic *ecs.BasicEntity, space *SpaceComponent, trackRotation bool) {
	cam.tracking = cameraEntity{basic, space}
	cam.trackRotation = trackRotation
}

// X returns the X-coordinate of the location of the Camera.
func (cam *CameraSystem) X() float32 {
	return cam.x
}

// Y returns the Y-coordinate of the location of the Camera.
func (cam *CameraSystem) Y() float32 {
	return cam.y
}

// Z returns the Z-coordinate of the location of the Camera.
func (cam *CameraSystem) Z() float32 {
	return cam.z
}

// Angle returns the angle (in degrees) at which the Camera is rotated.
func (cam *CameraSystem) Angle() float32 {
	return cam.angle
}

func (cam *CameraSystem) moveX(value float32) {
	if cam.x+(value*engo.GetGlobalScale().X) > CameraBounds.Max.X*engo.GetGlobalScale().X {
		cam.x = CameraBounds.Max.X * engo.GetGlobalScale().X
	} else if cam.x+(value*engo.GetGlobalScale().X) < CameraBounds.Min.X*engo.GetGlobalScale().X {
		cam.x = CameraBounds.Min.X * engo.GetGlobalScale().X
	} else {
		cam.x += value * engo.GetGlobalScale().X
	}
}

func (cam *CameraSystem) moveY(value float32) {
	if cam.y+(value*engo.GetGlobalScale().Y) > CameraBounds.Max.Y*engo.GetGlobalScale().Y {
		cam.y = CameraBounds.Max.Y * engo.GetGlobalScale().Y
	} else if cam.y+(value*engo.GetGlobalScale().Y) < CameraBounds.Min.Y*engo.GetGlobalScale().Y {
		cam.y = CameraBounds.Min.Y * engo.GetGlobalScale().Y
	} else {
		cam.y += value * engo.GetGlobalScale().Y
	}
}

func (cam *CameraSystem) zoom(value float32) {
	cam.zoomTo(cam.z + value)
}

func (cam *CameraSystem) rotate(value float32) {
	cam.rotateTo(cam.angle + value)
}

func (cam *CameraSystem) moveToX(location float32) {
	cam.x = mgl32.Clamp(location*engo.GetGlobalScale().X, CameraBounds.Min.X*engo.GetGlobalScale().X, CameraBounds.Max.X*engo.GetGlobalScale().X)
}

func (cam *CameraSystem) moveToY(location float32) {
	cam.y = mgl32.Clamp(location*engo.GetGlobalScale().Y, CameraBounds.Min.Y*engo.GetGlobalScale().Y, CameraBounds.Max.Y*engo.GetGlobalScale().Y)
}

func (cam *CameraSystem) zoomTo(zoomLevel float32) {
	cam.z = mgl32.Clamp(zoomLevel, MinZoom, MaxZoom)
}

func (cam *CameraSystem) rotateTo(rotation float32) {
	cam.angle = math.Mod(rotation, 360)
}

func (cam *CameraSystem) centerCam(x, y, z float32) {
	cam.moveToX(x)
	cam.moveToY(y)
	cam.zoomTo(z)
}

// CameraAxis is the axis at which the Camera can/has to move.
type CameraAxis uint8

const (
	// XAxis is the x-axis of the camera
	XAxis CameraAxis = iota
	// YAxis is the y-axis of the camera.
	YAxis
	// ZAxis is the z-axis of the camera.
	ZAxis
	// Angle is the angle the camera is rotated by.
	Angle
)

// CameraMessage is a message that can be sent to the Camera (and other Systemers),
// to indicate movement.
type CameraMessage struct {
	Axis        CameraAxis
	Value       float32
	Incremental bool
	Duration    time.Duration
	speed       float32
}

// Type implements the engo.Message interface.
func (CameraMessage) Type() string {
	return "CameraMessage"
}

// KeyboardScroller is a System that allows for scrolling when certain keys are pressed.
type KeyboardScroller struct {
	ScrollSpeed                  float32
	horizontalAxis, verticalAxis string
	keysMu                       sync.RWMutex
}

// Priority implememts the ecs.Prioritizer interface.
func (*KeyboardScroller) Priority() int { return KeyboardScrollerPriority }

// Remove does nothing because the KeyboardScroller system has no entities. It implements the
// ecs.System interface.
func (*KeyboardScroller) Remove(ecs.BasicEntity) {}

// Update updates the camera based on keyboard input.
func (c *KeyboardScroller) Update(dt float32) {
	c.keysMu.RLock()
	defer c.keysMu.RUnlock()

	m := engo.Point{
		X: engo.Input.Axis(c.horizontalAxis).Value(),
		Y: engo.Input.Axis(c.verticalAxis).Value(),
	}
	n, _ := m.Normalize()
	engo.Mailbox.Dispatch(CameraMessage{Axis: XAxis, Value: n.X * c.ScrollSpeed * dt, Incremental: true})
	engo.Mailbox.Dispatch(CameraMessage{Axis: YAxis, Value: n.Y * c.ScrollSpeed * dt, Incremental: true})
}

// BindKeyboard sets the vertical and horizontal axes used by the KeyboardScroller.
func (c *KeyboardScroller) BindKeyboard(hori, vert string) {
	c.keysMu.Lock()

	c.verticalAxis = vert
	c.horizontalAxis = hori

	defer c.keysMu.Unlock()
}

// NewKeyboardScroller creates a new KeyboardScroller system using the provided scrollSpeed,
// and horizontal and vertical axes.
func NewKeyboardScroller(scrollSpeed float32, hori, vert string) *KeyboardScroller {
	kbs := &KeyboardScroller{
		ScrollSpeed: scrollSpeed,
	}

	kbs.BindKeyboard(hori, vert)

	return kbs
}

// EntityScroller scrolls the camera to the position of a entity using its space component.
type EntityScroller struct {
	*SpaceComponent
	TrackingBounds engo.AABB
	Rotation       bool
}

// New adjusts CameraBounds to the bounds of EntityScroller.
func (c *EntityScroller) New(*ecs.World) {
	offsetX, offsetY := engo.GameWidth()/2, engo.GameWidth()/2

	CameraBounds.Min.X = c.TrackingBounds.Min.X + (offsetX / engo.GetGlobalScale().X)
	CameraBounds.Min.Y = c.TrackingBounds.Min.Y + (offsetY / engo.GetGlobalScale().Y)

	CameraBounds.Max.X = c.TrackingBounds.Max.X - (offsetX / engo.GetGlobalScale().X)
	CameraBounds.Max.Y = c.TrackingBounds.Max.Y - (offsetY / engo.GetGlobalScale().Y)
}

// Priority implements the ecs.Prioritizer interface.
func (*EntityScroller) Priority() int { return EntityScrollerPriority }

// Remove does nothing because the EntityScroller system has no entities. This implements
// the ecs.System interface.
func (*EntityScroller) Remove(ecs.BasicEntity) {}

// Update moves the camera to the center of the space component.
// Values are automatically clamped to TrackingBounds by the camera.
func (c *EntityScroller) Update(dt float32) {
	width, height := c.SpaceComponent.Width, c.SpaceComponent.Height

	pos := c.SpaceComponent.Position
	trackToX := pos.X + width/2
	trackToY := pos.Y + height/2

	engo.Mailbox.Dispatch(CameraMessage{Axis: XAxis, Value: trackToX, Incremental: false})
	engo.Mailbox.Dispatch(CameraMessage{Axis: YAxis, Value: trackToY, Incremental: false})
	if c.Rotation {
		engo.Mailbox.Dispatch(CameraMessage{Axis: Angle, Value: c.SpaceComponent.Rotation, Incremental: false})
	}
}

// EdgeScroller is a System that allows for scrolling when the cursor is near the edges of
// the window.
type EdgeScroller struct {
	ScrollSpeed float32
	EdgeMargin  float32
}

// Priority implements the ecs.Prioritizer interface.
func (*EdgeScroller) Priority() int { return EdgeScrollerPriority }

// Remove does nothing because EdgeScroller has no entities. It implements the ecs.System
// interface.
func (*EdgeScroller) Remove(ecs.BasicEntity) {}

// Update moves the camera based on the position of the mouse. If the mouse is on the edge
// of the screen, the camera moves towards that edge.
// TODO: Warning doesn't get the cursor position
func (c *EdgeScroller) Update(dt float32) {
	curX, curY := engo.CursorPos()
	maxX, maxY := engo.GameWidth(), engo.GameHeight()

	if curX < c.EdgeMargin && curY < c.EdgeMargin {
		s := math.Sqrt(2)
		engo.Mailbox.Dispatch(CameraMessage{Axis: XAxis, Value: -c.ScrollSpeed * dt / s, Incremental: true})
		engo.Mailbox.Dispatch(CameraMessage{Axis: YAxis, Value: -c.ScrollSpeed * dt / s, Incremental: true})
	} else if curX < c.EdgeMargin && curY > maxY-c.EdgeMargin {
		s := math.Sqrt(2)
		engo.Mailbox.Dispatch(CameraMessage{Axis: XAxis, Value: -c.ScrollSpeed * dt / s, Incremental: true})
		engo.Mailbox.Dispatch(CameraMessage{Axis: YAxis, Value: c.ScrollSpeed * dt / s, Incremental: true})
	} else if curX > maxX-c.EdgeMargin && curY < c.EdgeMargin {
		s := math.Sqrt(2)
		engo.Mailbox.Dispatch(CameraMessage{Axis: XAxis, Value: c.ScrollSpeed * dt / s, Incremental: true})
		engo.Mailbox.Dispatch(CameraMessage{Axis: YAxis, Value: -c.ScrollSpeed * dt / s, Incremental: true})
	} else if curX > maxX-c.EdgeMargin && curY > maxY-c.EdgeMargin {
		s := math.Sqrt(2)
		engo.Mailbox.Dispatch(CameraMessage{Axis: XAxis, Value: c.ScrollSpeed * dt / s, Incremental: true})
		engo.Mailbox.Dispatch(CameraMessage{Axis: YAxis, Value: c.ScrollSpeed * dt / s, Incremental: true})
	} else if curX < c.EdgeMargin {
		engo.Mailbox.Dispatch(CameraMessage{Axis: XAxis, Value: -c.ScrollSpeed * dt, Incremental: true})
	} else if curX > maxX-c.EdgeMargin {
		engo.Mailbox.Dispatch(CameraMessage{Axis: XAxis, Value: c.ScrollSpeed * dt, Incremental: true})
	} else if curY < c.EdgeMargin {
		engo.Mailbox.Dispatch(CameraMessage{Axis: YAxis, Value: -c.ScrollSpeed * dt, Incremental: true})
	} else if curY > maxY-c.EdgeMargin {
		engo.Mailbox.Dispatch(CameraMessage{Axis: YAxis, Value: c.ScrollSpeed * dt, Incremental: true})
	}
}

// MouseZoomer is a System that allows for zooming when the scroll wheel is used.
type MouseZoomer struct {
	ZoomSpeed float32
}

// Priority implements the ecs.Prioritizer interface.
func (*MouseZoomer) Priority() int { return MouseZoomerPriority }

// Remove does nothing because MouseZoomer has no entities. This implements the
// ecs.System interface.
func (*MouseZoomer) Remove(ecs.BasicEntity) {}

// Update zooms the camera in and out based on the movement of the scroll wheel.
func (c *MouseZoomer) Update(float32) {
	if engo.Input.Mouse.ScrollY != 0 {
		engo.Mailbox.Dispatch(CameraMessage{Axis: ZAxis, Value: engo.Input.Mouse.ScrollY * c.ZoomSpeed, Incremental: true})
	}
}

// MouseRotator is a System that allows for rotating the camera based on pressing
// down the scroll wheel.
type MouseRotator struct {
	// RotationSpeed indicates the speed at which the rotation should happen. This is being used together with the
	// movement by the mouse on the X-axis, to compute the actual rotation.
	RotationSpeed float32

	oldX    float32
	pressed bool
}

// Priority implements the ecs.Prioritizer interface.
func (*MouseRotator) Priority() int { return MouseRotatorPriority }

// Remove does nothing because MouseRotator has no entities. This implements the ecs.System
// interface.
func (*MouseRotator) Remove(ecs.BasicEntity) {}

// Update rotates the camera if the scroll wheel is pressed down.
func (c *MouseRotator) Update(float32) {
	if engo.Input.Mouse.Button == engo.MouseButtonMiddle && engo.Input.Mouse.Action == engo.Press {
		c.pressed = true
	}

	if engo.Input.Mouse.Action == engo.Release {
		c.pressed = false
	}

	if c.pressed {
		engo.Mailbox.Dispatch(CameraMessage{Axis: Angle, Value: (c.oldX - engo.Input.Mouse.X) * -c.RotationSpeed, Incremental: true})
	}

	c.oldX = engo.Input.Mouse.X
}
