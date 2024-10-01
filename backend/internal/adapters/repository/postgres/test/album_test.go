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

type AlbumSuite struct {
	suite.Suite
}

func NewAlbumRepository() (ports.IAlbumRepository, sqlmock.Sqlmock) {
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	conn := sqlx.NewDb(db, "pgx")
	repo := postgres.NewPostgresAlbumRepository(conn)
	return repo, mock
}

type AlbumCreateSuite struct {
	AlbumSuite
}

func (s *AlbumCreateSuite) SuccessRepositoryMock(mock sqlmock.Sqlmock, album domain.Album) {
	pgAlbum := entity.NewPgAlbum(album)
	queryString := InsertQueryString(pgAlbum, "albums")
	mock.ExpectExec(queryString).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(postgres.AlbumInsertQuery).
		WillReturnResult(sqlmock.NewResult(1, 1))

	expectedRows := sqlmock.NewRows(EntityColumns(pgAlbum)).
		AddRow(EntityValues(pgAlbum)...)
	mock.ExpectQuery(postgres.AlbumGetByIDInternalQuery).
		WithArgs(pgAlbum.ID).
		WillReturnRows(expectedRows)
}

func (s *AlbumCreateSuite) TestSuccess(t provider.T) {
	t.Parallel()
	t.Title("Repository Album create test success")
	repo, mock := NewAlbumRepository()
	album := builder.NewAlbumBuilder().Default().Build()
	musicianID := uuid.New()
	s.SuccessRepositoryMock(mock, album)

	albumResult, err := repo.Create(context.Background(), album, musicianID)

	t.Assert().Nil(err)
	t.Assert().Equal(album, albumResult)
}

func (s *AlbumCreateSuite) DuplicateRepositoryMock(mock sqlmock.Sqlmock, user domain.Album) {
	pgAlbum := entity.NewPgAlbum(user)
	queryString := InsertQueryString(pgAlbum, "albums")
	mock.ExpectExec(queryString).
		WillReturnError(&pgconn.PgError{Code: pgerrcode.UniqueViolation})
}

func (s *AlbumCreateSuite) TestDuplicate(t provider.T) {
	t.Parallel()
	t.Title("Repository Album create test duplicate")
	repo, mock := NewAlbumRepository()
	album := builder.NewAlbumBuilder().Default().Build()
	musicianID := uuid.New()
	s.DuplicateRepositoryMock(mock, album)

	albumResult, err := repo.Create(context.Background(), album, musicianID)

	t.Assert().ErrorIs(err, ports.ErrAlbumDuplicate)
	t.Assert().Equal(domain.Album{}, albumResult)
}

func TestAlbumCreateSuite(t *testing.T) {
	suite.RunNamedSuite(t, "AlbumCreateRepository", new(AlbumCreateSuite))
}

type AlbumGetAllSuite struct {
	AlbumSuite
}

func (s *AlbumGetAllSuite) SuccessRepositoryMock(mock sqlmock.Sqlmock, album domain.Album) {
	pgAlbum := entity.NewPgAlbum(album)
	expectedRows := sqlmock.NewRows(EntityColumns(pgAlbum)).
		AddRow(EntityValues(pgAlbum)...)
	mock.ExpectQuery(postgres.AlbumGetAllQuery).
		WillReturnRows(expectedRows)
}

func (s *AlbumGetAllSuite) TestSuccess(t provider.T) {
	t.Parallel()
	t.Title("Repository Album get all test success")
	repo, mock := NewAlbumRepository()
	album := builder.NewAlbumBuilder().Default().Build()
	s.SuccessRepositoryMock(mock, album)

	albums, err := repo.GetAll(context.Background())

	t.Assert().Nil(err)
	t.Assert().Equal(album, albums[0])
}

func (s *AlbumGetAllSuite) InternalErrorRepositoryMock(mock sqlmock.Sqlmock, album domain.Album) {
	mock.ExpectQuery(postgres.AlbumGetAllQuery).
		WillReturnError(sql.ErrNoRows)
}

func (s *AlbumGetAllSuite) TestInternalError(t provider.T) {
	t.Parallel()
	t.Title("Repository Album get all test internal error")
	repo, mock := NewAlbumRepository()
	album := builder.NewAlbumBuilder().Default().Build()
	s.InternalErrorRepositoryMock(mock, album)

	albums, err := repo.GetAll(context.Background())

	t.Assert().Nil(albums)
	t.Assert().ErrorIs(err, ports.ErrInternalAlbumRepo)
}

func TestAlbumGetAllSuite(t *testing.T) {
	suite.RunNamedSuite(t, "AlbumGetAllRepository", new(AlbumGetAllSuite))
}

type AlbumGetByMusicianIDSuite struct {
	AlbumSuite
}

func (s *AlbumGetByMusicianIDSuite) SuccessRepositoryMock(mock sqlmock.Sqlmock, album domain.Album, musicianID uuid.UUID) {
	pgAlbum := entity.NewPgAlbum(album)
	expectedRows := sqlmock.NewRows(EntityColumns(pgAlbum)).
		AddRow(EntityValues(pgAlbum)...)
	mock.ExpectQuery(postgres.AlbumGetByMusicianIDQuery).
		WillReturnRows(expectedRows)
}

func (s *AlbumGetByMusicianIDSuite) TestSuccess(t provider.T) {
	t.Parallel()
	t.Title("Repository Album get by musician id test success")
	repo, mock := NewAlbumRepository()
	album := builder.NewAlbumBuilder().Default().Build()
	musicianID := uuid.New()
	s.SuccessRepositoryMock(mock, album, musicianID)

	albums, err := repo.GetByMusicianID(context.Background(), musicianID)

	t.Assert().Nil(err)
	t.Assert().Equal(album, albums[0])
}

func (s *AlbumGetByMusicianIDSuite) InternalErrorRepositoryMock(mock sqlmock.Sqlmock, album domain.Album, musicianID uuid.UUID) {
	mock.ExpectQuery(postgres.AlbumGetByMusicianIDQuery).
		WillReturnError(sql.ErrNoRows)
}

func (s *AlbumGetByMusicianIDSuite) TestInternalError(t provider.T) {
	t.Parallel()
	t.Title("Repository Album get by musician id test internal error")
	repo, mock := NewAlbumRepository()
	album := builder.NewAlbumBuilder().Default().Build()
	musicianID := uuid.New()
	s.InternalErrorRepositoryMock(mock, album, musicianID)

	albums, err := repo.GetByMusicianID(context.Background(), musicianID)

	t.Assert().ErrorIs(err, ports.ErrInternalAlbumRepo)
	t.Assert().Nil(albums)
}

func TestAlbumGetByMusicianIDSuite(t *testing.T) {
	suite.RunNamedSuite(t, "AlbumGetGetByMusicianIDRepository", new(AlbumGetByMusicianIDSuite))
}

type AlbumGetOwnSuite struct {
	AlbumSuite
}

func (s *AlbumGetOwnSuite) SuccessRepositoryMock(mock sqlmock.Sqlmock, album domain.Album, musicianID uuid.UUID) {
	pgAlbum := entity.NewPgAlbum(album)
	expectedRows := sqlmock.NewRows(EntityColumns(pgAlbum)).
		AddRow(EntityValues(pgAlbum)...)
	mock.ExpectQuery(postgres.AlbumGetOwnQuery).
		WillReturnRows(expectedRows)
}

func (s *AlbumGetOwnSuite) TestSuccess(t provider.T) {
	t.Parallel()
	t.Title("Repository Album get own test success")
	repo, mock := NewAlbumRepository()
	album := builder.NewAlbumBuilder().Default().Build()
	musicianID := uuid.New()
	s.SuccessRepositoryMock(mock, album, musicianID)

	albums, err := repo.GetOwn(context.Background(), musicianID)

	t.Assert().Nil(err)
	t.Assert().Equal(album, albums[0])
}

func (s *AlbumGetOwnSuite) InternalErrorRepositoryMock(mock sqlmock.Sqlmock, album domain.Album, musicianID uuid.UUID) {
	mock.ExpectQuery(postgres.AlbumGetOwnQuery).
		WillReturnError(sql.ErrNoRows)
}

func (s *AlbumGetOwnSuite) TestInternalError(t provider.T) {
	t.Parallel()
	t.Title("Repository Album get own test internal error")
	repo, mock := NewAlbumRepository()
	album := builder.NewAlbumBuilder().Default().Build()
	musicianID := uuid.New()
	s.InternalErrorRepositoryMock(mock, album, musicianID)

	albums, err := repo.GetByMusicianID(context.Background(), musicianID)

	t.Assert().ErrorIs(err, ports.ErrInternalAlbumRepo)
	t.Assert().Nil(albums)
}

func TestAlbumGetOwnSuite(t *testing.T) {
	suite.RunNamedSuite(t, "AlbumGetOwnRepository", new(AlbumGetOwnSuite))
}

type AlbumGetByIDSuite struct {
	AlbumSuite
}

func (s *AlbumGetByIDSuite) SuccessRepositoryMock(mock sqlmock.Sqlmock, album domain.Album) {
	pgAlbum := entity.NewPgAlbum(album)
	expectedRows := sqlmock.NewRows(EntityColumns(pgAlbum)).
		AddRow(EntityValues(pgAlbum)...)
	mock.ExpectQuery(postgres.AlbumGetByIDQuery).
		WithArgs(album.ID).
		WillReturnRows(expectedRows)
}

func (s *AlbumGetByIDSuite) TestSuccess(t provider.T) {
	t.Parallel()
	t.Title("Repository Album get by id test success")
	repo, mock := NewAlbumRepository()
	album := builder.NewAlbumBuilder().Default().Build()
	s.SuccessRepositoryMock(mock, album)

	resultAlbum, err := repo.GetByID(context.Background(), album.ID)

	t.Assert().Nil(err)
	t.Assert().Equal(album, resultAlbum)
}

func (s *AlbumGetByIDSuite) NotFoundRepositoryMock(mock sqlmock.Sqlmock, album domain.Album) {
	mock.ExpectQuery(postgres.AlbumGetByIDQuery).
		WillReturnError(sql.ErrNoRows)
}

func (s *AlbumGetByIDSuite) TestNotFound(t provider.T) {
	t.Parallel()
	t.Title("Repository Album get by id test not found")
	repo, mock := NewAlbumRepository()
	album := builder.NewAlbumBuilder().Default().Build()
	s.NotFoundRepositoryMock(mock, album)

	resultAlbum, err := repo.GetByID(context.Background(), album.ID)

	t.Assert().ErrorIs(err, ports.ErrAlbumIDNotFound)
	t.Assert().Equal(resultAlbum, domain.Album{})
}

func TestAlbumGetByIDSuite(t *testing.T) {
	suite.RunNamedSuite(t, "AlbumGetByIDRepository", new(AlbumGetByIDSuite))
}

type AlbumPublishSuite struct {
	AlbumSuite
}

func (s *AlbumPublishSuite) SuccessRepositoryMock(mock sqlmock.Sqlmock, album domain.Album) {
	pgAlbum := entity.NewPgAlbum(album)
	expectedRows := sqlmock.NewRows(EntityColumns(pgAlbum)).
		AddRow(EntityValues(pgAlbum)...)
	mock.ExpectQuery(postgres.AlbumGetByIDInternalQuery).
		WillReturnRows(expectedRows)

	queryString := UpdateQueryString(pgAlbum, "albums")
	mock.ExpectExec(queryString).
		WillReturnResult(sqlmock.NewResult(1, 1))
}

func (s *AlbumPublishSuite) TestSuccess(t provider.T) {
	t.Parallel()
	t.Title("Repository Album publish test success")
	repo, mock := NewAlbumRepository()
	album := builder.NewAlbumBuilder().Default().Build()
	s.SuccessRepositoryMock(mock, album)

	err := repo.Publish(context.Background(), album.ID)

	t.Assert().Nil(err)
}

func (s *AlbumPublishSuite) NotFoundRepositoryMock(mock sqlmock.Sqlmock, album domain.Album) {
	mock.ExpectQuery(postgres.AlbumGetByIDInternalQuery).
		WillReturnError(sql.ErrNoRows)
}

func (s *AlbumPublishSuite) TestNotFound(t provider.T) {
	t.Parallel()
	t.Title("Repository Album publish test not found")
	repo, mock := NewAlbumRepository()
	album := builder.NewAlbumBuilder().Default().Build()
	s.NotFoundRepositoryMock(mock, album)

	err := repo.Publish(context.Background(), album.ID)

	t.Assert().ErrorIs(err, ports.ErrAlbumIDNotFound)
}

func TestAlbumPublishSuite(t *testing.T) {
	suite.RunNamedSuite(t, "AlbumPublishRepository", new(AlbumPublishSuite))
}
