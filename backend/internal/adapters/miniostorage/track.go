package miniostorage

import (
	"context"
	"github.com/hanoys/sigma-music/internal/ports"
	"github.com/minio/minio-go"
)

type TrackStorage struct {
	client     *minio.Client
	bucketName string
}

func NewTrackStorage(client *minio.Client, bucketName string) *TrackStorage {
	return &TrackStorage{client: client, bucketName: bucketName}
}

func (ts *TrackStorage) PutTrack(ctx context.Context, req ports.PutTrackReq) error {
	_, err := ts.client.PutObjectWithContext(ctx, ts.bucketName, req.TrackID, req.TrackBLOB, req.TrackSize, minio.PutObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}
