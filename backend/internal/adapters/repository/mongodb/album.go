package mongodb

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/guregu/null/v5"
	entity2 "github.com/hanoys/sigma-music/internal/adapters/repository/mongodb/entity"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
	"github.com/hanoys/sigma-music/internal/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type MongoAlbumRepository struct {
	db *mongo.Collection
}

func NewMongoAlbumRepository(db *mongo.Database) *MongoAlbumRepository {
	return &MongoAlbumRepository{
		db: db.Collection(AlbumCollection),
	}
}

func (ar *MongoAlbumRepository) Create(ctx context.Context, album domain.Album, musicianID uuid.UUID) (domain.Album, error) {
	session, err := ar.db.Database().Client().StartSession()
	if err != nil {
		return domain.Album{}, ports.ErrInternalAlbumRepo
	}

	defer session.EndSession(ctx)

	err = mongo.WithSession(ctx, session, func(sessionContext mongo.SessionContext) error {
		MongoAlbum := entity2.NewMongoAlbum(album)
		_, err := ar.db.InsertOne(sessionContext, MongoAlbum)
		if err != nil {
			if mongo.IsDuplicateKeyError(err) {
				return util.WrapError(ports.ErrAlbumDuplicate, err)
			}
			return util.WrapError(ports.ErrInternalAlbumRepo, err)
		}

		_, err = ar.db.Database().Collection(AlbumMusicianCollection).InsertOne(sessionContext, struct {
			ID         uuid.UUID `bson:"_id"`
			AlbumID    uuid.UUID `bson:"album_id"`
			MusicianID uuid.UUID `bson:"musician_id"`
		}{
			ID:         uuid.New(),
			AlbumID:    album.ID,
			MusicianID: musicianID,
		})
		if err != nil {
			if mongo.IsDuplicateKeyError(err) {
				return util.WrapError(ports.ErrAlbumDuplicate, err)
			}
			return util.WrapError(ports.ErrInternalAlbumRepo, err)
		}

		return nil
	})

	return ar.GetByID(ctx, album.ID)
}

func (ar *MongoAlbumRepository) GetAll(ctx context.Context) ([]domain.Album, error) {
	cursor, err := ar.db.Find(ctx, bson.M{"published": true})
	if err != nil {
		return nil, util.WrapError(ports.ErrInternalAlbumRepo, err)
	}

	var mgAlbumArray []entity2.MongoAlbum
	err = cursor.All(ctx, &mgAlbumArray)
	if err != nil {
		return nil, err
	}

	albums := make([]domain.Album, len(mgAlbumArray))
	for i, album := range mgAlbumArray {
		albums[i] = album.ToDomain()
	}

	return albums, nil
}

func (ar *MongoAlbumRepository) GetByMusicianID(ctx context.Context, musicianID uuid.UUID) ([]domain.Album, error) {
	cursor, err := ar.db.Database().Collection(AlbumMusicianCollection).Find(ctx, bson.M{"musician_id": musicianID})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, util.WrapError(ports.ErrAlbumIDNotFound, err)
		}
		return nil, util.WrapError(ports.ErrInternalAlbumRepo, err)
	}

	var mgAlbumsMusicians []struct {
		AlbumID    uuid.UUID `bson:"album_id"`
		MusicianID uuid.UUID `bson:"musician_id"`
	}
	err = cursor.All(ctx, &mgAlbumsMusicians)
	if err != nil {
		return nil, err
	}

	var albums []domain.Album
	for _, entry := range mgAlbumsMusicians {
		a, _ := ar.GetByID(ctx, entry.AlbumID)
		albums = append(albums, a)
	}

	return albums, nil
}

func (ar *MongoAlbumRepository) GetOwn(ctx context.Context, musicianID uuid.UUID) ([]domain.Album, error) {
	cursor, err := ar.db.Database().Collection(AlbumMusicianCollection).Find(ctx, bson.M{"musician_id": musicianID})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, util.WrapError(ports.ErrAlbumIDNotFound, err)
		}
		return nil, util.WrapError(ports.ErrInternalAlbumRepo, err)
	}

	var mgAlbumsMusicians []struct {
		AlbumID    uuid.UUID `bson:"album_id"`
		MusicianID uuid.UUID `bson:"musician_id"`
	}
	err = cursor.All(ctx, &mgAlbumsMusicians)
	if err != nil {
		return nil, err
	}

	var albums []domain.Album
	for _, entry := range mgAlbumsMusicians {
		a, _ := ar.GetByID(ctx, entry.AlbumID)
		albums = append(albums, a)
	}

	return albums, nil
}

func (ar *MongoAlbumRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.Album, error) {
	result := ar.db.FindOne(ctx, bson.M{"_id": id})

	var mgAlbum entity2.MongoAlbum
	if err := result.Decode(&mgAlbum); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Album{}, util.WrapError(ports.ErrAlbumIDNotFound, err)
		}
		return domain.Album{}, util.WrapError(ports.ErrInternalAlbumRepo, err)
	}
	return mgAlbum.ToDomain(), nil
}

func (ar *MongoAlbumRepository) Publish(ctx context.Context, id uuid.UUID) error {
	album, err := ar.GetByID(ctx, id)
	if err != nil {
		return err
	}

	album.Published = true
	album.ReleaseDate = null.TimeFrom(time.Now())
	var mongoAlbum = entity2.NewMongoAlbum(album)
	_, err = ar.db.ReplaceOne(ctx, bson.M{"_id": id}, mongoAlbum)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return util.WrapError(ports.ErrAlbumIDNotFound, err)
		}
		return util.WrapError(ports.ErrInternalAlbumRepo, err)
	}

	return nil
}
