package structures

import "time"

type File struct {
	ID         int64
	Identifier string
	Filename   string
	Path       string
	CreatedAt  time.Time
	ExpiresAt  time.Time
	Email      string
}
