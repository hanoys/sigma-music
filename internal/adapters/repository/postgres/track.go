package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	entity2 "github.com/hanoys/sigma-music/internal/adapters/repository/postgres/entity"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
	"github.com/hanoys/sigma-music/internal/util"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

const (
	TrackGetAllQuery          = "SELECT t.id, t.album_id, t.name, t.url, t.image_url FROM tracks t JOIN albums a ON t.album_id = a.id WHERE a.published = TRUE"
	TrackDeleteQuery          = "DELETE FROM tracks WHERE id = $1"
	TrackGetByIDQuery         = "SELECT t.id, t.album_id, t.name, t.url, t.image_url FROM tracks t JOIN albums a ON t.album_id = a.id WHERE t.id = $1 AND a.published = TRUE"
	TrackGetByIDInternalQuery = "SELECT id, album_id, name, url, image_url FROM tracks WHERE id = $1"
	TrackGetUserFavorites     = "SELECT t.id, t.album_id, t.name, t.url, t.image_url FROM tracks t JOIN favorite f on t.id = f.track_id WHERE f.user_id = $1"
	TrackGetByAlbumID         = "SELECT t.id, t.album_id, t.name, t.url, t.image_url FROM tracks t JOIN albums a ON t.album_id = a.id WHERE a.published = TRUE AND a.id = $1"
	TrackGetByMusicianID      = "SELECT t.id, t.album_id, t.name, t.url, t.image_url FROM tracks t JOIN albums a ON t.album_id = a.id JOIN album_musician am on a.id = am.album_id JOIN musicians m on am.musician_id = m.id WHERE published = TRUE and m.id = $1"
	TrackGetOwn               = "SELECT t.id, t.album_id, t.name, t.url, t.image_url FROM tracks t JOIN albums a ON t.album_id = a.id JOIN album_musician am on a.id = am.album_id JOIN musicians m on am.musician_id = m.id WHERE m.id = $1"
	TrackInsertFavorite       = "INSERT INTO favorite(user_id, track_id) VALUES ($1, $2)"
)

type PostgresTrackRepository struct {
	connection *sqlx.DB
}

func NewPostgresTrackRepository(connection *sqlx.DB) *PostgresTrackRepository {
	return &PostgresTrackRepository{connection: connection}
}

func (tr *PostgresTrackRepository) Create(ctx context.Context, track domain.Track) (domain.Track, error) {
	pgTrack := entity2.NewPgTrack(track)
	queryString := entity2.InsertQueryString(pgTrack, "tracks")
	_, err := tr.connection.NamedExecContext(ctx, queryString, pgTrack)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return domain.Track{}, util.WrapError(ports.ErrTrackDuplicate, err)
			}
		}
		return domain.Track{}, util.WrapError(ports.ErrInternalTrackRepo, err)
	}

	var createdTrack entity2.PgTrack
	err = tr.connection.GetContext(ctx, &createdTrack, TrackGetByIDInternalQuery, pgTrack.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Track{}, util.WrapError(ports.ErrTrackIDNotFound, err)
		}
		return domain.Track{}, util.WrapError(ports.ErrInternalTrackRepo, err)
	}

	return createdTrack.ToDomain(), nil
}

func (tr *PostgresTrackRepository) Update(ctx context.Context, track domain.Track) (domain.Track, error) {
	pgTrack := entity2.NewPgTrack(track)
	queryString := entity2.UpdateQueryString(pgTrack, "tracks")
	_, err := tr.connection.NamedExecContext(ctx, queryString, pgTrack)
	if err != nil {
		return domain.Track{}, util.WrapError(ports.ErrTrackUpdate, err)
	}

	var updatedTrack entity2.PgTrack
	err = tr.connection.GetContext(ctx, &updatedTrack, TrackGetByIDInternalQuery, pgTrack.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Track{}, util.WrapError(ports.ErrTrackIDNotFound, err)
		}
		return domain.Track{}, util.WrapError(ports.ErrInternalTrackRepo, err)
	}

	return updatedTrack.ToDomain(), nil
}

func (tr *PostgresTrackRepository) GetAll(ctx context.Context) ([]domain.Track, error) {
	var tracks []entity2.PgTrack
	err := tr.connection.SelectContext(ctx, &tracks, TrackGetAllQuery)
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
	var foundTrack entity2.PgTrack
	err := tr.connection.GetContext(ctx, &foundTrack, TrackGetByIDQuery, trackID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Track{}, util.WrapError(ports.ErrTrackIDNotFound, err)
		}
		return domain.Track{}, util.WrapError(ports.ErrInternalTrackRepo, err)
	}

	return foundTrack.ToDomain(), nil
}

func (tr *PostgresTrackRepository) Delete(ctx context.Context, trackID uuid.UUID) (domain.Track, error) {
	var deletedTrack entity2.PgTrack
	err := tr.connection.GetContext(ctx, &deletedTrack, TrackGetByIDInternalQuery, trackID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Track{}, util.WrapError(ports.ErrTrackIDNotFound, err)
		}
		return domain.Track{}, util.WrapError(ports.ErrInternalTrackRepo, err)
	}

	_, err = tr.connection.ExecContext(ctx, TrackDeleteQuery, trackID)
	if err != nil {
		return domain.Track{}, util.WrapError(ports.ErrTrackDelete, err)
	}

	return deletedTrack.ToDomain(), nil
}

func (tr *PostgresTrackRepository) GetUserFavorites(ctx context.Context, userID uuid.UUID) ([]domain.Track, error) {
	var tracks []entity2.PgTrack
	err := tr.connection.SelectContext(ctx, &tracks, TrackGetUserFavorites, userID)
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
	_, err := tr.connection.ExecContext(ctx, TrackInsertFavorite, userID, trackID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return util.WrapError(ports.ErrTrackDuplicate, err)
			case pgerrcode.ForeignKeyViolation:
				return util.WrapError(ports.ErrTrackIDNotFound, err)
			}
		}
		return util.WrapError(ports.ErrInternalTrackRepo, err)
	}

	return nil
}

func (tr *PostgresTrackRepository) GetByAlbumID(ctx context.Context, albumID uuid.UUID) ([]domain.Track, error) {
	var tracks []entity2.PgTrack
	err := tr.connection.SelectContext(ctx, &tracks, TrackGetByAlbumID, albumID)
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
	var tracks []entity2.PgTrack
	err := tr.connection.SelectContext(ctx, &tracks, TrackGetByMusicianID, musicianID)
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
	var tracks []entity2.PgTrack
	err := tr.connection.SelectContext(ctx, &tracks, TrackGetOwn, musicianID)
	if err != nil {
		return nil, util.WrapError(ports.ErrInternalTrackRepo, err)
	}

	domainTracks := make([]domain.Track, len(tracks))
	for i, track := range tracks {
		domainTracks[i] = track.ToDomain()
	}

	return domainTracks, nil
}
