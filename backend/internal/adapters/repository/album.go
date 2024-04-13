package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/adapters/repository/entity"
	"github.com/hanoys/sigma-music/internal/domain"
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
		return domain.Album{}, err
	}

	var createdUser entity.PgAlbum
	err = ar.db.GetContext(ctx, &createdUser, albumGetByIDQuery, pgAlbum.ID)
	if err != nil {
		return domain.Album{}, err
	}

	return createdUser.ToDomain(), nil
}

func (ar *PostgresAlbumRepository) GetByMusicianID(ctx context.Context, musicianID uuid.UUID) ([]domain.Album, error) {
	var albums []entity.PgAlbum
	err := ar.db.SelectContext(ctx, &albums, albumGetByMusicianIDQuery, musicianID)
	if err != nil {
		return nil, err
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
		return domain.Album{}, err
	}

	return foundAlbum.ToDomain(), nil
}

func (ar *PostgresAlbumRepository) Publish(ctx context.Context, id uuid.UUID) error {
	var foundAlbum entity.PgAlbum
	err := ar.db.GetContext(ctx, &foundAlbum, albumGetByIDQuery, id)
	if err != nil {
		return err
	}

	foundAlbum.Published = true
	updateQuery := entity.UpdateQueryString(foundAlbum, "albums")
	_, err = ar.db.NamedExecContext(ctx, updateQuery, foundAlbum)
	if err != nil {
		return err
	}

	return nil
}