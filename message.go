package engo

import (
	"sync"
)

//A MessageHandler is used to dispatch a message to the subscribed handler.
type MessageHandler func(msg Message)

// MessageHandlerId is used to track handlers, each handler will get a unique ID
type MessageHandlerId uint64

var currentHandlerID MessageHandlerId

func init() {
	currentHandlerID = 0
}

func getNewHandlerID() MessageHandlerId {
	currentHandlerID++
	return currentHandlerID
}

type HandlerIDPair struct {
	MessageHandlerId
	MessageHandler
}

// A Message is used to send messages within the MessageManager
type Message interface {
	Type() string
}

// MessageManager manages messages and subscribed handlers
type MessageManager struct {
	// this mutex will prevent race
	// conditions on listeners and
	// sync its state across the game
	sync.RWMutex
	listeners        map[string][]HandlerIDPair
	handlersToRemove map[string][]MessageHandlerId
}

// Dispatch sends a message to all subscribed handlers of the message's type
// To prevent any data races, be aware that these listeners occur as callbacks and can be
// executed at any time. If variables are altered in the handler, utilize channels, locks,
// semaphores, or any other method necessary to ensure the memory is not altered by multiple
// functions simultaneously.
func (mm *MessageManager) Dispatch(message Message) {
	mm.RLock()
	mm.clearRemovedHandlers()
	handlers := make([]MessageHandler, len(mm.listeners[message.Type()]))
	pairs := mm.listeners[message.Type()]
	for i := range pairs {
		handlers[i] = pairs[i].MessageHandler
	}
	mm.RUnlock()

	for _, handler := range handlers {
		handler(message)
	}

}

// Listen subscribes to the specified message type and calls the specified handler when fired
func (mm *MessageManager) Listen(messageType string, handler MessageHandler) MessageHandlerId {
	mm.Lock()
	defer mm.Unlock()
	if mm.listeners == nil {
		mm.listeners = make(map[string][]HandlerIDPair)
	}
	handlerID := getNewHandlerID()
	newHandlerIDPair := HandlerIDPair{MessageHandlerId: handlerID, MessageHandler: handler}
	mm.listeners[messageType] = append(mm.listeners[messageType], newHandlerIDPair)
	return handlerID
}

// ListenOnce is a convenience wrapper around StopListen() to only listen to a specified message once
func (mm *MessageManager) ListenOnce(messageType string, handler MessageHandler) {
	handlerID := MessageHandlerId(0)
	handlerID = mm.Listen(messageType, func(msg Message) {
		handler(msg)
		mm.StopListen(messageType, handlerID)
	})
}

// StopListen removes a previously added handler from the listener queue
func (mm *MessageManager) StopListen(messageType string, handlerID MessageHandlerId) {
	if mm.handlersToRemove == nil {
		mm.handlersToRemove = make(map[string][]MessageHandlerId)
	}
	mm.handlersToRemove[messageType] = append(mm.handlersToRemove[messageType], handlerID)
}

// Will deleted all queued handlers that are scheduled for removal due to StopListen()
func (mm *MessageManager) clearRemovedHandlers() {
	for messageType, handlerList := range mm.handlersToRemove {
		for _, handlerID := range handlerList {
			mm.removeHandler(messageType, handlerID)
		}
	}
	mm.handlersToRemove = make(map[string][]MessageHandlerId)
}

// Removes a single handler from the handler queue, called during cleanup of all handlers scheduled for removal
func (mm *MessageManager) removeHandler(messageType string, handlerID MessageHandlerId) {
	indexOfHandler := -1
	for i, activeHandler := range mm.listeners[messageType] {
		if activeHandler.MessageHandlerId == handlerID {
			indexOfHandler = i
			break
		}
	}
	// A handler might have already been removed during a previous Dispatch(), no action necessary
	if indexOfHandler == -1 {
		return
	}
	mm.listeners[messageType] = append(mm.listeners[messageType][:indexOfHandler], mm.listeners[messageType][indexOfHandler+1:]...)
}

// WindowResizeMessage is a message that's being dispatched whenever the game window is being resized by the gamer
type WindowResizeMessage struct {
	OldWidth, OldHeight int
	NewWidth, NewHeight int
}

// Type returns the type of the current object "WindowResizeMessage"
func (WindowResizeMessage) Type() string { return "WindowResizeMessage" }

// TextMessage is a message that is dispatched whenever a character is typed on the
// keyboard. This is not the same as a keypress, as it returns the rune of the
// character typed by the user, which could be a combination of keypresses.
type TextMessage struct {
	Char rune
}

// Type returns the type of the message, "TextMessage"
func (TextMessage) Type() string { return "TextMessage" }
