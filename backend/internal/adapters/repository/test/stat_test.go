package test

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/adapters/repository"
	"github.com/hanoys/sigma-music/internal/ports"
	"testing"
)

func TestStatRepository(t *testing.T) {
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

	t.Run("add stat", func(t *testing.T) {
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

		repo := repository.NewPostgresStatRepository(db)
		err = repo.Add(context.Background(), uuid.New(), uuid.New())
		if !errors.Is(err, ports.ErrInternalStatRepo) {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("musicians stat", func(t *testing.T) {
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

		repo := repository.NewPostgresStatRepository(db)
		stat, err := repo.GetMostListenedMusicians(context.Background(), uuid.New(), 10)
		if len(stat) != 0 {
			t.Errorf("unexpected len stat len: %v", len(stat))
		}
	})

	t.Run("genre stat", func(t *testing.T) {
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

		repo := repository.NewPostgresStatRepository(db)
		stat, err := repo.GetListenedGenres(context.Background(), uuid.New())
		if len(stat) != 0 {
			t.Errorf("unexpected len stat len: %v", len(stat))
		}
	})
}
