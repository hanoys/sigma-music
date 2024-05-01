package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/adapters/repository/entity"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
	"github.com/hanoys/sigma-music/internal/utill"
	"github.com/jmoiron/sqlx"
)

const (
	genreGetByIDQuery     = "SELECT * FROM genres WHERE id = $1"
	genreGetAllQuery      = "SELECT * FROM genres"
	genreAddForTrackQuery = "INSERT INTO track_genre (track_id, genre_id) VALUES ($1, $2)"
)

type PostgresGenreRepository struct {
	db *sqlx.DB
}

func NewPostgresGenreRepository(db *sqlx.DB) *PostgresGenreRepository {
	return &PostgresGenreRepository{db: db}
}

func (gr *PostgresGenreRepository) GetAll(ctx context.Context) ([]domain.Genre, error) {
	var genres []entity.PgGenre
	err := gr.db.SelectContext(ctx, &genres, genreGetAllQuery)
	if err != nil {
		return nil, utill.WrapError(ports.ErrInternalGenreRepo, err)
	}

	domainGenres := make([]domain.Genre, len(genres))
	for i, genre := range genres {
		domainGenres[i] = genre.ToDomain()
	}

	return domainGenres, nil
}

func (gr *PostgresGenreRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.Genre, error) {
	var foundGenre entity.PgGenre
	err := gr.db.GetContext(ctx, &foundGenre, genreGetByIDQuery, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Genre{}, utill.WrapError(ports.ErrGenreIDNotFound, err)
		}
		return domain.Genre{}, utill.WrapError(ports.ErrInternalGenreRepo, err)
	}

	return foundGenre.ToDomain(), nil
}

func (gr *PostgresGenreRepository) AddForTrack(ctx context.Context, trackID uuid.UUID, genresID []uuid.UUID) error {
	tx := gr.db.MustBegin()
	for _, genreID := range genresID {
		tx.MustExec(genreAddForTrackQuery, trackID, genreID)
	}

	err := tx.Commit()
	if err != nil {
		return utill.WrapError(ports.ErrInternalGenreRepo, err)
	}

	return nil
}
