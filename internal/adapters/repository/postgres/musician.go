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
	MusicianGetAllQuery       = "SELECT * FROM musicians"
	MusicianGetByIDQuery      = "SELECT * FROM musicians WHERE id = $1"
	MusicianGetByNameQuery    = "SELECT * FROM musicians WHERE name = $1"
	MusicianGetByEmailQuery   = "SELECT * FROM musicians WHERE email = $1"
	MusicianGetByAlbumIDQuery = "SELECT m.id, m.name, m.email, m.salt, m.password, m.country, m.description, m.image_url FROM musicians m JOIN public.album_musician am on m.id = am.musician_id WHERE album_id = $1"
	MusicianGetByTrackIDQuery = "SELECT m.id, m.name, m.email, m.salt, m.password, m.country, m.description, m.image_url FROM musicians m JOIN public.album_musician am on m.id = am.musician_id JOIN public.tracks t ON am.album_id = t.album_id WHERE t.id = $1"
)

type PostgresMusicianRepository struct {
	connection *sqlx.DB
}

func NewPostgresMusicianRepository(connection *sqlx.DB) *PostgresMusicianRepository {
	return &PostgresMusicianRepository{connection: connection}
}

func (mr *PostgresMusicianRepository) Create(ctx context.Context, musician domain.Musician) (domain.Musician, error) {
	pgMusician := entity2.NewPgMusician(musician)
	queryString := entity2.InsertQueryString(pgMusician, "musicians")
	_, err := mr.connection.NamedExecContext(ctx, queryString, pgMusician)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return domain.Musician{}, util.WrapError(ports.ErrMusicianDuplicate, err)
			}
		}
		return domain.Musician{}, util.WrapError(ports.ErrInternalMusicianRepo, err)
	}

	var createdMusician entity2.PgMusician
	err = mr.connection.GetContext(ctx, &createdMusician, MusicianGetByIDQuery, pgMusician.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Musician{}, util.WrapError(ports.ErrMusicianIDNotFound, err)
		}
		return domain.Musician{}, util.WrapError(ports.ErrInternalMusicianRepo, err)
	}

	return createdMusician.ToDomain(), nil
}

func (mr *PostgresMusicianRepository) Update(ctx context.Context, musician domain.Musician) (domain.Musician, error) {
	pgMusician := entity2.NewPgMusician(musician)
	queryString := entity2.UpdateQueryString(pgMusician, "musicians")
	_, err := mr.connection.NamedExecContext(ctx, queryString, pgMusician)
	if err != nil {
		return domain.Musician{}, util.WrapError(ports.ErrMusicianUpdate, err)
	}

	var updatedMusician entity2.PgMusician
	err = mr.connection.GetContext(ctx, &updatedMusician, MusicianGetByIDQuery, pgMusician.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Musician{}, util.WrapError(ports.ErrMusicianIDNotFound, err)
		}
		return domain.Musician{}, util.WrapError(ports.ErrInternalMusicianRepo, err)
	}

	return updatedMusician.ToDomain(), nil
}

func (mr *PostgresMusicianRepository) GetAll(ctx context.Context) ([]domain.Musician, error) {
	var musicians []entity2.PgMusician
	err := mr.connection.SelectContext(ctx, &musicians, MusicianGetAllQuery)
	if err != nil {
		return nil, util.WrapError(ports.ErrInternalMusicianRepo, err)
	}

	domainMusicians := make([]domain.Musician, len(musicians))
	for i, musician := range musicians {
		domainMusicians[i] = musician.ToDomain()
	}

	return domainMusicians, nil
}

func (mr *PostgresMusicianRepository) GetByID(ctx context.Context, musicianID uuid.UUID) (domain.Musician, error) {
	var foundMusician entity2.PgMusician
	err := mr.connection.GetContext(ctx, &foundMusician, MusicianGetByIDQuery, musicianID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Musician{}, util.WrapError(ports.ErrMusicianIDNotFound, err)
		}
		return domain.Musician{}, util.WrapError(ports.ErrInternalMusicianRepo, err)
	}

	return foundMusician.ToDomain(), nil
}

func (mr *PostgresMusicianRepository) GetByName(ctx context.Context, name string) (domain.Musician, error) {
	var foundMusician entity2.PgMusician
	err := mr.connection.GetContext(ctx, &foundMusician, MusicianGetByNameQuery, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Musician{}, util.WrapError(ports.ErrMusicianNameNotFound, err)
		}
		return domain.Musician{}, util.WrapError(ports.ErrInternalMusicianRepo, err)
	}

	return foundMusician.ToDomain(), nil
}

func (mr *PostgresMusicianRepository) GetByEmail(ctx context.Context, email string) (domain.Musician, error) {
	var foundMusician entity2.PgMusician
	err := mr.connection.GetContext(ctx, &foundMusician, MusicianGetByEmailQuery, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Musician{}, util.WrapError(ports.ErrMusicianEmailNotFound, err)
		}
		return domain.Musician{}, util.WrapError(ports.ErrInternalMusicianRepo, err)
	}

	return foundMusician.ToDomain(), nil
}

func (mr *PostgresMusicianRepository) GetByAlbumID(ctx context.Context, albumID uuid.UUID) (domain.Musician, error) {
	var foundMusician entity2.PgMusician
	err := mr.connection.GetContext(ctx, &foundMusician, MusicianGetByAlbumIDQuery, albumID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Musician{}, util.WrapError(ports.ErrMusicianIDNotFound, err)
		}
		return domain.Musician{}, util.WrapError(ports.ErrInternalMusicianRepo, err)
	}

	return foundMusician.ToDomain(), nil
}

func (mr *PostgresMusicianRepository) GetByTrackID(ctx context.Context, trackID uuid.UUID) (domain.Musician, error) {
	var foundMusician entity2.PgMusician
	err := mr.connection.GetContext(ctx, &foundMusician, MusicianGetByTrackIDQuery, trackID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Musician{}, util.WrapError(ports.ErrMusicianIDNotFound, err)
		}
		return domain.Musician{}, util.WrapError(ports.ErrInternalMusicianRepo, err)
	}

	return foundMusician.ToDomain(), nil
}
