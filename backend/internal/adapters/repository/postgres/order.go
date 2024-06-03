package postgres

import (
	"context"
	"database/sql"
	"errors"
	entity2 "github.com/hanoys/sigma-music/internal/adapters/repository/postgres/entity"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
	"github.com/hanoys/sigma-music/internal/util"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

const (
	orderGetByID = "SELECT * FROM orders WHERE id = $1"
)

type PostgresOrderRepository struct {
	db *sqlx.DB
}

func NewPostgresOrderRepository(db *sqlx.DB) *PostgresOrderRepository {
	return &PostgresOrderRepository{db: db}
}

func (or *PostgresOrderRepository) Create(ctx context.Context, order domain.Order) (domain.Order, error) {
	pgOrder := entity2.NewPgOrder(order)
	queryString := entity2.InsertQueryString(pgOrder, "orders")
	_, err := or.db.NamedExecContext(ctx, queryString, pgOrder)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return domain.Order{}, util.WrapError(ports.ErrOrderDuplicate, err)
			}
		}
		return domain.Order{}, util.WrapError(ports.ErrInternalOrderRepo, err)
	}

	var createdOrder entity2.PgOrder
	err = or.db.GetContext(ctx, &createdOrder, orderGetByID, pgOrder.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Order{}, util.WrapError(ports.ErrOrderIDNotFound, err)
		}
		return domain.Order{}, util.WrapError(ports.ErrInternalOrderRepo, err)
	}

	return createdOrder.ToDomain(), nil
}
