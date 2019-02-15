package engo

import "testing"

type testMessageCounter struct {
	counter, counter2 int
}

func (testMessageCounter) Type() string {
	return "testMessageCounter"
}

func TestMessageCounterSimple(t *testing.T) {
	mailbox := &MessageManager{}
	msg := testMessageCounter{}
	mailbox.Listen("testMessageCounter", func(message Message) {
		m, ok := message.(*testMessageCounter)
		if !ok {
			t.Error("Message should be of type testMessageCounter")
		}
		m.counter++
	})
	mailbox.Dispatch(&msg)
	if msg.counter != 1 {
		t.Error("Message should have been received 1 times by now")
	}
	mailbox.Dispatch(&msg)
	if msg.counter != 2 {
		t.Error("Message should have been received 2 times by now")
	}
}

func TestMessageCounterWithRemoval(t *testing.T) {
	mailbox := &MessageManager{}
	msg := testMessageCounter{}
	handlerID := mailbox.Listen("testMessageCounter", func(message Message) {
		m, ok := message.(*testMessageCounter)
		if !ok {
			t.Error("Message should be of type testMessageCounter")
		}
		m.counter++
	})
	mailbox.Dispatch(&msg)
	if msg.counter != 1 {
		t.Error("Message should have been received 1 times by now")
	}

	mailbox.StopListen("testMessageCounter", handlerID)

	mailbox.Dispatch(&msg)
	if msg.counter != 1 {
		t.Error("Message should have been received exactly 1 times since its handler was removed from listeners")
	}
}

func TestMessageListenOnce(t *testing.T) {
	mailbox := &MessageManager{}
	msg := testMessageCounter{}
	mailbox.ListenOnce("testMessageCounter", func(message Message) {
		m, ok := message.(*testMessageCounter)
		if !ok {
			t.Error("Message should be of type testMessageCounter")
		}
		m.counter++
	})
	mailbox.Listen("testMessageCounter", func(message Message) {
		m, ok := message.(*testMessageCounter)
		if !ok {
			t.Error("Message should be of type TestMessageCounter")
		}
		m.counter2++
	})
	mailbox.Dispatch(&msg)
	if msg.counter != 1 {
		t.Error("Message should have been received 1 times by now")
	}
	if msg.counter2 != 1 {
		t.Error("Message should have been recieved by second listener 1 time by now")
	}
	mailbox.Dispatch(&msg)
	if msg.counter != 1 {
		t.Error("Message should have been received exactly 1 times since its been added by ListenOnce()")
	}
	if msg.counter2 != 2 {
		t.Error("Message should have been recieved exactly 2 times by second listener")
	}
}

func TestMessageRemoveNotExistAndMultipleMessageTypes(t *testing.T) {
	mailbox := &MessageManager{}
	msg := testMessageCounter{}
	resizeMsg := WindowResizeMessage{}
	textMsg := TextMessage{}
	mailbox.Listen("testMessageCounter", func(message Message) {
		m, ok := message.(*testMessageCounter)
		if !ok {
			t.Error("Message should be of type TestMessageCounter")
		}
		m.counter++
	})
	mailbox.StopListen("TestingNonExistantMessage", 10)
	mailbox.Dispatch(&resizeMsg)
	mailbox.Dispatch(&textMsg)
	if msg.counter != 0 {
		t.Error("Message counter should be 0 since no messages were dispatched to it")
	}
	mailbox.Dispatch(&msg)
	if msg.counter != 1 {
		t.Error("Message counter should be 1. Only one message was dispatched to it")
	}
}
