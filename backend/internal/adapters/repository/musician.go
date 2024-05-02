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
	musicianGetByIDQuery    = "SELECT * FROM musicians WHERE id = $1"
	musicianGetByNameQuery  = "SELECT * FROM musicians WHERE name = $1"
	musicianGetByEmailQuery = "SELECT * FROM musicians WHERE email = $1"
)

type PostgresMusicianRepository struct {
	db *sqlx.DB
}

func NewPostgresMusicianRepository(db *sqlx.DB) *PostgresMusicianRepository {
	return &PostgresMusicianRepository{db: db}
}

func (mr *PostgresMusicianRepository) Create(ctx context.Context, musician domain.Musician) (domain.Musician, error) {
	pgMusician := entity.NewPgMusician(musician)
	queryString := entity.InsertQueryString(pgMusician, "musicians")
	_, err := mr.db.NamedExecContext(ctx, queryString, pgMusician)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return domain.Musician{}, util.WrapError(ports.ErrMusicianDuplicate, err)
			}
		}
		return domain.Musician{}, util.WrapError(ports.ErrInternalMusicianRepo, err)
	}

	var createdUser entity.PgMusician
	err = mr.db.GetContext(ctx, &createdUser, musicianGetByIDQuery, pgMusician.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Musician{}, util.WrapError(ports.ErrMusicianIDNotFound, err)
		}
		return domain.Musician{}, util.WrapError(ports.ErrInternalMusicianRepo, err)
	}

	return createdUser.ToDomain(), nil
}

func (mr *PostgresMusicianRepository) GetByID(ctx context.Context, musicianID uuid.UUID) (domain.Musician, error) {
	var foundMusician entity.PgMusician
	err := mr.db.GetContext(ctx, &foundMusician, musicianGetByIDQuery, musicianID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Musician{}, util.WrapError(ports.ErrMusicianIDNotFound, err)
		}
		return domain.Musician{}, util.WrapError(ports.ErrInternalMusicianRepo, err)
	}

	return foundMusician.ToDomain(), nil
}

func (mr *PostgresMusicianRepository) GetByName(ctx context.Context, name string) (domain.Musician, error) {
	var foundMusician entity.PgMusician
	err := mr.db.GetContext(ctx, &foundMusician, musicianGetByNameQuery, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Musician{}, util.WrapError(ports.ErrMusicianNameNotFound, err)
		}
		return domain.Musician{}, util.WrapError(ports.ErrInternalMusicianRepo, err)
	}

	return foundMusician.ToDomain(), nil
}

func (mr *PostgresMusicianRepository) GetByEmail(ctx context.Context, email string) (domain.Musician, error) {
	var foundMusician entity.PgMusician
	err := mr.db.GetContext(ctx, &foundMusician, musicianGetByEmailQuery, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Musician{}, util.WrapError(ports.ErrMusicianEmailNotFound, err)
		}
		return domain.Musician{}, util.WrapError(ports.ErrInternalMusicianRepo, err)
	}

	return foundMusician.ToDomain(), nil
}
