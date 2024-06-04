package user

import "auth-api/src/internal/domain/events"

type UserRegisteredEvent struct {
	Email             string
	NeedsVerification bool
}

func (e *UserRegisteredEvent) GetType() events.EventType {
	return events.UserRegistered
}
