package test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/hanoys/sigma-music/internal/adapters/repository/postgres"
	"github.com/hanoys/sigma-music/internal/adapters/repository/postgres/entity"
	"github.com/hanoys/sigma-music/internal/adapters/repository/postgres/test/builder"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type UserSuite struct {
	suite.Suite
}

func NewUserRepository() (ports.IUserRepository, sqlmock.Sqlmock) {
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	conn := sqlx.NewDb(db, "pgx")
	repo := postgres.NewPostgresUserRepository(conn)
	return repo, mock
}

type UserCreateSuite struct {
	UserSuite
}

func (s *UserCreateSuite) SuccessRepositoryMock(mock sqlmock.Sqlmock, user domain.User) {
	pgUser := entity.NewPgUser(user)
	queryString := InsertQueryString(pgUser, "users")
	mock.ExpectExec(queryString).
		WithArgs(EntityValues(pgUser)...).
		WillReturnResult(sqlmock.NewResult(1, 1))

	expectedRows := sqlmock.NewRows(EntityColumns(pgUser)).
		AddRow(EntityValues(pgUser)...)
	mock.ExpectQuery(postgres.UserGetByIDQuery).
		WithArgs(pgUser.ID).
		WillReturnRows(expectedRows)
}

func (s *UserCreateSuite) TestSuccess(t provider.T) {
	t.Parallel()
	repo, mock := NewUserRepository()
	user := builder.NewUserBuilder().Default().Build()
	s.SuccessRepositoryMock(mock, user)

	userResult, err := repo.Create(context.Background(), user)

	t.Assert().Nil(err)
	t.Assert().Equal(user, userResult)
}

func (s *UserCreateSuite) DuplicateRepositoryMock(mock sqlmock.Sqlmock, user domain.User) {
	pgUser := entity.NewPgUser(user)
	queryString := InsertQueryString(pgUser, "users")
	mock.ExpectExec(queryString).
		WithArgs(EntityValues(pgUser)...).
		WillReturnError(&pgconn.PgError{Code: pgerrcode.UniqueViolation})
}

func (s *UserCreateSuite) TestDuplicate(t provider.T) {
	t.Parallel()
	repo, mock := NewUserRepository()
	user := builder.NewUserBuilder().Default().Build()
	s.DuplicateRepositoryMock(mock, user)

	userResult, err := repo.Create(context.Background(), user)

	t.Assert().ErrorIs(err, ports.ErrUserDuplicate)
	t.Assert().Equal(domain.User{}, userResult)
}

func TestUserCreateSuite(t *testing.T) {
	suite.RunNamedSuite(t, "UserCreateRepository", new(UserCreateSuite))
}

type UserGetAllSuite struct {
	UserSuite
}

func (s *UserGetAllSuite) SuccessRepositoryMock(mock sqlmock.Sqlmock, user domain.User) {
	pgUser := entity.NewPgUser(user)
	expectedRows := sqlmock.NewRows(EntityColumns(pgUser)).
		AddRow(EntityValues(pgUser)...)
	mock.ExpectQuery(postgres.UserGetAllQuery).
		WillReturnRows(expectedRows)
}

func (s *UserGetAllSuite) TestSuccess(t provider.T) {
	t.Parallel()
	t.Title("Repository User get all test success")
	repo, mock := NewUserRepository()
	user := builder.NewUserBuilder().Default().Build()
	s.SuccessRepositoryMock(mock, user)

	users, err := repo.GetAll(context.Background())

	t.Assert().Nil(err)
	t.Assert().Equal(user, users[0])
}

func (s *UserGetAllSuite) InternalErrorRepositoryMock(mock sqlmock.Sqlmock, album domain.User) {
	mock.ExpectQuery(postgres.UserGetAllQuery).
		WillReturnError(sql.ErrNoRows)
}

func (s *UserGetAllSuite) TestInternalError(t provider.T) {
	t.Parallel()
	t.Title("Repository User get all test internal error")
	repo, mock := NewUserRepository()
	album := builder.NewUserBuilder().Default().Build()
	s.InternalErrorRepositoryMock(mock, album)

	albums, err := repo.GetAll(context.Background())

	t.Assert().Nil(albums)
	t.Assert().ErrorIs(err, ports.ErrInternalUserRepo)
}

func TestUserGetAllSuite(t *testing.T) {
	suite.RunNamedSuite(t, "UserGetAllRepository", new(UserGetAllSuite))
}

type UserGetByIDSuite struct {
	UserSuite
}

func (s *UserGetByIDSuite) SuccessRepositoryMock(mock sqlmock.Sqlmock, user domain.User) {
	pgUser := entity.NewPgUser(user)
	expectedRows := sqlmock.NewRows(EntityColumns(pgUser)).
		AddRow(EntityValues(pgUser)...)
	mock.ExpectQuery(postgres.UserGetByIDQuery).
		WithArgs(user.ID).
		WillReturnRows(expectedRows)
}

func (s *UserGetByIDSuite) TestSuccess(t provider.T) {
	t.Parallel()
	repo, mock := NewUserRepository()
	user := builder.NewUserBuilder().Default().Build()
	s.SuccessRepositoryMock(mock, user)

	resultUser, err := repo.GetByID(context.Background(), user.ID)

	t.Assert().Nil(err)
	t.Assert().Equal(user, resultUser)
}

func (s *UserGetByIDSuite) NotFoundRepositoryMock(mock sqlmock.Sqlmock, user domain.User) {
	mock.ExpectQuery(postgres.UserGetByIDQuery).
		WillReturnError(sql.ErrNoRows)
}

func (s *UserGetByIDSuite) TestNotFound(t provider.T) {
	t.Parallel()
	repo, mock := NewUserRepository()
	user := builder.NewUserBuilder().Default().Build()
	s.NotFoundRepositoryMock(mock, user)

	resultUser, err := repo.GetByID(context.Background(), user.ID)

	t.Assert().ErrorIs(err, ports.ErrUserIDNotFound)
	t.Assert().Equal(resultUser, domain.User{})
}

func TestUserGetByIDSuite(t *testing.T) {
	suite.RunNamedSuite(t, "UserGetByIDRepository", new(UserGetByIDSuite))
}

type UserGetByNameSuite struct {
	UserSuite
}

func (s *UserGetByNameSuite) SuccessRepositoryMock(mock sqlmock.Sqlmock, user domain.User) {
	pgUser := entity.NewPgUser(user)
	expectedRows := sqlmock.NewRows(EntityColumns(pgUser)).
		AddRow(EntityValues(pgUser)...)
	mock.ExpectQuery(postgres.UserGetByNameQuery).
		WithArgs(user.Name).
		WillReturnRows(expectedRows)
}

func (s *UserGetByNameSuite) TestSuccess(t provider.T) {
	t.Parallel()
	repo, mock := NewUserRepository()
	user := builder.NewUserBuilder().Default().Build()
	s.SuccessRepositoryMock(mock, user)

	resultUser, err := repo.GetByName(context.Background(), user.Name)

	t.Assert().Nil(err)
	t.Assert().Equal(user, resultUser)
}

func (s *UserGetByNameSuite) NotFoundRepositoryMock(mock sqlmock.Sqlmock, user domain.User) {
	mock.ExpectQuery(postgres.UserGetByNameQuery).
		WillReturnError(sql.ErrNoRows)
}

func (s *UserGetByNameSuite) TestNotFound(t provider.T) {
	t.Parallel()
	repo, mock := NewUserRepository()
	user := builder.NewUserBuilder().Default().Build()
	s.NotFoundRepositoryMock(mock, user)

	resultUser, err := repo.GetByName(context.Background(), user.Name)

	t.Assert().ErrorIs(err, ports.ErrUserNameNotFound)
	t.Assert().Equal(resultUser, domain.User{})
}

func TestUserGetByNameSuite(t *testing.T) {
	suite.RunNamedSuite(t, "UserGetByNameRepository", new(UserGetByNameSuite))
}

type UserGetByEmailSuite struct {
	UserSuite
}

func (s *UserGetByEmailSuite) SuccessRepositoryMock(mock sqlmock.Sqlmock, user domain.User) {
	pgUser := entity.NewPgUser(user)
	expectedRows := sqlmock.NewRows(EntityColumns(pgUser)).
		AddRow(EntityValues(pgUser)...)
	mock.ExpectQuery(postgres.UserGetByEmailQuery).
		WithArgs(user.Email).
		WillReturnRows(expectedRows)
}

func (s *UserGetByEmailSuite) TestSuccess(t provider.T) {
	t.Parallel()
	repo, mock := NewUserRepository()
	user := builder.NewUserBuilder().Default().Build()
	s.SuccessRepositoryMock(mock, user)

	resultUser, err := repo.GetByEmail(context.Background(), user.Email)

	t.Assert().Nil(err)
	t.Assert().Equal(user, resultUser)
}

func (s *UserGetByEmailSuite) NotFoundRepositoryMock(mock sqlmock.Sqlmock, user domain.User) {
	mock.ExpectQuery(postgres.UserGetByEmailQuery).
		WillReturnError(sql.ErrNoRows)
}

func (s *UserGetByEmailSuite) TestNotFound(t provider.T) {
	t.Parallel()
	repo, mock := NewUserRepository()
	user := builder.NewUserBuilder().Default().Build()
	s.NotFoundRepositoryMock(mock, user)

	resultUser, err := repo.GetByEmail(context.Background(), user.Email)

	t.Assert().ErrorIs(err, ports.ErrUserEmailNotFound)
	t.Assert().Equal(resultUser, domain.User{})
}

func TestUserGetByEmailSuite(t *testing.T) {
	suite.RunNamedSuite(t, "UserGetByEmailRepository", new(UserGetByEmailSuite))
}

type UserGetByPhoneSuite struct {
	UserSuite
}

func (s *UserGetByPhoneSuite) SuccessRepositoryMock(mock sqlmock.Sqlmock, user domain.User) {
	pgUser := entity.NewPgUser(user)
	expectedRows := sqlmock.NewRows(EntityColumns(pgUser)).
		AddRow(EntityValues(pgUser)...)
	mock.ExpectQuery(postgres.UserGetByPhoneQuery).
		WithArgs(user.Phone).
		WillReturnRows(expectedRows)
}

func (s *UserGetByPhoneSuite) TestSuccess(t provider.T) {
	t.Parallel()
	repo, mock := NewUserRepository()
	user := builder.NewUserBuilder().Default().Build()
	s.SuccessRepositoryMock(mock, user)

	resultUser, err := repo.GetByPhone(context.Background(), user.Phone)

	t.Assert().Nil(err)
	t.Assert().Equal(user, resultUser)
}

func (s *UserGetByPhoneSuite) NotFoundRepositoryMock(mock sqlmock.Sqlmock, user domain.User) {
	mock.ExpectQuery(postgres.UserGetByPhoneQuery).
		WillReturnError(sql.ErrNoRows)
}

func (s *UserGetByPhoneSuite) TestNotFound(t provider.T) {
	t.Parallel()
	repo, mock := NewUserRepository()
	user := builder.NewUserBuilder().Default().Build()
	s.NotFoundRepositoryMock(mock, user)

	resultUser, err := repo.GetByPhone(context.Background(), user.Phone)

	t.Assert().ErrorIs(err, ports.ErrUserPhoneNotFound)
	t.Assert().Equal(resultUser, domain.User{})
}

func TestUserGetByPhoneSuite(t *testing.T) {
	suite.RunNamedSuite(t, "UserGetByPhoneRepository", new(UserGetByPhoneSuite))
}
