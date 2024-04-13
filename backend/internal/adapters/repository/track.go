package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/adapters/repository/entity"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/jmoiron/sqlx"
)

const (
	trackDeleteQuery  = "DELETE FROM tracks WHERE id = $1"
	trackGetByIDQuery = "SELECT * FROM tracks WHERE id = $1"
)

type PostgresTrackRepository struct {
	db *sqlx.DB
}

func NewPostgresTrackRepository(db *sqlx.DB) *PostgresTrackRepository {
	return &PostgresTrackRepository{db: db}
}

func (tr *PostgresTrackRepository) Create(ctx context.Context, track domain.Track) (domain.Track, error) {
	pgTrack := entity.NewPgTrack(track)
	queryString := entity.InsertQueryString(pgTrack, "tracks")
	_, err := tr.db.NamedExecContext(ctx, queryString, pgTrack)
	if err != nil {
		return domain.Track{}, err
	}

	var createdTrack entity.PgTrack
	err = tr.db.GetContext(ctx, &createdTrack, trackGetByIDQuery, pgTrack.ID)
	if err != nil {
		return domain.Track{}, err
	}

	return createdTrack.ToDomain(), nil
}

func (tr *PostgresTrackRepository) Delete(ctx context.Context, trackID uuid.UUID) (domain.Track, error) {
	var deletedTrack entity.PgTrack
	err := tr.db.GetContext(ctx, &deletedTrack, trackGetByIDQuery, trackID)
	if err != nil {
		return domain.Track{}, err
	}

	_, err = tr.db.ExecContext(ctx, trackDeleteQuery, trackID)
	if err != nil {
		return domain.Track{}, err
	}

	return deletedTrack.ToDomain(), nil
}
