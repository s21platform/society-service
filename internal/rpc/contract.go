package rpc

import (
	"github.com/s21platform/society-service/internal/model"
)

type DbRepo interface {
	CreateGroup(socData *model.SocietyData) (int, error)
	GetAccessLevel() (*[]model.AccessLevel, error)
	GetPermissions() (*[]model.GetPermissions, error)
	GetSocietyWithOffset(data *model.WithOffsetData) (*[]model.SocietyWithOffsetData, error)
}
