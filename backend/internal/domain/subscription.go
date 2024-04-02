package domain

import "time"

type Subscription struct {
	ID             int
	UserID         int
	OrderID        int
	StartDate      time.Time
	ExpirationDate time.Time
}
