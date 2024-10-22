package entity

import (
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
)

type GormUser struct {
	ID       uuid.UUID `gorm:"column: id"`
	Name     string    `gorm:"column: name"`
	Email    string    `gorm:"column: email"`
	Phone    string    `gorm:"column: phone"`
	Password string    `gorm:"column: password"`
	Salt     string    `gorm:"column: salt"`
	Country  string    `gorm:"column: country"`
}

func (u *GormUser) ToDomain() domain.User {
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

func NewGORMUser(user domain.User) GormUser {
	return GormUser{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Phone:    user.Phone,
		Password: user.Password,
		Salt:     user.Salt,
		Country:  user.Country,
	}
}
