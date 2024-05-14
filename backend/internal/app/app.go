package app

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/adapters/repository"
	"github.com/hanoys/sigma-music/internal/domain"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

type Config struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
}

const (
	maxConn         = 100
	maxConnIdleTime = 1 * time.Minute
	maxConnLifetime = 3 * time.Minute
)

func NewPostgresDB(cfg *Config) (*sqlx.DB, error) {
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

func Run() {
	conn, err := NewPostgresDB(&Config{
		Host:     "localhost",
		Port:     "5432",
		Database: "sigmamusic",
		User:     "sigma",
		Password: "sigma",
	})

	if err != nil {
		log.Fatal(err)
	}
	r := repository.NewPostgresUserRepository(conn)

	var newUser = domain.User{
		ID:       uuid.New(),
		Name:     "CreatedUser",
		Email:    "CreatedUser@mail.com",
		Phone:    "+71111111111",
		Password: "CreatedUserPassword",
		Country:  "CreatedUserCountry",
	}

	created, err := r.GetByName(context.Background(), newUser.Name)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Created user: %v\n", created)

}
