package test

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/adapters/repository"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
	"testing"
	"time"
)

var newSubscription = domain.Subscription{
	ID:             uuid.New(),
	UserID:         uuid.New(),
	OrderID:        uuid.New(),
	StartDate:      time.Now(),
	ExpirationDate: time.Now(),
}

func TestSubscriptionRepository(t *testing.T) {
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

	t.Run("create subscription", func(t *testing.T) {
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

		repo := repository.NewPostgresSubscriptionRepository(db)
		_, err = repo.Create(ctx, newSubscription)
		if !errors.Is(err, ports.ErrInternalSubRepo) {
			t.Errorf("unexpected error: %v", err)
		}
	})
}
