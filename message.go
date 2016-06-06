package engo

//A MessageHandler is used to dispatch a message to the subscribed handler.
type MessageHandler func(msg Message)

// A Message is used to send messages within the MessageManager
type Message interface {
	Type() string
}

// MessageManager manages messages and subscribed handlers
type MessageManager struct {
	listeners map[string][]MessageHandler
}

// Dispatch sends a message to all subscribed handlers of the message's type
func (mm *MessageManager) Dispatch(message Message) {
	handlers := mm.listeners[message.Type()]

	for _, handler := range handlers {
		handler(message)
	}
}

// Listen subscribes to the specified message type and calls the specified handler when fired
func (mm *MessageManager) Listen(messageType string, handler MessageHandler) {
	if mm.listeners == nil {
		mm.listeners = make(map[string][]MessageHandler)
	}
	mm.listeners[messageType] = append(mm.listeners[messageType], handler)
}

// WindowResizeMessage is a message that's being dispatched whenever the game window is being resized by the gamer
type WindowResizeMessage struct {
	OldWidth, OldHeight int
	NewWidth, NewHeight int
}

// Type returns the type of the current object "WindowResizeMessage"
func (WindowResizeMessage) Type() string { return "WindowResizeMessage" }
