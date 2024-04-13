package repository

import (
	"context"
	"github.com/hanoys/sigma-music/internal/adapters/repository/entity"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/jmoiron/sqlx"
)

const (
	userGetByUniqueQuery = "SELECT * FROM users WHERE $1 = $2"
)

type PostgresUserRepository struct {
	db *sqlx.DB
}

func NewPostgresUserRepository(db *sqlx.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (ur *PostgresUserRepository) Create(ctx context.Context, user domain.User) (domain.User, error) {
	pgUser := entity.NewPgUser(user)
	queryString := entity.InsertQueryString(pgUser, "users")
	_, err := ur.db.NamedExecContext(ctx, queryString, pgUser)
	if err != nil {
		return domain.User{}, err
	}

	var createdUser entity.PgUser
	err = ur.db.GetContext(ctx, &createdUser, userGetByUniqueQuery, "id", pgUser.ID)
	if err != nil {
		return domain.User{}, err
	}

	return createdUser.ToDomain(), nil
}

func (ur *PostgresUserRepository) GetByName(ctx context.Context, name string) (domain.User, error) {
	var foundUser entity.PgUser
	err := ur.db.GetContext(ctx, &foundUser, userGetByUniqueQuery, "name", name)
	if err != nil {
		return domain.User{}, err
	}

	return foundUser.ToDomain(), nil
}

func (ur *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	var foundUser entity.PgUser
	err := ur.db.GetContext(ctx, &foundUser, userGetByUniqueQuery, "email", email)
	if err != nil {
		return domain.User{}, err
	}

	return foundUser.ToDomain(), nil
}

func (ur *PostgresUserRepository) GetByPhone(ctx context.Context, phone string) (domain.User, error) {
	var foundUser entity.PgUser
	err := ur.db.GetContext(ctx, &foundUser, userGetByUniqueQuery, "phone", phone)
	if err != nil {
		return domain.User{}, err
	}

	return foundUser.ToDomain(), nil
}
