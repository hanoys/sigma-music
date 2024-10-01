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

type CommentSuite struct {
	suite.Suite
}

type CommentCreateSuite struct {
	CommentSuite
}

func (s *CommentCreateSuite) SuccessRepositoryMock(mock sqlmock.Sqlmock, comment domain.Comment) {
	pgComment := entity.NewPgComment(comment)
	queryString := InsertQueryString(pgComment, "comments")
	mock.ExpectExec(queryString).
		WithArgs(EntityValues(pgComment)...).
		WillReturnResult(sqlmock.NewResult(1, 1))

	expectedRows := sqlmock.NewRows(EntityColumns(pgComment)).
		AddRow(EntityValues(pgComment)...)
	mock.ExpectQuery(postgres.CommentGetByIDQuery).
		WithArgs(pgComment.ID).
		WillReturnRows(expectedRows)
}

func (s *CommentCreateSuite) TestSuccess(t provider.T) {
	t.Parallel()
	t.Title("Repository Comment create test success")
	repo, mock := NewCommentRepository()
	comment := builder.NewCommentBuilder().Default().Build()
	s.SuccessRepositoryMock(mock, comment)

	commentResult, err := repo.Create(context.Background(), comment)

	t.Assert().Nil(err)
	t.Assert().Equal(comment, commentResult)
}

func (s *CommentCreateSuite) DuplicateRepositoryMock(mock sqlmock.Sqlmock, comment domain.Comment) {
	pgComment := entity.NewPgComment(comment)
	queryString := InsertQueryString(pgComment, "comments")
	mock.ExpectExec(queryString).
		WithArgs(EntityValues(pgComment)...).
		WillReturnError(&pgconn.PgError{Code: pgerrcode.UniqueViolation})
}

func (s *CommentCreateSuite) TestDuplicate(t provider.T) {
	t.Parallel()
	t.Title("Repository Comment create test duplicate")
	repo, mock := NewCommentRepository()
	comment := builder.NewCommentBuilder().Default().Build()
	s.DuplicateRepositoryMock(mock, comment)

	commentResult, err := repo.Create(context.Background(), comment)

	t.Assert().ErrorIs(err, ports.ErrCommentDuplicate)
	t.Assert().Equal(domain.Comment{}, commentResult)
}

func TestCommentCreateSuite(t *testing.T) {
	suite.RunNamedSuite(t, "CommentCreateRepository", new(CommentCreateSuite))
}

func NewCommentRepository() (ports.ICommentRepository, sqlmock.Sqlmock) {
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	conn := sqlx.NewDb(db, "pgx")
	repo := postgres.NewPostgresCommentRepository(conn)
	return repo, mock
}

type CommentGetByUserIDSuite struct {
	CommentSuite
}

func (s *CommentGetByUserIDSuite) SuccessRepositoryMock(mock sqlmock.Sqlmock, comment domain.Comment) {
	pgComment := entity.NewPgComment(comment)
	expectedRows := sqlmock.NewRows(EntityColumns(pgComment)).
		AddRow(EntityValues(pgComment)...)
	mock.ExpectQuery(postgres.CommentGetByUserIDQuery).
		WithArgs(comment.UserID).
		WillReturnRows(expectedRows)
}

func (s *CommentGetByUserIDSuite) TestSuccess(t provider.T) {
	t.Parallel()
	t.Title("Repository Comment get by user id test success")
	repo, mock := NewCommentRepository()
	comment := builder.NewCommentBuilder().Default().Build()
	s.SuccessRepositoryMock(mock, comment)

	comments, err := repo.GetByUserID(context.Background(), comment.UserID)

	t.Assert().Nil(err)
	t.Assert().Equal(comment, comments[0])
}

func (s *CommentGetByUserIDSuite) InternalErrorRepositoryMock(mock sqlmock.Sqlmock, comment domain.Comment) {
	mock.ExpectQuery(postgres.CommentGetByUserIDQuery).
		WithArgs(comment.UserID).
		WillReturnError(sql.ErrNoRows)
}

func (s *CommentGetByUserIDSuite) TestInternalError(t provider.T) {
	t.Parallel()
	t.Title("Repository Comment get by user id test internal error")
	repo, mock := NewCommentRepository()
	comment := builder.NewCommentBuilder().Default().Build()
	s.InternalErrorRepositoryMock(mock, comment)

	comments, err := repo.GetByUserID(context.Background(), comment.UserID)

	t.Assert().Nil(comments)
	t.Assert().ErrorIs(err, ports.ErrInternalCommentRepo)
}

func TestCommentGetByUserIDSuite(t *testing.T) {
	suite.RunNamedSuite(t, "CommentGetByUserIDRepository", new(CommentGetByUserIDSuite))
}

type CommentGetByTrackIDSuite struct {
	CommentSuite
}

func (s *CommentGetByTrackIDSuite) SuccessRepositoryMock(mock sqlmock.Sqlmock, comment domain.Comment) {
	pgComment := entity.NewPgComment(comment)
	expectedRows := sqlmock.NewRows(EntityColumns(pgComment)).
		AddRow(EntityValues(pgComment)...)
	mock.ExpectQuery(postgres.CommentGetByTrackIDQuery).
		WithArgs(comment.TrackID).
		WillReturnRows(expectedRows)
}

func (s *CommentGetByTrackIDSuite) TestSuccess(t provider.T) {
	t.Parallel()
	t.Title("Repository Comment get by track id test success")
	repo, mock := NewCommentRepository()
	comment := builder.NewCommentBuilder().Default().Build()
	s.SuccessRepositoryMock(mock, comment)

	comments, err := repo.GetByTrackID(context.Background(), comment.TrackID)

	t.Assert().Nil(err)
	t.Assert().Equal(comment, comments[0])
}

func (s *CommentGetByTrackIDSuite) InternalErrorRepositoryMock(mock sqlmock.Sqlmock, comment domain.Comment) {
	mock.ExpectQuery(postgres.CommentGetByTrackIDQuery).
		WithArgs(comment.TrackID).
		WillReturnError(sql.ErrNoRows)
}

func (s *CommentGetByTrackIDSuite) TestInternalError(t provider.T) {
	t.Parallel()
	t.Title("Repository Comment get by track id test internal error")
	repo, mock := NewCommentRepository()
	comment := builder.NewCommentBuilder().Default().Build()
	s.InternalErrorRepositoryMock(mock, comment)

	comments, err := repo.GetByUserID(context.Background(), comment.TrackID)

	t.Assert().Nil(comments)
	t.Assert().ErrorIs(err, ports.ErrInternalCommentRepo)
}

func TestCommentGetByTrackIDSuite(t *testing.T) {
	suite.RunNamedSuite(t, "CommentGetByTrackIDRepository", new(CommentGetByTrackIDSuite))
}
