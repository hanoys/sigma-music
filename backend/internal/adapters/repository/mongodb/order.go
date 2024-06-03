package mongodb

import (
	"context"
	entity2 "github.com/hanoys/sigma-music/internal/adapters/repository/mongodb/entity"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
	"github.com/hanoys/sigma-music/internal/util"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoOrderRepository struct {
	db *mongo.Collection
}

func NewMongoOrderRepository(db *mongo.Database) *MongoOrderRepository {
	return &MongoOrderRepository{db: db.Collection(OrderCollection)}
}

func (or *MongoOrderRepository) Create(ctx context.Context, order domain.Order) (domain.Order, error) {
	mongoOrder := entity2.NewMongoOrder(order)
	_, err := or.db.InsertOne(ctx, mongoOrder)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return domain.Order{}, util.WrapError(ports.ErrOrderDuplicate, err)
		}
		return domain.Order{}, util.WrapError(ports.ErrInternalOrderRepo, err)
	}

	return order, nil
}
