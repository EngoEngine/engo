package engo

import "testing"

type testMessageCounter struct {
	counter int
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
	handlerId := mailbox.Listen("testMessageCounter", func(message Message) {
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

	mailbox.StopListen("testMessageCounter", handlerId)

	mailbox.Dispatch(&msg)
	if msg.counter != 1 {
		t.Error("Message should have been received exactly 1 times since its handler was removed from listeners")
	}
}

func TestRemovalOfNonexistentHandler(t *testing.T) {
	mailbox := &MessageManager{}
	err := mailbox.StopListen("testMessageCounter", MessageHandlerId(42))
	if err == nil {
		t.Error("StopListen should return an error in case the handler to remove doesn't exist")
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
	mailbox.Dispatch(&msg)
	if msg.counter != 1 {
		t.Error("Message should have been received 1 times by now")
	}

	mailbox.Dispatch(&msg)
	if msg.counter != 1 {
		t.Error("Message should have been received exactly 1 times since its been added by ListenOnce()")
	}
}
