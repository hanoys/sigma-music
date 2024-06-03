package entity

import (
	"github.com/Rhymond/go-money"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
	"time"
)

type PgOrder struct {
	ID         uuid.UUID `db:"id"`
	UserID     uuid.UUID `db:"user_id"`
	CreateTime time.Time `db:"create_time"`
	Price      float64   `db:"price"`
}

func (o *PgOrder) ToDomain() domain.Order {
	return domain.Order{
		ID:         o.ID,
		UserID:     o.UserID,
		CreateTime: o.CreateTime,
		Price:      money.NewFromFloat(o.Price, money.RUB),
	}
}

func NewPgOrder(order domain.Order) PgOrder {
	return PgOrder{
		ID:         order.ID,
		UserID:     order.UserID,
		CreateTime: order.CreateTime,
		Price:      order.Price.AsMajorUnits(),
	}
}
