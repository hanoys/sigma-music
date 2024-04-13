package domain

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID
	Name     string
	Email    string
	Phone    string
	Password string
	Country  string
}
