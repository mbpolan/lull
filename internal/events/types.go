package events

type Code int

const (
	EventNavigateLeft Code = iota
	EventNavigateRight
	EventNavigateUp
	EventNavigateDown
	EventStatusBarContextChange
)

// Payload is additional data sent with an event.
type Payload struct {
	Sender any
	Data   any
}

// StatusBarContextChangeData contains the list of shortcuts to show in the status bar when focus changes to another
// controlling primitive.
type StatusBarContextChangeData struct {
	Fields []StatusBarContextChangeSequence
}

// StatusBarContextChangeSequence is a key sequence that is recognized by the focused primitive.
type StatusBarContextChangeSequence struct {
	Label       string
	KeySequence string
}

// Listener is an interface for types that can subscribe to events.
type Listener interface {
	HandleEvent(code Code, payload Payload)
}
