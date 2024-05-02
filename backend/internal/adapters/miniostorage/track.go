package miniostorage

import (
	"context"
	"github.com/minio/minio-go/v7"
	"io"
	"net/url"
	"path/filepath"
)

type TrackStorage struct {
	client     *minio.Client
	bucketName string
}

func NewTrackStorage(client *minio.Client, bucketName string) *TrackStorage {
	return &TrackStorage{client: client, bucketName: bucketName}
}

func (ts *TrackStorage) PutTrack(ctx context.Context, filename string, track io.Reader) (url.URL, error) {
	_, err := ts.client.PutObject(ctx, ts.bucketName, filename, track, -1, minio.PutObjectOptions{})
	if err != nil {
		return url.URL{}, err
	}

	fileURL := url.URL{
		Scheme: "http",
		Host:   ts.client.EndpointURL().Host,
		Path:   filepath.Join(ts.bucketName, filename),
	}

	return fileURL, nil
}
