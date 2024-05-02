package adapters

import (
	"context"
	"encoding/json"
	"github.com/hanoys/sigma-music/internal/adapters/auth/ports"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/util"
	"github.com/redis/go-redis/v9"
	"time"
)

type TokenStorage struct {
	redisClient *redis.Client
}

func (ts *TokenStorage) Set(ctx context.Context, key string, payload domain.Payload, expiration time.Duration) error {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return util.WrapError(ports.ErrInternalTokenStorage, err)
	}

	_, err = ts.redisClient.Set(ctx, key, payloadJSON, expiration).Result()
	if err != nil {
		return util.WrapError(ports.ErrInternalTokenStorage, err)
	}

	return nil
}

func (ts *TokenStorage) Del(ctx context.Context, key string) error {
	ok, err := ts.redisClient.Del(ctx, key).Result()
	if err != nil || ok != 1 {
		return util.WrapError(ports.ErrInternalTokenStorage, err)
	}

	return nil
}
