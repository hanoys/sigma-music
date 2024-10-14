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
	"github.com/jmoiron/sqlx"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type GenreSuite struct {
	suite.Suite
}

func NewGenreRepository() (ports.IGenreRepository, sqlmock.Sqlmock) {
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	conn := sqlx.NewDb(db, "pgx")
	repo := postgres.NewPostgresGenreRepository(conn)
	return repo, mock
}

type GenreGetAllSuite struct {
	GenreSuite
}

func (s *GenreGetAllSuite) SuccessRepositoryMock(mock sqlmock.Sqlmock, genre domain.Genre) {
	pgGenre := entity.NewPgGenre(genre)
	expectedRows := sqlmock.NewRows(EntityColumns(pgGenre)).
		AddRow(EntityValues(pgGenre)...)
	mock.ExpectQuery(postgres.GenreGetAllQuery).
		WillReturnRows(expectedRows)
}

func (s *GenreGetAllSuite) TestSuccess(t provider.T) {
	t.Parallel()
	t.Title("Repository Genre get all test success")
	repo, mock := NewGenreRepository()
	genre := builder.NewGenreBuilder().Default().Build()
	s.SuccessRepositoryMock(mock, genre)

	genres, err := repo.GetAll(context.Background())

	t.Assert().Nil(err)
	t.Assert().Equal(genre, genres[0])
}

func (s *GenreGetAllSuite) InternalErrorRepositoryMock(mock sqlmock.Sqlmock, album domain.Album) {
	mock.ExpectQuery(postgres.GenreGetAllQuery).
		WillReturnError(sql.ErrNoRows)
}

func (s *GenreGetAllSuite) TestInternalError(t provider.T) {
	t.Parallel()
	t.Title("Repository Genre get all test internal error")
	repo, mock := NewGenreRepository()
	genre := builder.NewAlbumBuilder().Default().Build()
	s.InternalErrorRepositoryMock(mock, genre)

	genres, err := repo.GetAll(context.Background())

	t.Assert().Nil(genres)
	t.Assert().ErrorIs(err, ports.ErrInternalGenreRepo)
}

func TestGenreGetAllSuite(t *testing.T) {
	suite.RunNamedSuite(t, "GenreGetAllRepository", new(GenreGetAllSuite))
}

type GenreGetByIDSuite struct {
	GenreSuite
}

func (s *GenreGetByIDSuite) SuccessRepositoryMock(mock sqlmock.Sqlmock, genre domain.Genre) {
	pgGenre := entity.NewPgGenre(genre)
	expectedRows := sqlmock.NewRows(EntityColumns(pgGenre)).
		AddRow(EntityValues(pgGenre)...)
	mock.ExpectQuery(postgres.GenreGetByIDQuery).
		WithArgs(genre.ID).
		WillReturnRows(expectedRows)
}

func (s *GenreGetByIDSuite) TestSuccess(t provider.T) {
	t.Parallel()
	t.Title("Repository Genre get by id test success")
	repo, mock := NewGenreRepository()
	genre := builder.NewGenreBuilder().Default().Build()
	s.SuccessRepositoryMock(mock, genre)

	resultGenre, err := repo.GetByID(context.Background(), genre.ID)

	t.Assert().Nil(err)
	t.Assert().Equal(genre, resultGenre)
}

func (s *GenreGetByIDSuite) InternalErrorRepositoryMock(mock sqlmock.Sqlmock, genre domain.Genre) {
	mock.ExpectQuery(postgres.GenreGetByIDQuery).
		WillReturnError(sql.ErrNoRows)
}

func (s *GenreGetByIDSuite) TestInternalError(t provider.T) {
	t.Parallel()
	t.Title("Repository Genre get by id test internal error")
	repo, mock := NewGenreRepository()
	genre := builder.NewGenreBuilder().Default().Build()
	s.InternalErrorRepositoryMock(mock, genre)

	_, err := repo.GetByID(context.Background(), genre.ID)

	t.Assert().ErrorIs(err, ports.ErrGenreIDNotFound)
}

func TestGenreGetByIDSuite(t *testing.T) {
	suite.RunNamedSuite(t, "GenreGetByIDRepository", new(GenreGetByIDSuite))
}

type GenreGetByTrackIDSuite struct {
	GenreSuite
}

func (s *GenreGetByTrackIDSuite) SuccessRepositoryMock(mock sqlmock.Sqlmock, genre domain.Genre, trackID uuid.UUID) {
	pgGenre := entity.NewPgGenre(genre)
	expectedRows := sqlmock.NewRows(EntityColumns(pgGenre)).
		AddRow(EntityValues(pgGenre)...)
	mock.ExpectQuery(postgres.GenreGetByTrack).
		WithArgs(trackID).
		WillReturnRows(expectedRows)
}

func (s *GenreGetByTrackIDSuite) TestSuccess(t provider.T) {
	t.Parallel()
	t.Title("Repository Genre get by track id test success")
	repo, mock := NewGenreRepository()
	genre := builder.NewGenreBuilder().Default().Build()
	trackID := uuid.New()
	s.SuccessRepositoryMock(mock, genre, trackID)

	genres, err := repo.GetByTrackID(context.Background(), trackID)

	t.Assert().Nil(err)
	t.Assert().Equal(genre, genres[0])
}

func (s *GenreGetByTrackIDSuite) InternalErrorRepositoryMock(mock sqlmock.Sqlmock, album domain.Album, trackID uuid.UUID) {
	mock.ExpectQuery(postgres.GenreGetByTrack).
		WillReturnError(sql.ErrNoRows)
}

func (s *GenreGetByTrackIDSuite) TestInternalError(t provider.T) {
	t.Parallel()
	t.Title("Repository Genre get by track id test internal error")
	repo, mock := NewGenreRepository()
	genre := builder.NewAlbumBuilder().Default().Build()
	trackID := uuid.New()
	s.InternalErrorRepositoryMock(mock, genre, trackID)

	genres, err := repo.GetByTrackID(context.Background(), genre.ID)

	t.Assert().Nil(genres)
	t.Assert().ErrorIs(err, ports.ErrInternalGenreRepo)
}

func TestGenreGetByTrackIDSuite(t *testing.T) {
	suite.RunNamedSuite(t, "GenreGetByTrackIDRepository", new(GenreGetByTrackIDSuite))
}
