package engi

import (
	"sync"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/paked/engi/ecs"
)

var (
	MinZoom float32 = 0.25
	MaxZoom float32 = 3
)

// CameraSystem is a System that manages the state of the Camera
type cameraSystem struct {
	ecs.LinearSystem
	x, y, z  float32
	tracking *ecs.Entity // The entity that is currently being followed
}

func (*cameraSystem) Type() string { return "cameraSystem" }
func (*cameraSystem) Pre()         {}
func (*cameraSystem) Post()        {}

func (cam *cameraSystem) New(*ecs.World) {
	cam.x = WorldBounds.Max.X / 2
	cam.y = WorldBounds.Max.Y / 2
	cam.z = 1

	Mailbox.Listen("CameraMessage", func(msg Message) {
		cammsg, ok := msg.(CameraMessage)
		if !ok {
			return
		}

		if cammsg.Incremental {
			switch cammsg.Axis {
			case XAxis:
				cam.moveX(cammsg.Value)
			case YAxis:
				cam.moveY(cammsg.Value)
			case ZAxis:
				cam.zoom(cammsg.Value)
			}
		} else {
			switch cammsg.Axis {
			case XAxis:
				cam.moveToX(cammsg.Value)
			case YAxis:
				cam.moveToY(cammsg.Value)
			case ZAxis:
				cam.zoomTo(cammsg.Value)
			}
		}
	})
}

func (cam *cameraSystem) FollowEntity(entity *ecs.Entity) {
	cam.tracking = entity
	var space *SpaceComponent

	if _, ok := cam.tracking.ComponentFast(space).(*SpaceComponent); !ok {
		cam.tracking = nil
		return
	}
}

func (cam *cameraSystem) moveX(value float32) {
	cam.moveToX(cam.x + value)
}

func (cam *cameraSystem) moveY(value float32) {
	cam.moveToY(cam.y + value)
}

func (cam *cameraSystem) zoom(value float32) {
	cam.zoomTo(cam.z + value)
}

func (cam *cameraSystem) moveToX(location float32) {
	cam.x = mgl32.Clamp(location, WorldBounds.Min.X, WorldBounds.Max.X)
}

func (cam *cameraSystem) moveToY(location float32) {
	cam.y = mgl32.Clamp(location, WorldBounds.Min.Y, WorldBounds.Max.Y)
}

func (cam *cameraSystem) zoomTo(zoomLevel float32) {
	cam.z = mgl32.Clamp(zoomLevel, MinZoom, MaxZoom)
}

func (cam *cameraSystem) X() float32 {
	return cam.x
}

func (cam *cameraSystem) Y() float32 {
	return cam.y
}

func (cam *cameraSystem) Z() float32 {
	return cam.z
}

func (cam *cameraSystem) UpdateEntity(entity *ecs.Entity, dt float32) {
	if cam.tracking == nil {
		return
	}

	var space *SpaceComponent
	var ok bool

	if space, ok = cam.tracking.ComponentFast(space).(*SpaceComponent); !ok {
		return
	}

	cam.centerCam(space.Position.X+space.Width/2, space.Position.Y+space.Height/2, cam.z)
}

func (cam *cameraSystem) centerCam(x, y, z float32) {
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
)

// CameraMessage is a message that can be sent to the Camera (and other Systemers), to indicate movement
type CameraMessage struct {
	Axis        CameraAxis
	Value       float32
	Incremental bool
}

func (CameraMessage) Type() string {
	return "CameraMessage"
}

// KeyboardScroller is a System that allows for scrolling when certain keys are pressed
type KeyboardScroller struct {
	ScrollSpeed float32
	upKeys      []Key
	leftKeys    []Key
	downKeys    []Key
	rightKeys   []Key

	keysMu sync.RWMutex
}

func (*KeyboardScroller) Type() string             { return "KeyboardScroller" }
func (*KeyboardScroller) Priority() int            { return 10 }
func (*KeyboardScroller) AddEntity(*ecs.Entity)    {}
func (*KeyboardScroller) RemoveEntity(*ecs.Entity) {}
func (*KeyboardScroller) New(*ecs.World)           {}

func (c *KeyboardScroller) Update(dt float32) {
	c.keysMu.RLock()
	defer c.keysMu.RUnlock()

	for _, upKey := range c.upKeys {
		if Keys.Get(upKey).Down() {
			Mailbox.Dispatch(CameraMessage{YAxis, -c.ScrollSpeed * dt, true})
			break
		}
	}

	for _, rightKey := range c.rightKeys {
		if Keys.Get(rightKey).Down() {
			Mailbox.Dispatch(CameraMessage{XAxis, c.ScrollSpeed * dt, true})
			break
		}
	}

	for _, downKey := range c.downKeys {
		if Keys.Get(downKey).Down() {
			Mailbox.Dispatch(CameraMessage{YAxis, c.ScrollSpeed * dt, true})
			break
		}
	}

	for _, leftKey := range c.leftKeys {
		if Keys.Get(leftKey).Down() {
			Mailbox.Dispatch(CameraMessage{XAxis, -c.ScrollSpeed * dt, true})
			break
		}
	}
}

func (c *KeyboardScroller) BindKeyboard(up, right, down, left Key) {
	c.keysMu.Lock()
	defer c.keysMu.Unlock()

	c.upKeys = append(c.upKeys, up)
	c.rightKeys = append(c.rightKeys, right)
	c.downKeys = append(c.downKeys, down)
	c.leftKeys = append(c.leftKeys, left)
}

func NewKeyboardScroller(scrollSpeed float32, up, right, down, left Key) *KeyboardScroller {
	kbs := &KeyboardScroller{
		ScrollSpeed: scrollSpeed,
	}
	kbs.BindKeyboard(up, right, down, left)
	return kbs
}

// EdgeScroller is a System that allows for scrolling when the cursor is near the edges of
// the window
type EdgeScroller struct {
	ScrollSpeed float32
	EdgeMargin  float64
}

func (*EdgeScroller) Type() string             { return "EdgeScroller" }
func (*EdgeScroller) Priority() int            { return 10 }
func (*EdgeScroller) AddEntity(*ecs.Entity)    {}
func (*EdgeScroller) RemoveEntity(*ecs.Entity) {}
func (*EdgeScroller) New(*ecs.World)           {}

func (c *EdgeScroller) Update(dt float32) {
	curX, curY := window.GetCursorPos()
	maxX, maxY := window.GetSize()

	if curX < c.EdgeMargin {
		Mailbox.Dispatch(CameraMessage{XAxis, -c.ScrollSpeed * dt, true})
	} else if curX > float64(maxX)-c.EdgeMargin {
		Mailbox.Dispatch(CameraMessage{XAxis, c.ScrollSpeed * dt, true})
	}

	if curY < c.EdgeMargin {
		Mailbox.Dispatch(CameraMessage{YAxis, -c.ScrollSpeed * dt, true})
	} else if curY > float64(maxY)-c.EdgeMargin {
		Mailbox.Dispatch(CameraMessage{YAxis, c.ScrollSpeed * dt, true})
	}
}

// MouseZoomer is a System that allows for zooming when the scroll wheel is used
type MouseZoomer struct {
	ZoomSpeed float32
}

func (*MouseZoomer) Type() string             { return "MouseZoomer" }
func (*MouseZoomer) Priority() int            { return 10 }
func (*MouseZoomer) AddEntity(*ecs.Entity)    {}
func (*MouseZoomer) RemoveEntity(*ecs.Entity) {}
func (*MouseZoomer) New(*ecs.World)           {}

func (c *MouseZoomer) Update(dt float32) {
	if Mouse.ScrollY != 0 {
		Mailbox.Dispatch(CameraMessage{ZAxis, Mouse.ScrollY * c.ZoomSpeed, true})
	}
}
