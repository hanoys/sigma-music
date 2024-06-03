package test

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/adapters/repository/postgres"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
	"github.com/stretchr/testify/require"
	"testing"
)

var newUser = domain.User{
	ID:       uuid.New(),
	Name:     "CreatedUser",
	Email:    "CreatedUser@mail.com",
	Phone:    "+71111111111",
	Password: "CreatedUserPassword",
	Country:  "CreatedUserCountry",
}

func TestUserRepository(t *testing.T) {
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

	t.Run("create user", func(t *testing.T) {
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

		repo := postgres.NewPostgresUserRepository(db)
		createdUser, err := repo.Create(ctx, newUser)
		if err != nil {
			t.Errorf("unexcpected error: %v", err)
		}

		require.Equal(t, newUser, createdUser)
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
		repo := postgres.NewPostgresUserRepository(db)
		repo.Create(ctx, newUser)
		_, err = repo.Create(ctx, newUser)

		if !errors.Is(err, ports.ErrUserDuplicate) {
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
		repo := postgres.NewPostgresUserRepository(db)
		repo.Create(ctx, newUser)
		foundUser, err := repo.GetByID(ctx, newUser.ID)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		require.Equal(t, newUser, foundUser)
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
		repo := postgres.NewPostgresUserRepository(db)
		repo.Create(ctx, newUser)
		foundUser, err := repo.GetByEmail(ctx, newUser.Email)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		require.Equal(t, newUser, foundUser)
	})

	t.Run("find by phone", func(t *testing.T) {
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
		repo := postgres.NewPostgresUserRepository(db)
		repo.Create(ctx, newUser)
		foundUser, err := repo.GetByPhone(ctx, newUser.Phone)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		require.Equal(t, newUser, foundUser)
	})
}
