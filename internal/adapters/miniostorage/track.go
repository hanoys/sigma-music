package miniostorage

import (
	"context"
	"io"
	"net/url"
	"path/filepath"

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
		Host:   ts.client.EndpointURL().Host,
		Path:   filepath.Join(ts.bucketName, req.TrackID),
	}

	return fileURL, nil
}

func (ts *TrackStorage) UploadImage(ctx context.Context, image io.Reader, id string) (url.URL, error) {
	_, err := ts.client.PutObject(ctx, ts.bucketName, id, image, -1, minio.PutObjectOptions{})
	if err != nil {
		return url.URL{}, err
	}

	fileURL := url.URL{
		Scheme: "http",
		Host:   ts.client.EndpointURL().Host,
		Path:   filepath.Join(ts.bucketName, id) + "_image" + ".png",
	}

	return fileURL, nil
}

func (ts *TrackStorage) DeleteTrack(ctx context.Context, trackID uuid.UUID) error {
	return nil
}
