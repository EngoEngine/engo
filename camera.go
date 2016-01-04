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
	*ecs.System
	x, y, z  float32
	tracking *ecs.Entity // The entity that is currently being followed
}

func (cameraSystem) Type() string {
	return "cameraSystem"
}

func (cam *cameraSystem) New(*ecs.World) {
	cam.System = ecs.NewSystem()

	cam.x = WorldBounds.Max.X / 2
	cam.y = WorldBounds.Max.Y / 2
	cam.z = 1

	cam.AddEntity(ecs.NewEntity([]string{cam.Type()}))

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
	cam.y = mgl32.Clamp(location, WorldBounds.Min.X, WorldBounds.Max.Y)
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

func (cam *cameraSystem) Update(entity *ecs.Entity, dt float32) {
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

// KeyboardScroller is a Systemer that allows for scrolling when certain keys are pressed
type KeyboardScroller struct {
	*ecs.System
	scrollSpeed float32
	upKeys      []Key
	leftKeys    []Key
	downKeys    []Key
	rightKeys   []Key

	keysMu  sync.RWMutex
	isSetup bool
}

func (*KeyboardScroller) Type() string {
	return "KeyboardScroller"
}

func (c *KeyboardScroller) New(*ecs.World) {
	if !c.isSetup {
		c.System = ecs.NewSystem()
		c.isSetup = true
	}
}

func (c *KeyboardScroller) Update(entity *ecs.Entity, dt float32) {
	c.keysMu.RLock()
	defer c.keysMu.RUnlock()

	for _, upKey := range c.upKeys {
		if Keys.Get(upKey).Down() {
			Mailbox.Dispatch(CameraMessage{YAxis, -c.scrollSpeed * dt, true})
			break
		}
	}

	for _, rightKey := range c.rightKeys {
		if Keys.Get(rightKey).Down() {
			Mailbox.Dispatch(CameraMessage{XAxis, c.scrollSpeed * dt, true})
			break
		}
	}

	for _, downKey := range c.downKeys {
		if Keys.Get(downKey).Down() {
			Mailbox.Dispatch(CameraMessage{YAxis, c.scrollSpeed * dt, true})
			break
		}
	}

	for _, leftKey := range c.leftKeys {
		if Keys.Get(leftKey).Down() {
			Mailbox.Dispatch(CameraMessage{XAxis, -c.scrollSpeed * dt, true})
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
		scrollSpeed: scrollSpeed,
	}
	kbs.New(nil)
	kbs.BindKeyboard(up, right, down, left)
	kbs.AddEntity(ecs.NewEntity([]string{kbs.Type()}))
	return kbs
}

// EdgeScroller is a Systemer that allows for scrolling when the mouse is near the edges
type EdgeScroller struct {
	*ecs.System
	scrollSpeed float32
	margin      float64

	isSetup bool
}

func (*EdgeScroller) Type() string {
	return "EdgeScroller"
}

func (c *EdgeScroller) New(*ecs.World) {
	if !c.isSetup {
		c.System = ecs.NewSystem()
		c.isSetup = true
	}
}

func (c *EdgeScroller) Update(entity *ecs.Entity, dt float32) {
	curX, curY := window.GetCursorPos()
	maxX, maxY := window.GetSize()

	if curX < c.margin {
		Mailbox.Dispatch(CameraMessage{XAxis, -c.scrollSpeed * dt, true})
	} else if curX > float64(maxX)-c.margin {
		Mailbox.Dispatch(CameraMessage{XAxis, c.scrollSpeed * dt, true})
	}

	if curY < c.margin {
		Mailbox.Dispatch(CameraMessage{YAxis, -c.scrollSpeed * dt, true})
	} else if curY > float64(maxY)-c.margin {
		Mailbox.Dispatch(CameraMessage{YAxis, c.scrollSpeed * dt, true})
	}
}

func NewEdgeScroller(scrollSpeed float32, margin float64) *EdgeScroller {
	es := &EdgeScroller{
		scrollSpeed: scrollSpeed,
		margin:      margin,
	}
	es.New(nil)
	es.AddEntity(ecs.NewEntity([]string{es.Type()}))
	return es
}

// MouseZoomer is a Systemer that allows for zooming when the scroll wheel is used
type MouseZoomer struct {
	*ecs.System
	zoomSpeed float32

	isSetup bool
}

func (*MouseZoomer) Type() string {
	return "MouseZoomer"
}

func (c *MouseZoomer) New(*ecs.World) {
	if !c.isSetup {
		c.System = ecs.NewSystem()
		c.isSetup = true
	}
}

func (c *MouseZoomer) Update(entity *ecs.Entity, dt float32) {
	if Mouse.ScrollY != 0 {
		Mailbox.Dispatch(CameraMessage{ZAxis, Mouse.ScrollY * c.zoomSpeed, true})
	}
}

func NewMouseZoomer(zoomSpeed float32) *MouseZoomer {
	es := &MouseZoomer{
		zoomSpeed: zoomSpeed,
	}
	es.New(nil)
	es.AddEntity(ecs.NewEntity([]string{es.Type()}))
	return es
}
