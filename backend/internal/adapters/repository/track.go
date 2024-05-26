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
	trackGetAllQuery          = "SELECT t.id, t.album_id, t.name, t.url FROM tracks t JOIN albums a ON t.album_id = a.id WHERE a.published = TRUE"
	trackDeleteQuery          = "DELETE FROM tracks WHERE id = $1"
	trackGetByIDQuery         = "SELECT t.id, t.album_id, t.name, t.url FROM tracks t JOIN albums a ON t.album_id = a.id WHERE t.id = $1 AND a.published = TRUE"
	trackGetByIDInternalQuery = "SELECT id, album_id, name, url FROM tracks WHERE id = $1"
	trackGetUserFavorites     = "SELECT t.id, t.album_id, t.name, t.url FROM tracks t JOIN favorite f on t.id = f.track_id WHERE f.user_id = $1"
	trackGetByAlbumID         = "SELECT t.id, t.album_id, t.name, t.url FROM tracks t JOIN albums a ON t.album_id = a.id WHERE a.published = TRUE AND a.id = $1"
	trackGetByMusicianID      = "SELECT t.id, t.album_id, t.name, t.url FROM tracks t JOIN albums a ON t.album_id = a.id JOIN album_musician am on a.id = am.album_id JOIN musicians m on am.musician_id = m.id WHERE published = TRUE and m.id = $1"
	trackGetOwn               = "SELECT t.id, t.album_id, t.name, t.url FROM tracks t JOIN albums a ON t.album_id = a.id JOIN album_musician am on a.id = am.album_id JOIN musicians m on am.musician_id = m.id WHERE m.id = $1"
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
	err = tr.db.GetContext(ctx, &createdTrack, trackGetByIDInternalQuery, pgTrack.ID)
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
	err := tr.db.GetContext(ctx, &deletedTrack, trackGetByIDInternalQuery, trackID)
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

func (tr *PostgresTrackRepository) GetUserFavorites(ctx context.Context, userID uuid.UUID) ([]domain.Track, error) {
	var tracks []entity.PgTrack
	err := tr.db.SelectContext(ctx, &tracks, trackGetUserFavorites, userID)
	if err != nil {
		return nil, util.WrapError(ports.ErrInternalTrackRepo, err)
	}

	domainTracks := make([]domain.Track, len(tracks))
	for i, track := range tracks {
		domainTracks[i] = track.ToDomain()
	}

	return domainTracks, nil
}

func (tr *PostgresTrackRepository) AddToUserFavorites(ctx context.Context, trackID uuid.UUID, userID uuid.UUID) error {
	_, err := tr.db.ExecContext(ctx, "INSERT INTO favorite(user_id, track_id) VALUES ($1, $2)", userID, trackID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return util.WrapError(ports.ErrTrackDuplicate, err)
			}
		}
		return util.WrapError(ports.ErrInternalTrackRepo, err)
	}

	return nil
}

func (tr *PostgresTrackRepository) GetByAlbumID(ctx context.Context, albumID uuid.UUID) ([]domain.Track, error) {
	var tracks []entity.PgTrack
	err := tr.db.SelectContext(ctx, &tracks, trackGetByAlbumID, albumID)
	if err != nil {
		return nil, util.WrapError(ports.ErrInternalTrackRepo, err)
	}

	domainTracks := make([]domain.Track, len(tracks))
	for i, track := range tracks {
		domainTracks[i] = track.ToDomain()
	}

	return domainTracks, nil
}

func (tr *PostgresTrackRepository) GetByMusicianID(ctx context.Context, musicianID uuid.UUID) ([]domain.Track, error) {
	var tracks []entity.PgTrack
	err := tr.db.SelectContext(ctx, &tracks, trackGetByMusicianID, musicianID)
	if err != nil {
		return nil, util.WrapError(ports.ErrInternalTrackRepo, err)
	}

	domainTracks := make([]domain.Track, len(tracks))
	for i, track := range tracks {
		domainTracks[i] = track.ToDomain()
	}

	return domainTracks, nil
}

func (tr *PostgresTrackRepository) GetOwn(ctx context.Context, musicianID uuid.UUID) ([]domain.Track, error) {
	var tracks []entity.PgTrack
	err := tr.db.SelectContext(ctx, &tracks, trackGetOwn, musicianID)
	if err != nil {
		return nil, util.WrapError(ports.ErrInternalTrackRepo, err)
	}

	domainTracks := make([]domain.Track, len(tracks))
	for i, track := range tracks {
		domainTracks[i] = track.ToDomain()
	}

	return domainTracks, nil
}
