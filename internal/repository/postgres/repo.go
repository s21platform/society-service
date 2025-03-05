package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
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
	tx, err := r.connection.Beginx()
	if err != nil {
		return "", fmt.Errorf("failed to start transaction: %w", err)
	}

	var societyUUID string

	query, args, err := sq.Insert("society").
		Columns("name", "owner_uuid", "format_id", "post_permission_id", "is_search").
		Values(socData.Name, socData.OwnerUUID, socData.FormatID, socData.PostPermission, socData.IsSearch).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		_ = tx.Rollback()
		return "", fmt.Errorf("failed to build society insert query: %w", err)
	}

	err = tx.Get(&societyUUID, query, args...)
	if err != nil {
		_ = tx.Rollback()
		return "", fmt.Errorf("failed to insert society: %w", err)
	}

	query, args, err = sq.Insert("society_members").
		Columns("society_id", "user_uuid", "role", "payment_status").
		Values(societyUUID, socData.OwnerUUID, 1, 1).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		_ = tx.Rollback()
		return "", fmt.Errorf("failed to build society_members insert query: %w", err)
	}

	_, err = tx.Exec(query, args...)
	if err != nil {
		_ = tx.Rollback()
		return "", fmt.Errorf("failed to insert society member: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	return societyUUID, nil
}

func (r *Repository) GetSocietyInfo(societyUUID string) (*model.SocietyInfo, error) {
	var societyInfo model.SocietyInfo

	query := sq.Select(
		"s.name",
		"s.description",
		"s.owner_uuid",
		"s.photo_url",
		"s.format_id",
		"s.post_permission_id",
		"s.is_search",
	).
		From("society s").
		Where(sq.Eq{"s.id": societyUUID}).
		PlaceholderFormat(sq.Dollar)

	sqlString, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query: %w", err)
	}

	err = r.connection.GetContext(context.Background(), &societyInfo, sqlString, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get society info: %w", err)
	}

	if !societyInfo.Description.Valid {
		societyInfo.Description.String = ""
	}

	count, err := r.CountSubscribe(societyUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get count of subscribers: %w", err)
	}
	societyInfo.CountSubscribe = count

	tags, err := r.GetTags(societyUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}
	societyInfo.TagsID = tags

	return &societyInfo, nil
}

func (r *Repository) GetTags(societyUUID string) ([]int64, error) {
	query := sq.Select("tag_id").
		From("society_has_tags").
		Where(sq.Eq{"society_id": societyUUID})

	sqlString, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query GetTags: %w", err)
	}

	var tags []int64
	err = r.connection.Select(&tags, sqlString, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query GetTags: %w", err)
	}
	return tags, nil
}

func (r *Repository) CountSubscribe(societyUUID string) (int64, error) {
	query := sq.Select("count(*)").From("society_members").Where(sq.Eq{"society_id": societyUUID})
	sqlString, args, err := query.ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build SQL query CountSubscribe: %w", err)
	}
	var counts []int64
	err = r.connection.Select(&counts, sqlString, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to execute query CountSubscribe: %w", err)
	}
	return counts[0], nil
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
		PlaceholderFormat(sq.Dollar)

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
