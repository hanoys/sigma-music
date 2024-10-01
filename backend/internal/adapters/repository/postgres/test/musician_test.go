package test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
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

type MusicianSuite struct {
	suite.Suite
}

func NewMusicianRepository() (ports.IMusicianRepository, sqlmock.Sqlmock) {
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	conn := sqlx.NewDb(db, "pgx")
	repo := postgres.NewPostgresMusicianRepository(conn)
	return repo, mock
}

type MusicianCreateSuite struct {
	MusicianSuite
}

func (s *MusicianCreateSuite) SuccessRepositoryMock(mock sqlmock.Sqlmock, user domain.Musician) {
	pgMusician := entity.NewPgMusician(user)
	queryString := InsertQueryString(pgMusician, "musicians")
	mock.ExpectExec(queryString).
		WithArgs(EntityValues(pgMusician)...).
		WillReturnResult(sqlmock.NewResult(1, 1))

	expectedRows := sqlmock.NewRows(EntityColumns(pgMusician)).
		AddRow(EntityValues(pgMusician)...)
	mock.ExpectQuery(postgres.MusicianGetByIDQuery).
		WithArgs(pgMusician.ID).
		WillReturnRows(expectedRows)
}

func (s *MusicianCreateSuite) TestSuccess(t provider.T) {
	t.Parallel()
	t.Title("Repository Musician create test success")
	repo, mock := NewMusicianRepository()
	user := builder.NewMusicianBuilder().Default().Build()
	s.SuccessRepositoryMock(mock, user)

	userResult, err := repo.Create(context.Background(), user)

	t.Assert().Nil(err)
	t.Assert().Equal(user, userResult)
}

func (s *MusicianCreateSuite) DuplicateRepositoryMock(mock sqlmock.Sqlmock, user domain.Musician) {
	pgMusician := entity.NewPgMusician(user)
	queryString := InsertQueryString(pgMusician, "musicians")
	mock.ExpectExec(queryString).
		WithArgs(EntityValues(pgMusician)...).
		WillReturnError(&pgconn.PgError{Code: pgerrcode.UniqueViolation})
}

func (s *MusicianCreateSuite) TestDuplicate(t provider.T) {
	t.Parallel()
	repo, mock := NewMusicianRepository()
	user := builder.NewMusicianBuilder().Default().Build()
	s.DuplicateRepositoryMock(mock, user)

	userResult, err := repo.Create(context.Background(), user)

	t.Assert().ErrorIs(err, ports.ErrMusicianDuplicate)
	t.Assert().Equal(domain.Musician{}, userResult)
}

func TestMusicianCreateSuite(t *testing.T) {
	suite.RunNamedSuite(t, "MusicianCreateRepository", new(MusicianCreateSuite))
}

type MusicianGetAllSuite struct {
	MusicianSuite
}

func (s *MusicianGetAllSuite) SuccessRepositoryMock(mock sqlmock.Sqlmock, user domain.Musician) {
	pgMusician := entity.NewPgMusician(user)
	expectedRows := sqlmock.NewRows(EntityColumns(pgMusician)).
		AddRow(EntityValues(pgMusician)...)
	mock.ExpectQuery(postgres.MusicianGetAllQuery).
		WillReturnRows(expectedRows)
}

func (s *MusicianGetAllSuite) TestSuccess(t provider.T) {
	t.Parallel()
	repo, mock := NewMusicianRepository()
	user := builder.NewMusicianBuilder().Default().Build()
	s.SuccessRepositoryMock(mock, user)

	users, err := repo.GetAll(context.Background())

	t.Assert().Nil(err)
	t.Assert().Equal(user, users[0])
}

func (s *MusicianGetAllSuite) InternalErrorRepositoryMock(mock sqlmock.Sqlmock, album domain.Musician) {
	mock.ExpectQuery(postgres.MusicianGetAllQuery).
		WillReturnError(sql.ErrNoRows)
}

func (s *MusicianGetAllSuite) TestInternalError(t provider.T) {
	t.Parallel()
	repo, mock := NewMusicianRepository()
	album := builder.NewMusicianBuilder().Default().Build()
	s.InternalErrorRepositoryMock(mock, album)

	albums, err := repo.GetAll(context.Background())

	t.Assert().Nil(albums)
	t.Assert().ErrorIs(err, ports.ErrInternalMusicianRepo)
}

func TestMusicianGetAllSuite(t *testing.T) {
	suite.RunNamedSuite(t, "MusicianGetAllRepository", new(MusicianGetAllSuite))
}

type MusicianGetByIDSuite struct {
	MusicianSuite
}

func (s *MusicianGetByIDSuite) SuccessRepositoryMock(mock sqlmock.Sqlmock, user domain.Musician) {
	pgMusician := entity.NewPgMusician(user)
	expectedRows := sqlmock.NewRows(EntityColumns(pgMusician)).
		AddRow(EntityValues(pgMusician)...)
	mock.ExpectQuery(postgres.MusicianGetByIDQuery).
		WithArgs(user.ID).
		WillReturnRows(expectedRows)
}

func (s *MusicianGetByIDSuite) TestSuccess(t provider.T) {
	t.Parallel()
	repo, mock := NewMusicianRepository()
	user := builder.NewMusicianBuilder().Default().Build()
	s.SuccessRepositoryMock(mock, user)

	resultMusician, err := repo.GetByID(context.Background(), user.ID)

	t.Assert().Nil(err)
	t.Assert().Equal(user, resultMusician)
}

func (s *MusicianGetByIDSuite) NotFoundRepositoryMock(mock sqlmock.Sqlmock, user domain.Musician) {
	mock.ExpectQuery(postgres.MusicianGetByIDQuery).
		WillReturnError(sql.ErrNoRows)
}

func (s *MusicianGetByIDSuite) TestNotFound(t provider.T) {
	t.Parallel()
	repo, mock := NewMusicianRepository()
	user := builder.NewMusicianBuilder().Default().Build()
	s.NotFoundRepositoryMock(mock, user)

	resultMusician, err := repo.GetByID(context.Background(), user.ID)

	t.Assert().ErrorIs(err, ports.ErrMusicianIDNotFound)
	t.Assert().Equal(resultMusician, domain.Musician{})
}

func TestMusicianGetByIDSuite(t *testing.T) {
	suite.RunNamedSuite(t, "MusicianGetByIDRepository", new(MusicianGetByIDSuite))
}

type MusicianGetByNameSuite struct {
	MusicianSuite
}

func (s *MusicianGetByNameSuite) SuccessRepositoryMock(mock sqlmock.Sqlmock, user domain.Musician) {
	pgMusician := entity.NewPgMusician(user)
	expectedRows := sqlmock.NewRows(EntityColumns(pgMusician)).
		AddRow(EntityValues(pgMusician)...)
	mock.ExpectQuery(postgres.MusicianGetByNameQuery).
		WithArgs(user.Name).
		WillReturnRows(expectedRows)
}

func (s *MusicianGetByNameSuite) TestSuccess(t provider.T) {
	t.Parallel()
	repo, mock := NewMusicianRepository()
	user := builder.NewMusicianBuilder().Default().Build()
	s.SuccessRepositoryMock(mock, user)

	resultMusician, err := repo.GetByName(context.Background(), user.Name)

	t.Assert().Nil(err)
	t.Assert().Equal(user, resultMusician)
}

func (s *MusicianGetByNameSuite) NotFoundRepositoryMock(mock sqlmock.Sqlmock, user domain.Musician) {
	mock.ExpectQuery(postgres.MusicianGetByNameQuery).
		WillReturnError(sql.ErrNoRows)
}

func (s *MusicianGetByNameSuite) TestNotFound(t provider.T) {
	t.Parallel()
	repo, mock := NewMusicianRepository()
	user := builder.NewMusicianBuilder().Default().Build()
	s.NotFoundRepositoryMock(mock, user)

	resultMusician, err := repo.GetByName(context.Background(), user.Name)

	t.Assert().ErrorIs(err, ports.ErrMusicianNameNotFound)
	t.Assert().Equal(resultMusician, domain.Musician{})
}

func TestMusicianGetByNameSuite(t *testing.T) {
	suite.RunNamedSuite(t, "MusicianGetByNameRepository", new(MusicianGetByNameSuite))
}

type MusicianGetByEmailSuite struct {
	MusicianSuite
}

func (s *MusicianGetByEmailSuite) SuccessRepositoryMock(mock sqlmock.Sqlmock, user domain.Musician) {
	pgMusician := entity.NewPgMusician(user)
	expectedRows := sqlmock.NewRows(EntityColumns(pgMusician)).
		AddRow(EntityValues(pgMusician)...)
	mock.ExpectQuery(postgres.MusicianGetByEmailQuery).
		WithArgs(user.Email).
		WillReturnRows(expectedRows)
}

func (s *MusicianGetByEmailSuite) TestSuccess(t provider.T) {
	t.Parallel()
	repo, mock := NewMusicianRepository()
	user := builder.NewMusicianBuilder().Default().Build()
	s.SuccessRepositoryMock(mock, user)

	resultMusician, err := repo.GetByEmail(context.Background(), user.Email)

	t.Assert().Nil(err)
	t.Assert().Equal(user, resultMusician)
}

func (s *MusicianGetByEmailSuite) NotFoundRepositoryMock(mock sqlmock.Sqlmock, user domain.Musician) {
	mock.ExpectQuery(postgres.MusicianGetByEmailQuery).
		WillReturnError(sql.ErrNoRows)
}

func (s *MusicianGetByEmailSuite) TestNotFound(t provider.T) {
	t.Parallel()
	repo, mock := NewMusicianRepository()
	user := builder.NewMusicianBuilder().Default().Build()
	s.NotFoundRepositoryMock(mock, user)

	resultMusician, err := repo.GetByEmail(context.Background(), user.Email)

	t.Assert().ErrorIs(err, ports.ErrMusicianEmailNotFound)
	t.Assert().Equal(resultMusician, domain.Musician{})
}

func TestMusicianGetByEmailSuite(t *testing.T) {
	suite.RunNamedSuite(t, "MusicianGetByEmailRepository", new(MusicianGetByEmailSuite))
}

type MusicianGetByAlbumIDSuite struct {
	MusicianSuite
}

func (s *MusicianGetByAlbumIDSuite) SuccessRepositoryMock(mock sqlmock.Sqlmock, user domain.Musician, albumID uuid.UUID) {
	pgMusician := entity.NewPgMusician(user)
	expectedRows := sqlmock.NewRows(EntityColumns(pgMusician)).
		AddRow(EntityValues(pgMusician)...)
	mock.ExpectQuery(postgres.MusicianGetByAlbumIDQuery).
		WithArgs(albumID).
		WillReturnRows(expectedRows)
}

func (s *MusicianGetByAlbumIDSuite) TestSuccess(t provider.T) {
	t.Parallel()
	repo, mock := NewMusicianRepository()
	user := builder.NewMusicianBuilder().Default().Build()
	albumID := uuid.New()
	s.SuccessRepositoryMock(mock, user, albumID)

	resultMusician, err := repo.GetByAlbumID(context.Background(), albumID)

	t.Assert().Nil(err)
	t.Assert().Equal(user, resultMusician)
}

func (s *MusicianGetByAlbumIDSuite) NotFoundRepositoryMock(mock sqlmock.Sqlmock, user domain.Musician) {
	mock.ExpectQuery(postgres.MusicianGetByAlbumIDQuery).
		WillReturnError(sql.ErrNoRows)
}

func (s *MusicianGetByAlbumIDSuite) TestNotFound(t provider.T) {
	t.Parallel()
	repo, mock := NewMusicianRepository()
	user := builder.NewMusicianBuilder().Default().Build()
	albumID := uuid.New()
	s.NotFoundRepositoryMock(mock, user)

	resultMusician, err := repo.GetByAlbumID(context.Background(), albumID)

	t.Assert().ErrorIs(err, ports.ErrMusicianIDNotFound)
	t.Assert().Equal(resultMusician, domain.Musician{})
}

func TestMusicianGetByAlbumIDSuite(t *testing.T) {
	suite.RunNamedSuite(t, "MusicianGetByAlbumIDRepository", new(MusicianGetByAlbumIDSuite))
}

type MusicianGetByTrackIDSuite struct {
	MusicianSuite
}

func (s *MusicianGetByTrackIDSuite) SuccessRepositoryMock(mock sqlmock.Sqlmock, user domain.Musician, albumID uuid.UUID) {
	pgMusician := entity.NewPgMusician(user)
	expectedRows := sqlmock.NewRows(EntityColumns(pgMusician)).
		AddRow(EntityValues(pgMusician)...)
	mock.ExpectQuery(postgres.MusicianGetByTrackIDQuery).
		WithArgs(albumID).
		WillReturnRows(expectedRows)
}

func (s *MusicianGetByTrackIDSuite) TestSuccess(t provider.T) {
	t.Parallel()
	repo, mock := NewMusicianRepository()
	user := builder.NewMusicianBuilder().Default().Build()
	albumID := uuid.New()
	s.SuccessRepositoryMock(mock, user, albumID)

	resultMusician, err := repo.GetByTrackID(context.Background(), albumID)

	t.Assert().Nil(err)
	t.Assert().Equal(user, resultMusician)
}

func (s *MusicianGetByTrackIDSuite) NotFoundRepositoryMock(mock sqlmock.Sqlmock, user domain.Musician) {
	mock.ExpectQuery(postgres.MusicianGetByTrackIDQuery).
		WillReturnError(sql.ErrNoRows)
}

func (s *MusicianGetByTrackIDSuite) TestNotFound(t provider.T) {
	t.Parallel()
	repo, mock := NewMusicianRepository()
	user := builder.NewMusicianBuilder().Default().Build()
	albumID := uuid.New()
	s.NotFoundRepositoryMock(mock, user)

	resultMusician, err := repo.GetByTrackID(context.Background(), albumID)

	t.Assert().ErrorIs(err, ports.ErrMusicianIDNotFound)
	t.Assert().Equal(resultMusician, domain.Musician{})
}

func TestMusicianGetByTrackIDSuite(t *testing.T) {
	suite.RunNamedSuite(t, "MusicianGetByTrackIDRepository", new(MusicianGetByTrackIDSuite))
}
