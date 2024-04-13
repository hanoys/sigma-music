package ports

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
)

var (
	ErrSubscriptionCreate = errors.New("can't create subscription: internal error")
)

type ISubscriptionRepository interface {
	Create(ctx context.Context, sub domain.Subscription) (domain.Subscription, error)
}

type CreateSubscriptionReq struct {
	UserID  uuid.UUID
	OrderID uuid.UUID
}

type ISubscriptionService interface {
	Create(ctx context.Context, subReq CreateSubscriptionReq) (domain.Subscription, error)
}
