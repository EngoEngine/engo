package common

import (
	"engo.io/ecs"
	"engo.io/engo"
)

type HitGroup byte

type HitBox struct {
	x, y, w, h float32
}

func (a HitBox) Hit(b HitBox) bool {
	if a.x > b.x+b.w {
		return false
	}
	if a.y > b.y+b.h {
		return false
	}
	if b.x > a.x+a.w {
		return false
	}
	if b.y > a.y+a.h {
		return false
	}
	return true
}

//Minimum step a needs to take to get off of B
func (a HitBox) MinimumStepOffD(b HitBox) (float32, float32) {
	angle := 0 // top
	dist := a.y + a.h - b.y

	// right
	if b.x+b.w-a.x < dist {
		dist = b.x + b.w - a.x
		angle = 1
	}
	// bottom
	if b.y+b.h-a.y < dist {
		dist = b.y + b.h - a.y
		angle = 2
	}
	//left
	if a.x+a.w-b.x < dist {
		return b.x - (a.x + a.w), 0
	}
	switch angle {
	case 0:
		return 0, -dist
	case 1:
		return dist, 0
	default:
		return 0, dist
	}

}

type Hitable interface {
	ID() uint64
	GetHitBox() HitBox
	HitGroups() (HitGroup, HitGroup)
	Push(float32, float32)
}

type HitMessage struct {
	Mainob  Hitable
	Groupob Hitable
}

func (HitMessage) Type() string {
	return "HitMessage"
}

type HitSystem struct {
	Solid    HitGroup
	Entities []Hitable
}

func (hs *HitSystem) Add(h Hitable) {
	hs.Entities = append(hs.Entities, h)
}

func (hs *HitSystem) Remove(be ecs.BasicEntity) {
	id := be.ID()
	del := -1
	for k, v := range hs.Entities {
		if v.ID() == id {
			del = k
			break
		}
	}
	hs.Entities = append(hs.Entities[:del], hs.Entities[del+1:]...)
}

//Dispatch is to enable testing without the full engo system running
func (hs *HitSystem) Dispatch(m engo.Message) {
	if engo.Mailbox != nil {
		engo.Mailbox.Dispatch(m)
	}
}

func (hs *HitSystem) Update(dt float32) {
	for _, v := range hs.Entities {
		mg, _ := v.HitGroups()
		vhb := v.GetHitBox()
		if mg == 0 {
			continue
		}
		for _, gob := range hs.Entities {
			if gob == v {
				continue
			}
			_, gg := gob.HitGroups()
			if mg&gg == 0 {
				continue
			}
			ghb := gob.GetHitBox()
			if !vhb.Hit(ghb) {
				continue
			}
			//is a Hit
			hs.Dispatch(HitMessage{Mainob: v, Groupob: gob})
			if mg&gg&hs.Solid != 0 {
				dx, dy := vhb.MinimumStepOffD(ghb)
				v.Push(dx, dy)
			}
		}

	}

}
