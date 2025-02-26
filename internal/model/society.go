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
