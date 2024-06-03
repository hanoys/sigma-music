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

type MongoCommentRepository struct {
	db *mongo.Collection
}

func NewMongoCommentRepository(db *mongo.Database) *MongoCommentRepository {
	return &MongoCommentRepository{
		db: db.Collection(CommentCollection),
	}
}

func (cr *MongoCommentRepository) Create(ctx context.Context, comment domain.Comment) (domain.Comment, error) {
	session, err := cr.db.Database().Client().StartSession()
	if err != nil {
		return domain.Comment{}, nil
	}

	defer session.EndSession(ctx)

	err = mongo.WithSession(ctx, session, func(sessionContext mongo.SessionContext) error {
		MongoAlbum := entity2.NewMongoComment(comment)
		_, err := cr.db.InsertOne(sessionContext, MongoAlbum)
		if err != nil {
			if mongo.IsDuplicateKeyError(err) {
				return util.WrapError(ports.ErrCommentDuplicate, err)
			}
			return util.WrapError(ports.ErrInternalCommentRepo, err)
		}

		return nil
	})

	return comment, nil
}

func (cr *MongoCommentRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Comment, error) {
	cursor, err := cr.db.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, util.WrapError(ports.ErrAlbumIDNotFound, err)
		}
		return nil, util.WrapError(ports.ErrInternalAlbumRepo, err)
	}

	var mgCommentArray []entity2.MongoComment
	err = cursor.All(ctx, &mgCommentArray)
	if err != nil {
		return nil, err
	}

	comments := make([]domain.Comment, len(mgCommentArray))
	for i, comment := range mgCommentArray {
		comments[i] = comment.ToDomain()
	}

	return comments, nil
}

func (cr *MongoCommentRepository) GetByTrackID(ctx context.Context, trackID uuid.UUID) ([]domain.Comment, error) {
	cursor, err := cr.db.Find(ctx, bson.M{"track_id": trackID})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, util.WrapError(ports.ErrAlbumIDNotFound, err)
		}
		return nil, util.WrapError(ports.ErrInternalAlbumRepo, err)
	}

	var mgCommentArray []entity2.MongoComment
	err = cursor.All(ctx, &mgCommentArray)
	if err != nil {
		return nil, err
	}

	comments := make([]domain.Comment, len(mgCommentArray))
	for i, comment := range mgCommentArray {
		comments[i] = comment.ToDomain()
	}

	return comments, nil
}
