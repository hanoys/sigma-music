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

type MongoTrackRepository struct {
	db *mongo.Collection
}

func NewMongoTrackRepository(db *mongo.Database) *MongoTrackRepository {
	return &MongoTrackRepository{db: db.Collection(TrackCollection)}
}

func (tr *MongoTrackRepository) Create(ctx context.Context, track domain.Track) (domain.Track, error) {
	session, err := tr.db.Database().Client().StartSession()
	if err != nil {
		return domain.Track{}, ports.ErrInternalTrackRepo
	}

	defer session.EndSession(ctx)

	err = mongo.WithSession(ctx, session, func(sessionContext mongo.SessionContext) error {
		mongoTrack := entity2.NewMongoTrack(track)
		_, err := tr.db.InsertOne(sessionContext, mongoTrack)
		if err != nil {
			if mongo.IsDuplicateKeyError(err) {
				return util.WrapError(ports.ErrTrackDuplicate, err)
			}
			return util.WrapError(ports.ErrInternalTrackRepo, err)
		}

		return nil
	})

	return tr.GetByID(ctx, track.ID)
}

func (tr *MongoTrackRepository) GetAll(ctx context.Context) ([]domain.Track, error) {
	cursor, err := tr.db.Find(ctx, bson.M{})
	if err != nil {
		return nil, util.WrapError(ports.ErrInternalTrackRepo, err)
	}

	var mgTrackArray []entity2.MongoTrack
	err = cursor.All(ctx, &mgTrackArray)
	if err != nil {
		return nil, err
	}

	tracks := make([]domain.Track, len(mgTrackArray))
	for i, track := range mgTrackArray {
		tracks[i] = track.ToDomain()
	}

	return tracks, nil
}

func (tr *MongoTrackRepository) GetByID(ctx context.Context, trackID uuid.UUID) (domain.Track, error) {
	result := tr.db.FindOne(ctx, bson.M{"_id": trackID})

	var mgTrack entity2.MongoTrack
	if err := result.Decode(&mgTrack); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Track{}, util.WrapError(ports.ErrTrackIDNotFound, err)
		}
		return domain.Track{}, util.WrapError(ports.ErrInternalTrackRepo, err)
	}
	return mgTrack.ToDomain(), nil
}

func (tr *MongoTrackRepository) Delete(ctx context.Context, trackID uuid.UUID) (domain.Track, error) {
	track, err := tr.GetByID(ctx, trackID)
	if err != nil {
		return track, err
	}

	_, err = tr.db.DeleteOne(ctx, bson.M{"_id": trackID})
	if err != nil {
		return domain.Track{}, util.WrapError(ports.ErrTrackDelete, err)
	}

	return track, nil
}

func (tr *MongoTrackRepository) GetUserFavorites(ctx context.Context, userID uuid.UUID) ([]domain.Track, error) {
	cursor, err := tr.db.Database().Collection(Favorite).Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, util.WrapError(ports.ErrTrackIDNotFound, err)
		}
		return nil, util.WrapError(ports.ErrInternalTrackRepo, err)
	}

	var favorites []struct {
		TrackID    uuid.UUID `bson:"track_id"`
		MusicianID uuid.UUID `bson:"user_id"`
	}
	err = cursor.All(ctx, &favorites)
	if err != nil {
		return nil, err
	}

	var tracks []domain.Track
	for _, entry := range favorites {
		t, _ := tr.GetByID(ctx, entry.TrackID)
		tracks = append(tracks, t)
	}

	return tracks, nil
}

func (tr *MongoTrackRepository) AddToUserFavorites(ctx context.Context, trackID uuid.UUID, userID uuid.UUID) error {
	_, err := tr.db.Database().Collection(Favorite).InsertOne(ctx, struct {
		TrackID uuid.UUID `bson:"track_id"`
		UserID  uuid.UUID `bson:"user_id"`
	}{
		TrackID: trackID,
		UserID:  userID,
	})

	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return util.WrapError(ports.ErrTrackDuplicate, err)
		}
		return util.WrapError(ports.ErrInternalTrackRepo, err)
	}

	return nil
}

func (tr *MongoTrackRepository) GetByAlbumID(ctx context.Context, albumID uuid.UUID) ([]domain.Track, error) {
	cursor, err := tr.db.Find(ctx, bson.M{"album_id": albumID})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, util.WrapError(ports.ErrTrackIDNotFound, err)
		}
		return nil, util.WrapError(ports.ErrInternalTrackRepo, err)
	}

	var mgTracks []entity2.MongoTrack
	err = cursor.All(ctx, &mgTracks)
	if err != nil {
		return nil, err
	}

	var tracks = make([]domain.Track, len(mgTracks))
	for i, entry := range mgTracks {
		tracks[i] = entry.ToDomain()
	}

	return tracks, nil
}

func (tr *MongoTrackRepository) GetByMusicianID(ctx context.Context, musicianID uuid.UUID) ([]domain.Track, error) {
	cursor, err := tr.db.Database().Collection(AlbumMusicianCollection).Find(ctx, bson.M{"musician_id": musicianID})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, util.WrapError(ports.ErrTrackIDNotFound, err)
		}
		return nil, util.WrapError(ports.ErrInternalTrackRepo, err)
	}

	var mgAlbumsMusicians []struct {
		AlbumID    uuid.UUID `bson:"album_id"`
		MusicianID uuid.UUID `bson:"musician_id"`
	}
	err = cursor.All(ctx, &mgAlbumsMusicians)
	if err != nil {
		return nil, err
	}

	var tracks []domain.Track
	for _, entry := range mgAlbumsMusicians {
		a, _ := tr.GetByAlbumID(ctx, entry.AlbumID)
		for _, track := range a {
			tracks = append(tracks, track)
		}
	}

	return tracks, nil
}

func (tr *MongoTrackRepository) GetOwn(ctx context.Context, musicianID uuid.UUID) ([]domain.Track, error) {
	cursor, err := tr.db.Database().Collection(AlbumMusicianCollection).Find(ctx, bson.M{"musician_id": musicianID})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, util.WrapError(ports.ErrTrackIDNotFound, err)
		}
		return nil, util.WrapError(ports.ErrInternalTrackRepo, err)
	}

	var mgAlbumsMusicians []struct {
		AlbumID    uuid.UUID `bson:"album_id"`
		MusicianID uuid.UUID `bson:"musician_id"`
	}
	err = cursor.All(ctx, &mgAlbumsMusicians)
	if err != nil {
		return nil, err
	}

	var tracks []domain.Track
	for _, entry := range mgAlbumsMusicians {
		a, _ := tr.GetByAlbumID(ctx, entry.AlbumID)
		for _, track := range a {
			tracks = append(tracks, track)
		}
	}

	return tracks, nil
}
