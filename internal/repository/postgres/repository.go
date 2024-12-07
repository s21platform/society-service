package postgres

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/s21platform/society-service/internal/model"

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
	var lastId int
	err := r.connection.QueryRowx("INSERT INTO societies(name, description, is_private, direction_id, owner_uuid, photo_url, access_id) VALUES ($1,$2,$3,$4,$5, $6, $7) RETURNING id",
		socData.Name,
		socData.Description,
		socData.IsPrivate,
		socData.DirectionId,
		socData.OwnerId,
		"https://storage.yandexcloud.net/space21/avatars/default/logo-discord.jpeg",
		socData.AccessLevelId).Scan(&lastId)
	if err != nil {
		return 0, err
	}
	return lastId, nil
}

func (r *Repository) GetAccessLevel() (*[]model.AccessLevel, error) {
	var data []model.AccessLevel
	err := r.connection.Select(&data, "SELECT id, level_name FROM access_level")
	if err != nil {
		return nil, fmt.Errorf("r.connection.Select: %v", err)
	}

	return &data, nil
}

func (r *Repository) GetPermissions() (*[]model.GetPermissions, error) {
	var data []model.GetPermissions
	err := r.connection.Select(&data, "SELECT id, name, description FROM user_permissions")
	if err != nil {
		return nil, fmt.Errorf("failed to get permission: %v", err)
	}

	return &data, nil
}

func (r *Repository) GetSocietyWithOffset(socData *model.WithOffsetData) (*[]model.SocietyWithOffsetData, error) {
	var data []model.SocietyWithOffsetData
	query := "SELECT name, photo_url avatar_link, s.id society_id, " +
		"CASE WHEN ss.user_uuid = $1 THEN true ELSE false END AS is_member " +
		"FROM societies s " +
		"LEFT JOIN societies_subscribers ss ON s.id = ss.society_id AND ss.user_uuid = $1 " +
		"WHERE ($2 = '' OR name ILIKE $2) " +
		"OFFSET $3 LIMIT $4"

	err := r.connection.Select(&data, query, socData.Uuid, "%"+socData.Name+"%", socData.Offset, socData.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get permission: %v", err)
	}

	return &data, err
}

func (r *Repository) GetSocietyInfo(id int64) (*model.SocietyInfo, error) {
	var data model.SocietyInfo
	query := "SELECT name, description, owner_uuid, photo_url, is_private FROM societies WHERE id = $1"
	err := r.connection.Get(&data, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get society info: %v", err)
	}

	return &data, nil
}

func (r *Repository) SubscribeToSociety(id int64, uuid string) (bool, error) {
	_, err := r.connection.Exec("INSERT INTO societies_subscribers (society_id, user_uuid) VALUES ($1, $2)", id, uuid)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return false, errors.New("user already subscribed to this society")
		}
		return false, fmt.Errorf("r.connection.Exec: %v", err)
	}

	return true, nil
}

func (r *Repository) UnsubscribeFromSociety(id int64, uuid string) (bool, error) {
	if _, err := r.connection.Exec("DELETE FROM societies_subscribers WHERE society_id = $1 AND user_uuid = $2", id, uuid); err != nil {
		return false, fmt.Errorf("r.connection.Exec: %v", err)
	}

	return true, nil
}
