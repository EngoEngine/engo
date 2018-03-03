package engo

import "testing"

type testMessage struct {
	success bool
}

func (testMessage) Type() string {
	return "testMessage"
}

func TestMessages(t *testing.T) {
	mailbox := &MessageManager{}
	msg := testMessage{}
	mailbox.Listen("testMessage", func(message Message) {
		m, ok := message.(*testMessage)
		if ok {
			m.success = true
		}
	})
	mailbox.Dispatch(&msg)
	if !msg.success {
		t.Error("huh")
	}
}
