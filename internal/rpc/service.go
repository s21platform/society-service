package rpc

import (
	"context"
	"fmt"

	"github.com/s21platform/society-service/internal/config"

	society "github.com/s21platform/society-proto/society-proto"
	"github.com/s21platform/society-service/internal/model"
)

type Server struct {
	society.UnimplementedSocietyServiceServer
	dbR DbRepo
}

func New(repo DbRepo) *Server {
	return &Server{
		dbR: repo,
	}
}
func (s *Server) CreateSociety(ctx context.Context, in *society.SetSocietyIn) (*society.SetSocietyOut, error) {
	uuid, ok := ctx.Value(config.KeyUUID).(string)

	if !ok {
		return nil, fmt.Errorf("uuid not found in metadata")
	}

	SocietyData := model.SocietyData{
		Name:          in.Name,
		Description:   in.Description,
		IsPrivate:     in.IsPrivate,
		DirectionId:   in.DirectionId,
		OwnerId:       uuid,
		AccessLevelId: in.AccessLevelId,
	}
	id, err := s.dbR.CreateGroup(&SocietyData)

	if err != nil {
		return nil, err
	}
	out := &society.SetSocietyOut{SocietyId: int64(id)}
	return out, err
}

func (s *Server) GetAccessLevel(context.Context, *society.EmptySociety) (*society.GetAccessLevelOut, error) {
	data, err := s.dbR.GetAccessLevel()
	if err != nil {
		return nil, fmt.Errorf("s.dbR.GetAccessLevel %v", err)
	}

	out := society.GetAccessLevelOut{
		Levels: make([]*society.AccessLevel, len(data.AccessLevel)),
	}
	for i := range data.AccessLevel {
		level := &society.AccessLevel{
			Id:          data.AccessLevel[i].Id,
			AccessLevel: data.AccessLevel[i].AccessLevel,
		}
		out.Levels[i] = level
	}

	return &out, err
}

func (s *Server) GetPermissions(context.Context, *society.EmptySociety) (*society.GetPermissionsOut, error) {
	data, err := s.dbR.GetPermissions()
	if err != nil {
		return nil, fmt.Errorf("failed to get permission: %v", err)
	}

	out := society.GetPermissionsOut{
		Permissions: make([]*society.Permission, len(*data)),
	}

	for a, i := range *data {
		level := &society.Permission{
			Id:          i.Id,
			Name:        i.Name,
			Description: i.Description,
		}
		out.Permissions[a] = level
	}

	return &out, err
}
