package common

import (
	"testing"

	"engo.io/ecs"
	"engo.io/engo"
)

type EveryComp struct {
	ecs.BasicEntity
	AnimationComponent
	MouseComponent
	RenderComponent
	SpaceComponent
	CollisionComponent
	AudioComponent
}

type TestInterfaceScene struct {
	failed bool
	reason string
}

func (*TestInterfaceScene) Type() string { return "TestInterfaceScene" }

func (*TestInterfaceScene) Preload() {}

func (s *TestInterfaceScene) Setup(u engo.Updater) {
	w, _ := u.(*ecs.World)

	asys := AnimationSystem{}
	var a *Animationable
	var nota *NotAnimationable
	w.AddSystemInterface(&asys, a, nota)

	msys := MouseSystem{}
	var m *Mouseable
	var notm *NotMouseable
	w.AddSystemInterface(&msys, m, notm)

	rsys := RenderSystem{}
	var r *Renderable
	var notr *NotRenderable
	w.AddSystemInterface(&rsys, r, notr)

	csys := CollisionSystem{}
	var c *Collisionable
	var notc *NotCollisionable
	w.AddSystemInterface(&csys, c, notc)

	audsys := AudioSystem{}
	var aud *Audioable
	var notaud *NotAudioable
	w.AddSystemInterface(&audsys, aud, notaud)

	e := &EveryComp{BasicEntity: ecs.NewBasic()}
	w.AddEntity(e)

	if len(asys.entities) != 1 {
		s.failed = true
		s.reason = "did not add entity to animation system"
		return
	}
	asys.Remove(e.BasicEntity)
	if len(asys.entities) != 0 {
		s.failed = true
		s.reason = "did not remove entry from animation system"
		return
	}

	if len(msys.entities) != 1 {
		s.failed = true
		s.reason = "did not add entity to mouse system"
		return
	}
	msys.Remove(e.BasicEntity)
	if len(msys.entities) != 0 {
		s.failed = true
		s.reason = "did not remove entry from mouse system"
		return
	}

	if len(rsys.entities) != 1 {
		s.failed = true
		s.reason = "did not add entity to render system"
		return
	}
	rsys.Remove(e.BasicEntity)
	if len(rsys.entities) != 0 {
		s.failed = true
		s.reason = "did not remove entry from render system"
		return
	}

	if len(csys.entities) != 1 {
		s.failed = true
		s.reason = "did not add entity to collision system"
		return
	}
	csys.Remove(e.BasicEntity)
	if len(csys.entities) != 0 {
		s.failed = true
		s.reason = "did not remove entry from collision system"
		return
	}

	if len(audsys.entities) != 1 {
		s.failed = true
		s.reason = "did not add entity to audio system"
		return
	}
	audsys.Remove(e.BasicEntity)
	if len(audsys.entities) != 0 {
		s.failed = true
		s.reason = "did not remove entry from audio system"
		return
	}
}

// TestEveryInterface Creates an Everything component and tries to add and then remove it from each system to each system using AddByInterface.
// I can't test adding things that don't work as the code won't compile
func TestEveryInterface(t *testing.T) {
	s := &TestInterfaceScene{}
	engo.Run(engo.RunOptions{
		NoRun:        true,
		HeadlessMode: true,
	}, s)
	if s.failed {
		t.Errorf("failed to test every interface. Reason: %v", s.reason)
	}
}
