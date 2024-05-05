package domain

import "github.com/google/uuid"

type Musician struct {
	ID          uuid.UUID
	Name        string
	Email       string
	Password    string
	Salt        string
	Country     string
	Description string
}
