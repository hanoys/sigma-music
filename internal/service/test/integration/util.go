package integrationtest

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/golang-migrate/migrate"
	migratepg "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	testpg "github.com/testcontainers/testcontainers-go/modules/postgres"
)

var (
	DatabaseName = "sigma-music"
	UserName     = "sigma-music"
	Password     = "sigma-music"
)

func newPostgresContainer(ctx context.Context) (*testpg.PostgresContainer, error) {
	container, err := testpg.Run(
		ctx,
		"docker.io/postgres:16-alpine",
		testpg.WithDatabase(DatabaseName),
		testpg.WithUsername(UserName),
		testpg.WithPassword(Password),
		testpg.BasicWaitStrategies(),
		testpg.WithSQLDriver("pgx"),
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

func isPreviousTestsFailed() bool {
	if os.Getenv("UNIT_SUCCESS") == "1" {
		return false
	}

	return true
}
