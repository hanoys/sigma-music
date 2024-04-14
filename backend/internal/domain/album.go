package domain

import (
	"github.com/google/uuid"
	"github.com/guregu/null/v5"
)

type Album struct {
	ID          uuid.UUID
	Name        string // len
	Description string // len
	Published   bool
	ReleaseDate null.Time
}
