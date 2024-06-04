package code

import "time"

type Code struct {
	Value      string
	ExpiresAt  time.Time
	Identifier string
}

func (c *Code) IsExpired() bool {
	return c.ExpiresAt.Before(time.Now())
}
