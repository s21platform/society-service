package model

type SocietyData struct {
	Name          string
	Description   string
	IsPrivate     bool
	DirectionId   int64
	AccessLevelId int64
}

type AccessLevel struct {
	Id          int64  `db:"id"`
	AccessLevel string `db:"level_name"`
}

type AccessLevelData struct {
	AccessLevel []AccessLevel
}
