package gorm

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/adapters/repository/postgres/gorm/entity"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
	"github.com/hanoys/sigma-music/internal/util"
	"gorm.io/gorm"
)

type PostgresUserRepository struct {
	connection *gorm.DB
}

func NewPostgresUserRepository(connection *gorm.DB) *PostgresUserRepository {
	return &PostgresUserRepository{
		connection: connection,
	}
}

func (ur *PostgresUserRepository) Create(ctx context.Context, user domain.User) (domain.User, error) {
	gormUser := entity.NewGORMUser(user)
	result := ur.connection.Table("users").WithContext(ctx).Create(gormUser)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return domain.User{}, util.WrapError(ports.ErrUserDuplicate, result.Error)
		}

		return domain.User{}, util.WrapError(ports.ErrInternalUserRepo, result.Error)
	}

	var createdUser entity.GormUser
	result = ur.connection.Table("users").WithContext(ctx).Take(&createdUser, "id = ?", gormUser.ID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return domain.User{}, util.WrapError(ports.ErrUserIDNotFound, result.Error)
		}

		return domain.User{}, util.WrapError(ports.ErrInternalUserRepo, result.Error)
	}

	return createdUser.ToDomain(), nil
}

func (ur *PostgresUserRepository) GetAll(ctx context.Context) ([]domain.User, error) {
	var users []entity.GormUser
	result := ur.connection.Table("users").WithContext(ctx).Find(&users)
	if result.Error != nil {
		return nil, util.WrapError(ports.ErrInternalUserRepo, result.Error)
	}

	domainUsers := make([]domain.User, len(users))
	for i, user := range users {
		domainUsers[i] = user.ToDomain()
	}

	return domainUsers, nil
}

func (ur *PostgresUserRepository) GetByID(ctx context.Context, userID uuid.UUID) (domain.User, error) {
	var foundUser entity.GormUser
	result := ur.connection.Table("users").WithContext(ctx).Take(&foundUser, "id = ?", userID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return domain.User{}, util.WrapError(ports.ErrUserIDNotFound, result.Error)
		}

		return domain.User{}, util.WrapError(ports.ErrInternalUserRepo, result.Error)
	}

	return foundUser.ToDomain(), nil
}

func (ur *PostgresUserRepository) GetByName(ctx context.Context, name string) (domain.User, error) {
	var foundUser entity.GormUser
	result := ur.connection.Table("users").WithContext(ctx).Take(&foundUser, "name = ?", name)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return domain.User{}, util.WrapError(ports.ErrUserIDNotFound, result.Error)
		}

		return domain.User{}, util.WrapError(ports.ErrInternalUserRepo, result.Error)
	}

	return foundUser.ToDomain(), nil
}

func (ur *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	var foundUser entity.GormUser
	result := ur.connection.Table("users").WithContext(ctx).Take(&foundUser, "email = ?", email)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return domain.User{}, util.WrapError(ports.ErrUserIDNotFound, result.Error)
		}

		return domain.User{}, util.WrapError(ports.ErrInternalUserRepo, result.Error)
	}

	return foundUser.ToDomain(), nil
}

func (ur *PostgresUserRepository) GetByPhone(ctx context.Context, phone string) (domain.User, error) {
	var foundUser entity.GormUser
	result := ur.connection.Table("users").WithContext(ctx).Take(&foundUser, "phone = ?", phone)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return domain.User{}, util.WrapError(ports.ErrUserIDNotFound, result.Error)
		}

		return domain.User{}, util.WrapError(ports.ErrInternalUserRepo, result.Error)
	}

	return foundUser.ToDomain(), nil
}
