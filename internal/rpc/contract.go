package rpc

import (
	society "github.com/s21platform/society-proto/society-proto"
	"github.com/s21platform/society-service/internal/model"
)

type DbRepo interface {
	CreateGroup(socData *model.SocietyData) (int, error)
	GetAccessLevel() (*society.GetAccessLevelOut, error)
}
