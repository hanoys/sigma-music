package gin

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	ginapi "github.com/hanoys/sigma-music/internal/adapters/delivery/api/gin"
	"github.com/hanoys/sigma-music/internal/adapters/hash"
	"github.com/hanoys/sigma-music/internal/adapters/repository/postgres"
	"github.com/hanoys/sigma-music/internal/app/config"
	"github.com/hanoys/sigma-music/internal/service"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

var (
        // Counter for total HTTP requests
    requestCount = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "gin_http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method"},
    )

     requestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "gin_http_request_duration_seconds",
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

// Middleware to collect metrics
func prometheusMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Start time for request duration
        startTime := time.Now()
        
        // Process request
        c.Next()

        // Calculate metrics
        // path := c.FullPath()
        method := c.Request.Method
        duration := time.Since(startTime).Seconds()

        // Update Prometheus metrics
        requestCount.WithLabelValues(method).Inc()
        requestDuration.WithLabelValues(method).Observe(duration)
    }
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
	router := gin.Default()
	router.Use(prometheusMiddleware())
	services := &ginapi.Services{
		AuthService:     nil,
		AlbumService:    nil,
		MusicianService: nil,
		UserService:     userService,
		TrackService:    nil,
		CommentService:  nil,
		GenreService:    nil,
	}
	apiRouter := router.Group("")
	_ = ginapi.NewUserHandler(
		apiRouter,
		logger,
		services,
	)

	apiRouter.GET("/metrics", gin.WrapH(promhttp.Handler()))

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
