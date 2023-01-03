package events

type Code int

const (
	EventNavigateLeft Code = iota
	EventNavigateRight
	EventNavigateUp
	EventNavigateDown
)

// Payload is additional data sent with an event.
type Payload struct {
	Sender any
	Data   any
}

// Listener is an interface for types that can subscribe to events.
type Listener interface {
	HandleEvent(code Code, payload Payload)
}
