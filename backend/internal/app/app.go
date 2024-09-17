package app

import (
	"context"
	"fmt"
	"github.com/hanoys/sigma-music/internal/adapters/auth"
	"github.com/hanoys/sigma-music/internal/adapters/auth/adapters"
	"github.com/hanoys/sigma-music/internal/adapters/delivery/console"
	"github.com/hanoys/sigma-music/internal/adapters/hash"
	"github.com/hanoys/sigma-music/internal/adapters/miniostorage"
	"github.com/hanoys/sigma-music/internal/adapters/repository/mongodb"
	"github.com/hanoys/sigma-music/internal/adapters/repository/postgres"
	"github.com/hanoys/sigma-music/internal/app/config"
	"github.com/hanoys/sigma-music/internal/ports"
	"github.com/hanoys/sigma-music/internal/service"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"strings"
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

type LoggerConfig struct {
	LogLevel string
}

const (
	maxConn         = 100
	maxConnIdleTime = 1 * time.Minute
	maxConnLifetime = 3 * time.Minute
)

type MongoConfig struct {
	Database string
	User     string
	Password string
	Url      string
}

func NewMongoDB(cfg *MongoConfig) (*mongo.Database, error) {
	ctx := context.Background()
	opts := options.Client().ApplyURI(cfg.Url)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		fmt.Printf("failed to connect mongodb db")
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		fmt.Printf("failed to ping mongodb db")
		return nil, err
	}

	return client.Database(cfg.Database), nil
}

func NewPostgresDB(cfg *PostgresConfig) (*PostgresDBConnections, error) {
	guestConnectionString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s",
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
	Sub      ports.ISubscriptionRepository
	Track    ports.ITrackRepository
}

func Run() {
	cfg, err := config.GetConfig(".env.local")

	if err != nil {
		log.Println("config error:", err)
		return
	}

	logger, err := NewLogger(&LoggerConfig{LogLevel: cfg.Logger.LogLevel})
	if err != nil {
		log.Println("logger error:", err)
		return
	}

	repositories := Repositories{}
	switch cfg.DB.Type {
	case "postgres":
		dbConn, err := NewPostgresDB(&PostgresConfig{
			Host:     cfg.DB.Postgres.Host,
			Port:     cfg.DB.Postgres.Port,
			Database: cfg.DB.Postgres.Name,
			User:     cfg.DB.Postgres.User,
			Password: cfg.DB.Postgres.Password,
		})

		if err != nil {
			logger.Fatal("Error connecting postgres", zap.Error(err))
			return
		}

		repositories.User = postgres.NewPostgresUserRepository(dbConn)
		repositories.Musician = postgres.NewPostgresMusicianRepository(dbConn)
		repositories.Album = postgres.NewPostgresAlbumRepository(dbConn)
		repositories.Comment = postgres.NewPostgresCommentRepository(dbConn)
		repositories.Genre = postgres.NewPostgresGenreRepository(dbConn)
		repositories.Stat = postgres.NewPostgresStatRepository(dbConn)
		repositories.Sub = postgres.NewPostgresSubscriptionRepository(dbConn)
		repositories.Track = postgres.NewPostgresTrackRepository(dbConn)

	case "mongodb":
		dbConn, err := NewMongoDB(&MongoConfig{
			Database: cfg.DB.Mongodb.Database,
			User:     cfg.DB.Mongodb.User,
			Password: cfg.DB.Mongodb.Password,
			Url:      cfg.DB.Mongodb.URL,
		})
		if err != nil {
			logger.Fatal("Error connecting mongodb", zap.Error(err))
			return
		}

		repositories.User = mongodb.NewMongoUserRepository(dbConn)
		repositories.Musician = mongodb.NewMongoMusicianRepository(dbConn)
		repositories.Album = mongodb.NewMongoAlbumRepository(dbConn)
		repositories.Comment = mongodb.NewMongoCommentRepository(dbConn)
		repositories.Genre = mongodb.NewMongoGenreRepository(dbConn)
		repositories.Stat = mongodb.NewMongoStatRepository(dbConn)
		repositories.Sub = mongodb.NewMongoSubscriptionRepository(dbConn)
		repositories.Track = mongodb.NewMongoTrackRepository(dbConn)
	default:
		logger.Fatal("Error unknown database name", zap.Error(err),
			zap.String("Database name", cfg.DB.Type))
		return
	}

	redisClient, err := NewRedisClient(&RedisConfig{
		Host: cfg.Redis.Host,
		Port: cfg.Redis.Port,
	})

	if err != nil {
		logger.Fatal("Error connecting redis", zap.Error(err))
		return
	}

	minioClient, err := NewMinioClient(&MinioConfig{
		Endpoint:     cfg.Minio.Endpoint,
		BucketName:   cfg.Minio.BucketName,
		RootUser:     cfg.Minio.RootUser,
		RootPassword: cfg.Minio.RootPassword,
	})

	if err != nil {
		logger.Fatal("Error connecting minio", zap.Error(err))
		return
	}

	userRepo := repositories.User
	musicianRepo := repositories.Musician
	albumRepo := repositories.Album
	commentRepo := repositories.Comment
	genreRepo := repositories.Genre
	statRepo := repositories.Stat
	subRepo := repositories.Sub
	trackRepo := repositories.Track

	tokenStorage := adapters.NewTokenStorage(redisClient)
	tokenProvider := auth.NewProvider(tokenStorage, &auth.ProviderConfig{
		AccessTokenExpTime:  cfg.JWT.AccessTokenExpTime,
		RefreshTokenExpTime: cfg.JWT.RefreshTokenExpTime,
		SecretKey:           cfg.JWT.SecretKey,
	})
	hashProvider := hash.NewHashPasswordProvider()
	trackStorage := miniostorage.NewTrackStorage(minioClient, cfg.Minio.BucketName)

	authService := service.NewAuthorizationService(userRepo, musicianRepo, tokenProvider, hashProvider, logger)
	userService := service.NewUserService(userRepo, hashProvider, logger)
	musicianService := service.NewMusicianService(musicianRepo, hashProvider, logger)
	albumService := service.NewAlbumService(albumRepo, logger)
	commentService := service.NewCommentService(commentRepo, logger)
	genreService := service.NewGenreService(genreRepo, logger)
	statService := service.NewStatService(statRepo, genreService, musicianService, logger)
	subService := service.NewSubscriptionService(subRepo, logger)
	trackService := service.NewTrackService(trackRepo, trackStorage, genreService, logger)

	cons := console.NewConsole(console.NewHandler(console.HandlerParams{
		AlbumService:        albumService,
		AuthService:         authService,
		CommentService:      commentService,
		GenreService:        genreService,
		MusicianService:     musicianService,
		StatService:         statService,
		SubscriptionService: subService,
		TrackService:        trackService,
		UserService:         userService,
	}))

	cons.Start()
}
