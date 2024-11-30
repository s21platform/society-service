package model

type SocietyData struct {
	Name          string
	Description   string
	IsPrivate     bool
	DirectionId   int64
	OwnerId       string
	AccessLevelId int64
}

type AccessLevel struct {
	Id          int64  `db:"id"`
	AccessLevel string `db:"level_name"`
}

type AccessLevelData struct {
	AccessLevel []AccessLevel
}

type GetPermissions struct {
	Id          int64  `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
}

type SocietyWithOffset struct {
	Society []SocietyWithOffsetData
	Total   int64
}

type SocietyWithOffsetData struct {
	Name       string
	AvatarLink string
	SocietyId  int64
	IsMember   bool
}

type WithOffsetData struct {
	Limit  int64
	Offset int64
	Name   string
	Uuid   string
}
