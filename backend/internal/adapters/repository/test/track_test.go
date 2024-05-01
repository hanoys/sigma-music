package test

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/adapters/repository"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
	"testing"
)

var newTrack = domain.Track{
	ID:      uuid.New(),
	AlbumID: uuid.New(),
	Name:    "Track",
	URL:     "URL",
}

func TestTrackRepository(t *testing.T) {
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

	t.Run("create track", func(t *testing.T) {
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

		repo := repository.NewPostgresTrackRepository(db)
		_, err = repo.Create(ctx, newTrack)
		if !errors.Is(err, ports.ErrInternalTrackRepo) {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("delete track", func(t *testing.T) {
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

		repo := repository.NewPostgresTrackRepository(db)
		_, err = repo.Delete(ctx, uuid.New())
		if !errors.Is(err, ports.ErrTrackIDNotFound) {
			t.Errorf("unexpected error: %v", err)
		}
	})
}
