package domain

import (
	"github.com/google/uuid"
	"github.com/guregu/null/v5"
)

type Album struct {
	ID          uuid.UUID
	Name        string
	Description string
	Published   bool
	ReleaseDate null.Time
	ImageURL    null.String
}
