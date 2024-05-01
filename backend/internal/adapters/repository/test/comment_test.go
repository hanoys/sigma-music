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

var newComment = domain.Comment{
	ID:      uuid.New(),
	UserID:  uuid.New(),
	TrackID: uuid.New(),
	Stars:   0,
	Text:    "text",
}

func TestCommentRepository(t *testing.T) {
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

	t.Run("create comment", func(t *testing.T) {
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

		repo := repository.NewPostgresCommentRepository(db)
		_, err = repo.Create(ctx, newComment)
		if !errors.Is(err, ports.ErrInternalCommentRepo) {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("get by user id comment", func(t *testing.T) {
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

		repo := repository.NewPostgresCommentRepository(db)
		comments, err := repo.GetByUserID(ctx, uuid.New())
		if len(comments) != 0 {
			t.Errorf("unexpected comment count: %v", len(comments))
		}
	})

	t.Run("get by track id comment", func(t *testing.T) {
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

		repo := repository.NewPostgresCommentRepository(db)
		comments, err := repo.GetByUserID(ctx, uuid.New())
		if len(comments) != 0 {
			t.Errorf("unexpected comment count: %v", len(comments))
		}
	})
}
