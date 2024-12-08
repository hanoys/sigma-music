package miniostorage

import (
	"context"
	"io"
	"net/url"
	"path/filepath"

	"github.com/minio/minio-go/v7"
)

type AlbumImageStorage struct {
	client     *minio.Client
	bucketName string
}

func NewAlbumImageStorage(client *minio.Client, bucketName string) *AlbumImageStorage {
	return &AlbumImageStorage{client, bucketName}
}

func (a *AlbumImageStorage) UploadImage(ctx context.Context, image io.Reader, id string) (url.URL, error) {
	_, err := a.client.PutObject(ctx, a.bucketName, id, image, -1, minio.PutObjectOptions{})
	if err != nil {
		return url.URL{}, err
	}

	fileURL := url.URL{
		Scheme: "http",
		Host:   a.client.EndpointURL().Host,
		Path:   filepath.Join(a.bucketName, id) + ".png",
	}

	return fileURL, nil
}
