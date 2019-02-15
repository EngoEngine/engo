# Textbox Demo

## What does it do?
It demonstrates how one can get user text input and print it to the screen.

## What are important aspects of the code?
To listen to the character callbacks from the back-end (glfw, js, etc) subscribe
to engo.TextMessage

```go
engo.Mailbox.Listen("TextMessage", func(msg engo.Message) {
  m, ok := msg.(engo.TextMessage)
  if !ok {
    return
  }
  t.runeLock.Lock()
  t.runesToAdd = append(t.runesToAdd, m.Char)
  t.runeLock.Unlock()
})
```

Once you're subscribed to the message, you can get the characters as they're
typed as m.Char
