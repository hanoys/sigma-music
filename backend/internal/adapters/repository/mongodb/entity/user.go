package entity

import (
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
)

type MongoUser struct {
	ID       uuid.UUID `bson:"_id"`
	Name     string    `bson:"name"`
	Email    string    `bson:"email"`
	Phone    string    `bson:"phone"`
	Password string    `bson:"password"`
	Salt     string    `bson:"salt"`
	Country  string    `bson:"country"`
}

func (u *MongoUser) ToDomain() domain.User {
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

func NewMongoUser(user domain.User) MongoUser {
	return MongoUser{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Phone:    user.Phone,
		Password: user.Password,
		Salt:     user.Salt,
		Country:  user.Country,
	}
}
