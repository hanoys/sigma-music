package test

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/adapters/repository"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
	"github.com/stretchr/testify/require"
	"testing"
)

var newMusician = domain.Musician{
	ID:       uuid.New(),
	Name:     "CreatedMusician",
	Email:    "CreatedMusician@mail.com",
	Password: "CreatedMusicianPassword",
	Country:  "CreatedMusicianCountry",
}

func TestMusicianRepository(t *testing.T) {
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

	t.Run("create musician", func(t *testing.T) {
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

		repo := repository.NewPostgresMusicianRepository(db)
		createdMusician, err := repo.Create(ctx, newMusician)
		if err != nil {
			t.Errorf("unexcpected error: %v", err)
		}

		require.Equal(t, newMusician, createdMusician)
	})

	t.Run("find by id duplicate", func(t *testing.T) {
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
		repo := repository.NewPostgresMusicianRepository(db)
		repo.Create(ctx, newMusician)
		_, err = repo.Create(ctx, newMusician)

		if !errors.Is(err, ports.ErrMusicianDuplicate) {
			t.Errorf("unexpected error %v", err)
		}
	})

	t.Run("find by id", func(t *testing.T) {
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
		repo := repository.NewPostgresMusicianRepository(db)
		repo.Create(ctx, newMusician)
		foundMusician, err := repo.GetByID(ctx, newMusician.ID)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		require.Equal(t, newMusician, foundMusician)
	})

	t.Run("find by email", func(t *testing.T) {
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
		repo := repository.NewPostgresMusicianRepository(db)
		repo.Create(ctx, newMusician)
		foundMusician, err := repo.GetByEmail(ctx, newMusician.Email)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		require.Equal(t, newMusician, foundMusician)
	})
}
