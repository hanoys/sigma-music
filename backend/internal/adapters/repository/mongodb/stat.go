package mongodb

import (
	"context"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoStatRepository struct {
	db *mongo.Collection
}

func NewMongoStatRepository(db *mongo.Database) *MongoStatRepository {
	return &MongoStatRepository{db: db.Collection(StatCollection)}
}

func (sr *MongoStatRepository) Add(ctx context.Context, userID uuid.UUID, trackID uuid.UUID) error {
	return nil
}

func (sr *MongoStatRepository) GetMostListenedMusicians(ctx context.Context, userID uuid.UUID, maxCnt int) ([]domain.UserMusiciansStat, error) {
	return nil, nil
}

func (sr *MongoStatRepository) GetListenedGenres(ctx context.Context, userID uuid.UUID) ([]domain.UserGenresStat, error) {
	return nil, nil
}
