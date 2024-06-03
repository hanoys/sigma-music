package test

import (
	"context"
	"errors"
	"github.com/Rhymond/go-money"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/adapters/repository/postgres"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
	"testing"
	"time"
)

var newOrder = domain.Order{
	ID:         uuid.New(),
	UserID:     uuid.New(),
	CreateTime: time.Time{},
	Price:      money.New(100, money.RUB),
}

func TestOrderRepository(t *testing.T) {
	ctx := context.Background()
	container, err := newPostgresContainer(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// Clean up the container after the test is complete
	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	})

	url, err := container.ConnectionString(ctx)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("create order", func(t *testing.T) {
		t.Cleanup(func() {
			err = container.Restore(ctx)
			if err != nil {
				t.Fatal(err)
			}
		})

		db, err := newPostgresDB(url)
		if err != nil {
			t.Fatal(err)
		}
		defer db.Close()

		repo := postgres.NewPostgresOrderRepository(db)
		_, err = repo.Create(ctx, newOrder)
		if !errors.Is(err, ports.ErrInternalOrderRepo) {
			t.Errorf("unexpected error: %v", err)
		}
	})
}
