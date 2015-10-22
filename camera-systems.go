package engi

import (
	"sync"
)

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
