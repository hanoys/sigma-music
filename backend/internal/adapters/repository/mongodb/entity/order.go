package entity

import (
	"github.com/Rhymond/go-money"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
	"time"
)

type MongoOrder struct {
	ID         uuid.UUID `bson:"_id"`
	UserID     uuid.UUID `bson:"user_id"`
	CreateTime time.Time `bson:"create_time"`
	Price      float64   `bson:"price"`
}

func (o *MongoOrder) ToDomain() domain.Order {
	return domain.Order{
		ID:         o.ID,
		UserID:     o.UserID,
		CreateTime: o.CreateTime,
		Price:      money.NewFromFloat(o.Price, money.RUB),
	}
}

func NewMongoOrder(order domain.Order) MongoOrder {
	return MongoOrder{
		ID:         order.ID,
		UserID:     order.UserID,
		CreateTime: order.CreateTime,
		Price:      order.Price.AsMajorUnits(),
	}
}
