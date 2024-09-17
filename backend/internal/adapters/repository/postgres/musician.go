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
	musicianGetAllQuery       = "SELECT * FROM musicians"
	musicianGetByIDQuery      = "SELECT * FROM musicians WHERE id = $1"
	musicianGetByNameQuery    = "SELECT * FROM musicians WHERE name = $1"
	musicianGetByEmailQuery   = "SELECT * FROM musicians WHERE email = $1"
	musicianGetByAlbumIDQuery = "SELECT m.id, m.name, m.email, m.salt, m.password, m.country, m.description FROM musicians m JOIN public.album_musician am on m.id = am.musician_id WHERE album_id = $1"
	musicianGetByTrackIDQuery = "SELECT m.id, m.name, m.email, m.salt, m.password, m.country, m.description FROM musicians m JOIN public.album_musician am on m.id = am.musician_id JOIN public.tracks t ON am.album_id = t.album_id WHERE t.id = $1"
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
	err = mr.connection.GetContext(ctx, &createdMusician, musicianGetByIDQuery, pgMusician.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Musician{}, util.WrapError(ports.ErrMusicianIDNotFound, err)
		}
		return domain.Musician{}, util.WrapError(ports.ErrInternalMusicianRepo, err)
	}

	return createdMusician.ToDomain(), nil
}

func (mr *PostgresMusicianRepository) GetAll(ctx context.Context) ([]domain.Musician, error) {
	var musicians []entity2.PgMusician
	err := mr.connection.SelectContext(ctx, &musicians, musicianGetAllQuery)
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
	err := mr.connection.GetContext(ctx, &foundMusician, musicianGetByIDQuery, musicianID)
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
	err := mr.connection.GetContext(ctx, &foundMusician, musicianGetByNameQuery, name)
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
	err := mr.connection.GetContext(ctx, &foundMusician, musicianGetByEmailQuery, email)
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
	err := mr.connection.GetContext(ctx, &foundMusician, musicianGetByAlbumIDQuery, albumID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Musician{}, util.WrapError(ports.ErrMusicianEmailNotFound, err)
		}
		return domain.Musician{}, util.WrapError(ports.ErrInternalMusicianRepo, err)
	}

	return foundMusician.ToDomain(), nil
}

func (mr *PostgresMusicianRepository) GetByTrackID(ctx context.Context, trackID uuid.UUID) (domain.Musician, error) {
	var foundMusician entity2.PgMusician
	err := mr.connection.GetContext(ctx, &foundMusician, musicianGetByTrackIDQuery, trackID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Musician{}, util.WrapError(ports.ErrMusicianEmailNotFound, err)
		}
		return domain.Musician{}, util.WrapError(ports.ErrInternalMusicianRepo, err)
	}

	return foundMusician.ToDomain(), nil
}
