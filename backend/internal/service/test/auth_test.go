package test

import (
	"context"
	"errors"
	"github.com/google/uuid"
	mocks2 "github.com/hanoys/sigma-music/internal/adapters/auth/mocks"
	mocks3 "github.com/hanoys/sigma-music/internal/adapters/hash/mocks"
	"github.com/hanoys/sigma-music/internal/adapters/repository/mocks"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
	"github.com/hanoys/sigma-music/internal/service"
	"testing"
)

var credentials = ports.LogInCredentials{
	Name:     "test",
	Password: "test",
	Role:     domain.UserRole,
}

var foundUser = domain.User{
	ID:       uuid.New(),
	Name:     "test",
	Email:    "test",
	Phone:    "test",
	Password: "test",
	Salt:     "test",
	Country:  "test",
}

func TestAuthServiceLogIn(t *testing.T) {
	tests := []struct {
		name              string
		userRepoMock      func(repository *mocks.UserRepository)
		musicianRepoMock  func(repository *mocks.MusicianRepository)
		tokenProviderMock func(provider *mocks2.TokenProvider)
		hashProviderMock  func(provider *mocks3.HashPasswordProvider)
		expected          error
	}{
		{
			name: "success login",
			userRepoMock: func(repository *mocks.UserRepository) {
				repository.
					On("GetByName", context.Background(), credentials.Name).
					Return(foundUser, nil)
			},
			musicianRepoMock: func(repository *mocks.MusicianRepository) {

			},
			tokenProviderMock: func(provider *mocks2.TokenProvider) {
				provider.
					On("NewSession", context.Background(), domain.Payload{
						UserID: foundUser.ID,
						Role:   credentials.Role,
					}).Return(domain.TokenPair{}, nil)
			},
			hashProviderMock: func(provider *mocks3.HashPasswordProvider) {
				provider.
					On("ComparePasswordWithHash", credentials.Password, domain.SaltedPassword{
						HashPassword: foundUser.Password,
						Salt:         foundUser.Salt,
					}).Return(true)
			},
			expected: nil,
		},
		{
			name: "fail login: incorrect password",
			userRepoMock: func(repository *mocks.UserRepository) {
				repository.
					On("GetByName", context.Background(), credentials.Name).
					Return(foundUser, nil)
			},
			musicianRepoMock: func(repository *mocks.MusicianRepository) {

			},
			tokenProviderMock: func(provider *mocks2.TokenProvider) {

			},
			hashProviderMock: func(provider *mocks3.HashPasswordProvider) {
				provider.
					On("ComparePasswordWithHash", credentials.Password, domain.SaltedPassword{
						HashPassword: foundUser.Password,
						Salt:         foundUser.Salt,
					}).Return(false)
			},
			expected: ports.ErrIncorrectPassword,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			userRepository := mocks.NewUserRepository(t)
			musicianRepository := mocks.NewMusicianRepository(t)
			tokenProvider := mocks2.NewTokenProvider(t)
			hashProvider := mocks3.NewHashPasswordProvider(t)
			authService := service.NewAuthorizationService(userRepository, musicianRepository, tokenProvider, hashProvider)
			test.userRepoMock(userRepository)
			test.musicianRepoMock(musicianRepository)
			test.tokenProviderMock(tokenProvider)
			test.hashProviderMock(hashProvider)

			_, err := authService.LogIn(context.Background(), credentials)
			if !errors.Is(err, test.expected) {
				t.Errorf("got %v, want %v", err, test.expected)
			}
		})
	}
}

var tokenString = "testTokenString"

func TestAuthServiceVerifyToken(t *testing.T) {
	tests := []struct {
		name              string
		userRepoMock      func(repository *mocks.UserRepository)
		musicianRepoMock  func(repository *mocks.MusicianRepository)
		tokenProviderMock func(provider *mocks2.TokenProvider)
		hashProviderMock  func(provider *mocks3.HashPasswordProvider)
		expected          error
	}{
		{
			name: "success verification",
			userRepoMock: func(repository *mocks.UserRepository) {
			},
			musicianRepoMock: func(repository *mocks.MusicianRepository) {
			},
			tokenProviderMock: func(provider *mocks2.TokenProvider) {
				provider.
					On("VerifyToken", context.Background(), tokenString).
					Return(domain.Payload{}, nil)
			},
			hashProviderMock: func(provider *mocks3.HashPasswordProvider) {
			},
			expected: nil,
		},
		{
			name: "fail verification",
			userRepoMock: func(repository *mocks.UserRepository) {
			},
			musicianRepoMock: func(repository *mocks.MusicianRepository) {
			},
			tokenProviderMock: func(provider *mocks2.TokenProvider) {
				provider.
					On("VerifyToken", context.Background(), tokenString).
					Return(domain.Payload{}, ports.ErrTokenProviderInvalidToken)
			},
			hashProviderMock: func(provider *mocks3.HashPasswordProvider) {
			},
			expected: ports.ErrTokenProviderInvalidToken,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			userRepository := mocks.NewUserRepository(t)
			musicianRepository := mocks.NewMusicianRepository(t)
			tokenProvider := mocks2.NewTokenProvider(t)
			hashProvider := mocks3.NewHashPasswordProvider(t)
			authService := service.NewAuthorizationService(userRepository, musicianRepository, tokenProvider, hashProvider)
			test.userRepoMock(userRepository)
			test.musicianRepoMock(musicianRepository)
			test.tokenProviderMock(tokenProvider)
			test.hashProviderMock(hashProvider)

			_, err := authService.VerifyToken(context.Background(), tokenString)
			if !errors.Is(err, test.expected) {
				t.Errorf("got %v, want %v", err, test.expected)
			}
		})
	}
}
