package structures

import "time"

type File struct {
	ID int64
	Identifier string
	Name string
	Path string
	CreatedAt time.Time
	ExpiresAt time.Time
	Email string
}