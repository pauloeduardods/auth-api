package events

type EventType string

type Event interface {
	GetType() EventType
	Validate() error
}

type EventHandler interface {
	Handle(event Event) error
}

type EventDispatcher interface {
	Register(eventType EventType, handler EventHandler)
	Dispatch(event Event) error
}
