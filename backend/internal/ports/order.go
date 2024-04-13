package ports

import (
	"context"
	"errors"
	"github.com/Rhymond/go-money"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
)

var (
	ErrOrderCreate = errors.New("can't create order: internal error")
)

type IOrderRepository interface {
	Create(ctx context.Context, order domain.Order) (domain.Order, error)
}

type CreateOrderReq struct {
	UserID uuid.UUID
	Price  *money.Money
}

type IOrderService interface {
	Create(ctx context.Context, orderReq CreateOrderReq) (domain.Order, error)
}
