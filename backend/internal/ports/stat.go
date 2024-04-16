package ports

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
)

var (
	ErrInternalStatRepo = errors.New("internal statistics repository error")
)

type IStatRepository interface {
	Add(ctx context.Context, userID uuid.UUID, trackID uuid.UUID) error
	GetMostListenedMusicians(ctx context.Context, userID uuid.UUID, maxCnt int) ([]domain.UserMusiciansStat, error)
	GetListenedGenres(ctx context.Context, userID uuid.UUID) ([]domain.UserGenresStat, error)
}

type IStatService interface {
	Add(ctx context.Context, userID uuid.UUID, trackID uuid.UUID) error
	FormReport(ctx context.Context, userID uuid.UUID) (domain.ListenReport, error)
}
