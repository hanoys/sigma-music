package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/adapters/repository/entity"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
	"github.com/hanoys/sigma-music/internal/utill"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

const (
	commentGetByIDQuery      = "SELECT * FROM comments WHERE id = $1"
	commentGetByUserIDQuery  = "SELECT * FROM comments WHERE user_id = $1"
	commentGetByTrackIDQuery = "SELECT * FROM comments WHERE track_id = $1"
)

type PostgresCommentRepository struct {
	db *sqlx.DB
}

func NewPostgresCommentRepository(db *sqlx.DB) *PostgresCommentRepository {
	return &PostgresCommentRepository{db: db}
}

func (cr *PostgresCommentRepository) Create(ctx context.Context, comment domain.Comment) (domain.Comment, error) {
	pgComment := entity.NewPgComment(comment)
	queryString := entity.InsertQueryString(pgComment, "comments")
	_, err := cr.db.NamedExecContext(ctx, queryString, pgComment)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return domain.Comment{}, utill.WrapError(ports.ErrCommentDuplicate, err)
			}
		}
		return domain.Comment{}, utill.WrapError(ports.ErrInternalCommentRepo, err)
	}

	var createdTrack entity.PgComment
	err = cr.db.GetContext(ctx, &createdTrack, commentGetByIDQuery, pgComment.ID)
	if err != nil {
		return domain.Comment{}, utill.WrapError(ports.ErrCommentIDNotFound, err)
	}

	return createdTrack.ToDomain(), nil
}

func (cr *PostgresCommentRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Comment, error) {
	var comments []entity.PgComment
	err := cr.db.SelectContext(ctx, &comments, commentGetByUserIDQuery, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utill.WrapError(ports.ErrCommentByUserIDNotFound, err)
		}
		return nil, utill.WrapError(ports.ErrInternalCommentRepo, err)
	}

	domainComments := make([]domain.Comment, len(comments))
	for i, comment := range comments {
		domainComments[i] = comment.ToDomain()
	}

	return domainComments, nil
}

func (cr *PostgresCommentRepository) GetByTrackID(ctx context.Context, trackID uuid.UUID) ([]domain.Comment, error) {
	var comments []entity.PgComment
	err := cr.db.SelectContext(ctx, &comments, commentGetByTrackIDQuery, trackID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utill.WrapError(ports.ErrCommentByTrackIDNotFound, err)
		}
		return nil, utill.WrapError(ports.ErrInternalCommentRepo, err)
	}

	domainComments := make([]domain.Comment, len(comments))
	for i, comment := range comments {
		domainComments[i] = comment.ToDomain()
	}

	return domainComments, nil
}
