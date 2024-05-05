package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/adapters/repository/entity"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
	"github.com/hanoys/sigma-music/internal/util"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

const (
	trackGetAllQuery  = "SELECT * FROM tracks"
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
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return domain.Track{}, util.WrapError(ports.ErrTrackDuplicate, err)
			}
		}
		return domain.Track{}, util.WrapError(ports.ErrInternalTrackRepo, err)
	}

	var createdTrack entity.PgTrack
	err = tr.db.GetContext(ctx, &createdTrack, trackGetByIDQuery, pgTrack.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Track{}, util.WrapError(ports.ErrTrackIDNotFound, err)
		}
		return domain.Track{}, util.WrapError(ports.ErrInternalTrackRepo, err)
	}

	return createdTrack.ToDomain(), nil
}

func (tr *PostgresTrackRepository) GetAll(ctx context.Context) ([]domain.Track, error) {
	var tracks []entity.PgTrack
	err := tr.db.SelectContext(ctx, &tracks, trackGetAllQuery)
	if err != nil {
		return nil, util.WrapError(ports.ErrInternalTrackRepo, err)
	}

	domainTracks := make([]domain.Track, len(tracks))
	for i, track := range tracks {
		domainTracks[i] = track.ToDomain()
	}

	return domainTracks, nil
}

func (tr *PostgresTrackRepository) GetByID(ctx context.Context, trackID uuid.UUID) (domain.Track, error) {
	var foundTrack entity.PgTrack
	err := tr.db.GetContext(ctx, &foundTrack, trackGetByIDQuery, trackID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Track{}, util.WrapError(ports.ErrTrackIDNotFound, err)
		}
		return domain.Track{}, util.WrapError(ports.ErrInternalTrackRepo, err)
	}

	return foundTrack.ToDomain(), nil
}

func (tr *PostgresTrackRepository) Delete(ctx context.Context, trackID uuid.UUID) (domain.Track, error) {
	var deletedTrack entity.PgTrack
	err := tr.db.GetContext(ctx, &deletedTrack, trackGetByIDQuery, trackID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Track{}, util.WrapError(ports.ErrTrackIDNotFound, err)
		}
		return domain.Track{}, util.WrapError(ports.ErrInternalTrackRepo, err)
	}

	_, err = tr.db.ExecContext(ctx, trackDeleteQuery, trackID)
	if err != nil {
		return domain.Track{}, util.WrapError(ports.ErrTrackDelete, err)
	}

	return deletedTrack.ToDomain(), nil
}
