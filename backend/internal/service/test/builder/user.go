package builder

import (
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
)

type UserServiceCreateRequestBuilder struct {
	obj ports.UserServiceCreateRequest
}

func NewUserServiceCreateRequestBuilder() *UserServiceCreateRequestBuilder {
	return new(UserServiceCreateRequestBuilder)
}

func (b *UserServiceCreateRequestBuilder) Build() ports.UserServiceCreateRequest {
	return b.obj
}

func (b *UserServiceCreateRequestBuilder) Default() *UserServiceCreateRequestBuilder {
	b.obj = ports.UserServiceCreateRequest{
		Name:     "test",
		Email:    "test@mail.com",
		Phone:    "+79999999999",
		Password: "test",
		Country:  "Russia",
	}
	return b
}

func (b *UserServiceCreateRequestBuilder) SetName(name string) *UserServiceCreateRequestBuilder {
	b.obj.Name = name
	return b
}

func (b *UserServiceCreateRequestBuilder) SetEmail(email string) *UserServiceCreateRequestBuilder {
	b.obj.Email = email
	return b
}

func (b *UserServiceCreateRequestBuilder) SetPhone(phone string) *UserServiceCreateRequestBuilder {
	b.obj.Phone = phone
	return b
}

func (b *UserServiceCreateRequestBuilder) SetPassword(password string) *UserServiceCreateRequestBuilder {
	b.obj.Password = password
	return b
}

func (b *UserServiceCreateRequestBuilder) SetCountry(country string) *UserServiceCreateRequestBuilder {
	b.obj.Country = country
	return b
}

type UserBuilder struct {
	obj domain.User
}

func NewUserBuilder() *UserBuilder {
	return new(UserBuilder)
}

func (b *UserBuilder) Build() domain.User {
	return b.obj
}

func (b *UserBuilder) Default() *UserBuilder {
	b.obj = domain.User{
		ID:       uuid.New(),
		Name:     "test",
		Email:    "test@mail.com",
		Phone:    "+79999999999",
		Password: "test",
		Country:  "Russia",
		Salt:     "test",
	}
	return b
}

func (b *UserBuilder) SetID(id uuid.UUID) *UserBuilder {
	b.obj.ID = id
	return b
}

func (b *UserBuilder) SetName(name string) *UserBuilder {
	b.obj.Name = name
	return b
}

func (b *UserBuilder) SetEmail(email string) *UserBuilder {
	b.obj.Email = email
	return b
}

func (b *UserBuilder) SetPhone(phone string) *UserBuilder {
	b.obj.Phone = phone
	return b
}

func (b *UserBuilder) SetPassword(password string) *UserBuilder {
	b.obj.Password = password
	return b
}

func (b *UserBuilder) SetSalt(salt string) *UserBuilder {
	b.obj.Salt = salt
	return b
}

func (b *UserBuilder) SetCountry(country string) *UserBuilder {
	b.obj.Country = country
	return b
}
