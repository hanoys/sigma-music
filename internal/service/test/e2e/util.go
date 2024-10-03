package e2e

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"time"

	"github.com/golang-migrate/migrate"
	migratepg "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/testcontainers/testcontainers-go"
	testminio "github.com/testcontainers/testcontainers-go/modules/minio"
	testpg "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	DatabaseName = "sigma-music"
	UserName     = "sigma-music"
	Password     = "sigma-music"
)

func newPostgresContainer(ctx context.Context) (*testpg.PostgresContainer, error) {
	container, err := testpg.RunContainer(
		ctx,
		testpg.WithDatabase(DatabaseName),
		testpg.WithUsername(UserName),
		testpg.WithPassword(Password),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)

	if err != nil {
		return nil, fmt.Errorf("creating container failure: %w", err)
	}

	_, path, _, ok := runtime.Caller(0)
	if !ok {
		return nil, fmt.Errorf("failed to get caller path")
	}

	sourceUrl := "file://" + filepath.Dir(path) + "/migrations"
	url, err := container.ConnectionString(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get postgres db url: %s", err)
	}

	db, err := newPostgresDB(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres db: %s", err)
	}
	defer db.Close()

	driver, err := migratepg.WithInstance(db.DB, &migratepg.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to get db driver from instance: %s", err)
	}

	mig, err := migrate.NewWithDatabaseInstance(sourceUrl, DatabaseName, driver)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrator driver: %s", err)
	}

	err = mig.Up()
	if err != nil {
		return nil, fmt.Errorf("failed to up migrations: %s", err)
	}

	return container, nil
}

const (
	maxConn         = 100
	maxConnIdleTime = 1 * time.Minute
	maxConnLifetime = 3 * time.Minute
)

func newPostgresDB(url string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("pgx", url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect postgres db: %s", err)
	}

	db.SetMaxOpenConns(maxConn)
	db.SetConnMaxLifetime(maxConnLifetime)
	db.SetConnMaxIdleTime(maxConnIdleTime)

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping postgres db: %s", err)
	}

	return db, nil
}

var (
	MinioEndpoint     = "minio:9000"
	MinioBucketName   = "music"
	MinioRootUser     = "sigma"
	MinioRootPassword = "sigmamusic"
)

func newMinioContainer(ctx context.Context) (*testminio.MinioContainer, error) {
	minioContainer, err := testminio.RunContainer(ctx,
		testcontainers.WithImage("docker.io/minio/minio"),
		testminio.WithUsername(MinioRootUser),
		testminio.WithPassword(MinioRootPassword),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start minio container: %s", err)
	}
	return minioContainer, nil
}

func newMinioClient(url string) (*minio.Client, error) {
	minioClient, err := minio.New(url, &minio.Options{
		Creds:  credentials.NewStaticV4(MinioRootUser, MinioRootPassword, ""),
		Secure: false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %s", err)
	}

	ctx := context.Background()
	err = minioClient.MakeBucket(ctx, MinioBucketName, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(ctx, MinioBucketName)
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
		]}`, MinioBucketName)
	err = minioClient.SetBucketPolicy(ctx, MinioBucketName, policy)
	if err != nil {
		return nil, fmt.Errorf("failed to set bucket public policy: %s", err)
	}

	return minioClient, nil
}
