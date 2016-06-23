package common

import (
	"log"
	"sync"
	"time"

	"engo.io/ecs"
	"engo.io/engo"
	"github.com/engoengine/math"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	MouseRotatorPriority     = 100
	MouseZoomerPriority      = 110
	EdgeScrollerPriority     = 120
	KeyboardScrollerPriority = 130
	EntityScrollerPriority   = 140
)

var (
	MinZoom float32 = 0.25
	MaxZoom float32 = 3

	CameraBounds engo.AABB
)

type cameraEntity struct {
	*ecs.BasicEntity
	*SpaceComponent
}

// CameraSystem is a System that manages the state of the virtual camera
type CameraSystem struct {
	x, y, z  float32
	tracking cameraEntity // The entity that is currently being followed

	// angle is the angle of the camera, in degrees (not radians!)
	angle float32

	longTasks map[CameraAxis]*CameraMessage
}

func (cam *CameraSystem) New(*ecs.World) {
	if CameraBounds.Max.X == 0 && CameraBounds.Max.Y == 0 {
		CameraBounds.Max = engo.Point{engo.GameWidth(), engo.GameHeight()}
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

func (cam *CameraSystem) Remove(ecs.BasicEntity) {}

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
}

func (cam *CameraSystem) FollowEntity(basic *ecs.BasicEntity, space *SpaceComponent) {
	cam.tracking = cameraEntity{basic, space}
}

// X returns the X-coordinate of the location of the Camera
func (cam *CameraSystem) X() float32 {
	return cam.x
}

// Y returns the Y-coordinate of the location of the Camera
func (cam *CameraSystem) Y() float32 {
	return cam.y
}

// Z returns the Z-coordinate of the location of the Camera
func (cam *CameraSystem) Z() float32 {
	return cam.z
}

// Angle returns the angle (in degrees) at which the Camera is rotated
func (cam *CameraSystem) Angle() float32 {
	return cam.angle
}

func (cam *CameraSystem) moveX(value float32) {
	cam.moveToX(cam.x + value)
}

func (cam *CameraSystem) moveY(value float32) {
	cam.moveToY(cam.y + value)
}

func (cam *CameraSystem) zoom(value float32) {
	cam.zoomTo(cam.z + value)
}

func (cam *CameraSystem) rotate(value float32) {
	cam.rotateTo(cam.angle + value)
}

func (cam *CameraSystem) moveToX(location float32) {
	cam.x = mgl32.Clamp(location, CameraBounds.Min.X, CameraBounds.Max.X)
}

func (cam *CameraSystem) moveToY(location float32) {
	cam.y = mgl32.Clamp(location, CameraBounds.Min.Y, CameraBounds.Max.Y)
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

// CameraAxis is the axis at which the Camera can/has to move
type CameraAxis uint8

const (
	XAxis CameraAxis = iota
	YAxis
	ZAxis
	Angle
)

// CameraMessage is a message that can be sent to the Camera (and other Systemers), to indicate movement
type CameraMessage struct {
	Axis        CameraAxis
	Value       float32
	Incremental bool
	Duration    time.Duration
	speed       float32
}

func (CameraMessage) Type() string {
	return "CameraMessage"
}

// KeyboardScroller is a System that allows for scrolling when certain keys are pressed
type KeyboardScroller struct {
	ScrollSpeed                  float32
	horizontalAxis, verticalAxis string
	keysMu                       sync.RWMutex
}

func (*KeyboardScroller) Priority() int          { return KeyboardScrollerPriority }
func (*KeyboardScroller) Remove(ecs.BasicEntity) {}

func (c *KeyboardScroller) Update(dt float32) {
	c.keysMu.RLock()
	defer c.keysMu.RUnlock()

	vert := engo.Input.Axis(c.verticalAxis)
	engo.Mailbox.Dispatch(CameraMessage{Axis: YAxis, Value: vert.Value() * c.ScrollSpeed * dt, Incremental: true})

	hori := engo.Input.Axis(c.horizontalAxis)
	engo.Mailbox.Dispatch(CameraMessage{Axis: XAxis, Value: hori.Value() * c.ScrollSpeed * dt, Incremental: true})
}

func (c *KeyboardScroller) BindKeyboard(hori, vert string) {
	c.keysMu.Lock()

	c.verticalAxis = vert
	c.horizontalAxis = hori

	defer c.keysMu.Unlock()
}

func NewKeyboardScroller(scrollSpeed float32, hori, vert string) *KeyboardScroller {
	kbs := &KeyboardScroller{
		ScrollSpeed: scrollSpeed,
	}

	kbs.BindKeyboard(hori, vert)

	return kbs
}

// EntityScroller scrolls the camera to the position of a entity using its space component
type EntityScroller struct {
	*SpaceComponent
	TrackingBounds engo.AABB
}

// New adjusts CameraBounds to the bounds of EntityScroller
func (c *EntityScroller) New(*ecs.World) {
	offsetX, offsetY := engo.CanvasWidth()/2, engo.CanvasHeight()/2
	// This is to account for upper left origin
	distY := c.TrackingBounds.Max.Y - c.TrackingBounds.Min.Y

	CameraBounds.Min.X = c.TrackingBounds.Min.X + offsetX
	CameraBounds.Min.Y = (c.TrackingBounds.Max.Y - distY) + offsetY

	CameraBounds.Max.X = c.TrackingBounds.Max.X - offsetX
	CameraBounds.Max.Y = (c.TrackingBounds.Min.Y + distY) - offsetY
}

func (*EntityScroller) Priority() int          { return EntityScrollerPriority }
func (*EntityScroller) Remove(ecs.BasicEntity) {}

// Update moves the camera to the center of the space component
// Values are automatically clamped to TrackingBounds by the camera
func (c *EntityScroller) Update(dt float32) {
	width, height := c.SpaceComponent.Width, c.SpaceComponent.Height

	pos := c.SpaceComponent.Position
	trackToX := pos.X + width/2
	trackToY := pos.Y + height/2

	engo.Mailbox.Dispatch(CameraMessage{Axis: XAxis, Value: trackToX, Incremental: false})
	engo.Mailbox.Dispatch(CameraMessage{Axis: YAxis, Value: trackToY, Incremental: false})
}

// EdgeScroller is a System that allows for scrolling when the cursor is near the edges of
// the window
type EdgeScroller struct {
	ScrollSpeed float32
	EdgeMargin  float32
}

func (*EdgeScroller) Priority() int          { return EdgeScrollerPriority }
func (*EdgeScroller) Remove(ecs.BasicEntity) {}

// TODO: Warning doesn't get the cursor position
func (c *EdgeScroller) Update(dt float32) {
	curX, curY := engo.CursorPos()
	maxX, maxY := engo.CanvasWidth(), engo.CanvasHeight()

	if curX < c.EdgeMargin {
		engo.Mailbox.Dispatch(CameraMessage{Axis: XAxis, Value: -c.ScrollSpeed * dt, Incremental: true})
	} else if curX > maxX-c.EdgeMargin {
		engo.Mailbox.Dispatch(CameraMessage{Axis: XAxis, Value: c.ScrollSpeed * dt, Incremental: true})
	}

	if curY < c.EdgeMargin {
		engo.Mailbox.Dispatch(CameraMessage{Axis: YAxis, Value: -c.ScrollSpeed * dt, Incremental: true})
	} else if curY > maxY-c.EdgeMargin {
		engo.Mailbox.Dispatch(CameraMessage{Axis: YAxis, Value: c.ScrollSpeed * dt, Incremental: true})
	}
}

// MouseZoomer is a System that allows for zooming when the scroll wheel is used
type MouseZoomer struct {
	ZoomSpeed float32
}

func (*MouseZoomer) Priority() int          { return MouseZoomerPriority }
func (*MouseZoomer) Remove(ecs.BasicEntity) {}

func (c *MouseZoomer) Update(float32) {
	if engo.Input.Mouse.ScrollY != 0 {
		engo.Mailbox.Dispatch(CameraMessage{Axis: ZAxis, Value: engo.Input.Mouse.ScrollY * c.ZoomSpeed, Incremental: true})
	}
}

// MouseRotator is a System that allows for zooming when the scroll wheel is used
type MouseRotator struct {
	// RotationSpeed indicates the speed at which the rotation should happen. This is being used together with the
	// movement by the mouse on the X-axis, to compute the actual rotation.
	RotationSpeed float32

	oldX    float32
	pressed bool
}

func (*MouseRotator) Priority() int          { return MouseRotatorPriority }
func (*MouseRotator) Remove(ecs.BasicEntity) {}

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
