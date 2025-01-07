package rpc

import (
	"github.com/s21platform/society-service/internal/model"
)

type DbRepo interface {
	CreateGroup(socData *model.SocietyData) (int, error)
	GetAccessLevel() (*[]model.AccessLevel, error)
	GetPermissions() (*[]model.GetPermissions, error)
	GetSocietyWithOffset(data *model.WithOffsetData) (*[]model.SocietyWithOffsetData, error)
	GetSocietyInfo(id int64) (*model.SocietyInfo, error)
	SubscribeToSociety(id int64, uuid string) (bool, error)
	UnsubscribeFromSociety(id int64, uuid string) (bool, error)
	GetSocietiesForUser(uuid string, uuidUser string) (*[]model.SocietyWithOffsetData, error)
}
