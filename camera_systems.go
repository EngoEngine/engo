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

func (c *KeyboardScroller) Name() string {
	return "KeyboardScroller"
}

func (c *KeyboardScroller) New() {
	if !c.isSetup {
		c.System = &System{}
		c.isSetup = true
	}
}

func (c *KeyboardScroller) Update(entity *Entity, dt float32) {
	c.keysMu.RLock()
	defer c.keysMu.RUnlock()

	for _, upKey := range c.upKeys {
		if Keys.Get(upKey).Down() {
			Cam.MoveY(-c.scrollSpeed * dt)
			break
		}
	}

	for _, rightKey := range c.rightKeys {
		if Keys.Get(rightKey).Down() {
			Cam.MoveX(c.scrollSpeed * dt)
			break
		}
	}

	for _, downKey := range c.downKeys {
		if Keys.Get(downKey).Down() {
			Cam.MoveY(c.scrollSpeed * dt)
			break
		}
	}

	for _, leftKey := range c.leftKeys {
		if Keys.Get(leftKey).Down() {
			Cam.MoveX(-c.scrollSpeed * dt)
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
	kbs.AddEntity(NewEntity([]string{kbs.Name()}))
	return kbs
}

// EdgeScroller is a Systemer that allows for scrolling when the mouse is near the edges
type EdgeScroller struct {
	*System
	scrollSpeed float32
	margin      float64

	isSetup bool
}

func (c *EdgeScroller) Name() string {
	return "EdgeScroller"
}

func (c *EdgeScroller) New() {
	if !c.isSetup {
		c.System = &System{}
		c.isSetup = true
	}
}

func (c *EdgeScroller) Update(entity *Entity, dt float32) {
	curX, curY := window.GetCursorPos()
	maxX, maxY := window.GetSize()

	if curX < c.margin {
		Cam.MoveX(-c.scrollSpeed * dt)
	} else if curX > float64(maxX)-c.margin {
		Cam.MoveX(c.scrollSpeed * dt)
	}

	if curY < c.margin {
		Cam.MoveY(-c.scrollSpeed * dt)
	} else if curY > float64(maxY)-c.margin {
		Cam.MoveY(c.scrollSpeed * dt)
	}
}

func NewEdgeScroller(scrollSpeed float32, margin float64) *EdgeScroller {
	es := &EdgeScroller{
		scrollSpeed: scrollSpeed,
		margin:      margin,
	}
	es.New()
	es.AddEntity(NewEntity([]string{es.Name()}))
	return es
}
