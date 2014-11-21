package engi

type Event interface {
	Type() string
	Data() interface{}
}

type BasicEvent struct {
	Info interface{}
}

func (basic BasicEvent) Type() string {
	return "BasicEvent"
}

func (basic BasicEvent) Data() interface{} {
	return basic.Info
}
