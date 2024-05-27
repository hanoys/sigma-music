package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
	"go.uber.org/zap"
	"time"
)

type SubscriptionService struct {
	repository ports.ISubscriptionRepository
	logger     *zap.Logger
}

func NewSubscriptionService(repo ports.ISubscriptionRepository, logger *zap.Logger) *SubscriptionService {
	return &SubscriptionService{
		repository: repo,
		logger:     logger,
	}
}

func (ss *SubscriptionService) Create(ctx context.Context, subReq ports.CreateSubscriptionReq) (domain.Subscription, error) {
	newSubscription := domain.Subscription{
		ID:             uuid.New(),
		UserID:         subReq.UserID,
		OrderID:        subReq.OrderID,
		StartDate:      time.Now(),
		ExpirationDate: time.Now().Add(time.Hour * 30),
	}

	sub, err := ss.repository.Create(ctx, newSubscription)
	if err != nil {
		ss.logger.Error("Failed to create subscription", zap.Error(err),
			zap.String("User ID", subReq.UserID.String()))

		return domain.Subscription{}, err
	}

	ss.logger.Info("Subscription successfully created", zap.String("Subscription ID", newSubscription.ID.String()),
		zap.String("User ID", newSubscription.UserID.String()))

	return sub, nil
}
