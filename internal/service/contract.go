//go:generate mockgen -destination=mock_contract_test.go -package=${GOPACKAGE} -source=contract.go
package service

import (
	"github.com/s21platform/society-service/internal/model"
)

type DbRepo interface {
	CreateSociety(socData *model.SocietyData) (string, error)
	GetSocietyInfo(societyUUID string) (*model.SocietyInfo, error)
}
