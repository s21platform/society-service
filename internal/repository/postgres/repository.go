package postgres

import (
	"fmt"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/google/uuid"

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

func (r *Repository) CreateSociety(socData *model.SocietyData) (string, error) {
	var societyUUIDStr string
	err := r.connection.QueryRowx("INSERT INTO society(name, owner_uuid, format_id, post_permission_id, is_search)"+
		"VALUES ($1, $2, $3, $4, $5) RETURNING id",
		socData.Name,
		socData.OwnerUUID,
		socData.FormatID,
		socData.PostPermission,
		socData.IsSearch).Scan(&societyUUIDStr)
	if err != nil {
		return "", err
	}

	societyUUID, err := uuid.Parse(societyUUIDStr)
	if err != nil {
		return "", err
	}

	return societyUUID.String(), nil
}

func (r *Repository) GetSocietyWithOffset(data *model.WithOffsetData) (*[]model.SocietyWithOffsetData, error) {
	var out []model.SocietyWithOffsetData

	baseQuery, args, err := sq.Select(
		"name",
		"photo_url",
		"s.id AS society_id",
		"CASE WHEN ss.user_uuid = ? THEN true ELSE false END AS is_member",
	).
		From("societies s").
		LeftJoin("societies_subscribers ss ON s.id = ss.society_id AND ss.user_uuid = ?", data.Uuid).
		Where(sq.Or{
			sq.Expr("? = ''", data.Name),
			sq.Expr("name ILIKE ?", "%"+data.Name+"%"),
		}).
		Offset(uint64(data.Offset)).
		Limit(uint64(data.Limit)).PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	err = r.connection.Select(&out, baseQuery, args...)
	if err != nil {
		return nil, err
	}

	return &out, err
}

func (r *Repository) GetCountSocietyWithOffset(socData *model.WithOffsetData) (int64, error) {
	var count int64

	baseQuery := sq.Select(
		"name",
		"photo_url",
		"s.id AS society_id",
		"CASE WHEN ss.user_uuid = ? THEN true ELSE false END AS is_member",
	).
		From("societies s").
		LeftJoin("societies_subscribers ss ON s.id = ss.society_id AND ss.user_uuid = ?", socData.Uuid).
		Where(sq.Or{
			sq.Expr("? = ''", socData.Name),
			sq.Expr("name ILIKE ?", "%"+socData.Name+"%"),
		}).
		PlaceholderFormat(sq.Dollar)

	countQuery, args, err := sq.Select("COUNT(*)").
		FromSelect(baseQuery, "test").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return 0, err
	}

	err = r.connection.Get(&count, countQuery, args...)
	if err != nil {
		return 0, err
	}
	return count, nil
}
