package rog

import "testing"

type NumActor float32

func (a NumActor) Act() float32 {
	return float32(a)
}

func TestAdd(t *testing.T) {
	q := NewQueue()
	q.Add(NumActor(1), 100)
	if q.Next() != NumActor(1) {
		t.Errorf("should return added event")
	}
}

func TestEmpty(t *testing.T) {
	q := NewQueue()
	if q.Next() != nil {
		t.Errorf("should return null when no events are available")
	}
}

func TestPop(t *testing.T) {
	q := NewQueue()
	q.Add(NumActor(1), 0)
	q.Next()
	if q.Next() != nil {
		t.Errorf("should remove returned events")
	}
}

func TestRemove(t *testing.T) {
	q := NewQueue()
	q.Add(NumActor(1), 0)
	q.Add(NumActor(2), 0)
	if q.Remove(NumActor(1)) != true {
		t.Errorf("should remove events")
	}
	if q.Next() != NumActor(2) {
		t.Errorf("should remove events")
	}
}

func TestExisting(t *testing.T) {
	q := NewQueue()
	q.Add(NumActor(1), 0)
	if q.Remove(NumActor(2)) != false {
		t.Errorf("should survive removal of non-existant events")
	}
	if q.Next() != NumActor(1) {
		t.Errorf("should survive removal of non-existant events")
	}
}

func TestSorted(t *testing.T) {
	q := NewQueue()
	q.Add(NumActor(2), 10)
	q.Add(NumActor(1), 5)
	q.Add(NumActor(3), 15)
	if q.Next() != NumActor(1) {
		t.Errorf("should return events sorted")
	}
	if q.Next() != NumActor(2) {
		t.Errorf("should return events sorted")
	}
	if q.Next() != NumActor(3) {
		t.Errorf("should return events sorted")
	}
}

func TestTime(t *testing.T) {
	q := NewQueue()
	q.Add(NumActor(2), 10)
	q.Add(NumActor(1), 5)
	q.Add(NumActor(3), 15)
	q.Next()
	q.Next()
	q.Next()
	if q.Time() != 15 {
		t.Errorf("should compute elapsed time")
	}
}

func TestSameTimes(t *testing.T) {
	q := NewQueue()
	q.Add(NumActor(2), 10)
	q.Add(NumActor(1), 10)
	q.Add(NumActor(3), 10)
	if q.Next() != NumActor(2) {
		t.Errorf("should maintain event order for same timestamps")
	}
	if q.Next() != NumActor(1) {
		t.Errorf("should maintain event order for same timestamps")
	}
	if q.Next() != NumActor(3) {
		t.Errorf("should maintain event order for same timestamps")
	}
}
