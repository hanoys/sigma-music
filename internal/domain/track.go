package domain

import (
	"github.com/google/uuid"
)

type Track struct {
	ID       uuid.UUID
	AlbumID  uuid.UUID
	Name     string
	URL      string
}
