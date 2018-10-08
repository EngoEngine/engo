package engo

import (
	"errors"
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
	listeners map[string][]HandlerIDPair
}

// Dispatch sends a message to all subscribed handlers of the message's type
func (mm *MessageManager) Dispatch(message Message) {
	mm.Lock()
	handlers := mm.listeners[message.Type()]
	mm.Unlock()

	mm.RLock()
	defer mm.RUnlock()
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

// StopListen removes a previously added handler from the listener queue
func (mm *MessageManager) StopListen(messageType string, handlerId MessageHandlerId) error {
	// first, we need to find the handler to remove or return an error if that's not possible
	indexOfHandler := -1
	for i, activeHandler := range mm.listeners[messageType] {
		if activeHandler.MessageHandlerId == handlerId {
			indexOfHandler = i
			break
		}
	}
	if indexOfHandler == -1 {
		return errors.New("Trying to remove handler that does not exist")
	}
	mm.listeners[messageType] = append(mm.listeners[messageType][:indexOfHandler], mm.listeners[messageType][indexOfHandler+1:]...)
	return nil
}

func (mm *MessageManager) ListenOnce(messageType string, handler MessageHandler) {
	handlerId := MessageHandlerId(0)
	handlerId = mm.Listen(messageType, func(msg Message) {
		handler(msg)
		mm.StopListen(messageType, handlerId)
	})
}

// WindowResizeMessage is a message that's being dispatched whenever the game window is being resized by the gamer
type WindowResizeMessage struct {
	OldWidth, OldHeight int
	NewWidth, NewHeight int
}

// Type returns the type of the current object "WindowResizeMessage"
func (WindowResizeMessage) Type() string { return "WindowResizeMessage" }
