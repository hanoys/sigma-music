package domain

import (
	"github.com/Rhymond/go-money"
	"github.com/google/uuid"
	"time"
)

type Order struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	CreateTime time.Time
	Price      *money.Money // > 0
}
