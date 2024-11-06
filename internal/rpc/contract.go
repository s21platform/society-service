package rpc

import (
	"github.com/s21platform/society-service/internal/model"
)

type DbRepo interface {
	CreateGroup(socData *model.SocietyData) (int, error)
	GetAccessLevel(data *model.AccessLevelData) error
}
