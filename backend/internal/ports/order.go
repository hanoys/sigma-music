package ports

import (
	"context"
	"errors"
	"github.com/Rhymond/go-money"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
)

var (
	ErrOrderDuplicate    = errors.New("order duplicate error")
	ErrOrderIDNotFound   = errors.New("order with such id not found")
	ErrInternalOrderRepo = errors.New("internal order repository error")
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
