package test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/adapters/repository/postgres"
	"github.com/hanoys/sigma-music/internal/ports"
	"github.com/jmoiron/sqlx"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type StatSuite struct {
	suite.Suite
}

func NewStatRepository() (ports.IStatRepository, sqlmock.Sqlmock) {
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	conn := sqlx.NewDb(db, "pgx")
	repo := postgres.NewPostgresStatRepository(conn)
	return repo, mock
}

type StatAddSuite struct {
	StatSuite
}

func (s *StatAddSuite) SuccessRepositoryMock(mock sqlmock.Sqlmock) {
	mock.ExpectExec(postgres.StatAddQuery).
		WillReturnResult(sqlmock.NewResult(1, 1))
}

func (s *StatAddSuite) TestSuccess(t provider.T) {
	t.Parallel()
	repo, mock := NewStatRepository()
	s.SuccessRepositoryMock(mock)

	err := repo.Add(context.Background(), uuid.New(), uuid.New(), uuid.New())

	t.Assert().Nil(err)
}

func TestStatAddSuite(t *testing.T) {
	suite.RunNamedSuite(t, "StatAddRepository", new(StatAddSuite))
}
