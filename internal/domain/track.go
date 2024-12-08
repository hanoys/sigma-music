package domain

import (
	"github.com/google/uuid"
	"github.com/guregu/null/v5"
)

type Track struct {
	ID       uuid.UUID
	AlbumID  uuid.UUID
	Name     string
	URL      string
	ImageURL null.String
}
