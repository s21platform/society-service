package postgres

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	sq "github.com/Masterminds/squirrel"
	society "github.com/s21platform/society-proto/society-proto"
	"github.com/s21platform/society-service/internal/config"
	"github.com/s21platform/society-service/internal/model"
)

type Repository struct {
	connection *sqlx.DB
}

func connect(cfg *config.Config) (*Repository, error) {
	conStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Database, cfg.Postgres.Host, cfg.Postgres.Port)

	db, err := sqlx.Connect("postgres", conStr)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	return &Repository{
		connection: db,
	}, nil
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
	var societyUUID string

	query := sq.Insert("society").
		Columns("name", "owner_uuid", "format_id", "post_permission_id", "is_search").
		Values(socData.Name, socData.OwnerUUID, socData.FormatID, socData.PostPermission, socData.IsSearch).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		RunWith(r.connection)

	err := query.QueryRow().Scan(&societyUUID)
	if err != nil {
		return "", fmt.Errorf("failed to insert society: %v", err)
	}

	query = sq.Insert("society_members").
		Columns("society_id", "user_uuid", "role", "payment_status").
		Values(societyUUID, socData.OwnerUUID, 1, 1).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.connection)

	_, err = query.Exec()
	if err != nil {
		return "", fmt.Errorf("failed to insert society member: %v", err)
	}

	return societyUUID, nil
}

func (r *Repository) GetSocietyInfo(societyUUID string) (*model.SocietyInfo, error) {
	var societyInfo model.SocietyInfo
	var tags pq.Int64Array

	query := sq.Select(
		"s.name",
		"s.description",
		"s.owner_uuid",
		"s.photo_url",
		"s.format_id",
		"s.post_permission_id",
		"s.is_search",
		"COALESCE(COUNT(mr.user_uuid), 0) AS count_subscribe",
		"ARRAY_REMOVE(ARRAY_AGG(sha.tag_id), NULL) AS tags_id",
	).
		From("society s").
		LeftJoin("members_requests mr ON s.id = mr.society_id AND mr.status_id = 1").
		LeftJoin("society_has_tags sha ON s.id = sha.society_id AND sha.is_active = TRUE").
		Where(sq.Eq{"s.id": societyUUID}).
		GroupBy(
			"s.id",
			"s.name",
			"s.description",
			"s.owner_uuid",
			"s.photo_url",
			"s.format_id",
			"s.post_permission_id",
			"s.is_search",
		).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query: %w", err)
	}

	row := r.connection.QueryRow(sql, args...)
	err = row.Scan(
		&societyInfo.Name,
		&societyInfo.Description,
		&societyInfo.OwnerUUID,
		&societyInfo.PhotoURL,
		&societyInfo.FormatID,
		&societyInfo.PostPermission,
		&societyInfo.IsSearch,
		&societyInfo.CountSubscribe,
		&tags,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan society info: %w", err)
	}

	societyInfo.TagsID = tags
	return &societyInfo, nil
}

func (r *Repository) UpdateSociety(societyData *society.UpdateSocietyIn, peerUUID string) error {
	isAllowed, err := isOwnerAdminModerator(peerUUID, societyData.SocietyUUID, r)
	if err != nil {
		return fmt.Errorf("failed to check user permissions: %w", err)
	}
	if !isAllowed {
		return fmt.Errorf("failed to user is not Owner, Admin or Moderator")
	}

	query := sq.Update("society").
		Set("name", societyData.Name).
		Set("description", societyData.Description).
		Set("photo_url", societyData.PhotoURL).
		Set("format_id", societyData.FormatID).
		Set("post_permission_id", societyData.PostPermission).
		Set("is_search", societyData.IsSearch).
		Where(sq.Eq{"id": societyData.SocietyUUID}).
		PlaceholderFormat(sq.Dollar) // Гарантируем использование $1, $2...

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build SQL query: %w", err)
	}

	_, err = r.connection.Exec(sql, args...)
	if err != nil {
		return fmt.Errorf("failed to update society: %w", err)
	}

	return nil
}

func isOwnerAdminModerator(peerUUID, societyUUID string, r *Repository) (bool, error) {
	var role int

	query := sq.Select("role").
		From("society_members").
		Where(sq.Eq{"society_id": societyUUID, "user_uuid": peerUUID}).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return false, fmt.Errorf("failed to build SQL query: %w", err)
	}

	err = r.connection.QueryRow(sql, args...).Scan(&role)
	if err != nil {
		return false, fmt.Errorf("failed to fetch user role: %w", err)
	}

	return role == 1 || role == 2 || role == 3, nil
}
