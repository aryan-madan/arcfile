package structures

import "time"

type File struct {
	Identifier string
	Filename   string
	UUID       string
	CreatedAt  time.Time
	ExpiresAt  time.Time
	Email      string
}
