package test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/adapters/hash"
	"github.com/hanoys/sigma-music/internal/adapters/repository/mocks"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
	"github.com/hanoys/sigma-music/internal/service"
	"github.com/hanoys/sigma-music/internal/service/test/builder"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type UserSuite struct {
	suite.Suite
	logger       *zap.Logger
	hashProvider *hash.HashPasswordProvider
}

func (s *UserSuite) BeforeEach(t provider.T) {
	loggerBuilder := zap.NewDevelopmentConfig()
	loggerBuilder.Level = zap.NewAtomicLevelAt(zap.FatalLevel)
	s.logger, _ = loggerBuilder.Build()

	s.hashProvider = hash.NewHashPasswordProvider()
}

type UserRegisterSuite struct {
	UserSuite
}

func (s *UserRegisterSuite) CorrectRepositoryMock(repository *mocks.UserRepository, req ports.UserServiceCreateRequest) {
	repository.
		On("Create", context.Background(), mock.AnythingOfType("domain.User")).
		Return(domain.User{}, nil).
		On("GetByName", context.Background(), req.Name).
		Return(domain.User{}, ports.ErrUserNameNotFound).
		On("GetByEmail", context.Background(), req.Email).
		Return(domain.User{}, ports.ErrUserEmailNotFound).
		On("GetByPhone", context.Background(), req.Phone).
		Return(domain.User{}, ports.ErrUserPhoneNotFound)

}

func (s *UserRegisterSuite) TestCorrect(t provider.T) {
	t.Parallel()
	t.Title("User register test correct")
	req := builder.NewUserServiceCreateRequestBuilder().Default().Build()
	repository := mocks.NewUserRepository(t)
	userService := service.NewUserService(repository, s.hashProvider, s.logger)
	s.CorrectRepositoryMock(repository, req)

	_, err := userService.Register(context.Background(), req)

	t.Assert().Nil(err)
}

func (s *UserRegisterSuite) NameExistsRepositoryMock(repository *mocks.UserRepository, req ports.UserServiceCreateRequest) {
	repository.
		On("GetByName", context.Background(), req.Name).
		Return(domain.User{}, nil)
}

func (s *UserRegisterSuite) TestNameExists(t provider.T) {
	t.Parallel()
	t.Title("User register test name exists")
	req := builder.NewUserServiceCreateRequestBuilder().Default().Build()
	repository := mocks.NewUserRepository(t)
	userService := service.NewUserService(repository, s.hashProvider, s.logger)
	s.NameExistsRepositoryMock(repository, req)

	_, err := userService.Register(context.Background(), req)

	t.Assert().ErrorIs(err, ports.ErrUserWithSuchNameAlreadyExists)
}

func (s *UserRegisterSuite) EmailExistsRepositoryMock(repository *mocks.UserRepository, req ports.UserServiceCreateRequest) {
	repository.
		On("GetByName", context.Background(), req.Name).
		Return(domain.User{}, ports.ErrUserNameNotFound).
		On("GetByEmail", context.Background(), req.Email).
		Return(domain.User{}, nil)
}

func (s *UserRegisterSuite) TestEmailExists(t provider.T) {
	t.Parallel()
	t.Title("User register test email exists")
	req := builder.NewUserServiceCreateRequestBuilder().Default().Build()
	repository := mocks.NewUserRepository(t)
	userService := service.NewUserService(repository, s.hashProvider, s.logger)
	s.EmailExistsRepositoryMock(repository, req)

	_, err := userService.Register(context.Background(), req)

	t.Assert().ErrorIs(err, ports.ErrUserWithSuchEmailAlreadyExists)
}

func (s *UserRegisterSuite) PhoneExistsRepositoryMock(repository *mocks.UserRepository, req ports.UserServiceCreateRequest) {
	repository.
		On("GetByName", context.Background(), req.Name).
		Return(domain.User{}, ports.ErrUserNameNotFound).
		On("GetByEmail", context.Background(), req.Email).
		Return(domain.User{}, ports.ErrUserEmailNotFound).
		On("GetByPhone", context.Background(), req.Phone).
		Return(domain.User{}, nil)
}

func (s *UserRegisterSuite) TestPhoneExists(t provider.T) {
	t.Parallel()
	t.Title("User register test phone exists")
	req := builder.NewUserServiceCreateRequestBuilder().Default().Build()
	repository := mocks.NewUserRepository(t)
	userService := service.NewUserService(repository, s.hashProvider, s.logger)
	s.PhoneExistsRepositoryMock(repository, req)

	_, err := userService.Register(context.Background(), req)

	t.Assert().ErrorIs(err, ports.ErrUserWithSuchPhoneAlreadyExists)
}

func TestUserRegisterSuite(t *testing.T) {
	suite.RunSuite(t, new(UserRegisterSuite))
}

type UserGetAllSuite struct {
	UserSuite
}

func (s *UserGetAllSuite) RepositoryErrorRepositoryMock(repository *mocks.UserRepository) {
	repository.
		On("GetAll", context.Background()).
		Return(nil, ports.ErrInternalUserRepo)
}

func (s *UserGetAllSuite) TestRepositoryError(t provider.T) {
	t.Parallel()
	t.Title("User get all test error")
	repository := mocks.NewUserRepository(t)
	userService := service.NewUserService(repository, s.hashProvider, s.logger)
	s.RepositoryErrorRepositoryMock(repository)

	_, err := userService.GetAll(context.Background())

	t.Assert().ErrorIs(err, ports.ErrInternalUserRepo)
}

func (s *UserGetAllSuite) SuccessRepositoryMock(repository *mocks.UserRepository) {
	repository.
		On("GetAll", context.Background()).
		Return([]domain.User{}, nil)
}

func (s *UserGetAllSuite) TestSuccess(t provider.T) {
	t.Parallel()
	t.Title("User get all test success")
	repository := mocks.NewUserRepository(t)
	userService := service.NewUserService(repository, s.hashProvider, s.logger)
	s.SuccessRepositoryMock(repository)

	_, err := userService.GetAll(context.Background())

	t.Assert().Nil(err)
}

func TestUserGetAllSuite(t *testing.T) {
	suite.RunSuite(t, new(UserGetAllSuite))
}

type UserGetByIdSuite struct {
	UserSuite
}

func (s *UserGetByIdSuite) IdNotFoundRepositoryMock(repository *mocks.UserRepository) {
	repository.
		On("GetByID", context.Background(), mock.AnythingOfType("uuid.UUID")).
		Return(domain.User{}, ports.ErrUserIDNotFound)
}

func (s *UserGetByIdSuite) TestIDNotFound(t provider.T) {
	t.Parallel()
	t.Title("User get by id test id not found")
	repository := mocks.NewUserRepository(t)
	userService := service.NewUserService(repository, s.hashProvider, s.logger)
	s.IdNotFoundRepositoryMock(repository)

	_, err := userService.GetById(context.Background(), uuid.New())

	t.Assert().ErrorIs(err, ports.ErrUserIDNotFound)
}

func (s *UserGetByIdSuite) SuccessRepositoryMock(repository *mocks.UserRepository) {
	repository.
		On("GetByID", context.Background(), mock.AnythingOfType("uuid.UUID")).
		Return(domain.User{}, nil)
}

func (s *UserGetByIdSuite) TestSuccess(t provider.T) {
	t.Parallel()
	t.Title("User get by id test success")
	repository := mocks.NewUserRepository(t)
	userService := service.NewUserService(repository, s.hashProvider, s.logger)
	s.SuccessRepositoryMock(repository)

	_, err := userService.GetById(context.Background(), uuid.New())

	t.Assert().Nil(err)
}

func TestUserGetByIdSuite(t *testing.T) {
	suite.RunSuite(t, new(UserGetByIdSuite))
}

// UserGetByNameSuite
type UserGetByNameSuite struct {
	UserSuite
}

func (s *UserGetByNameSuite) NameNotFoundRepositoryMock(repository *mocks.UserRepository) {
	repository.
		On("GetByName", context.Background(), mock.AnythingOfType("string")).
		Return(domain.User{}, ports.ErrUserNameNotFound)
}

func (s *UserGetByNameSuite) TestNameNotFound(t provider.T) {
	t.Parallel()
	t.Title("User get by name test not found")
	repository := mocks.NewUserRepository(t)
	userService := service.NewUserService(repository, s.hashProvider, s.logger)
	s.NameNotFoundRepositoryMock(repository)

	_, err := userService.GetByName(context.Background(), "")

	t.Assert().ErrorIs(err, ports.ErrUserNameNotFound)
}

func (s *UserGetByNameSuite) SuccessRepositoryMock(repository *mocks.UserRepository) {
	repository.
		On("GetByName", context.Background(), mock.AnythingOfType("string")).
		Return(domain.User{}, nil)
}

func (s *UserGetByNameSuite) TestSuccess(t provider.T) {
	t.Parallel()
	t.Title("User get by name test success")
	repository := mocks.NewUserRepository(t)
	userService := service.NewUserService(repository, s.hashProvider, s.logger)
	s.SuccessRepositoryMock(repository)

	_, err := userService.GetByName(context.Background(), "")

	t.Assert().Nil(err)
}

func TestUserGetByNameSuite(t *testing.T) {
	suite.RunSuite(t, new(UserGetByNameSuite))
}

// UserGetByEmailSuite
type UserGetByEmailSuite struct {
	UserSuite
}

func (s *UserGetByEmailSuite) EmailNotFoundRepositoryMock(repository *mocks.UserRepository) {
	repository.
		On("GetByEmail", context.Background(), mock.AnythingOfType("string")).
		Return(domain.User{}, ports.ErrUserEmailNotFound)
}

func (s *UserGetByEmailSuite) TestEmailNotFoundError(t provider.T) {
	t.Parallel()
	t.Title("User get by email test not found")
	repository := mocks.NewUserRepository(t)
	userService := service.NewUserService(repository, s.hashProvider, s.logger)
	s.EmailNotFoundRepositoryMock(repository)

	_, err := userService.GetByEmail(context.Background(), "")

	t.Assert().ErrorIs(err, ports.ErrUserEmailNotFound)
}

func (s *UserGetByEmailSuite) SuccessRepositoryMock(repository *mocks.UserRepository) {
	repository.
		On("GetByEmail", context.Background(), mock.AnythingOfType("string")).
		Return(domain.User{}, nil)
}

func (s *UserGetByEmailSuite) TestSuccess(t provider.T) {
	t.Parallel()
	t.Title("User get by email test success")
	repository := mocks.NewUserRepository(t)
	userService := service.NewUserService(repository, s.hashProvider, s.logger)
	s.SuccessRepositoryMock(repository)

	_, err := userService.GetByEmail(context.Background(), "")

	t.Assert().Nil(err)
}

func TestUserGetByEmailSuite(t *testing.T) {
	suite.RunSuite(t, new(UserGetByEmailSuite))
}

type UserGetByPhoneSuite struct {
	UserSuite
}

func (s *UserGetByPhoneSuite) PhoneNotFoundRepositoryMock(repository *mocks.UserRepository) {
	repository.
		On("GetByPhone", context.Background(), mock.AnythingOfType("string")).
		Return(domain.User{}, ports.ErrUserPhoneNotFound)
}

func (s *UserGetByPhoneSuite) TestRepositoryError(t provider.T) {
	t.Parallel()
	t.Title("User get by phone test error")
	repository := mocks.NewUserRepository(t)
	userService := service.NewUserService(repository, s.hashProvider, s.logger)
	s.PhoneNotFoundRepositoryMock(repository)

	_, err := userService.GetByPhone(context.Background(), "")

	t.Assert().ErrorIs(err, ports.ErrUserPhoneNotFound)
}

func (s *UserGetByPhoneSuite) SuccessRepositoryMock(repository *mocks.UserRepository) {
	repository.
		On("GetByPhone", context.Background(), mock.AnythingOfType("string")).
		Return(domain.User{}, nil)
}

func (s *UserGetByPhoneSuite) TestSuccess(t provider.T) {
	t.Parallel()
	t.Title("User get by phone test success")
	repository := mocks.NewUserRepository(t)
	userService := service.NewUserService(repository, s.hashProvider, s.logger)
	s.SuccessRepositoryMock(repository)

	_, err := userService.GetByPhone(context.Background(), "")

	t.Assert().Nil(err)
}

func TestUserGetByPhoneSuite(t *testing.T) {
	suite.RunSuite(t, new(UserGetByPhoneSuite))
}
