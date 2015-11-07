package engi

import (
	"sync"
)

// KeyboardScroller is a Systemer that allows for scrolling when certain keys are pressed
type KeyboardScroller struct {
	*System
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

func (c *KeyboardScroller) New() {
	if !c.isSetup {
		c.System = NewSystem()
		c.isSetup = true
	}
}

func (c *KeyboardScroller) Update(entity *Entity, dt float32) {
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
	kbs.New()
	kbs.BindKeyboard(up, right, down, left)
	kbs.AddEntity(NewEntity([]string{kbs.Type()}))
	return kbs
}

// EdgeScroller is a Systemer that allows for scrolling when the mouse is near the edges
type EdgeScroller struct {
	*System
	scrollSpeed float32
	margin      float64

	isSetup bool
}

func (*EdgeScroller) Type() string {
	return "EdgeScroller"
}

func (c *EdgeScroller) New() {
	if !c.isSetup {
		c.System = NewSystem()
		c.isSetup = true
	}
}

func (c *EdgeScroller) Update(entity *Entity, dt float32) {
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
	es.New()
	es.AddEntity(NewEntity([]string{es.Type()}))
	return es
}

// MouseZoomer is a Systemer that allows for zooming when the scroll wheel is used
type MouseZoomer struct {
	*System
	zoomSpeed float32

	isSetup bool
}

func (*MouseZoomer) Type() string {
	return "MouseZoomer"
}

func (c *MouseZoomer) New() {
	if !c.isSetup {
		c.System = NewSystem()
		c.isSetup = true
	}
}

func (c *MouseZoomer) Update(entity *Entity, dt float32) {
	if Mouse.ScrollY != 0 {
		Mailbox.Dispatch(CameraMessage{ZAxis, Mouse.ScrollY * c.zoomSpeed, true})
	}
}

func NewMouseZoomer(zoomSpeed float32) *MouseZoomer {
	es := &MouseZoomer{
		zoomSpeed: zoomSpeed,
	}
	es.New()
	es.AddEntity(NewEntity([]string{es.Type()}))
	return es
}
