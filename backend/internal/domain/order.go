package domain

import (
	"github.com/Rhymond/go-money"
	"time"
)

type Order struct {
	ID         int
	UserID     int
	CreateTime time.Time
	Price      *money.Money
}
