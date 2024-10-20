package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/adapters/repository/postgres/entity"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
	"github.com/hanoys/sigma-music/internal/util"
	"github.com/jmoiron/sqlx"
)

const (
	GenreGetByIDQuery        = "SELECT * FROM genres WHERE id = $1"
	GenreGetAllQuery         = "SELECT * FROM genres"
	genreDeleteForTrackQuery = "DELETE * FROM track_genre WHERE track_id=$1"
	genreAddForTrackQuery    = "INSERT INTO track_genre (track_id, genre_id) VALUES ($1, $2)"
	GenreGetByTrack          = "SELECT g.id, g.name FROM genres g JOIN public.track_genre tg on g.id = tg.genre_id WHERE tg.track_id = $1"
)

type PostgresGenreRepository struct {
	connection *sqlx.DB
}

func NewPostgresGenreRepository(connection *sqlx.DB) *PostgresGenreRepository {
	return &PostgresGenreRepository{connection: connection}
}

func (gr *PostgresGenreRepository) GetAll(ctx context.Context) ([]domain.Genre, error) {
	var genres []entity.PgGenre
	err := gr.connection.SelectContext(ctx, &genres, GenreGetAllQuery)
	if err != nil {
		return nil, util.WrapError(ports.ErrInternalGenreRepo, err)
	}

	domainGenres := make([]domain.Genre, len(genres))
	for i, genre := range genres {
		domainGenres[i] = genre.ToDomain()
	}

	return domainGenres, nil
}

func (gr *PostgresGenreRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.Genre, error) {
	var foundGenre entity.PgGenre
	err := gr.connection.GetContext(ctx, &foundGenre, GenreGetByIDQuery, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Genre{}, util.WrapError(ports.ErrGenreIDNotFound, err)
		}
		return domain.Genre{}, util.WrapError(ports.ErrInternalGenreRepo, err)
	}

	return foundGenre.ToDomain(), nil
}

func (gr *PostgresGenreRepository) AddForTrack(ctx context.Context, trackID uuid.UUID, genresID []uuid.UUID) error {
	tx := gr.connection.MustBegin()
    tx.MustExec(genreDeleteForTrackQuery, trackID)

	for _, genreID := range genresID {
		tx.MustExec(genreAddForTrackQuery, trackID, genreID)
	}

	err := tx.Commit()
	if err != nil {
		return util.WrapError(ports.ErrInternalGenreRepo, err)
	}

	return nil
}

func (gr *PostgresGenreRepository) GetByTrackID(ctx context.Context, trackID uuid.UUID) ([]domain.Genre, error) {
	var genres []entity.PgGenre
	err := gr.connection.SelectContext(ctx, &genres, GenreGetByTrack, trackID)
	if err != nil {
		return nil, util.WrapError(ports.ErrInternalGenreRepo, err)
	}

	domainGenres := make([]domain.Genre, len(genres))
	for i, genre := range genres {
		domainGenres[i] = genre.ToDomain()
	}

	return domainGenres, nil
}
