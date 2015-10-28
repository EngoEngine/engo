package engi

// UnpauseComponent is a component that indicates whether or not the Entity should be affected by
// system-wide pauses.
type UnpauseComponent struct{}

func (UnpauseComponent) Name() string {
	return "UnpauseComponent"
}

// PauseSystem is a Systemer that listens for Pause messages, and then pauses the entire world
type PauseSystem struct {
	*System
	// TODO: can we handle this any other way?
	World *World
}

func (ps *PauseSystem) New() {
	ps.System = &System{}
	Mailbox.Listen("PauseMessage", func(message Message) {
		pm, ok := message.(PauseMessage)
		if !ok {
			return
		}
		ps.World.paused = pm.Pause
	})
}

func (*PauseSystem) Update(*Entity, float32) {}

func (PauseSystem) Name() string {
	return "PauseSystem"
}

// PauseMessage is a message that is sent to (un)pause the world
type PauseMessage struct {
	Pause bool
}

func (PauseMessage) Type() string {
	return "PauseMessage"
}
