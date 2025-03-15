package model

import "database/sql"

type SocietyData struct {
	Name           string
	FormatID       int64
	PostPermission int64
	IsSearch       bool
	OwnerUUID      string
}

type SocietyInfo struct {
	Name           string         `db:"name"`
	Description    sql.NullString `db:"description"`
	OwnerUUID      string         `db:"owner_uuid"`
	PhotoURL       string         `db:"photo_url"`
	FormatID       int64          `db:"format_id"`
	PostPermission int64          `db:"post_permission_id"`
	IsSearch       bool           `db:"is_search"`
	CountSubscribe int64          `db:"-"`
	TagsID         []int64        `db:"-"`
}

type Role []struct {
	Role int `db:"role"`
}
