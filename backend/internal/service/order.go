package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
	"time"
)

type OrderService struct {
	repository ports.IOrderRepository
}

func NewOrderService(repo ports.IOrderRepository) *OrderService {
	return &OrderService{repository: repo}
}

func (os *OrderService) Create(ctx context.Context, orderReq ports.CreateOrderReq) (domain.Order, error) {
	newOrder := domain.Order{
		ID:         uuid.New(),
		UserID:     orderReq.UserID,
		CreateTime: time.Now(),
		Price:      orderReq.Price,
	}

	return os.repository.Create(ctx, newOrder)
}
