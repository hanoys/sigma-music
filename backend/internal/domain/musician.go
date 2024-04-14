package domain

import "github.com/google/uuid"

type Musician struct {
	ID          uuid.UUID
	Name        string // len
	Email       string // len
	Password    string // len
	Country     string // len
	Description string // len
}
