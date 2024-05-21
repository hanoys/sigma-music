package app

import (
	"context"
	"fmt"
	"github.com/hanoys/sigma-music/internal/adapters/auth"
	"github.com/hanoys/sigma-music/internal/adapters/auth/adapters"
	"github.com/hanoys/sigma-music/internal/adapters/delivery/console"
	"github.com/hanoys/sigma-music/internal/adapters/hash"
	"github.com/hanoys/sigma-music/internal/adapters/miniostorage"
	"github.com/hanoys/sigma-music/internal/adapters/repository"
	"github.com/hanoys/sigma-music/internal/app/config"
	"github.com/hanoys/sigma-music/internal/service"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"time"
)

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
	Endpoint     string
	BucketName   string
	RootUser     string
	RootPassword string
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

func NewMinioClient(cfg *MinioConfig) (*minio.Client, error) {
	minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.RootUser, cfg.RootPassword, ""),
		Secure: false,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create minio client")
	}

	ctx := context.Background()
	err = minioClient.MakeBucket(ctx, cfg.BucketName, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(ctx, cfg.BucketName)
		if errBucketExists != nil || !exists {
			return nil, errors.Wrap(errBucketExists, "failed to make minio bucket")
		}
	}
	policy := fmt.Sprintf(`{
		"Version":"2012-10-17",
		"Statement":[{
			"Effect":"Allow",
			"Principal":"*",
			"Action":["s3:GetObject"],
			"Resource":["arn:aws:s3:::%s/*"]}
		]}`, cfg.BucketName)
	err = minioClient.SetBucketPolicy(ctx, cfg.BucketName, policy)
	if err != nil {
		return nil, errors.Wrap(err, "failed to set bucket public policy")
	}

	return minioClient, nil
}

func Run() {
	cfg, err := config.GetConfig(".env.local")

	if err != nil {
		fmt.Println("config error:", err)
		return
	}

	dbConn, err := NewPostgresDB(&PostgresConfig{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		Database: cfg.DB.Name,
		User:     cfg.DB.User,
		Password: cfg.DB.Password,
	})

	if err != nil {
		fmt.Println("postgres error:", err)
		return
	}

	redisClient, err := NewRedisClient(&RedisConfig{
		Host: cfg.Redis.Host,
		Port: cfg.Redis.Port,
	})

	if err != nil {
		fmt.Println("redis error:", err)
		return
	}

	minioClient, err := NewMinioClient(&MinioConfig{
		Endpoint:     cfg.Minio.Endpoint,
		BucketName:   cfg.Minio.BucketName,
		RootUser:     cfg.Minio.RootUser,
		RootPassword: cfg.Minio.RootPassword,
	})

	if err != nil {
		fmt.Println("minio error:", err)
		return
	}

	userRepo := repository.NewPostgresUserRepository(dbConn)
	musicianRepo := repository.NewPostgresMusicianRepository(dbConn)
	albumRepo := repository.NewPostgresAlbumRepository(dbConn)
	commentRepo := repository.NewPostgresCommentRepository(dbConn)
	genreRepo := repository.NewPostgresGenreRepository(dbConn)
	orderRepo := repository.NewPostgresOrderRepository(dbConn)
	statRepo := repository.NewPostgresStatRepository(dbConn)
	subRepo := repository.NewPostgresSubscriptionRepository(dbConn)
	trackRepo := repository.NewPostgresTrackRepository(dbConn)

	tokenStorage := adapters.NewTokenStorage(redisClient)
	tokenProvider := auth.NewProvider(tokenStorage, &auth.ProviderConfig{
		AccessTokenExpTime:  cfg.JWT.AccessTokenExpTime,
		RefreshTokenExpTime: cfg.JWT.RefreshTokenExpTime,
		SecretKey:           cfg.JWT.SecretKey,
	})
	hashProvider := hash.NewHashPasswordProvider()
	trackStorage := miniostorage.NewTrackStorage(minioClient, cfg.Minio.BucketName)

	authService := service.NewAuthorizationService(userRepo, musicianRepo, tokenProvider, hashProvider)
	userService := service.NewUserService(userRepo, hashProvider)
	musicianService := service.NewMusicianService(musicianRepo, hashProvider)
	albumService := service.NewAlbumService(albumRepo)
	commentService := service.NewCommentService(commentRepo)
	genreService := service.NewGenreService(genreRepo)
	orderService := service.NewOrderService(orderRepo)
	statService := service.NewStatService(statRepo, genreService, musicianService)
	subService := service.NewSubscriptionService(subRepo)
	trackService := service.NewTrackService(trackRepo, trackStorage, genreService)

	cons := console.NewConsole(console.NewHandler(console.HandlerParams{
		AlbumService:        albumService,
		AuthService:         authService,
		CommentService:      commentService,
		GenreService:        genreService,
		MusicianService:     musicianService,
		OrderService:        orderService,
		StatService:         statService,
		SubscriptionService: subService,
		TrackService:        trackService,
		UserService:         userService,
	}))

	cons.Start()
}
