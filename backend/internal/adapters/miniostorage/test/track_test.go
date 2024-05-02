package test

import (
	"bytes"
	"context"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/adapters/miniostorage"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

func TestTrackStorage(t *testing.T) {
	ctx := context.Background()
	container, err := newMinioContainer(ctx)
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

	minioClient, err := newMinioClient(url)
	if err != nil {
		t.Fatal(err)
	}

	store := miniostorage.NewTrackStorage(minioClient, BucketName)

	_, path, _, ok := runtime.Caller(0)
	require.Equal(t, ok, true)
	filesPath := filepath.Dir(path) + "/files"

	t.Run("test put track", func(t *testing.T) {
		testFilename := "MorgenshternPABLO.mp3"
		data, err := os.ReadFile(filepath.Join(filesPath, testFilename))
		if err != nil {
			t.Errorf("failed to read file %s: %s", testFilename, err)
		}

		fileURL, err := store.PutTrack(ctx, uuid.New().String(), bytes.NewReader(data))
		if err != nil {
			t.Errorf("failed to save file to minio: %v", err)
		}

		resp, err := http.Get(fileURL.String())
		if err != nil {
			t.Errorf("failed to download saved file: %v", err)
		}
		defer resp.Body.Close()

		savedData, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("failed to read response data: %s", err)
		}
		require.Equal(t, reflect.DeepEqual(data, savedData), true)
	})
}
