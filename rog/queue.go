package rog

type Queue struct {
	time       float32
	actors     []Actor
	actorTimes []float32
}

func NewQueue() *Queue {
	return &Queue{0, make([]Actor, 0, 1), make([]float32, 0, 1)}
}

func (q *Queue) Time() float32 {
	return q.time
}

func (q *Queue) Clear() {
	q.actors = make([]Actor, 0)
	q.actorTimes = make([]float32, 0)
}

func (q *Queue) Add(actor Actor, time float32) {
	index := len(q.actorTimes)
	for i, t := range q.actorTimes {
		if t > time {
			index = i
			break
		}
	}

	q.actors = append(q.actors, nil)
	copy(q.actors[index+1:], q.actors[index:])
	q.actors[index] = actor

	q.actorTimes = append(q.actorTimes, 0)
	copy(q.actorTimes[index+1:], q.actorTimes[index:])
	q.actorTimes[index] = time
}

func (q *Queue) Get() Actor {
	if len(q.actors) == 0 {
		return nil
	}

	var time float32
	time, q.actorTimes = q.actorTimes[0], q.actorTimes[1:]
	if time > 0 {
		q.time += time
		for i := 0; i < len(q.actorTimes); i++ {
			q.actorTimes[i] -= time
		}
	}

	var actor Actor
	actor, q.actors = q.actors[0], q.actors[1:]
	return actor
}

func (q *Queue) Remove(actor Actor) bool {
	for i, v := range q.actors {
		if v == actor {
			copy(q.actors[i:], q.actors[i+1:])
			q.actors[len(q.actors)-1] = nil
			q.actors = q.actors[:len(q.actors)-1]

			copy(q.actorTimes[i:], q.actorTimes[i+1:])
			q.actorTimes[len(q.actorTimes)-1] = 0
			q.actorTimes = q.actorTimes[:len(q.actorTimes)-1]
			return true
		}
	}
	return false
}
