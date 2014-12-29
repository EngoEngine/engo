package engi

type Handler func(i interface{})

type Message interface {
	Type() string
}

type CollisionMessage struct {
	Entity *Entity
	To     *Entity
}

func (collision CollisionMessage) Type() string {
	return "CollisionMessage"
}

type MessageManager struct {
	listeners map[string][]Handler
}

func (mm *MessageManager) Dispatch(name string, message interface{}) {
	handlers := mm.listeners[name]

	for _, handler := range handlers {
		handler(message)
	}
}

func (mm *MessageManager) Listen(messageType string, handler Handler) {
	if mm.listeners == nil {
		mm.listeners = make(map[string][]Handler)
	}
	mm.listeners[messageType] = append(mm.listeners[messageType], handler)
}
