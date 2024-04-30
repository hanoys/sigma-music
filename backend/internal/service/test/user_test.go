package test

import (
	"context"
	"errors"
	"github.com/hanoys/sigma-music/internal/adapters/repository/mocks"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
	"github.com/hanoys/sigma-music/internal/service"
	"github.com/stretchr/testify/mock"
	"testing"
)

var createUserReq = ports.UserServiceCreateRequest{
	Name:     "TestUser",
	Email:    "email",
	Phone:    "+79999999999",
	Password: "password",
	Country:  "testCountry",
}

func TestUserServiceRegister(t *testing.T) {
	tests := []struct {
		name           string
		repositoryMock func(repository *mocks.UserRepository)
		req            ports.UserServiceCreateRequest
		expected       error
	}{
		{
			name: "register success",
			req:  createUserReq,
			repositoryMock: func(repository *mocks.UserRepository) {
				repository.
					On("Create", context.Background(), mock.AnythingOfType("domain.User")).
					Return(domain.User{}, nil).
					On("GetByName", context.Background(), createUserReq.Name).
					Return(domain.User{}, ports.ErrUserNameNotFound).
					On("GetByEmail", context.Background(), createUserReq.Email).
					Return(domain.User{}, ports.ErrUserEmailNotFound).
					On("GetByPhone", context.Background(), createUserReq.Phone).
					Return(domain.User{}, ports.ErrUserPhoneNotFound)
			},
			expected: nil,
		},
		{
			name: "user name already exists",
			req:  createUserReq,
			repositoryMock: func(repository *mocks.UserRepository) {
				repository.
					On("GetByName", context.Background(), createUserReq.Name).
					Return(domain.User{}, nil)
			},
			expected: ports.ErrUserWithSuchNameAlreadyExists,
		},
		{
			name: "user email already exists",
			req:  createUserReq,
			repositoryMock: func(repository *mocks.UserRepository) {
				repository.
					On("GetByName", context.Background(), createUserReq.Name).
					Return(domain.User{}, ports.ErrUserNameNotFound).
					On("GetByEmail", context.Background(), createUserReq.Email).
					Return(domain.User{}, nil)
			},
			expected: ports.ErrUserWithSuchEmailAlreadyExists,
		},
		{
			name: "user phone already exists",
			req:  createUserReq,
			repositoryMock: func(repository *mocks.UserRepository) {
				repository.
					On("GetByName", context.Background(), createUserReq.Name).
					Return(domain.User{}, ports.ErrUserNameNotFound).
					On("GetByEmail", context.Background(), createUserReq.Email).
					Return(domain.User{}, ports.ErrUserEmailNotFound).
					On("GetByPhone", context.Background(), createUserReq.Phone).
					Return(domain.User{}, nil)
			},
			expected: ports.ErrUserWithSuchPhoneAlreadyExists,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			userRepository := mocks.NewUserRepository(t)
			userService := service.NewUserService(userRepository)
			test.repositoryMock(userRepository)

			_, err := userService.Register(context.Background(), test.req)
			if !errors.Is(err, test.expected) {
				t.Errorf("got %v, want %v", err, test.expected)
			}
		})
	}
}
