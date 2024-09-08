package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/s21platform/society-service/internal/config"
	"log"
	"time"
)

type Repository struct {
	connection *sql.DB
}

func connect(cfg *config.Config) (*Repository, error) {
	// Connect db
	conStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Database, cfg.Postgres.Host, cfg.Postgres.Port)

	db, err := sql.Open("postgres", conStr)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	// Сhecking connection db
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("db.Ping: %w", err)
	}

	return &Repository{db}, nil
}

func New(cfg *config.Config) (*Repository, error) {
	var err error

	var repo *Repository

	for i := 0; i < 5; i++ {
		repo, err = connect(cfg)
		if err == nil {
			return repo, nil
		}

		log.Println(err)
		time.Sleep(500 * time.Millisecond)
	}

	return nil, err
}

func (r *Repository) Close() {
	r.connection.Close()
}
