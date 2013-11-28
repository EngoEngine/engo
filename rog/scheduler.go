package rog

type Actor interface {
	Act() float32
}

type Scheduler struct {
	lock int
}

func NewScheduler() *Scheduler {
	return new(Scheduler)
}

func (s *Scheduler) Lock() {
	s.lock++
}

func (s *Scheduler) Unlock() {
	if s.lock == 0 {
		panic("Cannot unlock unlocked Scheduler")
	}
	s.lock--
}

func (s *Scheduler) Update() {
	if s.lock == 0 {
		actor := s.next()
		if actor == nil {
			return
		}
		actor.Act()
	}
}

func (s *Scheduler) next() Actor {
	return nil
}
