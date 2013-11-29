package rog

// The Actor interface represents an entity that can take actions
// of a certain duration.
type Actor interface {
	// The return value determines how long the action will take.
	// Returning a value less than 0 removes the actor.
	Act() float32
}

// Queue schedules actors by how long their actions take.
// The Queue can be recursively locked in order to, for example,
// listen for user input or show a menu.
type Queue struct {
	lock       int
	time       float32
	actors     []Actor
	actorTimes []float32
}

func NewQueue() *Queue {
	return &Queue{0, 0, make([]Actor, 0, 1), make([]float32, 0, 1)}
}

// Lock the queue.
func (q *Queue) Lock() {
	q.lock++
}

// Unlock one lock of the queue.
func (q *Queue) Unlock() {
	if q.lock == 0 {
		panic("Cannot unlock unlocked Queue")
	}
	q.lock--
}

// Update the queue. If it isn't locked, grab the next actor,
// run it, and possibly requeue it.
func (q *Queue) Update() {
	if q.lock == 0 {
		actor := q.Next()
		if actor != nil {
			q.Add(actor, actor.Act())
		}
	}
}

// Return the current time of the queue.
func (q *Queue) Time() float32 {
	return q.time
}

// Add an actor to the queue with a delta time.
func (q *Queue) Add(actor Actor, time float32) {
	if time < 0 {
		return
	}

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

// Clear and unlock the queue.
func (q *Queue) Clear() {
	q.actors = make([]Actor, 0, 1)
	q.actorTimes = make([]float32, 0, 1)
	q.lock = 0
}

// Remove an actor from the queue.
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

// Return the next actor in the queue.
func (q *Queue) Next() Actor {
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
