package test

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/testcontainers/testcontainers-go"
	testminio "github.com/testcontainers/testcontainers-go/modules/minio"
)

var (
	AccessKey  = "key"
	SecretKey  = "secretKey"
	BucketName = "music"
)

func newMinioContainer(ctx context.Context) (*testminio.MinioContainer, error) {
	minioContainer, err := testminio.RunContainer(ctx,
		testcontainers.WithImage("docker.io/minio/minio"),
		testminio.WithUsername(AccessKey),
		testminio.WithPassword(SecretKey),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start minio container: %s", err)
	}
	return minioContainer, nil
}

func newMinioClient(url string) (*minio.Client, error) {
	minioClient, err := minio.New(url, &minio.Options{
		Creds:  credentials.NewStaticV4(AccessKey, SecretKey, ""),
		Secure: false,
	})

	ctx := context.Background()
	err = minioClient.MakeBucket(ctx, BucketName, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(ctx, BucketName)
		if errBucketExists != nil || !exists {
			return nil, fmt.Errorf("failed to make minio bucket: %s", errBucketExists)
		}
	}
	policy := fmt.Sprintf(`{
		"Version":"2012-10-17",
		"Statement":[{
			"Effect":"Allow",
			"Principal":"*",
			"Action":["s3:GetObject"],
			"Resource":["arn:aws:s3:::%s/*"]}
		]}`, BucketName)
	err = minioClient.SetBucketPolicy(ctx, BucketName, policy)
	if err != nil {
		return nil, fmt.Errorf("failed to set bucket public policy: %s", err)
	}

	return minioClient, nil
}
