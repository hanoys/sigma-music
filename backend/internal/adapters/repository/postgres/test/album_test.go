package test

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/adapters/repository/postgres"
	"github.com/hanoys/sigma-music/internal/ports"
	"testing"
)

func TestAlbumRepository(t *testing.T) {
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

	t.Run("get album by musician id", func(t *testing.T) {
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
		defer func() {
			db.GuestConnection.Close()
			db.UserConnection.Close()
			db.MusicianConnection.Close()
		}()

		repo := postgres.NewPostgresAlbumRepository(db)
		albums, err := repo.GetByMusicianID(ctx, uuid.New())
		if len(albums) != 0 {
			t.Errorf("unexpected albums count: %v", len(albums))
		}
	})

	t.Run("get album by id", func(t *testing.T) {
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
		defer func() {
			db.GuestConnection.Close()
			db.UserConnection.Close()
			db.MusicianConnection.Close()
		}()

		repo := postgres.NewPostgresAlbumRepository(db)
		_, err = repo.GetByID(ctx, uuid.New())
		if !errors.Is(err, ports.ErrAlbumIDNotFound) {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("publish album", func(t *testing.T) {
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
		defer func() {
			db.GuestConnection.Close()
			db.UserConnection.Close()
			db.MusicianConnection.Close()
		}()

		repo := postgres.NewPostgresAlbumRepository(db)
		err = repo.Publish(ctx, uuid.New())
		if !errors.Is(err, ports.ErrAlbumIDNotFound) {
			t.Errorf("unexpected error: %v", err)
		}
	})
}
