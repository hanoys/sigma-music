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
	userGetAllQuery     = "SELECT * FROM users"
	userGetByIDQuery    = "SELECT * FROM users WHERE id = $1"
	userGetByNameQuery  = "SELECT * FROM users WHERE name = $1"
	userGetByEmailQuery = "SELECT * FROM users WHERE email = $1"
	userGetByPhoneQuery = "SELECT * FROM users WHERE phone = $1"
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
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return domain.User{}, util.WrapError(ports.ErrUserDuplicate, err)
			}
		}
		return domain.User{}, util.WrapError(ports.ErrInternalUserRepo, err)
	}

	var createdUser entity.PgUser
	err = ur.db.GetContext(ctx, &createdUser, userGetByIDQuery, pgUser.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, util.WrapError(ports.ErrUserIDNotFound, err)
		}
		return domain.User{}, util.WrapError(ports.ErrInternalUserRepo, err)
	}

	return createdUser.ToDomain(), nil
}

func (ur *PostgresUserRepository) GetAll(ctx context.Context) ([]domain.User, error) {
	var users []entity.PgUser
	err := ur.db.SelectContext(ctx, &users, userGetAllQuery)
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
	var foundUser entity.PgUser
	err := ur.db.GetContext(ctx, &foundUser, userGetByIDQuery, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, util.WrapError(ports.ErrUserIDNotFound, err)
		}
		return domain.User{}, util.WrapError(ports.ErrInternalUserRepo, err)
	}

	return foundUser.ToDomain(), nil
}

func (ur *PostgresUserRepository) GetByName(ctx context.Context, name string) (domain.User, error) {
	var foundUser entity.PgUser
	err := ur.db.GetContext(ctx, &foundUser, userGetByNameQuery, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, util.WrapError(ports.ErrUserNameNotFound, err)
		}
		return domain.User{}, util.WrapError(ports.ErrInternalUserRepo, err)
	}

	return foundUser.ToDomain(), nil
}

func (ur *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	var foundUser entity.PgUser
	err := ur.db.GetContext(ctx, &foundUser, userGetByEmailQuery, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, util.WrapError(ports.ErrUserEmailNotFound, err)
		}
		return domain.User{}, util.WrapError(ports.ErrInternalUserRepo, err)
	}

	return foundUser.ToDomain(), nil
}

func (ur *PostgresUserRepository) GetByPhone(ctx context.Context, phone string) (domain.User, error) {
	var foundUser entity.PgUser
	err := ur.db.GetContext(ctx, &foundUser, userGetByPhoneQuery, phone)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, util.WrapError(ports.ErrUserPhoneNotFound, err)
		}
		return domain.User{}, util.WrapError(ports.ErrInternalUserRepo, err)
	}

	return foundUser.ToDomain(), nil
}
