package repository

import (
	"context"
	"github.com/hanoys/sigma-music/internal/adapters/repository/entity"
	"github.com/hanoys/sigma-music/internal/domain"
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
	pgOrder := entity.NewPgOrder(order)
	queryString := entity.InsertQueryString(pgOrder, "orders")
	_, err := or.db.NamedExecContext(ctx, queryString, pgOrder)
	if err != nil {
		return domain.Order{}, err
	}

	var createdOrder entity.PgOrder
	err = or.db.GetContext(ctx, &createdOrder, orderGetByID, pgOrder.ID)
	if err != nil {
		return domain.Order{}, err
	}

	return createdOrder.ToDomain(), nil
}
