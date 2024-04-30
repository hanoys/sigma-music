package ports

import (
	"context"
	"errors"
	"github.com/hanoys/sigma-music/internal/domain"
	"time"
)

var (
	ErrInternalTokenStorage = errors.New("internal token storage error")
)

type ITokenStorage interface {
	Set(ctx context.Context, key string, payload domain.Payload, expiration time.Duration) error
	Del(ctx context.Context, key string) error
}
