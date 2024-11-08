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
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

var (
	// Counter for total HTTP requests
	requestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "echo_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method"},
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "echo_http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)
)

func init() {
    // Register metrics with Prometheus
    prometheus.MustRegister(requestCount, requestDuration)
}


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
	router.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c) // Call the next handler
			duration := time.Since(start).Seconds()

			// Increment the request counter
			requestCount.WithLabelValues(c.Request().Method).Inc()

			// Record the request duration
			requestDuration.WithLabelValues(c.Request().Method).Observe(duration)

			return err
		}
	})
	services := &echoapi.Services{
		AuthService:     nil,
		AlbumService:    nil,
		MusicianService: nil,
		UserService:     userService,
		TrackService:    nil,
		CommentService:  nil,
		GenreService:    nil,
	}
	router.Use(middleware.Logger())
	_ = echoapi.NewUserHandler(
		router,
		logger,
		services,
	)

	router.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	server := http.Server{
		Handler:      router,
		Addr:         ":8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
	}

	if err = server.ListenAndServe(); err != nil {
		log.Fatalf("error while listening: %v", err)
	}
}
