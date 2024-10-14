package domain

import "github.com/google/uuid"

type Comment struct {
	ID      uuid.UUID
	UserID  uuid.UUID
	TrackID uuid.UUID
	Stars   int
	Text    string
}
