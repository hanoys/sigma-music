package repository

import (
	"context"
	"github.com/hanoys/sigma-music/internal/adapters/repository/entity"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/jmoiron/sqlx"
)

const (
	musicianGetByUniqueQuery = "SELECT * FROM musicians WHERE $1 = $2"
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
		return domain.Musician{}, err
	}

	var createdUser entity.PgMusician
	err = mr.db.GetContext(ctx, &createdUser, musicianGetByUniqueQuery, "id", pgMusician.ID)
	if err != nil {
		return domain.Musician{}, err
	}

	return createdUser.ToDomain(), nil
}

func (mr *PostgresMusicianRepository) GetByName(ctx context.Context, name string) (domain.Musician, error) {
	var foundMusician entity.PgMusician
	err := mr.db.GetContext(ctx, &foundMusician, musicianGetByUniqueQuery, "name", name)
	if err != nil {
		return domain.Musician{}, err
	}

	return foundMusician.ToDomain(), nil
}

func (mr *PostgresMusicianRepository) GetByEmail(ctx context.Context, email string) (domain.Musician, error) {
	var foundMusician entity.PgMusician
	err := mr.db.GetContext(ctx, &foundMusician, musicianGetByUniqueQuery, "email", email)
	if err != nil {
		return domain.Musician{}, err
	}

	return foundMusician.ToDomain(), nil
}
