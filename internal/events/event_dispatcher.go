package events

var instance EventDispatcher

// EventDispatcher is a synchronous manager that allows senders to dispatch events to subscribers.
type EventDispatcher struct {
	listeners map[Code][]Listener
}

// Setup initializes the event dispatcher. This should be called once on start up.
func Setup() {
	instance = EventDispatcher{
		listeners: map[Code][]Listener{},
	}
}

// Dispatcher returns an instance of EventDispatcher.
func Dispatcher() *EventDispatcher {
	return &instance
}

// Subscribe adds a listener that will be invoked when an event with a matching code is posted.
func (ed *EventDispatcher) Subscribe(listener Listener, codes []Code) {
	for _, i := range codes {
		if ed.listeners[i] == nil {
			ed.listeners[i] = []Listener{}
		}

		ed.listeners[i] = append(ed.listeners[i], listener)
	}
}

// PostSimple posts an event without any additional event data.
func (ed *EventDispatcher) PostSimple(code Code, sender any) {
	ed.Post(code, sender, nil)
}

// Post posts an event with additional event data.
func (ed *EventDispatcher) Post(code Code, sender any, data any) {
	listeners := ed.listeners[code]
	if listeners == nil {
		return
	}

	for _, l := range listeners {
		l.HandleEvent(code, Payload{
			Sender: sender,
			Data:   data,
		})
	}
}
