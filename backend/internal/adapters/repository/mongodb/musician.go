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

type MongoMusicianRepository struct {
	db *mongo.Collection
}

func NewMongoMusicianRepository(db *mongo.Database) *MongoMusicianRepository {
	return &MongoMusicianRepository{db: db.Collection(MusicianCollection)}
}

func (mr *MongoMusicianRepository) Create(ctx context.Context, musician domain.Musician) (domain.Musician, error) {
	session, err := mr.db.Database().Client().StartSession()
	if err != nil {
		return domain.Musician{}, ports.ErrInternalMusicianRepo
	}

	defer session.EndSession(ctx)

	err = mongo.WithSession(ctx, session, func(sessionContext mongo.SessionContext) error {
		mongoUser := entity2.NewMongoMusician(musician)
		_, err := mr.db.InsertOne(sessionContext, mongoUser)
		if err != nil {
			if mongo.IsDuplicateKeyError(err) {
				return util.WrapError(ports.ErrMusicianDuplicate, err)
			}
			return util.WrapError(ports.ErrInternalMusicianRepo, err)
		}

		return nil
	})

	return mr.GetByID(ctx, musician.ID)
}

func (mr *MongoMusicianRepository) GetAll(ctx context.Context) ([]domain.Musician, error) {
	cursor, err := mr.db.Find(ctx, bson.M{})
	if err != nil {
		return nil, util.WrapError(ports.ErrInternalMusicianRepo, err)
	}

	var mgMusicianArray []entity2.MongoMusician
	err = cursor.All(ctx, &mgMusicianArray)
	if err != nil {
		return nil, err
	}

	musicians := make([]domain.Musician, len(mgMusicianArray))
	for i, musician := range mgMusicianArray {
		musicians[i] = musician.ToDomain()
	}

	return musicians, nil
}

func (mr *MongoMusicianRepository) GetByID(ctx context.Context, musicianID uuid.UUID) (domain.Musician, error) {
	result := mr.db.FindOne(ctx, bson.M{"_id": musicianID})

	var mgMusician entity2.MongoMusician
	if err := result.Decode(&mgMusician); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Musician{}, util.WrapError(ports.ErrMusicianIDNotFound, err)
		}
		return domain.Musician{}, util.WrapError(ports.ErrInternalMusicianRepo, err)
	}
	return mgMusician.ToDomain(), nil
}

func (mr *MongoMusicianRepository) GetByName(ctx context.Context, name string) (domain.Musician, error) {
	result := mr.db.FindOne(ctx, bson.M{"name": name})

	var mgMusician entity2.MongoMusician
	if err := result.Decode(&mgMusician); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Musician{}, util.WrapError(ports.ErrMusicianIDNotFound, err)
		}
		return domain.Musician{}, util.WrapError(ports.ErrInternalMusicianRepo, err)
	}
	return mgMusician.ToDomain(), nil
}

func (mr *MongoMusicianRepository) GetByEmail(ctx context.Context, email string) (domain.Musician, error) {
	result := mr.db.FindOne(ctx, bson.M{"email": email})

	var mgMusician entity2.MongoMusician
	if err := result.Decode(&mgMusician); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Musician{}, util.WrapError(ports.ErrMusicianIDNotFound, err)
		}
		return domain.Musician{}, util.WrapError(ports.ErrInternalMusicianRepo, err)
	}
	return mgMusician.ToDomain(), nil
}

func (mr *MongoMusicianRepository) GetByAlbumID(ctx context.Context, albumID uuid.UUID) (domain.Musician, error) {
	result := mr.db.Database().Collection(AlbumMusicianCollection).FindOne(ctx, bson.M{"album_id": albumID})

	var albumMusician struct {
		AlbumID    uuid.UUID `bson:"album_id"`
		MusicianID uuid.UUID `bson:"musician_id"`
	}
	if err := result.Decode(&albumMusician); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Musician{}, util.WrapError(ports.ErrMusicianIDNotFound, err)
		}
		return domain.Musician{}, util.WrapError(ports.ErrInternalMusicianRepo, err)
	}

	return mr.GetByID(ctx, albumMusician.MusicianID)
}

func (mr *MongoMusicianRepository) GetByTrackID(ctx context.Context, trackID uuid.UUID) (domain.Musician, error) {
	result := mr.db.Database().Collection(TrackCollection).FindOne(ctx, bson.M{"_id": trackID})
	var mgTrack entity2.MongoTrack
	if err := result.Decode(&mgTrack); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Musician{}, util.WrapError(ports.ErrMusicianIDNotFound, err)
		}
		return domain.Musician{}, util.WrapError(ports.ErrInternalMusicianRepo, err)
	}

	return mr.GetByAlbumID(ctx, mgTrack.AlbumID)
}
