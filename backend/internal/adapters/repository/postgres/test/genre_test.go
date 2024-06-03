package test

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/adapters/repository/postgres"
	"github.com/hanoys/sigma-music/internal/ports"
	"testing"
)

func TestGenreRepository(t *testing.T) {
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

	t.Run("get all genre", func(t *testing.T) {
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

		repo := postgres.NewPostgresGenreRepository(db)
		genres, err := repo.GetAll(context.Background())
		if err != nil {
			t.Errorf("unexcpected error: %v", err)
		}

		if len(genres) != 0 {
			t.Errorf("len is not zero")
		}

		t.Run("get by id genre", func(t *testing.T) {
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

			repo := postgres.NewPostgresGenreRepository(db)
			_, err = repo.GetByID(context.Background(), uuid.New())

			if !errors.Is(err, ports.ErrGenreIDNotFound) {
				t.Errorf("unexcpected error: %v", err)
			}
		})
	})
}
