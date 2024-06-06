package user

import "auth-api/src/internal/events"

const (
	UserRegistered events.EventType = "UserRegistered"
	// UserConfirmed  events.EventType = "UserConfirmed"
)

type UserRegisteredEvent struct {
	Email             string
	NeedsVerification bool
}

func (e *UserRegisteredEvent) GetType() events.EventType {
	return UserRegistered
}

func (e *UserRegisteredEvent) Validate() error {
	if e.Email == "" {
		return ErrInvalidEmail
	}
	return nil
}
