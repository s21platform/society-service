package postgres

import (
	"fmt"
	"log"
	"time"

	society "github.com/s21platform/society-proto/society-proto"

	"github.com/lib/pq"

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
	err := r.connection.QueryRowx("INSERT INTO society(name, description, owner_uuid, format_id, post_permission_id, is_search)"+
		"VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		socData.Name,
		"",
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
	_, err = r.connection.Exec("INSERT INTO society_members(society_id, user_uuid, role, payment_status)"+
		"VALUES ($1, $2, 1, 1)", societyUUID, socData.OwnerUUID)

	if err != nil {
		return "", err
	}

	return societyUUID.String(), nil
}

func (r *Repository) GetSocietyInfo(societyUUID string) (*model.SocietyInfo, error) {
	var societyInfo model.SocietyInfo

	query := `SELECT 
		s.name, 
		s.description, 
		s.owner_uuid, 
		s.photo_url, 
		s.format_id, 
		s.post_permission_id, 
		s.is_search, 
		COALESCE(COUNT(mr.user_uuid), 0) AS count_subscribe, 
		ARRAY_REMOVE(ARRAY_AGG(sha.tag_id), NULL) AS tags_id
	FROM society s
	LEFT JOIN members_requests mr ON s.id = mr.society_id AND mr.status_id = 1 -- 1 означает, что пользователь принят в сообщество
	LEFT JOIN society_has_tags sha ON s.id = sha.society_id AND sha.is_active = TRUE
	WHERE s.id = $1
	GROUP BY s.id;`

	row := r.connection.QueryRow(query, societyUUID)

	var tags pq.Int64Array

	err := row.Scan(
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
		return nil, err
	}

	societyInfo.TagsID = tags

	return &societyInfo, nil
}

func (r *Repository) UpdateSociety(societyData *society.UpdateSocietyIn, peerUUID string) error {
	if !isOwnerAdminModerator(peerUUID, societyData.SocietyUUID, r) {
		return fmt.Errorf("faild to peer not Owner, Admin or Moderator to update society")
	}
	_, err := r.connection.Exec("update society set name = $1, description = $2,"+
		"photo_url = $3, format_id = $4, post_permission_id = $5, is_search = $6 where id = $7",
		societyData.Name,
		societyData.Description,
		societyData.PhotoURL,
		societyData.FormatID,
		societyData.PostPermission,
		societyData.IsSearch,
		societyData.SocietyUUID)
	if err != nil {
		return fmt.Errorf("faild to update society: %w", err)
	}
	return nil
}

func isOwnerAdminModerator(peerUUID, societyUUID string, r *Repository) bool {
	var role int
	err := r.connection.QueryRowx("select role from society_members where society_id = $1 and user_uuid = $2", societyUUID, peerUUID).Scan(&role)
	if err != nil {
		return false
	}
	if role == 1 || role == 2 || role == 3 {
		return true
	}
	return false
}
