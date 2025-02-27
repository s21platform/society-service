package model

type SocietyData struct {
	Name           string
	FormatID       int64
	PostPermission int64
	IsSearch       bool
	OwnerUUID      string
}

type SocietyInfo struct {
	Name           string
	Description    string
	OwnerUUID      string
	PhotoURL       string
	FormatID       int64
	PostPermission int64
	IsSearch       bool
	CountSubscribe int64
	TagsID         []int64
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
	IsPrivate   bool   `db:"is_private"`
}

type WithOffsetData struct {
	Limit  int64
	Offset int64
	Name   string
	Uuid   string
}

//
//type AccessLevel struct {
//	Id          int64  `db:"id"`
//	AccessLevel string `db:"level_name"`
//}
//
//type AccessLevelData struct {
//	AccessLevel []AccessLevel
//}
//
//type GetPermissions struct {
//	Id          int64  `db:"id"`
//	Name        string `db:"name"`
//	Description string `db:"description"`
//}
//
//type SocietyWithOffset struct {
//	Society []SocietyWithOffsetData
//	Total   int64
//}
//
//type SocietyWithOffsetData struct {
//	Name       string `db:"name"`
//	AvatarLink string `db:"photo_url"`
//	SocietyId  int64  `db:"society_id"`
//	IsMember   bool   `db:"is_member"`
//	IsPrivate  bool   `db:"is_private"`
//}
//
//type WithOffsetData struct {
//	Limit  int64
//	Offset int64
//	Name   string
//	Uuid   string
//}
//
//type SocietyInfo struct {
//	Name             string `db:"name"`
//	Description      string `db:"description"`
//	OwnerId          string `db:"owner_uuid"`
//	PhotoUrl         string `db:"photo_url"`
//	IsPrivate        bool   `db:"is_private"`
//	CountSubscribers int64  `db:"count_subscribers"`
//}
//
//type UsersForSociety struct {
//	Name       string `db:"name"`
//	AvatarLink string `db:"avatar_link"`
//	Uuid       string `db:"uuid"`
//}
//
//type IdSociety struct {
//	ID int64 `db:"society_id"`
//}
//
//type SocietiesForUser struct {
//	Name        string `db:"name"`
//	Description string `db:"description"`
//}
