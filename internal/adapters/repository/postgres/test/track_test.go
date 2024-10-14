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

type TrackSuite struct {
	suite.Suite
}

func NewTrackRepository() (ports.ITrackRepository, sqlmock.Sqlmock) {
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	conn := sqlx.NewDb(db, "pgx")
	repo := postgres.NewPostgresTrackRepository(conn)
	return repo, mock
}

type TrackCreateSuite struct {
	TrackSuite
}

func (s *TrackCreateSuite) SuccessRepositoryMock(mock sqlmock.Sqlmock, track domain.Track) {
	pgTrack := entity.NewPgTrack(track)
	queryString := InsertQueryString(pgTrack, "tracks")
	mock.ExpectExec(queryString).
		WillReturnResult(sqlmock.NewResult(1, 1))

	expectedRows := sqlmock.NewRows(EntityColumns(pgTrack)).
		AddRow(EntityValues(pgTrack)...)
	mock.ExpectQuery(postgres.TrackGetByIDInternalQuery).
		WillReturnRows(expectedRows)
}

func (s *TrackCreateSuite) TestSuccess(t provider.T) {
	t.Parallel()
	repo, mock := NewTrackRepository()
	track := builder.NewTrackBuilder().Default().Build()
	s.SuccessRepositoryMock(mock, track)

	trackResult, err := repo.Create(context.Background(), track)

	t.Assert().Nil(err)
	t.Assert().Equal(track, trackResult)
}

func (s *TrackCreateSuite) DuplicateRepositoryMock(mock sqlmock.Sqlmock, track domain.Track) {
	pgTrack := entity.NewPgTrack(track)
	queryString := InsertQueryString(pgTrack, "tracks")
	mock.ExpectExec(queryString).
		WillReturnError(&pgconn.PgError{Code: pgerrcode.UniqueViolation})
}

func (s *TrackCreateSuite) TestDuplicate(t provider.T) {
	t.Parallel()
	repo, mock := NewTrackRepository()
	track := builder.NewTrackBuilder().Default().Build()
	s.DuplicateRepositoryMock(mock, track)

	trackResult, err := repo.Create(context.Background(), track)

	t.Assert().ErrorIs(err, ports.ErrTrackDuplicate)
	t.Assert().Equal(domain.Track{}, trackResult)
}

func TestTrackCreateSuite(t *testing.T) {
	suite.RunNamedSuite(t, "TrackCreateRepository", new(TrackCreateSuite))
}

type TrackGetAllSuite struct {
	TrackSuite
}

func (s *TrackGetAllSuite) SuccessRepositoryMock(mock sqlmock.Sqlmock, track domain.Track) {
	pgTrack := entity.NewPgTrack(track)
	expectedRows := sqlmock.NewRows(EntityColumns(pgTrack)).
		AddRow(EntityValues(pgTrack)...)
	mock.ExpectQuery(postgres.TrackGetAllQuery).
		WillReturnRows(expectedRows)
}

func (s *TrackGetAllSuite) TestSuccess(t provider.T) {
	t.Parallel()
	repo, mock := NewTrackRepository()
	track := builder.NewTrackBuilder().Default().Build()
	s.SuccessRepositoryMock(mock, track)

	tracks, err := repo.GetAll(context.Background())

	t.Assert().Nil(err)
	t.Assert().Equal(track, tracks[0])
}

func (s *TrackGetAllSuite) InternalErrorRepositoryMock(mock sqlmock.Sqlmock, album domain.Track) {
	mock.ExpectQuery(postgres.TrackGetAllQuery).
		WillReturnError(sql.ErrNoRows)
}

func (s *TrackGetAllSuite) TestInternalError(t provider.T) {
	t.Parallel()
	repo, mock := NewTrackRepository()
	album := builder.NewTrackBuilder().Default().Build()
	s.InternalErrorRepositoryMock(mock, album)

	albums, err := repo.GetAll(context.Background())

	t.Assert().Nil(albums)
	t.Assert().ErrorIs(err, ports.ErrInternalTrackRepo)
}

func TestTrackGetAllSuite(t *testing.T) {
	suite.RunNamedSuite(t, "TrackGetAllRepository", new(TrackGetAllSuite))
}

type TrackGetByIDSuite struct {
	TrackSuite
}

func (s *TrackGetByIDSuite) SuccessRepositoryMock(mock sqlmock.Sqlmock, track domain.Track) {
	pgTrack := entity.NewPgTrack(track)
	expectedRows := sqlmock.NewRows(EntityColumns(pgTrack)).
		AddRow(EntityValues(pgTrack)...)
	mock.ExpectQuery(postgres.TrackGetByIDQuery).
		WithArgs(track.ID).
		WillReturnRows(expectedRows)
}

func (s *TrackGetByIDSuite) TestSuccess(t provider.T) {
	t.Parallel()
	repo, mock := NewTrackRepository()
	track := builder.NewTrackBuilder().Default().Build()
	s.SuccessRepositoryMock(mock, track)

	resultTrack, err := repo.GetByID(context.Background(), track.ID)

	t.Assert().Nil(err)
	t.Assert().Equal(track, resultTrack)
}

func (s *TrackGetByIDSuite) NotFoundRepositoryMock(mock sqlmock.Sqlmock, track domain.Track) {
	mock.ExpectQuery(postgres.TrackGetByIDQuery).
		WillReturnError(sql.ErrNoRows)
}

func (s *TrackGetByIDSuite) TestNotFound(t provider.T) {
	t.Parallel()
	repo, mock := NewTrackRepository()
	track := builder.NewTrackBuilder().Default().Build()
	s.NotFoundRepositoryMock(mock, track)

	resultTrack, err := repo.GetByID(context.Background(), track.ID)

	t.Assert().ErrorIs(err, ports.ErrTrackIDNotFound)
	t.Assert().Equal(resultTrack, domain.Track{})
}

func TestTrackGetByIDSuite(t *testing.T) {
	suite.RunNamedSuite(t, "TrackGetByIDRepository", new(TrackGetByIDSuite))
}

type TrackGetByAlbumIDSuite struct {
	TrackSuite
}

func (s *TrackGetByAlbumIDSuite) SuccessRepositoryMock(mock sqlmock.Sqlmock, track domain.Track) {
	pgTrack := entity.NewPgTrack(track)
	expectedRows := sqlmock.NewRows(EntityColumns(pgTrack)).
		AddRow(EntityValues(pgTrack)...)
	mock.ExpectQuery(postgres.TrackGetByAlbumID).
		WillReturnRows(expectedRows)
}

func (s *TrackGetByAlbumIDSuite) TestSuccess(t provider.T) {
	t.Parallel()
	repo, mock := NewTrackRepository()
	track := builder.NewTrackBuilder().Default().Build()
	s.SuccessRepositoryMock(mock, track)

	resultTracks, err := repo.GetByAlbumID(context.Background(), uuid.New())

	t.Assert().Nil(err)
	t.Assert().Equal(track, resultTracks[0])
}

func (s *TrackGetByAlbumIDSuite) InternalErrorRepositoryMock(mock sqlmock.Sqlmock, track domain.Track) {
	mock.ExpectQuery(postgres.TrackGetByAlbumID).
		WillReturnError(sql.ErrNoRows)
}

func (s *TrackGetByAlbumIDSuite) TestNotFound(t provider.T) {
	t.Parallel()
	repo, mock := NewTrackRepository()
	track := builder.NewTrackBuilder().Default().Build()
	s.InternalErrorRepositoryMock(mock, track)

	resultTracks, err := repo.GetByAlbumID(context.Background(), track.ID)

	t.Assert().ErrorIs(err, ports.ErrInternalTrackRepo)
	t.Assert().Nil(resultTracks)
}

func TestTrackGetByAlbumIDSuite(t *testing.T) {
	suite.RunNamedSuite(t, "TrackGetByAlbumIDRepository", new(TrackGetByAlbumIDSuite))
}

type TrackGetByMusicianIDSuite struct {
	TrackSuite
}

func (s *TrackGetByMusicianIDSuite) SuccessRepositoryMock(mock sqlmock.Sqlmock, track domain.Track) {
	pgTrack := entity.NewPgTrack(track)
	expectedRows := sqlmock.NewRows(EntityColumns(pgTrack)).
		AddRow(EntityValues(pgTrack)...)
	mock.ExpectQuery(postgres.TrackGetByMusicianID).
		WillReturnRows(expectedRows)
}

func (s *TrackGetByMusicianIDSuite) TestSuccess(t provider.T) {
	t.Parallel()
	repo, mock := NewTrackRepository()
	track := builder.NewTrackBuilder().Default().Build()
	s.SuccessRepositoryMock(mock, track)

	resultTracks, err := repo.GetByMusicianID(context.Background(), uuid.New())

	t.Assert().Nil(err)
	t.Assert().Equal(track, resultTracks[0])
}

func (s *TrackGetByMusicianIDSuite) InternalErrorRepositoryMock(mock sqlmock.Sqlmock, track domain.Track) {
	mock.ExpectQuery(postgres.TrackGetByMusicianID).
		WillReturnError(sql.ErrNoRows)
}

func (s *TrackGetByMusicianIDSuite) TestNotFound(t provider.T) {
	t.Parallel()
	repo, mock := NewTrackRepository()
	track := builder.NewTrackBuilder().Default().Build()
	s.InternalErrorRepositoryMock(mock, track)

	resultTracks, err := repo.GetByMusicianID(context.Background(), track.ID)

	t.Assert().ErrorIs(err, ports.ErrInternalTrackRepo)
	t.Assert().Nil(resultTracks)
}

func TestTrackGetByMusicianIDSuite(t *testing.T) {
	suite.RunNamedSuite(t, "TrackGetByMusicianIDRepository", new(TrackGetByMusicianIDSuite))
}

type TrackGetOwnSuite struct {
	TrackSuite
}

func (s *TrackGetOwnSuite) SuccessRepositoryMock(mock sqlmock.Sqlmock, track domain.Track) {
	pgTrack := entity.NewPgTrack(track)
	expectedRows := sqlmock.NewRows(EntityColumns(pgTrack)).
		AddRow(EntityValues(pgTrack)...)
	mock.ExpectQuery(postgres.TrackGetOwn).
		WillReturnRows(expectedRows)
}

func (s *TrackGetOwnSuite) TestSuccess(t provider.T) {
	t.Parallel()
	repo, mock := NewTrackRepository()
	track := builder.NewTrackBuilder().Default().Build()
	s.SuccessRepositoryMock(mock, track)

	resultTracks, err := repo.GetOwn(context.Background(), uuid.New())

	t.Assert().Nil(err)
	t.Assert().Equal(track, resultTracks[0])
}

func (s *TrackGetOwnSuite) InternalErrorRepositoryMock(mock sqlmock.Sqlmock, track domain.Track) {
	mock.ExpectQuery(postgres.TrackGetOwn).
		WillReturnError(sql.ErrNoRows)
}

func (s *TrackGetOwnSuite) TestNotFound(t provider.T) {
	t.Parallel()
	repo, mock := NewTrackRepository()
	track := builder.NewTrackBuilder().Default().Build()
	s.InternalErrorRepositoryMock(mock, track)

	resultTracks, err := repo.GetOwn(context.Background(), track.ID)

	t.Assert().ErrorIs(err, ports.ErrInternalTrackRepo)
	t.Assert().Nil(resultTracks)
}

func TestTrackGetOwnSuite(t *testing.T) {
	suite.RunNamedSuite(t, "TrackGetOwnRepository", new(TrackGetOwnSuite))
}

type TrackGetUserFavoritesSuite struct {
	TrackSuite
}

func (s *TrackGetUserFavoritesSuite) SuccessRepositoryMock(mock sqlmock.Sqlmock, track domain.Track) {
	pgTrack := entity.NewPgTrack(track)
	expectedRows := sqlmock.NewRows(EntityColumns(pgTrack)).
		AddRow(EntityValues(pgTrack)...)
	mock.ExpectQuery(postgres.TrackGetUserFavorites).
		WillReturnRows(expectedRows)
}

func (s *TrackGetUserFavoritesSuite) TestSuccess(t provider.T) {
	t.Parallel()
	repo, mock := NewTrackRepository()
	track := builder.NewTrackBuilder().Default().Build()
	s.SuccessRepositoryMock(mock, track)

	resultTracks, err := repo.GetUserFavorites(context.Background(), uuid.New())

	t.Assert().Nil(err)
	t.Assert().Equal(track, resultTracks[0])
}

func (s *TrackGetUserFavoritesSuite) InternalErrorRepositoryMock(mock sqlmock.Sqlmock, track domain.Track) {
	mock.ExpectQuery(postgres.TrackGetUserFavorites).
		WillReturnError(sql.ErrNoRows)
}

func (s *TrackGetUserFavoritesSuite) TestNotFound(t provider.T) {
	t.Parallel()
	repo, mock := NewTrackRepository()
	track := builder.NewTrackBuilder().Default().Build()
	s.InternalErrorRepositoryMock(mock, track)

	resultTracks, err := repo.GetUserFavorites(context.Background(), track.ID)

	t.Assert().ErrorIs(err, ports.ErrInternalTrackRepo)
	t.Assert().Nil(resultTracks)
}

func TestTrackGetUserFavoritesSuite(t *testing.T) {
	suite.RunNamedSuite(t, "TrackGetUserFavoritesRepository", new(TrackGetUserFavoritesSuite))
}

type TrackAddToUserFavoritesSuite struct {
	TrackSuite
}

func (s *TrackAddToUserFavoritesSuite) SuccessRepositoryMock(mock sqlmock.Sqlmock, track domain.Track) {
	mock.ExpectExec(postgres.TrackInsertFavorite).
		WillReturnResult(sqlmock.NewResult(1, 1))
}

func (s *TrackAddToUserFavoritesSuite) TestSuccess(t provider.T) {
	t.Parallel()
	repo, mock := NewTrackRepository()
	track := builder.NewTrackBuilder().Default().Build()
	s.SuccessRepositoryMock(mock, track)

	err := repo.AddToUserFavorites(context.Background(), uuid.New(), uuid.New())

	t.Assert().Nil(err)
}

func (s *TrackAddToUserFavoritesSuite) DuplicateRepositoryMock(mock sqlmock.Sqlmock, track domain.Track) {
	mock.ExpectExec(postgres.TrackInsertFavorite).
		WillReturnError(&pgconn.PgError{Code: pgerrcode.UniqueViolation})
}

func (s *TrackAddToUserFavoritesSuite) TestDuplicate(t provider.T) {
	t.Parallel()
	repo, mock := NewTrackRepository()
	track := builder.NewTrackBuilder().Default().Build()
	s.DuplicateRepositoryMock(mock, track)

	err := repo.AddToUserFavorites(context.Background(), uuid.New(), uuid.New())

	t.Assert().ErrorIs(err, ports.ErrTrackDuplicate)
}

func TestTrackAddToUserFavoritesSuite(t *testing.T) {
	suite.RunNamedSuite(t, "TrackGetAddToUserFavoritesRepository", new(TrackAddToUserFavoritesSuite))
}
