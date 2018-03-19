package common

import (
	"engo.io/ecs"
	"engo.io/engo"
)

type HitGroup byte

//HitBox is currently a simple bounds rect.
//I hope to add a possibility of internal rects, for more fine grain hit tests at some point
type HitBox struct {
	x, y, w, h float32
}

//Hit tests whether this Hitbox, has an area collision with another HitBox
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

//MinimumStepOffD calculates the mininum step this HitBox needs to take to get off of 'b'
//Return dx, dy
func (a HitBox) MinimumStepOffD(b HitBox) (float32, float32) {
	//angle and dist are the current winners
	//0 = up, 1 = right ...

	// top
	angle := 0
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

//Hitable is the main interface for the HitSystem
//Entities can implement the methods by containing Components that implement them
type Hitable interface {
	ID() uint64
	GetHitBox() HitBox
	HitGroups() (HitGroup, HitGroup)
	Shunt(float32, float32)
}

//HitMessage the message for the engo Messagebox
type HitMessage struct {
	Mainob  Hitable
	Groupob Hitable
}

//Type fulfils the engo Messagebox interface
func (HitMessage) Type() string {
	return "HitMessage"
}

//HitSystem is the replacement system for the CollisionSystem
//It runs on interfaces instead of directly working on Components
type HitSystem struct {
	Solid    HitGroup
	Entities []Hitable
}

//Add adds a hitable entity to the system
func (hs *HitSystem) Add(h Hitable) {
	hs.Entities = append(hs.Entities, h)
}

//Remove removes a hitable entity from the system
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

//Update, checks all 'main' entities against all other entities, to see if any collide
//if they do dispatches a message
//if the collision is 'Solid' moves the main item off of the other.
func (hs *HitSystem) Update(dt float32) {
	for _, v := range hs.Entities {
		vmain, vgrp := v.HitGroups()
		vhb := v.GetHitBox()
		if vmain == 0 {
			continue
		}
		for _, e2 := range hs.Entities {
			if e2 == v {
				continue
			}
			e2main, e2grp := e2.HitGroups()
			if vmain&e2grp == 0 {
				continue
			}
			ghb := e2.GetHitBox()
			if !vhb.Hit(ghb) {
				continue
			}
			//is a Hit
			if engo.Mailbox != nil { // Workaround for safe testing
				engo.Mailbox.Dispatch(HitMessage{Mainob: v, Groupob: e2})
			}

			if vmain&e2grp&hs.Solid != 0 {
				dx, dy := vhb.MinimumStepOffD(ghb)
				if e2main&vgrp&hs.Solid != 0 {
					//collision is between equals
					v.Shunt(dx/2, dy/2)
					e2.Shunt(-dx/2, -dy/2)
					engo.Mailbox.Dispatch(HitMessage{Mainob: e2, Groupob: v})
				} else {
					v.Shunt(dx, dy)
				}
			}
		}

	}

}
