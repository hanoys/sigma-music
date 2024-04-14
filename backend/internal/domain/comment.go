package domain

import "github.com/google/uuid"

type Comment struct {
	ID      uuid.UUID
	UserID  uuid.UUID
	TrackID uuid.UUID
	Stars   int    // > 0 && < 5
	Text    string // len
}
