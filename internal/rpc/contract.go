package rpc

import (
	"github.com/s21platform/society-service/internal/model"
)

type DbRepo interface {
	CreateGroup(socData *model.SocietyData) (int, error)
	GetAccessLevel() (*model.AccessLevelData, error)
	GetPermissions() (*model.GetPermissionsData, error)
}
