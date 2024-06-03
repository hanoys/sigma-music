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

type MongoGenreRepository struct {
	db *mongo.Collection
}

func NewMongoGenreRepository(db *mongo.Database) *MongoGenreRepository {
	r := &MongoGenreRepository{
		db: db.Collection(GenreCollection),
	}
	r.Fill()

	return r
}

func (gr *MongoGenreRepository) Fill() {
	for _, genreName := range genresNames {
		mongoGenre := entity2.NewMongoGenre(domain.Genre{
			ID:   uuid.New(),
			Name: genreName,
		})
		gr.db.InsertOne(context.Background(), mongoGenre)
	}
}

func (gr *MongoGenreRepository) GetAll(ctx context.Context) ([]domain.Genre, error) {
	cursor, err := gr.db.Find(ctx, bson.M{})
	if err != nil {
		return nil, util.WrapError(ports.ErrInternalGenreRepo, err)
	}

	var mgGenreArray []entity2.MongoGenre
	err = cursor.All(ctx, &mgGenreArray)
	if err != nil {
		return nil, err
	}

	genres := make([]domain.Genre, len(mgGenreArray))
	for i, genre := range mgGenreArray {
		genres[i] = genre.ToDomain()
	}

	return genres, nil
}

func (gr *MongoGenreRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.Genre, error) {
	result := gr.db.FindOne(ctx, bson.M{"_id": id})

	var mgGenre entity2.MongoGenre
	if err := result.Decode(&mgGenre); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Genre{}, util.WrapError(ports.ErrAlbumIDNotFound, err)
		}
		return domain.Genre{}, util.WrapError(ports.ErrInternalAlbumRepo, err)
	}
	return mgGenre.ToDomain(), nil
}

func (gr *MongoGenreRepository) AddForTrack(ctx context.Context, trackID uuid.UUID, genresID []uuid.UUID) error {
	session, err := gr.db.Database().Client().StartSession()
	if err != nil {
		return ports.ErrInternalGenreRepo
	}

	defer session.EndSession(ctx)

	err = mongo.WithSession(ctx, session, func(sessionContext mongo.SessionContext) error {
		for _, genreID := range genresID {
			_, err := gr.db.Database().Collection(TrackGenreCollection).InsertOne(sessionContext, struct {
				TrackID uuid.UUID `bson:"track_id"`
				GenreID uuid.UUID `bson:"genre_id"`
			}{
				TrackID: trackID,
				GenreID: genreID,
			})

			if err != nil {
				return util.WrapError(ports.ErrInternalAlbumRepo, err)
			}
		}

		return nil
	})

	return err
}

func (gr *MongoGenreRepository) GetByTrackID(ctx context.Context, trackID uuid.UUID) ([]domain.Genre, error) {
	cursor, err := gr.db.Database().Collection(TrackGenreCollection).Find(ctx, bson.M{"track_id": trackID})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, util.WrapError(ports.ErrGenreIDNotFound, err)
		}
		return nil, util.WrapError(ports.ErrInternalGenreRepo, err)
	}

	var mgTrackGenres []struct {
		TrackID uuid.UUID `bson:"track_id"`
		GenreID uuid.UUID `bson:"genre_id"`
	}
	err = cursor.All(ctx, &mgTrackGenres)
	if err != nil {
		return nil, err
	}

	var genres []domain.Genre
	for _, entry := range mgTrackGenres {
		a, _ := gr.GetByID(ctx, entry.GenreID)
		genres = append(genres, a)
	}

	return genres, nil
}
