package postgres

import (
	"context"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func (ur *UserRepository) Create(ctx context.Context, user domain.User) (domain.User, error) {
	var createdUser domain.User
	err := ur.db.QueryRow(ctx,
		"INSERT INTO users(id, name, email, phone, password, country) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *",
		user.ID, user.Name, user.Email, user.Phone, user.Password, user.Country).Scan(
		&createdUser.ID, &createdUser.Name, &createdUser.Email, &createdUser.Phone, &createdUser.Password,
		&createdUser.Country)

	if err != nil {
		return domain.User{}, err
	}

	return createdUser, nil
}

func (ur *UserRepository) getByUniqueColumn(ctx context.Context, column string, value string) (domain.User, error) {
	var foundUser domain.User
	err := ur.db.QueryRow(ctx,
		"SELECT * FROM users WHERE $1 = $2", column, value).Scan(
		&foundUser.ID, &foundUser.Name, &foundUser.Email, &foundUser.Phone, &foundUser.Password,
		&foundUser.Country)

	if err != nil {
		return domain.User{}, err
	}

	return foundUser, nil
}

func (ur *UserRepository) GetByName(ctx context.Context, name string) (domain.User, error) {
	return ur.getByUniqueColumn(ctx, "name", name)
}

func (ur *UserRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	return ur.getByUniqueColumn(ctx, "email", email)
}

func (ur *UserRepository) GetByPhone(ctx context.Context, phone string) (domain.User, error) {
	return ur.getByUniqueColumn(ctx, "phone", phone)
}
