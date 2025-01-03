package miniostorage

import (
	"context"
	"io"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go/v7"
)

type MusicianImageStorage struct {
	client     *minio.Client
	bucketName string
}

func NewMusicianImageStorage(client *minio.Client, bucketName string) *MusicianImageStorage {
	return &MusicianImageStorage{client, bucketName}
}

func (m *MusicianImageStorage) UploadImage(ctx context.Context, image io.Reader, id string) (url.URL, error) {
	object_name := id + "_image.jpg"
	_, err := m.client.PutObject(ctx, m.bucketName, object_name, image, -1, minio.PutObjectOptions{})
	if err != nil {
		return url.URL{}, err
	}

	fileURL := url.URL{
		Scheme: "http",
		Host:   strings.Replace(m.client.EndpointURL().Host, "minio", "localhost", 1),
		Path:   filepath.Join(m.bucketName, object_name),
	}

	return fileURL, nil
}
