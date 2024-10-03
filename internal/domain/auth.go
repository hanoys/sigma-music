package domain

import "github.com/google/uuid"

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type Payload struct {
	UserID uuid.UUID
	Role   int
}
