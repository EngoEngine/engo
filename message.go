package engi

type Message interface {
	Type() string
}

type CollisionMessage struct {
	Entity *Entity
}

func (collision CollisionMessage) Type() string {
	return "CollisionMessage"
}

type MessageManager struct {
	listeners map[string][]Systemer
}

func (mm *MessageManager) Dispatch(message Message) {
	systems := mm.listeners[message.Type()]

	for _, system := range systems {
		// println("Has message")
		system.Push(message)
	}
}

func (mm *MessageManager) Listen(messageType string, system Systemer) {
	if mm.listeners == nil {
		mm.listeners = make(map[string][]Systemer)
	}
	mm.listeners[messageType] = append(mm.listeners[messageType], system)
}
