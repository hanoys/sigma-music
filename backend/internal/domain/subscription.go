package domain

import (
	"github.com/google/uuid"
	"time"
)

type Subscription struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	OrderID        uuid.UUID
	StartDate      time.Time
	ExpirationDate time.Time
}
