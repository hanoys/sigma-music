package test

import (
	"context"
	"testing"

	"github.com/hanoys/sigma-music/internal/adapters/auth/mocks"
	mocks2 "github.com/hanoys/sigma-music/internal/adapters/hash/mocks"
	mocks3 "github.com/hanoys/sigma-music/internal/adapters/repository/mocks"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
	"github.com/hanoys/sigma-music/internal/service"
	"github.com/hanoys/sigma-music/internal/service/test/builder"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"go.uber.org/zap"
)

type AuthSuite struct {
	suite.Suite
	logger       *zap.Logger
	hashProvider *mocks2.HashPasswordProvider
}

func (s *AuthSuite) BeforeEach(t provider.T) {
	loggerBuilder := zap.NewDevelopmentConfig()
	loggerBuilder.Level = zap.NewAtomicLevelAt(zap.FatalLevel)
	s.logger, _ = loggerBuilder.Build()
	s.hashProvider = mocks2.NewHashPasswordProvider(t)
}

type AuthLogInSuite struct {
	AuthSuite
}

func (s *AuthLogInSuite) CorrectRepositoryMock(userRepository *mocks3.UserRepository,
	musicianRepository *mocks3.MusicianRepository, tokenProvider *mocks.TokenProvider,
	user domain.User) {
	userRepository.
		On("GetByName", context.Background(), user.Name).
		Return(user, nil)

	tokenProvider.
		On("NewSession", context.Background(), domain.Payload{
			UserID: user.ID,
			Role:   domain.UserRole,
		}).Return(domain.TokenPair{}, nil)

	s.hashProvider.
		On("ComparePasswordWithHash", user.Password, domain.SaltedPassword{
			HashPassword: user.Password,
			Salt:         user.Salt,
		}).Return(true)
}

func (s *AuthLogInSuite) TestCorrect(t provider.T) {
	t.Parallel()
	t.Title("Auth login test correct")
	user := builder.NewUserBuilder().Default().Build()
	loginCred := builder.NewLoginCredentialsMother(user.Name, user.Password).Create()
	tokenProvider := mocks.NewTokenProvider(t)
	userRepository := mocks3.NewUserRepository(t)
	musicianRepository := mocks3.NewMusicianRepository(t)
	authService := service.NewAuthorizationService(userRepository, musicianRepository,
		tokenProvider, s.hashProvider, s.logger)
	s.CorrectRepositoryMock(userRepository, musicianRepository, tokenProvider, user)

	_, err := authService.LogIn(context.Background(), loginCred)

	t.Assert().Nil(err)
}

func (s *AuthLogInSuite) ErrorRepositoryMock(userRepository *mocks3.UserRepository,
	musicianRepository *mocks3.MusicianRepository, tokenProvider *mocks.TokenProvider,
	user domain.User) {
	userRepository.
		On("GetByName", context.Background(), user.Name).
		Return(domain.User{}, ports.ErrUserNameNotFound)

	musicianRepository.
		On("GetByName", context.Background(), user.Name).
		Return(domain.Musician{}, ports.ErrMusicianNameNotFound)

}

func (s *AuthLogInSuite) TestError(t provider.T) {
	t.Parallel()
	t.Title("Auth login test correct")
	user := builder.NewUserBuilder().Default().Build()
	loginCred := builder.NewLoginCredentialsMother(user.Name, user.Password).Create()
	tokenProvider := mocks.NewTokenProvider(t)
	userRepository := mocks3.NewUserRepository(t)
	musicianRepository := mocks3.NewMusicianRepository(t)
	authService := service.NewAuthorizationService(userRepository, musicianRepository,
		tokenProvider, s.hashProvider, s.logger)
	s.ErrorRepositoryMock(userRepository, musicianRepository, tokenProvider, user)

	_, err := authService.LogIn(context.Background(), loginCred)

	t.Assert().NotNil(err)
}

func TestAuthLogInSuite(t *testing.T) {
	suite.RunSuite(t, new(AuthLogInSuite))
}

type AuthLogOutSuite struct {
	AuthSuite
}

func (s *AuthLogOutSuite) CorrectRepositoryMock(userRepository *mocks3.UserRepository,
	musicianRepository *mocks3.MusicianRepository, tokenProvider *mocks.TokenProvider,
	tokenString string) {
	tokenProvider.
		On("CloseSession", context.Background(), tokenString).
		Return(nil)
}

func (s *AuthLogOutSuite) TestCorrect(t provider.T) {
	t.Parallel()
	t.Title("Auth logout test correct")
	tokenString := "tokenstring"
	tokenProvider := mocks.NewTokenProvider(t)
	userRepository := mocks3.NewUserRepository(t)
	musicianRepository := mocks3.NewMusicianRepository(t)
	authService := service.NewAuthorizationService(userRepository, musicianRepository,
		tokenProvider, s.hashProvider, s.logger)
	s.CorrectRepositoryMock(userRepository, musicianRepository, tokenProvider, tokenString)

	err := authService.LogOut(context.Background(), tokenString)

	t.Assert().Nil(err)
}

func (s *AuthLogOutSuite) TokenExpiredRepositoryMock(userRepository *mocks3.UserRepository,
	musicianRepository *mocks3.MusicianRepository, tokenProvider *mocks.TokenProvider,
	tokenString string, payload domain.Payload) {
	tokenProvider.
		On("CloseSession", context.Background(), tokenString).
		Return(ports.ErrInternalTokenProvider)

	tokenProvider.
		On("VerifyToken", context.Background(), tokenString).
		Return(payload, ports.ErrTokenProviderExpiredToken)
}

func (s *AuthLogOutSuite) TestTokenExpired(t provider.T) {
	t.Parallel()
	t.Title("Auth logout test token expired")
	tokenString := "tokenstring"
	payload := builder.NewPayloadBuilder().Default().Build()
	tokenProvider := mocks.NewTokenProvider(t)
	userRepository := mocks3.NewUserRepository(t)
	musicianRepository := mocks3.NewMusicianRepository(t)
	authService := service.NewAuthorizationService(userRepository, musicianRepository,
		tokenProvider, s.hashProvider, s.logger)
	s.TokenExpiredRepositoryMock(userRepository, musicianRepository, tokenProvider, tokenString, payload)

	err := authService.LogOut(context.Background(), tokenString)

	t.Assert().ErrorIs(err, ports.ErrInternalTokenProvider)
}

func TestAuthLogOutSuite(t *testing.T) {
	suite.RunSuite(t, new(AuthLogOutSuite))
}

type AuthVerifyTokenSuite struct {
	AuthSuite
}

func (s *AuthVerifyTokenSuite) CorrectRepositoryMock(tokenProvider *mocks.TokenProvider,
	tokenString string, payload domain.Payload) {
	tokenProvider.
		On("VerifyToken", context.Background(), tokenString).
		Return(payload, nil)
}

func (s *AuthVerifyTokenSuite) TestCorrect(t provider.T) {
	t.Parallel()
	t.Title("Auth verify token test correct")
	tokenString := "tokenstring"
	payload := builder.NewPayloadBuilder().Default().Build()
	tokenProvider := mocks.NewTokenProvider(t)
	userRepository := mocks3.NewUserRepository(t)
	musicianRepository := mocks3.NewMusicianRepository(t)
	authService := service.NewAuthorizationService(userRepository, musicianRepository,
		tokenProvider, s.hashProvider, s.logger)
	s.CorrectRepositoryMock(tokenProvider, tokenString, payload)

	servicePayload, err := authService.VerifyToken(context.Background(), tokenString)

	t.Assert().Equal(payload.UserID, servicePayload.UserID)
	t.Assert().Equal(payload.Role, servicePayload.Role)
	t.Assert().Nil(err)
}

func (s *AuthVerifyTokenSuite) InvalidTokenRepositoryMock(tokenProvider *mocks.TokenProvider,
	tokenString string, payload domain.Payload) {
	tokenProvider.
		On("VerifyToken", context.Background(), tokenString).
		Return(domain.Payload{}, ports.ErrTokenProviderInvalidToken)
}

func (s *AuthVerifyTokenSuite) TestInvalidToken(t provider.T) {
	t.Parallel()
	t.Title("Auth verify token test invalid token")
	tokenString := "tokenstring"
	payload := builder.NewPayloadBuilder().Default().Build()
	tokenProvider := mocks.NewTokenProvider(t)
	userRepository := mocks3.NewUserRepository(t)
	musicianRepository := mocks3.NewMusicianRepository(t)
	authService := service.NewAuthorizationService(userRepository, musicianRepository,
		tokenProvider, s.hashProvider, s.logger)
	s.InvalidTokenRepositoryMock(tokenProvider, tokenString, payload)

	servicePayload, err := authService.VerifyToken(context.Background(), tokenString)

	t.Assert().Equal(domain.Payload{}, servicePayload)
	t.Assert().ErrorIs(err, ports.ErrInvalidToken)
}

func TestAuthVerifyTokenSuite(t *testing.T) {
	suite.RunSuite(t, new(AuthVerifyTokenSuite))
}

type AuthRefreshTokenSuite struct {
	AuthSuite
}

func (s *AuthRefreshTokenSuite) CorrectRepositoryMock(tokenProvider *mocks.TokenProvider,
	refreshTokenString string, tokenPair domain.TokenPair) {
	tokenProvider.
		On("RefreshSession", context.Background(), refreshTokenString).
		Return(tokenPair, nil)
}

func (s *AuthRefreshTokenSuite) TestCorrect(t provider.T) {
	t.Parallel()
	t.Title("Auth refresh token test correct")
	tokenString := "tokenstring"
	tokenPair := builder.NewTokenPairBuilder().Default().Build()
	tokenProvider := mocks.NewTokenProvider(t)
	userRepository := mocks3.NewUserRepository(t)
	musicianRepository := mocks3.NewMusicianRepository(t)
	authService := service.NewAuthorizationService(userRepository, musicianRepository,
		tokenProvider, s.hashProvider, s.logger)
	s.CorrectRepositoryMock(tokenProvider, tokenString, tokenPair)

	serviceTokenPair, err := authService.RefreshToken(context.Background(), tokenString)

	t.Assert().Equal(tokenPair, serviceTokenPair)
	t.Assert().Nil(err)
}

func TestAuthRefreshTokenSuite(t *testing.T) {
	suite.RunSuite(t, new(AuthRefreshTokenSuite))
}
