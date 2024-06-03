package mongodb

import (
	"context"
	"errors"
	"github.com/google/uuid"
	entity2 "github.com/hanoys/sigma-music/internal/adapters/repository/mongodb/entity"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
	"github.com/hanoys/sigma-music/internal/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoUserRepository struct {
	db *mongo.Collection
}

func NewMongoUserRepository(db *mongo.Database) *MongoUserRepository {
	return &MongoUserRepository{db: db.Collection(UserCollection)}
}

func (ur *MongoUserRepository) Create(ctx context.Context, user domain.User) (domain.User, error) {
	session, err := ur.db.Database().Client().StartSession()
	if err != nil {
		return domain.User{}, ports.ErrInternalUserRepo
	}

	defer session.EndSession(ctx)

	err = mongo.WithSession(ctx, session, func(sessionContext mongo.SessionContext) error {
		mongoUser := entity2.NewMongoUser(user)
		_, err := ur.db.InsertOne(sessionContext, mongoUser)
		if err != nil {
			if mongo.IsDuplicateKeyError(err) {
				return util.WrapError(ports.ErrUserDuplicate, err)
			}
			return util.WrapError(ports.ErrInternalUserRepo, err)
		}

		return nil
	})

	return ur.GetByID(ctx, user.ID)
}

func (ur *MongoUserRepository) GetAll(ctx context.Context) ([]domain.User, error) {
	cursor, err := ur.db.Find(ctx, bson.M{})
	if err != nil {
		return nil, util.WrapError(ports.ErrInternalUserRepo, err)
	}

	var mgUserArray []entity2.MongoUser
	err = cursor.All(ctx, &mgUserArray)
	if err != nil {
		return nil, err
	}

	users := make([]domain.User, len(mgUserArray))
	for i, user := range mgUserArray {
		users[i] = user.ToDomain()
	}

	return users, nil
}

func (ur *MongoUserRepository) GetByID(ctx context.Context, userID uuid.UUID) (domain.User, error) {
	result := ur.db.FindOne(ctx, bson.M{"_id": userID})

	var mgUser entity2.MongoUser
	if err := result.Decode(&mgUser); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.User{}, util.WrapError(ports.ErrUserIDNotFound, err)
		}
		return domain.User{}, util.WrapError(ports.ErrInternalUserRepo, err)
	}
	return mgUser.ToDomain(), nil
}

func (ur *MongoUserRepository) GetByName(ctx context.Context, name string) (domain.User, error) {
	result := ur.db.FindOne(ctx, bson.M{"name": name})

	var mgUser entity2.MongoUser
	if err := result.Decode(&mgUser); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.User{}, util.WrapError(ports.ErrUserIDNotFound, err)
		}
		return domain.User{}, util.WrapError(ports.ErrInternalUserRepo, err)
	}
	return mgUser.ToDomain(), nil
}

func (ur *MongoUserRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	result := ur.db.FindOne(ctx, bson.M{"email": email})

	var mgUser entity2.MongoUser
	if err := result.Decode(&mgUser); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.User{}, util.WrapError(ports.ErrUserIDNotFound, err)
		}
		return domain.User{}, util.WrapError(ports.ErrInternalUserRepo, err)
	}
	return mgUser.ToDomain(), nil
}

func (ur *MongoUserRepository) GetByPhone(ctx context.Context, phone string) (domain.User, error) {
	result := ur.db.FindOne(ctx, bson.M{"phone": phone})

	var mgUser entity2.MongoUser
	if err := result.Decode(&mgUser); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.User{}, util.WrapError(ports.ErrUserIDNotFound, err)
		}
		return domain.User{}, util.WrapError(ports.ErrInternalUserRepo, err)
	}
	return mgUser.ToDomain(), nil
}
