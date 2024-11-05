package echo

import (
	"log"
	"net/http"
	"time"

	echoapi "github.com/hanoys/sigma-music/internal/adapters/delivery/api/echo"
	"github.com/hanoys/sigma-music/internal/adapters/hash"
	"github.com/hanoys/sigma-music/internal/adapters/repository/postgres"
	"github.com/hanoys/sigma-music/internal/app/config"
	"github.com/hanoys/sigma-music/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	userRepo := postgres.NewPostgresUserRepository(dbConn)
	hashProvider := hash.NewHashPasswordProvider()
	userService := service.NewUserService(userRepo, hashProvider, logger)
	router := echo.New()
	services := &echoapi.Services{
		AuthService:     nil,
		AlbumService:    nil,
		MusicianService: nil,
		UserService:     userService,
		TrackService:    nil,
		CommentService:  nil,
		GenreService:    nil,
	}
	apiRouter := router.Group("api/v1")
	_ = echoapi.NewUserHandler(
		apiRouter,
		logger,
		services,
	)

	apiRouter.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	server := http.Server{
		Handler:      router,
		Addr:         ":8082",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
	}

	if err = server.ListenAndServe(); err != nil {
		log.Fatalf("error while listening: %v", err)
	}
}
