package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v5"
	entity2 "github.com/hanoys/sigma-music/internal/adapters/repository/postgres/entity"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
	"github.com/hanoys/sigma-music/internal/util"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

const (
	AlbumGetAllQuery          = "SELECT * FROM albums WHERE published = TRUE"
	AlbumGetByMusicianIDQuery = "SELECT a.id, a.name, a.description, a.published, a.release_date FROM album_musician JOIN public.albums a on a.id = album_musician.album_id WHERE musician_id = $1 AND published = TRUE"
	AlbumGetOwnQuery          = "SELECT a.id, a.name, a.description, a.published, a.release_date FROM album_musician JOIN public.albums a on a.id = album_musician.album_id WHERE musician_id = $1"
	AlbumGetByIDQuery         = "SELECT * FROM albums WHERE id = $1 AND published = TRUE"
	AlbumGetByIDInternalQuery = "SELECT * FROM albums WHERE id = $1"
	AlbumInsertQuery          = "INSERT INTO album_musician(musician_id, album_id) VALUES ($1, $2)"
)

type PostgresAlbumRepository struct {
	connection *sqlx.DB
}

func NewPostgresAlbumRepository(connection *sqlx.DB) *PostgresAlbumRepository {
	return &PostgresAlbumRepository{connection: connection}
}

func (ar *PostgresAlbumRepository) Create(ctx context.Context, album domain.Album, musicianID uuid.UUID) (domain.Album, error) {
	pgAlbum := entity2.NewPgAlbum(album)
	queryString := entity2.InsertQueryString(pgAlbum, "albums")
	_, err := ar.connection.NamedExecContext(ctx, queryString, pgAlbum)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return domain.Album{}, ports.ErrAlbumDuplicate
			}
		}
		return domain.Album{}, util.WrapError(ports.ErrInternalAlbumRepo, err)
	}

	_, err = ar.connection.ExecContext(ctx, AlbumInsertQuery, musicianID, album.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return domain.Album{}, ports.ErrAlbumDuplicate
			}
		}
		return domain.Album{}, util.WrapError(ports.ErrInternalAlbumRepo, err)
	}

	var createdUser entity2.PgAlbum
	err = ar.connection.GetContext(ctx, &createdUser, AlbumGetByIDInternalQuery, pgAlbum.ID)
	if err != nil {
		return domain.Album{}, util.WrapError(ports.ErrAlbumIDNotFound, err)
	}

	return createdUser.ToDomain(), nil
}

func (ar *PostgresAlbumRepository) GetAll(ctx context.Context) ([]domain.Album, error) {
	var albums []entity2.PgAlbum
	err := ar.connection.SelectContext(ctx, &albums, AlbumGetAllQuery)
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
	var albums []entity2.PgAlbum
	err := ar.connection.SelectContext(ctx, &albums, AlbumGetByMusicianIDQuery, musicianID)
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
	var albums []entity2.PgAlbum
	err := ar.connection.SelectContext(ctx, &albums, AlbumGetOwnQuery, musicianID)
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
	var foundAlbum entity2.PgAlbum
	err := ar.connection.GetContext(ctx, &foundAlbum, AlbumGetByIDQuery, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Album{}, util.WrapError(ports.ErrAlbumIDNotFound, err)
		}
		return domain.Album{}, util.WrapError(ports.ErrInternalAlbumRepo, err)
	}

	return foundAlbum.ToDomain(), nil
}

func (ar *PostgresAlbumRepository) Publish(ctx context.Context, id uuid.UUID) error {
	var foundAlbum entity2.PgAlbum
	err := ar.connection.GetContext(ctx, &foundAlbum, AlbumGetByIDInternalQuery, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return util.WrapError(ports.ErrAlbumIDNotFound, err)
		}
		return util.WrapError(ports.ErrInternalAlbumRepo, err)
	}

	foundAlbum.Published = true
	foundAlbum.ReleaseDate = null.TimeFrom(time.Now())
	updateQuery := entity2.UpdateQueryString(foundAlbum, "albums")
	_, err = ar.connection.NamedExecContext(ctx, updateQuery, foundAlbum)
	if err != nil {
		return util.WrapError(ports.ErrAlbumPublish, err)
	}

	return nil
}
