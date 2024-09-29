package test

import (
	"context"

	"github.com/hanoys/sigma-music/internal/adapters/auth/mocks"
	mocks2 "github.com/hanoys/sigma-music/internal/adapters/hash/mocks"
	mocks3 "github.com/hanoys/sigma-music/internal/adapters/repository/mocks"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"go.uber.org/zap"
)

type AuthSuite struct {
	suite.Suite
	logger        *zap.Logger
	tokenProvider *mocks.TokenProvider
	hashProvider  *mocks2.HashPasswordProvider
}

func (s *AuthSuite) BeforeEach(t provider.T) {
	loggerBuilder := zap.NewDevelopmentConfig()
	loggerBuilder.Level = zap.NewAtomicLevelAt(zap.FatalLevel)
	s.logger, _ = loggerBuilder.Build()
	s.tokenProvider = mocks.NewTokenProvider(t)
	s.hashProvider = mocks2.NewHashPasswordProvider(t)
}

type AuthLogInSuite struct {
	AuthSuite
}

func (s *AuthLogInSuite) CorrectRepositoryMock(userRepository *mocks3.UserRepository, musicianRepository *mocks3.MusicianRepository) {
	userRepository.
		On("GetByName", context.Background(), credentials.Name).
		Return(foundUser, nil)

	s.tokenProvider.
		On("NewSession", context.Background(), domain.Payload{
			UserID: foundUser.ID,
			Role:   domain.UserRole,
		}).Return(domain.TokenPair{}, nil)

	s.hashProvider.
		On("ComparePasswordWithHash", credentials.Password, domain.SaltedPassword{
			HashPassword: foundUser.Password,
			Salt:         foundUser.Salt,
		}).Return(true)
}

func (s *AuthLogInSuite) TestCorrect(t provider.T) {
    t.Parallel()
}
