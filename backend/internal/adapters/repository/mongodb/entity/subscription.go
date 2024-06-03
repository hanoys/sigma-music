package entity

import (
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
	"time"
)

type MongoSubscription struct {
	ID             uuid.UUID `bson:"_id"`
	UserID         uuid.UUID `bson:"user_id"`
	OrderID        uuid.UUID `bson:"order_id"`
	StartDate      time.Time `bson:"start_date"`
	ExpirationDate time.Time `bson:"expiration_date"`
}

func (s *MongoSubscription) ToDomain() domain.Subscription {
	return domain.Subscription{
		ID:             s.ID,
		UserID:         s.UserID,
		OrderID:        s.OrderID,
		StartDate:      s.StartDate,
		ExpirationDate: s.ExpirationDate,
	}
}

func NewMongoSubscription(subscription domain.Subscription) MongoSubscription {
	return MongoSubscription{
		ID:             subscription.ID,
		UserID:         subscription.UserID,
		OrderID:        subscription.OrderID,
		StartDate:      subscription.StartDate,
		ExpirationDate: subscription.ExpirationDate,
	}
}
