package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/guregu/null/v5"
	"github.com/hanoys/sigma-music/internal/adapters/repository/entity"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
	"github.com/hanoys/sigma-music/internal/util"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
	"time"
)

const (
	albumGetAllQuery          = "SELECT * FROM albums WHERE published = TRUE"
	albumGetByMusicianIDQuery = "SELECT a.id, a.name, a.description, a.published, a.release_date FROM album_musician JOIN public.albums a on a.id = album_musician.album_id WHERE musician_id = $1 AND published = TRUE"
	albumGetOwnQuery          = "SELECT a.id, a.name, a.description, a.published, a.release_date FROM album_musician JOIN public.albums a on a.id = album_musician.album_id WHERE musician_id = $1"
	albumGetByIDQuery         = "SELECT * FROM albums WHERE id = $1 AND published = TRUE"
	albumGetByIDInternalQuery = "SELECT * FROM albums WHERE id = $1"
	albumInsertQuery          = "INSERT INTO album_musician(musician_id, album_id) VALUES ($1, $2)"
)

type PostgresAlbumRepository struct {
	db *sqlx.DB
}

func NewPostgresAlbumRepository(db *sqlx.DB) *PostgresAlbumRepository {
	return &PostgresAlbumRepository{db: db}
}

func (ar *PostgresAlbumRepository) Create(ctx context.Context, album domain.Album, musicianID uuid.UUID) (domain.Album, error) {
	pgAlbum := entity.NewPgAlbum(album)
	queryString := entity.InsertQueryString(pgAlbum, "albums")
	_, err := ar.db.NamedExecContext(ctx, queryString, pgAlbum)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return domain.Album{}, ports.ErrAlbumDuplicate
			}
		}
		return domain.Album{}, util.WrapError(ports.ErrInternalAlbumRepo, err)
	}

	_, err = ar.db.ExecContext(ctx, albumInsertQuery, musicianID, album.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return domain.Album{}, ports.ErrAlbumDuplicate
			}
		}
		return domain.Album{}, util.WrapError(ports.ErrInternalAlbumRepo, err)
	}

	var createdUser entity.PgAlbum
	err = ar.db.GetContext(ctx, &createdUser, albumGetByIDInternalQuery, pgAlbum.ID)
	if err != nil {
		return domain.Album{}, util.WrapError(ports.ErrAlbumIDNotFound, err)
	}

	return createdUser.ToDomain(), nil
}

func (ar *PostgresAlbumRepository) GetAll(ctx context.Context) ([]domain.Album, error) {
	var albums []entity.PgAlbum
	err := ar.db.SelectContext(ctx, &albums, albumGetAllQuery)
	if err != nil {
		return nil, util.WrapError(ports.ErrInternalAlbumRepo, err)
	}

	domainAlbums := make([]domain.Album, len(albums))
	for i, track := range albums {
		domainAlbums[i] = track.ToDomain()
	}

	return domainAlbums, nil
}

func (ar *PostgresAlbumRepository) GetByMusicianID(ctx context.Context, musicianID uuid.UUID) ([]domain.Album, error) {
	var albums []entity.PgAlbum
	err := ar.db.SelectContext(ctx, &albums, albumGetByMusicianIDQuery, musicianID)
	if err != nil {
		return nil, util.WrapError(ports.ErrInternalAlbumRepo, err)
	}

	domainAlbums := make([]domain.Album, len(albums))
	for i, album := range albums {
		domainAlbums[i] = album.ToDomain()
	}

	return domainAlbums, nil
}

func (ar *PostgresAlbumRepository) GetOwn(ctx context.Context, musicianID uuid.UUID) ([]domain.Album, error) {
	var albums []entity.PgAlbum
	err := ar.db.SelectContext(ctx, &albums, albumGetOwnQuery, musicianID)
	if err != nil {
		return nil, util.WrapError(ports.ErrInternalAlbumRepo, err)
	}

	domainAlbums := make([]domain.Album, len(albums))
	for i, album := range albums {
		domainAlbums[i] = album.ToDomain()
	}

	return domainAlbums, nil
}

func (ar *PostgresAlbumRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.Album, error) {
	var foundAlbum entity.PgAlbum
	err := ar.db.GetContext(ctx, &foundAlbum, albumGetByIDQuery, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Album{}, util.WrapError(ports.ErrAlbumIDNotFound, err)
		}
		return domain.Album{}, util.WrapError(ports.ErrInternalAlbumRepo, err)
	}

	return foundAlbum.ToDomain(), nil
}

func (ar *PostgresAlbumRepository) Publish(ctx context.Context, id uuid.UUID) error {
	var foundAlbum entity.PgAlbum
	err := ar.db.GetContext(ctx, &foundAlbum, albumGetByIDInternalQuery, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return util.WrapError(ports.ErrAlbumIDNotFound, err)
		}
		return util.WrapError(ports.ErrInternalAlbumRepo, err)
	}

	foundAlbum.Published = true
	foundAlbum.ReleaseDate = null.TimeFrom(time.Now())
	updateQuery := entity.UpdateQueryString(foundAlbum, "albums")
	_, err = ar.db.NamedExecContext(ctx, updateQuery, foundAlbum)
	if err != nil {
		return util.WrapError(ports.ErrAlbumPublish, err)
	}

	return nil
}
