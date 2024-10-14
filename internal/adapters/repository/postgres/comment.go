package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	entity2 "github.com/hanoys/sigma-music/internal/adapters/repository/postgres/entity"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
	"github.com/hanoys/sigma-music/internal/util"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

const (
	CommentGetByIDQuery      = "SELECT * FROM comments WHERE id = $1"
	CommentGetByUserIDQuery  = "SELECT * FROM comments WHERE user_id = $1"
	CommentGetByTrackIDQuery = "SELECT * FROM comments WHERE track_id = $1"
)

type PostgresCommentRepository struct {
	connection *sqlx.DB
}

func NewPostgresCommentRepository(connection *sqlx.DB) *PostgresCommentRepository {
	return &PostgresCommentRepository{connection: connection}
}

func (cr *PostgresCommentRepository) Create(ctx context.Context, comment domain.Comment) (domain.Comment, error) {
	pgComment := entity2.NewPgComment(comment)
	queryString := entity2.InsertQueryString(pgComment, "comments")
	_, err := cr.connection.NamedExecContext(ctx, queryString, pgComment)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return domain.Comment{}, util.WrapError(ports.ErrCommentDuplicate, err)
			}
		}
		return domain.Comment{}, util.WrapError(ports.ErrInternalCommentRepo, err)
	}

	var createdTrack entity2.PgComment
	err = cr.connection.GetContext(ctx, &createdTrack, CommentGetByIDQuery, pgComment.ID)
	if err != nil {
		return domain.Comment{}, util.WrapError(ports.ErrCommentIDNotFound, err)
	}

	return createdTrack.ToDomain(), nil
}

func (cr *PostgresCommentRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Comment, error) {
	var comments []entity2.PgComment
	err := cr.connection.SelectContext(ctx, &comments, CommentGetByUserIDQuery, userID)
	if err != nil {
		return nil, util.WrapError(ports.ErrInternalCommentRepo, err)
	}

	domainComments := make([]domain.Comment, len(comments))
	for i, comment := range comments {
		domainComments[i] = comment.ToDomain()
	}

	return domainComments, nil
}

func (cr *PostgresCommentRepository) GetByTrackID(ctx context.Context, trackID uuid.UUID) ([]domain.Comment, error) {
	var comments []entity2.PgComment
	err := cr.connection.SelectContext(ctx, &comments, CommentGetByTrackIDQuery, trackID)
	if err != nil {
		return nil, util.WrapError(ports.ErrInternalCommentRepo, err)
	}

	domainComments := make([]domain.Comment, len(comments))
	for i, comment := range comments {
		domainComments[i] = comment.ToDomain()
	}

	return domainComments, nil
}
