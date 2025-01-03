package miniostorage

import (
	"context"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/ports"
	"github.com/minio/minio-go/v7"
)

type TrackStorage struct {
	client     *minio.Client
	bucketName string
}

func NewTrackStorage(client *minio.Client, bucketName string) *TrackStorage {
	return &TrackStorage{client: client, bucketName: bucketName}
}

func (ts *TrackStorage) PutTrack(ctx context.Context, req ports.PutTrackReq) (url.URL, error) {
	_, err := ts.client.PutObject(ctx, ts.bucketName, req.TrackID, req.TrackBLOB, -1, minio.PutObjectOptions{})
	if err != nil {
		return url.URL{}, err
	}

	fileURL := url.URL{
		Scheme: "http",
		Host:   strings.Replace(ts.client.EndpointURL().Host, "minio", "localhost", 1),
		Path:   filepath.Join(ts.bucketName, req.TrackID),
	}

	return fileURL, nil
}

func (ts *TrackStorage) DeleteTrack(ctx context.Context, trackID uuid.UUID) error {
	return nil
}
