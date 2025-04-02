//go:generate mockgen -destination=mock_contract_test.go -package=${GOPACKAGE} -source=contract.go
package service

import (
	"context"

	society "github.com/s21platform/society-proto/society-proto"
	"github.com/s21platform/society-service/internal/model"
)

type DbRepo interface {
	CreateSociety(ctx context.Context, socData *model.SocietyData) (string, error)
	GetSocietyInfo(ctx context.Context, societyUUID string) (*model.SocietyInfo, error)
	UpdateSociety(ctx context.Context, societyData *society.UpdateSocietyIn) error
	IsOwnerAdminModerator(ctx context.Context, peerUUID, societyUUID string) (int, error)
	GetTags(ctx context.Context, societyUUID string) ([]int64, error)
	CountSubscribe(ctx context.Context, societyUUID string) (int64, error)
	GetSocietyWithOffset(data *model.WithOffsetData) (*[]model.SocietyWithOffsetData, error)
	GetCountSocietyWithOffset(data *model.WithOffsetData) (int64, error)
}
