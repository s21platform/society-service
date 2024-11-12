package rpc

import (
	"context"
	"fmt"

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
	SocietyData := model.SocietyData{
		Name:          in.Name,
		Description:   in.Description,
		IsPrivate:     in.IsPrivate,
		DirectionId:   in.DirectionId,
		AccessLevelId: in.AccessLevelId,
	}
	id, err := s.dbR.CreateGroup(&SocietyData)

	if err != nil {
		return nil, err
	}
	out := &society.SetSocietyOut{SocietyId: int64(id)}
	return out, err
}

func (s *Server) GetAccessLevel(context.Context, *society.Empty) (*society.GetAccessLevelOut, error) {
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
		return nil, fmt.Errorf("s.dbR.GetPermissions %v", err)
	}

	out := society.GetPermissionsOut{
		Permissions: make([]*society.GetPermissions, len(data.GetPermissions)),
	}

	for i := range data.GetPermissions {
		level := &society.GetPermissions{
			id:          data.GetPermissions[i].Id,
			name:        data.GetPermissions[i].Name,
			description: data.GetPermissions[i].Description,
		}
		out.Permissions[i] = level
	}

	return &out, err
}
