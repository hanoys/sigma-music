package repository

import (
	"context"
	"github.com/hanoys/sigma-music/internal/adapters/repository/entity"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/jmoiron/sqlx"
)

const (
	subscriptionGetByID = "SELECT * FROM subscriptions WHERE id = $1"
)

type PostgresSubscriptionRepository struct {
	db *sqlx.DB
}

func NewPostgresSubscriptionRepository(db *sqlx.DB) *PostgresSubscriptionRepository {
	return &PostgresSubscriptionRepository{db: db}
}

func (sr *PostgresSubscriptionRepository) Create(ctx context.Context, sub domain.Subscription) (domain.Subscription, error) {
	pgSubscription := entity.NewPgSuscription(sub)
	queryString := entity.InsertQueryString(pgSubscription, "subscriptions")
	_, err := sr.db.NamedExecContext(ctx, queryString, pgSubscription)
	if err != nil {
		return domain.Subscription{}, err
	}

	var createdSubscription entity.PgSubscription
	err = sr.db.GetContext(ctx, &createdSubscription, subscriptionGetByID, pgSubscription.ID)
	if err != nil {
		return domain.Subscription{}, err
	}

	return createdSubscription.ToDomain(), nil
}
