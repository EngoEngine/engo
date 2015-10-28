package engi

type MessageHandler func(msg Message)

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
	listeners map[string][]MessageHandler
}

func (mm *MessageManager) Dispatch(message Message) {
	handlers := mm.listeners[message.Type()]

	for _, handler := range handlers {
		handler(message)
	}
}

func (mm *MessageManager) Listen(messageType string, handler MessageHandler) {
	if mm.listeners == nil {
		mm.listeners = make(map[string][]MessageHandler)
	}
	mm.listeners[messageType] = append(mm.listeners[messageType], handler)
}
