package events

type EventType string

const (
	UserRegistered EventType = "user_registered"
)

type Event interface {
	GetType() EventType
}

type EventHandler interface {
	Handle(event Event) error
}

type EventDispatcher interface {
	Register(eventType EventType, handler EventHandler)
	Dispatch(event Event) error
}

type UserRegisteredEvent struct {
	Email             string
	NeedsVerification bool
}

func (e *UserRegisteredEvent) GetType() EventType {
	return UserRegistered
}
