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
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

const (
	albumGetByMusicianIDQuery = "SELECT (a.id, a.name, a.description, a.published, a.release_date) FROM album_musician JOIN public.albums a on a.id = album_musician.album_id WHERE musician_id = $1"
	albumGetByIDQuery         = "SELECT * FROM albums WHERE id = $1"
)

type PostgresAlbumRepository struct {
	db *sqlx.DB
}

func NewPostgresAlbumRepository(db *sqlx.DB) *PostgresAlbumRepository {
	return &PostgresAlbumRepository{db: db}
}

func (ar *PostgresAlbumRepository) Create(ctx context.Context, album domain.Album) (domain.Album, error) {
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
		return domain.Album{}, utill.WrapError(ports.ErrInternalAlbumRepo, err)
	}

	var createdUser entity.PgAlbum
	err = ar.db.GetContext(ctx, &createdUser, albumGetByIDQuery, pgAlbum.ID)
	if err != nil {
		return domain.Album{}, utill.WrapError(ports.ErrAlbumIDNotFound, err)
	}

	return createdUser.ToDomain(), nil
}

func (ar *PostgresAlbumRepository) GetByMusicianID(ctx context.Context, musicianID uuid.UUID) ([]domain.Album, error) {
	var albums []entity.PgAlbum
	err := ar.db.SelectContext(ctx, &albums, albumGetByMusicianIDQuery, musicianID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utill.WrapError(ports.ErrAlbumByMusicianIDNotFound, err)
		}
		return nil, utill.WrapError(ports.ErrInternalAlbumRepo, err)
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
			return domain.Album{}, utill.WrapError(ports.ErrAlbumIDNotFound, err)
		}
		return domain.Album{}, utill.WrapError(ports.ErrInternalAlbumRepo, err)
	}

	return foundAlbum.ToDomain(), nil
}

func (ar *PostgresAlbumRepository) Publish(ctx context.Context, id uuid.UUID) error {
	var foundAlbum entity.PgAlbum
	err := ar.db.GetContext(ctx, &foundAlbum, albumGetByIDQuery, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utill.WrapError(ports.ErrAlbumIDNotFound, err)
		}
		return utill.WrapError(ports.ErrInternalAlbumRepo, err)
	}

	foundAlbum.Published = true
	updateQuery := entity.UpdateQueryString(foundAlbum, "albums")
	_, err = ar.db.NamedExecContext(ctx, updateQuery, foundAlbum)
	if err != nil {
		return utill.WrapError(ports.ErrAlbumPublish, err)
	}

	return nil
}
