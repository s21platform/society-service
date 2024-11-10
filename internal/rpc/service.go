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
		out.Levels[i].Id = data.AccessLevel[i].Id
		out.Levels[i].AccessLevel = data.AccessLevel[i].AccessLevel

	}
	return &out, err
}
