package config

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/JeremyLoy/config"
	"github.com/hanoys/sigma-music/internal/ports"
	"github.com/jmoiron/sqlx"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Config struct {
	DB struct {
		Type     string `yaml:"type"`
		Postgres struct {
			User     string `yaml:"user"`
			Password string `yaml:"password"`
			Host     string `yaml:"host"`
			Port     string `yaml:"port"`
			Name     string `yaml:"name"`
		} `yaml:"postgres"`
	} `yaml:"db"`

	JWT struct {
		AccessTokenExpTime  int64  `yaml:"access_expiration_time"`
		RefreshTokenExpTime int64  `yaml:"refresh_expiration_time"`
		SecretKey           string `yaml:"secret"`
	} `yaml:"jwt"`

	Redis struct {
		Host string `config:"REDIS_HOST"`
		Port string `config:"REDIS_PORT"`
	}

	Minio struct {
		Endpoint             string `config:"MINIO_ENDPOINT"`
		TrackBucketName      string `config:"TRACK_MINIO_BUCKET_NAME"`
		AlbumImageBucketName string `config:"ALBUM_IMAGE_MINIO_BUCKET_NAME"`
		RootUser             string `config:"MINIO_ROOT_USER"`
		RootPassword         string `config:"MINIO_ROOT_PASSWORD"`
	}

	Logger struct {
		LogLevel string `yaml:"level"`
	} `yaml:"log"`
}

func GetConfig(configPath string) (*Config, error) {
	var conf Config

	buf, err := os.ReadFile("./config/config.yml")
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(buf, &conf)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(buf, &conf.DB)
	if err != nil {
		return nil, err
	}

	err = config.From(configPath).To(&conf.Redis)
	if err != nil {
		return nil, err
	}

	err = config.From(configPath).To(&conf.Minio)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}

type PostgresConfig struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
}

type RedisConfig struct {
	Host string
	Port string
}

type MinioConfig struct {
	Endpoint             string
	TrackBucketName      string
	AlbumImageBucketName string
	RootUser             string
	RootPassword         string
}

type LoggerConfig struct {
	LogLevel string
}

const (
	maxConn         = 100
	maxConnIdleTime = 1 * time.Minute
	maxConnLifetime = 3 * time.Minute
)

func NewPostgresDB(cfg *PostgresConfig) (*sqlx.DB, error) {
	connectionString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Database,
		cfg.Password,
	)

	db, err := sqlx.Connect("pgx", connectionString)
	if err != nil {
		fmt.Printf("failed to connect postgres db: %s", connectionString)
		return nil, err
	}

	db.SetMaxOpenConns(maxConn)
	db.SetConnMaxLifetime(maxConnLifetime)
	db.SetConnMaxIdleTime(maxConnIdleTime)

	err = db.Ping()
	if err != nil {
		fmt.Printf("failed to ping postgres db: %s", connectionString)
		return nil, err
	}

	return db, nil
}

func NewRedisClient(cfg *RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: cfg.Host + ":" + cfg.Port,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return client, nil
}

func minioCreateBucket(ctx context.Context, minioClient *minio.Client, bucketName string) error {
	err := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists != nil || !exists {
			return errors.Wrap(errBucketExists, "failed to make minio bucket")
		}
	}
	policy := fmt.Sprintf(`{
		"Version":"2012-10-17",
		"Statement":[{
			"Effect":"Allow",
			"Principal":"*",
			"Action":["s3:GetObject"],
			"Resource":["arn:aws:s3:::%s/*"]}
		]}`, bucketName)
	err = minioClient.SetBucketPolicy(ctx, bucketName, policy)
	if err != nil {
		return errors.Wrap(err, "failed to set bucket public policy")
	}

	return nil
}

func NewMinioClient(cfg *MinioConfig) (*minio.Client, error) {
	minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.RootUser, cfg.RootPassword, ""),
		Secure: false,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create minio client")
	}

	ctx := context.Background()
	minioCreateBucket(ctx, minioClient, cfg.TrackBucketName)
	minioCreateBucket(ctx, minioClient, cfg.AlbumImageBucketName)
	return minioClient, nil
}

func NewLogger(cfg *LoggerConfig) (*zap.Logger, error) {
	var logLevel zap.AtomicLevel
	if strings.ToLower(cfg.LogLevel) == "info" {
		logLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
	} else if strings.ToLower(cfg.LogLevel) == "error" {
		logLevel = zap.NewAtomicLevelAt(zap.ErrorLevel)
	} else if strings.ToLower(cfg.LogLevel) == "fatal" {
		logLevel = zap.NewAtomicLevelAt(zap.FatalLevel)
	} else {
		return nil, fmt.Errorf("unknown log level: %s", cfg.LogLevel)
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	conf := zap.Config{
		Level:             logLevel,
		Development:       false,
		DisableCaller:     true,
		DisableStacktrace: true,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig:     encoderCfg,
		OutputPaths: []string{
			"./log/sigma-music.log",
		},
		ErrorOutputPaths: []string{
			"./log/sigma-music.log",
		},
	}

	return zap.Must(conf.Build()), nil
}

type Repositories struct {
	User     ports.IUserRepository
	Musician ports.IMusicianRepository
	Album    ports.IAlbumRepository
	Comment  ports.ICommentRepository
	Genre    ports.IGenreRepository
	Stat     ports.IStatRepository
	Track    ports.ITrackRepository
}
