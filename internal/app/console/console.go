package console

import (
	"log"

	"github.com/hanoys/sigma-music/internal/adapters/auth"
	"github.com/hanoys/sigma-music/internal/adapters/auth/adapters"
	consd "github.com/hanoys/sigma-music/internal/adapters/delivery/console"
	"github.com/hanoys/sigma-music/internal/adapters/hash"
	"github.com/hanoys/sigma-music/internal/adapters/miniostorage"
	"github.com/hanoys/sigma-music/internal/adapters/repository/postgres"
	"github.com/hanoys/sigma-music/internal/app/config"
	"github.com/hanoys/sigma-music/internal/service"
	"go.uber.org/zap"
)

func Run() {
	cfg, err := config.GetConfig(".env.local")

	if err != nil {
		log.Println("config error:", err)
		return
	}

	logger, err := config.NewLogger(&config.LoggerConfig{LogLevel: cfg.Logger.LogLevel})
	if err != nil {
		log.Println("logger error:", err)
		return
	}

	repositories := config.Repositories{}
	switch cfg.DB.Type {
	case "postgres":
		dbConn, err := config.NewPostgresDB(&config.PostgresConfig{
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
		repositories.Track = postgres.NewPostgresTrackRepository(dbConn)
	default:
		logger.Fatal("Error unknown database name", zap.Error(err),
			zap.String("Database name", cfg.DB.Type))
		return
	}

	redisClient, err := config.NewRedisClient(&config.RedisConfig{
		Host: cfg.Redis.Host,
		Port: cfg.Redis.Port,
	})

	if err != nil {
		logger.Fatal("Error connecting redis", zap.Error(err))
		return
	}

	minioClient, err := config.NewMinioClient(&config.MinioConfig{
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
	trackService := service.NewTrackService(trackRepo, trackStorage, genreService, logger)

	cons := consd.NewConsole(consd.NewHandler(consd.HandlerParams{
		AlbumService:    albumService,
		AuthService:     authService,
		CommentService:  commentService,
		GenreService:    genreService,
		MusicianService: musicianService,
		StatService:     statService,
		TrackService:    trackService,
		UserService:     userService,
	}))

	cons.Start()
}
