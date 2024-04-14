package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
	"time"
)

type SubscriptionService struct {
	repository ports.ISubscriptionRepository
}

func NewSubscriptionService(repo ports.ISubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{repository: repo}
}

func (ss *SubscriptionService) Create(ctx context.Context, subReq ports.CreateSubscriptionReq) (domain.Subscription, error) {
	newSubscription := domain.Subscription{
		ID:             uuid.New(),
		UserID:         subReq.UserID,
		OrderID:        subReq.OrderID,
		StartDate:      time.Now(),
		ExpirationDate: time.Now().Add(time.Hour * 30),
	}

	return ss.repository.Create(ctx, newSubscription)
}
