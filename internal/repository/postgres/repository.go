package postgres

import (
	"fmt"
	"github.com/s21platform/society-service/internal/model"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Импортируем данную библиотеку для работы с бд.
	"github.com/s21platform/society-service/internal/config"
)

type Repository struct {
	connection *sqlx.DB
}

func connect(cfg *config.Config) (*Repository, error) {
	// Connect db
	conStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Database, cfg.Postgres.Host, cfg.Postgres.Port)

	db, err := sqlx.Connect("postgres", conStr)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
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

func (r *Repository) CreateGroup(socData *model.SocietyData) (int, error) {
	tx, err := r.connection.Beginx()
	if err != nil {
		return 0, err
	}
	var lastId int
	err = tx.QueryRowx("INSERT INTO societies(name, description, is_private, direction_id, access_level) VALUES ($1,$2,$3,$4,$5) RETURNING id", socData.Name, socData.Description, socData.IsPrivate, socData.DirectionId, socData.AccessLevelId).Scan(&lastId)
	if err != nil {
		return 0, err
	}
	return lastId, nil
}
