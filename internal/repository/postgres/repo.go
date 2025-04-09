package postgres

import (
	"context"
	"database/sql"
	"errors"
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

func (r *Repository) CreateSociety(ctx context.Context, socData *model.SocietyData) (string, error) {
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

	err = tx.GetContext(ctx, &societyUUID, query, args...)
	if err != nil {
		_ = tx.Rollback()
		return "", fmt.Errorf("failed to insert society: %w", err)
	}

	query, args, err = sq.Insert("society_members").
		Columns("society_id", "user_uuid", "role", "payment_status").
		Values(societyUUID, socData.OwnerUUID, "1", "1").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		_ = tx.Rollback()
		return "", fmt.Errorf("failed to build society_members insert query: %w", err)
	}

	_, err = tx.ExecContext(ctx, query, args...)
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

func (r *Repository) GetSocietyInfo(ctx context.Context, societyUUID string) (*model.SocietyInfo, error) {
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

	err = r.connection.GetContext(ctx, &societyInfo, sqlString, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get society info: %w", err)
	}
	return &societyInfo, nil
}

func (r *Repository) GetTags(ctx context.Context, societyUUID string) ([]int64, error) {
	query := sq.Select("tag_id").From("society_has_tags").Where(sq.Eq{"society_id": societyUUID})
	sqlString, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query GetTags: %w", err)
	}

	var tags []int64
	err = r.connection.SelectContext(ctx, &tags, sqlString, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query GetTags: %w", err)
	}
	return tags, nil
}

func (r *Repository) CountSubscribe(ctx context.Context, societyUUID string) (int64, error) {
	query := sq.Select("count(*)").From("society_members").Where(sq.Eq{"society_id": societyUUID})
	sqlString, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build SQL query CountSubscribe: %w", err)
	}

	var count int64
	err = r.connection.GetContext(ctx, &count, sqlString, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to execute query CountSubscribe: %w", err)
	}

	return count, nil
}

func (r *Repository) UpdateSociety(ctx context.Context, societyData *society.UpdateSocietyIn) error {
	query := sq.Update("society").
		Set("name", societyData.Name).
		Set("description", societyData.Description).
		Set("format_id", societyData.FormatID).
		Set("post_permission_id", societyData.PostPermission).
		Set("is_search", societyData.IsSearch).
		Where(sq.Eq{"id": societyData.SocietyUUID}).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build SQL query: %w", err)
	}

	_, err = r.connection.ExecContext(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to update society: %w", err)
	}

	return nil
}

func (r *Repository) IsOwnerAdminModerator(ctx context.Context, peerUUID, societyUUID string) (int, error) {
	query, args, err := sq.Select("role").
		From("society_members").
		Where(sq.Eq{"society_id": societyUUID, "user_uuid": peerUUID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build SQL query: %w", err)
	}

	var result model.Role

	err = sqlx.GetContext(ctx, r.connection, &result, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, fmt.Errorf("failed to execute query isOwnerAdminModerator: %w", err)
	}

	return result.Role, nil
}
