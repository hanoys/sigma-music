package miniostorage

import (
	"context"
	"io"
	"net/url"
	"path/filepath"
	"strings"

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
	object_name := id + "_image.jpg"
	_, err := a.client.PutObject(ctx, a.bucketName, object_name, image, -1, minio.PutObjectOptions{})
	if err != nil {
		return url.URL{}, err
	}

	fileURL := url.URL{
		Scheme: "http",
		Host:   strings.Replace(a.client.EndpointURL().Host, "minio", "localhost", 1),
		Path:   filepath.Join(a.bucketName, object_name),
	}

	return fileURL, nil
}
