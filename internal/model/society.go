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
	Name       string `db:"name"`
	AvatarLink string `db:"photo_url"`
	SocietyId  int64  `db:"society_id"`
	IsMember   bool   `db:"is_member"`
	IsPrivate  bool   `db:"is_private"`
}

type WithOffsetData struct {
	Limit  int64
	Offset int64
	Name   string
	Uuid   string
}

type SocietyInfo struct {
	Name        string `db:"name"`
	Description string `db:"description"`
	OwnerId     string `db:"owner_uuid"`
	PhotoUrl    string `db:"photo_url"`
	IsPrivate   bool   `db:"is_private"`
}

type UsersForSociety struct {
	Name       string `db:"name"`
	AvatarLink string `db:"avatar_link"`
	Uuid       string `db:"uuid"`
}

type IdSociety struct {
	ID int64 `db:"society_id"`
}

type SocietiesForUser struct {
	Name        string `db:"name"`
	Description string `db:"description"`
}
