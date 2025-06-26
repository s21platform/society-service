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
	CanEditSociety bool           `db:"-"`
}

type SocietyWithOffset struct {
	Society []SocietyWithOffsetData
	Total   int64
}

type SocietyWithOffsetData struct {
	SocietyUUID string `db:"id"`
	Name        string `db:"name"`
	PhotoURL    string `db:"photo_url"`
	IsMember    bool   `db:"is_member"`
	FormatId    int64  `db:"format_id"`
}

type WithOffsetData struct {
	Limit  int64
	Offset int64
	Name   string
	Uuid   string
}

type Role struct {
	Role int `db:"role"`
}
