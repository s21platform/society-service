package postgres

import (
	"fmt"
	"log"
	"time"

	"github.com/lib/pq"

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

//	var data model.SocietyInfo
//	query := "SELECT name, " +
//		"description, " +
//		"owner_uuid, " +
//		"photo_url, " +
//		"is_private, " +
//		"COALESCE(count_s, 0) " +
//		"AS count_subscribers FROM societies s LEFT JOIN (SELECT society_id, count(*) AS count_s from societies_subscribers GROUP BY society_id) ss ON s.id = ss.society_id " +
//		"WHERE id = $1"
//	err := r.connection.Get(&data, query, id)
//	if err != nil {
//		return nil, fmt.Errorf("failed to get society info: %v", err)
//	}
//
//	return &data, nil
//}
//
//func (r *Repository) GetAccessLevel() (*[]model.AccessLevel, error) {
//	var data []model.AccessLevel
//	err := r.connection.Select(&data, "SELECT id, level_name FROM access_level")
//	if err != nil {
//		return nil, fmt.Errorf("r.connection.Select: %v", err)
//	}
//
//	return &data, nil
//}
//
//func (r *Repository) GetPermissions() (*[]model.GetPermissions, error) {
//	var data []model.GetPermissions
//	err := r.connection.Select(&data, "SELECT id, name, description FROM user_permissions")
//	if err != nil {
//		return nil, fmt.Errorf("failed to get permission: %v", err)
//	}
//
//	return &data, nil
//}
//
//func (r *Repository) GetSocietyWithOffset(socData *model.WithOffsetData) (*[]model.SocietyWithOffsetData, error) {
//	var data []model.SocietyWithOffsetData
//	query := "SELECT name, photo_url, s.id society_id, " +
//		"CASE WHEN ss.user_uuid = $1 THEN true ELSE false END AS is_member " +
//		"FROM societies s " +
//		"LEFT JOIN societies_subscribers ss ON s.id = ss.society_id AND ss.user_uuid = $1 " +
//		"WHERE ($2 = '' OR name ILIKE $2) "
//
//	err := r.connection.Select(&data, query+"OFFSET $3 LIMIT $4", socData.Uuid, "%"+socData.Name+"%", socData.Offset, socData.Limit)
//	if err != nil {
//		return nil, fmt.Errorf("failed to get society with offset: %v", err)
//	}
//	return &data, err
//}
//
//func (r *Repository) GetCountSocietyWithOffset(socData *model.WithOffsetData) (int64, error) {
//	var count int64
//	query := "SELECT name, photo_url, s.id society_id, " +
//		"CASE WHEN ss.user_uuid = $1 THEN true ELSE false END AS is_member " +
//		"FROM societies s " +
//		"LEFT JOIN societies_subscribers ss ON s.id = ss.society_id AND ss.user_uuid = $1 " +
//		"WHERE ($2 = '' OR name ILIKE $2) "
//
//	queryCount := "with test as  (" + query + ")" +
//		"SELECT count(*) FROM test"
//
//	err := r.connection.Get(&count, queryCount, socData.Uuid, "%"+socData.Name+"%")
//	if err != nil {
//		return 0, fmt.Errorf("failed to get count society with offset: %v", err)
//	}
//	return count, nil
//}
//
//
//func (r *Repository) SubscribeToSociety(id int64, uuid string) (bool, error) {
//	_, err := r.connection.Exec("INSERT INTO societies_subscribers (society_id, user_uuid) VALUES ($1, $2)", id, uuid)
//	if err != nil {
//		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
//			return false, nil
//		}
//		return false, fmt.Errorf("r.connection.Exec: %v", err)
//	}
//
//	return true, nil
//}
//
//func (r *Repository) UnsubscribeFromSociety(id int64, uuid string) (bool, error) {
//	var data []model.IdSociety
//	err := r.connection.Select(&data, "SELECT society_id FROM societies_subscribers WHERE society_id = $1 AND user_uuid = $2", id, uuid)
//	if err != nil {
//		return false, fmt.Errorf("failed to select subscribe from society: %v", err)
//	}
//	if len(data) == 0 {
//		return false, nil
//	}
//	if _, err := r.connection.Exec("DELETE FROM societies_subscribers WHERE society_id = $1 AND user_uuid = $2", id, uuid); err != nil {
//		return false, fmt.Errorf("r.connection.Exec: %v", err)
//	}
//
//	return true, nil
//}
//
//func (r *Repository) GetSocietiesForUser(uuid string, uuidUser string) (*[]model.SocietyWithOffsetData, error) {
//	var data []model.SocietyWithOffsetData
//	err := r.connection.Select(&data, "SELECT name, photo_url, ss.id AS society_id, user_uuid = $1 as is_member, is_private FROM societies s JOIN societies_subscribers ss ON s.id = ss.society_id WHERE user_uuid = $2", uuid, uuidUser)
//	if err != nil {
//		return nil, fmt.Errorf("failed to select societies for user: %v", err)
//	}
//
//	return &data, nil
//}
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
