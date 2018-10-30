package engo

import (
	"sync"
)

//A MessageHandler is used to dispatch a message to the subscribed handler.
type MessageHandler func(msg Message)

// in order to track handlers, each handler will get a unique ID
type MessageHandlerId uint64

var currentHandlerId MessageHandlerId

func init() {
	currentHandlerId = 0
}

func getNewHandlerId() MessageHandlerId {
	currentHandlerId++
	return currentHandlerId
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
	mm.Lock()
	mm.clearRemovedHandlers()
	handlers := mm.listeners[message.Type()]
	mm.Unlock()

	for _, handler := range handlers {
		handler.MessageHandler(message)
	}

}

// Listen subscribes to the specified message type and calls the specified handler when fired
func (mm *MessageManager) Listen(messageType string, handler MessageHandler) MessageHandlerId {
	mm.Lock()
	defer mm.Unlock()
	if mm.listeners == nil {
		mm.listeners = make(map[string][]HandlerIDPair)
	}
	handlerID := getNewHandlerId()
	newHandlerIdPair := HandlerIDPair{MessageHandlerId: handlerID, MessageHandler: handler}
	mm.listeners[messageType] = append(mm.listeners[messageType], newHandlerIdPair)
	return handlerID
}

// ListenOnce is a convenience wrapper around StopListen() to only listen to a specified message once
func (mm *MessageManager) ListenOnce(messageType string, handler MessageHandler) {
	handlerId := MessageHandlerId(0)
	handlerId = mm.Listen(messageType, func(msg Message) {
		handler(msg)
		mm.StopListen(messageType, handlerId)
	})
}

// StopListen removes a previously added handler from the listener queue
func (mm *MessageManager) StopListen(messageType string, handlerId MessageHandlerId) {
	if mm.handlersToRemove == nil {
		mm.handlersToRemove = make(map[string][]MessageHandlerId)
	}
	mm.handlersToRemove[messageType] = append(mm.handlersToRemove[messageType], handlerId)
}

// Will deleted all queued handlers that are scheduled for removal due to StopListen()
func (mm *MessageManager) clearRemovedHandlers() {
	for messageType, handlerList := range mm.handlersToRemove {
		for _, handlerId := range handlerList {
			mm.removeHandler(messageType, handlerId)
		}
	}
	mm.handlersToRemove = make(map[string][]MessageHandlerId)
}

// Removes a single handler from the handler queue, called during cleanup of all handlers scheduled for removal
func (mm *MessageManager) removeHandler(messageType string, handlerId MessageHandlerId) {
	indexOfHandler := -1
	for i, activeHandler := range mm.listeners[messageType] {
		if activeHandler.MessageHandlerId == handlerId {
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
