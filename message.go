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

// Remove removes a specific listener
func (mm *MessageManager) Remove(messageType string, handler MessageHandler) {
	// TODO: ...
}

// Remove removes all listeners of a message type
func (mm *MessageManager) RemoveType(messageType string) {
	delete(mm.listeners, messageType)
}

// WindowResizeMessage is a message that's being dispatched whenever the game window is being resized by the gamer
type WindowResizeMessage struct {
	OldWidth, OldHeight int
	NewWidth, NewHeight int
}

// IterationUpdateMessage will be used with mailbox
type IterationUpdateMessage struct {
	Delta float32
}

// PreparationMessage will be used with mailbox
type PreparationMessage struct{}

// Type of message that will be sent
func (WindowResizeMessage) Type() string {
	return "engo.WindowResizeMessage"
}

// Type of message that will be sent
func (IterationUpdateMessage) Type() string {
	return "engo.IterationUpdateMessage"
}

// Type of message that will be sent
func (PreparationMessage) Type() string {
	return "engo.PreparationMessage"
}
