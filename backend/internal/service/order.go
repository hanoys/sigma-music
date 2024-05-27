package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
	"go.uber.org/zap"
	"time"
)

type OrderService struct {
	repository ports.IOrderRepository
	logger     *zap.Logger
}

func NewOrderService(repo ports.IOrderRepository, logger *zap.Logger) *OrderService {
	return &OrderService{
		repository: repo,
		logger:     logger,
	}
}

func (os *OrderService) Create(ctx context.Context, orderReq ports.CreateOrderReq) (domain.Order, error) {
	newOrder := domain.Order{
		ID:         uuid.New(),
		UserID:     orderReq.UserID,
		CreateTime: time.Now(),
		Price:      orderReq.Price,
	}

	order, err := os.repository.Create(ctx, newOrder)
	if err != nil {
		os.logger.Error("Failed to create order", zap.Error(err),
			zap.String("User ID", orderReq.UserID.String()), zap.String("Price", orderReq.Price.Display()))

		return domain.Order{}, err
	}

	os.logger.Info("Order successfully created", zap.String("Order ID", newOrder.ID.String()),
		zap.String("User ID", orderReq.UserID.String()), zap.String("Price", orderReq.Price.Display()))

	return order, nil
}
