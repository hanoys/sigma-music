package domain

import (
	"github.com/google/uuid"
	"github.com/guregu/null/v5"
)

type Musician struct {
	ID          uuid.UUID
	Name        string
	Email       string
	Password    string
	Salt        string
	Country     string
	Description string
	ImageURL    null.String
}
