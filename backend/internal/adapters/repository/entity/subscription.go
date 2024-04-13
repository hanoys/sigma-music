package entity

import (
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
	"time"
)

type PgSubscription struct {
	ID             uuid.UUID `db:"id"`
	UserID         uuid.UUID `db:"user_id"`
	OrderID        uuid.UUID `db:"order_id"`
	StartDate      time.Time `db:"start_date"`
	ExpirationDate time.Time `db:"expiration_date"`
}

func (s *PgSubscription) ToDomain() domain.Subscription {
	return domain.Subscription{
		ID:             s.ID,
		UserID:         s.UserID,
		OrderID:        s.OrderID,
		StartDate:      s.StartDate,
		ExpirationDate: s.ExpirationDate,
	}
}

func NewPgSuscription(subscription domain.Subscription) PgSubscription {
	return PgSubscription{
		ID:             subscription.ID,
		UserID:         subscription.UserID,
		OrderID:        subscription.OrderID,
		StartDate:      subscription.StartDate,
		ExpirationDate: subscription.ExpirationDate,
	}
}
