package entity

import (
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
)

type PgUser struct {
	ID       uuid.UUID `db:"id"`
	Name     string    `db:"name"`
	Email    string    `db:"email"`
	Phone    string    `db:"phone"`
	Password string    `db:"password"`
	Salt     string    `db:"salt"`
	Country  string    `db:"country"`
}

func (u *PgUser) ToDomain() domain.User {
	return domain.User{
		ID:       u.ID,
		Name:     u.Name,
		Email:    u.Email,
		Phone:    u.Phone,
		Password: u.Password,
		Salt:     u.Salt,
		Country:  u.Country,
	}
}

func NewPgUser(user domain.User) PgUser {
	return PgUser{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Phone:    user.Phone,
		Password: user.Password,
		Salt:     user.Salt,
		Country:  user.Country,
	}
}
