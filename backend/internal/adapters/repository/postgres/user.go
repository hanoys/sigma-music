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
	UserGetAllQuery     = "SELECT * FROM users"
	UserGetByIDQuery    = "SELECT * FROM users WHERE id = $1"
	UserGetByNameQuery  = "SELECT * FROM users WHERE name = $1"
	UserGetByEmailQuery = "SELECT * FROM users WHERE email = $1"
	UserGetByPhoneQuery = "SELECT * FROM users WHERE phone = $1"
)

type PostgresUserRepository struct {
	connection *sqlx.DB
}

func NewPostgresUserRepository(connection *sqlx.DB) *PostgresUserRepository {
	return &PostgresUserRepository{connection: connection}
}

func (ur *PostgresUserRepository) Create(ctx context.Context, user domain.User) (domain.User, error) {
	pgUser := entity2.NewPgUser(user)
	queryString := entity2.InsertQueryString(pgUser, "users")

	_, err := ur.connection.NamedExecContext(ctx, queryString, pgUser)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return domain.User{}, util.WrapError(ports.ErrUserDuplicate, err)
			}
		}
		return domain.User{}, util.WrapError(ports.ErrInternalUserRepo, err)
	}

	var createdUser entity2.PgUser
	err = ur.connection.GetContext(ctx, &createdUser, UserGetByIDQuery, pgUser.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, util.WrapError(ports.ErrUserIDNotFound, err)
		}
		return domain.User{}, util.WrapError(ports.ErrInternalUserRepo, err)
	}

	return createdUser.ToDomain(), nil
}

func (ur *PostgresUserRepository) GetAll(ctx context.Context) ([]domain.User, error) {
	var users []entity2.PgUser
	err := ur.connection.SelectContext(ctx, &users, UserGetAllQuery)
	if err != nil {
		return nil, util.WrapError(ports.ErrInternalUserRepo, err)
	}

	domainUsers := make([]domain.User, len(users))
	for i, user := range users {
		domainUsers[i] = user.ToDomain()
	}

	return domainUsers, nil
}

func (ur *PostgresUserRepository) GetByID(ctx context.Context, userID uuid.UUID) (domain.User, error) {
	var foundUser entity2.PgUser
	err := ur.connection.GetContext(ctx, &foundUser, UserGetByIDQuery, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, util.WrapError(ports.ErrUserIDNotFound, err)
		}
		return domain.User{}, util.WrapError(ports.ErrInternalUserRepo, err)
	}

	return foundUser.ToDomain(), nil
}

func (ur *PostgresUserRepository) GetByName(ctx context.Context, name string) (domain.User, error) {
	var foundUser entity2.PgUser
	err := ur.connection.GetContext(ctx, &foundUser, UserGetByNameQuery, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, util.WrapError(ports.ErrUserNameNotFound, err)
		}
		return domain.User{}, util.WrapError(ports.ErrInternalUserRepo, err)
	}

	return foundUser.ToDomain(), nil
}

func (ur *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	var foundUser entity2.PgUser
	err := ur.connection.GetContext(ctx, &foundUser, UserGetByEmailQuery, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, util.WrapError(ports.ErrUserEmailNotFound, err)
		}
		return domain.User{}, util.WrapError(ports.ErrInternalUserRepo, err)
	}

	return foundUser.ToDomain(), nil
}

func (ur *PostgresUserRepository) GetByPhone(ctx context.Context, phone string) (domain.User, error) {
	var foundUser entity2.PgUser
	err := ur.connection.GetContext(ctx, &foundUser, UserGetByPhoneQuery, phone)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, util.WrapError(ports.ErrUserPhoneNotFound, err)
		}
		return domain.User{}, util.WrapError(ports.ErrInternalUserRepo, err)
	}

	return foundUser.ToDomain(), nil
}
